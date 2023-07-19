package service

import "github.com/google/wire"

var ServiceProvider = wire.NewSet(NewTestService)
