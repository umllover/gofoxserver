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
	*model.Userextrainfo
	*model.Userattr
	Id int
	RoomId int // roomId 就是tableid
	Status int //当前游戏状态
	ChairId int //当前椅子
	sync.RWMutex
}

func NewUser(UserId int) *User {
	return &User{Id : UserId}
}

func (u *User) GetUid() int{
	return u.Id
}

