package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

// ExtractTranslatedTextFromHtml extracts the translated text from a given response body in JSON format.
// The response body is expected to be a nested JSON array where the first element contains the translated text.
// Returns the translated text or an error if the format is unexpected.
func ExtractTranslatedTextFromHtml(respBody []byte) ([]string, error) {
	var data [][]string
	err := json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 && len(data[0]) > 0 {
		return data[0], nil
	}
	return nil, errors.New("unexpected response format")
}

// ExtractTranslatedText extracts the translated text from a given response body in JSON format.
// The response body is expected to be a nested JSON array where the first element contains the translated text.
// Returns the translated text or an error if the format is unexpected.
func ExtractTranslatedText(respBody []byte) ([]string, error) {
	var data [][]string
	err := json.Unmarshal(respBody, &data)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 && len(data[0]) > 0 {
		return SplitWithSeparator(data[0][0]), nil
	}
	return nil, errors.New("unexpected response format")
}

// DecodeUnicode decodes Unicode escape sequences in a string.
// For example, it converts `\u003c` to `<`.
// Returns the decoded string or an error if decoding fails.
func DecodeUnicode(text string) (string, error) {
	jsonStr := fmt.Sprintf(`"%s"`, text)
	var output string
	err := json.Unmarshal([]byte(jsonStr), &output)
	if err != nil {
		return "", err
	}
	return output, nil
}

// ExtractTranslatedTextFromJson extracts the translated text from a JSON response with a "translation" field.
// The response is expected to contain a "translation" field with the translated text.
// Returns the translated text or an error if unmarshalling fails.
func ExtractTranslatedTextFromJson(respBody []byte) ([]string, error) {
	var result struct {
		Translation string `json:"translation"`
	}
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}
	return SplitWithSeparator(result.Translation), nil
}

// ExtractTranslatedTextFromArray extracts the translated text from a JSON array where each element is a sentence.
// The function expects the first layer of the JSON array to be a list of sentences.
// It concatenates and returns the translated sentences as a single string.
func ExtractTranslatedTextFromArray(data []byte) ([]string, error) {
	var rawData []interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return nil, err
	}

	sentences, ok := rawData[0].([]interface{})
	if !ok {
		return nil, errors.New("unexpected format: cannot extract first layer")
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

	return SplitWithSeparator(builder.String()), nil
}

// GetRandomValue returns a random value from a slice of type T.
// T is a generic type, so the function works for any type of slice.
func GetRandomValue[T any](data []T) T {
	return data[rand.Intn(len(data))]
}

// GetConditionalRandomValue returns a random value from the data slice if the condition is true,
// otherwise it returns a value from the defaultData slice. If the condition is false, the default value is returned.
func GetConditionalRandomValue[T any](defaultData, data []T, condition bool) T {
	if condition && len(data) > 0 {
		return GetRandomValue(data)
	}
	if condition {
		return GetRandomValue(defaultData)
	}
	return defaultData[0]
}

// JoinWithSeparator joins elements of a slice with a separator.
func JoinWithSeparator(input []string, separator ...string) string {
	return strings.Join(input, getSeparator(separator))
}

// SplitWithSeparator splits a string into a slice using a separator.
func SplitWithSeparator(input string, separator ...string) []string {
	return strings.Split(input, getSeparator(separator))
}

// getSeparator returns the first separator if provided or default to "\n".
func getSeparator(separator []string) string {
	if len(separator) > 0 && separator[0] != "" {
		return separator[0]
	}
	return "\n"
}
