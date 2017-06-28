package pk_base

import (
	"mj/common/utils"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"strconv"


)

type DataManager interface {

	InitRoom(UserCnt int)
	GetRoomId() int
	CanOperatorRoom(uid int) bool

	BeforeStartGame(UserCnt int)
	StartGameing()
	AfterStartGame()

	NormalEnd()
	DismissEnd()

	SendPersonalTableTip(*user.User)
	SendStatusPlay(u *user.User)
	SendStatusReady(u *user.User)
/*
BeforeStartGame(UserCnt int)
	StartGameing()
	AfterStartGame()
	SendPersonalTableTip(*user.User)
	SendStatusPlay(u *user.User)
	NotifySendCard(u *user.User, cbCardData int, bSysOut bool)
	EstimateUserRespond(int, int, int) bool
	DispatchCardData(int, bool) int
	HasOperator(ChairId, OperateCode int) bool
	HasCard(ChairId, cardIdx int) bool
	CheckUserOperator(*user.User, int, *mj_hz_msg.C2G_HZMJ_OperateCard) (int, int)
	UserChiHu(wTargetUser, userCnt int)
	WeaveCard(cbTargetAction, wTargetUser int)
	RemoveCardByOP(wTargetUser, ChoOp int) bool
	CallOperateResult(wTargetUser, cbTargetAction int)
	ZiMo(u *user.User)
	AnGang(u *user.User, cbOperateCode int, cbOperateCard []int) int
	NormalEnd()
	DismissEnd()
	GetTrusteeOutCard(wChairID int) int
	CanOperatorRoom(uid int) bool
	SendStatusReady(u *user.User)

	GetResumeUser() int
	GetGangStatus() int
	GetUserCardIndex(ChairId int) []int
	GetCurrentUser() int //当前出牌用户
	GetRoomId() int
	GetProvideUser() int
	IsActionDone() bool

 */

}


type LogicManager interface {

}

////////////////////////////////////////////
//全局变量
// TODO 增加 默认(错误)值 参数
func getGlobalVar(key string) string {
	if globalVar, ok := base.GlobalVarCache.Get(key); ok {
		return globalVar.V
	}
	return ""
}

func getGlobalVarFloat64(key string) float64 {
	if value := getGlobalVar(key); value != "" {
		v, _ := strconv.ParseFloat(value, 10)
		return v
	}
	return 0
}

func getGlobalVarInt64(key string, val int64) int64 {
	if value := getGlobalVar(key); value != "" {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return v
		}
	}
	return val
}

func getGlobalVarInt(key string) int {
	if value := getGlobalVar(key); value != "" {
		v, _ := strconv.Atoi(value)
		return v
	}
	return 0
}

func getGlobalVarIntList(key string) []int {
	if value := getGlobalVar(key); value != "" {
		return utils.GetStrIntList(value)
	}
	return nil
}
