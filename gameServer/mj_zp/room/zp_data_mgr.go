package room

import (
	"math"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"mj/common/msg/mj_zp_msg"

	"encoding/json"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/timer"
	"github.com/lovelly/leaf/util"
)

type ZP_RoomData struct {
	*mj_base.RoomData
	ChaHuaTime *timer.Timer //插花时间
	ZhuaHuaCnt int          //抓花个数
	WithZiCard bool         //带字牌
	ScoreType  int          //算分制式

	FollowCard   []int       //跟牌
	IsFollowCard bool        //是否跟牌
	FlowerCnt    [4]int      //补花数
	ChaHuaMap    map[int]int //插花数

	HuKindType      []int       //胡牌类型
	HuKindScore     map[int]int //特殊胡牌分
	ZhuaHuaScore    int         //插花得分
	FollowCardScore []int       //跟牌得分
}

func NewDataMgr(id, uid, configIdx int, name string, temp *base.GameServiceOption, base *ZP_base, set string) *ZP_RoomData {
	r := new(ZP_RoomData)
	r.ChaHuaMap = make(map[int]int)
	r.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base)

	//房间游戏设置
	info := make(map[string]interface{})
	err := json.Unmarshal([]byte(set), &info)
	if err != nil {
		log.Error("at NewDataMgr error:%s", err.Error())
		return nil
	}

	getData, ok := info["ZhuaHua"].(float64)
	if !ok {
		log.Error("zpmj at NewDataMgr [ZhuaHua] error")
		return nil
	}
	r.ZhuaHuaCnt = int(getData)

	getData2, ok := info["WithZiCard"].(bool)
	if !ok {
		log.Error("zpmj at NewDataMgr [WithZiCard] error")
		return nil
	}
	r.WithZiCard = getData2

	getData3, ok := info["ScoreType"].(float64)
	if !ok {
		log.Error("zpmj at NewDataMgr [ScoreType] error")
		return nil
	}
	r.ScoreType = int(getData3)

	return r
}

func (room *ZP_RoomData) InitRoom(UserCnt int) {
	//初始化
	log.Debug("初始化漳浦房间")
	room.RepertoryCard = make([]int, room.GetCfg().MaxRepertory)
	room.CardIndex = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.CardIndex[i] = make([]int, room.GetCfg().MaxIdx)
	}
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
	for i := 0; i < UserCnt; i++ {
		room.HeapCardInfo[i] = make([]int, 2)
	}

	room.LeftCardCount = room.GetCfg().MaxRepertory
	room.UserActionDone = false
	room.SendStatus = Not_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideGangUser = INVALID_CHAIR
	room.HistoryScores = make([]*mj_base.HistoryScore, UserCnt)

	//设置漳浦麻将牌数据
	room.EndLeftCount = 16
	room.FollowCard = make([]int, 60)
}

func (room *ZP_RoomData) BeforeStartGame(UserCnt int) {
	room.InitRoom(UserCnt)
}

func (room *ZP_RoomData) StartGameing() {
	log.Debug("开始漳浦游戏")
	if room.MjBase.TimerMgr.GetPlayCount() == 0 {
		room.MjBase.UserMgr.SendMsgAll(&mj_zp_msg.G2C_MJZP_GetChaHua{})
		//room.ChaHuaTime = room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OutCardTime)*time.Second, func() {
		//room.ChaHuaTime = room.MjBase.AfterFunc(time.Duration(0)*time.Second, func() {
		log.Debug("超时插花")
		//洗牌
		room.StartDispatchCard()
		//向客户端发牌
		room.SendGameStart()
		//开局补花
		room.InitBuHua()
		//庄家开局动作
		room.InitBankerAction()
		//检查自摸
		room.CheckZiMo()
		//})
	} else {
		room.StartDispatchCard()
		//向客户端发牌
		room.SendGameStart()
		//开局补花
		room.InitBuHua()
		//庄家开局动作
		room.InitBankerAction()
		//检查自摸
		room.CheckZiMo()
	}
}

func (room *ZP_RoomData) AfterStartGame() {

}

