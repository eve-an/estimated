package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/config"
	"github.com/eve-an/estimated/internal/di"
)

func main() {
	config, err := config.LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	handler, err := di.InitializeApp(config)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:    config.ServerAddress,
		Handler: handler,
	}

	slog.Info("starting server", "config", config)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
