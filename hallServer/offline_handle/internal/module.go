package internal

import (
	"mj/gameServer/base"

	"time"

	"container/list"

	"mj/common/msg"

	"mj/hallServer/game_list"

	"sort"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

var (
	skeleton  = base.NewSkeleton()
	ChanRPC   = skeleton.ChanRPCServer
)

type MatchModule struct {
	*module.Skeleton
}

func (m *MatchModule) OnInit() {
	m.Skeleton = skeleton
}

func (m *MatchModule) OnDestroy() {

}

