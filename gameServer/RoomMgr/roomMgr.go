package RoomMgr

import (
	. "mj/common/cost"
	"mj/common/msg"
	"sync"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
)

type IRoom interface {
	GetChanRPC() *chanrpc.Server
	GetRoomId() int
	GetBirefInfo() *msg.RoomInfo
}

var (
	mgrLock sync.RWMutex
	Rooms   = make(map[int]IRoom)
)

func AddRoom(r IRoom) {
	cluster.Broadcast(HallPrefix, "notifyNewRoom", r.GetBirefInfo())
	mgrLock.Lock()
	defer mgrLock.Unlock()
	Rooms[r.GetRoomId()] = r
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
