package handlers

import (
	"fmt"
	"net/http"
)

func (app *Application) ClearHandler(w http.ResponseWriter, r *http.Request) {
	totalDeleted := app.sessions.DeleteAll()

	app.json(w, http.StatusOK, map[string]string{
		"message": fmt.Sprintf("deleted %d", totalDeleted),
	})
}
