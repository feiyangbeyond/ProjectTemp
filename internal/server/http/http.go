package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"template/internal/router"
	"template/pkg/config"
	"template/pkg/log"
	"template/pkg/middeware/cors"

	"github.com/gin-gonic/gin"
)

// NewEngine creates a new instance of gin.Engine.
// It sets up the necessary middlewares and returns the engine.
func NewEngine(conf *config.Config, logger log.Logger) *gin.Engine {
	// Create a new gin engine
	r := gin.New()

	// Use gin recovery middleware to recover from any panics during request processing
	r.Use(gin.Recovery())

	// Use gin logger middleware to log every HTTP request
	r.Use(gin.Logger())

	// Use CORS middleware to handle Cross-Origin Resource Sharing
	r.Use(cors.Cors())

	// Return the gin engine
	return r
}

type Server struct {
	*http.Server
	lis     net.Listener
	err     error
	network string
	address string
	ws      *router.WsRouter
	http    *router.HttpRouter
}

// NewServer creates a new Server instance.
// It takes a config.Config, *gin.Engine, router.WsRouter, and router.HttpRouter as parameters.
// It returns a pointer to the Server instance.
func NewServer(config *config.Config, engine *gin.Engine, wsr *router.WsRouter, httpr *router.HttpRouter) *Server {
	s := &Server{
		network: "tcp",
		address: ":0",
		ws:      wsr,
		http:    httpr,
	}

	// If the Server Port is specified in the config, update the address accordingly.
	if config.Server.Port != 0 {
		s.address = fmt.Sprintf(":%d", config.Server.Port)
	}

	// Create a new http.Server instance with the provided engine as the handler.
	// Set the ReadTimeout, WriteTimeout, and MaxHeaderBytes values.
	s.Server = &http.Server{
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s
}

// Start starts the HTTP server.
func (s *Server) Start(ctx context.Context) error {
	// Listen for network connections
	lis, err := net.Listen(s.network, s.address)
	if err != nil {
		s.err = err
		return err
	}
	s.lis = lis

	// Set the base context for the server
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}

	// Log the server address
	log.Infof("[HTTP] server listening on: %s", s.lis.Addr().String())

	// Serve incoming requests
	if err := s.Serve(s.lis); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop stops the HTTP server gracefully.
// It shuts down the server by calling the Shutdown method with the provided context.
func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}
