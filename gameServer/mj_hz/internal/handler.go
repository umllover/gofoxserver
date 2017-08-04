package internal

import (
	"mj/common/msg/mj_hz_msg"
	"mj/common/register"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/gate"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	// c 2 s
	reg.RegisterC2S(&mj_hz_msg.C2G_HZMJ_HZOutCard{}, C2G_OutCard)
	reg.RegisterC2S(&mj_hz_msg.C2G_HZMJ_OperateCard{}, C2G_OperateCard)

}

func C2G_OutCard(args []interface{}) {
	recvMsg := args[0].(*mj_hz_msg.C2G_HZMJ_HZOutCard)
	agent := args[1].(gate.Agent)
	u := agent.UserData().(*user.User)
	r := getRoom(u.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OutCard", u, recvMsg.CardData)
	}
}

func C2G_OperateCard(args []interface{}) {
	recvMsg := args[0].(*mj_hz_msg.C2G_HZMJ_OperateCard)
	agent := args[1].(gate.Agent)
	u := agent.UserData().(*user.User)
	r := getRoom(u.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OperateCard", u, recvMsg.OperateCode, recvMsg.OperateCard)
	}
}
