// Package go_translate implements translation services for various providers, including Google Translate.
package go_translate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

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
func NewGoogleTranslateService(client *http.Client, opts *TranslateOptions) *GoogleTranslateService {
	return &GoogleTranslateService{
		client: client,
		opts:   opts,
	}
}

// It tries multiple endpoints based on the API type (HTML, PaGtx, ClientGtx, etc.) and returns the translated text.
// It returns an error if all translation attempts fail.
func (s *GoogleTranslateService) translate(ctx context.Context, texts []string, target string) ([]string, error) {
	var translatedText []string
	var err error
	googleApiType := s.opts.GoogleAPIType
	if s.opts.GoogleAPIType == TypeRandom {
		googleApiType = utils.GetRandomValue([]GoogleAPIType{TypeHtml, TypeClientGtx, TypeClientDictChromeEx, TypePaGtx})
	}
	endpoint := GoogleUrls[googleApiType]
	// Try each API endpoint
	switch googleApiType {
	case TypeHtml:
		translatedText, err = s.callTranslateHTML(ctx, texts, target, endpoint)
	case TypeClientGtx, TypeClientDictChromeEx:
		translatedText, err = s.callTranslateGet(ctx, texts, target, endpoint, googleApiType == TypeClientGtx)
	case TypePaGtx:
		translatedText, err = s.callTranslatePa(ctx, texts, target, endpoint)
	}
	// If translation is successful, return the result
	if err == nil && translatedText != nil {
		return translatedText, nil
	}
	log.Printf("[ERROR] API %s failed: %v", googleApiType, err)
	return nil, errors.New("unable to translate text, all APIs failed")
}

// TranslateBatchText translates the provided text into the target language using the configured provider and API type.
// It returns an error if all translation attempts fail.
func (s *GoogleTranslateService) TranslateText(ctx context.Context, texts []string, target string, detectedLangCode ...string) ([]string, error) {
	return s.translate(ctx, texts, target)
}

// callTranslateHTML makes a POST request to the HTML API endpoint and returns the translated text.
func (s *GoogleTranslateService) callTranslateHTML(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
	body := buildGoogleHTMLBody(texts, target)
	headers := map[string]string{
		"User-Agent":     utils.GetConditionalRandomValue(DefaultUserAgents, s.opts.CustomUserAgents, s.opts.UseRandomUserAgents),
		"Content-Type":   "application/json+protobuf",
		"X-Goog-API-Key": GOOGLE_API_KEY_TRANSLATE,
	}
	respBytes, err := s.doRequest(ctx, "POST", endpoint, headers, nil, []byte(body))
	if err != nil {
		return nil, err
	}
	return utils.ExtractTranslatedTextFromHtml(respBytes)
}

// callTranslateGet makes a GET request to the Google Translate API (client-gtx or client-dict) and returns the translated text.
func (s *GoogleTranslateService) callTranslateGet(ctx context.Context, texts []string, target, endpoint string, isGtx bool) ([]string, error) {
	text := utils.JoinWithSeparator(texts)
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
	respBytes, err := s.doRequest(ctx, "GET", fullURL, headers, params, nil)
	if err != nil {
		return nil, err
	}
	if isGtx {
		return utils.ExtractTranslatedTextFromArray(respBytes)
	}
	return utils.ExtractTranslatedText(respBytes)
}

// callTranslatePa makes a GET request to the PaGtx API endpoint and returns the translated text.
func (s *GoogleTranslateService) callTranslatePa(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
	params := url.Values{
		"query.target_language": {target},
		"key":                   {GOOGLE_API_KEY_TRANSLATE_PA},
		"query.text":            {utils.JoinWithSeparator(texts)},
	}
	headers := map[string]string{
		"User-Agent": utils.GetConditionalRandomValue(DefaultUserAgents, s.opts.CustomUserAgents, s.opts.UseRandomUserAgents),
	}
	respBytes, err := s.doRequest(ctx, "GET", endpoint, headers, params, nil)
	if err != nil {
		return nil, err
	}
	return utils.ExtractTranslatedTextFromJson(respBytes)
}

// doRequest is a helper function that performs the HTTP request (GET or POST) with the given method, URL, headers, parameters, and body.
// It returns the response body as a byte slice or an error if the request fails.
func (s *GoogleTranslateService) doRequest(ctx context.Context, method, endpoint string, headers map[string]string, params url.Values, body []byte) ([]byte, error) {
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
	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
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
func buildGoogleHTMLBody(texts []string, target string) string {
	var quoted string
	if len(texts) == 1 {
		quoted = `"` + texts[0] + `"`
	} else {
		quoted = `"` + strings.Join(texts, `","`) + `"`
	}
	return fmt.Sprintf(`[[[%s],"auto","%s"],"wt_lib"]`, quoted, target)
}
