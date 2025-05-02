package handlers

import (
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type registerHandler struct {
	logger *slog.Logger
}

func newRegisterHandler(logger *slog.Logger) *registerHandler {
	return &registerHandler{
		logger: logger,
	}
}

func (rh *registerHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	key, err := httpx.SessionKeyFromContext(r.Context())
	if err != nil {
		rh.logger.Error("could not found assosiated session", "err", err)
		httpx.WriteJSON(w, http.StatusInternalServerError, httpx.APIResponse{
			Status: httpx.StatusError,
			Data:   "could not found assosiated session",
		})
		return
	}

	httpx.WriteJSON(w, http.StatusOK, httpx.APIResponse{
		Status: httpx.StatusSuccess,
		Data:   key,
	})
}

func (rh *registerHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/register", rh.HandleRegister)

	return r
}
