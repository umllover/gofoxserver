package internal

import (
	"mj/common/msg"
	"mj/hallServer/base"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

func NewUserHandle(a gate.Agent) gate.UserHandler {
	log.Debug("at NewUserHandle === ")
	m := new(UserModule)
	m.Skeleton = base.NewSkeleton()
	m.ChanRPC = m.Skeleton.ChanRPCServer
	m.closeCh = make(chan bool, 1)
	m.a = a
	RegisterHandler(m)
	m.OnInit()
	return m
}

type UserModule struct {
	*module.Skeleton
	ChanRPC *chanrpc.Server
	closeCh chan bool
	a       gate.Agent
}

func (m *UserModule) OnInit() {

}

func (m *UserModule) OnDestroy() {

}

func (m *UserModule) Run() {
	go func() {
		m.Skeleton.Run(m.closeCh)
		m.OnDestroy()
	}()
}

func (m *UserModule) Close(Reason int) {
	m.a.WriteMsg(&msg.L2C_KickOut{Reason: Reason})
	m.a.SetReason(Reason)
	m.a.Close()
	m.closeCh <- true
}

func (m *UserModule) GetChanRPC() *chanrpc.Server {
	return m.ChanRPC
}
