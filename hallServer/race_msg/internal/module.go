package internal

import (
	"mj/hallServer/base"

	"github.com/lovelly/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	InitRaceMsg()
}

func (m *Module) OnDestroy() {

}
