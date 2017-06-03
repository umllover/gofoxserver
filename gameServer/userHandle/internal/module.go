package internal

import (
	"github.com/lovelly/leaf/module"
	"mj/gameServer/base"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/gate"
	"sync"
)

var (
	Users = make(map[int]struct{}) //key is userId
	UsersLock sync.RWMutex
)

func NewUserHandle(a gate.Agent) *module.Skeleton {
	log.Debug("at NewUserHandle === ")
	m := new(Module)
	m.Skeleton = base.NewSkeleton()
	m.ChanRPC = m.Skeleton.ChanRPCServer

	m.a = a
	RegisterHandler(m)
	m.OnInit()
	return m.Skeleton
}


type Module struct {
	*module.Skeleton
	ChanRPC *chanrpc.Server
	a gate.Agent
}

func (m *Module) OnInit() {


}

func (m *Module) OnDestroy() {

}

func (m *Module)Close(){
	m.a.Close()
}



func HasUser(uid int) bool{
	UsersLock.RLock()
	defer  UsersLock.RUnlock()
	_, ok := Users[uid]
	return ok
}

func AddUser(uid int) {
	UsersLock.Lock()
	defer  UsersLock.Unlock()
	Users[uid]= struct {}{}
}

func DelUser(uid int) {
	UsersLock.Lock()
	defer  UsersLock.Unlock()
	delete(Users, uid)
}



