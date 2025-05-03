package handlers

import (
	"net/http"

	"github.com/eve-an/estimated/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Routes() http.Handler
}

type Application struct {
	votesHandler   *votesHandler
	sessionHandler *sessionHandler
	Middleware     *middleware.Middleware
}

func NewApplication(
	votesHandler *votesHandler,
	sessionHandler *sessionHandler,
	middleware *middleware.Middleware,
) *Application {
	return &Application{
		votesHandler:   votesHandler,
		sessionHandler: sessionHandler,
		Middleware:     middleware,
	}
}

func (app *Application) RegisterAPIRoutes(r chi.Router) http.Handler {
	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/votes", app.votesHandler.Routes())
		r.Mount("/register", app.sessionHandler.Routes())
	})

	return r
}
