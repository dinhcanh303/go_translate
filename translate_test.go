package go_translate

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTranslateBatchText(t *testing.T) {
	type TranslateTestCase struct {
		opts           *TranslateOptions
		input          []string
		targetLang     string
		detectedLang   string
		expectedOutput []string
	}
	tcs := map[string]TranslateTestCase{
		"google case 1": {
			opts: &TranslateOptions{Provider: ProviderGoogle,
				GoogleAPIType:        TypeMix,
				UseRandomUserAgents:  true,
				UseRandomServiceUrls: true,
				AddToken:             true,
				HTTPClient: &http.Client{
					Timeout: 15 * time.Second,
					Transport: &http.Transport{
						MaxIdleConns:        2000,
						MaxIdleConnsPerHost: 2000,
						IdleConnTimeout:     100 * time.Second,
					},
				}},
			input:          []string{"Thank you for using our package.", "한국어", "I'm fine", "我认为我们需要拭目以待。美联储加息可能会让市场更加动荡😑😑😑😑"},
			detectedLang:   "auto",
			targetLang:     "vi",
			expectedOutput: []string{"Cảm ơn bạn đã sử dụng gói dịch vụ của chúng tôi.", "Hàn Quốc", "Tôi ổn", "Tôi nghĩ chúng ta cần phải chờ xem. Việc Fed tăng lãi suất có thể khiến thị trường biến động hơn 😑😑😑😑"},
		},
		"google case 2": {
			opts:           &TranslateOptions{Provider: "google", GoogleAPIType: TypeClientGtx},
			input:          []string{"Thank you for using our package."},
			detectedLang:   "auto",
			targetLang:     "vi",
			expectedOutput: []string{"Cảm ơn bạn đã sử dụng gói của chúng tôi."},
		},
		"microsoft case 1": {
			opts:           &TranslateOptions{Provider: "microsoft"},
			input:          []string{"Thank you for using our package."},
			detectedLang:   "en",
			targetLang:     "vi",
			expectedOutput: []string{"Cảm ơn bạn đã sử dụng gói của chúng tôi."},
		},
	}
	ctx := context.Background()
	for scenario, tc := range tcs {
		tc := tc
		t.Run(scenario, func(t *testing.T) {
			translator, err := NewTranslator(tc.opts)
			require.Nil(t, err)
			require.NotNil(t, translator)

			fmt.Println(tc.input)
			result, err := translator.TranslateText(ctx, tc.input, tc.targetLang)
			require.Nil(t, err)
			require.Equal(t, tc.expectedOutput, result)
		})
	}
}
