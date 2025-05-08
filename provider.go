// Package go_translate defines supported translation providers.
package go_translate

// Provider defines the name of a translation service provider.
type Provider string

const (
	// ProviderGoogle represents the Google Translate provider.
	ProviderGoogle Provider = "google"

	// ProviderMicrosoft represents the Microsoft Translator provider.
	ProviderMicrosoft Provider = "microsoft"
)
