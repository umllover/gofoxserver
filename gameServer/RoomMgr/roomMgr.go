package RoomMgr

import (
	. "mj/common/cost"
	"mj/common/msg"
	"sync"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/log"
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
	cluster.Broadcast(HallPrefix, "notifyNewRoom", r.GetBirefInfo())
	mgrLock.Lock()
	defer mgrLock.Unlock()
	if _, ok := Rooms[r.GetRoomId()]; ok {
		log.Debug("at AddRoom doeble add, roomid:%v", r.GetRoomId())
		r.Destroy(r.GetRoomId())
		return false
	}
	Rooms[r.GetRoomId()] = r
	return true
}

func GetRoom(id int) IRoom {
	mgrLock.RLock()
	defer mgrLock.RUnlock()
	return Rooms[id]
}

func DelRoom(id int) {
	cluster.Broadcast(HallPrefix, "notifyDelRoom", id)
	mgrLock.Lock()
	defer mgrLock.Unlock()
	delete(Rooms, id)
}
