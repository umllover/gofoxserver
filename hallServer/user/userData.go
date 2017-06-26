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

//扣砖石
func (u *User) SubCurrency(sub int) bool {
	u.Lock()
	defer u.Unlock()
	if u.Currency < sub {
		return false
	}

	err := model.UsertokenOp.UpdateWithMap(u.Id, map[string]interface{}{
		"Currency": u.Currency,
	})
	if err != nil {
		log.Error("at SubCurrency UpdateWithMap error, %v,  sub Currency:%v", err.Error(), sub)
	}
	return true
}

//加砖石
func (u *User) AddCurrency(add int) bool {
	u.Lock()
	defer u.Unlock()
	u.Currency += add
	err := model.UsertokenOp.UpdateWithMap(u.Id, map[string]interface{}{
		"Currency": u.Currency,
	})
	if err != nil {
		log.Error("at AddCurrency UpdateWithMap error, %v,  sub Currency:%v", err.Error(), add)
		return false
	}
	return true
}
