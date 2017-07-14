package internal

import (
	"mj/common/msg"
	"mj/common/register"

	"github.com/lovelly/leaf/gate"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	reg.RegisterC2S(&msg.C2L_QuickMatch{}, QuickMatch)
}

func QuickMatch(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_QuickMatch)
	agent := args[1].(gate.Agent)
	DefaultMachModule.AddMatchPlayer(recvMsg.KindID, &MachPlayer{ch: agent.ChanRPC()})
}
