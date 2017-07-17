package internal

import (
	"mj/gameServer/user"

	"github.com/lovelly/leaf/gate"
)

func init() {
	// c 2 s
	//handlerC2S(&mj_hz_msg.C2G_HZMJ_HZOutCard{}, SSSShowCard)
	//handlerC2S(&mj_hz_msg.C2G_HZMJ_OperateCard{}, OperateCard)

}
func SSSShowCard(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.GetChanRPC().Go("CallScore", args[0], user)
	}

}
