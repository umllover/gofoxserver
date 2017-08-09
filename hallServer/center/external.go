package center

import (
	"mj/hallServer/center/internal"

	"github.com/lovelly/leaf/chanrpc"
)

var (
	Module  = new(internal.Module)
	ChanRPC = internal.ChanRPC
)

func SetGameListRpc(rpc *chanrpc.Server) {
	internal.GamelistRpc = rpc
}

//发送消息给本服服务器上的玩家
func SendToThisNodeUser(uid int64, funcName string, data interface{}) {
	ChanRPC.Go("SendMsgToSelfNotdeUser", uid, funcName, data)
}

func SendMsgToHallUser(uid int64, data interface{}) {
	ChanRPC.Go("SendMsgToHallUser", uid, data)
}

func SetOfflineHandler(fn func(htype string, uid int64, data interface{}, Notify bool) bool) {
	internal.AddOfflineHandler = fn
}
