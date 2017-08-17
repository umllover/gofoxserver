package mj_base

import (
	"math"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/mj_hz_msg"
	"mj/common/utils"
	. "mj/gameServer/common/mj"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"strconv"
	"time"

	"strings"

	"mj/gameServer/conf"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/timer"
	"github.com/lovelly/leaf/util"
)

func NewDataMgr(id int, uid int64, configIdx int, name string, temp *base.GameServiceOption, base *Mj_base, info *msg.L2G_CreatorRoom) *RoomData {
	r := new(RoomData)
	r.ID = id
	if name == "" {
		r.Name = temp.RoomName
	} else {
		r.Name = name
	}
	r.CreatorUid = uid
	r.CreatorNodeId = info.CreatorNodeId
	r.Source = temp.Source
	r.IniSource = temp.IniScore
	r.OtherInfo = info.OtherInfo //客户端动态的配置信息
	r.MjBase = base
	r.ConfigIdx = configIdx
	return r
}

//当一张桌子理解
type RoomData struct {
	ID            int
	Name          string //房间名字
	CreatorUid    int64  //创建房间的人
	CreatorNodeId int    //创建房间者的NodeId
	MjBase        *Mj_base
	ConfigIdx     int //配置索引

	IsResponse        []bool //标记是否对吃碰杠胡做出过动作
	PerformAction     []int  //记住玩家出的动作， 用来等待优先级更高的玩家
	Source            int    //底分
	IniSource         int    //初始分数
	IsGoldOrGameScore int    //金币场还是积分场 0 标识 金币场 1 标识 积分场
	Password          string // 密码
	ProvideCard       int    //供应扑克
	ResumeUser        int    //还原用户
	ProvideUser       int    //供应用户
	OutCardCount      int    //出牌数目
	EndLeftCount      int    //荒庄牌数
	LastCatchCardUser int    //最后一个摸牌的用户
	MinusHeadCount    int    //头部空缺
	MinusLastCount    int    //尾部空缺
	TingCnt           []int  //听牌个数

	SiceCount       int                //色子大小
	UserActionDone  bool               //操作完成
	SendStatus      int                //发牌状态
	GangStatus      int                //杠牌状态
	GangOutCard     bool               //杠后出牌
	ProvideGangUser int                //供杠用户
	GangCard        []bool             //杠牌状态
	GangCount       []int              //杠牌次数
	RepertoryCard   []int              //库存扑克
	UserGangScore   []int              //游戏中杠的输赢
	ChiHuKind       []int              //吃胡结果
	ChiHuRight      []int              //胡牌类型
	UserAction      []int              //用户动作
	OperateCard     [][]int            //操作扑克
	ChiPengCount    []int              //吃碰杠次数
	CardIndex       [][]int            //用户扑克[GAME_PLAYER][MAX_INDEX]
	WeaveItemArray  [][]*msg.WeaveItem //组合扑克
	DiscardCard     [][]int            //丢弃记录
	OutCardData     int                //出牌扑克
	OutCardUser     int                //当前出牌用户
	HeapHead        int                //堆立头部
	HeapTail        int                //堆立尾部
	HeapCardInfo    [][]int            //堆牌信息
	SendCardData    int                //发牌扑克
	HistorySe       *HistoryScore      //历史积分
	CurrentUser     int                //当前操作用户
	Ting            []bool             //是否听牌
	BankerUser      int                //庄家用户
	HuOfCard        int                //胡的牌

	FlowerCnt    []int                  //补花数
	FlowerCard   [][]int                //记录花牌
	ChangeBanker bool                   //庄家是否变动
	OtherInfo    map[string]interface{} //客户端动态的配置信息

	BanUser    []int   //是否出牌禁忌
	BanCardCnt [][]int //禁忌卡牌

	//timer
	OperateTime []*timer.Timer //操作定时器
}

func (room *RoomData) GetDataMgr() DataManager {
	return room.MjBase.GetDataMgr()
}

func (room *RoomData) InitRoomOne() {
	room.HistorySe = &HistoryScore{AllScore: make([]int, room.MjBase.UserMgr.GetMaxPlayerCnt())}
}

func (room *RoomData) GetUserScore(chairid int) int {
	if chairid > room.MjBase.UserMgr.GetMaxPlayerCnt() {
		return 0
	}
	return room.HistorySe.AllScore[chairid]
}

func (room *RoomData) ResetUserScore() {
	room.HistorySe.AllScore = make([]int, room.MjBase.UserMgr.GetMaxPlayerCnt())
	room.HistorySe.DetailScore = [][]int{}
}

func (room *RoomData) GetCreator() int64 {
	return room.CreatorUid
}

func (room *RoomData) GetCreatorNodeId() int {
	return room.CreatorNodeId
}

func (room *RoomData) GetCfg() *MJ_CFG {
	return GetCfg(room.ConfigIdx)
}

func (room *RoomData) CanOperatorRoom(uid int64) bool {
	if uid == room.CreatorUid {
		return true
	}
	return false
}

func (room *RoomData) GetCurrentUser() int {
	return room.CurrentUser
}

func (room *RoomData) GetGangStatus() int {
	return room.GangStatus
}

func (room *RoomData) GetProvideUser() int {
	return room.ProvideUser
}

func (room *RoomData) GetResumeUser() int {
	return room.ResumeUser
}

func (room *RoomData) GetRoomId() int {
	return room.ID
}
func (room *RoomData) SendPersonalTableTip(u *user.User) {
	TableTip := &msg.G2C_PersonalTableTip{
		TableOwnerUserID:  room.CreatorUid,                                               //桌主 I D
		PlayerCnt:         room.MjBase.UserMgr.GetMaxPlayerCnt(),                         //玩家数量
		DrawCountLimit:    room.MjBase.TimerMgr.GetMaxPlayCnt(),                          //局数限制
		DrawTimeLimit:     room.MjBase.TimerMgr.GetTimeLimit(),                           //时间限制
		PlayCount:         room.MjBase.TimerMgr.GetPlayCount(),                           //已玩局数
		PlayTime:          int(room.MjBase.TimerMgr.GetCreatrTime() - time.Now().Unix()), //已玩时间
		CellScore:         room.Source,                                                   //游戏底分
		IniScore:          room.IniSource,                                                //初始分数
		ServerID:          strconv.Itoa(room.ID),                                         //房间编号
		PayType:           room.MjBase.UserMgr.GetPayType(),                              //支付类型
		IsJoinGame:        0,                                                             //是否参与游戏 todo  tagPersonalTableParameter
		IsGoldOrGameScore: room.IsGoldOrGameScore,                                        //金币场还是积分场 0 标识 金币场 1 标识 积分场
		OtherInfo:         room.OtherInfo,
		LeaveInfo:         room.MjBase.UserMgr.GetLeaveInfo(u.Id), //请求离家的玩家的信息
		KindID:            room.MjBase.Temp.KindID,
	}
	u.WriteMsg(TableTip)
}

func (room *RoomData) SendStatusReady(u *user.User) {
	StatusFree := &msg.G2C_StatusFree{}
	StatusFree.CellScore = room.Source                                     //基础积分
	StatusFree.TimeOutCard = room.MjBase.TimerMgr.GetTimeOutCard()         //出牌时间
	StatusFree.TimeOperateCard = room.MjBase.TimerMgr.GetTimeOperateCard() //操作时间
	StatusFree.CreateTime = room.MjBase.TimerMgr.GetCreatrTime()           //开始时间
	StatusFree.TurnScore = room.HistorySe.AllScore
	StatusFree.CollectScore = room.HistorySe.DetailScore
	StatusFree.PlayerCount = room.MjBase.UserMgr.GetBeginPlayer() //玩家人数
	StatusFree.MaCount = 0                                        //码数
	StatusFree.PlayCount = room.MjBase.TimerMgr.GetPlayCount()    //已玩局数
	StatusFree.CountLimit = room.MjBase.TimerMgr.GetMaxPlayCnt()  //局数限制
	u.WriteMsg(StatusFree)
}

func (room *RoomData) SendStatusPlay(u *user.User) {
	StatusPlay := &msg.G2C_StatusPlay{}
	//自定规则
	StatusPlay.TimeOutCard = room.MjBase.TimerMgr.GetTimeOutCard()
	StatusPlay.TimeOperateCard = room.MjBase.TimerMgr.GetTimeOperateCard()
	StatusPlay.CreateTime = room.MjBase.TimerMgr.GetCreatrTime()

	//规则
	StatusPlay.PlayerCount = room.MjBase.UserMgr.GetBeginPlayer()
	StatusPlay.PlayCount = room.MjBase.TimerMgr.GetPlayCount()
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	//游戏变量
	StatusPlay.BankerUser = room.BankerUser
	StatusPlay.CurrentUser = room.CurrentUser
	StatusPlay.CellScore = room.Source
	StatusPlay.MagicIndex = room.MjBase.LogicMgr.GetMagicIndex()
	StatusPlay.Trustee = room.MjBase.UserMgr.GetTrustees()
	StatusPlay.HuCardCount = make([]int, room.GetCfg().MaxCount)
	StatusPlay.HuCardData = make([][]int, room.GetCfg().MaxCount)
	StatusPlay.OutCardDataEx = make([]int, room.GetCfg().MaxCount)
	StatusPlay.CardCount = make([]int, UserCnt)

	//状态变量
	StatusPlay.ActionCard = room.ProvideCard
	StatusPlay.LeftCardCount = room.GetLeftCard()
	StatusPlay.ActionMask = room.UserAction[u.ChairId]

	StatusPlay.Ting = room.Ting
	//当前能胡的牌
	StatusPlay.OutCardCount = 0 //room.MjBase.LogicMgr.AnalyseTingCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], StatusPlay.OutCardDataEx, StatusPlay.HuCardCount, StatusPlay.HuCardData)

	//历史记录
	StatusPlay.OutCardUser = room.OutCardUser
	StatusPlay.OutCardData = room.OutCardData
	StatusPlay.DiscardCard = room.DiscardCard
	for _, v := range room.DiscardCard {
		StatusPlay.DiscardCount = append(StatusPlay.DiscardCount, len(v))
	}

	StatusPlay.WeaveItemArray = room.WeaveItemArray
	for _, v := range room.WeaveItemArray {
		StatusPlay.WeaveItemCount = append(StatusPlay.WeaveItemCount, len(v))
	}

	//堆立信息
	StatusPlay.HeapHead = room.HeapHead
	StatusPlay.HeapTail = room.HeapTail
	StatusPlay.HeapCardInfo = room.HeapCardInfo

	//扑克数据
	for j := 0; j < UserCnt; j++ {
		StatusPlay.CardCount[j] = room.MjBase.LogicMgr.GetCardCount(room.CardIndex[j])
	}

	StatusPlay.CardData = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[u.ChairId])
	if room.CurrentUser == u.ChairId {
		StatusPlay.SendCardData = room.SendCardData
	} else {
		StatusPlay.SendCardData = 0x00
	}

	//历史积分
	StatusPlay.TurnScore = room.HistorySe.AllScore
	StatusPlay.CollectScore = room.HistorySe.DetailScore
	u.WriteMsg(StatusPlay)
}

