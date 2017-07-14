package user

import (
	"mj/common/msg"
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
	Rooms     map[int]*model.CreateRoomInfo
	Records   map[int]*model.TokenRecord
	Times     map[int]int64 //永久次数
	DayTimes  map[int]int64 //每日次数
	WeekTimes map[int]int64 //周次数
	Id        int
	sync.RWMutex
}

func NewUser(UserId int) *User {
	u := &User{Id: UserId}
	u.Rooms = make(map[int]*model.CreateRoomInfo)
	u.Records = make(map[int]*model.TokenRecord)
	return u
}

func (u *User) GetUid() int {
	return u.Id
}

func (u *User) AddRooms(r *model.CreateRoomInfo) {
	model.CreateRoomInfoOp.Insert(r)
	u.Lock()
	defer u.Unlock()
	u.Rooms[r.RoomId] = r
}

func (u *User) DelRooms(id int) {
	u.Lock()
	_, ok := u.Rooms[id]
	if ok {
		delete(u.Rooms, id)
	}
	u.Unlock()
	model.CreateRoomInfoOp.Delete(id)
}

func (u *User) GetRoom(id int) *model.CreateRoomInfo {
	u.RLock()
	defer u.RUnlock()
	return u.Rooms[id]
}

func (u *User) GetRoomInfo() []*msg.CreatorRoomInfo {
	u.RLock()
	defer u.RUnlock()
	info := make([]*msg.CreatorRoomInfo, 0)
	for _, v := range u.Rooms {
		RoomInfo := &msg.CreatorRoomInfo{}
		RoomInfo.Status = v.Status
		RoomInfo.CreatorTime = v.CreateTime.Unix()
		RoomInfo.RoomName = v.RoomName
		RoomInfo.RoomID = v.RoomId
		info = append(info, RoomInfo)
	}
	return info
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

//增加扣钱计入
func (u *User) AddRecord(tr *model.TokenRecord) bool {
	u.Lock()
	u.Records[tr.RoomId] = tr
	u.Unlock()
	_, err := model.TokenRecordOp.Insert(tr)
	if err != nil {
		log.Debug("ad TokenRecordOp error :%s", err.Error())
		return false
	}
	return true
}

//删除扣钱记录
func (u *User) DelRecord(id int) error {
	u.Lock()
	r, ok := u.Records[id]
	if ok {
		delete(u.Records, id)
	}
	u.Unlock()
	return model.TokenRecordOp.Delete(r.RoomId, r.UserId)
}

func (u *User) GetRecord(id int) *model.TokenRecord {
	u.RLock()
	defer u.RUnlock()
	return u.Records[id]
}

func (u *User) DelGameLockInfo() {
	u.KindID = 0
	u.ServerID = 0
	u.EnterIP = ""
	u.GameNodeID = 0
	err := model.GamescorelockerOp.UpdateWithMap(u.Id, map[string]interface{}{
		"GameNodeID": "",
		"EnterIP":    "",
		"KindID":     0,
		"ServerID":   0,
	})
	if err != nil {
		log.Error("at EnterRoom  updaye .Gamescorelocker error:%s", err.Error())
	}
}
func (u *User) SetVip() {

}
