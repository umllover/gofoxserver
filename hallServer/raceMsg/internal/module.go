package internal

import (
	"mj/hallServer/base"

	"time"

	"github.com/lovelly/leaf/module"
)

const HorseRaceInterval = 5 * 60

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	//m.Skeleton.AfterFunc(HorseRaceInterval*time.Second, m.StartHorseRaceLamp)
}

func (m *Module) OnDestroy() {

}

func (m *Module) StartHorseRaceLamp() {
	StartHorseRaceLamp() // 跑马灯
	m.Skeleton.AfterFunc(HorseRaceInterval*time.Second, m.StartHorseRaceLamp)
}
