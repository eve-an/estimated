package handlers

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/db"
	"github.com/eve-an/estimated/internal/httpx"
	"github.com/eve-an/estimated/internal/model"
	"github.com/go-chi/chi/v5"
)

type votesHandler struct {
	logger *slog.Logger
	store  db.VoteEntryStore
}

func NewVotesHandler(logger *slog.Logger, store db.VoteEntryStore) *votesHandler {
	return &votesHandler{
		logger: logger,
		store:  store,
	}
}

func (s *votesHandler) Add(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.ReadRequestBody(r)
	if err != nil {
		s.logger.Error("reading body failed", "err", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "failed to read body",
		})
		return
	}

	var votes []model.VoteEntry
	if err := httpx.ParseJSON(body, &votes); err != nil {
		s.logger.Warn("invalid request body", "err", err)
		httpx.WriteJSON(w, http.StatusBadRequest, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "invalid JSON format",
		})
		return
	}

	key, err := httpx.SessionKeyFromContext(r.Context())
	if err != nil {
		s.logger.Warn("client has no session key", "err", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "no session found",
		})
		return
	}

	for _, vote := range votes {
		if err := s.store.Add(key, vote); err != nil {
			s.logger.Error("failed to add vote", "err", err)
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

func (s *votesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	votes, err := s.store.List()
	if err != nil {
		s.logger.Error("could not fetch all votes from the store", "err", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "could not fetch all votes from the store",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.APIResponse{
		Status: httpx.StatusSuccess,
		Data:   votes,
	})
}

func (s *votesHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	n, err := s.store.Clear()
	if err != nil {
		s.logger.Error("could not delete all votes from the store", "err", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
			Status: httpx.StatusError,
			Error:  "could not delete all votes from store",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.APIResponse{
		Status: httpx.StatusSuccess,
		Data:   n,
	})
}

func (s *votesHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", s.Add)
	r.Get("/", s.GetAll)
	r.Delete("/", s.DeleteAll)

	return r
}
