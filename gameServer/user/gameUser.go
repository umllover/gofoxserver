package user

import (
	"mj/gameServer/db/model"
	"github.com/lovelly/leaf/gate"
	"sync"
	"mj/common/msg"
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
	UserLimit int64  //限制行为
	sync.RWMutex
}

func NewUser(UserId int) *User {
	return &User{Id : UserId}
}

func (u *User) GetUid() int{
	return u.Id
}

func (u *User) SendSysMsg(ty int, context string) {
	u.WriteMsg(&msg.SysMsg{
		ClientID:u.Id,
		Type:ty,
		Context:context,
	})
}


/////////////////////////
//关键函数加锁
func (u *User) GetRoomCard() int {
	u.RLock()
	defer u.RUnlock()
	return u.RoomCard
}

func (u *User) GetCurrency() int {
	u.RLock()
	defer u.RUnlock()
	return u.Currency
}
