package user

import (
	"mj/hallServer/db/model"
	"sync"
)

type User struct {
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
	return &User{Id : UserId}
}
