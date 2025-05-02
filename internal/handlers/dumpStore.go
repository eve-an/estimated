package handlers

import (
	"net/http"
)

func (app *Application) DumpStoreHandler(w http.ResponseWriter, r *http.Request) {
	votes := app.sessions.GetAllVotings()
	app.jsonIndent(w, http.StatusOK, votes)
}
