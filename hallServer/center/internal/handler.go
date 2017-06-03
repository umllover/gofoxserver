package internal

import (
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
	//"mj/common/cost"
	"mj/hallServer/conf"
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
func NotifyOtherNodeLogin(args []interface{}){
	uid := args[0].(int)
	ServerName := args[1].(string)
	OtherUsers[uid] = ServerName
}

//玩家在别的节点登出了
func NotifyOtherNodelogout(args []interface{}){
	uid := args[0].(int)
	delete(OtherUsers, uid)
}

//发消息给别的玩家
func GoMsgToUser(args []interface{})  {
	uid := args[0].(int)
	FuncName :=  args[1].(string)
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

//异步回调消息给别的玩家
func AsyncCallUser (args []interface{})  {

}