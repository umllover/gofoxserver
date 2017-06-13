package internal

import (
	"mj/common/base"
	"mj/gameServer/center"
	"mj/gameServer/user"
	"sync"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

var (
	skeleton  = base.NewSkeleton()
	Users     = make(map[int]*user.User) //key is userId
	UsersLock sync.RWMutex
)

type MgrModule struct {
	*module.Skeleton
}

func (m *MgrModule) OnInit() {
	m.Skeleton = skeleton

}

func (m *MgrModule) OnDestroy() {
	log.Debug("at server close offline user ")
}

func getUser(uid int) (*user.User, bool) {
	UsersLock.RLock()
	defer UsersLock.RUnlock()
	u, ok := Users[uid]
	return u, ok
}

func AddUser(uid int, u *user.User) {
	center.ChanRPC.Go("SelfNodeAddPlayer", uid, u.ChanRPC())
	UsersLock.Lock()
	defer UsersLock.Unlock()
	Users[uid] = u
}

func DelUser(uid int) {
	center.ChanRPC.Go("SelfNodeDelPlayer", uid)
	UsersLock.Lock()
	defer UsersLock.Unlock()
	delete(Users, uid)
}
