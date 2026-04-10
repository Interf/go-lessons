package unpack

import (
	"errors"
	"unicode"
)

var ErrUnpackDigitFirst = errors.New("first char in digit")
var ErrUnpackFoundNumber = errors.New("number found in digit")
var ErrUnpackRawEscapeCharWithoutDigit = errors.New("raw escape char without digit")

func Unpack(str string) (string, error) {

	if str == "" {
		return "", nil
	}

	var result []rune
	var prev rune
	var isPrevDigit bool

	var isRawString bool
	isRawString = str[0] == '`' && str[len(str)-1] == '`'

	if isRawString {

		return unpackRaw(str)

	}

	for index, char := range str {

		if unicode.IsDigit(char) {
			if prev == 0 {
				return "", ErrUnpackDigitFirst
			}
			if isPrevDigit {
				return "", ErrUnpackFoundNumber
			}

			countRepeat := int(char - '0')

			if countRepeat == 0 {
				result = result[:len(result)-1]
			} else {
				for i := 0; i < countRepeat-1; i++ {

					if index-2 >= 0 && index-2 < len(result) && str[index-2] == '\\' {
						result = append(result, '\\')
					}

					result = append(result, prev)
				}
			}

			isPrevDigit = true
		} else {
			result = append(result, char)
			prev = char
			isPrevDigit = false
		}

	}

	return string(result), nil
}

func unpackRaw(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	str = str[1 : len(str)-1]

	var result []rune
	var prev rune
	var isPrevDigit bool
	var isPrevEscapeChar bool

	for index, char := range str {

		if char == '\\' && !isPrevEscapeChar {
			isPrevEscapeChar = true
			continue
		}

		if unicode.IsDigit(char) && !isPrevEscapeChar {
			if prev == 0 {
				return "", ErrUnpackDigitFirst
			}
			if isPrevDigit {
				return "", ErrUnpackFoundNumber
			}

			countRepeat := int(char - '0')

			if countRepeat == 0 {
				result = result[:len(result)-1]
			} else {
				for i := 0; i < countRepeat-1; i++ {
					result = append(result, prev)
				}
			}

			isPrevDigit = true

		} else {

			if isPrevEscapeChar && !unicode.IsDigit(char) && index+1 < len(str) && !unicode.IsDigit(rune(str[index+1])) {
				return "", ErrUnpackRawEscapeCharWithoutDigit
			}

			result = append(result, char)
			prev = char
			isPrevDigit = false
		}

		isPrevEscapeChar = false

	}

	return string(result), nil
}
