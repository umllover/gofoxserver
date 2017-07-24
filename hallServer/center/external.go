package center

import (
	. "mj/common/cost"
	"mj/hallServer/center/internal"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/nsq/cluster"
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

//发送消息给游戏服
func SendMsgToGame(svrid int, data interface{}) {
	cluster.Go(GetGameSvrName(svrid), data)
}

func BroadcastToGame(data interface{}) {
	cluster.Broadcast(GamePrefix, data)
}

func AsynCallGame(svrid int, chanAsynRet chan *chanrpc.RetInfo, data interface{}, cb interface{}) {
	cluster.AsynCall(GetGameSvrName(svrid), chanAsynRet, data, cb)
}

//发消息给大厅
func SendMsgToHall(svrid int, data interface{}) {
	cluster.Go(GetHallSvrName(svrid), data)
}

func BroadcastToHall(data interface{}) {
	cluster.Broadcast(HallPrefix, data)
}

func AsynCallHall(svrid int, chanAsynRet chan *chanrpc.RetInfo, data interface{}, cb interface{}) {
	cluster.AsynCall(GetHallSvrName(svrid), chanAsynRet, data, cb)
}
