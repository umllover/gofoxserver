package userHandle

import (
	"mj/hallServer/userHandle/internal"

	"mj/hallServer/user"

	"github.com/lovelly/leaf/gate"
)

var (
	UserMgr = new(internal.MgrModule)
)

func NewUserHandle(a gate.Agent) gate.UserHandler {
	return internal.NewUserHandle(a)
}

func GetUser(uid int) *user.User {
	return internal.GetUser(uid)
}

func ForEachUser(f func(u *user.User)) {
	UserMgr.ForEachUser(f)
}
