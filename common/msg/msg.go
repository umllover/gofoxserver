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
	Processor.Register(&L2C_UserIndividual{})
	Processor.Register(&C2L_GetRoomList{})
	Processor.Register(&C2L_QuickMatch{})

	Processor.Register(&C2L_ReqCreatorRoomRecord{})
	Processor.Register(&C2L_ReqRoomPlayerBrief{})
	Processor.Register(&L2C_CreatorRoomRecord{})
	Processor.Register(&L2C_RoomPlayerBrief{})

	//game
	Processor.Register(&G2C_LogonFinish{})
	Processor.Register(&G2C_ConfigServer{})
	Processor.Register(&G2C_ConfigFinish{})
	Processor.Register(&G2C_UserEnter{})
	Processor.Register(&C2G_GR_LogonMobile{})
	Processor.Register(&C2G_GR_UserChairReq{})
	Processor.Register(&C2L_CreateTable{})
	Processor.Register(&L2C_CreateTableFailure{})
	Processor.Register(&L2C_CreateTableSucess{})
	Processor.Register(&C2L_SearchServerTable{})
	Processor.Register(&L2C_SearchResult{})
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
	Processor.Register(&C2G_HostlDissumeRoom{})
	Processor.Register(&G2C_CancelTable{})
	Processor.Register(&G2C_PersonalTableEnd{})
	Processor.Register(&C2G_LoadRoom{})
	Processor.Register(&G2C_LoadRoomOk{})

	//chat
	Processor.Register(&C2G_GameChart_ToAll{})
	Processor.Register(&G2C_GameChart_ToAll{})

}

type ShowErrCode struct {
	ErrorCode      int
	DescribeString string
}
