package go_translate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTranslator_Google(t *testing.T) {
	opts := &TranslateOptions{Provider: "google"}
	translator := NewTranslator(opts)
	require.NotNil(t, translator)
	result, err := translator.TranslateText("Thank you for using our package.", "vi")
	require.Nil(t, err)
	require.Equal(t, result, "Cảm ơn bạn đã sử dụng gói dịch vụ của chúng tôi.")
}

func TestNewTranslator_Microsoft(t *testing.T) {
	opts := &TranslateOptions{Provider: "microsoft"}
	translator := NewTranslator(opts)
	require.NotNil(t, translator)
	result, err := translator.TranslateText("Thank you for using our package.", "vi")
	require.Nil(t, err)
	require.Equal(t, result, "Cảm ơn bạn đã sử dụng gói của chúng tôi.")
}