func (room *RoomData) NotifySendCard(u *user.User, cbCardData int, bSysOut bool) {
	//设置变量
	room.SendStatus = OutCard_Send
	room.SendCardData = 0
	room.UserAction[u.ChairId] = WIK_NULL

	//出牌记录
	room.OutCardUser = u.ChairId
	room.OutCardData = cbCardData

	//构造数据
	OutCard := &mj_hz_msg.G2C_HZMJ_OutCard{}
	OutCard.OutCardUser = u.ChairId
	OutCard.OutCardData = cbCardData
	OutCard.SysOut = bSysOut

	//发送消息
	room.MjBase.UserMgr.SendMsgAll(OutCard)
	room.ProvideUser = u.ChairId
	room.ProvideCard = cbCardData

	//丢弃扑克记录
	room.DiscardCard[u.ChairId] = append(room.DiscardCard[u.ChairId], cbCardData)

	//用户切换
	room.CurrentUser = (u.ChairId + 1) % room.MjBase.UserMgr.GetMaxPlayerCnt()
}

func (room *RoomData) GetUserCardIndex(ChairId int) []int {
	return room.CardIndex[ChairId]
}

//检测是否可以做某个操作
func (room *RoomData) HasOperator(ChairId, OperateCode int) bool {

	if room.UserAction[ChairId] == WIK_NULL {
		log.Error("+++++++++++++ room.UserAction:%v", room.UserAction)
		log.Error("room.UserAction[ChairId] == WIK_NULL, ChairId:%d", ChairId)
		return false
	}

	if OperateCode != WIK_NULL && ((room.UserAction[ChairId] & OperateCode) == 0) {
		//log.Error("HasOperator return false, ChairId=%d, OperateCode=%d, UserAction=%v", ChairId, OperateCode, room.UserAction[ChairId])
		return false
	}

	return true
}

//手中是否存在某张牌
func (room *RoomData) HasCard(ChairId, cardIdx int) bool {
	if cardIdx > room.GetCfg().MaxIdx {
		return false
	}
	return room.CardIndex[ChairId][cardIdx] > 0
}

//设置用户相应牌的操作 ,返回是否可以操作
func (room *RoomData) CheckUserOperator(u *user.User, userCnt, OperateCode int, OperateCard []int) (int, int) {
	if room.IsResponse[u.ChairId] {
		return -1, u.ChairId
	}

	room.IsResponse[u.ChairId] = true
	room.PerformAction[u.ChairId] = OperateCode
	room.OperateCard[u.ChairId] = make([]int, 4)
	if len(OperateCard) < 2 {
		room.BuildOpCard(u.ChairId, OperateCode, OperateCard[0])
	} else {
		for i, card := range OperateCard {
			room.OperateCard[u.ChairId][i] = card
		}
	}

	u.UserLimit = 0
	//放弃操作
	if OperateCode == WIK_NULL {
		////禁止这轮吃胡
		if room.HasOperator(u.ChairId, WIK_CHI_HU) {
			u.UserLimit |= LimitChiHu
		}
	}

	cbTargetAction := OperateCode
	wTargetUser := u.ChairId
	//执行判断
	for i := 0; i < userCnt; i++ {
		//获取动作
		cbUserAction := room.UserAction[i]
		if room.IsResponse[i] {
			cbUserAction = room.PerformAction[i]
		}

		//优先级别
		cbUserActionRank := room.MjBase.LogicMgr.GetUserActionRank(cbUserAction)
		cbTargetActionRank := room.MjBase.LogicMgr.GetUserActionRank(cbTargetAction)

		//动作判断
		if cbUserActionRank > cbTargetActionRank {
			wTargetUser = i
			cbTargetAction = cbUserAction
		}
	}

	if room.IsResponse[wTargetUser] == false { //最高权限的人没响应
		log.Error("CheckUserOperator wTargetUser=%d, ChairId=%d, room.IsResponse=%v", wTargetUser, u.ChairId, room.IsResponse)
		return -1, u.ChairId
	}

	if cbTargetAction == WIK_NULL {
		room.IsResponse = make([]bool, userCnt)
		room.UserAction = make([]int, userCnt)
		room.OperateCard = make([][]int, userCnt)
		room.PerformAction = make([]int, userCnt)
		return cbTargetAction, wTargetUser
	}
	//走到这里一定是所有人都响应完了
	return cbTargetAction, wTargetUser
}

func (room *RoomData) BuildOpCard(ChairId, OperateCode, opcard int) {
	if OperateCode&WIK_LEFT != 0 {
		room.OperateCard[ChairId][0] = opcard
		room.OperateCard[ChairId][1] = opcard + 1
		room.OperateCard[ChairId][2] = opcard + 2
	} else if OperateCode&WIK_CENTER != 0 {
		room.OperateCard[ChairId][0] = opcard - 1
		room.OperateCard[ChairId][1] = opcard
		room.OperateCard[ChairId][2] = opcard + 1
	} else if OperateCode&WIK_RIGHT != 0 {
		room.OperateCard[ChairId][0] = opcard - 2
		room.OperateCard[ChairId][1] = opcard - 1
		room.OperateCard[ChairId][2] = opcard
	} else {
		room.OperateCard[ChairId][0] = opcard
		room.OperateCard[ChairId][1] = opcard
		room.OperateCard[ChairId][2] = opcard
		if OperateCode&WIK_GANG != 0 {
			room.OperateCard[ChairId][3] = opcard
		}
	}
}

func (room *RoomData) GetOpCard(ChairId, OperateCode int) int {
	if OperateCode&WIK_LEFT != 0 {
		return room.OperateCard[ChairId][0]
	} else if OperateCode&WIK_CENTER != 0 {
		return room.OperateCard[ChairId][1]
	} else if OperateCode&WIK_RIGHT != 0 {
		return room.OperateCard[ChairId][2]
	} else {
		return room.OperateCard[ChairId][0]
	}
}

func (room *RoomData) UserChiHu(wTargetUser, userCnt int) {
	//结束信息
	wChiHuUser := room.BankerUser
	for i := 0; i < userCnt; i++ {
		wChiHuUser = (room.BankerUser + i) % userCnt
		//过虑判断
		if (room.PerformAction[wChiHuUser] & WIK_CHI_HU) == 0 { //一跑多响
			continue
		}

		//胡牌判断
		pWeaveItem := room.WeaveItemArray[wChiHuUser]
		hu, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[wChiHuUser], pWeaveItem, room.OperateCard[wTargetUser][0])
		if hu {
			room.ChiHuKind[wChiHuUser] = WIK_CHI_HU
		}

		//todo room.ChiHuRight[wChiHuUser]
		//插入扑克
		if room.ChiHuKind[wChiHuUser] != WIK_NULL {
			wTargetUser = wChiHuUser
		}
	}
}

//组合要操作的牌
func (room *RoomData) WeaveCard(cbTargetAction, wTargetUser int) {
	//变量定义
	cbTargetCard := room.GetOpCard(wTargetUser, cbTargetAction)

	//出牌变量
	room.SendStatus = Gang_Send
	room.SendCardData = 0
	room.OutCardUser = INVALID_CHAIR
	room.OutCardData = 0

	//组合扑克
	Wrave := &msg.WeaveItem{}
	Wrave.Param = WIK_GANERAL
	Wrave.CenterCard = cbTargetCard
	Wrave.WeaveKind = cbTargetAction
	Wrave.CardData = make([]int, 4)
	if room.ProvideUser == INVALID_CHAIR {
		Wrave.ProvideUser = wTargetUser
	} else {
		Wrave.ProvideUser = room.ProvideUser
	}

	Wrave.CardData = util.CopySlicInt(room.OperateCard[wTargetUser])

	log.Debug("###############杠牌：%v", Wrave)
	room.WeaveItemArray[wTargetUser] = append(room.WeaveItemArray[wTargetUser], Wrave)
}

func (room *RoomData) RemoveCardByOP(wTargetUser, ChoOp int) bool {
	opCard := room.OperateCard[wTargetUser]
	var card []int
	switch ChoOp {
	case WIK_LEFT:
		card = opCard[1:3]
	case WIK_RIGHT:
		card = opCard[0:2]
	case WIK_CENTER:
		card = []int{opCard[0], opCard[2]}
	case WIK_PENG:
		card = []int{opCard[0], opCard[0]}
	case WIK_GANG: //杠牌操作
		card = []int{opCard[0], opCard[0], opCard[0]}
	default:
		return false
	}
	//删除扑克
	if !room.MjBase.LogicMgr.RemoveCardByArr(room.CardIndex[wTargetUser], card) {
		log.Error("not foud card at RemoveCardByOP :%v :%v", card, room.CardIndex[wTargetUser])
		return false
	}
	room.ChiPengCount[wTargetUser]++
	return true
}

func (room *RoomData) RemoveDiscardInfo() {
	//吃碰杠胡后移除供牌用户的最后一个丢牌记录
	room.DiscardCard[room.ProvideUser] = utils.IntSliceDelete(room.DiscardCard[room.ProvideUser], len(room.DiscardCard[room.ProvideUser])-1)
}

func (room *RoomData) AnGang(u *user.User, cbOperateCode int, cbOperateCard []int) int {
	room.SendStatus = Gang_Send
	//变量定义
	var cbWeave *msg.WeaveItem
	cbCardIndex := room.MjBase.LogicMgr.SwitchToCardIndex(cbOperateCard[0])
	wProvideUser := u.ChairId
	cbGangKind := WIK_MING_GANG
	//杠牌处理
	if room.CardIndex[u.ChairId][cbCardIndex] == 1 {
		//寻找组合
		for _, v := range room.WeaveItemArray[u.ChairId] {
			if (v.CenterCard == cbOperateCard[0]) && (v.WeaveKind == WIK_PENG) {
				cbWeave = v
				break
			}
		}

		//没找到明杠
		if cbWeave == nil {
			return 0
		}
		cbGangKind = WIK_MING_GANG

		//组合扑克
		cbWeave.Param = WIK_MING_GANG
		cbWeave.WeaveKind = cbOperateCode
		cbWeave.CenterCard = cbOperateCard[0]
		cbWeave.CardData[3] = cbOperateCard[0]

		//杠牌得分
		wProvideUser = cbWeave.ProvideUser
	} else {
		//扑克效验

		if room.CardIndex[u.ChairId][cbCardIndex] != 4 {
			return 0
		}
		cbGangKind = WIK_AN_GANG
		cbWeave = &msg.WeaveItem{}
		cbWeave.Param = WIK_AN_GANG
		cbWeave.ProvideUser = u.ChairId
		cbWeave.WeaveKind = cbOperateCode
		cbWeave.CenterCard = cbOperateCard[0]
		cbWeave.CardData = make([]int, 4)
		for j := 0; j < 4; j++ {
			cbWeave.CardData[j] = cbOperateCard[0]
		}
		room.WeaveItemArray[u.ChairId] = append(room.WeaveItemArray[u.ChairId], cbWeave)
	}

	//删除扑克
	room.CardIndex[u.ChairId][cbCardIndex] = 0
	room.GangStatus = cbGangKind
	room.ProvideGangUser = wProvideUser
	room.GangCard[u.ChairId] = true
	room.GangCount[u.ChairId]++

	//发送消息
	room.GetDataMgr().SendOperateResult(u, cbWeave)
	return cbGangKind
}

