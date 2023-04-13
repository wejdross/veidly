package helpers

import (
	"strings"
)

func ExtractLanguageFromLocale(lc string) (string, error) {
	parts := strings.Split(lc, "_")
	if len(parts) != 2 || parts[0] == "" {
		return "en", nil
	}
	return parts[0], nil
}
