package common

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
	StartGame()
	SendGameStart(gameLogic LogicManager, userMgr UserManager)
	SendPersonalTableTip(*user.User)
	SendStatusReady(*user.User)
	SendStatusPlay(*user.User, int, int, LogicManager)
	OnUserOutCard(int, int, bool) int
	NotifySendCard(u *user.User, cbCardData int, userMgr UserManager, bSysOut bool)
	EstimateUserRespond(wCenterUser int, cbCenterCard int, EstimatKind int, userMgr UserManager) bool
	DispatchCardData(int, bool) bool
	HasOperator(ChairId, OperateCode int) bool
	HasCard(ChairId, cardIdx int) bool
	CheckUserOperator(*user.User, int, *mj_hz_msg.C2G_HZMJ_OperateCard, LogicManager) (int, int)
	UserChiHu(wTargetUser, userCnt int, gameLogic LogicManager)
	WeaveCard(cbTargetAction, wTargetUser int)
	RemoveCardByOP(wTargetUser, ChoOp int, gameLogic LogicManager) bool
	OperateResult(wTargetUser, cbTargetAction int, userMgr UserManager, gameLogic LogicManager)
	ZiMo(u *user.User, gameLogic LogicManager)
	AnGang(u *user.User, cbOperateCode int, cbOperateCard []int, userMgr UserManager, gameLogic LogicManager) int
	NormalEnd(userMgr UserManager, gameLogic LogicManager, template *base.GameServiceOption)
	DismissEnd(userMgr UserManager, gameLogic LogicManager)
	OnUserTrustee(wChairID int, bTrustee bool) bool
	GetResumeUser() int
	GetGangStatus() int
	GetUserCardIndex(ChairId int) []int
	GetMaxPayCnt() int
	GetCurPayInt() int
	GetCreatrTime() int64
	GetCurrentUser() int //当前出牌用户
	GetRoomId() int
	GetCountLimit() int
	GetTimeLimit() int
	GetProvideUser() int
}

type BaseManager interface {
	Destroy(int)
	RoomRun(int)
	Skeleton() *module.Skeleton
	GetChanRPC() *chanrpc.Server
}

type UserManager interface {
	Sit(*user.User, int) int
	Standup(*user.User) bool
	CanOperatorRoom(int) bool
	ForEachUser(fn func(*user.User))
	LeaveRoom(*user.User)
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
}

type LogicManager interface {
	AnalyseTingCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbOutCardData, cbHuCardCount []int, cbHuCardData [][]int) int
	GetCardCount(cbCardIndex []int) int
	SwitchToCardData([]int) []int
	IsValidCard(int) bool
	RemoveCard(cbCardIndex []int, cbRemoveCard int) bool
	RemoveCardByArr(cbCardIndex, cbRemoveCard []int) bool
	EstimatePengCard(cbCardIndex []int, cbCurrentCard int) int
	EstimateGangCard(cbCardIndex []int, cbCurrentCard int) int
	SwitchToCardIndex(cbCardData int) int
	GetUserActionRank(cbUserAction int) int
	AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbCurrentCard int, ChiHuRight int, b4HZHu bool) int
	AnalyseGangCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount, cbProvideCard int, gcr *TagGangCardResult) int
	RandCardList(cbCardBuffer, OriDataArray []int)

	GetMagicIndex() int
	SetMagicIndex(int)
}

type TimerManager interface {
	GetPlayCount() int
	AddPlayCount()
	StartCreatorTimer(skeleton *module.Skeleton, cb func())
	StartPlayingTimer(skeleton *module.Skeleton, cb func())
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