//发送操作结果
func (room *RoomData) SendOperateResult(u *user.User, wrave *msg.WeaveItem) {
	OperateResult := &mj_hz_msg.G2C_HZMJ_OperateResult{}
	OperateResult.ProvideUser = wrave.ProvideUser
	OperateResult.OperateCode = wrave.WeaveKind
	for i := 0; i < len(OperateResult.OperateCard); i++ {
		OperateResult.OperateCard[i] = wrave.CardData[i]
	}
	if u != nil {
		OperateResult.OperateUser = u.ChairId
	} else {
		OperateResult.OperateUser = wrave.OperateUser
		OperateResult.ActionMask = wrave.ActionMask
	}
	log.Debug("############# SendOperateResult OperateUser=%d, ProvideUser=%d", OperateResult.OperateUser, OperateResult.ProvideUser)
	log.Debug("############# SendOperateResult ActionMask=%d, OperateCode=%d, OperateCard=%v", OperateResult.ActionMask, OperateResult.OperateCode, OperateResult.OperateCard)
	room.MjBase.UserMgr.SendMsgAll(OperateResult)
}

func (room *RoomData) ZiMo(u *user.User) {
	//普通胡牌
	pWeaveItem := room.WeaveItemArray[u.ChairId]
	if !room.MjBase.LogicMgr.RemoveCard(room.CardIndex[u.ChairId], room.SendCardData) {
		log.Error("not foud card at Operater")
		return
	}
	hu, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], pWeaveItem, room.SendCardData)
	if hu {
		room.ChiHuKind[u.ChairId] = WIK_CHI_HU
	}

	//todo  room.ChiHuRight[u.ChairId]
	//结束信息

	room.ProvideCard = room.SendCardData
	return
}

func (room *RoomData) CallOperateResult(wTargetUser, cbTargetAction int) {
	//构造结果
	wrave := &msg.WeaveItem{}
	wrave.OperateUser = wTargetUser
	wrave.WeaveKind = cbTargetAction
	if room.ProvideUser == INVALID_CHAIR {
		wrave.ProvideUser = wTargetUser
	} else {
		wrave.ProvideUser = room.ProvideUser
	}

	wrave.CardData = util.CopySlicInt(room.OperateCard[wTargetUser])
	wrave.CenterCard = room.GetOpCard(wTargetUser, cbTargetAction)

	//用户状态
	room.ResetUserOperate()

	//如果非杠牌
	if cbTargetAction != WIK_GANG {
		room.ProvideUser = INVALID_CHAIR
		room.ProvideCard = 0

		gcr := &TagGangCardResult{}
		room.UserAction[wTargetUser] |= room.MjBase.LogicMgr.AnalyseGangCard(room.CardIndex[wTargetUser], room.WeaveItemArray[wTargetUser], 0, gcr)

		if room.CheckTingCard(wTargetUser) {
			wrave.ActionMask |= room.UserAction[wTargetUser]
		}
	}

	//发送消息
	room.GetDataMgr().SendOperateResult(nil, wrave)

	//设置用户
	room.CurrentUser = wTargetUser

	//杠牌处理
	if cbTargetAction == WIK_GANG {
		room.GangStatus = WIK_FANG_GANG
		if room.ProvideUser == INVALID_CHAIR {
			room.ProvideGangUser = wTargetUser
		} else {
			room.ProvideGangUser = room.ProvideUser
		}
		room.GangCard[wTargetUser] = true
		room.GangCount[wTargetUser]++

	}
	return
}

//重置用户状态
func (room *RoomData) ResetUserOperate() {
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	room.IsResponse = make([]bool, UserCnt)
	room.UserAction = make([]int, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
}

//重置用户状态
func (room *RoomData) ResetUserOperateEx(u *user.User) {
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	room.UserAction = make([]int, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
}

//响应判断
func (room *RoomData) EstimateUserRespond(wCenterUser int, cbCenterCard int, EstimatKind int) bool {
	log.Debug("EstimateUserRespond start")
	//变量定义
	bAroseAction := false

	//用户状态
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	room.UserAction = make([]int, UserCnt)

	//动作判断
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		//用户过滤
		if wCenterUser == u.ChairId {
			//log.Debug("at EstimateUserRespond ======== wCenterUser:%v", wCenterUser)
			return
		}

		//托管了不响应
		if room.MjBase.UserMgr.IsTrustee(u.ChairId) {
			log.Debug("at EstimateUserRespond ======== IsTrustee ChairId:%v", u.ChairId)
			return
		}

		//出牌类型
		if EstimatKind == EstimatKind_OutCard {
			//吃碰判断
			log.Debug(".UserLimit&LimitPen %v, %v", u.UserLimit, u.UserLimit&LimitPeng)
			if u.UserLimit&LimitPeng == 0 {
				//碰牌判断
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimatePengCard(room.CardIndex[u.ChairId], cbCenterCard)
			}

			//杠牌判断
			log.Debug(".room.LeftCardCount > room.EndLeftCount %v, %v", room.IsEnoughCard(), u.UserLimit&LimitGang)
			if room.IsEnoughCard() && u.UserLimit&LimitGang == 0 {
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimateGangCard(room.CardIndex[u.ChairId], cbCenterCard)
			}
		}

		//检查抢杠胡
		if EstimatKind == EstimatKind_GangCard {
			//只有庄家和闲家之间才能放炮
			MagicCard := room.MjBase.LogicMgr.SwitchToCardData(room.MjBase.LogicMgr.GetMagicIndex())
			if room.MjBase.LogicMgr.GetMagicIndex() == room.GetCfg().MaxIdx || (room.MjBase.LogicMgr.GetMagicIndex() != room.GetCfg().MaxIdx && cbCenterCard != MagicCard) {
				if u.UserLimit&LimitChiHu == 0 {
					//吃胡判断
					hu, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], cbCenterCard)
					if hu {
						room.UserAction[u.ChairId] |= WIK_CHI_HU
					}
				}
			}
		}

		//结果判断
		if room.UserAction[u.ChairId] != WIK_NULL {
			bAroseAction = true
		}
	})

	//log.Debug("AAAAAAAAAAAAAAAA : %v", bAroseAction)
	//结果处理
	if bAroseAction {
		//设置变量
		room.ProvideUser = wCenterUser
		room.ProvideCard = cbCenterCard
		room.ResumeUser = room.CurrentUser
		room.CurrentUser = INVALID_CHAIR

		//发送提示
		room.GetDataMgr().SendOperateNotify(nil, 0)
		return true
	}

	//if room.GangStatus != WIK_GANERAL {
	//	room.GangOutCard = true
	//	room.GangStatus = WIK_GANERAL
	//	room.ProvideGangUser = INVALID_CHAIR
	//} else {
	//	room.GangOutCard = false
	//}

	return false
}

//发送提示
func (room *RoomData) SendOperateNotify(u *user.User, card int) {
	if u != nil {
		u.WriteMsg(&mj_hz_msg.G2C_HZMJ_OperateNotify{
			ActionMask: room.UserAction[u.ChairId],
			ActionCard: card,
		})
	} else {
		room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
			if room.UserAction[u.ChairId] != WIK_NULL {
				log.Debug("########### EstimateUserRespond ActionMask %v ###########", room.UserAction[u.ChairId])
				u.WriteMsg(&mj_hz_msg.G2C_HZMJ_OperateNotify{
					ActionMask: room.UserAction[u.ChairId],
					ActionCard: room.ProvideCard,
				})
			}
		})
	}
}

//派发扑克
func (room *RoomData) DispatchCardData(wCurrentUser int, bTail bool) int {
	log.Error("at  base DispatchCardData ...................... ")
	//状态效验
	if room.SendStatus == Not_Send {
		log.Error("at DispatchCardData f room.SendStatus == Not_Send")
		return -1
	}

	//荒庄结束
	if !room.IsEnoughCard() {
		room.ProvideUser = INVALID_CHAIR
		return 1
	}

	//发送扑克
	room.ProvideCard = room.GetSendCard(bTail, room.MjBase.UserMgr.GetMaxPlayerCnt())
	room.SendCardData = room.ProvideCard
	room.LastCatchCardUser = wCurrentUser

	u := room.MjBase.UserMgr.GetUserByChairId(wCurrentUser)
	if u == nil {
		log.Error("at DispatchCardData not foud user ")
	}

	//清除禁止胡牌的牌
	u.UserLimit &= ^LimitChiHu
	u.UserLimit &= ^LimitPeng
	u.UserLimit &= ^LimitGang

	//设置变量
	room.OutCardUser = INVALID_CHAIR
	room.OutCardData = 0
	room.CurrentUser = wCurrentUser
	room.ProvideUser = wCurrentUser
	room.GangOutCard = false

	if bTail { //从尾部取牌，说明玩家杠牌了,计算分数
		room.CallGangScore()
		room.GangStatus = WIK_GANERAL
		room.ProvideGangUser = INVALID_CHAIR
	}

	//加牌
	room.CardIndex[wCurrentUser][room.MjBase.LogicMgr.SwitchToCardIndex(room.ProvideCard)]++
	//room.UserCatchCardCount[wCurrentUser]++;

	if !room.MjBase.UserMgr.IsTrustee(wCurrentUser) {
		//胡牌判断
		room.CardIndex[wCurrentUser][room.MjBase.LogicMgr.SwitchToCardIndex(room.SendCardData)]--
		//log.Debug("before %v ", room.UserAction[wCurrentUser])
		hu, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser], room.SendCardData)
		if hu {
			room.UserAction[wCurrentUser] |= WIK_CHI_HU
		}

		//log.Debug("after %v ", room.UserAction[wCurrentUser])
		room.CardIndex[wCurrentUser][room.MjBase.LogicMgr.SwitchToCardIndex(room.SendCardData)]++

		//杠牌判断
		if room.IsEnoughCard() && !room.Ting[wCurrentUser] {
			GangCardResult := &TagGangCardResult{}
			room.UserAction[wCurrentUser] |= room.MjBase.LogicMgr.AnalyseGangCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser], room.ProvideCard, GangCardResult)
		}
	}

	//听牌判断
	room.CheckTingCard(wCurrentUser)

	if room.SendStatus != BuHua_Send {
		room.GetDataMgr().SendCardToCli(u, bTail)
	}
	log.Debug("User Action === %v , %d", room.UserAction, room.UserAction[wCurrentUser])

	room.UserActionDone = false
	if room.MjBase.UserMgr.IsTrustee(wCurrentUser) {
		room.UserActionDone = true
	}

	if room.UserAction[wCurrentUser] != WIK_NULL {
		room.GetDataMgr().SendOperateNotify(u, room.ProvideCard)
	}
	return 0
}

