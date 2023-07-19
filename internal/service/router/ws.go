package router

import (
	"fmt"
	"net/http"
	"time"

	"deviceback/v3/internal/handler"
	"deviceback/v3/pkg/log"
	"deviceback/v3/pkg/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsRouter struct {
	router *gin.Engine
	cm     *ws.ClientManager
	log    *log.Helper

	wh *handler.TestHandler
}

func NewWsRouter(engine *gin.Engine, logger log.Logger) *WsRouter {
	r := &WsRouter{
		router: engine,
		cm:     ws.GetClientManagerInstance(),
		log:    log.NewHelper(logger),
	}
	r.initRouter()
	go r.cm.Start()

	return r
}

func (r *WsRouter) initRouter() {
	// ws事件
	ws.Register("conn.push", r.wh.ConnPush)
	ws.Register("heartbeat", r.wh.Heartbeat)
	ws.Register("ping", r.wh.Ping)

	// ws路由
	r.router.Group("./")
	{
		r.router.GET("/ws", r.serveWs)
	}
}

func (r *WsRouter) serveWs(c *gin.Context) {

	upgrader := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

	// 升级协议
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	fmt.Println("webSocket 建立连接:", conn.RemoteAddr().String())

	currentTime := uint64(time.Now().Unix())
	client := ws.NewClient(conn.RemoteAddr().String(), conn, currentTime)

	go client.Read()
	go client.Write()

	// 用户连接事件
	r.cm.Register <- client
}
