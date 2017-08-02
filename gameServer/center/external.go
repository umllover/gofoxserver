package center

import (
	"mj/gameServer/center/internal"
)

var (
	Module  = new(internal.Module)
	ChanRPC = internal.ChanRPC
)

//发送消息给本服服务器上的玩家
func SendToThisNodeUser(uid int, funcName string, data interface{}) {
	ChanRPC.Go("SendMsgToSelfNotdeUser", uid, funcName, data)
}
