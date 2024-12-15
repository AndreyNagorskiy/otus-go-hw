package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const zeroRune = '0'

func Unpack(input string) (string, error) {
	var resBuilder strings.Builder
	var lastChar rune

	for _, char := range input {
		if unicode.IsDigit(char) {
			if lastChar == 0 || unicode.IsDigit(lastChar) {
				return "", ErrInvalidString
			}

			if char == zeroRune && resBuilder.Len() > 0 {
				runesRes := []rune(resBuilder.String())
				if len(runesRes) > 0 {
					runesRes = runesRes[:len(runesRes)-1]
				}

				resBuilder.Reset()
				resBuilder.WriteString(string(runesRes))

				continue
			}

			repeatCount, err := strconv.Atoi(string(char))
			if err != nil {
				continue
			}

			resBuilder.WriteString(strings.Repeat(string(lastChar), repeatCount-1))
		} else {
			resBuilder.WriteRune(char)
		}
		lastChar = char
	}

	return resBuilder.String(), nil
}
