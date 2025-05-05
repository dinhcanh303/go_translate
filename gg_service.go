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
	for i, endpoint := range GoogleUrls {
		var translated string
		var err error
		switch i {
		case 1, 2:
			translated, err = s.callTranslateGet(text, target, endpoint, i == 2)
		case 3:
			translated, err = s.callTranslatePa(text, target, endpoint)
		case 4:
			translated, err = s.callTranslateHTML(text, target, endpoint)
		}

		if err == nil && translated != "" {
			return translated, nil
		}
		log.Printf("[ERROR] API %d failed: %v", i, err)
	}
	return "", errors.New("unable to translate text, all APIs failed")
}

func (s *GoogleTranslateService) callTranslateHTML(text, target, endpoint string) (string, error) {
	body := fmt.Sprintf(`[[["%s"],"auto","%s"],"wt_lib"]`, text, target)
	headers := map[string]string{
		"Accept":         "/",
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
	sv := DefaultServiceUrls[0]
	if s.opts.RandomServiceUrls && len(s.opts.ServiceUrls) > 0 {
		sv = utils.GetRandom(s.opts.ServiceUrls)
	}
	if s.opts.RandomServiceUrls {
		sv = utils.GetRandom(DefaultServiceUrls)
	}
	fullURL := "https://" + sv + endpoint

	params := url.Values{
		"tl": {target},
		"q":  {text},
	}
	if s.opts.AddToken {
		params.Set("tk", utils.GgTokenGenerate(text))
	}
	userAgent := UserAgents[0]
	if s.opts.RandomUserAgents {
		userAgent = utils.GetRandom(UserAgents)
	}
	headers := map[string]string{
		"User-Agent": userAgent,
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
	respBytes, err := s.doRequest("GET", endpoint, nil, params, nil)
	if err != nil {
		return "", err
	}
	return utils.ExtractTranslatedTextJson(respBytes)
}

func (s *GoogleTranslateService) doRequest(method, endpoint string, headers map[string]string, params url.Values, body []byte) ([]byte, error) {
	reqURL := endpoint
	if params != nil {
		u, err := url.Parse(endpoint)
		if err != nil {
			return nil, err
		}
		u.RawQuery = params.Encode()
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
