// Package go_translate provides a unified interface for working with multiple translation providers
// such as Google Translate and Microsoft Translator.
package go_translate

// Translator is a generic interface for text translation services.
//
// Implementations may include Google Translate, Microsoft Translator, or other providers.
// The method TranslateText translates input text into the target language.
type Translator interface {
	// TranslateText translates the input `text` into the `target` language code (e.g., "en", "vi").
	// Optionally, a detected source language code can be provided to skip language detection.
	TranslateText(text, target string, detectedLangCode ...string) (string, error)
}

// NewTranslator returns a Translator implementation based on the given TranslateOptions.
//
// If no options are provided, it defaults to using Google Translate with HTML API type.
// This function panics if the provider is unsupported.
func NewTranslator(opts ...*TranslateOptions) Translator {
	var opts1 *TranslateOptions
	if len(opts) > 0 && opts[0] != nil {
		opts1 = opts[0]
	}
	// Apply default values if missing
	if opts1.Provider == "" {
		opts1.Provider = ProviderGoogle
	}
	if opts1.GoogleAPIType == "" {
		opts1.GoogleAPIType = TypeHtml
	}
	// Create appropriate service based on provider
	switch opts1.Provider {
	case ProviderGoogle:
		return NewGoogleTranslateService(opts1)
	case ProviderMicrosoft:
		return NewMicrosoftTranslateService(opts1)
	default:
		panic("unsupported provider: " + string(opts1.Provider))
	}
}
