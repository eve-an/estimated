//go:build wireinject
// +build wireinject

package di

import (
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/eve-an/estimated/internal/api/handlers"
	"github.com/eve-an/estimated/internal/api/middleware"
	"github.com/eve-an/estimated/internal/config"
	"github.com/eve-an/estimated/internal/infra/notify"
	"github.com/eve-an/estimated/internal/infra/store"
	"github.com/eve-an/estimated/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/wire"
)

var SingletonSet = wire.NewSet(
	ProvideLogger,
)

var StoreSet = wire.NewSet(
	store.NewSessionStore,
	wire.Bind(new(service.VoteStore), new(*store.SessionStore)),
)

var HandlerSet = wire.NewSet(
	handlers.NewVotesHandler,
	handlers.NewSessionHandler,
	handlers.NewEventHandler,
)

var HTTPSet = wire.NewSet(
	ProvideRouter,
	middleware.NewMiddleware,
)

var ServiceSet = wire.NewSet(
	service.NewVoteService,
	service.NewEventService,
)

var NotifierSet = wire.NewSet(
	notify.NewSessionNotifier,
	wire.Bind(new(notify.Notifier), new(*notify.SessionNotifier)),
)

var ApplicationSet = wire.NewSet(
	SingletonSet,
	StoreSet,
	ServiceSet,
	NotifierSet,
	HTTPSet,
	HandlerSet,
)

func ProvideLogger() *slog.Logger {
	levelStr := strings.ToUpper(os.Getenv("LOG_LEVEL"))

	debugLevel := slog.LevelInfo
	switch levelStr {
	case "DEBUG":
		debugLevel = slog.LevelDebug
	case "WARN":
		debugLevel = slog.LevelWarn
	case "ERROR":
		debugLevel = slog.LevelError
	default:
		debugLevel = slog.LevelInfo
	}

	s := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: debugLevel,
	}))
	slog.SetDefault(s)

	return s
}

func ProvideRouter(
	votesHandler *handlers.VotesHandler,
	sessionHandler *handlers.SessionHandler,
	eventHandler *handlers.EventHandler,
	mw *middleware.Middleware,
	config *config.Config,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(mw.Logging)
	r.Use(chiMiddleware.Recoverer)
	r.Use(mw.AddSessionCookie)
	r.Use(chiMiddleware.Timeout(time.Duration(config.ServerTimeout) * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/votes", votesHandler.Routes())
		r.Mount("/register", sessionHandler.Routes())
		r.Mount("/events", eventHandler.Routes())
	})

	return r
}

func InitializeApp(config *config.Config) (http.Handler, error) {
	wire.Build(ApplicationSet)

	return nil, nil
}
