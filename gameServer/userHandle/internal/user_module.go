package internal

import (
	"mj/gameServer/base"

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
	closeCh chan bool
	ChanRPC *chanrpc.Server
	a       gate.Agent
}

func (m *UserModule) Run() {
	go func() {
		m.Skeleton.Run(m.closeCh)
		m.OnDestroy()
	}()
}

func (m *UserModule) Close(Reason int) {
	defer func() {
		m.a.Close()
		m.closeCh <- true
	}()
	m.UserOffline()
}

func (m *UserModule) GetChanRPC() *chanrpc.Server {
	return m.ChanRPC
}

///////////////////////////////////
func (m *UserModule) OnInit() {

}

func (m *UserModule) OnDestroy() {

}