func (room *RoomData) SendCardToCli(u *user.User, bTail bool) {
	//构造数据
	SendCard := &mj_hz_msg.G2C_HZMJ_SendCard{}
	SendCard.SendCardUser = room.CurrentUser
	SendCard.CurrentUser = room.CurrentUser
	SendCard.Tail = bTail
	SendCard.ActionMask = room.UserAction[room.CurrentUser]
	SendCard.CardData = room.ProvideCard
	//发送数据
	u.WriteMsg(SendCard)

	SendCardOther := &mj_hz_msg.G2C_HZMJ_SendCard{}
	SendCardOther.SendCardUser = room.CurrentUser
	SendCardOther.CurrentUser = room.CurrentUser
	SendCardOther.Tail = bTail
	SendCardOther.ActionMask = room.UserAction[room.CurrentUser]
	SendCardOther.CardData = 0
	room.MjBase.UserMgr.SendMsgAllNoSelf(u.Id, SendCardOther)
}

func (room *RoomData) BeforeStartGame(UserCnt int) {
	room.InitRoom(UserCnt)
}

func (room *RoomData) StartGameing() {
	gameLogic := room.MjBase.LogicMgr
	//万能牌设置
	if room.GetCfg().MagicCard != 0 {
		gameLogic.SetMagicIndex(gameLogic.SwitchToCardIndex(room.GetCfg().MagicCard))
	}

	//洗牌
	gameLogic.RandCardList(room.RepertoryCard, GetCardByIdx(room.ConfigIdx))

	log.Debug("======房间Id：%d", room.ID)
	//选取庄家
	room.ElectionBankerUser()

	//发牌
	room.StartDispatchCard()

	//开局补花
	room.InitBuHua()
}

func (room *RoomData) AfterStartGame() {
	log.Debug("============= AfterStartGame")
	//检查自摸
	room.CheckTingCard(room.BankerUser)
	//通知客户端开始了
	room.SendGameStart()
	//检测起手动作
	room.InitBankerAction()
}

func (room *RoomData) ResetGameAfterRenewal() {
	room.ResetUserScore() //重置用户所有积分
}

func (room *RoomData) InitRoom(UserCnt int) {
	//初始化
	log.Debug("初始化房间")
	room.RepertoryCard = make([]int, room.GetCfg().MaxRepertory)
	room.CardIndex = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.CardIndex[i] = make([]int, room.GetCfg().MaxIdx)
	}
	room.FlowerCnt = make([]int, UserCnt)
	room.ChiHuKind = make([]int, UserCnt)
	room.ChiPengCount = make([]int, UserCnt)
	room.GangCard = make([]bool, UserCnt) //杠牌状态
	room.GangCount = make([]int, UserCnt)
	room.Ting = make([]bool, UserCnt)
	room.UserAction = make([]int, UserCnt)
	room.DiscardCard = make([][]int, UserCnt)
	room.UserGangScore = make([]int, UserCnt)
	room.WeaveItemArray = make([][]*msg.WeaveItem, UserCnt)
	room.ChiHuRight = make([]int, UserCnt)
	room.HeapCardInfo = make([][]int, UserCnt)
	room.IsResponse = make([]bool, UserCnt)
	room.PerformAction = make([]int, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
	room.BanUser = make([]int, UserCnt)
	room.BanCardCnt = make([][]int, UserCnt)
	room.TingCnt = make([]int, UserCnt)
	room.FlowerCard = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.HeapCardInfo[i] = make([]int, 2)
		room.BanCardCnt[i] = make([]int, 9)
	}
	room.UserActionDone = false
	room.SendStatus = Not_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideGangUser = INVALID_CHAIR
	room.MinusLastCount = 0
	room.MinusHeadCount = room.GetCfg().MaxRepertory
	room.OutCardCount = 0
	room.HuOfCard = 0
}

func (room *RoomData) GetSice() (int, int) {
	Sice1 := util.RandInterval(1, 6)
	Sice2 := util.RandInterval(1, 6)
	minSice := int(math.Min(float64(Sice1), float64(Sice2)))
	return Sice2<<8 | Sice1, minSice
}

//开局补花
func (room *RoomData) InitBuHua() {
	log.Debug("开局补花")
	if room.GetCfg().HuaIndex == 0 {
		log.Debug("not hua card at InitBuHua")
	}
	playerIndex := room.BankerUser
	playerCNT := room.MjBase.UserMgr.GetMaxPlayerCnt()
	for i := 0; i < playerCNT; i++ {
		if playerIndex > 3 {
			playerIndex = 0
		}
		room.GetDataMgr().CheckHuaCard(playerIndex, playerCNT, true)
		playerIndex++
	}
}

func (room *RoomData) CheckHuaCard(playerIndex, playerCNT int, IsInitFlower bool) {
	logic := room.MjBase.LogicMgr
	for j := room.GetCfg().MaxIdx - room.GetCfg().HuaIndex; j < room.GetCfg().MaxIdx; j++ {
		if room.CardIndex[playerIndex][j] == 1 {
			index := j
			for {
				NewCard := room.GetSendCard(true, playerCNT)
				newCardIndex := logic.SwitchToCardIndex(NewCard)
				ReplaceCard := logic.SwitchToCardData(index)
				room.GetDataMgr().SendReplaceCard(playerIndex, ReplaceCard, NewCard, IsInitFlower)
				log.Debug("玩家%d,j:%d 补花：%x，新牌：%x", playerIndex, j, logic.SwitchToCardData(index), NewCard)
				room.FlowerCnt[playerIndex]++
				room.FlowerCard[playerIndex] = append(room.FlowerCard[playerIndex], ReplaceCard)
				if newCardIndex < (room.GetCfg().MaxIdx - room.GetCfg().HuaIndex) {
					room.CardIndex[playerIndex][j]--
					room.CardIndex[playerIndex][newCardIndex]++
					if playerIndex == room.BankerUser || !IsInitFlower {
						room.SendCardData = NewCard
						room.ProvideCard = NewCard
					}
					break
				} else {
					index = newCardIndex
				}
			}
		}
	}
	return
}

//用户补花
func (room *RoomData) OnUserReplaceCard(u *user.User, CardData int) bool {
	if room.GetCfg().HuaIndex <= 0 {
		log.Debug("atOnUserReplaceCard  HuaIndex === 0")
		return true
	}

	log.Debug("[用户补花开始] 用户：%d补花：%d", u.ChairId, CardData)
	gameLogic := room.MjBase.LogicMgr
	if gameLogic.RemoveCard(room.CardIndex[u.ChairId], CardData) == false {
		log.Error("[用户补花] 用户：%d补花失败", u.ChairId)
		return false
	}

	//记录补花
	room.FlowerCnt[u.ChairId]++
	room.FlowerCard[u.ChairId] = append(room.FlowerCard[u.ChairId], CardData)

	//是否花杠
	if room.FlowerCnt[u.ChairId] == 8 {
		room.UserAction[u.ChairId] |= WIK_CHI_HU
	}

	//状态变量
	room.SendStatus = BuHua_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideGangUser = INVALID_CHAIR

	//派发扑克
	room.GetDataMgr().DispatchCardData(u.ChairId, true)
	room.GetDataMgr().SendReplaceCard(u.ChairId, CardData, room.SendCardData, false)

	log.Debug("[用户补花结束] 用户：%d,花牌：%x 新牌：%x", u.ChairId, CardData, room.SendCardData)
	return true
}

func (room *RoomData) SendReplaceCard(ReplaceUser, ReplaceCard, NewCard int, IsInitFlower bool) {
	room.MjBase.UserMgr.SendMsgAll(&msg.G2C_ReplaceCard{
		ReplaceUser:  ReplaceUser,
		ReplaceCard:  ReplaceCard,
		NewCard:      NewCard,
		IsInitFlower: IsInitFlower,
	})
}

//开始发牌
func (room *RoomData) StartDispatchCard() {
	log.Debug("begin StartDispatchCard")
	userMgr := room.MjBase.UserMgr
	gameLogic := room.MjBase.LogicMgr

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, US_PLAYING)
	})

	var minSice int
	UserCnt := userMgr.GetMaxPlayerCnt()
	room.SiceCount, minSice = room.GetSice()

	//分发扑克
	userMgr.ForEachUser(func(u *user.User) {
		for i := 0; i < room.GetCfg().MaxCount-1; i++ {
			cardIdx := gameLogic.SwitchToCardIndex(room.GetHeadCard())
			room.CardIndex[u.ChairId][cardIdx]++
		}
	})

	room.SendCardData = room.GetHeadCard()
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]++
	room.ProvideCard = room.SendCardData
	room.ProvideUser = room.BankerUser
	room.CurrentUser = room.BankerUser

	//TODO 测试用
	log.Debug("begin reoakce test card ======= %v ", conf.Test)
	if conf.Test {
		log.Debug("begin reoakce test card ======= ")
		room.ReplaceCard()
	}
	//起手糊
	//newCard := make([]int, room.GetCfg().MaxIdx)
	//newCard[gameLogic.SwitchToCardIndex(0x5)] = 3
	//newCard[gameLogic.SwitchToCardIndex(0x8)] = 3
	//newCard[gameLogic.SwitchToCardIndex(0x11)] = 3
	//newCard[gameLogic.SwitchToCardIndex(0x13)] = 3
	//newCard[gameLogic.SwitchToCardIndex(0x21)] = 2
	//room.SendCardData = 0x21
	//room.ProvideCard = 0x21
	//room.CardIndex[room.BankerUser] = newCard
	//room.RepertoryCard[55] = 0x1

	//堆立信息
	SiceCount := LOBYTE(room.SiceCount) + HIBYTE(room.SiceCount)
	TakeChairID := (room.BankerUser + SiceCount - 1) % UserCnt
	TakeCount := room.GetCfg().MaxRepertory - room.GetLeftCard()
	for i := 0; i < UserCnt; i++ {
		//计算数目
		var ValidCount int
		if i == 0 {
			ValidCount = HEAP_FULL_COUNT - room.HeapCardInfo[TakeChairID][1] - (minSice)*2
		} else {
			ValidCount = HEAP_FULL_COUNT - room.HeapCardInfo[TakeChairID][1]
		}

		RemoveCount := int(math.Min(float64(ValidCount), float64(TakeCount)))

		//提取扑克
		TakeCount -= RemoveCount
		if i == 0 {
			room.HeapCardInfo[TakeChairID][1] += RemoveCount
		} else {
			room.HeapCardInfo[TakeChairID][0] += RemoveCount
		}

		//完成判断
		if TakeCount == 0 {
			room.HeapHead = TakeChairID
			room.HeapTail = (room.BankerUser + SiceCount - 1) % UserCnt
			break
		}
		//切换索引
		TakeChairID = (TakeChairID + UserCnt - 1) % UserCnt
	}

	room.UserAction = make([]int, UserCnt)

	return
}

