//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"template/internal/data"
	"template/internal/handler"
	"template/internal/router"
	"template/internal/server"
	"template/internal/service"
	"template/pkg/app"
	"template/pkg/config"
	"template/pkg/log"

	"github.com/google/wire"
)

func wireApp(*config.Config, log.Logger) (*app.App, func(), error) {
	panic(wire.Build(server.Provider, service.Provider, router.Provider, handler.Provider, data.Provider, newApp))
}
