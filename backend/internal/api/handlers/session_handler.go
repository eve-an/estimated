package handlers

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/api"
	"github.com/eve-an/estimated/internal/api/dto"
	"github.com/eve-an/estimated/internal/infra/session"
	"github.com/eve-an/estimated/internal/service"
	"github.com/go-chi/chi/v5"
)

type SessionHandler struct {
	logger        *slog.Logger
	nameGenerator service.NameGenerator
}

func NewSessionHandler(
	logger *slog.Logger,
	nameGenerator service.NameGenerator,
) *SessionHandler {
	return &SessionHandler{
		logger:        logger,
		nameGenerator: nameGenerator,
	}
}

func (rh *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	key, err := session.FromContext(r.Context())
	if err != nil {
		rh.logger.Error("could not found assosiated session", "err", err)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "could not found assosiated session",
		})
		return
	}

	name, err := rh.nameGenerator.NameFor(key)
	if err != nil {
		rh.logger.Error("could not generate name", "err", err, "sesion_key", key)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "could not generate name",
		})
		return
	}

	api.WriteJSON(w, http.StatusOK, dto.APIResponse{
		Status: dto.StatusSuccess,
		Data:   name,
	})
}

func (rh *SessionHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", rh.Create)

	return r
}