//庄家开局动作
func (room *RoomData) InitBankerAction() {
	gameLogic := room.MjBase.LogicMgr
	userMgr := room.MjBase.UserMgr

	//杠牌判断
	gangCardResult := &TagGangCardResult{}
	room.UserAction[room.BankerUser] |= gameLogic.AnalyseGangCard(room.CardIndex[room.BankerUser], nil, 0, gangCardResult)
	//胡牌判断
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]--
	hu, _ := gameLogic.AnalyseChiHuCard(room.CardIndex[room.BankerUser], []*msg.WeaveItem{}, room.SendCardData)
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]++
	log.Debug("InitBankerAction hu=%v, UserAction=%v, gangCardResult=%v", hu, room.UserAction[room.BankerUser], gangCardResult)
	log.Debug("InitBankerAction user card :=== %v", room.CardIndex[room.BankerUser])
	if hu {
		room.UserAction[room.BankerUser] |= WIK_CHI_HU
	}
	if room.UserAction[room.BankerUser] != WIK_NULL {
		u := userMgr.GetUserByChairId(room.BankerUser)
		actionCard := room.SendCardData
		if room.UserAction[room.BankerUser]&WIK_GANG != 0 {
			actionCard = gangCardResult.CardData[0]
		}
		room.SendOperateNotify(u, actionCard)
	}
}

func (room *RoomData) ReplaceCard() {
	base.GameTestpaiCache.LoadAll()
	for _, v := range base.GameTestpaiCache.All() {
		log.Debug("======%d======%d======%d======%d======%d======%d", v.KindID, room.MjBase.Temp.KindID, v.ServerID, room.MjBase.Temp.ServerID, v.IsAcivate, v.RoomID, room.ID)
		if v.KindID == room.MjBase.Temp.KindID && v.ServerID == room.MjBase.Temp.ServerID && v.IsAcivate == 1 && room.ID == v.RoomID {
			log.Debug("11111111111111111111111111111")
			chairIds := utils.GetStrIntList(v.ChairId, "#")
			if len(chairIds) < 1 {
				break
			}
			var TmpRepertoryCard []int //玩家手里的牌
			cards := strings.Split(v.Cards, "#")
			for i, _ := range cards {
				cards[i] = strings.Replace(cards[i], " ", "", -1)
				cards[i] = strings.Replace(cards[i], ",", "，", -1)
				cards[i] = strings.TrimLeft(cards[i], "，")
				cards[i] = strings.TrimRight(cards[i], "，")
			}
			if len(cards) < len(chairIds) {
				break
			}

			max := room.GetCfg().MaxCount
			mim := room.GetCfg().MaxCount - 1
			testCards := make([]int, 0)
			for idx, value := range chairIds {
				card := utils.GetStrIntSixteenList(cards[idx], "，")
				log.Debug("用户手牌：%v", card)
				if value == 0 && max < len(card) {
					log.Error("给庄家设置的手牌数目过多")
					return
				} else if value != 0 && mim < len(card) {
					log.Error("给其他用户设置的手牌数目过多")
					return
				}
				testCards = append(testCards, card...)
			}
			isRight := room.CheckUserCard(v.KindID, testCards)
			if isRight == false {
				break
			}
			for idx, chair := range chairIds {
				card := utils.GetStrIntSixteenList(cards[idx], "，")
				room.SetUserCard(chair, card)
			}

			bankUser := room.BankerUser
			TmpRepertoryCard = room.GetUserCard(bankUser)
			m := GetCardByIdx(room.ConfigIdx)
			log.Debug("库存的牌%v", m)
			log.Debug("TmpRepertoryCard:%d  %v", len(TmpRepertoryCard), TmpRepertoryCard)

			tempCard := make([]int, len(m))

			room.MjBase.LogicMgr.RandCardList(tempCard, m)

			log.Debug("删除前%d   %v", len(tempCard), tempCard)
			for _, v := range TmpRepertoryCard {
				for idx, card := range tempCard {
					if v == card {
						tempCard = utils.IntSliceDelete(tempCard, idx)
						break
					}
				}
			}
			log.Debug("删除后 %d   %v", len(tempCard), tempCard)

			room.RepertoryCard = make([]int, 0)
			for _, k := range tempCard {
				room.RepertoryCard = append(room.RepertoryCard, k)
			}

			if len(room.RepertoryCard) != room.MinusHeadCount {
				log.Debug("len(room.RepertoryCard) != room.MinusHeadCount")
			}

			for _, v := range TmpRepertoryCard {
				room.RepertoryCard = append(room.RepertoryCard, v)
			}

			log.Debug("库存牌%v", room.RepertoryCard)
			if len(room.RepertoryCard) != room.GetCfg().MaxRepertory {
				log.Debug(" len(room.RepertoryCard) != room.GetCfg().MaxRepertory ")
			}
		}
	}
}

//注意这个函数仅供调试用
func (room *RoomData) SetUserCard(charirID int, cards []int) {
	if charirID >= room.MjBase.UserMgr.GetMaxPlayerCnt() {
		return
	}
	log.Debug("begin SetUserCard len:%d---%v", len(room.CardIndex[charirID]), room.CardIndex[charirID])
	log.Debug("begin chairId %v =========card : %v", charirID, cards)
	gameLogic := room.MjBase.LogicMgr

	inc := 0
	userCard := room.CardIndex[charirID]
	temp := util.CopySlicInt(userCard)
	for idx, cnt := range temp {
		for i := 0; i < cnt; i++ {
			if inc >= len(cards) {
				log.Debug("%d======%d", inc, len(cards))
				break
			}
			userCard[idx]--
			userCard[gameLogic.SwitchToCardIndex(cards[inc])]++
			inc++
		}
	}

	log.Debug("end SetUserCard %v", userCard)
}

//获取用户手牌用于调试
func (room *RoomData) GetUserCard(bankUser int) []int {
	var userCard = make([]int, 0)
	log.Debug("%v", room.CardIndex)
	var Poker int
	for i := 0; i < room.MjBase.UserMgr.GetMaxPlayerCnt(); i++ {
		for key, value := range room.CardIndex[i] {
			if value != 0 {
				for j := 0; j < value; j++ {
					Poker = room.MjBase.LogicMgr.SwitchToCardData(key)
					userCard = append(userCard, Poker)
				}

			}
		}
	}
	count := 0
	var sendCard = make([]int, 0)
	for index, value := range room.CardIndex[bankUser] {
		if value != 0 {
			for j := 0; j < value; j++ {
				Poker = room.MjBase.LogicMgr.SwitchToCardData(index)
				sendCard = append(sendCard, Poker)
				count++
			}
		}
	}
	room.SendCardData = sendCard[count-1]
	log.Debug("Poker Number:%d The Poker is:%v=====%d sendCard:%d----bankUser:%d", count, sendCard, sendCard[count-1], room.SendCardData, bankUser)

	log.Debug("userCard %v", userCard)
	return userCard
}

//用户手牌是否设置错误用于调试
func (room *RoomData) CheckUserCard(KindID int, testCards []int) bool {
	mycard := make(map[int]int)
	for _, value := range testCards {
		mycard[value]++
		if (value <= 0x37 && KindID == 391) || (value <= 0x37 && KindID == 390) {
			if mycard[value] > 4 {
				log.Error("手牌设置出错 cards  ==== card :%d  ## :%v", value, mycard)
				return false
			}
		}
		if value <= 0x29 && KindID == 389 {
			if mycard[value] > 4 {
				log.Error("手牌设置出错 cards  ==== card :%d  ## :%v", value, mycard)
				return false
			}
		}
		if value > 0x29 && KindID == 389 {
			if mycard[value] > 4 {
				log.Error("手牌设置出错 cards  ==== card :%d  ## :%v", value, mycard)
				return false
			}
		}

		if value > 0x37 {
			if mycard[value] > 1 {
				log.Error("手牌设置出错 cards  ==== card :%d  ## :%v", value, mycard)
				return false
			}
		}
		if value < 1 || (value > 0x09 && value < 0x11) || (value > 0x19 && value < 0x21) {
			log.Error("手牌设置出错 cards  ==== card :%d  ## :%v", value, mycard)
			return false
		}
		if KindID == 389 && ((value > 0x29 && value < 0x35) || value > 0x35) {
			log.Error("红中麻将首牌设置出错 cards  ==== card :%d  ## :%v", value, mycard)
			return false
		}

		if KindID == 391 && ((value > 0x29 && value < 0x31) || (value > 0x37 && value < 0x41) || (value > 0x48)) {
			log.Error("漳浦麻将首牌设置出错 cards  ==== card :%d  ## :%v", value, mycard)
			return false
		}
	}
	return true
}

//听牌判断
//func (room *RoomData) CheckTingCard(chairID int) bool {
//	CheckUser := room.MjBase.UserMgr.GetUserByChairId(chairID)
//	if CheckUser == nil {
//		log.Error("at CheckTingCard not found user %d", chairID)
//		return false
//	}
//	if room.Ting[chairID] == false {
//		HuData := &msg.G2C_Hu_Data{OutCardData: make([]int, room.GetCfg().MaxCount), HuCardCount: make([]int, room.GetCfg().MaxCount), HuCardData: make([][]int, room.GetCfg().MaxCount), HuCardRemainingCount: make([][]int, room.GetCfg().MaxCount)}
//		Count := room.MjBase.LogicMgr.AnalyseTingCard(room.CardIndex[chairID], []*msg.WeaveItem{}, HuData.OutCardData, HuData.HuCardCount, HuData.HuCardData)
//		log.Debug("at CheckTingCard Count=%d", Count)
//		if Count > 0 {
//			HuData.OutCardCount = Count
//			room.UserAction[chairID] |= WIK_LISTEN
//			for i := 0; i < room.GetCfg().MaxCount; i++ {
//				if HuData.HuCardCount[i] > 0 {
//					for j := 0; j < HuData.HuCardCount[i]; j++ {
//						HuData.HuCardRemainingCount[i] = append(HuData.HuCardRemainingCount[i], room.GetRemainingCount(chairID, HuData.HuCardData[i][j]))
//					}
//				} else {
//					break
//				}
//			}
//			log.Debug("at CheckTingCard HuData=%v", HuData)
//			CheckUser.WriteMsg(HuData)
//			return true
//		}
//	}
//	return false
//}
func (room *RoomData) CheckTingCard(chairID int) bool {
	return false
}

