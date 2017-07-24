package internal

import (
	"errors"
	"mj/common/consul"
	"mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/conf"
	"mj/gameServer/user"

	"mj/common/register"

	"mj/gameServer/RoomMgr"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	reg.RegisterRpc("SelfNodeAddPlayer", SelfNodeAddPlayer)
	reg.RegisterRpc("SelfNodeDelPlayer", SelfNodeDelPlayer)
	reg.RegisterRpc("SendMsgToSelfNotdeUser", SendMsgToSelfNotdeUser)
	reg.RegisterRpc("HanldeFromGameMsg", HanldeFromGameMsg)
	reg.RegisterRpc("ServerFaild", serverFaild)
	reg.RegisterRpc("ServerStart", serverStart)

	reg.RegisterS2S(&msg.S2S_GetPlayerInfo{}, GetPlayerInfo)
	reg.RegisterS2S(&msg.S2S_NotifyOtherNodeLogin{}, NotifyOtherNodeLogin)
	reg.RegisterS2S(&msg.S2S_NotifyOtherNodelogout{}, NotifyOtherNodelogout)

	// 登录服发来的协议
	reg.RegisterS2S(&msg.S2S_CloseRoom{}, SREQCloseRoom)
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
func HanldeFromGameMsg(args []interface{}) {
	SendMsgToSelfNotdeUser(args)
}

func GetPlayerInfo(args []interface{}) (interface{}, error) {
	uid := args[0].(int64)
	log.Debug("at GetPlayerInfo uid:%d", uid)
	ch, chok := Users[uid]
	if !chok {
		return nil, errors.New("not foud user ch")
	}
	us, err := ch.TimeOutCall1("GetUser", 5)
	if err != nil {
		return nil, err
	}

	u, ok := us.(*user.User)
	if !ok {
		return nil, errors.New("user data error")
	}

	gu := map[string]interface{}{
		"Id":          u.Id,
		"NickName":    u.NickName,
		"Currency":    u.Currency,
		"RoomCard":    u.RoomCard,
		"FaceID":      u.FaceID,
		"CustomID":    u.CustomID,
		"HeadImgUrl":  u.HeadImgUrl,
		"Experience":  u.Experience,
		"Gender":      u.Gender,
		"WinCount":    u.WinCount,
		"LostCount":   u.LostCount,
		"DrawCount":   u.DrawCount,
		"FleeCount":   u.FleeCount,
		"UserRight":   u.UserRight,
		"Score":       u.Score,
		"Revenue":     u.Revenue,
		"InsureScore": u.InsureScore,
		"MemberOrder": u.MemberOrder,
	}
	return gu, nil
}

func SREQCloseRoom(args []interface{}) (interface{}, error) {
	recvMsg := args[0].(*msg.S2S_CloseRoom)
	room := RoomMgr.GetRoom(recvMsg.RoomID)
	if room == nil {
		return nil, errors.New("not foud room")
	}

	room.GetChanRPC().Go("DissumeRoom", nil)
	return nil, nil
}

//新的节点启动了
func serverStart(args []interface{}) {
	svr := args[0].(*consul.CacheInfo)
	log.Debug("%s on line", svr.Csid)
	cluster.AddClient(&cluster.NsqClient{Addr: svr.Host, ServerName: svr.Csid})
}

//节点关闭了
func serverFaild(args []interface{}) {
	svr := args[0].(*consul.CacheInfo)
	log.Debug("%s off line", svr.Csid)
	cluster.RemoveClient(svr.Csid)
}
