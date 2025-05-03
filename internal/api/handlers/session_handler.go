package handlers

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/api"
	"github.com/eve-an/estimated/internal/api/dto"
	"github.com/eve-an/estimated/internal/infra/session"
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
	key, err := session.FromContext(r.Context())
	if err != nil {
		rh.logger.Error("could not found assosiated session", "err", err)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "could not found assosiated session",
		})
		return
	}

	api.WriteJSON(w, http.StatusOK, dto.APIResponse{
		Status: dto.StatusSuccess,
		Data:   key,
	})
}

func (rh *sessionHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", rh.Create)

	return r
}
