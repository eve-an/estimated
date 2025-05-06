package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/api"
	"github.com/eve-an/estimated/internal/api/dto"
	"github.com/eve-an/estimated/internal/infra/session"
	"github.com/eve-an/estimated/internal/service"
	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	logger *slog.Logger

	eventService service.EventService
}

func NewEventHandler(
	logger *slog.Logger,
	eventService service.EventService,
) *EventHandler {
	return &EventHandler{
		logger:       logger,
		eventService: eventService,
	}
}

func (e *EventHandler) setSSEHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*") // change for better in prod
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func (e *EventHandler) EventHandler(w http.ResponseWriter, r *http.Request) {
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

		response := dto.APIResponse{
			Status: dto.StatusSuccess,
			Data:   responseDTO,
		}

		data, err := json.Marshal(response)
		if err != nil {
			e.logger.Error("error while marshalling event response", "err", err, "response_dto", responseDTO)
			continue
		}

		if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
			e.logger.Error("could not write event to client", "err", err, "session_key", token)
			continue
		}

		w.(http.Flusher).Flush()
	}

	api.WriteJSON(w, http.StatusOK, dto.APIResponse{
		Status: dto.StatusSuccess,
		Data:   "ok",
	})
}

func (e *EventHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", e.EventHandler)

	return r
}
