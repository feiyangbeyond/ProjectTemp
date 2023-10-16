package router

import (
	"github.com/google/wire"
)

var Provider = wire.NewSet(NewHttpRouter, NewWsRouter)
