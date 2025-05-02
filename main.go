package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/eve-an/estimated/internal/handlers"
	appMiddleware "github.com/eve-an/estimated/internal/middleware"
	"github.com/eve-an/estimated/internal/session"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	sessionStore := session.NewSessionStore()
	middlewares := appMiddleware.NewMiddleware(logger, sessionStore)

	app := handlers.NewApplication(
		logger,
		sessionStore,
	)

	r := chi.NewRouter()
	r.Use(middlewares.AddSessionCookie)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middlewares.Logging)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(45 * time.Minute))

	r.Post("/register", app.RegisterHandler)
	r.Post("/submit", app.SubmitHandler)
	r.Get("/events", app.EventHandler)
	r.Delete("/clear", app.ClearHandler)
	r.Get("/dump", app.DumpStoreHandler)

	logger.Info("starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
