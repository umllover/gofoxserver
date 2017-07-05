package mj_base

import (
	"mj/common/msg"
	"mj/common/utils"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"strconv"
)

type DataManager interface {
	BeforeStartGame(UserCnt int)                                     //开始前的处理
	StartGameing()                                                   //游戏开始种的处理
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
	ZiMo(u *user.User)                                               //自摸处理
	AnGang(u *user.User, cbOperateCode int, cbOperateCard []int) int //暗杠处理
	NormalEnd()                                                      //正常结束
	DismissEnd()                                                     //解散结束
	GetTrusteeOutCard(wChairID int) int                              //获取托管的牌
	CanOperatorRoom(uid int) bool                                    //能否操作房间
	SendStatusReady(u *user.User)                                    //发送准备
	GetChaHua(u *user.User, setCount int)                            //获取插花
	OnUserReplaceCard(u *user.User, CardData int) bool               //替换牌
	OnUserListenCard(u *user.User, bListenCard bool) bool            //听牌
	RecordFollowCard(cbCenterCard int) bool                          //记录跟牌
	RecordOutCarCnt() int                                            //记录出牌数
	OnZhuaHua(CenterUser int) (CardData []int, BuZhong []int)        //抓花 扎码出库

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

	GetCardColor(cbCardData int) int
	GetCardValue(cbCardData int) int

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
