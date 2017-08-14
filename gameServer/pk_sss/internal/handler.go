package internal

import (
	"mj/common/msg"
	"mj/common/msg/pk_sss_msg"
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
	handlerC2S(&pk_sss_msg.C2G_SSS_Open_Card{}, SSSShowCard)
	handlerC2S(&pk_sss_msg.C2G_SSS_TRUSTEE{}, TRUSTEE)
}
func SSSShowCard(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("ShowCard", args[0], user)
	}

}

// 用户托管
func TRUSTEE(args []interface{}) {
	recvMsg := args[0].(*pk_sss_msg.C2G_SSS_TRUSTEE)
	agent := args[1].(gate.Agent)
	u := agent.UserData().(*user.User)

	r := getRoom(u.RoomId)
	if r != nil {
		r.GetChanRPC().Go("Trustee", recvMsg, u)
	}
}
