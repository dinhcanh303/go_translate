package go_translate

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/dinhcanh303/go_translate/utils"
)

type MicrosoftTranslateService struct {
	client *http.Client
	opts   *TranslateOptions
}

func NewMicrosoftTranslateService(opts *TranslateOptions) *MicrosoftTranslateService {
	return &MicrosoftTranslateService{
		client: &http.Client{Timeout: 5 * time.Second},
		opts:   opts,
	}
}

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
