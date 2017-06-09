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
	Processor.Register(&C2L_User_Individual{})

	//game
	Processor.Register(&G2C_LogonFinish{})
	Processor.Register(&G2C_ConfigServer{})
	Processor.Register(&G2C_ConfigFinish{})
	Processor.Register(&G2C_UserEnter{})
	Processor.Register(&C2G_GR_LogonMobile{})
	Processor.Register(&C2G_GR_UserChairReq{})
	Processor.Register(&C2G_CreateTable{})
	Processor.Register(&G2C_CreateTableFailure{})
	Processor.Register(&G2C_CreateTableSucess{})
	Processor.Register(&C2L_SearchServerTable{})
	Processor.Register(&G2C_SearchResult{})
	Processor.Register(&C2G_UserSitdown{})
	Processor.Register(&G2C_UserStatus{})
	Processor.Register(&C2G_GameOption{})
	Processor.Register(&G2C_PersonalTableTip{})
	Processor.Register(&G2C_Record{})
	Processor.Register(&G2C_StatusFree{})
	Processor.Register(&G2C_StatusPlay{})
	Processor.Register(&C2G_REQUserInfo{})
	Processor.Register(&G2C_GameStatus{})
	Processor.Register(&C2G_REQUserChairInfo{})
	Processor.Register(&G2C_LogonFailur{})
	Processor.Register(&C2G_UserStandup{})
	Processor.Register(&C2G_UserReady{})
	Processor.Register(&G2C_Hu_Data{})
	Processor.Register(&SysMsg{})

	//HZMJ msg
	Processor.Register(&G2C_HZMG_GameStart{})
	Processor.Register(&G2C_HZMJ_GameConclude{})
	Processor.Register(&C2G_HZMJ_HZOutCard{})
	Processor.Register(&G2C_HZMJ_OutCard{})
	Processor.Register(&G2C_HZMJ_OperateNotify{})
	Processor.Register(&G2C_HZMJ_SendCard{})
	Processor.Register(&C2G_HZMJ_OperateCard{})
	Processor.Register(&G2C_HZMJ_OperateResult{})
	Processor.Register(&G2C_HZMJ_Trustee{})

	//chat
	Processor.Register(&C2G_GameChart_ToAll{})
	Processor.Register(&G2C_GameChart_ToAll{})
}

type ShowErrCode struct {
	ErrorCode      int
	DescribeString string
}
