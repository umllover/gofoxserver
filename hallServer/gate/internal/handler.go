package internal

import (
	"github.com/lovelly/leaf/gate"
	"mj/hallServer/UserData"
)

func init() {
	//msg.Processor.SetHandler(&msg.C2F_CheckLogin{}, handleCheckLogin)
}

func onAgentInit(agent gate.Agent) {

}

func onAgentDestroy(agent gate.Agent)() {
	userId, ok := agent.UserData().(int);
	if ok {
		UserData.ChanRPC.Go("delUser",userId)
	}
}

