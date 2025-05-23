// Package go_translate implements translation services for various providers, including Google Translate.
package go_translate

import (
	"context"
	"errors"
	"fmt"
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
		"X-Goog-API-Key": s.opts.GoogleAPIKeyTranslateHtml,
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
		"key":                   {s.opts.GoogleAPIKeyTranslatePa},
		"query.text":            {utils.JoinWithSeparator(texts)},
	}
	headers := map[string]string{
		"User-Agent": utils.GetConditionalRandomValue(DefaultUserAgents, s.opts.CustomUserAgents, s.opts.UseRandomUserAgents),
	}

	return s.executeAPIRequest(ctx, "GET", endpoint, headers, params, nil, utils.ExtractTranslatedTextFromJson)
}

// callTranslatePa makes a GET request to the PaGtx API endpoint and returns the translated text.
func (s *GoogleTranslateService) callTranslateDic(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
	params := url.Values{
		"language": {target},
		"key":      {s.opts.GoogleAPIKeyTranslateDic},
		"term":     {utils.JoinWithSeparator(texts)},
	}
	headers := map[string]string{
		"User-Agent": utils.GetConditionalRandomValue(DefaultUserAgents, s.opts.CustomUserAgents, s.opts.UseRandomUserAgents),
		"x-referer":  "chrome-extension://mgijmajocgfcbeboacabfgobmjgjcoja",
	}

	return s.executeAPIRequest(ctx, "GET", endpoint, headers, params, nil, utils.ExtractTranslatedTextFromGGDic)
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

func buildGoogleHTMLBody(texts []string, target string) string {
	var quoted string
	if len(texts) == 1 {
		quoted = `"` + texts[0] + `"`
	} else {
		quoted = `"` + strings.Join(texts, `","`) + `"`
	}
	return fmt.Sprintf(`[[[%s],"auto","%s"],"wt_lib"]`, quoted, target)
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
		TypeDictionary: func(ctx context.Context, texts []string, target, endpoint string) ([]string, error) {
			return s.callTranslateDic(ctx, texts, target, endpoint)
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
	respBytes, err := utils.DoRequest(s.client, ctx, method, endpoint, headers, params, body)
	if err != nil {
		return nil, err
	}
	return extractFunc(respBytes)
}
