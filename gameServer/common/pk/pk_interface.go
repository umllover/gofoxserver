package pk

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

	// 游戏开始
	BeforeStartGame(UserCnt int)
	StartGameing()
	AfterStartGame()

	// 游戏结束
	NormalEnd()
	DismissEnd()

	SendPersonalTableTip(*user.User)
	SendStatusPlay(u *user.User)
	SendStatusReady(u *user.User)

	AddScoreTimes(u *user.User, scoreTimes int)
	AddScore(u *user.User, score int)
}

type LogicManager interface {
	RandCardList(cbCardBuffer, OriDataArray []int)
	SortCardList(cardData []int, cardCount int)
	GetCardValue(CardData int) int
	GetCardColor(CardData int) int


	/*
		AnalyseTingCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbOutCardData, cbHuCardCount []int, cbHuCardData [][]int) int
		GetCardCount(cbCardIndex []int) int
		RemoveCard(cbCardIndex []int, cbRemoveCard int) bool
		RemoveCardByArr(cbCardIndex, cbRemoveCard []int) bool
		EstimatePengCard(cbCardIndex []int, cbCurrentCard int) int
		EstimateGangCard(cbCardIndex []int, cbCurrentCard int) int
		GetUserActionRank(cbUserAction int) int
		AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbCurrentCard int, ChiHuRight int, b4HZHu bool) int
		AnalyseGangCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbProvideCard int, gcr *TagGangCardResult) int
		GetUserCards(cbCardIndex []int) (cbCardData []int)
		SwitchToCardData(cbCardIndex int) int
		SwitchToCardIndex(cbCardData int) int
		IsValidCard(card int) bool

		GetMagicIndex() int
		SetMagicIndex(int)

	*/
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
