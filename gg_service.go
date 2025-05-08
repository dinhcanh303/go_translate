// Package go_translate implements translation services for various providers, including Google Translate.
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

// GoogleTranslateService is a concrete implementation of the Translator interface for Google Translate.
// It supports multiple API endpoints and handles requests for different Google Translate API types.
type GoogleTranslateService struct {
	client *http.Client      // HTTP client used for making API requests
	opts   *TranslateOptions // Options for configuring the translation service
}

// NewGoogleTranslateService creates a new instance of GoogleTranslateService with the given options.
// The HTTP client timeout is set to 5 seconds.
func NewGoogleTranslateService(opts *TranslateOptions) *GoogleTranslateService {
	return &GoogleTranslateService{
		client: &http.Client{Timeout: 5 * time.Second},
		opts:   opts,
	}
}

// TranslateText translates the provided text into the target language using the configured provider and API type.
// It tries multiple endpoints based on the API type (HTML, PaGtx, ClientGtx, etc.) and returns the translated text.
// It returns an error if all translation attempts fail.
func (s *GoogleTranslateService) TranslateText(text, target string, detectedLangCode ...string) (string, error) {
	// Try each API endpoint
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
		// If translation is successful, return the result
		if err == nil && translated != "" {
			return translated, nil
		}
		log.Printf("[ERROR] API %s failed: %v", t, err)
	}
	return "", errors.New("unable to translate text, all APIs failed")
}

// callTranslateHTML makes a POST request to the HTML API endpoint and returns the translated text.
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

// callTranslateGet makes a GET request to the Google Translate API (client-gtx or client-dict) and returns the translated text.
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

// callTranslatePa makes a GET request to the PaGtx API endpoint and returns the translated text.
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

// doRequest is a helper function that performs the HTTP request (GET or POST) with the given method, URL, headers, parameters, and body.
// It returns the response body as a byte slice or an error if the request fails.
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
