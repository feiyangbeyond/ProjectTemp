package handler

import (
	"context"

	"template/internal/service"
	"template/pkg/log"
	"template/pkg/util"
	"template/pkg/ws"

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

func (h *TestHandler) ConnPush(userId string, message []byte) []byte {
	ws.PushMsg(userId, "push.msg", message)
	ws.PushMsgAll("push.msg.all", message)

	return util.MakeWsResp(200, "", nil)
}
