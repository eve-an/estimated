package handlers

import (
	"net/http"
)

func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	app.json(w, http.StatusOK, map[string]string{"name": app.getToken(r)})
}
