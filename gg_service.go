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
	googleApiType := s.opts.GoogleAPIType
	if s.opts.GoogleAPIType == TypeRandom {
		googleApiType = utils.GetRandomValue(GoogleAPITypeSupport)
	}
	endpoint := GoogleUrls[googleApiType]
	handlers := s.getAPIHandlers()
	handler, ok := handlers[googleApiType]
	if !ok {
		return nil, errors.New("unsupported Google API type: " + string(googleApiType))
	}
	translatedText, err := handler(ctx, texts, target, endpoint)
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
	return s.executeAPIRequest(ctx, "POST", endpoint, headers, nil, []byte(body), utils.ExtractTranslatedTextFromHtml)
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
	extractFunc := utils.ExtractTranslatedText
	if isGtx {
		extractFunc = utils.ExtractTranslatedTextFromArray
	}
	return s.executeAPIRequest(ctx, "GET", fullURL, headers, params, nil, extractFunc)
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

	return s.executeAPIRequest(ctx, "GET", endpoint, headers, params, nil, utils.ExtractTranslatedTextFromJson)
}

// callTranslateGet makes a GET request to the Google Translate API (client-gtx or client-dict) and returns the translated text.
func (s *GoogleTranslateService) callTranslateSequential(ctx context.Context, texts []string, target string) ([]string, error) {
	handlers := s.getAPIHandlers()
	for apiType, endpoint := range GoogleUrls {
		handler, ok := handlers[apiType]
		if !ok {
			continue // skip unsupported apiTypes
		}
		translatedText, err := handler(ctx, texts, target, endpoint)
		if err == nil && translatedText != nil {
			return translatedText, nil
		}
		log.Printf("[ERROR] Sequential API %s failed: %v", apiType, err)
	}
	return nil, errors.New("unable to translate text, all APIs failed")
}
func (s *GoogleTranslateService) callTranslateMix(
	ctx context.Context,
	texts []string,
	target string,
) ([]string, error) {
	// googleApiType := utils.GetRandomValue(GoogleAPITypeSupport)
	googleApiType := TypeHtml
	endpoint := GoogleUrls[googleApiType]
	handlers := s.getAPIHandlers()
	handler, ok := handlers[googleApiType]
	if !ok {
		return nil, errors.New("unsupported Google API type: " + string(googleApiType))
	}
	translatedText, err := handler(ctx, texts, target, endpoint)
	if err == nil && translatedText != nil {
		return translatedText, nil
	}
	var remainHandlers = make(map[GoogleAPIType]apiHandler)
	exclude := map[GoogleAPIType]struct{}{
		googleApiType:  {},
		TypeSequential: {},
		TypeMix:        {},
	}
	for k, v := range handlers {
		if _, skip := exclude[k]; skip {
			continue
		}
		remainHandlers[k] = v
	}
	for apiType, endpoint := range GoogleUrls {
		handler, ok := remainHandlers[apiType]
		if !ok {
			continue // skip unsupported apiTypes
		}
		translatedText, err := handler(ctx, texts, target, endpoint)
		if err == nil && translatedText != nil {
			return translatedText, nil
		}
		log.Printf("[ERROR] Sequential API %s failed: %v", apiType, err)
	}
	return nil, errors.New("unable to translate text, all APIs failed")
}

// doRequest is a helper function that performs the HTTP request (GET or POST) with the given method, URL, headers, parameters, and body.
// It returns the response body as a byte slice or an error if the request fails.
func (s *GoogleTranslateService) doRequest(ctx context.Context, method, endpoint string, headers map[string]string, params url.Values, body []byte) ([]byte, error) {
	reqURL := buildRequestURL(endpoint, params)
	reqBody := buildRequestBody(body)

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

	if err := handleHTTPError(resp); err != nil {
		return nil, err
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

// handleHTTPError checks the HTTP response status and returns an error if it's not successful.
func handleHTTPError(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}
	return nil
}

// buildRequestURL constructs the full request URL with parameters.
func buildRequestURL(endpoint string, params url.Values) string {
	if params == nil {
		return endpoint
	}
	u, _ := url.Parse(endpoint)
	u.RawQuery = u.Query().Encode() + "&" + params.Encode()
	return u.String()
}

// buildRequestBody creates the request body from a byte slice.
func buildRequestBody(body []byte) io.Reader {
	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}
	return reqBody
}

type apiHandler func(ctx context.Context, texts []string, target, endpoint string) ([]string, error)

func (s *GoogleTranslateService) getAPIHandlers() map[GoogleAPIType]apiHandler {
	return map[GoogleAPIType]apiHandler{
		TypeHtml: func(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
			return s.callTranslateHTML(ctx, texts, target, endpoint)
		},
		TypeClientGtx: func(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
			return s.callTranslateGet(ctx, texts, target, endpoint, true)
		},
		TypeClientDictChromeEx: func(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
			return s.callTranslateGet(ctx, texts, target, endpoint, false)
		},
		TypePaGtx: func(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
			return s.callTranslatePa(ctx, texts, target, endpoint)
		},
		TypeSequential: func(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
			return s.callTranslateSequential(ctx, texts, target)
		},
		TypeMix: func(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
			return s.callTranslateMix(ctx, texts, target)
		},
	}
}

// executeAPIRequest handles the common logic for making API requests.
func (s *GoogleTranslateService) executeAPIRequest(
	ctx context.Context,
	method string,
	endpoint string,
	headers map[string]string,
	params url.Values,
	body []byte,
	extractFunc func([]byte) ([]string, error),
) ([]string, error) {
	respBytes, err := s.doRequest(ctx, method, endpoint, headers, params, body)
	if err != nil {
		return nil, err
	}
	return extractFunc(respBytes)
}
