package user

import (
	"mj/hallServer/db/model"
	"sync"

	"github.com/lovelly/leaf/gate"
)

type User struct {
	gate.Agent
	*model.Accountsinfo
	*model.Accountsmember
	*model.Gamescorelocker
	*model.Gamescoreinfo
	*model.Userattr
	*model.Userextrainfo
	Id int
	sync.RWMutex
}

func NewUser(UserId int) *User {
	return &User{Id: UserId}
}
