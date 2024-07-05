package translation

import (
	"strings"
)

func Translate(word string, language string) string {
	if sanitize(word) != "hello" {
		return ""
	}

	switch sanitize(language) {
	case "english":
		return "hello"
	case "german":
		return "hallo"
	case "finnish":
		return "hei"
	case "french":
		return "bonjour"
	default:
		return ""
	}
}

func sanitize(word string) string {
	w := strings.TrimSpace(word)
	return strings.ToLower(w)
}
