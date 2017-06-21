package room_base

import (
	"mj/gameServer/common"
	"sync"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

//房间管理类

/// 房间里面的玩家管理
type RoomBase struct {
	// module 必须字段
	*module.Skeleton
	ChanRPC  *chanrpc.Server //接受客户端消息的chan
	CloseSig chan bool
	wg       sync.WaitGroup //
}

func NewRoomBase() common.BaseManager {
	r := new(RoomBase)
	skeleton := &module.Skeleton{
		GoLen:              1000,
		TimerDispatcherLen: 1000,
		AsynCallLen:        1000,
		ChanRPCServer:      chanrpc.NewServer(1000),
	}
	skeleton.Init()
	r.Skeleton = skeleton
	r.ChanRPC = skeleton.ChanRPCServer
	return r
}

func (r *RoomBase) GetSkeleton() *module.Skeleton {
	return r.Skeleton
}

func (r *RoomBase) RoomRun(id int) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Recover(r)
			}
		}()

		log.Debug("room Room start run ID:%d", id)
		r.Run(r.CloseSig)
		log.Debug("room Room End run ID:%d", id)
	}()
}

func (r *RoomBase) Destroy(id int) {
	defer func() {
		if r := recover(); r != nil {
			log.Recover(r)
		}
	}()
	r.CloseSig <- true
	log.Debug("room Room Destroy ok %d", id)
}

func (r *RoomBase) GetChanRPC() *chanrpc.Server {
	return r.ChanRPC
}
