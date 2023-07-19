package ws

import (
	"fmt"
	"sync"
	"time"
)

var clientManager = NewClientManager()

// ClientManager 连接管理
type ClientManager struct {
	Clients     map[*Client]bool   // 全部的连接
	ClientsLock sync.RWMutex       // 读写锁
	Users       map[string]*Client // 登录的用户 // userid
	UserLock    sync.RWMutex       // 读写锁
	Register    chan *Client       // 连接连接处理
	Unregister  chan *Client       // 断开连接处理程序
	Broadcast   chan []byte        // 广播 向全部成员发送数据
}

func NewClientManager() (clientManager *ClientManager) {
	clientManager = &ClientManager{
		Clients:    make(map[*Client]bool),
		Users:      make(map[string]*Client),
		Register:   make(chan *Client, 100),
		Unregister: make(chan *Client, 100),
		Broadcast:  make(chan []byte, 100),
	}

	return
}

/**************************  manager  ***************************************/

func (manager *ClientManager) InClient(client *Client) (ok bool) {
	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	// 连接存在，在添加
	_, ok = manager.Clients[client]

	return
}

// GetClients
func (manager *ClientManager) GetClients() (clients map[*Client]bool) {

	clients = make(map[*Client]bool)

	manager.ClientsRange(func(client *Client, value bool) (result bool) {
		clients[client] = value

		return true
	})

	return
}

// 遍历
func (manager *ClientManager) ClientsRange(f func(client *Client, value bool) (result bool)) {

	manager.ClientsLock.RLock()
	defer manager.ClientsLock.RUnlock()

	for key, value := range manager.Clients {
		result := f(key, value)
		if result == false {
			return
		}
	}

	return
}

// GetClientsLen
func (manager *ClientManager) GetClientsLen() (clientsLen int) {

	clientsLen = len(manager.Clients)

	return
}

// 添加客户端
func (manager *ClientManager) AddClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	manager.Clients[client] = true
}

// 删除客户端
func (manager *ClientManager) DelClients(client *Client) {
	manager.ClientsLock.Lock()
	defer manager.ClientsLock.Unlock()

	if _, ok := manager.Clients[client]; ok {
		delete(manager.Clients, client)
	}
}

// 获取用户的连接
func (manager *ClientManager) GetUserClient(userId string) (client *Client) {

	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()

	if value, ok := manager.Users[userId]; ok {
		client = value
	}

	return
}

// GetUsersLen
func (manager *ClientManager) GetUsersLen() (userLen int) {
	userLen = len(manager.Users)

	return
}

// 添加用户
func (manager *ClientManager) AddUsers(key string, client *Client) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	manager.Users[key] = client
}

// 删除用户
func (manager *ClientManager) DelUsers(client *Client) (result bool) {
	manager.UserLock.Lock()
	defer manager.UserLock.Unlock()

	if value, ok := manager.Users[client.UserId]; ok {
		// 判断是否为相同的用户
		if value.Addr != client.Addr {

			return
		}
		delete(manager.Users, client.UserId)
		result = true
	}

	return
}

// 获取用户的key
func (manager *ClientManager) GetUserList() (userList []string) {
	userList = make([]string, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for id := range manager.Users {
		userList = append(userList, id)
	}

	return
}

// 获取用户的client
func (manager *ClientManager) GetUserClients() (clients []*Client) {
	clients = make([]*Client, 0)
	manager.UserLock.RLock()
	defer manager.UserLock.RUnlock()
	for _, v := range manager.Users {
		clients = append(clients, v)
	}

	return
}

// 向全部成员发送数据
func (manager *ClientManager) sendAll(message []byte, ignoreClient *Client) {
	clients := manager.GetUserClients()
	for _, conn := range clients {
		if conn != ignoreClient {
			conn.SendMsg(message)
		}
	}
}

// 用户建立连接事件
func (manager *ClientManager) ClientRegister(client *Client) {
	manager.AddClients(client)

	// 连接存在，在添加
	if manager.InClient(client) {
		manager.AddUsers(client.UserId, client)
	}

	fmt.Println("ClientRegister 用户建立连接", client.UserId, client.Addr)
	// 连接成功，处理第一次连接
	HandleEvent(client, "conn.push", nil)
}

// 用户断开连接
func (manager *ClientManager) ClientUnregister(client *Client) {
	manager.DelClients(client)

	// 删除用户连接
	deleteResult := manager.DelUsers(client)
	if !deleteResult {
		// 不是当前连接的客户端
		return
	}

	_ = client.Socket.Close()
	fmt.Println("客户端取消注册成功", client.Addr, client.UserId)
}

// 管道处理程序
func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.Register:
			// 建立连接事件
			manager.ClientRegister(conn)

		case conn := <-manager.Unregister:
			// 断开连接事件
			manager.ClientUnregister(conn)

		case message := <-manager.Broadcast:
			// 广播事件
			clients := manager.GetClients()
			for conn := range clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
				}
			}
		}
	}
}

/**************************  manager info  ***************************************/

func ClientRegister(client *Client) {
	clientManager.Register <- client
}

func StartClientManager() {
	clientManager.start()
}

// 获取用户所在的连接
func GetUserClient(userId string) (client *Client) {
	client = clientManager.GetUserClient(userId)
	return
}

// 定时清理超时连接
func ClearTimeoutConnections() {
	ticker := time.NewTicker(time.Minute)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case t := <-ticker.C:
			clients := clientManager.GetClients()
			for client := range clients {
				if client.IsHeartbeatTimeout(uint64(t.Unix())) {
					fmt.Println("心跳时间超时 关闭连接", client.Addr, client.UserId, client.HeartbeatTime)
					client.Offline()
				}
			}
		}
	}
}

// 获取全部用户
func GetUserList() (userList []string) {
	userList = clientManager.GetUserList()
	return
}
