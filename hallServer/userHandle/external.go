package userHandle

import (
	"mj/hallServer/user"
	"mj/hallServer/userHandle/internal"

	"github.com/lovelly/leaf/gate"
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
