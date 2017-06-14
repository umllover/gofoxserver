package internal

import (
	"mj/gameServer/RoomMgr"
	"mj/gameServer/base"
	"mj/gameServer/mj_hz/room"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
)

var (
	skeleton        = base.NewSkeleton()
	ChanRPC         = skeleton.ChanRPCServer
	clientCount int = 0
	wTableCount int = 0
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	m.Skeleton.AfterFunc(10*time.Second, m.checkUpdate)
}

func (m *Module) OnDestroy() {

}

func (m *Module) GetChanRPC() *chanrpc.Server {
	return ChanRPC
}

func (m *Module) checkUpdate() {
	//todo 向大厅报告人数
	m.Skeleton.AfterFunc(10*time.Second, m.checkUpdate)
}

func (m *Module) GetClientCount() int {
	return wTableCount
}

func AddClientCount() {
	clientCount++
}

func (m *Module) GetTableCount() int {
	return wTableCount
}

func AddTableCount() {
	wTableCount++
}

func addRoom(r *room.Room) {
	RoomMgr.AddRoom(r)
	AddTableCount()
}

func delRoom(id int) {
	RoomMgr.DelRoom(id)
}

func getRoom(id int) *room.Room {
	r, _ := RoomMgr.GetRoom(id).(*room.Room)
	return r
}
