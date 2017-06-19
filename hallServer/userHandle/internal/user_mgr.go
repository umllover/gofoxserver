package internal

import (
	"mj/common/base"
	"mj/hallServer/center"
	"sync"

	"mj/hallServer/user"

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

func (m *MgrModule) ForEachUser(f func(u *user.User)) {
	UsersLock.RLock()
	defer UsersLock.RUnlock()
	for _, u := range Users {
		if u != nil {
			f(u)
		}
	}
}

func GetUser(uid int) *user.User {
	UsersLock.RLock()
	defer UsersLock.RUnlock()
	u, _ := Users[uid]
	return u
}

func HasUser(uid int) bool {
	UsersLock.RLock()
	defer UsersLock.RUnlock()
	_, ok := Users[uid]
	return ok
}

func AddUser(uid int, u *user.User) {
	log.Debug("AddUser: %d ===", uid)
	center.ChanRPC.Go("SelfNodeAddPlayer", uid, u.ChanRPC())
	UsersLock.Lock()
	defer UsersLock.Unlock()
	Users[uid] = u
}

func DelUser(uid int) {
	log.Debug("deluser %d ===", uid)
	center.ChanRPC.Go("SelfNodeDelPlayer", uid)
	UsersLock.Lock()
	defer UsersLock.Unlock()
	delete(Users, uid)
}
