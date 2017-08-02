package internal

import (
	"errors"
	"mj/common/cost"
	"mj/common/msg"
	"mj/common/register"
	"mj/gameServer/RoomMgr"
	"mj/gameServer/conf"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	reg.RegisterRpc("SelfNodeAddPlayer", SelfNodeAddPlayer)
	reg.RegisterRpc("SelfNodeDelPlayer", SelfNodeDelPlayer)
	reg.RegisterRpc("SendMsgToSelfNotdeUser", SendMsgToSelfNotdeUser)
	reg.RegisterRpc("HanldeFromHallMsg", HanldeFromHallMsg)
	reg.RegisterRpc("ServerFaild", serverFaild)
	reg.RegisterRpc("ServerStart", serverStart)

	reg.RegisterS2S(&msg.S2S_NotifyOtherNodeLogin{}, NotifyOtherNodeLogin)
	reg.RegisterS2S(&msg.S2S_NotifyOtherNodelogout{}, NotifyOtherNodelogout)

	// 登录服发来的协议
	reg.RegisterS2S(&msg.S2S_CloseRoom{}, SREQCloseRoom)
	reg.RegisterS2S(&msg.L2G_CreatorRoom{}, CreatorRoom)
}

//玩家在本服节点登录
func SelfNodeAddPlayer(args []interface{}) {
	uid := args[0].(int64)
	ch := args[1].(*chanrpc.Server)
	Users[uid] = ch
	cluster.Broadcast(cost.GamePrefix, &msg.S2S_NotifyOtherNodeLogin{
		Uid:        uid,
		ServerName: conf.ServerName(),
	})
}

//本服玩家登出
func SelfNodeDelPlayer(args []interface{}) {
	uid := args[0].(int64)
	delete(Users, uid)
	cluster.Broadcast(cost.GamePrefix, &msg.S2S_NotifyOtherNodelogout{
		Uid: uid,
	})
}

//玩家在别的节点登录了
func NotifyOtherNodeLogin(args []interface{}) {
	recvMsg := args[0].(*msg.S2S_NotifyOtherNodeLogin)
	OtherUsers[recvMsg.Uid] = recvMsg.ServerName
	log.Debug("user %d login on %s", recvMsg.Uid, recvMsg.ServerName)
}

//玩家在别的节点登出了
func NotifyOtherNodelogout(args []interface{}) {
	recvMsg := args[0].(*msg.S2S_NotifyOtherNodelogout)
	log.Debug("user %d logout on %s", recvMsg.Uid, OtherUsers[recvMsg.Uid])
	delete(OtherUsers, recvMsg.Uid)
}

//处理来自游戏服的消息
func HanldeFromHallMsg(args []interface{}) {
	SendMsgToSelfNotdeUser(args)
}

//登录服发来的删除房间协议
func SREQCloseRoom(args []interface{}) (interface{}, error) {
	recvMsg := args[0].(*msg.S2S_CloseRoom)
	room := RoomMgr.GetRoom(recvMsg.RoomID)
	if room == nil {
		return nil, errors.New("not foud room")
	}

	room.GetChanRPC().Go("DissumeRoom", nil)
	return nil, nil
}

//大厅服发来的创建房间
func CreatorRoom(args []interface{}) (interface{}, error) {
	if LoadRoom(args[0]) {
		return nil, nil
	} else {
		return nil, errors.New("creator room faild ")
	}
}
