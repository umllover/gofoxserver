package internal

import (
	"mj/common/msg"
	"mj/common/msg/mj_zp_msg"
	"mj/gameServer/user"
	"reflect"

	"github.com/lovelly/leaf/gate"
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}) {
	ChanRPC.Register(id, f)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	// c 2 s
	handlerC2S(&mj_zp_msg.C2G_ZPMJ_OutCard{}, ZPOutCard)
	//handlerC2S(&mj_zp_msg.C2G_ZPMJ_OperateCard{}, OperateCard)
	//handlerC2S(&mj_zp_msg.C2G_ZPMJ_OperateCard{}, room.)
}

func ZPOutCard(args []interface{}) {
	recvMsg := args[0].(*mj_zp_msg.C2G_ZPMJ_OutCard)
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OutCard", user, recvMsg.CardData)
	}
}

func OperateCard(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OperateCard", args[0], user)
	}
}
