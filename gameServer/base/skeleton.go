package base

import (
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
	"mj/gameServer/conf"
)

func NewSkeleton() *module.Skeleton {
	skeleton := &module.Skeleton{
		GoLen:              conf.GoLen,
		TimerDispatcherLen: conf.TimerDispatcherLen,
		AsynCallLen:        conf.AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(conf.ChanRPCLen),
	}

	skeleton.Init()
	return skeleton
}
