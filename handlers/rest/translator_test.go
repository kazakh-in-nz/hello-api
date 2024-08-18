//go:build unit

package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kazakh-in-nz/hello-api/handlers/rest"
)

type stubbedSvc struct{}

func (s *stubbedSvc) Translate(word string, language string) string {
	if word == "foo" {
		return "bar"
	}

	return ""
}

func TestTranslateAPI(t *testing.T) {
	tt := []struct {
		name                string
		endpoint            string
		statusCode          int
		expectedLanguage    string
		expectedTranslation string
	}{
		{
			name:                "english translation",
			endpoint:            "/foo",
			statusCode:          http.StatusOK,
			expectedLanguage:    "english",
			expectedTranslation: "bar",
		},
		{
			name:                "german translation",
			endpoint:            "/foo?language=german",
			statusCode:          http.StatusOK,
			expectedLanguage:    "german",
			expectedTranslation: "bar",
		},
		{
			name:                "not found",
			endpoint:            "/hello?language=dutch",
			statusCode:          http.StatusNotFound,
			expectedLanguage:    "",
			expectedTranslation: "",
		},
	}

	underTest := rest.NewTranslatorHandler(&stubbedSvc{})
	handler := http.HandlerFunc(underTest.TranslateHandler)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.endpoint, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)
			if rec.Code != tc.statusCode {
				t.Errorf("expected status code %d, got %d", tc.statusCode, rec.Code)
			}

			var resp rest.Resp

			_ = json.Unmarshal(rec.Body.Bytes(), &resp)
			if resp.Language != tc.expectedLanguage {
				t.Errorf("expected language %s, got %s", tc.expectedLanguage, resp.Language)
			}

			if resp.Translation != tc.expectedTranslation {
				t.Errorf("expected translation %s, got %s", tc.expectedTranslation, resp.Translation)
			}
		})
	}
}
