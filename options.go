package go_translate

import "net/http"

// TranslateOptions defines configurable options for the translation service.
type TranslateOptions struct {

	// Provider specifies which translation provider to use (e.g., Google or Microsoft).
	Provider Provider

	//HTTPClient config
	HTTPClient *http.Client

	// GoogleAPIType specifies the API type to use for Google Translate (e.g., "html" || "pa-gtx" || "client-gtx" || "client-dict").
	GoogleAPIType GoogleAPIType

	// UseRandomUserAgents enables random selection of User-Agent headers for each request. (Only Google)
	UseRandomUserAgents bool

	// UseRandomServiceUrls enables random selection of base service URLs (e.g., multiple Google endpoints).
	UseRandomServiceUrls bool

	// AddToken indicates whether a token should be added to the request (used for some unofficial Google APIs).
	AddToken bool

	// CustomServiceUrls provides a list of service URLs to override the default service urls (used if random is enabled).
	CustomServiceUrls []string

	// CustomUserAgents provides a list of User-Agent strings to use (used if random is enabled).
	CustomUserAgents []string
}
