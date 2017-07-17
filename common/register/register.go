package register

import (
	"mj/common/msg"

	"fmt"
	"reflect"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/nsq/cluster"
)

type Reg struct {
	rpc *chanrpc.Server
}

func (r Reg) RegisterRpc(id interface{}, f interface{}) {
	r.rpc.Register(id, f)
}

func (r Reg) RegisterC2S(id interface{}, f interface{}) {
	msg.Processor.SetRouter(id, r.rpc)
	r.rpc.Register(reflect.TypeOf(id), f)
}

func (r Reg) RegisterS2S(id interface{}, f interface{}) {
	msgType := reflect.TypeOf(id)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("function id %v:  msgType == nil || msgType.Kind() != reflect.Ptr", id))
	}

	msgID := msgType.Elem().Name()

	cluster.SetRouter(msgID, r.rpc)
	r.rpc.Register(msgID, f)
}

func NewRegister(rpc *chanrpc.Server) *Reg {
	return &Reg{rpc: rpc}
}
