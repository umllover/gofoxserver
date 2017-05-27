package internal

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/db/model"
	"mj/gameServer/user"
	"reflect"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/gate"
	"mj/hallServer/UserData"
)

var userDatach = UserData.ChanRPC

////注册rpc 消息
func handleRpc(id interface{}, f interface{}, fType int) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.RegisterFromType(id, f, fType)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {

	handlerC2S(&msg.C2L_Login{}, handleMBLogin)
}


func handleMBLogin(args []interface{}) {

}



///////
func loadUser(UserID int) (*user.User, bool){

}

func createUser(UserID int)  (*user.User, bool) {
	U := user.NewUser(UserID)


	return U, true
}

func BuildClientMsg(retMsg *msg.CMD_MB_LogonSuccess, user *user.User){

}
