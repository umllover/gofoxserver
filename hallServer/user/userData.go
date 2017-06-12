package user

import (
	"mj/hallServer/db/model"
	"sync"

	"github.com/lovelly/leaf/gate"
)

type User struct {
	gate.Agent
	*model.Accountsmember
	*model.Gamescorelocker
	*model.Gamescoreinfo
	*model.Userattr
	*model.Usertoken
	*model.Userextrainfo
	Id int
	sync.RWMutex
}

func NewUser(UserId int) *User {
	return &User{Id: UserId}
}
