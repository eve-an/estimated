package handlers

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/httpx"
	"github.com/eve-an/estimated/internal/session"
	"github.com/go-chi/chi/v5"
)

type SessionStore interface {
	Push(vote session.VoteEntry) error
}

type submitHandler struct {
	logger *slog.Logger
	store  SessionStore
}

func newSubmitHandler(logger *slog.Logger, store SessionStore) *submitHandler {
	return &submitHandler{
		logger: logger,
		store:  store,
	}
}

func (s *submitHandler) HandleAdd(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.ReadRequestBody(r)
	if err != nil {
		s.logger.Error("reading body failed", "err", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "failed to read body",
		})
		return
	}

	var votes []session.VoteEntry
	if err := httpx.ParseJSON(body, &votes); err != nil {
		s.logger.Warn("invalid request body", "err", err)
		httpx.WriteJSON(w, http.StatusBadRequest, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "invalid JSON format",
		})
		return
	}

	for _, vote := range votes {
		if err := s.store.Push(vote); err != nil {
			s.logger.Error("failed to push vote", "err", err)
			httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
				Status: httpx.StatusError,
				Error:  "failed to store vote",
			})
			return
		}
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.APIResponse{
		Status: httpx.StatusSuccess,
		Data:   votes,
	})
}

func (s *submitHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/add", s.HandleAdd)

	return r
}
