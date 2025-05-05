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
go get github.com/dinhcanh303/go_translate
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

    result, err := translator.TranslateText("Hello world","vi")
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

    result, err := translator.TranslateText("Hello world","vi")
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
    Provider          Provider // "google" or "microsoft"
    AddToken          bool     // For Google Translate token generation
    RandomUserAgents  bool     // Randomize user-agent header
    RandomServiceUrls bool     // Shuffle between service URLs
    ServiceUrls       []string // Custom backend servers (Google only)
  }
```
## Note
- Using the free Microsoft Translate API does not support automatic language detection, the default is en (English), if you want to automatically detect the language you can use a model that can detect the language. I am using the fasttext and opencc model to detect the language
- Server example: 
  -- https://github.com/dinhcanh303/language_detection
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
	result, err := translator.TranslateText(text, "vi", resp.DetectedLang)
	require.Nil(t, err)
	require.Equal(t, result, "Chào thế giới")
}
```
## 📄 License

- MIT License

## 🙌 Contributing

- Contributions are welcome! Feel free to fork, create issues or open pull requests.