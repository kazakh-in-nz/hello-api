package translation

import (
	"strings"
)

type StaticService struct{}

func NewStaticService() *StaticService {
	return &StaticService{}
}

func (s *StaticService) Translate(word string, language string) string {
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
	case "russian":
		return "привет"
	default:
		return ""
	}
}

func sanitize(word string) string {
	w := strings.TrimSpace(word)
	return strings.ToLower(w)
}