//获得插花
func (room *ZP_RoomData) GetChaHua(u *user.User, setCount int) {
	room.ChaHuaMap[u.ChairId] = setCount
	if len(room.ChaHuaMap) == 4 {
		room.StartDispatchCard()
		//向客户端发牌
		room.SendGameStart()
		//开局补花
		room.InitBuHua()
		//庄家开局动作
		room.InitBankerAction()
		//检查自摸
		room.CheckZiMo()
	}
}

//用户补花
func (room *ZP_RoomData) OnUserReplaceCard(u *user.User, CardData int) bool {
	gameLogic := room.MjBase.LogicMgr
	if gameLogic.RemoveCard(room.CardIndex[u.ChairId], CardData) == false {
		log.Debug("[用户补花] 用户：%d补花失败", u.ChairId)
		return false
	}

	//记录补花
	room.FlowerCnt[u.ChairId]++

	//是否花杠
	if room.FlowerCnt[u.ChairId] == 8 {
		room.MjBase.OnEventGameConclude(u.ChairId, u, GER_NORMAL)
	}

	//状态变量
	room.SendStatus = BuHua_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideUser = INVALID_CHAIR

	//派发扑克
	room.DispatchCardData(u.ChairId, true)

	outData := &mj_zp_msg.G2C_MJZP_ReplaceCard{}
	outData.IsInitFlower = false
	outData.ReplaceUser = u.ChairId
	outData.ReplaceCard = CardData
	outData.NewCard = room.SendCardData
	room.MjBase.UserMgr.SendMsgAll(&outData)

	log.Debug("[用户补花] 用户：%d,花牌：%x 新牌：%x", u.ChairId, CardData, room.SendCardData)
	return true
}

//用户听牌
func (room *ZP_RoomData) OnUserListenCard(u *user.User, bListenCard bool) bool {
	gameLogic := room.MjBase.LogicMgr

	if bListenCard {
		if WIK_LISTEN == gameLogic.AnalyseTingCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], nil, nil, nil, room.GetCfg().MaxCount) {
			room.Ting[u.ChairId] = true
			//发给消息
			room.MjBase.UserMgr.SendMsgAllNoSelf(u.GetUid(), &mj_zp_msg.G2C_MJZP_ListenCard{
				ListenUser: u.ChairId,
				IsListen:   true,
			})

			//计算胡几张字
			sendData := &mj_zp_msg.G2C_MJZP_ListenCard{}
			sendData.ListenUser = u.ChairId
			sendData.IsListen = true
			res := gameLogic.GetHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], sendData.HuCardData, room.GetCfg().MaxCount)
			sendData.HuCardCount = res
			u.WriteMsg(sendData)
		} else {
			return false
		}
	} else {
		room.Ting[u.ChairId] = false
		sendData := &mj_zp_msg.G2C_MJZP_ListenCard{}
		sendData.ListenUser = u.ChairId
		sendData.IsListen = false
		room.MjBase.UserMgr.SendMsgAll(sendData)
		return true
	}
	return false
}

//剔除大字
func (room *ZP_RoomData) RemoveAllZiCar(NewDataArray, OriDataArray []int) {
	for _, v := range OriDataArray {
		if v >= 0x31 && v <= 0x43 {
			continue
		}
		NewDataArray = append(NewDataArray, v)
	}
}

//开局补花
func (room *ZP_RoomData) InitBuHua() {
	log.Debug("开局补花")
	playerIndex := room.BankerUser
	playerCNT := room.MjBase.UserMgr.GetMaxPlayerCnt()
	for i := 0; i < playerCNT; i++ {
		if playerIndex >= 3 {
			playerIndex = 0
		}

		outData := &mj_zp_msg.G2C_MJZP_ReplaceCard{}
		outData.ReplaceUser = playerIndex
		outData.IsInitFlower = true

		for j := room.GetCfg().MaxIdx - room.GetCfg().HuaIndex; j < room.GetCfg().MaxIdx; j++ {
			if room.CardIndex[playerIndex][j] == 1 {
				index := j
				for {
					outData.NewCard = room.GetSendCard(true, playerCNT)
					newCardIndex := SwitchToCardIndex(outData.NewCard)
					outData.ReplaceCard = SwitchToCardData(index)
					room.MjBase.UserMgr.SendMsgAll(&outData)
					log.Debug("玩家%d,j:%d 补花：%x，新牌：%x", playerIndex, j, SwitchToCardData(index), outData.NewCard)
					room.FlowerCnt[playerIndex]++
					room.CardIndex[playerIndex][newCardIndex]++
					if newCardIndex < (room.GetCfg().MaxIdx - room.GetCfg().HuaIndex) {
						room.CardIndex[playerIndex][j]--
						break
					} else {
						index = newCardIndex
					}
				}
			}
		}
		playerIndex++
	}
}

