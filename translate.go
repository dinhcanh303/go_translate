package go_translate

type Translator interface {
	TranslateText(text, target string, detectedLangCode ...string) (string, error)
}

func NewTranslator(opts *TranslateOptions) Translator {
	if opts == nil {
		opts = &TranslateOptions{
			Provider:          ProviderGoogle,
			AddToken:          false,
			RandomUserAgents:  false,
			RandomServiceUrls: false,
		}
	}
	switch opts.Provider {
	case ProviderGoogle:
		return NewGoogleTranslateService(opts)
	case ProviderMicrosoft:
		return NewMicrosoftTranslateService(opts)
	default:
		panic("unsupported provider")
	}
}
