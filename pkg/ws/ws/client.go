package ws

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"deviceback/v3/internal/model"

	"github.com/gorilla/websocket"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime = 6 * 60
)

// Client 用户连接
type Client struct {
	UserId string // userId

	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	ExitC         chan struct{}   // 退出信号
	IsOffline     bool
}

// NewClient 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime uint64, userId string) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 100),
		UserId:        userId,
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
		ExitC:         make(chan struct{}),
	}

	return
}

// 读取客户端数据
func (c *Client) Read() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write panic", string(debug.Stack()), r)
		}
	}()

	defer func() {
		fmt.Println("read协程关闭", c)
	}()

	defer func() {
		c.Offline()
	}()

	for {
		if c.IsOffline {
			break
		}

		t, message, err := c.Socket.ReadMessage()
		if err != nil || t == websocket.CloseMessage {
			break
		}

		// 收到ping
		if t == websocket.TextMessage && strings.ToUpper(string(message)) == "PING" {
			c.Heartbeat(uint64(time.Now().Unix()))
			c.SendMsg([]byte("PONG"))
			continue
		}

		// 处理程序
		go c.processData(message)
	}
}

// 向客户端写数据
func (c *Client) Write() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("write panic", string(debug.Stack()), r)
		}
	}()

	defer func() {
		fmt.Println("write协程关闭", c)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				return
			}
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		case <-c.ExitC:
			return
		}
	}
}

func (c *Client) Offline() {
	if c.IsOffline {
		return
	}

	c.IsOffline = true
	c.ExitC <- struct{}{}
	// 离线并断开连接
	clientManager.Unregister <- c
}

func (c *Client) LogoutWithMsg() {
	PushMsg(c.UserId, model.EventLogout, nil)
	time.Sleep(500 * time.Millisecond)
	c.Offline()
}

// SendMsg 发送数据
func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()

	c.Send <- msg
}

// 用户心跳
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime
	return
}

// IsHeartbeatTimeout 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	if c.HeartbeatTime+heartbeatExpirationTime <= currentTime {
		timeout = true
	}
	return
}

// 处理数据
func (c *Client) processData(message []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("处理数据 stop", r)
		}
	}()

	req := model.WsProto{}

	err := json.Unmarshal(message, &req)
	if err != nil {
		c.SendMsg([]byte("invalid data structure"))
		return
	}

	data := HandleEvent(c, req.Event, req.Data)

	resp := model.WsProto{
		Seq:   req.Seq,
		Event: req.Event,
		Data:  data,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return
	}

	c.SendMsg(b)

	return
}
