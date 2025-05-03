package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Routes() http.Handler
}

type Application struct {
	votesHandler   *votesHandler
	sessionHandler *sessionHandler
	eventHandler   *eventHandler
}

func NewApplication(
	votesHandler *votesHandler,
	sessionHandler *sessionHandler,
	eventHandler *eventHandler,
) *Application {
	return &Application{
		votesHandler:   votesHandler,
		sessionHandler: sessionHandler,
		eventHandler:   eventHandler,
	}
}

func (app *Application) RegisterAPIRoutes(r chi.Router) http.Handler {
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/votes", app.votesHandler.Routes())
		r.Mount("/register", app.sessionHandler.Routes())
		r.Mount("/events", app.eventHandler.Routes())
	})

	return r
}
