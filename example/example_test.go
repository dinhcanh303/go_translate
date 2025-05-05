package example

import (
	"context"
	"testing"

	"github.com/dinhcanh303/go_translate"
	"github.com/dinhcanh303/go_translate/example/grpc_client"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	client, err := grpc_client.NewGRPCLanguageDetectionClient("127.0.0.1:50055")
	require.Nil(t, err)
	require.NotNil(t, client)
	text := "Hello world"
	resp, err := client.DetectLanguage(context.Background(), text)
	require.Nil(t, err)
	require.NotNil(t, resp)
	translator := go_translate.NewTranslator(&go_translate.TranslateOptions{
		Provider: "microsoft",
	})
	require.NotNil(t, translator)
	result, err := translator.TranslateText(text, "vi", resp.DetectedLang)
	require.Nil(t, err)
	require.Equal(t, result, "Chào thế giới")
}
