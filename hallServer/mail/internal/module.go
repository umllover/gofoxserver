package internal

import (
	"mj/common/msg"
	"mj/hallServer/base"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
)

type MachPlayer struct {
	ch      *chanrpc.Server
	EndTime int64
	Uid     int64
}

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type MailModule struct {
	*module.Skeleton
	rooms map[int]*msg.RoomInfo
}

func (m *MailModule) OnInit() {
	m.Skeleton = skeleton
}

func (m *MailModule) OnDestroy() {

}