//向客户端发牌
func (room *RoomData) SendGameStart() {
	//构造变量
	GameStart := &mj_hz_msg.G2C_HZMG_GameStart{}
	GameStart.BankerUser = room.BankerUser
	GameStart.SiceCount = room.SiceCount
	GameStart.HeapHead = room.HeapHead
	GameStart.HeapTail = room.HeapTail
	GameStart.HeapCardInfo = room.HeapCardInfo
	GameStart.PlayCount = room.MjBase.TimerMgr.GetPlayCount()
	//发送数据
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		GameStart.UserAction = room.UserAction[u.ChairId]
		GameStart.CardData = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[u.ChairId])
		u.WriteMsg(GameStart)
	})
}

//正常结束房间
func (room *RoomData) NormalEnd(cbReason int) {
	//变量定义
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	GameConclude := &mj_hz_msg.G2C_GameConclude{}
	GameConclude.ChiHuKind = make([]int, UserCnt)
	GameConclude.CardCount = make([]int, UserCnt)
	GameConclude.HandCardData = make([][]int, UserCnt)
	GameConclude.GameScore = make([]int, UserCnt)
	GameConclude.GangScore = make([]int, UserCnt)
	GameConclude.Revenue = make([]int, UserCnt)
	GameConclude.ChiHuRight = make([]int, UserCnt)
	GameConclude.MaCount = make([]int, UserCnt)
	GameConclude.MaData = make([]int, UserCnt)

	for i := range GameConclude.HandCardData {
		GameConclude.HandCardData[i] = make([]int, room.GetCfg().MaxCount)
	}

	GameConclude.SendCardData = room.SendCardData
	GameConclude.LeftUser = INVALID_CHAIR
	//结束信息
	for i := 0; i < UserCnt; i++ {
		GameConclude.ChiHuKind[i] = room.ChiHuKind[i]
		//权位过滤
		if room.ChiHuKind[i] == WIK_CHI_HU {
			room.FiltrateRight(i, &room.ChiHuRight[i])
			GameConclude.ChiHuRight[i] = room.ChiHuRight[i]
		}
		GameConclude.HandCardData[i] = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[i])
		GameConclude.CardCount[i] = len(GameConclude.HandCardData[i])
	}

	//计算胡牌输赢分
	UserGameScore := make([]int, UserCnt)
	WinCount := room.CalHuPaiScore(UserGameScore)

	////拷贝码数据
	//GameConclude.MaCount = make([]int, 0)
	//nCount := 0
	//if nCount > 1 {
	//	nCount++
	//}
	//for i := 0; i < nCount; i++ {
	//	GameConclude.MaData[i] = room.RepertoryCard[room.MinusLastCount+i]
	//}

	//积分变量
	ScoreInfoArray := make([]*msg.TagScoreInfo, UserCnt)

	GameConclude.ProvideUser = room.ProvideUser
	GameConclude.ProvideCard = room.ProvideCard

	//统计积分
	DetailScore := make([]int, room.MjBase.UserMgr.GetMaxPlayerCnt())
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		if u.Status != US_PLAYING {
			return
		}
		GameConclude.GameScore[u.ChairId] = UserGameScore[u.ChairId]
		//胡牌分算完后再加上杠的输赢分就是玩家本轮最终输赢分
		if WinCount > 0 { //流局不算杠的分数
			GameConclude.GameScore[u.ChairId] += room.UserGangScore[u.ChairId]
			GameConclude.GangScore[u.ChairId] = room.UserGangScore[u.ChairId]
		}

		ScoreInfoArray[u.ChairId] = &msg.TagScoreInfo{}
		ScoreInfoArray[u.ChairId].Revenue = GameConclude.Revenue[u.ChairId]
		ScoreInfoArray[u.ChairId].Score = GameConclude.GameScore[u.ChairId]
		if ScoreInfoArray[u.ChairId].Score > 0 {
			ScoreInfoArray[u.ChairId].Type = SCORE_TYPE_WIN
		} else {
			ScoreInfoArray[u.ChairId].Type = SCORE_TYPE_LOSE
		}

		//历史积分
		room.HistorySe.AllScore[u.ChairId] += GameConclude.GameScore[u.ChairId]
		DetailScore[u.ChairId] = GameConclude.GameScore[u.ChairId]

		//设置玩家积分
		u.Score = int64(room.HistorySe.AllScore[u.ChairId])
	})

	room.HistorySe.DetailScore = append(room.HistorySe.DetailScore, DetailScore)
	GameConclude.Reason = cbReason
	GameConclude.AllScore = room.HistorySe.AllScore
	GameConclude.DetailScore = room.HistorySe.DetailScore
	//发送数据
	room.MjBase.UserMgr.SendMsgAll(GameConclude)

	//写入积分 todo
	room.MjBase.UserMgr.WriteTableScore(ScoreInfoArray, room.MjBase.UserMgr.GetMaxPlayerCnt(), HZMJ_CHANGE_SOURCE)
}

//解散结束
func (room *RoomData) DismissEnd(cbReason int) {
	//变量定义
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	GameConclude := &mj_hz_msg.G2C_GameConclude{}
	GameConclude.ChiHuKind = make([]int, UserCnt)
	GameConclude.CardCount = make([]int, UserCnt)
	GameConclude.HandCardData = make([][]int, UserCnt)
	GameConclude.GameScore = make([]int, UserCnt)
	GameConclude.GangScore = make([]int, UserCnt)
	GameConclude.Revenue = make([]int, UserCnt)
	GameConclude.ChiHuRight = make([]int, UserCnt)
	GameConclude.MaCount = make([]int, UserCnt)
	GameConclude.MaData = make([]int, UserCnt)
	for i, _ := range GameConclude.HandCardData {
		GameConclude.HandCardData[i] = make([]int, room.GetCfg().MaxCount)
	}

	room.BankerUser = INVALID_CHAIR

	GameConclude.SendCardData = room.SendCardData

	//用户扑克
	if len(room.CardIndex) > 0 { //没开始就结束情况下小于0
		for i := 0; i < UserCnt; i++ {
			if len(room.CardIndex[i]) > 0 {
				GameConclude.HandCardData[i] = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[i])
				GameConclude.CardCount[i] = len(GameConclude.HandCardData[i])
			}
		}
	}
	GameConclude.Reason = cbReason
	//发送信息
	room.MjBase.UserMgr.SendMsgAll(GameConclude)
}

func (room *RoomData) GetRemainingCount(ChairId int, cbCardData int) int {
	cbIndex := room.MjBase.LogicMgr.SwitchToCardIndex(cbCardData)
	Count := 0
	for i := room.MinusLastCount; i < room.GetCfg().MaxRepertory-room.MinusHeadCount; i++ {
		if room.RepertoryCard[i] == cbCardData {
			Count++
		}
	}

	for id, cards := range room.CardIndex {
		if id == ChairId {
			continue
		}
		Count += cards[cbIndex]
	}

	return Count
}

//权位过滤
func (room *RoomData) FiltrateRight(wWinner int, chr *int) {
	//自摸
	if wWinner == room.ProvideUser {
		*chr |= CHR_ZI_MO
	} else if room.GangStatus == WIK_MING_GANG {
		*chr |= CHR_QIANG_GANG_HU
	} else {
		log.Debug("AT FiltrateRight")
	}
	return
}

//算分
func (room *RoomData) CalHuPaiScore(EndScore []int) int {
	CellScore := room.Source
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	UserScore := make([]int, UserCnt) //玩家手上分
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		if u.Status != US_PLAYING {
			return
		}
		UserScore[u.ChairId] = int(u.Score)
	})

	WinUser := make([]int, UserCnt)
	WinCount := 0

	for i := 0; i < UserCnt; i++ {
		if WIK_CHI_HU == room.ChiHuKind[(room.BankerUser+i)%UserCnt] {
			WinUser[WinCount] = (room.BankerUser + i) % UserCnt
			WinCount++
		}
	}

	if WinCount > 0 {
		//有人胡牌
		bZiMo := room.ProvideUser == WinUser[0]
		if bZiMo {
			for i := 0; i < UserCnt; i++ {
				if i != WinUser[0] {
					EndScore[i] -= CellScore * 2
					EndScore[WinUser[0]] += CellScore * 2
				}
			}
		} else {
			//抢杠
			for i := 0; i < WinCount; i++ {
				for j := 0; j < UserCnt; j++ {
					if j != WinUser[i] {
						EndScore[WinUser[i]] += CellScore
					}
				}
				EndScore[room.ProvideUser] -= EndScore[WinUser[i]]
			}
		}
		//计算抓花
		room.CalcZhuahua(WinUser)

		//谁胡谁当庄
		room.BankerUser = WinUser[0]
		if WinCount > 1 { //多个玩家胡牌，放炮者当庄
			room.BankerUser = room.ProvideUser
		}
	} else { //荒庄
		room.BankerUser = room.LastCatchCardUser //最后一个摸牌的人当庄
	}

	return WinCount
}

func (room *RoomData) CalcZhuahua(winUser []int) {
}

//计算税收  暂时没有这个 功能
func (room *RoomData) CalculateRevenue(ChairId, lScore int) int {
	return 0
}

//取得扑克
func (room *RoomData) GetSendCard(bTail bool, UserCnt int) int {
	//发送扑克

	var cbSendCardData int
	if bTail {
		cbSendCardData = room.GetLastCard()
	} else {
		cbSendCardData = room.GetHeadCard()
	}

	//堆立信息

	if !bTail {
		//切换索引
		cbHeapCount := room.HeapCardInfo[room.HeapHead][0] + room.HeapCardInfo[room.HeapHead][1]
		if cbHeapCount == HEAP_FULL_COUNT {
			room.HeapHead = (room.HeapHead + UserCnt - 1) % len(room.HeapCardInfo)
		}
		room.HeapCardInfo[room.HeapHead][0]++
	} else {
		//切换索引
		cbHeapCount := room.HeapCardInfo[room.HeapTail][0] + room.HeapCardInfo[room.HeapTail][1]
		if cbHeapCount == HEAP_FULL_COUNT {
			room.HeapTail = (room.HeapTail + 1) % len(room.HeapCardInfo)
		}
		room.HeapCardInfo[room.HeapTail][1]++
	}

	return cbSendCardData
}

