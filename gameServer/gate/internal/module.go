package internal

import (
	"github.com/lovelly/leaf/gate"
	"mj/gameServer/conf"
	"mj/common/msg"
	"mj/gameServer/login"
)

type Module struct {
	*gate.Gate
}

func (m *Module) OnInit() {
	m.Gate = &gate.Gate{
		MaxConnNum:         conf.Server.MaxConnNum,
		PendingWriteNum:    conf.PendingWriteNum,
		MaxMsgLen:          conf.MaxMsgLen,
		WSAddr:             conf.Server.WSAddr,
		HTTPTimeout:        conf.HTTPTimeout,
		CertFile:           conf.Server.CertFile,
		KeyFile:            conf.Server.KeyFile,
		TCPAddr:            conf.Server.TCPAddr,
		LenMsgLen:          conf.LenMsgLen,
		LittleEndian:       conf.LittleEndian,
		Processor:          msg.Processor,
		GoLen:              conf.AgentGoLen,
		TimerDispatcherLen: conf.AgentTimerDispatcherLen,
		AsynCallLen:        conf.AgentAsynCallLen,
		//ChanRPCLen:         conf.AgentChanRPCLen,
		OnAgentInit:        onAgentInit,
		OnAgentDestroy:     onAgentDestroy,
		AgentChanRPC: 		login.ChanRPC,
	}

}
