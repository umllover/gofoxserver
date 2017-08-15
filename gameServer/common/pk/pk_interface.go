package pk

import (
	"mj/gameServer/user"
)

type DataManager interface {
	OnCreateRoom()
	InitRoom(UserCnt int)
	GetRoomId() int
	CanOperatorRoom(uid int64) bool

	// 游戏开始
	BeforeStartGame(UserCnt int)
	StartGameing()
	AfterStartGame()
	ResetGameAfterRenewal()

	// 游戏结束
	NormalEnd(Reason int)
	DismissEnd(Reason int)

	SendPersonalTableTip(*user.User)
	SendStatusPlay(u *user.User)
	SendStatusReady(u *user.User)

	// 叫分 加注 亮牌
	CallScore(u *user.User, scoreTimes int)
	AddScore(u *user.User, score int)
	OpenCard(u *user.User, cardType int, cardData []int)

	// 明牌
	ShowCard(u *user.User)
	// 托管
	OtherOperation(args []interface{})
	GetUserScore(int) int
	InitRoomOne()
	//
	GetCreatorNodeId() int
	GetCreator() int64
}

type LogicManager interface {
	RandCardList(cbCardBuffer, OriDataArray []int)
	SortCardList(cardData []int, cardCount int)
	GetCardValue(CardData int) int
	GetCardColor(CardData int) int

	CompareCard(firstCardData []int, lastCardData []int) bool
	GetCardType(cardData []int) int
	GetCardTimes(cardType int) int

	CompareCardWithParam(firstCardData []int, lastCardData []int, args []interface{}) (int, bool)
	// 以下接口不通用
	RemoveCardList(cbRemoveCard []int, cbCardData []int) ([]int, bool)
	SetParamToLogic(args interface{}) // 设置算法必要参数

}
