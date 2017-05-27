package user

import (
	"mj/gameServer/db/model"
	"github.com/lovelly/leaf/gate"
	"sync"
)

type User struct {
	gate.Agent
	*model.Accountsinfo
	*model.Accountsmember
	*model.Gamescorelocker
	*model.Gamescoreinfo
	*model.Userroomcard
	*model.Userextrainfo
	Id int
	RoomId int
	sync.RWMutex
}

func NewUser(UserId int) *User {
	return &User{Id : UserId}
}