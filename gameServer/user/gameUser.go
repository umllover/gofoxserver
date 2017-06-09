package user

import (
	"mj/common/msg"
	"mj/gameServer/db/model"
	"sync"

	"github.com/lovelly/leaf/gate"
)

type User struct {
	gate.Agent
	*model.Accountsinfo
	*model.Accountsmember
	*model.Gamescorelocker
	*model.Gamescoreinfo
	*model.Userextrainfo
	*model.Userattr
	Id         int
	RoomId     int   // roomId 就是tableid
	Status     int   //当前游戏状态
	ChairId    int   //当前椅子
	UserLimit  int64 //限制行为
	ChatRoomId int   //聊天房间ID
	sync.RWMutex
}

func NewUser(UserId int) *User {
	return &User{Id: UserId}
}

func (u *User) GetUid() int {
	return u.Id
}

func (u *User) SendSysMsg(ty int, context string) {
	u.WriteMsg(&msg.SysMsg{
		ClientID: u.Id,
		Type:     ty,
		Context:  context,
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
