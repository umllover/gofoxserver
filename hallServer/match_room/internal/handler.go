package internal

import (
	"mj/common/msg"
	"reflect"

	"github.com/lovelly/leaf/gate"
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}) {
	ChanRPC.Register(id, f)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handlerC2S(&msg.C2L_QuickMatch{}, QuickMatch)
}

func QuickMatch(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_QuickMatch)
	agent := args[1].(gate.Agent)
	DefaultMachModule.AddMatchPlayer(recvMsg.KindID, &MachPlayer{ch: agent.ChanRPC()})
}