func (room *RoomData) CallGangScore() {
	lcell := room.Source
	if room.GangStatus == WIK_FANG_GANG { //放杠一家扣分
		room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
			if u.Status != US_PLAYING {
				return
			}
			if u.ChairId != room.CurrentUser {
				room.UserGangScore[room.ProvideGangUser] -= lcell
				room.UserGangScore[room.CurrentUser] += lcell
			}
		})
	} else if room.GangStatus == WIK_MING_GANG { //明杠每家出1倍
		room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
			if u.Status != US_PLAYING {
				return
			}
			if u.ChairId != room.CurrentUser {
				room.UserGangScore[u.ChairId] -= lcell
				room.UserGangScore[room.CurrentUser] += lcell
			}
		})

		//记录明杠次数
	} else if room.GangStatus == WIK_AN_GANG { //暗杠每家出2倍
		room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
			if u.Status != US_PLAYING {
				return
			}
			if u.ChairId != room.CurrentUser {
				room.UserGangScore[u.ChairId] -= 2 * lcell
				room.UserGangScore[room.CurrentUser] += 2 * lcell
			}
		})
	}
}

func (room *RoomData) IsActionDone() bool {
	return room.UserActionDone
}

func (room *RoomData) GetTrusteeOutCard(wChairID int) int {
	cardindex := INVALID_BYTE
	if room.SendCardData != 0 {
		cardindex = room.MjBase.LogicMgr.SwitchToCardIndex(room.SendCardData)
	} else {
		for i := room.GetCfg().MaxIdx - 1; i > 0; i-- {
			if room.CardIndex[wChairID][i] > 0 {
				cardindex = i
				break
			}
		}
	}
	return cardindex
}

//插花
func (room *RoomData) GetChaHua(u *user.User, setCount int) {
}

//用户听牌
func (room *RoomData) OnUserListenCard(u *user.User, bListenCard bool) bool {
	return true
}

//记录分饼
func (room *RoomData) RecordFollowCard(wTargetUser, cbCenterCard int) bool {
	return true
}

//记录出牌数
func (room *RoomData) RecordOutCarCnt() int {
	room.OutCardCount++
	return room.OutCardCount
}

//出牌禁忌
func (room *RoomData) RecordBanCard(OperateCode, ChairId int) {
	room.BanUser[ChairId] |= OperateCode
}

//吃啥打啥
func (room *RoomData) OutOfChiCardRule(CardData, ChairId int) bool {
	return true
}

////////////////////////////////////////////////////////////////////////////////////////////////////
//胡牌算分类型

//地胡
func (room *RoomData) IsDiHu(pAnalyseItem *TagAnalyseItem) int {
	if room.BankerUser == room.CurrentUser {
		return 0
	}
	//吃碰杠失效
	for k, v := range room.CardIndex {
		if k == room.CurrentUser {
			if len(v) != room.GetCfg().MaxCount {
				return 0
			}
		}
		if len(v) != room.GetCfg().MaxCount-1 {
			return 0
		}
	}

	if room.OutCardCount > 0 && room.OutCardCount <= 4 {
		return CHR_DI_HU
	}

	return 0
}

//天胡
func (room *RoomData) IsTianHu(pAnalyseItem *TagAnalyseItem) int {
	if room.BankerUser != room.CurrentUser {
		return 0
	}

	//吃碰杠失效
	for _, v := range pAnalyseItem.IsAnalyseGet {
		if v == false {
			return 0
		}
	}

	if room.OutCardCount != 0 {
		return 0
	}

	return CHR_TIAN_HU
}

//杠上开花
func (room *RoomData) IsGangKaiHua(pAnalyseItem *TagAnalyseItem, WeaveItem []*msg.WeaveItem) int {
	if room.CurrentUser != room.ProvideUser {
		return 0
	}
	index := len(WeaveItem)
	if index < 1 {
		return 0
	}
	index = index - 1

	log.Debug("########## pAnalyseItem.WeaveKind:%v", pAnalyseItem.WeaveKind)
	log.Debug("########## pAnalyseItem.IsAnalyseGet:%v len:%d", pAnalyseItem.IsAnalyseGet, len(pAnalyseItem.WeaveKind))
	if pAnalyseItem.WeaveKind[index] == WIK_GANG && pAnalyseItem.IsAnalyseGet[index] == false {
		return CHR_GANG_SHANG_HUA
	}
	return 0
}

//花上开花
func (room *RoomData) IsHuaKaiHua(pAnalyseItem *TagAnalyseItem) int {
	if room.CurrentUser != room.ProvideUser {
		return 0
	}

	if room.FlowerCnt[room.CurrentUser] == 0 {
		return 0
	}

	if room.SendStatus != BuHua_Send {
		return 0
	}

	return CHR_HUA_SHANG_HUA
}

//海底捞针
func (room *RoomData) IsHaiDiLaoYue(pAnalyseItem *TagAnalyseItem) int {
	if room.ProvideUser != room.CurrentUser {
		return 0
	}

	if room.GetLeftCard() != room.EndLeftCount {
		return 0
	}
	return CHR_HAI_DI_LAO_ZHEN
}

//字牌刻字
func (room *RoomData) IsKeZi(pAnalyseItem *TagAnalyseItem) (int, int, int) {

	type1Cnt := 0
	type2Cnt := 0
	for k, v := range pAnalyseItem.CenterCard {
		cardColor := v >> 4
		if cardColor == 3 && pAnalyseItem.WeaveKind[k] == WIK_PENG {
			cardColor := pAnalyseItem.CenterCard[k] >> 4
			cardValue := pAnalyseItem.CenterCard[k] & MASK_VALUE
			if cardColor == 3 {
				if cardValue > 4 {
					type1Cnt++
				} else {
					type2Cnt++
				}
			}
		}
	}

	if type1Cnt > 0 || type2Cnt > 0 {
		return CHR_ZI_KE_PAI, type1Cnt, type2Cnt
	}
	return 0, 0, 0
}

//花杠
func (room *RoomData) IsHuaGang(pAnalyseItem *TagAnalyseItem, FlowerCnt []int) int {
	if FlowerCnt[room.CurrentUser] == 8 {
		return CHR_HUA_GANG
	}
	return 0
}

//暗刻
func (room *RoomData) IsAnKe(pAnalyseItem *TagAnalyseItem) int {
	anKeCount := 0
	for k, v := range pAnalyseItem.WeaveKind {
		if v == WIK_PENG && pAnalyseItem.IsAnalyseGet[k] == true {
			anKeCount++
		}
	}
	if anKeCount == 3 {
		return CHR_SAN_AN_KE
	} else if anKeCount == 4 {
		return CHR_SI_AN_KE
	} else if anKeCount == 5 {
		return CHR_WU_AN_KE
	}
	return 0
}

//无花字
func (room *RoomData) IsWuHuaZi(pAnalyseItem *TagAnalyseItem, FlowerCnt []int) int {
	if FlowerCnt[room.CurrentUser] != 0 {
		return 0
	}

	for _, v := range pAnalyseItem.CenterCard {
		cardColor := v >> 4
		if cardColor == 3 {
			return 0
		}
	}

	if pAnalyseItem.CardEye>>4 == 3 {
		return 0
	}
	return CHR_WU_HUA_ZI
}

//对对胡
func (room *RoomData) IsDuiDuiHu(pAnalyseItem *TagAnalyseItem) int {
	for k, v := range pAnalyseItem.WeaveKind {
		if v&(WIK_PENG|WIK_GANG) == 0 && pAnalyseItem.IsAnalyseGet[k] == false {
			return 0
		}
		if v&(WIK_LEFT|WIK_RIGHT|WIK_CENTER) != 0 {
			return 0
		}
	}
	return CHR_PENG_PENG
}

//小四喜
func (room *RoomData) IsXiaoSiXi(pAnalyseItem *TagAnalyseItem) int {
	var colorCount [4]int
	for _, v := range pAnalyseItem.CenterCard {
		cardColor := v >> 4
		if cardColor == 3 {
			cardV := v & MASK_VALUE
			if cardV > 4 {
				return 0
			}
			colorCount[cardV-1] = 1
		}
	}

	if colorCount[0]+colorCount[1]+colorCount[2]+colorCount[3] == 3 {
		for k, v := range colorCount {
			if v == 0 {
				if k+1 == pAnalyseItem.CardEye&MASK_VALUE {
					return CHR_XIAO_SI_XI
				}
			}
		}
	}
	return 0
}

//大四喜
func (room *RoomData) IsDaSiXi(pAnalyseItem *TagAnalyseItem) int {
	var colorCount [4]int
	for _, v := range pAnalyseItem.CenterCard {
		cardColor := v >> 4
		if cardColor == 3 {
			cardV := v & MASK_VALUE
			if cardV > 4 {
				continue
			}
			colorCount[cardV-1] = 1
		}
	}

	if colorCount[0]+colorCount[1]+colorCount[2]+colorCount[3] == 4 {
		return CHR_DA_SI_XI
	}
	return 0
}

//小三元
func (room *RoomData) IsXiaoSanYuan(pAnalyseItem *TagAnalyseItem) int {
	var colorCount [3]int
	for _, v := range pAnalyseItem.CenterCard {
		cardColor := v >> 4
		if cardColor == 3 {
			cardV := (v & MASK_VALUE) - 5
			if cardV < 0 {
				continue
			}
			colorCount[cardV] = 1
		}
	}

	if colorCount[0]+colorCount[1]+colorCount[2] == 2 {
		for k, v := range colorCount {
			if v == 0 {
				if k == (pAnalyseItem.CardEye&MASK_VALUE)-5 {
					return CHR_XIAO_SAN_YUAN
				}
			}
		}
	}
	return 0
}

//大三元
func (room *RoomData) IsDaSanYuan(pAnalyseItem *TagAnalyseItem) int {
	var colorCount [3]int
	for _, v := range pAnalyseItem.CenterCard {
		cardColor := v >> 4
		if cardColor == 3 {
			cardV := (v & MASK_VALUE) - 5
			if cardV < 0 {
				continue
			}
			colorCount[cardV] = 1
		}
	}

	if colorCount[0]+colorCount[1]+colorCount[2] == 3 {
		return CHR_DA_SAN_YUAN
	}
	return 0
}

