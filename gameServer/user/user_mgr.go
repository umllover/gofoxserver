package user

import (
	"sync"

	"github.com/lovelly/leaf/chanrpc"
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

func GetUser(uid int64) (*User, bool) {
	UsersLock.RLock()
	defer UsersLock.RUnlock()
	u, ok := Users[uid]
	return u, ok
}

func AddUser(uid int64, u *User) {
	CenterRpc.Go("SelfNodeAddPlayer", uid, u.ChanRPC())
	UsersLock.Lock()
	defer UsersLock.Unlock()
	Users[uid] = u
}

func DelUser(uid int64) {
	CenterRpc.Go("SelfNodeDelPlayer", uid)
	UsersLock.Lock()
	defer UsersLock.Unlock()
	delete(Users, uid)
}
