//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/eve-an/estimated/internal/config"
	"github.com/eve-an/estimated/internal/db"
	"github.com/eve-an/estimated/internal/handlers"
	"github.com/eve-an/estimated/internal/httpx"
	internalMiddleware "github.com/eve-an/estimated/internal/middleware"
	"github.com/eve-an/estimated/internal/notify"
	"github.com/eve-an/estimated/internal/session"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/wire"
)

var voteEntryStoreSet = wire.NewSet(
	session.NewSessionStore,
	wire.Bind(new(db.VoteEntryStore), new(*session.SessionStore)),
)

var newSessionNotifierStoreSet = wire.NewSet(
	notify.NewSessionNotifier,
	wire.Bind(new(notify.Notifier), new(*notify.SessionNotifier)),
)

func provideLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func provideHTTPServer(config *config.Config) *http.Server {
	return &http.Server{
		Addr: config.ServerAddress,
	}
}

func provideMux(
	app *handlers.Application,
	mw *internalMiddleware.Middleware,
	config *config.Config,
) http.Handler {
	r := chi.NewRouter()
	r.Use(mw.AddSessionCookie)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Logging)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(time.Duration(config.ServerTimeout)))

	return app.RegisterAPIRoutes(r)
}

func InitializeApp(config *config.Config) (*httpx.Server, error) {
	wire.Build(
		provideHTTPServer,
		httpx.NewServer,
		handlers.NewApplication,

		handlers.NewSessionHandler,
		handlers.NewVotesHandler,
		handlers.NewEventHandler,

		provideMux,
		provideLogger,

		internalMiddleware.NewMiddleware,

		newSessionNotifierStoreSet,
		voteEntryStoreSet,
	)

	return nil, nil
}
