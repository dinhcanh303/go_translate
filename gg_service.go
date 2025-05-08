package go_translate

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dinhcanh303/go_translate/utils"
)

type GoogleTranslateService struct {
	client *http.Client
	opts   *TranslateOptions
}

func NewGoogleTranslateService(opts *TranslateOptions) *GoogleTranslateService {
	return &GoogleTranslateService{
		client: &http.Client{Timeout: 5 * time.Second},
		opts:   opts,
	}
}

func (s *GoogleTranslateService) TranslateText(text, target string, detectedLangCode ...string) (string, error) {
	for t, endpoint := range GoogleUrls {
		var translated string
		var err error
		switch t {
		case TypeHtml:
			translated, err = s.callTranslateHTML(text, target, endpoint)
		case TypeClientGtx, TypeClientDictChromeEx:
			translated, err = s.callTranslateGet(text, target, endpoint, t == TypeClientGtx)
		case TypePaGtx:
			translated, err = s.callTranslatePa(text, target, endpoint)
		}

		if err == nil && translated != "" {
			return translated, nil
		}
		log.Printf("[ERROR] API %s failed: %v", t, err)
	}
	return "", errors.New("unable to translate text, all APIs failed")
}

func (s *GoogleTranslateService) callTranslateHTML(text, target, endpoint string) (string, error) {
	body := fmt.Sprintf(`[[["%s"],"auto","%s"],"wt_lib"]`, text, target)
	headers := map[string]string{
		"User-Agent":     utils.GetConditionalRandomValue(DefaultUserAgents, s.opts.CustomUserAgents, s.opts.UseRandomUserAgents),
		"Content-Type":   "application/json+protobuf",
		"X-Goog-API-Key": GOOGLE_API_KEY_TRANSLATE,
	}
	respBytes, err := s.doRequest("POST", endpoint, headers, nil, []byte(body))
	if err != nil {
		return "", err
	}
	return utils.ExtractTranslatedText(respBytes)
}

func (s *GoogleTranslateService) callTranslateGet(text, target, endpoint string, isGtx bool) (string, error) {
	fullURL := "https://" + utils.GetConditionalRandomValue(DefaultServiceUrls, s.opts.CustomServiceUrls, s.opts.UseRandomServiceUrls) + endpoint
	params := url.Values{
		"tl": {target},
		"q":  {text},
	}
	if s.opts.AddToken {
		params.Set("tk", utils.GgTokenGenerate(text))
	}
	headers := map[string]string{
		"User-Agent": utils.GetConditionalRandomValue(DefaultUserAgents, s.opts.CustomUserAgents, s.opts.UseRandomUserAgents),
	}
	respBytes, err := s.doRequest("GET", fullURL, headers, params, nil)
	if err != nil {
		return "", err
	}
	if isGtx {
		return utils.ExtractTranslatedTextFromArray(respBytes)
	}
	return utils.ExtractTranslatedText(respBytes)
}

func (s *GoogleTranslateService) callTranslatePa(text, target, endpoint string) (string, error) {
	params := url.Values{
		"query.target_language": {target},
		"key":                   {GOOGLE_API_KEY_TRANSLATE_PA},
		"query.text":            {text},
	}
	headers := map[string]string{
		"User-Agent": utils.GetConditionalRandomValue(DefaultUserAgents, s.opts.CustomUserAgents, s.opts.UseRandomUserAgents),
	}
	respBytes, err := s.doRequest("GET", endpoint, headers, params, nil)
	if err != nil {
		return "", err
	}
	return utils.ExtractTranslatedTextFromJson(respBytes)
}

func (s *GoogleTranslateService) doRequest(method, endpoint string, headers map[string]string, params url.Values, body []byte) ([]byte, error) {
	reqURL := endpoint
	if params != nil {
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		u.RawQuery = u.Query().Encode() + "&" + params.Encode()
		reqURL = u.String()
	}

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}
