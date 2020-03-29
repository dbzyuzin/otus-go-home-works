package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inp string) (string, error) {
	inputReader := strings.NewReader(inp)

	var result strings.Builder
	for inputReader.Len() > 0 {
		current, _, _ := inputReader.ReadRune()
		next, _, _ := inputReader.ReadRune()

		if unicode.IsDigit(current) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(next) {
			count, _ := strconv.Atoi(string(next))
			result.WriteString(strings.Repeat(string(current), count))
		} else {
			_ = inputReader.UnreadRune()
			result.WriteRune(current)
		}
	}

	return result.String(), nil
}
