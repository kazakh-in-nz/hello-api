// go:build unit

package rest_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kazakh-in-nz/hello-api/handlers/rest"
)

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
			endpoint:            "/hello",
			statusCode:          http.StatusOK,
			expectedLanguage:    "english",
			expectedTranslation: "hello",
		},
		{
			name:                "german translation",
			endpoint:            "/hello?language=german",
			statusCode:          http.StatusOK,
			expectedLanguage:    "german",
			expectedTranslation: "hallo",
		},
		{
			name:                "not found",
			endpoint:            "/hello?language=dutch",
			statusCode:          http.StatusNotFound,
			expectedLanguage:    "",
			expectedTranslation: "",
		},
	}

	handler := http.HandlerFunc(rest.TranslateHandler)

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.endpoint, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.statusCode {
				t.Errorf("expected status code %d, got %d", tc.statusCode, rec.Code)
			}

			var resp rest.Resp
			json.Unmarshal(rec.Body.Bytes(), &resp)

			if resp.Language != tc.expectedLanguage {
				t.Errorf("expected language %s, got %s", tc.expectedLanguage, resp.Language)
			}

			if resp.Translation != tc.expectedTranslation {
				t.Errorf("expected translation %s, got %s", tc.expectedTranslation, resp.Translation)
			}
		})
	}
}
