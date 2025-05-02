package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/middleware"
	"github.com/eve-an/estimated/internal/session"
)

type Application struct {
	logger *slog.Logger

	sessions *session.SessionStore
}

func NewApplication(
	logger *slog.Logger,
	sessions *session.SessionStore,
) *Application {
	return &Application{
		logger,
		sessions,
	}
}

func (app *Application) json(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.logger.Error("failed to write JSON", "error", err)
	}
}

func (app *Application) jsonIndent(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	value, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		app.logger.Error("failed to write JSON", "error", err)
		return
	}

	fmt.Fprint(w, string(value))
}

func (app *Application) getToken(r *http.Request) string {
	cookie, err := r.Cookie(middleware.ClientTokenName)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func (app *Application) getSessionData(r *http.Request) *session.SessionData {
	token := app.getToken(r)
	if token == "" {
		return nil
	}

	return app.sessions.Get(token)
}
