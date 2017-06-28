package internal

import (
	"mj/common/msg"
	"mj/common/msg/nn_tb_msg"
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
	handlerC2S(&nn_tb_msg.C2G_TBNN_CallScore{}, TBNNCallScore)
	handlerC2S(&nn_tb_msg.C2G_TBNN_AddScore{}, TBNNAddScore)
	handlerC2S(&nn_tb_msg.C2G_TBNN_CallBanker{}, TBNNCallBanker)
	handlerC2S(&nn_tb_msg.C2G_TBNN_OxCard{}, TBNNOxCard)
	handlerC2S(&nn_tb_msg.C2G_TBNN_QIANG{}, TBNNQiang)
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
func TBNNCallBanker(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("CallBanker", args[0], user)
	}

}

func TBNNOxCard(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OxCard", args[0], user)
	}

}

func TBNNQiang(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("Qiang", args[0], user)
	}
}



