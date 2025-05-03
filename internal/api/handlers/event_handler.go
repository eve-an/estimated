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

type eventHandler struct {
	logger *slog.Logger

	eventService service.EventService
}

func NewEventHandler(
	logger *slog.Logger,
	eventService service.EventService,
) *eventHandler {
	return &eventHandler{
		logger: logger,
	}
}

func (e *eventHandler) setSSEHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*") // change for better in prod
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func (e *eventHandler) EventHandler(w http.ResponseWriter, r *http.Request) {
	e.setSSEHeader(w)

	token, err := session.FromContext(r.Context())
	if err != nil {
		e.logger.Error("could not found assosiated session", "err", err)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "could not found assosiated session",
		})
		return
	}

	c := e.eventService.Subscribe(r.Context(), token)

	for responseDTO := range c {
		if responseDTO == nil {
			break
		}

		api.WriteJSON(w, http.StatusOK, dto.APIResponse{
			Status: dto.StatusSuccess,
			Data:   responseDTO,
		})

		w.(http.Flusher).Flush()
	}

	api.WriteJSON(w, http.StatusOK, dto.APIResponse{
		Status: dto.StatusSuccess,
		Data:   "ok",
	})
}

func (e *eventHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", e.EventHandler)

	return r
}
