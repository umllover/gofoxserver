package cluster

import (
	"fmt"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/network/gob"
)

var (
	routeMap  = map[interface{}]*chanrpc.Client{}
	Processor = gob.NewProcessor()
)

func SetRoute(id interface{}, server *chanrpc.Server) {
	_, ok := routeMap[id]
	if ok {
		panic(fmt.Sprintf("function id %v: already set route", id))
	}

	routeMap[id] = server.Open(0)
}

type S2S_RequestMsg struct {
	RequestID uint32
	MsgID     interface{}
	CallType  uint8
	Args      []interface{}
}

type S2S_ResponseMsg struct {
	RequestID uint32
	Ret       interface{}
	Err       string
}
