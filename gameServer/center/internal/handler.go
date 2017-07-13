package internal

import (
	"mj/common/cost"
	"mj/gameServer/conf"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/log"
)

//中心模块 ， 投递消息给别的玩家， 或者别的服务器上的玩家
func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.Register(id, f)
}

func init() {
	handleRpc("SelfNodeAddPlayer", SelfNodeAddPlayer)         //暂时无效
	handleRpc("SelfNodeDelPlayer", SelfNodeDelPlayer)         //暂时无效
	handleRpc("NotifyOtherNodeLogin", NotifyOtherNodeLogin)   //暂时无效
	handleRpc("NotifyOtherNodelogout", NotifyOtherNodelogout) //暂时无效
	handleRpc("SendMsgToUser", GoMsgToUser)                   //暂时无效
}

//玩家在本服节点登录
func SelfNodeAddPlayer(args []interface{}) {
	log.Debug("at SelfNodeAddPlayer %v", args)
	uid := args[0].(int64)
	ch := args[1].(*chanrpc.Server)
	Users[uid] = ch
	cluster.Broadcast(cost.GamePrefix, "NotifyOtherNodeLogin", uid, conf.ServerName())
}

//本服玩家登出
func SelfNodeDelPlayer(args []interface{}) {
	log.Debug("at SelfNodeDelPlayer %v", args)
	uid := args[0].(int64)
	delete(Users, uid)
	cluster.Broadcast(cost.GamePrefix, "NotifyOtherNodelogout", uid)
}

//玩家在别的节点登录了
func NotifyOtherNodeLogin(args []interface{}) {
	log.Debug("at NotifyOtherNodeLogin %v", args)
	uid := args[0].(int64)
	ServerName := args[1].(string)
	OtherUsers[uid] = ServerName
}

//玩家在别的节点登出了
func NotifyOtherNodelogout(args []interface{}) {
	log.Debug("at NotifyOtherNodelogout %v", args)
	uid := args[0].(int64)
	delete(OtherUsers, uid)
}

//发消息给别的玩家
func GoMsgToUser(args []interface{}) {
	uid := args[0].(int64)
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