//庄家开局动作
func (room *ZP_RoomData) InitBankerAction() {
	log.Debug("庄家开局动作")
	userMgr := room.MjBase.UserMgr
	UserCnt := userMgr.GetMaxPlayerCnt()
	gameLogic := room.MjBase.LogicMgr
	room.UserAction = make([]int, UserCnt)

	//测试手牌
	var temp []int
	temp = make([]int, 42)
	temp[0] = 3
	temp[1] = 3
	temp[2] = 3
	temp[3] = 3
	temp[4] = 3
	temp[5] = 2
	room.CardIndex[room.BankerUser] = temp
	GetCardWordArray(room.CardIndex[room.BankerUser])

	log.Debug("---------------------------------------------------")
	gangCardResult := &mj_base.TagGangCardResult{}
	room.UserAction[room.BankerUser] |= gameLogic.AnalyseGangCard(room.CardIndex[room.BankerUser], nil, 0, gangCardResult)

	//胡牌判断
	chr := 0
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]--
	huKind, _ := gameLogic.AnalyseChiHuCard(room.CardIndex[room.BankerUser], []*msg.WeaveItem{}, room.SendCardData, chr, room.GetCfg().MaxCount, true)
	room.UserAction[room.BankerUser] |= huKind
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]++

	if room.UserAction[room.BankerUser] != 0 {
		outData := &mj_zp_msg.C2G_MJZP_OperateNotify{}
		outData.ActionCard = room.SendCardData
		outData.ActionMask = room.UserAction[room.BankerUser]
		userMgr.SendMsgAll(&outData)
	}
}

//发牌
func (room *ZP_RoomData) StartDispatchCard() {
	log.Debug("开始发牌")
	userMgr := room.MjBase.UserMgr
	gameLogic := room.MjBase.LogicMgr

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, US_PLAYING)
	})

	var minSice int
	UserCnt := userMgr.GetMaxPlayerCnt()
	room.SiceCount, minSice = room.GetSice()

	gameLogic.RandCardList(room.RepertoryCard, mj_base.GetCardByIdx(room.ConfigIdx))

	//剔除大字
	if room.WithZiCard == false {
		var tempCard []int
		room.RemoveAllZiCar(tempCard, room.RepertoryCard)
		room.RepertoryCard = tempCard
		log.Debug("剔除大字")
	}

	//分发扑克
	userMgr.ForEachUser(func(u *user.User) {
		for i := 0; i < room.GetCfg().MaxCount-1; i++ {
			room.LeftCardCount--
			room.MinusHeadCount++
			setIndex := SwitchToCardIndex(room.RepertoryCard[room.LeftCardCount])
			room.CardIndex[u.ChairId][setIndex]++
		}
	})

	OwnerUser, _ := userMgr.GetUserByUid(room.CreateUser)
	if room.BankerUser == INVALID_CHAIR && (room.MjBase.Temp.ServerType&GAME_GENRE_PERSONAL) != 0 { //房卡模式下先把庄家给房主
		if OwnerUser != nil {
			room.BankerUser = OwnerUser.ChairId
		} else {
			log.Error("get bamkerUser error at StartGame")
		}
	}

	if room.BankerUser == INVALID_CHAIR {
		room.BankerUser = util.RandInterval(0, UserCnt-1)
	}

	if room.BankerUser >= UserCnt {
		log.Error(" room.BankerUser >=UserCnt %d,  %d", room.BankerUser, UserCnt)
	}

	room.MinusHeadCount++
	room.SendCardData = room.RepertoryCard[room.LeftCardCount]
	room.LeftCardCount--

	room.CardIndex[room.BankerUser][SwitchToCardIndex(room.SendCardData)]++
	room.ProvideCard = room.SendCardData
	room.ProvideUser = room.BankerUser
	room.CurrentUser = room.BankerUser

	//堆立信息
	SiceCount := LOBYTE(room.SiceCount) + HIBYTE(room.SiceCount)
	TakeChairID := (room.BankerUser + SiceCount - 1) % UserCnt
	TakeCount := room.GetCfg().MaxRepertory - room.LeftCardCount
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
	return
}

