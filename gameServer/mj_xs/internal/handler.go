package internal

import (
	"mj/common/msg/mj_xs_msg"
	"mj/gameServer/user"

	"mj/common/register"

	"github.com/lovelly/leaf/gate"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	// c 2 s
	reg.RegisterC2S(&mj_xs_msg.C2G_MJXS_OutCard{}, HZOutCard)
	reg.RegisterC2S(&mj_xs_msg.C2G_MJXS_OperateCard{}, OperateCard)
	reg.RegisterC2S(&mj_xs_msg.C2G_MJXS_ReplaceCard{}, ReplaceCard)
}

func HZOutCard(args []interface{}) {
	recvMsg := args[0].(*mj_xs_msg.C2G_MJXS_OutCard)
	agent := args[1].(gate.Agent)
	player := agent.UserData().(*user.User)

	r := getRoom(player.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OutCard", player, recvMsg.CardData)
	}
}

func OperateCard(args []interface{}) {
	recvMsg := args[0].(*mj_xs_msg.C2G_MJXS_OperateCard)
	agent := args[1].(gate.Agent)
	player := agent.UserData().(*user.User)

	r := getRoom(player.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OperateCard", player, recvMsg.OperateCode, recvMsg.OperateCard)
	}
}
func ReplaceCard(args []interface{}) {
	recvMsg := args[0].(*mj_xs_msg.C2G_MJXS_ReplaceCard)
	agent := args[1].(gate.Agent)
	player := agent.UserData().(*user.User)

	r := getRoom(player.RoomId)
	if r != nil {
		r.GetChanRPC().Go("SetBuHua", player, recvMsg.CardData)
	}
}
