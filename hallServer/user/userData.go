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
	Rooms map[int]*model.CreateRoomInfo
	Id    int
	sync.RWMutex
}

func NewUser(UserId int) *User {
	u := &User{Id: UserId}
	u.Rooms = make(map[int]*model.CreateRoomInfo)
	return u
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

func (u *User) AddRooms(id int, r *model.CreateRoomInfo) {
	u.Lock()
	defer u.Unlock()
	u.Rooms[id] = r
}

func (u *User) DelRooms(id int) {
	u.Lock()
	defer u.Unlock()
	_, ok := u.Rooms[id]
	if ok {
		delete(u.Rooms, id)
		model.CreateRoomInfoOp.Delete(id)
	}
}

func (u *User) HasRoom(id int) bool {
	u.RLock()
	defer u.RUnlock()
	_, ok := u.Rooms[id]
	return ok
}

func (u *User) GetRoomCnt() int {
	return len(u.Rooms)
}
