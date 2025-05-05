package go_translate

type TranslateOptions struct {
	Provider          Provider
	AddToken          bool
	RandomUserAgents  bool
	RandomServiceUrls bool
	ServiceUrls       []string
}
