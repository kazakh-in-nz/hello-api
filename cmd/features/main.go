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

	mux := API(cfg)

	log.Printf("listening on %s\n", addr)

	log.Fatal(http.ListenAndServe(addr, mux))
}

func API(cfg config.Configuration) *http.ServeMux {
	mux := http.NewServeMux()

	var translationService rest.Translator

	if cfg.DatabaseURL != "" {
		db := translation.NewDatabaseService(cfg)
		translationService = db
	} else {
		translationService = translation.NewStaticService()
	}

	if cfg.LegacyEndpoint != "" {
		log.Printf("creating external translation client: %s", cfg.LegacyEndpoint)
		client := translation.NewHelloClient(cfg.LegacyEndpoint)
		translationService = translation.NewRemoteService(client)
	}

	translateHandler := rest.NewTranslatorHandler(translationService)

	mux.HandleFunc("/hello", translateHandler.TranslateHandler)
	mux.HandleFunc("/health", handlers.HealthCheck)

	return mux
}
