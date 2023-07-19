package router

import (
	"github.com/google/wire"
)

var RouterProvider = wire.NewSet(NewHttpRouter, NewWsRouter)
