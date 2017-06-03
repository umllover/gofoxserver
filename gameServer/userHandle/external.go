package userHandle

import (
	"mj/gameServer/userHandle/internal"
	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/module"
)

var (

)

func NewUserHandle(a gate.Agent) *module.Skeleton {
	return internal.NewUserHandle(a)
}



