package handler

import (
	"context"

	"deviceback/v3/internal/service"
	"deviceback/v3/pkg/log"
	"deviceback/v3/pkg/util"
	"deviceback/v3/pkg/ws"

	"github.com/gin-gonic/gin"
)

func NewTestHandler(testService *service.TestService, logger log.Logger) *TestHandler {
	return &TestHandler{
		ts:  testService,
		log: log.NewHelper(logger),
	}
}

type TestHandler struct {
	ts  *service.TestService
	log *log.Helper
}

func (h *TestHandler) Test(c *gin.Context) {
	err := h.ts.TestXxx(context.Background())
	if err != nil {
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// ============================================================================
// ===========================  ws  ===========================================
// ============================================================================

func (h *TestHandler) Heartbeat(userId string, msg []byte) []byte {
	return util.MakeWsResp(200, "心跳", nil)
}

func (h *TestHandler) Ping(userId string, message []byte) []byte {
	return util.MakeWsResp(200, "pong", nil)
}

func (h *TestHandler) ConnPush(userId string, message []byte) []byte {
	ws.PushMsg(userId, "xxx", message)
	return util.MakeWsResp(200, "", nil)
}