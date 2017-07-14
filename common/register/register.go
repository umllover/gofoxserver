package register

import (
	"mj/common/msg"

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
	r.rpc.Register(id, f)
}

func (r Reg) RegisterS2S(id interface{}, f interface{}) {
	cluster.Processor.SetRouter(id, r.rpc)
	r.rpc.Register(id, f)
}

func NewRegister(rpc *chanrpc.Server) *Reg {
	return &Reg{rpc: rpc}
}
