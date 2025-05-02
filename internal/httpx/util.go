package httpx

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var ErrInvalidJSON = errors.New("invalid JSON body")

func WriteJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func ReadRequestBody(r *http.Request) ([]byte, error) {
	return io.ReadAll(r.Body)
}

func ParseJSON(body []byte, dst any) error {
	if err := json.Unmarshal(body, dst); err != nil {
		return ErrInvalidJSON
	}
	return nil
}
