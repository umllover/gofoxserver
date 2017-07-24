package consul

import (
	"github.com/lovelly/leaf/chanrpc"
)

var (
	Module = new(ConsulModule)
)

func SetConfig(cfg Rgconfig) {
	Config = cfg
}

func SetSelfId(selfId string) {
	SelfId = selfId
}

func SetHookRpc(rpc *chanrpc.Server) {
	HookChanRpc = rpc
}
