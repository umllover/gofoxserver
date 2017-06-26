package room_base

import (
	"github.com/lovelly/leaf/chanrpc"
)

type Module interface {
	GetChanRPC() *chanrpc.Server
	GetClientCount() int
	GetTableCount() int
	OnDestroy()
	OnInit()
	Run(chan bool)
	CreateRoom(args ...interface{}) bool
}
