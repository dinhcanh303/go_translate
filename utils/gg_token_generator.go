package utils

import (
	"fmt"
	"math"
	"unicode/utf8"
)

func GgTokenGenerate(text string) string {
	tokenKeys := [2]int{406398, 2087938574}
	encodedChars := []int{}
	for i := 0; i < getTextLength(text); i++ {
		charCode := getCharCodeAt(text, i)
		if charCode < 128 {
			encodedChars = append(encodedChars, charCode)
		} else {
			if charCode < 2048 {
				encodedChars = append(encodedChars, charCode>>6|192)
			} else {
				if (charCode&64512) == 55296 && i+1 < getTextLength(text) && (getCharCodeAt(text, i+1)&64512) == 56320 {
					charCode = 65536 + ((charCode & 1023) << 10) + (getCharCodeAt(text, i+1) & 1023)
					i++
					encodedChars = append(encodedChars, charCode>>18|240)
					encodedChars = append(encodedChars, charCode>>12&63|128)
				} else {
					encodedChars = append(encodedChars, charCode>>12|224)
				}
				encodedChars = append(encodedChars, charCode>>6&63|128)
			}
			encodedChars = append(encodedChars, charCode&63|128)
		}
	}

	token := tokenKeys[0]
	for _, char := range encodedChars {
		token += char
		token = applyTransformation(token, "+-a^+6")
	}
	token = applyTransformation(token, "+-3^+b+-f")
	token ^= tokenKeys[1]
	if token < 0 {
		token = (token & 2147483647) + 2147483648
	}
	token = int(math.Mod(float64(token), 1000000))

	return fmt.Sprintf("%d.%d", token, token^tokenKeys[0])
}

func getCharCodeAt(text string, index int) int {
	r, _ := utf8.DecodeRuneInString(text[index:])
	return int(r)
}

func getTextLength(text string) int {
	return utf8.RuneCountInString(text)
}

func unsignedRightShift(value int, shift int) int {
	if shift >= 32 || shift < -32 {
		multiple := shift / 32
		shift -= multiple * 32
	}
	if shift < 0 {
		shift += 32
	}
	if shift == 0 {
		return ((value>>1)&0x7fffffff)*2 + ((value >> shift) & 1)
	}
	if value < 0 {
		value >>= 1
		value &= 0x7fffffff
		value |= 0x40000000
		value >>= (shift - 1)
	} else {
		value >>= shift
	}
	return value
}

func applyTransformation(value int, transformation string) int {
	for i := 0; i < len(transformation)-2; i += 3 {
		operation := transformation[i+2]
		var operationValue int
		if operation >= 'a' {
			operationValue = int(operation) - 87
		} else {
			operationValue = int(operation - '0')
		}
		if transformation[i+1] == '+' {
			operationValue = unsignedRightShift(value, operationValue)
		} else {
			operationValue = value << operationValue
		}
		if transformation[i] == '+' {
			value = (value + operationValue) & 0xffffffff
		} else {
			value = value ^ operationValue
		}
	}
	return value
}
