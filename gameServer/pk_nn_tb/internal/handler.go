package internal

import (
	"mj/common/msg/nn_tb_msg"
	"mj/gameServer/user"

	"mj/common/register"

	"github.com/lovelly/leaf/gate"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	// c 2 s
	reg.RegisterC2S(&nn_tb_msg.C2G_TBNN_CallScore{}, TBNNCallScore)
	reg.RegisterC2S(&nn_tb_msg.C2G_TBNN_AddScore{}, TBNNAddScore)
	//reg.RegisterC2S((&nn_tb_msg.C2G_TBNN_CallBanker{}, TBNNCallBanker)
	reg.RegisterC2S(&nn_tb_msg.C2G_TBNN_OpenCard{}, TBNNOpenCard)
	//reg.RegisterC2S((&nn_tb_msg.C2G_TBNN_QIANG{}, TBNNQiang)
}

func TBNNCallScore(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("CallScore", args[0], user)
	}

}
func TBNNAddScore(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("AddScore", args[0], user)
	}

}

/*
func TBNNCallBanker(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("CallBanker", args[0], user)
	}

}*/

func TBNNOpenCard(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OpenCard", args[0], user)
	}

}

/*
func TBNNQiang(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("Qiang", args[0], user)
	}
}
*/
