package internal

import (

	"mj/common/msg"
	"mj/gameServer/hzmj/room"
	"reflect"
	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/gate"
	"mj/gameServer/user"
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
	handlerC2S(&msg.C2G_CreateRoom{}, CreaterRoom)
	handlerC2S(&msg.C2G_HZOutCard{}, HZOutCard)

	handleRpc("DelRoom", DelRoom, chanrpc.FuncCommon)
}


func CreaterRoom(args []interface{}) {
	r  := room.NewRoom(ChanRPC)
	addRoom(r)
}

func HZOutCard (args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.ChanRPC.Go("OutCard", args[0])
	}
}


//////////////// rcp ///////////////////
func DelRoom(args []interface{}){
	id := args[0].(int)
	delRoom(id)
}


