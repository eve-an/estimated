package api

import (
	"encoding/json"
	"net/http"

	"github.com/eve-an/estimated/internal/api/dto"
)

func WriteJSON(w http.ResponseWriter, status int, resp dto.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
