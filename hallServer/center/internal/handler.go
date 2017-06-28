package internal

import (
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
	//"mj/common/cost"
	"mj/hallServer/conf"

	"errors"

	"mj/hallServer/user"

	"github.com/lovelly/leaf/log"
)

//中心模块 ， 投递消息给别的玩家， 或者别的服务器上的玩家
func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.Register(id, f)
}

func init() {
	handleRpc("SelfNodeAddPlayer", SelfNodeAddPlayer)
	handleRpc("SelfNodeDelPlayer", SelfNodeDelPlayer)
	handleRpc("NotifyOtherNodeLogin", NotifyOtherNodeLogin)
	handleRpc("NotifyOtherNodelogout", NotifyOtherNodelogout)
	handleRpc("SendMsgToUser", GoMsgToUser)
	handleRpc("AsyncCallUser", AsyncCallUser)
	handleRpc("GetPlayerInfo", GetPlayerInfo)
	handleRpc("SendMsgToSelfNotdeUser", SendMsgToSelfNotdeUser)
	handleRpc("HanldeFromGameMsg", HanldeFromGameMsg)
}

//玩家在本服节点登录
func SelfNodeAddPlayer(args []interface{}) {
	uid := args[0].(int)
	ch := args[1].(*chanrpc.Server)
	Users[uid] = ch
	//cluster.Broadcast(cost.HallPrefix,"NotifyOtherNodeLogin", uid, conf.ServerName())
}

//本服玩家登出
func SelfNodeDelPlayer(args []interface{}) {
	uid := args[0].(int)
	delete(Users, uid)
	//cluster.Broadcast(cost.HallPrefix,"NotifyOtherNodelogout", uid)
}

//玩家在别的节点登录了
func NotifyOtherNodeLogin(args []interface{}) {
	uid := args[0].(int)
	ServerName := args[1].(string)
	OtherUsers[uid] = ServerName
}

//玩家在别的节点登出了
func NotifyOtherNodelogout(args []interface{}) {
	uid := args[0].(int)
	delete(OtherUsers, uid)
}

//发消息给别的玩家
func GoMsgToUser(args []interface{}) {
	uid := args[0].(int)
	FuncName := args[1].(string)
	ch, ok := Users[uid]
	if ok {
		ch.Go(FuncName, args[2:]...)
		return
	}

	ServerName, ok1 := OtherUsers[uid]
	if ServerName == conf.ServerName() {
		log.Error("self server user not login .... ")
		return
	}

	if ok1 {
		cluster.Go(ServerName, "SendMsgToUser", args...)
	}
}

func SendMsgToSelfNotdeUser(args []interface{}) {
	uid := args[0].(int)
	FuncName := args[1].(string)
	ch, ok := Users[uid]
	if ok {
		ch.Go(FuncName, args[2:]...)
		return
	}
	log.Debug("at SendMsgToSelfNotdeUser player not in node")
}

//处理来自游戏服的消息
func HanldeFromGameMsg(args []interface{}) {
	SendMsgToSelfNotdeUser(args)
}

//异步回调消息给别的玩家
func AsyncCallUser(args []interface{}) {

}

func GetPlayerInfo(args []interface{}) (interface{}, error) {
	uid := args[0].(int)
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
		"UserRight":   u.Accountsmember.UserRight,
		"Score":       u.Score,
		"Revenue":     u.Revenue,
		"InsureScore": u.InsureScore,
		"MemberOrder": u.MemberOrder,
	}
	return gu, nil
}
