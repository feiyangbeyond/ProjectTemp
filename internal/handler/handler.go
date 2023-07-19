package handler

import "github.com/google/wire"

var HandlerProvider = wire.NewSet(NewTestHandler)
