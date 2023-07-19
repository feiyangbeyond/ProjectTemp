package util

import "encoding/json"

type WsResp struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func MakeWsResp(code uint32, msg string, data interface{}) []byte {
	resp := &WsResp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	b, _ := json.Marshal(resp)
	return b
}
