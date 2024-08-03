package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kazakh-in-nz/hello-api/handlers"
	"github.com/kazakh-in-nz/hello-api/handlers/rest"
	"github.com/kazakh-in-nz/hello-api/translation"
)

func main() {
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if addr == ":" {
		addr = ":8080"
	}

	mux := http.NewServeMux()

	translationSvc := translation.NewStaticService()
	translateHandler := rest.NewTranslatorHandler(translationSvc)
	mux.HandleFunc("/hello", translateHandler.TranslateHandler)
	mux.HandleFunc("/health", handlers.HealthCheck)

	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

type Resp struct {
	Language    string `json:"language"`
	Translation string `json:"translation"`
}
