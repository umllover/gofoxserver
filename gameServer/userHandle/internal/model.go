package internal

import (
	"mj/common/base"
	"time"

	"github.com/lovelly/leaf/log"

	"mj/gameServer/user"

	"github.com/lovelly/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
)

type MgrModule struct {
	*module.Skeleton
}

func (m *MgrModule) OnInit() {
	m.Skeleton = skeleton
}

func (m *MgrModule) OnDestroy() {
	log.Debug("at server close offline user ")
	ForEachUser(func(player *user.User) {
		log.Debug("111111111111111111 ")
		player.ChanRPC().Go("SvrShutdown")
	})
	time.Sleep(5 * time.Second)
}