//响应判断
func (room *ZP_RoomData) EstimateUserRespond(wCenterUser int, cbCenterCard int, EstimatKind int) bool {
	//变量定义
	bAroseAction := false

	//用户状态
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	room.UserAction = make([]int, UserCnt)

	//动作判断
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		//用户过滤
		if wCenterUser == u.ChairId || room.MjBase.UserMgr.IsTrustee(u.ChairId) {
			return
		}

		//出牌类型
		if EstimatKind == EstimatKind_OutCard {
			//吃碰判断
			if u.UserLimit&LimitPeng == 0 {
				//碰牌判断
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimatePengCard(room.CardIndex[u.ChairId], cbCenterCard)

				//吃牌判断
				eatUser := (wCenterUser + 4 - 1) % 4 //4==GAME_PLAYER
				if eatUser == u.ChairId {
					room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimateEatCard(room.CardIndex[u.ChairId], cbCenterCard)
				}
			}

			//杠牌判断
			if room.LeftCardCount > room.EndLeftCount && u.UserLimit&LimitGang == 0 {
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimateGangCard(room.CardIndex[u.ChairId], cbCenterCard)
			}

			//吃胡判断
			for i := 0; i < 4; i++ {
				if i == wCenterUser {
					continue
				}
				if u.UserLimit&LimitChiHu == 0 {
					chr := 0
					huKind, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], cbCenterCard, chr, room.GetCfg().MaxCount, false)
					if huKind == WIK_CHI_HU {
						room.UserAction[u.ChairId] |= WIK_CHI_HU
					}
				}
			}
		}

		//检查抢杠胡
		if EstimatKind == EstimatKind_GangCard {
			//只有庄家和闲家之间才能放炮
			MogicCard := room.MjBase.LogicMgr.SwitchToCardData(room.MjBase.LogicMgr.GetMagicIndex())
			if room.MjBase.LogicMgr.GetMagicIndex() == room.GetCfg().MaxIdx || (room.MjBase.LogicMgr.GetMagicIndex() != room.GetCfg().MaxIdx && cbCenterCard != MogicCard) {
				if u.UserLimit|LimitChiHu == 0 {
					//吃胡判断
					chr := 0
					huKind, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], cbCenterCard, chr, room.GetCfg().MaxCount, false)
					room.UserAction[u.ChairId] |= huKind
				}
			}
		}

		//结果判断
		if room.UserAction[u.ChairId] != WIK_NULL {
			bAroseAction = true
		}
	})

	//结果处理
	if bAroseAction {
		//设置变量
		room.ProvideUser = wCenterUser
		room.ProvideCard = cbCenterCard
		room.ResumeUser = room.CurrentUser
		room.CurrentUser = INVALID_CHAIR

		//发送提示
		room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
			if room.UserAction[u.ChairId] != WIK_NULL {
				u.WriteMsg(&mj_zp_msg.C2G_MJZP_OperateNotify{
					ActionMask: room.UserAction[u.ChairId],
					ActionCard: room.ProvideCard,
				})
			}
		})
		return true
	}

	if room.GangStatus != WIK_GANERAL {
		room.GangOutCard = true
		room.GangStatus = WIK_GANERAL
		room.ProvideGangUser = INVALID_CHAIR
	} else {
		room.GangOutCard = false
	}

	return false
}

