package mj_base

import (
	"mj/common/msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/user"
)

type DataManager interface {
	BeforeStartGame(UserCnt int)                                     //开始前的处理
	StartGameing()                                                   //游戏开始种的处理
	InitRoomOne()                                                    //
	AfterStartGame()                                                 //开始游戏的善后处理
	SendPersonalTableTip(*user.User)                                 //发送没开始前的场景信息
	SendStatusPlay(u *user.User)                                     //发送开始后的处理
	NotifySendCard(u *user.User, cbCardData int, bSysOut bool)       //通知发牌
	EstimateUserRespond(int, int, int) bool                          //检测能否吃碰杠胡
	DispatchCardData(int, bool) int                                  //派发扑克
	HasOperator(ChairId, OperateCode int) bool                       //是否存在操作
	HasCard(ChairId, cardIdx int) bool                               //检测是否存在牌
	CheckUserOperator(*user.User, int, int, []int) (int, int)        //检测吃碰杠胡
	UserChiHu(wTargetUser, userCnt int)                              //吃胡处理
	WeaveCard(cbTargetAction, wTargetUser int)                       //组合扑克
	RemoveCardByOP(wTargetUser, ChoOp int) bool                      //根据操作码删除扑克
	CallOperateResult(wTargetUser, cbTargetAction int)               //计算吃碰杠胡的操作结果
	ResetUserOperate()                                               //重置用户状态
	ZiMo(u *user.User)                                               //自摸处理
	AnGang(u *user.User, cbOperateCode int, cbOperateCard []int) int //暗杠处理
	NormalEnd(Reason int)                                            //正常结束
	DismissEnd(Reason int)                                           //解散结束
	GetTrusteeOutCard(wChairID int) int                              //获取托管的牌
	CanOperatorRoom(uid int64) bool                                  //能否操作房间
	SendStatusReady(u *user.User)                                    //发送准备
	GetChaHua(u *user.User, setCount int)                            //获取插花
	OnUserReplaceCard(u *user.User, CardData int) bool               //替换牌
	OnUserListenCard(u *user.User, bListenCard bool) bool            //听牌
	RecordFollowCard(cbCenterCard int) bool                          //记录跟牌
	RecordOutCarCnt() int                                            //记录出牌数
	OnZhuaHua(winUser []int) (CardData [][]int, BuZhong []int)       //抓花 扎码出库
	RecordBanCard(OperateCode, ChairId int)                          //记录出牌禁忌
	OutOfChiCardRule(CardData, ChairId int) bool                     //吃啥打啥
	SendOperateResult(u *user.User, wrave *msg.WeaveItem)            //通知操作结果
	ResetUserOperateEx(u *user.User)                                 //清除状态

	GetResumeUser() int
	GetGangStatus() int
	GetUserCardIndex(ChairId int) []int
	GetCurrentUser() int //当前出牌用户
	GetRoomId() int
	GetCreater() int64
	GetProvideUser() int
	IsActionDone() bool
	GetUserScore(chairid int) int

	//测试专用函数。 请勿用于生产
	SetUserCard(charirID int, cards []int)
}

type LogicManager interface {
	AnalyseTingCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbOutCardData, cbHuCardCount []int, cbHuCardData [][]int) int
	GetCardCount(cbCardIndex []int) int
	RemoveCard(cbCardIndex []int, cbRemoveCard int) bool
	RemoveCardByArr(cbCardIndex, cbRemoveCard []int) bool
	EstimatePengCard(cbCardIndex []int, cbCurrentCard int) int
	EstimateGangCard(cbCardIndex []int, cbCurrentCard int) int
	EstimateEatCard(cbCardIndex []int, cbCurrentCard int) int
	GetUserActionRank(cbUserAction int) int
	AnalyseChiHuCard([]int, []*msg.WeaveItem, int) (bool, []*TagAnalyseItem)
	AnalyseGangCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbProvideCard int, gcr *TagGangCardResult) int
	GetHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbHuCardData []int, MaxCount int) int
	RandCardList(cbCardBuffer, OriDataArray []int)
	GetUserCards(cbCardIndex []int) (cbCardData []int)
	SwitchToCardData(cbCardIndex int) int
	SwitchToCardIndex(cbCardData int) int
	IsValidCard(card int) bool
	GetHuOfCard() int //记录胡的那张牌

	GetCardColor(cbCardData int) int
	GetCardValue(cbCardData int) int

	GetMagicIndex() int
	SetMagicIndex(int)
}
