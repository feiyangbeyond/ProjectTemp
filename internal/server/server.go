package server

import (
	"template/internal/server/http"

	"github.com/google/wire"
)

var Provider = wire.NewSet(http.NewServer, http.NewEngine)
