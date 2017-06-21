package room_base

import (
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
)

type Module interface {
	Destroy(int)
	RoomRun(int)
	Skeleton() *module.Skeleton
	GetChanRPC() *chanrpc.Server
}
