//go:build wireinject
// +build wireinject

package main

import (
	"log/slog"
	"os"

	"github.com/eve-an/estimated/internal/db"
	"github.com/eve-an/estimated/internal/handlers"
	"github.com/eve-an/estimated/internal/middleware"
	"github.com/eve-an/estimated/internal/session"
	"github.com/google/wire"
)

var voteEntryStoreSet = wire.NewSet(
	session.NewSessionStore,
	wire.Bind(new(db.VoteEntryStore), new(*session.SessionStore)),
)

func ProvideLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func InitializeApp() (*handlers.Application, error) {
	wire.Build(
		handlers.NewApplication,

		handlers.NewSessionHandler,
		handlers.NewVotesHandler,
		ProvideLogger,
		voteEntryStoreSet,
		middleware.NewMiddleware,
	)

	return nil, nil
}
