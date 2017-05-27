package internal

import (
	"mj/common/msg"
	"mj/hallServer/user"
	//. "mj/common/cost"
	"reflect"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
)

var (
	userMap map[int]*user.User
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}, fType int) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.RegisterFromType(id, f, fType)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	userMap = make(map[int]*user.User)
	handleRpc("addUser", addUser, chanrpc.FuncCommon)
	handleRpc("delUser", delUser, chanrpc.FuncCommon)
	handleRpc("GetUserGameData", GetUserGameData, chanrpc.FuncCommon)
}

func addUser(args []interface{}){
	user :=  args[0].(*user.User)
	userMap[user.Id] = user
}

func delUser(args []interface{}) {
	userId :=  args[0].(int)
	user, ok := userMap[userId]
	if ok {
		delete(userMap, userId)
		_ = user
		//cluster.Go(GetGameSvrName(user.ServerID), "UserOffline", user)
	}
}

//call Func
func GetUserGameData(args []interface{}) (interface{}, error){
	return nil, nil
}




