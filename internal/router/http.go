package router

import (
	"deviceback/v3/internal/handler"

	"github.com/gin-gonic/gin"
)

type HttpRouter struct {
	router *gin.Engine
	th     *handler.TestHandler
}

func NewHttpRouter(engine *gin.Engine, testHandler *handler.TestHandler) *HttpRouter {
	r := &HttpRouter{
		router: engine,
		th:     testHandler,
	}
	r.initRouter()

	return r
}

func (r *HttpRouter) initRouter() {
	r.router.GET("/test", r.th.Test)
}
