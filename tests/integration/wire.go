//go:build wireinject
// +build wireinject

package integration

import (
	"net/http"

	"github.com/eve-an/estimated/internal/config"
	"github.com/eve-an/estimated/internal/di"
	"github.com/google/wire"
)

func InitializeTestApp(config *config.Config) (http.Handler, error) {
	wire.Build(di.ApplicationSet)
	return nil, nil
}
