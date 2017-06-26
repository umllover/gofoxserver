package userHandle

import (
	"mj/gameServer/userHandle/internal"

	"github.com/lovelly/leaf/gate"
)

var (
	UserMgr     = new(internal.MgrModule)
)

func NewUserHandle(a gate.Agent) gate.UserHandler {
	return internal.NewUserHandle(a)
}
