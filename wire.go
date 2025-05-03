//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/eve-an/estimated/internal/api"
	"github.com/eve-an/estimated/internal/api/handlers"
	internalMiddleware "github.com/eve-an/estimated/internal/api/middleware"
	"github.com/eve-an/estimated/internal/config"
	"github.com/eve-an/estimated/internal/infra/notify"
	"github.com/eve-an/estimated/internal/infra/store"
	"github.com/eve-an/estimated/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/wire"
)

var SingletonSet = wire.NewSet(
	provideLogger,
)

var StoreSet = wire.NewSet(
	store.NewSessionStore,
	wire.Bind(new(service.VoteStore), new(*store.SessionStore)),
)

var HandlerSet = wire.NewSet(
	handlers.NewVotesHandler,
	handlers.NewSessionHandler,
	handlers.NewEventHandler,
	handlers.NewApplication,
)

var HTTPSet = wire.NewSet(
	provideHTTPServer,
	provideRouter,
	api.NewServer,
	internalMiddleware.NewMiddleware,
)

var ServiceSet = wire.NewSet(
	service.NewVoteService,
	service.NewEventService,
)

var NotifierSet = wire.NewSet(
	notify.NewSessionNotifier,
	wire.Bind(new(notify.Notifier), new(*notify.SessionNotifier)),
)

func provideLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func provideHTTPServer(config *config.Config) *http.Server {
	return &http.Server{
		Addr:         config.ServerAddress,
		ReadTimeout:  time.Duration(config.ServerTimeout),
		WriteTimeout: time.Duration(config.ServerTimeout),
		IdleTimeout:  120 * time.Second,
	}
}

func provideRouter(
	app *handlers.Application,
	mw *internalMiddleware.Middleware,
	config *config.Config,
	logger *slog.Logger,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.Logging)
	r.Use(middleware.Recoverer)
	r.Use(mw.AddSessionCookie)
	r.Use(middleware.Timeout(time.Duration(config.ServerTimeout) * time.Second))

	return app.RegisterAPIRoutes(r)
}

func InitializeApp(config *config.Config) (*api.Server, error) {
	wire.Build(
		// Core dependencies
		SingletonSet,

		// Data stores
		StoreSet,

		// Services
		ServiceSet,

		// Notifiers
		NotifierSet,

		// HTTP components
		HTTPSet,

		// Handlers
		HandlerSet,
	)

	return nil, nil
}
