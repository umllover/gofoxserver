package user

import (
	"mj/common/msg"
	"mj/gameServer/db/model"
	"sync"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

type User struct {
	gate.Agent
	*model.Usertoken
	*model.Accountsmember
	*model.Gamescorelocker
	*model.Gamescoreinfo
	*model.Userextrainfo
	*model.Userattr
	Id         int
	RoomId     int   // roomId 就是tableid
	Status     int   //当前游戏状态
	offline    bool  //玩家是否在线
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
func (u *User) SubRoomCard(card int) {
	u.Lock()
	defer u.Unlock()
	if card < u.RoomCard {
		log.Error("card < u.RoomCar userId:%d", u.Id)
		u.RoomCard = 0
	}
	u.RoomCard -= card
}

func (u *User) GetRoomCard() int {
	u.RLock()
	defer u.RUnlock()
	return u.RoomCard
}

func (u *User) SubCurrency(menry int) {
	u.Lock()
	defer u.Unlock()
	if menry < u.Currency {
		log.Error("card < u.Currency userId:%d", u.Id)
		u.Currency = 0
	}
	u.Currency -= menry
}

func (u *User) GetCurrency() int {
	u.RLock()
	defer u.RUnlock()
	return u.Currency
}

func (u *User) IsOffline() bool {
	u.RLock()
	defer u.RUnlock()
	return u.offline
}

func (u *User) SetOffline(su bool) bool {
	u.Lock()
	defer u.Unlock()
	u.offline = su
	if u.RoomId != 0 {
		return false
	}
	return true
}

func (u *User) SetRoomId(id int) {
	u.Lock()
	defer u.Unlock()
	u.RoomId = id
}
