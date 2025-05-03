package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/db"
	"github.com/eve-an/estimated/internal/httpx"
	"github.com/eve-an/estimated/internal/model"
	"github.com/eve-an/estimated/internal/notify"
	"github.com/go-chi/chi/v5"
)

type eventHandler struct {
	logger *slog.Logger

	store           db.VoteEntryStore
	sessionNotifier *notify.SessionNotifier
}

func NewEventHandler(
	logger *slog.Logger,
	store db.VoteEntryStore,
	sessionNotifier *notify.SessionNotifier,
) *eventHandler {
	return &eventHandler{
		logger:          logger,
		store:           store,
		sessionNotifier: sessionNotifier,
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

	flusher, ok := w.(http.Flusher)
	if !ok {
		e.logger.Error("Streaming unsupported!")
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	token, err := httpx.SessionKeyFromContext(r.Context())
	if err != nil {
		e.logger.Error("could not found assosiated session", "err", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "could not found assosiated session",
		})
		return
	}

	notification := e.sessionNotifier.Subscribe(token)

Loop:
	for {
		select {
		case <-r.Context().Done():
			break Loop
		case <-notification:

			votes, err := e.store.List()
			if err != nil {
				e.logger.Error("could not extract votes from store", "err", err)
				continue
			}

			helper := struct {
				Name    string            `json:"name"`
				Votings []model.VoteEntry `json:"points"`
			}{
				Name:    token,
				Votings: votes,
			}

			data, err := json.Marshal(helper)
			if err != nil {
				e.logger.Error("could not marshal session votes", "err", err, "token", token)
				continue
			}

			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				e.logger.Error("could not write event to client", "err", err, "token", token)
				continue
			}

			flusher.Flush()
		}
	}

	if err := e.sessionNotifier.Unsubscribe(token); err != nil {
		e.logger.Error("could not find assosiated session", "err", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "could not find assosiated session",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.APIResponse{
		Status: httpx.StatusSuccess,
		Data:   "ok",
	})
}

func (e *eventHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", e.EventHandler)

	return r
}
