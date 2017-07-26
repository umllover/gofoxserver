package RoomMgr

import (
	. "mj/common/cost"
	"mj/common/msg"
	"sync"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
)

type IRoom interface {
	GetChanRPC() *chanrpc.Server
	GetRoomId() int
	GetBirefInfo() *msg.RoomInfo
	Destroy(int)
}

var (
	mgrLock sync.RWMutex
	Rooms   = make(map[int]IRoom)
)

func AddRoom(r IRoom) bool {
	mgrLock.Lock()
	if _, ok := Rooms[r.GetRoomId()]; ok {
		mgrLock.Unlock()
		log.Debug("at AddRoom doeble add, roomid:%v", r.GetRoomId())
		r.Destroy(r.GetRoomId())
		return false
	}
	Rooms[r.GetRoomId()] = r
	mgrLock.Unlock()
	return true
}

func GetRoom(id int) IRoom {
	mgrLock.RLock()
	defer mgrLock.RUnlock()
	return Rooms[id]
}

func DelRoom(id int) {
	log.Debug("at aaaaaaaaaaaaaaaaaaaaaa ")
	cluster.Broadcast(HallPrefix, &msg.S2S_notifyDelRoom{RoomID: id})
	mgrLock.Lock()
	defer mgrLock.Unlock()
	delete(Rooms, id)
}

func UpdateRoomToHall(data interface{}) {
	cluster.Broadcast(HallPrefix, data)
}

// 此函数有风险， 请注意 调用函数内不用mgrLock 锁， 此函数消耗也大， 请勿随意调用
func ForEachRoom(cb func(r IRoom)) {
	mgrLock.RLock()
	defer mgrLock.RUnlock()
	for _, v := range Rooms {
		cb(v)
	}
}
