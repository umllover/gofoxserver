package msg

import (
	//"gopkg.in/mgo.v2/bson"
	"github.com/lovelly/leaf/network/json"
	//"github.com/lovelly/leaf/cluster"
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
	Processor.Register(&L2C_ServerList{})
	Processor.Register(&L2C_ServerListFinish{})

	//game
	Processor.Register(&G2C_LogonSuccess{})
	Processor.Register(&C2G_GR_LogonMobile{})
	Processor.Register(&C2G_GR_UserChairReq{})
	Processor.Register(&C2G_CreateTable{})
	Processor.Register(&C2G_HZOutCard{})
	Processor.Register(&G2C_CreateTableFailure{})
	Processor.Register(&G2C_CreateTableSucess{})
	Processor.Register(&C2G_SearchServerTable{})
	Processor.Register(&G2C_SearchResult{})


}

type ShowErrCode struct {
	ErrCode int
}



