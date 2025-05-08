package go_translate

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dinhcanh303/go_translate/utils"
)

// MicrosoftTranslateService is a service for interacting with Microsoft's translation API.
type MicrosoftTranslateService struct {
	client *http.Client      // HTTP client to send requests
	opts   *TranslateOptions // Options that can be used for customizing translation behavior (e.g., API keys, etc.)
}

// NewMicrosoftTranslateService creates a new instance of MicrosoftTranslateService with the provided options.
func NewMicrosoftTranslateService(opts *TranslateOptions) *MicrosoftTranslateService {
	return &MicrosoftTranslateService{
		client: &http.Client{Timeout: 5 * time.Second},
		opts:   opts,
	}
}

// TranslateText performs the translation of the provided text into the target language using the Microsoft translation API.
// It also optionally accepts a detected language code if you want to specify the source language explicitly.
func (m *MicrosoftTranslateService) TranslateText(text, target string, detectedLangCode ...string) (string, error) {
	dir := "en/" + target
	if len(detectedLangCode) > 0 {
		dir = detectedLangCode[0] + "/" + target
	}
	formData := url.Values{
		"text":     {text},
		"dir":      {dir},
		"provider": {"microsoft"},
	}
	req, err := http.NewRequest("POST", MicrosoftServerUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := m.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	text, err = utils.DecodeUnicode(string(bodyBytes))
	if err != nil {
		return "", err
	}
	return text, nil
}