//正常结束房间
func (room *ZP_RoomData) NormalEnd() {
	//变量定义
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	GameConclude := &mj_zp_msg.G2C_ZPMJ_GameConclude{}
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
	room.ChiHuKind = make([]int, UserCnt)
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
	room.CalHuPaiScore(UserGameScore)

	//拷贝码数据
	GameConclude.MaCount = make([]int, 0)

	//积分变量
	ScoreInfoArray := make([]*msg.TagScoreInfo, UserCnt)

	GameConclude.ProvideUser = room.ProvideUser
	GameConclude.ProvideCard = room.ProvideCard

	//统计积分
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		if u.Status != US_PLAYING {
			return
		}
		GameConclude.GameScore[u.ChairId] = UserGameScore[u.ChairId]
		//胡牌分算完后再加上杠的输赢分就是玩家本轮最终输赢分
		GameConclude.GameScore[u.ChairId] += room.UserGangScore[u.ChairId]
		GameConclude.GangScore[u.ChairId] = room.UserGangScore[u.ChairId]

		//收税
		if GameConclude.GameScore[u.ChairId] > 0 && (room.MjBase.Temp.ServerType&GAME_GENRE_GOLD) != 0 {
			GameConclude.Revenue[u.ChairId] = room.CalculateRevenue(u.ChairId, GameConclude.GameScore[u.ChairId])
			GameConclude.GameScore[u.ChairId] -= GameConclude.Revenue[u.ChairId]
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
		if room.HistoryScores[u.ChairId] == nil {
			room.HistoryScores[u.ChairId] = &mj_base.HistoryScore{}
		}
		room.HistoryScores[u.ChairId].TurnScore = GameConclude.GameScore[u.ChairId]
		room.HistoryScores[u.ChairId].CollectScore += GameConclude.GameScore[u.ChairId]

	})

	//发送数据
	room.MjBase.UserMgr.SendMsgAll(GameConclude)

	//写入积分 todo
	room.MjBase.UserMgr.WriteTableScore(ScoreInfoArray, room.MjBase.UserMgr.GetMaxPlayerCnt(), ZPMJ_CHANGE_SOURCE)
}

//进行抓花
func (room *ZP_RoomData) OnZhuaHua(CenterUser int) (getData []int) {
	count := room.ZhuaHuaCnt
	if count == 0 {
		return nil
	}

	//抓花规则
	var getInedx [3]int
	index := [4][3]int{{1, 5, 9}, {0, 2, 6}, {0, 3, 7}, {0, 4, 8}}
	if room.BankerUser == CenterUser {
		getInedx = index[0]
	} else {
		v := math.Abs(float64(room.BankerUser - CenterUser))
		getInedx = index[int(v)]
	}

	sendData := &mj_zp_msg.G2C_ZPMJ_ZhuaHua{}
	for i := 0; i < count; i++ {
		room.LeftCardCount--
		cardData := room.RepertoryCard[room.LeftCardCount]
		cardColor := cardData & MASK_COLOR
		cardValue := cardData & MASK_VALUE
		if cardColor == 3 {
			//东南西北
			if cardValue < 5 {
				if cardValue == getInedx[0] || cardValue == getInedx[1] || cardValue == getInedx[2] {
					sendData.ZhongHua = append(sendData.ZhongHua, cardData)
					room.ZhuaHuaScore++
				}
			} else {
				//中发白
				temp := cardValue - 4
				if temp == getInedx[0] || temp == getInedx[1] || temp == getInedx[2] {
					sendData.ZhongHua = append(sendData.ZhongHua, cardData)
					room.ZhuaHuaScore++
				}
			}
		} else if cardColor >= 0 && cardColor <= 2 {
			if cardValue == getInedx[0] || cardValue == getInedx[1] || cardValue == getInedx[2] {
				sendData.ZhongHua = append(sendData.ZhongHua, cardData)
				room.ZhuaHuaScore++
			}
		}
		sendData.BuZhong = append(sendData.BuZhong, cardData)
	}
	return getData
}

