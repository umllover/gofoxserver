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
	handlerC2S(&mj_zp_msg.C2G_ZPMJ_OperateCard{}, OperateCard)
	handlerC2S(&mj_zp_msg.C2G_MJZP_SetChaHua{}, SetChaHua)
	handlerC2S(&mj_zp_msg.G2C_MJZP_ReplaceCard{}, SetBuHua)
	handlerC2S(&mj_zp_msg.C2G_MJZP_ListenCard{}, SetTingCard)
	handlerC2S(&mj_zp_msg.C2G_MJZP_Trustee{}, Trustee)
}

func ZPOutCard(args []interface{}) {
	recvMsg := args[0].(*mj_zp_msg.C2G_ZPMJ_OutCard)
	agent := args[1].(gate.Agent)
	u := agent.UserData().(*user.User)

	r := getRoom(u.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OutCard", u, recvMsg.CardData)
	}
}

func OperateCard(args []interface{}) {
	recvMsg := args[0].(*mj_zp_msg.C2G_ZPMJ_OperateCard)
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OperateCard", user, recvMsg.OperateCode, recvMsg.OperateCard)
	}
}

//插花
func SetChaHua(args []interface{}) {
	agent := args[1].(gate.Agent)
	u := agent.UserData().(*user.User)

	r := getRoom(u.RoomId)
	if r != nil {
		r.GetChanRPC().Go("SetChaHua", args[0], u)
	}
}

//补花
func SetBuHua(args []interface{}) {
	agent := args[1].(gate.Agent)
	u := agent.UserData().(*user.User)

	r := getRoom(u.RoomId)
	if r != nil {
		r.GetChanRPC().Go("SetBuHua", args[0], u)
	}
}

//听牌
func SetTingCard(args []interface{}) {
	agent := args[1].(gate.Agent)
	u := agent.UserData().(*user.User)

	r := getRoom(u.RoomId)
	if r != nil {
		r.GetChanRPC().Go("SetTingCard", args[0], u)
	}
}

//托管
func Trustee(args []interface{}) {
	agent := args[1].(gate.Agent)
	u := agent.UserData().(*user.User)

	r := getRoom(u.RoomId)
	if r != nil {
		r.GetChanRPC().Go("UserTrustee", args[0], u)
	}
}
