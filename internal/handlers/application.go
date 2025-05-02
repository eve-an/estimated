package handlers

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/session"
	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Routes() http.Handler
}

type Application struct {
	logger   *slog.Logger
	sessions *session.SessionStore

	submitHandler   Handler
	registerHandler Handler
}

func NewApplication(
	logger *slog.Logger,
	sessions *session.SessionStore,
	submitHandler Handler,
	registerHandler Handler,
) *Application {
	return &Application{
		logger:   logger,
		sessions: sessions,

		submitHandler:   submitHandler,
		registerHandler: registerHandler,
	}
}

func (app *Application) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/submit", app.submitHandler.Routes())
		r.Mount("/register", app.registerHandler.Routes())
	})
}
