package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

func ExtractTranslatedText(respBody []byte) (string, error) {
	var data [][]string
	err := json.Unmarshal(respBody, &data)
	if err != nil {
		return "", err
	}
	if len(data) > 0 && len(data[0]) > 0 {
		return data[0][0], nil
	}
	return "", errors.New("unexpected response format")
}

func DecodeUnicode(text string) (string, error) {
	jsonStr := fmt.Sprintf(`"%s"`, text)
	var output string
	err := json.Unmarshal([]byte(jsonStr), &output)
	if err != nil {
		return "", err
	}
	return output, nil
}

func ExtractTranslatedTextJson(respBody []byte) (string, error) {
	var result struct {
		Translation string `json:"translation"`
	}
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		return "", err
	}
	return result.Translation, nil
}
func ExtractTranslatedTextFromArray(data []byte) (string, error) {
	var rawData []interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return "", err
	}

	sentences, ok := rawData[0].([]interface{})
	if !ok {
		return "", errors.New("unexpected format: cannot extract first layer")
	}

	var builder strings.Builder
	for _, item := range sentences {
		segment, ok := item.([]interface{})
		if !ok || len(segment) == 0 {
			continue
		}

		first, ok := segment[0].(string)
		if ok {
			builder.WriteString(first)
		}
	}

	return builder.String(), nil
}

func GetRandom[T any](data []T) T {
	return data[rand.Intn(len(data))]
}
