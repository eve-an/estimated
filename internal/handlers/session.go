package handlers

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type sessionHandler struct {
	logger *slog.Logger
}

func NewSessionHandler(logger *slog.Logger) *sessionHandler {
	return &sessionHandler{
		logger: logger,
	}
}

func (rh *sessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	key, err := httpx.SessionKeyFromContext(r.Context())
	if err != nil {
		rh.logger.Error("could not found assosiated session", "err", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "could not found assosiated session",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.APIResponse{
		Status: httpx.StatusSuccess,
		Data:   key,
	})
}

func (rh *sessionHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", rh.Create)

	return r
}
