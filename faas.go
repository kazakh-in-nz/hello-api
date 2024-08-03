package faas

import (
	"net/http"

	"github.com/kazakh-in-nz/hello-api/handlers/rest"
	"github.com/kazakh-in-nz/hello-api/translation"
)

func Translate(w http.ResponseWriter, r *http.Request) {
	translationSvc := translation.NewStaticService()
	handler := rest.NewTranslatorHandler(translationSvc)

	handler.TranslateHandler(w, r)
}
