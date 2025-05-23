package go_translate

import (
	"context"
	"encoding/json"
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
	if m.opts.MicrosoftAPIType == TypeSmartLink {
		return m.callTranslateSmartLink(ctx, texts, target, detectedLangCode...)
	}
	return m.callTranslateEdge(ctx, texts, target)
}

// callTranslateEdge makes a POST request to the Edge API endpoint and returns the translated text.
func (m *MicrosoftTranslateService) callTranslateEdge(ctx context.Context, texts []string, target string) ([]string, error) {
	tokenBytes, err := utils.DoRequest(m.client, ctx, "GET", AuthEdgeUrl, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	var payload []map[string]string
	for _, text := range texts {
		payload = append(payload, map[string]string{
			"text": text,
		})
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	baseUrl := MicrosoftUrls[TypeEdge] + target
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": string(tokenBytes),
		"User-Agent":    utils.GetConditionalRandomValue(DefaultUserAgents, m.opts.CustomUserAgents, m.opts.UseRandomUserAgents),
	}
	resq, err := utils.DoRequest(m.client, ctx, "POST", baseUrl, header, nil, jsonPayload)
	if err != nil {
		return nil, err
	}
	return utils.ExtractTranslatedTextFromMCSEdge(resq)
}

// callTranslateSmartLink makes a POST request to the Microsoft translate API endpoint of smart link and returns the translated text.
func (m *MicrosoftTranslateService) callTranslateSmartLink(ctx context.Context, texts []string, target string, detectedLangCode ...string) ([]string, error) {
	// detectedLang , err := utils.DoRequest(m.client,ctx,"GET","")
	dir := "en/" + target
	if len(detectedLangCode) > 0 {
		dir = detectedLangCode[0] + "/" + target
	}

	formData := url.Values{
		"text":     {utils.JoinWithSeparator(texts)},
		"dir":      {dir},
		"provider": {"microsoft"},
	}
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
		"User-Agent":   utils.GetConditionalRandomValue(DefaultUserAgents, m.opts.CustomUserAgents, m.opts.UseRandomUserAgents),
	}
	resp, err := utils.DoRequest(m.client, ctx, "POST", MicrosoftServerUrl, header, formData, nil)
	if err != nil {
		return nil, err
	}
	text, err := utils.DecodeUnicode(string(resp))
	if err != nil {
		return nil, err
	}
	return utils.SplitWithSeparator(text), nil
}