//混一色
func (room *RoomData) IsHunYiSe(pAnalyseItem *TagAnalyseItem) int {
	cardColor := pAnalyseItem.CardEye >> 4
	var colorCount [4]int
	colorCount[cardColor] = 1
	for _, v := range pAnalyseItem.CenterCard {
		cardColor = v >> 4
		if colorCount[cardColor] == 0 {
			colorCount[cardColor] = 1
		}
	}
	if colorCount[0]+colorCount[1]+colorCount[2] == 1 && colorCount[3] == 1 {
		return CHR_HUN_YI_SE
	}

	return 0
}

//清一色
func (room *RoomData) IsQingYiSe(pAnalyseItem *TagAnalyseItem, FlowerCnt []int) int {
	cardColor := pAnalyseItem.CardEye & MASK_COLOR
	for _, v := range pAnalyseItem.CenterCard {
		if v&MASK_COLOR != cardColor {
			return 0
		}
	}

	if (FlowerCnt[room.CurrentUser] > 0) || (0x30 == cardColor) {
		return 0
	}
	return CHR_QING_YI_SE
}

//花一色
func (room *RoomData) IsHuaYiSe(pAnalyseItem *TagAnalyseItem, FlowerCnt []int) int {
	cardColor := pAnalyseItem.CardEye & MASK_COLOR
	for _, v := range pAnalyseItem.CenterCard {
		if v&MASK_COLOR != cardColor {
			return 0
		}
	}

	if FlowerCnt[room.CurrentUser] > 0 && 0x30 != cardColor {
		return CHR_HUA_YI_SE
	}

	return 0
}

//字一色
func (room *RoomData) IsZiYiSe(pAnalyseItem *TagAnalyseItem, FlowerCnt []int) int {
	if FlowerCnt[room.CurrentUser] != 0 {
		return 0
	}

	for _, v := range pAnalyseItem.CenterCard {
		cardColor := v >> 4
		if cardColor != 3 {
			return 0
		}
	}

	if (pAnalyseItem.CardEye >> 4) != 3 {
		return 0
	}

	return CHR_ZI_YI_SE
}

//门清
func (room *RoomData) IsMenQing(pAnalyseItem *TagAnalyseItem) int {

	for k, v := range pAnalyseItem.WeaveKind {
		if (v&(WIK_LEFT|WIK_CENTER|WIK_RIGHT)) != 0 || v == WIK_PENG {
			if pAnalyseItem.IsAnalyseGet[k] == false {
				return 0
			}
		} else if v == WIK_GANG {
			if !(pAnalyseItem.Param[k] == WIK_AN_GANG && pAnalyseItem.IsAnalyseGet[k] == false) {
				return 0
			}
		}
	}

	return CHR_PING_HU
}

//佰六
func (room *RoomData) IsBaiLiu(pAnalyseItem *TagAnalyseItem, FlowerCnt []int) int {

	if FlowerCnt[room.CurrentUser] > 0 {
		return 0
	}

	HuOfCard := room.HuOfCard
	for k, v := range pAnalyseItem.WeaveKind {
		if (v & (WIK_PENG | WIK_GANG)) > 0 {
			return 0
		} else {
			CenterColor := pAnalyseItem.CenterCard[k] >> 4
			CurColor := HuOfCard >> 4
			if CenterColor != CurColor {
				continue
			}

			if pAnalyseItem.CenterCard[k] == HuOfCard {
				//排除12和89
				if (pAnalyseItem.CardData[k][0]&MASK_VALUE == 1 && pAnalyseItem.CardData[k][1]&MASK_VALUE == 2) ||
					(pAnalyseItem.CardData[k][1]&MASK_VALUE == 8 && pAnalyseItem.CardData[k][2]&MASK_VALUE == 9) {
					continue
				}

				if pAnalyseItem.CardData[k][0] == HuOfCard && pAnalyseItem.CardData[k][2] == HuOfCard+2 {
					return CHR_BAI_LIU
				}

				if pAnalyseItem.CardData[k][0] == HuOfCard-2 && pAnalyseItem.CardData[k][2] == HuOfCard {
					return CHR_BAI_LIU
				}
			}

		}
	}

	return 0
}

//门清佰六
func (room *RoomData) IsMenQingBaiLiu(pAnalyseItem *TagAnalyseItem, FlowerCnt []int) int {
	if room.IsMenQing(pAnalyseItem) > 0 && room.IsBaiLiu(pAnalyseItem, FlowerCnt) > 0 {
		return CHR_QING_BAI_LIU
	}
	return 0
}

//字牌杠
func (room *RoomData) IsZiPaiGang(pAnalyseItem *TagAnalyseItem) (int, int, int) {
	type1Cnt := 0
	type2Cnt := 0
	for k, v := range pAnalyseItem.WeaveKind {
		if v == WIK_GANG && pAnalyseItem.Param[k] != WIK_AN_GANG {
			cardColor := pAnalyseItem.CenterCard[k] >> 4
			cardValue := pAnalyseItem.CenterCard[k] & MASK_VALUE
			if cardColor == 3 {
				if cardValue > 4 {
					type1Cnt++
				} else {
					type2Cnt++
				}
			}
		}
	}

	if type1Cnt > 0 || type2Cnt > 0 {
		return CHR_ZI_PAI_GANG, type1Cnt, type2Cnt
	}
	return 0, 0, 0
}

//胡尾张
func (room *RoomData) IsHuWeiZhang(pAnalyseItem *TagAnalyseItem) int {
	logic := room.MjBase.LogicMgr

	index := 0
	count := 0
	for k, v := range room.CardIndex[room.CurrentUser] {
		count += v
		if v > 0 {
			index = k
		}
	}

	if count != 1 {
		return 0
	}

	if pAnalyseItem.CardEye == logic.SwitchToCardData(index) {
		return CHR_WEI_ZHANG
	}
	return 0
}

//截头
func (room *RoomData) IsJieTou(pAnalyseItem *TagAnalyseItem, TingCnt []int) int {
	if TingCnt[room.CurrentUser] != 1 {
		return 0
	}

	HuOfCard := room.HuOfCard
	cardValue := HuOfCard & MASK_VALUE
	for k, v := range pAnalyseItem.WeaveKind {
		if v&(WIK_LEFT|WIK_CENTER|WIK_RIGHT) == 0 {
			continue
		} else if pAnalyseItem.IsAnalyseGet[k] {
			//1-2 胡3
			if cardValue == 3 && HuOfCard == pAnalyseItem.CardData[k][2] {
				return CHR_JIE_TOU
			}
			//8-9 胡7
			if cardValue == 7 && HuOfCard == pAnalyseItem.CardData[k][0] {
				return CHR_JIE_TOU
			}
		}

	}
	return 0
}

//空心
func (room *RoomData) IsKongXin(pAnalyseItem *TagAnalyseItem, TingCnt []int) int {
	if TingCnt[room.CurrentUser] != 1 {
		return 0
	}

	HuOfCard := room.HuOfCard
	for k, v := range pAnalyseItem.WeaveKind {
		if v&(WIK_LEFT|WIK_CENTER|WIK_RIGHT) == 0 {
			continue
		} else if pAnalyseItem.IsAnalyseGet[k] {
			if HuOfCard == pAnalyseItem.CardData[k][1] {
				return CHR_KONG_XIN
			}
		}
	}
	return 0
}

//单吊
func (room *RoomData) IsDanDiao(pAnalyseItem *TagAnalyseItem, TingCnt []int) int {
	if TingCnt[room.CurrentUser] != 1 {
		return 0
	}

	HuOfCard := room.HuOfCard
	if pAnalyseItem.CardEye == HuOfCard {
		return CHR_DAN_DIAO
	}
	return 0
}

//自摸
func (room *RoomData) IsZiMo() int {
	if room.OutCardCount == 0 {
		return 0
	}
	if room.CurrentUser == room.ProvideUser {
		return CHR_ZI_MO
	}
	return 0
}

//平胡
func (room *RoomData) IsPingHu() int {
	if room.OutCardCount == 0 {
		return 0
	}
	if room.CurrentUser != room.ProvideUser {
		return CHR_PING_HU
	}
	return 0
}

//胡牌算分类型
////////////////////////////////////////////////////////////////////////////////////////////////////

func (room *RoomData) OnZhuaHua(winUser []int) (CardData [][]int, BuZhong []int) {
	log.Error("at base OnZhuaHua")
	return
}

///// 非接口类型函数
func (room *RoomData) ClearStatus() {
	//用户状态
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	room.IsResponse = make([]bool, UserCnt)
	room.UserAction = make([]int, UserCnt)
	room.PerformAction = make([]int, UserCnt)
}

//是否足够的牌
func (room *RoomData) IsEnoughCard() bool {
	return room.GetLeftCard() > room.EndLeftCount
}

//获取剩下的牌
func (room *RoomData) GetLeftCard() int {
	return room.MinusHeadCount - room.MinusLastCount
}

//正常取牌
func (room *RoomData) GetHeadCard() int {
	room.MinusHeadCount--
	if room.MinusHeadCount < room.MinusLastCount {
		log.Error("at GetHeadCard room.MinusHeadCoun out index")
		return -1
	}
	log.Debug("get card == :%d", room.RepertoryCard[room.MinusHeadCount])
	return room.RepertoryCard[room.MinusHeadCount]
}

//用于补牌
func (room *RoomData) GetLastCard() (card int) {
	if room.MinusLastCount >= room.MinusHeadCount {
		log.Error("at GetLastCard room.MinusHeadCount out index")
		return -1
	}
	card = room.RepertoryCard[room.MinusLastCount]
	log.Debug("get card == :%d", card)
	room.MinusLastCount++
	return
}

//选举庄家
func (room *RoomData) ElectionBankerUser() {
	userMgr := room.MjBase.UserMgr
	OwnerUser, _ := room.MjBase.UserMgr.GetUserByUid(room.CreatorUid)
	if room.BankerUser == INVALID_CHAIR && room.MjBase.Temp.GameType == GAME_GENRE_ZhuanShi { //房卡模式下先把庄家给房主
		if OwnerUser != nil {
			room.BankerUser = OwnerUser.ChairId
		} else {
			log.Error("get bamkerUser error at StartGame")
		}
	} else {
		//庄家设置
		if room.BankerUser == room.ProvideUser {
			room.BankerUser = room.ProvideUser
			room.ChangeBanker = false
		} else {
			room.BankerUser = (room.BankerUser - 1 + userMgr.GetCurPlayerCnt()) % userMgr.GetCurPlayerCnt()
			room.ChangeBanker = true
		}
		//room.BankerUser, _ = utils.RandInt(0, room.MjBase.UserMgr.GetCurPlayerCnt())
	}
}

func (room *RoomData) IsHua(cardData int) bool {
	return cardData <= 0x48 && cardData >= 0x41
}
