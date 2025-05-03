package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	app, err := InitializeApp()
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(app.Middleware.AddSessionCookie)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(app.Middleware.Logging)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(45 * time.Minute))

	if err := http.ListenAndServe(":8080", app.RegisterAPIRoutes(r)); err != nil {
		panic(err)
	}
}
