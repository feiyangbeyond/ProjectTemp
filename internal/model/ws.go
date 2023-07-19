package model

type WsProto struct {
	Seq   string `json:"seq"`
	Event string `json:"event"`
	Data  []byte `json:"data"`
}

const (
	EventConnPush = "conn.push"
	EventLogout   = "logout"
)
