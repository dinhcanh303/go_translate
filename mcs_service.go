package go_translate

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/dinhcanh303/go_translate/utils"
)

// MicrosoftTranslateService is a service for interacting with Microsoft's translation API.
type MicrosoftTranslateService struct {
	client *http.Client      // HTTP client to send requests
	opts   *TranslateOptions // Options that can be used for customizing translation behavior (e.g., API keys, etc.)
}

// NewMicrosoftTranslateService creates a new instance of MicrosoftTranslateService with the provided options.
func NewMicrosoftTranslateService(client *http.Client, opts *TranslateOptions) *MicrosoftTranslateService {
	return &MicrosoftTranslateService{
		client: client,
		opts:   opts,
	}
}

// TranslateText performs the translation of the provided text into the target language using the Microsoft translation API.
// It also optionally accepts a detected language code if you want to specify the source language explicitly.
func (m *MicrosoftTranslateService) TranslateText(ctx context.Context, texts []string, target string, detectedLangCode ...string) ([]string, error) {
	return m.translate(ctx, texts, target, detectedLangCode...)
}

// TranslateText performs the translation of the provided text into the target language using the Microsoft translation API.
// It also optionally accepts a detected language code if you want to specify the source language explicitly.
func (m *MicrosoftTranslateService) translate(ctx context.Context, texts []string, target string, detectedLangCode ...string) ([]string, error) {
	dir := "en/" + target
	if len(detectedLangCode) > 0 {
		dir = detectedLangCode[0] + "/" + target
	}

	formData := url.Values{
		"text":     {utils.JoinWithSeparator(texts)},
		"dir":      {dir},
		"provider": {"microsoft"},
	}
	req, err := http.NewRequestWithContext(ctx, "POST", MicrosoftServerUrl, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := m.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	text, err := utils.DecodeUnicode(string(bodyBytes))
	if err != nil {
		return nil, err
	}
	return utils.SplitWithSeparator(text), nil
}
