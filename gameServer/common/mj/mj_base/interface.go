package mj_base

import (
	"mj/common/msg"
	"mj/common/utils"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"strconv"
)

type DataManager interface {
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
	CheckUserOperator(*user.User, int, int, []int) (int, int)
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
	GetChaHua(u *user.User, setCount int)
	OnUserReplaceCard(u *user.User, CardData int) bool
	OnUserListenCard(u *user.User, bListenCard bool) bool
	RecordFollowCard(cbCenterCard int) bool

	GetResumeUser() int
	GetGangStatus() int
	GetUserCardIndex(ChairId int) []int
	GetCurrentUser() int //当前出牌用户
	GetRoomId() int
	GetProvideUser() int
	IsActionDone() bool

	//测试专用函数。 请勿用于生产
	SetUserCard(charirID int, cards []int)
}

type LogicManager interface {
	AnalyseTingCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbOutCardData, cbHuCardCount []int, cbHuCardData [][]int, MaxCount int) int
	GetCardCount(cbCardIndex []int) int
	RemoveCard(cbCardIndex []int, cbRemoveCard int) bool
	RemoveCardByArr(cbCardIndex, cbRemoveCard []int) bool
	EstimatePengCard(cbCardIndex []int, cbCurrentCard int) int
	EstimateGangCard(cbCardIndex []int, cbCurrentCard int) int
	EstimateEatCard(cbCardIndex []int, cbCurrentCard int) int
	GetUserActionRank(cbUserAction int) int
	AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbCurrentCard int, ChiHuRight int, MaxCount int, b4HZHu bool) (int, []*TagAnalyseItem)
	AnalyseGangCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbProvideCard int, gcr *TagGangCardResult) int
	GetHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbHuCardData []int, MaxCount int) int
	RandCardList(cbCardBuffer, OriDataArray []int)
	GetUserCards(cbCardIndex []int) (cbCardData []int)
	SwitchToCardData(cbCardIndex int) int
	SwitchToCardIndex(cbCardData int) int
	IsValidCard(card int) bool

	GetMagicIndex() int
	SetMagicIndex(int)
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
