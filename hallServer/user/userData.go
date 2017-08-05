package user

import (
	"mj/common/msg"
	"mj/hallServer/db/model"

	"sync"

	"time"

	"mj/hallServer/db/model/base"

	datalog "mj/hallServer/log"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

//请注意， gameServer 只读属性， 不负责存库

type User struct {
	gate.Agent
	*base.FreeLimit
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

	//非入库字段
	MacKCodeTime *time.Time
	Id           int64
	sync.RWMutex
}

func NewUser(UserId int64) *User {
	u := &User{Id: UserId}
	u.Rooms = make(map[int]*model.CreateRoomInfo)
	u.Records = make(map[int]*model.TokenRecord)
	return u
}

func (u *User) GetUid() int64 {
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
	defer u.Unlock()
	_, ok := u.Rooms[id]
	if ok {
		delete(u.Rooms, id)
		model.CreateRoomInfoOp.Delete(id)
	}
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
		RoomInfo.KindID = v.KindId
		info = append(info, RoomInfo)
	}
	return info
}

func (u *User) GetRoomCnt() int {
	return len(u.Rooms)
}

func (u *User) EnoughCurrency(sub int) bool {
	u.Lock()
	defer u.Unlock()
	if u.Currency < sub {
		return false
	}

	return true
}

//扣砖石
func (u *User) SubCurrency(sub, subtype int) bool {
	u.Lock()
	defer u.Unlock()
	if u.Currency < sub {
		return false
	}

	consum := datalog.ConsumLog{}
	consum.AddConsumLogInfo(u.Id, subtype, sub)
	u.Currency -= sub
	err := model.UsertokenOp.UpdateWithMap(u.Id, map[string]interface{}{
		"Currency": u.Currency,
	})
	if err != nil {
		u.Currency += sub
		log.Error("at SubCurrency UpdateWithMap error, %v,  sub Currency:%v", err.Error(), sub)
	}

	u.UpdateUserAttr(map[string]interface{}{
		"Diamond": u.Currency,
	})
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
	//通知客户端
	u.UpdateUserAttr(map[string]interface{}{
		"Diamond": u.Currency,
	})
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

func (u *User) HasRecord(RoomId int) bool {
	u.Lock()
	u.Unlock()
	_, ok := u.Records[RoomId]
	return ok
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
	if u.EnterIP == "" && u.Roomid == 0 {
		return
	}

	log.Debug("at DelGameLockInfo######################### ")
	u.KindID = 0
	u.ServerID = 0
	u.GameNodeID = 0
	u.EnterIP = ""
	u.Roomid = 0
	err := model.GamescorelockerOp.UpdateWithMap(u.Id, map[string]interface{}{
		"KindID":     0,
		"ServerID":   0,
		"GameNodeID": 0,
		"EnterIP":    "",
		"roomid":     0,
	})
	if err != nil {
		log.Error("at EnterRoom  updaye .Gamescorelocker error:%s", err.Error())
	}
}

//同步变动属性给客户端
func (u *User) UpdateUserAttr(m map[string]interface{}) {
	u.WriteMsg(&msg.L2C_UpdateUserAttr{Data: m})
}

//判断是否免费
func (u *User) CheckFree() bool {
	u.Lock()
	defer u.Unlock()
	t := time.Now()
	tm := &t
	b := false
	for _, v := range base.FreeLimitCache.All() {
		if tm.Before(*v.FreeBegin) || tm.After(*v.FreeEnd) {
			continue
		}

		b = true
		log.Debug("b", b)
	}
	return b
}
