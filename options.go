package go_translate

type TranslateOptions struct {
	Provider             Provider
	GoogleAPIType        GoogleAPIType
	UseRandomUserAgents  bool
	UseRandomServiceUrls bool
	AddToken             bool
	CustomServiceUrls    []string
	CustomUserAgents     []string
}
