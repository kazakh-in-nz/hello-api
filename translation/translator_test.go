//go:build unit

package translation_test

import (
	"testing"

	"github.com/kazakh-in-nz/hello-api/translation"
)

func TestTranslate(t *testing.T) {
	tt := []struct {
		desc     string
		word     string
		language string
		expected string
	}{
		{"should pass if english language with known word", "hello", "english", "hello"},
		{"should pass if german language with known word", "hello", "german", "hallo"},
		{"should pass if french language with known word", "hello", "french", "bonjour"},
		{"should return empty string if uknown word with known language", "bye", "german", ""},
		{"should return translated word if word is of title case", "Hello", "german", "hallo"},
		{"should return translated word if word is of upper case", "HELLO", "german", "hallo"},
		{"should return translated word if language is of title case", "hello", "German", "hallo"},
		{"should return translated word if language is of upper case", "hello", "GERMAN", "hallo"},
		{"should pass if finnish language with known word", "hello", "finnish", "hei"},
		{"should return empty string if unknown language", "hello", "dutch", ""},
		{"should return empty string if unknown word and unkown language", "bye", "dutch", ""},
	}

	underTest := translation.NewStaticService()

	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			res := underTest.Translate(tc.word, tc.language)

			if res != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, res)
			}
		})
	}
}
