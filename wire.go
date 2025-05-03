//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"
	"os"

	"github.com/eve-an/estimated/internal/config"
	"github.com/eve-an/estimated/internal/db"
	"github.com/eve-an/estimated/internal/handlers"
	"github.com/eve-an/estimated/internal/middleware"
	"github.com/eve-an/estimated/internal/notify"
	"github.com/eve-an/estimated/internal/session"
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

func provideConfig() (*config.Config, error) {
	return config.LoadConfig("config.json")
}

func InitializeApp() (*handlers.Application, error) {
	wire.Build(
		handlers.NewApplication,

		handlers.NewSessionHandler,
		handlers.NewVotesHandler,
		handlers.NewEventHandler,

		provideConfig,
		provideLogger,

		middleware.NewMiddleware,

		newSessionNotifierStoreSet,
		voteEntryStoreSet,
	)

	return nil, nil
}
