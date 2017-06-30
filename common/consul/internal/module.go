package internal

import (
	"mj/common/base"

	"strings"

	"github.com/lovelly/leaf/log"
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
	log.Debug("at consul model OnInit")
	m.Skeleton = skeleton
	InitConsul("http")
	if Config.GetRegistSelf() {
		Register()
	}

	wn := Config.GetWatchSvrName()
	if wn != "" {
		list := strings.Split(wn, ",")
		for _, v := range list {
			WatchServices(v)
		}
	}
}

func (m *Module) OnDestroy() {
	deregDeregister()
	log.Debug("at consul model OnDestroy")
}
