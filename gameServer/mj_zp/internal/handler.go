package internal

import (
	"mj/common/msg/mj_zp_msg"
	"mj/gameServer/user"

	"mj/common/register"

	"github.com/lovelly/leaf/gate"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	// c 2 s
	reg.RegisterC2S(&mj_zp_msg.C2G_ZPMJ_OutCard{}, ZPOutCard)
	reg.RegisterC2S(&mj_zp_msg.C2G_ZPMJ_OperateCard{}, OperateCard)
	reg.RegisterC2S(&mj_zp_msg.C2G_MJZP_SetChaHua{}, SetChaHua)
	reg.RegisterC2S(&mj_zp_msg.C2G_MJZP_ReplaceCard{}, SetBuHua)
	reg.RegisterC2S(&mj_zp_msg.C2G_MJZP_ListenCard{}, SetTingCard)
	reg.RegisterC2S(&mj_zp_msg.C2G_MJZP_Trustee{}, Trustee)
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
	recvMsg := args[0].(*mj_zp_msg.C2G_MJZP_ReplaceCard)
	agent := args[1].(gate.Agent)
	u := agent.UserData().(*user.User)

	r := getRoom(u.RoomId)
	if r != nil {
		r.GetChanRPC().Go("SetBuHua", u, recvMsg.CardData)
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
