package RoomMgr

import (
	"sync"

	"github.com/lovelly/leaf/chanrpc"
)

type IRoom interface {
	GetChanRPC() *chanrpc.Server
	GetRoomId() int
}

var (
	mgrLock sync.RWMutex
	Rooms   = make(map[int]IRoom)
)

func AddRoom(r IRoom) {
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
	mgrLock.Lock()
	defer mgrLock.Unlock()
	delete(Rooms, id)
}
