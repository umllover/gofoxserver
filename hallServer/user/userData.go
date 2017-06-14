package user

import (
	"mj/hallServer/db/model"
	"sync"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

//请注意， gameServer 只读属性， 不负责存库

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

func (u *User) GetUid() int {
	return u.Id
}

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
