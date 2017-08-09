package cluster

import (
	"fmt"
	"mj/common/msg"

	"github.com/lovelly/leaf/log"

	"github.com/lovelly/leaf/chanrpc"
)

const (
	HallPrefixFmt = "HallSvr_%d" //房间服
	GamePrefixFmt = "GameSvr_%d" //游戏服
	HallPrefix    = "HallSvr"    //房间服
	GamePrefix    = "GameSvr"    //游戏服
)

func GetGameSvrName(sververId int) string {
	return fmt.Sprintf(GamePrefixFmt, sververId)
}
func GetHallSvrName(sververId int) string {
	return fmt.Sprintf(HallPrefixFmt, sververId)
}

//发消息给大厅服务器上的玩家
func SendMsgToHallUser(svrid int, uid int64, data interface{}) {
	bdate, err := Processor.Marshal(data)
	if err != nil {
		log.Error("at SendDataToHallUser error:%s", err.Error())
		return
	}

	Go(GetHallSvrName(svrid), &msg.S2S_HanldeFromUserMsg{Uid: uid, Data: bdate[0]})
}

//发送消息给游戏服
func SendMsgToGame(svrid int, data interface{}) {
	Go(GetGameSvrName(svrid), data)
}

func BroadcastToGame(data interface{}) {
	Broadcast(GamePrefix, data)
}

func AsynCallGame(svrid int, chanAsynRet chan *chanrpc.RetInfo, data interface{}, cb interface{}) {
	AsynCall(GetGameSvrName(svrid), chanAsynRet, data, cb)
}

//发消息给大厅
func SendMsgToHall(svrid int, data interface{}) {
	Go(GetHallSvrName(svrid), data)
}

func AsynCallHall(svrid int, chanAsynRet chan *chanrpc.RetInfo, data interface{}, cb interface{}) {
	AsynCall(GetHallSvrName(svrid), chanAsynRet, data, cb)
}

func Call1GameSvr(svrid int, data interface{}) (interface{}, error) {
	return Call1(GetGameSvrName(svrid), data)
}

func BroadcastToHall(data interface{}) {
	Broadcast(HallPrefix, data)
}
