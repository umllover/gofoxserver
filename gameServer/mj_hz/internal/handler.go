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
	reg.RegisterC2S(&mj_hz_msg.C2G_HZMJ_HZOutCard{}, HZOutCard)
	reg.RegisterC2S(&mj_hz_msg.C2G_HZMJ_OperateCard{}, OperateCard)

}

func HZOutCard(args []interface{}) {
	recvMsg := args[0].(*mj_hz_msg.C2G_HZMJ_HZOutCard)
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OutCard", user, recvMsg.CardData)
	}

}

func OperateCard(args []interface{}) {
	recvMsg := args[0].(*mj_hz_msg.C2G_HZMJ_OperateCard)
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OperateCard", user, recvMsg.OperateCode, recvMsg.OperateCard)
	}
}
