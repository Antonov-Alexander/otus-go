package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder
	var packedChar rune
	var isEscapedChar bool

	runes := []rune(input)
	runesLen := len(runes)

	for index, char := range runes {
		isEscape, isNumber, isLetter := getCharProps(char)
		isLast := index == runesLen-1
		isEmptyPackedChar := packedChar != 0

		if isNumber && isEmptyPackedChar {
			intVar, err := strconv.Atoi(string(char))
			if err != nil {
				return "", ErrInvalidString
			}

			if intVar != 0 {
				result.WriteString(strings.Repeat(string(packedChar), intVar))
			}

			packedChar = 0
			continue
		}

		if isEmptyPackedChar {
			isEscapedChar = false

			if isLetter || isEscape {
				result.WriteRune(packedChar)
				packedChar = 0
			} else {
				return "", ErrInvalidString
			}
		}

		if (isNumber && !isEscapedChar) || (isLetter && isEscapedChar) {
			return "", ErrInvalidString
		}

		if isEscape && !isEscapedChar {
			if isLast {
				return "", ErrInvalidString
			}

			isEscapedChar = true
			continue
		}

		if isLast {
			result.WriteRune(char)
			break
		}

		packedChar = char
	}

	return result.String(), nil
}

func getCharProps(char rune) (isEscape bool, isNumber bool, isLetter bool) {
	isEscape = char == 92                                                // "/" - 92
	isNumber = char >= 48 && char <= 57                                  // 0-9 - 48-57
	isLetter = (char >= 65 && char <= 90) || (char >= 97 && char <= 122) // A-Z - 65-90, a-z - 97-122
	return
}
