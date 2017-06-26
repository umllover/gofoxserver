package consul

import (
	"mj/common/consul/internal"
)

var (
	Module  = new(internal.Module)
	ChanRPC = internal.ChanRPC
)

func SetConfig(cfg internal.Rgconfig) {
	internal.Config = cfg
}

func SetSelfId(selfId string) {
	internal.SelfId = selfId
}

func AddinitiativeSvr(svrName ...string) {
	internal.InitiativeSvr = append(internal.InitiativeSvr, svrName...)
}
