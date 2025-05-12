# Go Translate

🚀 Go Translate is a free and unlimited go library that implemented Google Translate API and Microsoft Translate API.

## ✨ Features

- ✅ Supports multiple providers: Google, Microsoft
- 🔧 Customizable request headers, random user-agents, and token
- 🧪 Easy to extend with new providers
- 📦 Clean interface and modular design
---

## 📦 Installation

```bash
go get github.com/dinhcanh303/go_translate@latest
```

## 🧑‍💻 Usage

- Basic usage with Google Translate
```go
  package main
  import (
    "fmt"
    go_translate "github.com/dinhcanh303/go_translate"
  )
  func main() {
    translator := go_translate.NewTranslator(&go_translate.TranslateOptions{
      Provider: go_translate.ProviderGoogle,
    })

    result, err := translator.TranslateText(["Hello world","How are you"],"vi")
    if err != nil {
      fmt.Println("Error:", err)
      return
    }
    fmt.Println("Translated text:", result)
  }
```

- Basic usage with Microsoft Translate

```go
  package main
  import (
    "fmt"
    go_translate "github.com/dinhcanh303/go_translate"
  )
  func main() {
    translator := go_translate.NewTranslator(&go_translate.TranslateOptions{
      Provider: go_translate.ProviderMicrosoft,
    })

    result, err := translator.TranslateText(["Hello world","How are you"],"vi")
    if err != nil {
      fmt.Println("Error:", err)
      return
    }
    fmt.Println("Translated text:", result)
  }
```

## ⚙️ Options

```go
  type TranslateOptions struct {
	  // Provider specifies which translation provider to use (e.g., Google or Microsoft).
    Provider Provider

    //HTTPClient config
    HTTPClient *http.Client

    // GoogleAPIType specifies the API type to use for Google Translate (e.g., "html" || "pa-gtx" || "client-gtx" || "client-dict").
    GoogleAPIType GoogleAPIType

    // UseRandomUserAgents enables random selection of User-Agent headers for each request. (Only Google)
    UseRandomUserAgents bool

    // UseRandomServiceUrls enables random selection of base service URLs (e.g., multiple Google endpoints).
    UseRandomServiceUrls bool

    // AddToken indicates whether a token should be added to the request (used for some unofficial Google APIs).
    AddToken bool

    // CustomServiceUrls provides a list of service URLs to override the default service urls (used if random is enabled).
    CustomServiceUrls []string

    // CustomUserAgents provides a list of User-Agent strings to use (used if random is enabled).
    CustomUserAgents []string
  }
  const (
    // TypeHtml uses the standard "translate.google.com" HTML endpoint (unofficial, suitable for web-scraping style requests).
    TypeHtml GoogleAPIType = "html"

    // TypePaGtx uses the "translate-pa.googleapis.com" endpoint (used in some embedded Google services).
    TypePaGtx GoogleAPIType = "pa-gtx"

    // TypeClientGtx uses the "translate.googleapis.com" endpoint with the "client=gtx" query parameter.
    TypeClientGtx GoogleAPIType = "client-gtx"

    // TypeClientDictChromeEx uses "client=dict-chrome-ex" — used in Chrome dictionary extension.
    TypeClientDictChromeEx GoogleAPIType = "client-dict"

    //TypeRandom uses random multiple endpoint Google Translate API
    TypeRandom GoogleAPIType = "random"

    // TypeSequential tries APIs in a fixed order, one after another.
    TypeSequential GoogleAPIType = "sequential"

    // TypeMix of random and sequential (not implemented yet, placeholder for future use)
    TypeMix GoogleAPIType = "mix"
)

```
## Note
- Using the free Microsoft Translate API does not support automatic language detection, the default is en (English), if you want to automatically detect the language you can use a model that can detect the language. I am using the fasttext and opencc model to detect the language
- Server example: https://github.com/dinhcanh303/language_detection
- You can refer to the example folder for more information

```go
package example

import (
	"context"
	"testing"

	"github.com/dinhcanh303/go_translate"
	"github.com/dinhcanh303/go_translate/example/grpc_client"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	client, err := grpc_client.NewGRPCLanguageDetectionClient("127.0.0.1:50055") //Server model detect language
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
	result, err := translator.TranslateText([]string{text}, "vi", resp.DetectedLang)
	require.Nil(t, err)
	require.Equal(t, result, "Chào thế giới")
}
```
## 📄 License

- MIT License

## 🙌 Contributing

- Contributions are welcome! Feel free to fork, create issues or open pull requests.