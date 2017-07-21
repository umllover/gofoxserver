package userHandle

import (
	"mj/hallServer/userHandle/internal"

	"github.com/lovelly/leaf/gate"
)

func NewUserHandle(a gate.Agent) gate.UserHandler {
	return internal.NewUserHandle(a)
}
