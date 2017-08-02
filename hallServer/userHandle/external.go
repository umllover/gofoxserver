package userHandle

import (
	"mj/hallServer/user"
	"mj/hallServer/userHandle/internal"
	"time"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

var (
	UserMgr = new(internal.MgrModule)
)

func NewUserHandle(a gate.Agent) gate.UserHandler {
	return internal.NewUserHandle(a)
}

func ForEachUser(f func(u *user.User)) {
	internal.ForEachUser(f)
}

func KickOutUser() {
	log.Debug("at hallServer close, KickOutUser")
	ForEachUser(func(player *user.User) {
		player.ChanRPC().Go("SvrShutdown")
	})
	time.Sleep(5 * time.Second)
}
