package ws

import (
	"encoding/json"
	"sync"

	"deviceback/v3/internal/model"
)

type DisposeFunc func(userId string, message []byte) (data []byte)

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex
)

// Register 注册事件
func Register(key string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value

	return
}

func getHandler(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()

	value, ok = handlers[key]

	return
}

func HandleEvent(ctx *Client, event string, msg []byte) (data []byte) {
	var (
		ok bool
		f  DisposeFunc
	)

	// 获取处理函数
	if f, ok = getHandler(event); !ok {
		return
	}

	// 执行处理函数
	data = f(ctx.UserId, msg)
	return
}

func PushMsgAll(event string, data []byte) {
	msg := model.WsProto{
		Seq:   "s2c002",
		Event: event,
		Data:  data,
	}
	r, _ := json.Marshal(msg)
	clientManager.sendAll(r, nil)
}

func PushMsg(userId string, event string, data []byte) {
	client := clientManager.GetUserClient(userId)
	if client != nil {
		msg := model.WsProto{
			Seq:   "s2c001",
			Event: event,
			Data:  data,
		}
		r, _ := json.Marshal(msg)
		client.SendMsg(r)
	}
}
