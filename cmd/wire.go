//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"deviceback/v3/internal/data"
	"deviceback/v3/internal/handler"
	"deviceback/v3/internal/router"
	"deviceback/v3/internal/service"
	"deviceback/v3/pkg/config"
	"deviceback/v3/pkg/log"
	"deviceback/v3/pkg/server"

	"github.com/google/wire"
)

func wireApp(*config.Config, log.Logger) (*App, func(), error) {
	panic(wire.Build(server.ServerProvider, service.ServiceProvider, router.RouterProvider, handler.HandlerProvider, data.DataProvider, newApp))
}
