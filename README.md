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

## 📄 License

- MIT License

## 🙌 Contributing

- Contributions are welcome! Feel free to fork, create issues or open pull requests.