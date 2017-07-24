package internal

import (
	"mj/common/msg"
	"mj/common/msg/pk_ddz_msg"
	"mj/gameServer/user"
	"reflect"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
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
	handlerC2S(&pk_ddz_msg.C2G_DDZ_CallScore{}, CallScore)
	handlerC2S(&pk_ddz_msg.C2G_DDZ_OutCard{}, OutCard)
	handlerC2S(&pk_ddz_msg.C2G_DDZ_TRUSTEE{}, TRUSTEE)
	//handlerC2S(pk_ddz_msg.C2G_DDZ_SHOWCARD{}, ShowCard)
}

// 用户叫分
func CallScore(args []interface{}) {
	log.Debug("接受到客户端叫分信息")
	recvMsg := args[0].(*pk_ddz_msg.C2G_DDZ_CallScore)
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("CallScore", recvMsg, user)
	}
}

// 用户出牌
func OutCard(args []interface{}) {
	log.Debug("接受到客户端出牌信息")
	recvMsg := args[0].(*pk_ddz_msg.C2G_DDZ_OutCard)
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("OutCard", recvMsg, user)
	}
}

// 用户托管
func TRUSTEE(args []interface{}) {
	log.Debug("接受到客户端托管信息")
	recvMsg := args[0].(*pk_ddz_msg.C2G_DDZ_TRUSTEE)
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("Trustee", recvMsg, user)
	}
}

// 明牌
func ShowCard(args []interface{}) {
	log.Debug("接受到客户端明牌信息")
	recvMsg := args[0].(*pk_ddz_msg.C2G_DDZ_SHOWCARD)
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("ShowCard", recvMsg, user)
	}
}
