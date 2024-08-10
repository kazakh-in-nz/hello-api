package main

import (
	"log"
	"net/http"

	"github.com/kazakh-in-nz/hello-api/config"
	"github.com/kazakh-in-nz/hello-api/handlers"
	"github.com/kazakh-in-nz/hello-api/handlers/rest"
	"github.com/kazakh-in-nz/hello-api/translation"
)

func main() {
	cfg := config.LoadConfiguration()
	addr := cfg.Port

	mux := http.NewServeMux()

	var translationSvc rest.Translator
	translationSvc = translation.NewStaticService()
	if cfg.LegacyEndpoint != "" {
		log.Printf("creating external translation client: %s", cfg.LegacyEndpoint)
		client := translation.NewHelloClient(cfg.LegacyEndpoint)
		translationSvc = translation.NewRemoteService(client)
	}

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
