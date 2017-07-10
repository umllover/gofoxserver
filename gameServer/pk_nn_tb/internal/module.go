package internal

import (
	"mj/gameServer/RoomMgr"
	"mj/gameServer/base"
	"mj/gameServer/pk_nn_tb/room"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

func (m *Module) OnDestroy() {

}

func (m *Module) CreateRoom(args ...interface{}) bool {
	r := room.CreaterRoom(args)
	if r == nil {
		return false
	}

	return RoomMgr.AddRoom(r)
}

func (m *Module) GetChanRPC() *chanrpc.Server {
	return ChanRPC
}

func getRoom(id int) RoomMgr.IRoom {
	return RoomMgr.GetRoom(id)
}

func (m *Module) GetClientCount() int {
	return 0
}
func (m *Module) GetTableCount() int {
	return 0
}
