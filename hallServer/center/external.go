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

//发送消息给本服服务器上的玩家
func SendToThisNodeUser(uid int, funcName string, data interface{}) {
	ChanRPC.Go("SendMsgToSelfNotdeUser", uid, funcName, data)
}

//发送消息给其他hallSvr服务器上的玩家
func SendToUser(uid int, funcName string, data interface{}) {
	ChanRPC.Go("SendMsgToUser", uid, funcName, data)
}

//发送消息给游戏服
func SendMsgToGame(svrid int, funcName string, data ...interface{}) {
	cluster.Go(GetGameSvrName(svrid), funcName, data...)
}

func BroadcastToGame(funcName string, data ...interface{}) {
	cluster.Broadcast(GamePrefix, funcName, data...)
}

func AsynCallGame(svrid int, chanAsynRet chan *chanrpc.RetInfo, funcName string, data ...interface{}) {
	cluster.AsynCall(GetGameSvrName(svrid), chanAsynRet, funcName, data...)
}

//发消息给大厅
func SendMsgToHall(svrid int, funcName string, data ...interface{}) {
	cluster.Go(GetHallSvrName(svrid), funcName, data...)
}

func BroadcastToHall(funcName string, data ...interface{}) {
	cluster.Broadcast(HallPrefix, funcName, data...)
}

func AsynCallHall(svrid int, chanAsynRet chan *chanrpc.RetInfo, funcName string, data ...interface{}) {
	cluster.AsynCall(GetHallSvrName(svrid), chanAsynRet, funcName, data...)
}
