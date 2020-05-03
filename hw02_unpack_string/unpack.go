package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inp string) (string, error) {
	inputReader := strings.NewReader(inp)

	return UnpackFromReader(inputReader)
}

func UnpackFromReader(inputReader *strings.Reader) (string, error) {
	var result strings.Builder
	for {
		current, _, err := inputReader.ReadRune()
		if err == io.EOF {
			break
		}
		if unicode.IsDigit(current) {
			return "", ErrInvalidString
		}

		next, _, err := inputReader.ReadRune()
		if err == io.EOF {
			result.WriteRune(current)
			break
		}

		if unicode.IsDigit(next) {
			count, _ := strconv.Atoi(string(next))
			result.WriteString(strings.Repeat(string(current), count))
		} else {
			err = inputReader.UnreadRune()
			if err != nil {
				return "", fmt.Errorf("unread rune error: %w", err)
			}

			result.WriteRune(current)
		}
	}

	return result.String(), nil
}
