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
	Processor.Register(&L2C_ServerList{})
	Processor.Register(&L2C_ServerListFinish{})
	Processor.Register(&C2L_User_Individual{})
	Processor.Register(&L2C_UserIndividual{})
	Processor.Register(&C2L_GetRoomList{})
	Processor.Register(&C2L_QuickMatch{})
	Processor.Register(&L2C_CreateTableFailure{})

	Processor.Register(&C2L_ReqCreatorRoomRecord{})
	Processor.Register(&C2L_ReqRoomPlayerBrief{})
	Processor.Register(&L2C_CreatorRoomRecord{})
	Processor.Register(&L2C_RoomPlayerBrief{})

	Processor.Register(&C2L_DrawSahreAward{})
	Processor.Register(&L2C_DrawSahreAwardResult{})
	Processor.Register(&L2C_ActivityInfo{})
	Processor.Register(&C2L_SetElect{})
	Processor.Register(&L2C_SetElectResult{})
	Processor.Register(&L2C_RspTradeShopInfo{})
	Processor.Register(&L2C_GetRoomList{})
	Processor.Register(&L2C_QuickMatchOk{})
	Processor.Register(&C2L_DeleteRoom{})
	Processor.Register(&L2C_DeleteRoomResult{})
	Processor.Register(&C2L_ReqBindMaskCode{})
	Processor.Register(&L2C_ReqBindMaskCodeRsp{})
	Processor.Register(&C2L_SetPhoneNumber{})
	Processor.Register(&L2C_SetPhoneNumberRsp{})
	Processor.Register(&C2L_DianZhan{})
	Processor.Register(&L2C_DianZhanRsp{})

	Processor.Register(&C2L_RenewalFees{})
	Processor.Register(&L2C_RenewalFeesRsp{})
	Processor.Register(&G2C_RenewalFeesSuccess{})
	Processor.Register(&C2L_ChangeUserName{})
	Processor.Register(&L2C_ChangeUserNameRsp{})
	Processor.Register(&C2L_ChangeSign{})
	Processor.Register(&L2C_ChangeSignRsp{})
	Processor.Register(&L2C_KickOut{})
	Processor.Register(&L2C_UpdateUserAttr{})
	Processor.Register(&C2L_RechangerOk{})
	Processor.Register(&L2C_BeStar{})
	Processor.Register(&L2C_NotifyElectResult{})
	Processor.Register(&L2C_RechangerOk{})
	Processor.Register(&L2C_GMNotice{})
	Processor.Register(&C2L_GetRoomRecord{})
	Processor.Register(&L2C_RoomRecord{})
	Processor.Register(&C2L_GetUserRecords{})
	Processor.Register(&L2C_GetUserRecords{})

	//game
	Processor.Register(&G2C_LogonFinish{})
	Processor.Register(&G2C_ConfigServer{})
	Processor.Register(&G2C_ConfigFinish{})
	Processor.Register(&G2C_UserEnter{})
	Processor.Register(&C2G_GR_LogonMobile{})
	Processor.Register(&C2G_GR_UserChairReq{})
	Processor.Register(&C2L_CreateTable{})
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
	Processor.Register(&G2C_LogonFailure{})
	Processor.Register(&C2G_UserStandup{})
	Processor.Register(&C2G_UserReady{})
	Processor.Register(&G2C_Hu_Data{})
	Processor.Register(&SysMsg{})
	Processor.Register(&G2C_CancelTable{})
	Processor.Register(&G2C_PersonalTableEnd{})
	Processor.Register(&G2C_GameConclude{})
	Processor.Register(&G2C_UserSitDownRst{})
	Processor.Register(&G2C_KickOut{})
	Processor.Register(&C2G_LeaveRoom{})
	Processor.Register(&G2C_LeaveRoomRsp{})
	Processor.Register(&G2C_LeaveRoomBradcast{})
	Processor.Register(&G2C_ReplyRsp{})
	Processor.Register(&C2G_ReplyLeaveRoom{})
	Processor.Register(&C2L_ReqTimesInfo{})
	Processor.Register(&C2L_TimeSync{})
	Processor.Register(&L2C_TimeSync{})
	Processor.Register(&C2L_ReConnect{})
	Processor.Register(&L2C_ReConnectRsp{})
	Processor.Register(&G2C_RoomDissume{})
	Processor.Register(&G2C_OvertimeToNext{})

	//chat
	Processor.Register(&C2G_GameChart_ToAll{})
	Processor.Register(&G2C_GameChart_ToAll{})

}

type ShowErrCode struct {
	ErrorCode      int
	DescribeString string
}
