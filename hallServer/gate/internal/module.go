package internal

import (
	"mj/common/msg"
	"mj/hallServer/conf"
	"mj/hallServer/userHandle"

	"strings"

	"github.com/lovelly/leaf/gate"
)

type Module struct {
	*gate.Gate
}

func (m *Module) OnInit() {
	list := strings.Split(conf.Server.WSAddr, ":")
	var listenAddr string
	if len(list) > 1 {
		listenAddr = "0.0.0.0:" + list[1]
	}

	m.Gate = &gate.Gate{
		MaxConnNum:         conf.Server.MaxConnNum,
		PendingWriteNum:    conf.PendingWriteNum,
		MaxMsgLen:          conf.MaxMsgLen,
		WSAddr:             listenAddr,
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
		NewChanRPCFunc:     userHandle.NewUserHandle,
		OnAgentInit:        onAgentInit,
		OnAgentDestroy:     onAgentDestroy,
	}

}
