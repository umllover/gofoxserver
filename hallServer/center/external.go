package center

import (
	"mj/hallServer/center/internal"
)

var (
	Module  = new(internal.Module)
	ChanRPC = internal.ChanRPC
)

func SendMsgToThisNodeUser(uid int, funcName string, data ...interface{}) {
	ChanRPC.Go("SendMsgToSelfNotdeUser", uid, funcName, data)
}

func SendMsgToUser(uid int, funcName string, data ...interface{}) {
	ChanRPC.Go("SendMsgToUser", uid, funcName, data)
}
