package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eve-an/estimated/internal/session"
)

func (app *Application) EventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	w.Header().Set("Access-Control-Allow-Origin", "*") // change for better in prod
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	flusher, ok := w.(http.Flusher)
	if !ok {
		app.logger.Error("Streaming unsupported!")
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	sess := app.getSessionData(r)
	if sess == nil {
		app.logger.Error("session is empty", "request", fmt.Sprintf("%+v", r))
		return
	}

	token := app.getToken(r)

	app.logger.Info("starting event", "token", token)

	for {
		select {
		case <-r.Context().Done():
			app.logger.Info("ending event", "token", token)
			return
		case <-app.sessions.Updater:
			helper := struct {
				Name    string              `json:"name"`
				Votings []session.VoteEntry `json:"points"`
			}{
				Name:    token,
				Votings: app.sessions.GetAllVotings(),
			}

			data, err := json.Marshal(helper)
			if err != nil {
				app.logger.Error("could not marshal session votes", "err", err.Error(), "session", sess.Token)
				continue
			}

			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				app.logger.Error("could not write event to client", "err", err.Error(), "session", sess.Token)
				continue
			}

			flusher.Flush()
		}
	}
}
