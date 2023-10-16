package service

import "github.com/google/wire"

var Provider = wire.NewSet(NewTestService)
