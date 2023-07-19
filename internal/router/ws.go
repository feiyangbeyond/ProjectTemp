package router

import (
	"net/http"
	"time"

	"deviceback/v3/internal/handler"
	"deviceback/v3/internal/model"
	"deviceback/v3/pkg/log"
	"deviceback/v3/pkg/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsRouter struct {
	router *gin.Engine
	log    *log.Helper

	wh *handler.TestHandler
}

func NewWsRouter(engine *gin.Engine, logger log.Logger) *WsRouter {
	r := &WsRouter{
		router: engine,
		log:    log.NewHelper(logger),
	}
	r.initRouter()
	go ws.StartClientManager()
	go ws.ClearTimeoutConnections()

	return r
}

func (r *WsRouter) initRouter() {
	// ws事件
	ws.Register(model.EventConnPush, r.wh.ConnPush)

	// ws路由
	r.router.Group("./")
	{
		r.router.GET("/ws", r.serveWs)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (r *WsRouter) serveWs(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		http.NotFound(c.Writer, c.Request)
		return
	}

	oldClient := ws.GetUserClient(userId)
	if oldClient != nil {
		// 通知旧客户端下线
		oldClient.LogoutWithMsg()
	}

	// 升级协议
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	r.log.Info("webSocket 建立连接:", conn.RemoteAddr().String())
	currentTime := uint64(time.Now().Unix())
	client := ws.NewClient(conn.RemoteAddr().String(), conn, currentTime, userId)

	go client.Read()
	go client.Write()

	// 用户连接事件
	ws.ClientRegister(client)
}
