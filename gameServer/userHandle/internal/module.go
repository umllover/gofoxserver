package internal

import (
	"github.com/lovelly/leaf/module"
	"mj/gameServer/base"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	KModule  = new(Module)
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton

}

func (m *Module) OnDestroy() {

}


