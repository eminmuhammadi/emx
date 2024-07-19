package util

import (
	"net/url"
	"strings"
)

var specialChars = []string{"\"", "'", "`"}

func NormalizeQueryValue(value string) (string, error) {
	value, err := url.QueryUnescape(value)
	if err != nil {
		return "", err
	}

	for _, char := range specialChars {
		value = strings.TrimPrefix(value, char)
	}

	for _, char := range specialChars {
		value = strings.TrimSuffix(value, char)
	}

	return value, nil
}