//记录分饼
func (room *ZP_RoomData) RecordFollowCard(cbCenterCard int) bool {
	if room.IsFollowCard {
		return false
	}
	room.FollowCard = append(room.FollowCard, cbCenterCard)

	count := len(room.FollowCard) % 4
	if count == 0 {
		begin := count - 4
		for i := begin; i < count; i++ {
			if room.FollowCard[i] != cbCenterCard {
				room.IsFollowCard = true //取消跟牌
				return false
			}
		}
	}

	times := count / 4
	if times == 0 {
		times = 1
	}
	userCNT := room.MjBase.UserMgr.GetMaxPlayerCnt()
	for i := 0; i < userCNT; i++ {
		if i == room.BankerUser {
			room.FollowCardScore[room.BankerUser] -= 3 * times
			continue
		} else {
			room.FollowCardScore[i] += 1 * times
		}
	}

	return true
}

//设置用户相应牌的操作 ,返回是否可以操作
func (room *ZP_RoomData) CheckUserOperator(u *user.User, userCnt, OperateCode int, OperateCard []int) (int, int) {
	if room.IsResponse[u.ChairId] {
		return -1, u.ChairId
	}
	room.IsResponse[u.ChairId] = true
	room.PerformAction[u.ChairId] = OperateCode
	room.OperateCard[u.ChairId] = OperateCard

	u.UserLimit = 0
	//放弃操作
	if OperateCode == WIK_NULL {
		////禁止这轮吃胡
		if room.HasOperator(u.ChairId, WIK_CHI_HU) {
			u.UserLimit |= LimitChiHu
		}
		//禁止这轮碰
		if room.HasOperator(u.ChairId, WIK_PENG) {
			u.UserLimit |= LimitPeng
		}
		//禁止这轮杠
		if room.HasOperator(u.ChairId, WIK_PENG) {
			u.UserLimit |= LimitGang
		}
	}

	cbTargetAction := OperateCode
	wTargetUser := u.ChairId
	//执行判断
	for i := 0; i < userCnt; i++ {
		//获取动作
		cbUserAction := room.UserAction[i]
		if room.IsResponse[wTargetUser] {
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
		return -1, u.ChairId
	}

	if cbTargetAction == WIK_NULL {
		room.UserAction = make([]int, userCnt)
		room.OperateCard = make([][]int, userCnt)
		room.PerformAction = make([]int, userCnt)
		return cbTargetAction, wTargetUser
	}

	//走到这里一定是所有人都响应完了
	return cbTargetAction, wTargetUser
}

func (room *ZP_RoomData) ZiMo(u *user.User) {
	//普通胡牌
	pWeaveItem := room.WeaveItemArray[u.ChairId]
	if !room.MjBase.LogicMgr.RemoveCard(room.CardIndex[u.ChairId], room.SendCardData) {
		log.Error("not foud card at Operater")
		return
	}
	kind, TagAnalyseItem := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], pWeaveItem, room.SendCardData, room.ChiHuRight[u.ChairId], room.GetCfg().MaxCount, false)
	room.ChiHuKind[u.ChairId] = int(kind)
	room.ProvideCard = room.SendCardData

	//特殊胡牌算分
	room.SpecialCardKind(TagAnalyseItem)
	room.SpecialCardScore()
	return
}

func (room *ZP_RoomData) UserChiHu(wTargetUser, userCnt int) {
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
		chihuKind, TagAnalyseItem := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[wChiHuUser], pWeaveItem, room.OperateCard[wTargetUser][0], room.ChiHuRight[wChiHuUser], room.GetCfg().MaxCount, false)
		room.ChiHuKind[wChiHuUser] = chihuKind

		//特殊胡牌算分
		room.SpecialCardKind(TagAnalyseItem)
		room.SpecialCardScore()

		//插入扑克
		if room.ChiHuKind[wChiHuUser] != WIK_NULL {
			wTargetUser = wChiHuUser
		}
	}
}

