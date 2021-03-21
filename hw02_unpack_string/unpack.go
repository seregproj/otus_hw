package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func validateString(str string) error {
	if !utf8.ValidString(str) {
		return ErrInvalidString
	}

	if unicode.IsDigit([]rune(str)[:1][0]) {
		return ErrInvalidString
	}

	return nil
}

func Unpack(str string) (string, error) {
	err := validateString(str)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	var replyRune rune

	for i, v := range str {
		prevRune, _ := utf8.DecodeLastRune([]byte(str[:i]))

		switch {
		case unicode.IsDigit(v):
			switch {
			case replyRune != 0:
				multiplier, err := strconv.Atoi(string(v))
				if err != nil {
					return "", err
				}

				if multiplier != 0 {
					sb.WriteString(strings.Repeat(string(replyRune), multiplier))
				}

				replyRune = 0
			case prevRune == '\\':
				replyRune = v
			default:
				return "", ErrInvalidString
			}
		case v == '\\':
			switch {
			case i == len(str)-1 && prevRune != '\\':
				return "", ErrInvalidString
			case replyRune != 0:
				sb.WriteRune(replyRune)
				replyRune = 0
			case prevRune == '\\':
				replyRune = v
			}
		default:
			if prevRune == '\\' {
				return "", ErrInvalidString
			}

			if replyRune != 0 {
				sb.WriteRune(replyRune)
			}

			replyRune = v
		}
	}

	if replyRune != 0 {
		sb.WriteRune(replyRune)
	}

	return sb.String(), nil
}
