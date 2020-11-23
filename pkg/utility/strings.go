package utility

import (
	"crypto/rand"
	"fmt"
	"strings"
	"unicode"
)

// GenerateToken generates a random string of len length
func GenerateToken(len int) (string, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

// TrimWhitespace takes a string and removes any spaces
func TrimWhitespace(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if !unicode.IsSpace(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
