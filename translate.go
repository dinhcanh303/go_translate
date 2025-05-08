// Package go_translate provides a unified interface for working with multiple translation providers
// such as Google Translate and Microsoft Translator.
package go_translate

import (
	"errors"
	"net/http"
	"time"
)

// Translator is a generic interface for text translation services.
//
// Implementations may include Google Translate, Microsoft Translator, or other providers.
// The method TranslateText translates input text into the target language.
type Translator interface {
	// TranslateText translates the input `text` into the `target` language code (e.g., "en", "vi").
	// Optionally, a detected source language code can be provided to skip language detection.
	TranslateText(text []string, target string, detectedLangCode ...string) ([]string, error)
}

// NewTranslator returns a Translator implementation based on the given TranslateOptions.
//
// If no options are provided, it defaults to using Google Translate with HTML API type.
// This function panics if the provider is unsupported.
func NewTranslator(opts ...*TranslateOptions) (Translator, error) {
	options, err := validateOptions(opts...)
	if err != nil {
		return nil, err
	}
	var client *http.Client
	if options.HTTPClient != (&http.Client{}) {
		client = options.HTTPClient
	} else {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	// Create appropriate service based on provider
	switch options.Provider {
	case ProviderGoogle:
		return NewGoogleTranslateService(client, options), nil
	case ProviderMicrosoft:
		return NewMicrosoftTranslateService(client, options), nil
	default:
		return nil, errors.New("unsupported provider: " + string(options.Provider))
	}
}
func validateOptions(opts ...*TranslateOptions) (*TranslateOptions, error) {
	var options *TranslateOptions
	if len(opts) > 0 && opts[0] != nil {
		options = opts[0]
	}
	// Apply default values if missing
	if options.Provider == "" {
		options.Provider = ProviderGoogle
	}
	if options.Provider == ProviderGoogle {
		if options.GoogleAPIType == "" {
			options.GoogleAPIType = TypeHtml
		}
		// Set default GoogleAPIType
		validTypes := map[GoogleAPIType]struct{}{TypeHtml: {}, TypeClientGtx: {}, TypeClientDictChromeEx: {}, TypePaGtx: {}, TypeRandom: {}}
		if _, ok := validTypes[options.GoogleAPIType]; !ok {
			return nil, errors.New("unsupported Google API Type, please check list supported")
		}
	}

	return options, nil
}
