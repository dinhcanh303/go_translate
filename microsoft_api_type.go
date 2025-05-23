package go_translate

type MicrosoftAPIType string

const (
	// TypeSmartLink uses the standard "https://webmail.smartlinkcorp.com" endpoint
	TypeSmartLink MicrosoftAPIType = "smart-link"

	// TypeEdge uses the standard "https://api-edge.cognitive.microsofttranslator.com" endpoint
	TypeEdge MicrosoftAPIType = "edge"
)
