package pk_base

import (
	"mj/common/msg"
	"mj/common/utils"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"strconv"

	"mj/common/msg/mj_hz_msg"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/timer"
)

type DataManager interface {
	InitRoom(UserCnt int)
	StartDispatchCard(UserManager, LogicManager, *base.GameServiceOption)
	SendGameStart(gameLogic LogicManager, userMgr UserManager)
	SendPersonalTableTip(*user.User, TimerManager)
	SendStatusPlay(u *user.User, userMgr UserManager, gameLogic LogicManager, timerMgr TimerManager)
	NotifySendCard(u *user.User, cbCardData int, userMgr UserManager, bSysOut bool)
	EstimateUserRespond(int, int, int, UserManager, LogicManager) bool
	DispatchCardData(int, UserManager, LogicManager, bool) int
	HasOperator(ChairId, OperateCode int) bool
	HasCard(ChairId, cardIdx int) bool
	CheckUserOperator(*user.User, int, *mj_hz_msg.C2G_HZMJ_OperateCard, LogicManager) (int, int)
	UserChiHu(wTargetUser, userCnt int, gameLogic LogicManager)
	WeaveCard(cbTargetAction, wTargetUser int)
	RemoveCardByOP(wTargetUser, ChoOp int, gameLogic LogicManager) bool
	CallOperateResult(wTargetUser, cbTargetAction int, userMgr UserManager, gameLogic LogicManager)
	ZiMo(u *user.User, gameLogic LogicManager)
	AnGang(u *user.User, cbOperateCode int, cbOperateCard []int, userMgr UserManager, gameLogic LogicManager) int
	NormalEnd(userMgr UserManager, gameLogic LogicManager, template *base.GameServiceOption)
	DismissEnd(userMgr UserManager, gameLogic LogicManager)
	GetTrusteeOutCard(wChairID int, gameLogic LogicManager) int
	CanOperatorRoom(uid int) bool
	SendStatusReady(u *user.User, timerMgr TimerManager)

	GetResumeUser() int
	GetGangStatus() int
	GetUserCardIndex(ChairId int) []int
	GetCurrentUser() int //当前出牌用户
	GetRoomId() int
	GetProvideUser() int
	IsActionDone() bool
}

type BaseManager interface {
	Destroy(int)
	RoomRun(int)
	GetSkeleton() *module.Skeleton
	GetChanRPC() *chanrpc.Server
}

type UserManager interface {
	Sit(*user.User, int) int
	Standup(*user.User) bool
	ForEachUser(fn func(*user.User))
	LeaveRoom(*user.User) bool
	SetUsetStatus(*user.User, int)
	ReLogin(*user.User, int)
	IsAllReady() bool
	AddKickOutTimer(int, *timer.Timer)
	SendUserInfoToSelf(*user.User)
	SendMsgAll(data interface{})
	SendMsgAllNoSelf(selfid int, data interface{})
	WriteTableScore(source []*msg.TagScoreInfo, usercnt, Type int)

	GetCurPlayerCnt() int
	GetMaxPlayerCnt() int
	GetUserInfoByChairId(int) interface{}
	GetUserByChairId(int) *user.User
	GetUserByUid(userId int) (*user.User, int)
	SetUsetTrustee(chairId int, isTruste bool)
	IsTrustee(chairId int) bool
	GetTrustees() []bool
}

type LogicManager interface {
	AnalyseTingCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbOutCardData, cbHuCardCount []int, cbHuCardData [][]int) int
	GetCardCount(cbCardIndex []int) int
	RemoveCard(cbCardIndex []int, cbRemoveCard int) bool
	RemoveCardByArr(cbCardIndex, cbRemoveCard []int) bool
	EstimatePengCard(cbCardIndex []int, cbCurrentCard int) int
	EstimateGangCard(cbCardIndex []int, cbCurrentCard int) int
	GetUserActionRank(cbUserAction int) int
	RandCardList(cbCardBuffer, OriDataArray []int)
	GetUserCards(cbCardIndex []int) (cbCardData []int)
	SwitchToCardData(cbCardIndex int) int
	SwitchToCardIndex(cbCardData int) int
	IsValidCard(card int) bool

	GetMagicIndex() int
	SetMagicIndex(int)
}

type TimerManager interface {
	StartCreatorTimer(Skeleton *module.Skeleton, cb func())
	StartPlayingTimer(Skeleton *module.Skeleton, cb func())

	GetCountLimit() int
	GetTimeLimit() int
	GetPlayCount() int
	AddPlayCount()
	GetMaxPayCnt() int
	GetCreatrTime() int64
	GetTimeOutCard() int
	GetTimeOperateCard() int
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
