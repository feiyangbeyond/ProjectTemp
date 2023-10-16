package server

import (
	"context"
	"net/url"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type Endpointer interface {
	Endpoint() (*url.URL, error)
}
