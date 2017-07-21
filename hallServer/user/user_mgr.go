package user

import (
	"sync"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
)

var (
	Users     = make(map[int64]*User) //key is userId
	UsersLock sync.RWMutex
	CenterRpc *chanrpc.Server
)

func ForEachUser(f func(u *User)) {
	UsersLock.RLock()
	defer UsersLock.RUnlock()
	for _, u := range Users {
		if u != nil {
			f(u)
		}
	}
}

func GetUser(uid int64) *User {
	UsersLock.RLock()
	defer UsersLock.RUnlock()
	u, _ := Users[uid]
	return u
}

func HasUser(uid int64) bool {
	UsersLock.RLock()
	defer UsersLock.RUnlock()
	_, ok := Users[uid]
	return ok
}

func AddUser(uid int64, u *User) {
	log.Debug("AddUser: %d ===", uid)
	CenterRpc.Go("SelfNodeAddPlayer", uid, u.ChanRPC())
	UsersLock.Lock()
	defer UsersLock.Unlock()
	Users[uid] = u
}

func DelUser(uid int64) {
	log.Debug("deluser %d ===", uid)
	CenterRpc.Go("SelfNodeDelPlayer", uid)
	UsersLock.Lock()
	defer UsersLock.Unlock()
	delete(Users, uid)
}
