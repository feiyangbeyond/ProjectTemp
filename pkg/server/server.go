package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"deviceback/v3/internal/router"
	"deviceback/v3/pkg/config"
	"deviceback/v3/pkg/middeware/cors"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ServerProvider = wire.NewSet(NewServer, NewEngine)

type Server struct {
	port int
	r    *gin.Engine
	s    *http.Server
	ws   *router.WsRouter
	http *router.HttpRouter
}

func NewServer(config *config.Config, engine *gin.Engine, ws *router.WsRouter, http *router.HttpRouter) *Server {
	s := &Server{
		r:    engine,
		port: config.Server.Port,
		ws:   ws,
		http: http,
	}

	return s
}

func (s *Server) Start() error {
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", s.port),
		Handler:        s.r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.s = server
	return s.s.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}

func NewEngine() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Cors())

	return r
}
