package handler

import "github.com/google/wire"

var Provider = wire.NewSet(NewTestHandler)
