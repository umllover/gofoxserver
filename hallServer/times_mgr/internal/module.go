package internal

import (
	"mj/common/msg"
	"mj/gameServer/base"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
)

type MachPlayer struct {
	ch      *chanrpc.Server
	EndTime int64
	Uid     int
}

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type MatchModule struct {
	*module.Skeleton
	rooms map[int]*msg.RoomInfo
}

func (m *MatchModule) OnInit() {
	m.Skeleton = skeleton
	//m.Skeleton.AfterFunc(2*time.Second, m.Match)
}

func (m *MatchModule) OnDestroy() {

}
