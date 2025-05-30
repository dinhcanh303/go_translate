// Package go_translate defines Google Translate API type options used for configuring request behavior.
package go_translate

// GoogleAPIType defines the type of Google Translate API endpoint to use.
type GoogleAPIType string

const (
	// TypeHtml uses the standard "translate.google.com" HTML endpoint (unofficial, suitable for web-scraping style requests).
	TypeHtml GoogleAPIType = "html"

	// TypePaGtx uses the "translate-pa.googleapis.com" endpoint (used in some embedded Google services).
	TypePaGtx GoogleAPIType = "pa-gtx"

	// TypeClientGtx uses the "translate.googleapis.com" endpoint with the "client=gtx" query parameter.
	TypeClientGtx GoogleAPIType = "client-gtx"

	// TypeClientDictChromeEx uses "client=dict-chrome-ex" — used in Chrome dictionary extension.
	TypeClientDictChromeEx GoogleAPIType = "client-dict"

	// TypeDictionary uses "key" — used in Chrome dictionary extension.
	TypeDictionary GoogleAPIType = "dictionary"

	//TypeRandom uses random multiple endpoint Google Translate API
	TypeRandom GoogleAPIType = "random"

	// TypeSequential tries APIs in a fixed order, one after another.
	TypeSequential GoogleAPIType = "sequential"

	// TypeMix of random and sequential (not implemented yet, placeholder for future use)
	TypeMix GoogleAPIType = "mix"
)
