package internal

import (
	"github.com/lovelly/leaf/module"
	"mj/gameServer/base"
	"mj/gameServer/hzmj/room"
	"time"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/conf"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	rooms = make(map[int]*room.Room)
	clientCount int = 0
	wTableCount int = 0
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	m.Skeleton.AfterFunc(10 * time.Second, m.checkUpdate)
}

func (m *Module) OnDestroy() {
	for _, r := range rooms {
		r.Destroy()
	}
}

func(m *Module) GetChanRPC() (*chanrpc.Server){
	return ChanRPC
}

func (m *Module) checkUpdate() {
	//todo 向大厅报告人数
	m.Skeleton.AfterFunc(10 * time.Second, m.checkUpdate)
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
	rooms[r.GetRoomId()] = r
	AddTableCount()
	msg := &msg.RoomInfo{}
	msg.ServerID = r.ServerId
	msg.KindID = r.Kind
	msg.NodeId = conf.Server.NodeId
	msg.TableId = r.GetRoomId()
	cluster.Broadcast(HallPrefix,"notifyNewRoom", msg)
}

func delRoom(id int) {
	delete(rooms, id)
	cluster.Broadcast(HallPrefix, "notifyDelRoom", id)
}

func getRoom(id int) *room.Room{
	return rooms[id]
}