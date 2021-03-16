package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if !utf8.ValidString(str) {
		return "", ErrInvalidString
	}

	if unicode.IsDigit([]rune(str)[:1][0]) {
		return "", ErrInvalidString
	}

	var sb strings.Builder
	var replyRune rune

	for i, v := range str {
		prevRune, _ := utf8.DecodeLastRune([]byte(str[:i]))

		switch {
		case unicode.IsDigit(v):
			switch {
			case replyRune != 0:
				multiplier, _ := strconv.Atoi(string(v))

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
			if i == len(str)-1 && prevRune != '\\' {
				return "", ErrInvalidString
			}

			if replyRune != 0 {
				sb.WriteRune(replyRune)
				replyRune = 0
			} else if prevRune == '\\' {
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
