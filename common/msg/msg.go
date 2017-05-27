package msg

import (
	//"gopkg.in/mgo.v2/bson"
	"github.com/lovelly/leaf/network/json"
)

var (
	Processor = json.NewProcessor()
)

func init() {
	//hall
	Processor.Register(&C2L_Login{})
	Processor.Register(&ShowErrCode{})
	Processor.Register(&C2L_Regist{})
	Processor.Register(&L2C_LogonFailure{})
	Processor.Register(&L2C_LogonSuccess{})
	Processor.Register(&C2G_CreateRoom{})
	Processor.Register(&C2G_HZOutCard{})
	Processor.Register(&L2C_ServerList{})
	Processor.Register(&L2C_ServerListFinish{})

	//game
	Processor.Register(&G2C_LogonSuccess{})
	Processor.Register(&C2G_GR_LogonMobile{})
}

type ShowErrCode struct {
	ErrCode int
}



