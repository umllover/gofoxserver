package internal

import (
	"mj/hallServer/base"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
)

var (
	skeleton   = base.NewSkeleton()
	ChanRPC    = skeleton.ChanRPCServer
	Users      = make(map[int]*chanrpc.Server) //本服玩家
	OtherUsers = make(map[int]string)          //其他服登录的玩家  key is uid， values is NodeId
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {

}
