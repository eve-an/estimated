package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/eve-an/estimated/internal/session"
)

type requestData struct{}

func (app *Application) SubmitHandler(w http.ResponseWriter, r *http.Request) {
	sessionData := app.getSessionData(r)
	if sessionData == nil {
		app.logger.Warn("no session available when submitting")
		app.json(w, http.StatusBadRequest, map[string]string{"message": "not registered"})
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		app.logger.Error("reading request body failed", "err", err.Error())
		app.json(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	var votes []session.Voting
	if err := json.Unmarshal(data, &votes); err != nil {
		app.logger.Error("decoding voting response failed", "err", err.Error())
		app.json(w, http.StatusInternalServerError, map[string]string{"message": err.Error()})
		return
	}

	sessionData.Push(votes)

	app.json(w, http.StatusOK, map[string]string{"message": string(data)})
}
