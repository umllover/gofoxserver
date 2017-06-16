package internal

import (
	"mj/hallServer/base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

type HorseRaceMsg struct {
	startTime	int	// 起始时间
	endTime		int	// 结束时间
	intervalTime	int	// 间隔时间
	content		string	// 详情内容
}

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	log.Debug("测试")
}

func (m *Module) OnDestroy() {

}