//特殊胡牌类型及算分
func (room *ZP_RoomData) SpecialCardKind(TagAnalyseItem []*mj_base.TagAnalyseItem) {
	for _, v := range TagAnalyseItem {
		kind := 0

		kind = room.IsBaiLiu(v) //佰六
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_BL] = 6
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsDaSanYuan(v) //大三元
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_DSY] = 12
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsXiaoSanYuan(v) //小三元
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_XSY] = 6
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHunYiSe(v) //混一色
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_CYS] = 6
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsQingYiSe(v, nil) //清一色
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_QYS] = 24
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHuaYiSe(v) //花一色
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_HYS] = 12
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsGangKaiHua(v) //杠上开花
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_GSKH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHuaKaiHua(v) //花上开花
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_HSKH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsMenQing(v) //门清
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_MQQ] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsMenQingBaiLiu(v) //门清佰六
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_MQBL] = 0
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHuWeiZhang(v) //尾单吊
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_WDD] = 6
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsJieTou(v) //截头
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_DD] = 1
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsDuiDuiHu(v) //对对胡
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_DDH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsTianHu(v) //天胡
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_TH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsDiHu(v) //地胡
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_DH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsKeZi(v) //字牌刻字
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_ZPKZ] = 1
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHaiDiLaoYue(v) //海底捞针
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_HDLZ] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsWuHuaZi(v) //无花字
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_HDLZ] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsAnKe(v) //暗刻
		if kind > 0 {
			room.HuKindScore[IDX_SUB_SCORE_SANAK+kind/8] = 3 * (kind / 4) //2,8,16
			room.HuKindType = append(room.HuKindType, kind)
		}
		//todo,单吊
	}
}

//特殊胡牌分
func (room *ZP_RoomData) SpecialCardScore() {
	if room.ScoreType == GAME_TYPE_33 {
		return
	}

	if room.ScoreType == GAME_TYPE_48 {
		for k := range room.HuKindScore {
			switch k {
			case IDX_SUB_SCORE_ZPKZ:
				room.HuKindScore[k] = 1
			case IDX_SUB_SCORE_HDLZ:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_GSKH:
				room.HuKindScore[k] = 4
			case IDX_SUB_SCORE_HSKH:
				room.HuKindScore[k] = 4
			case IDX_SUB_SCORE_QYS:
				room.HuKindScore[k] = 32
			case IDX_SUB_SCORE_HYS:
				room.HuKindScore[k] = 16
			case IDX_SUB_SCORE_CYS:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_DSY:
				room.HuKindScore[k] = 16
			case IDX_SUB_SCORE_XSY:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_DDH:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_MQQ:
				room.HuKindScore[k] = 4
			case IDX_SUB_SCORE_BL:
				room.HuKindScore[k] = 4
			case IDX_SUB_SCORE_DH:
				room.HuKindScore[k] = 4
			case IDX_SUB_SCORE_TH:
				room.HuKindScore[k] = 4
			case IDX_SUB_SCORE_DD:
				room.HuKindScore[k] = 1
			case IDX_SUB_SCORE_WDD:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_MQBL:
				room.HuKindScore[k] = 12
			case IDX_SUB_SCORE_SANAK:
				room.HuKindScore[k] = 4
			case IDX_SUB_SCORE_SIAK:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_WUAK:
				room.HuKindScore[k] = 16
			}
		}

	} else if room.ScoreType == GAME_TYPE_88 {
		for k := range room.HuKindScore {
			switch k {
			case IDX_SUB_SCORE_ZPKZ:
				room.HuKindScore[k] = 1
			case IDX_SUB_SCORE_HDLZ:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_GSKH:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_HSKH:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_QYS:
				room.HuKindScore[k] = 32
			case IDX_SUB_SCORE_HYS:
				room.HuKindScore[k] = 16
			case IDX_SUB_SCORE_CYS:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_DSY:
				room.HuKindScore[k] = 16
			case IDX_SUB_SCORE_XSY:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_DDH:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_MQQ:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_BL:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_DH:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_TH:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_DD:
				room.HuKindScore[k] = 0
			case IDX_SUB_SCORE_WDD:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_MQBL:
				room.HuKindScore[k] = 12
			case IDX_SUB_SCORE_SANAK:
				room.HuKindScore[k] = 8
			case IDX_SUB_SCORE_SIAK:
				room.HuKindScore[k] = 16
			case IDX_SUB_SCORE_WUAK:
				room.HuKindScore[k] = 32
			}
		}
	}
}
