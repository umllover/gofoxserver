package room

import (
	"math"
	. "mj/common/cost"
	"mj/common/msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"mj/common/msg/mj_zp_msg"

	"encoding/json"

	"mj/common/utils"

	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/timer"
	"github.com/lovelly/leaf/util"
)

type ZP_RoomData struct {
	*mj_base.RoomData
	ChaHuaTime  *timer.Timer   //插花定时器
	OutCardTime *timer.Timer   //出牌定时器
	OperateTime []*timer.Timer //操作定时器

	ZhuaHuaCnt int  //抓花个数
	WithZiCard bool //带字牌
	ScoreType  int  //算分制式

	FollowCard   []int       //跟牌
	IsFollowCard bool        //是否跟牌
	FlowerCnt    [4]int      //补花数
	LianZhuang   int         //连庄次数
	ChaHuaMap    map[int]int //插花数
	HuKindType   []int       //胡牌类型
	TingCnt      [4]int      //听牌个数

	ZhuaHuaMap      [16]*mj_zp_msg.HuaUser   //插花数据
	HuKindScore     [4][COUNT_KIND_SCORE]int //特殊胡牌分
	ZhuaHuaScore    int                      //插花得分
	FollowCardScore []int                    //跟牌得分
	SumScore        [4]int                   //游戏总分
}

func NewDataMgr(id, uid, configIdx int, name string, temp *base.GameServiceOption, base *ZP_base, set string) *ZP_RoomData {
	r := new(ZP_RoomData)
	r.ChaHuaMap = make(map[int]int)
	r.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base)

	//房间游戏设置
	info := make(map[string]interface{})
	err := json.Unmarshal([]byte(set), &info)
	if err != nil {
		log.Error("zpmj at NewDataMgr error:%s", err.Error())
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
	log.Debug("zpmj at InitRoom")
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
	for i := 0; i < UserCnt; i++ {
		room.DiscardCard[i] = make([]int, 60)
	}
	room.UserGangScore = make([]int, UserCnt)
	room.WeaveItemArray = make([][]*msg.WeaveItem, UserCnt)
	room.ChiHuRight = make([]int, UserCnt)
	room.HeapCardInfo = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.HeapCardInfo[i] = make([]int, 2)
	}
	room.OperateTime = make([]*timer.Timer, UserCnt)

	room.LeftCardCount = room.GetCfg().MaxRepertory
	room.UserActionDone = false
	room.SendStatus = Not_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideGangUser = INVALID_CHAIR
	room.HistoryScores = make([]*HistoryScore, UserCnt)

	//设置漳浦麻将牌数据
	room.EndLeftCount = 16
	room.IsFollowCard = false
	room.TingCnt = [4]int{}
	room.FollowCard = room.FollowCard[0:0]
	room.ChaHuaMap = make(map[int]int)
	room.HuKindType = room.HuKindType[0:0]
	room.HuKindType = append(room.HuKindType, 1)
	room.FollowCardScore = make([]int, UserCnt)
	room.LianZhuang = 0
	room.ZhuaHuaScore = 0
	room.FlowerCnt = [4]int{}
	room.SumScore = [4]int{}
	room.BanCardCnt = [4][9]int{}
	room.HuKindScore = [4][COUNT_KIND_SCORE]int{}
	room.BanUser = [4]int{}

	room.IsResponse = make([]bool, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.OperateCard[i] = make([]int, 60)
	}
	log.Debug("len1 OperateCard: %d %d", len(room.OperateCard), len(room.OperateCard[1]))
	room.PerformAction = make([]int, UserCnt)
}

func (room *ZP_RoomData) BeforeStartGame(UserCnt int) {
	room.InitRoom(UserCnt)
}

func (room *ZP_RoomData) StartGameing() {
	log.Debug("开始漳浦游戏")
	if room.MjBase.TimerMgr.GetPlayCount() == 0 {
		room.MjBase.UserMgr.SendMsgAll(&mj_zp_msg.G2C_MJZP_NotifiChaHua{})

		room.ChaHuaTime = room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OperateCardTime)*time.Second, func() {
			log.Debug("超时插花")
			for i := 0; i < 4; i++ {
				_, ok := room.ChaHuaMap[i]
				if !ok {
					room.ChaHuaMap[i] = 0
				} else {
				}
			}

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
			//定时
			u := room.MjBase.UserMgr.GetUserByChairId(room.BankerUser)
			room.OutCardTimer(u)
		})
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
		//定时
		u := room.MjBase.UserMgr.GetUserByChairId(room.BankerUser)
		room.OutCardTimer(u)
	}
}

func (room *ZP_RoomData) AfterStartGame() {

}

//获得插花
func (room *ZP_RoomData) GetChaHua(u *user.User, setCount int) {
	room.ChaHuaMap[u.ChairId] = setCount

	sendData := &mj_zp_msg.G2C_MJZP_UserCharHua{}
	sendData.SetCount = setCount
	sendData.Chair = u.ChairId
	room.MjBase.UserMgr.SendMsgAll(sendData)
	if len(room.ChaHuaMap) == 4 {
		room.ChaHuaTime.Stop()

		room.StartDispatchCard()
		//向客户端发牌
		room.SendGameStart()
		//开局补花
		room.InitBuHua()
		//庄家开局动作
		room.InitBankerAction()
		//检查自摸
		room.CheckZiMo()
		//定时
		u := room.MjBase.UserMgr.GetUserByChairId(room.BankerUser)
		room.OutCardTimer(u)
	}
}

//用户补花
func (room *ZP_RoomData) OnUserReplaceCard(u *user.User, CardData int) bool {
	log.Debug("[用户补花开始] 用户：%d补花：%d", u.ChairId, CardData)
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
	room.ProvideGangUser = INVALID_CHAIR

	//派发扑克
	room.DispatchCardData(u.ChairId, true)

	outData := &mj_zp_msg.G2C_MJZP_ReplaceCard{}
	outData.IsInitFlower = false
	outData.ReplaceUser = u.ChairId
	outData.ReplaceCard = CardData
	outData.NewCard = room.SendCardData
	room.MjBase.UserMgr.SendMsgAll(outData)

	log.Debug("[用户补花结束] 用户：%d,花牌：%x 新牌：%x", u.ChairId, CardData, room.SendCardData)
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
			log.Error("zpmj at OnUserListenCard")
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
		if playerIndex > 3 {
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
					room.MjBase.UserMgr.SendMsgAll(outData)
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
	userMgr := room.MjBase.UserMgr
	UserCnt := userMgr.GetMaxPlayerCnt()
	gameLogic := room.MjBase.LogicMgr
	room.UserAction = make([]int, UserCnt)

	////测试手牌
	//var temp []int
	//temp = make([]int, 42)
	//temp[0] = 3
	//temp[1] = 3
	//temp[2] = 3
	//temp[3] = 3
	//temp[4] = 3
	//temp[5] = 2
	//room.CardIndex[room.BankerUser] = temp
	//GetCardWordArray(room.CardIndex[room.BankerUser])

	gangCardResult := &mj_base.TagGangCardResult{}
	room.UserAction[room.BankerUser] |= gameLogic.AnalyseGangCard(room.CardIndex[room.BankerUser], nil, 0, gangCardResult)

	//胡牌判断
	chr := 0
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]--
	huKind, _ := gameLogic.AnalyseChiHuCard(room.CardIndex[room.BankerUser], []*msg.WeaveItem{}, room.SendCardData, chr, room.GetCfg().MaxCount, true)
	room.UserAction[room.BankerUser] |= huKind
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]++

	if room.UserAction[room.BankerUser] != 0 {
		outData := &mj_zp_msg.G2C_MJZP_OperateNotify{}
		outData.ActionCard = room.SendCardData
		outData.ActionMask = room.UserAction[room.BankerUser]
		u := userMgr.GetUserByChairId(room.BankerUser)
		u.WriteMsg(u)
		//定时
		room.OperateCardTimer(u)
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
	log.Debug("牌大小：%v", room.RepertoryCard)
	log.Debug("牌型：%d", room.ConfigIdx)

	//剔除大字
	if room.WithZiCard == false {
		var tempCard []int
		room.RemoveAllZiCar(tempCard, room.RepertoryCard)
		room.RepertoryCard = tempCard
		room.LeftCardCount = room.GetCfg().MaxRepertory - 7*4
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
				//有禁止碰的牌
				if !(room.BanUser[u.ChairId]&LimitPeng != 0 && room.BanCardCnt[u.ChairId][LimitPeng] == cbCenterCard) {
					//碰牌判断
					room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimatePengCard(room.CardIndex[u.ChairId], cbCenterCard)
					room.BanCardCnt[u.ChairId][LimitPeng] = cbCenterCard
				}
			}

			//吃牌判断
			eatUser := (wCenterUser + 4 + 1) % 4 //4==GAME_PLAYER
			if eatUser == u.ChairId {
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimateEatCard(room.CardIndex[u.ChairId], cbCenterCard)
				room.BanCardCnt[u.ChairId][LimitChi] = cbCenterCard
			}

			//杠牌判断
			if room.LeftCardCount > room.EndLeftCount && u.UserLimit&LimitGang == 0 {
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimateGangCard(room.CardIndex[u.ChairId], cbCenterCard)
			}

			//吃胡判断
			if u.ChairId != wCenterUser {
				if u.UserLimit&LimitChiHu == 0 {
					//有禁止吃胡的牌
					if !(room.BanUser[u.ChairId]&LimitChiHu != 0 && room.BanCardCnt[u.ChairId][LimitChiHu] == cbCenterCard) {
						chr := 0
						huKind, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], cbCenterCard, chr, room.GetCfg().MaxCount, false)
						if huKind == WIK_CHI_HU {
							room.UserAction[u.ChairId] |= WIK_CHI_HU
						}
						room.BanCardCnt[u.ChairId][LimitChiHu] = cbCenterCard
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
			//抢杠胡特殊分
			room.HuKindScore[u.ChairId][IDX_SUB_SCORE_QGH] = 3
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
				u.WriteMsg(&mj_zp_msg.G2C_MJZP_OperateNotify{
					ActionMask: room.UserAction[u.ChairId],
					ActionCard: room.ProvideCard,
				})
				//定时
				room.OperateCardTimer(u)
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
	//清理变量
	room.ClearAllTimer()

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
		util.DeepCopy(&GameConclude.ScoreKind, &room.HuKindScore[i]) //游戏得分类型
	}

	//计算胡牌输赢分
	UserGameScore := make([]int, UserCnt)
	room.CalHuPaiScore(UserGameScore)

	//拷贝码数据
	GameConclude.MaCount = make([]int, 0)
	util.DeepCopy(&GameConclude.ZhuaHua, &room.ZhuaHuaMap) //抓花数据

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
		//胡牌分
		GameConclude.GameScore[u.ChairId] += room.UserGangScore[u.ChairId]
		GameConclude.GangScore[u.ChairId] += room.SumScore[u.ChairId]

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
			room.HistoryScores[u.ChairId] = &HistoryScore{}
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
func (room *ZP_RoomData) OnZhuaHua(CenterUser int) (CardData []int, BuZhong []int) {
	log.Debug("进行抓花 user:", CenterUser)
	count := room.ZhuaHuaCnt
	if count == 0 {
		return
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

	for i := 0; i < count; i++ {
		room.LeftCardCount--
		cardData := room.RepertoryCard[room.LeftCardCount]
		cardColor := room.MjBase.LogicMgr.GetCardColor(cardData)
		cardValue := room.MjBase.LogicMgr.GetCardValue(cardData)
		if cardColor == 3 {
			//东南西北
			if cardValue < 5 {
				if cardValue == getInedx[0] || cardValue == getInedx[1] || cardValue == getInedx[2] {
					CardData = append(CardData, cardData)
				}
			} else {
				//中发白
				temp := cardValue - 4
				if temp == getInedx[0] || temp == getInedx[1] || temp == getInedx[2] {
					CardData = append(CardData, cardData)
				}
			}
			room.ZhuaHuaScore++
		} else if cardColor >= 0 && cardColor <= 2 {
			if cardValue == getInedx[0] || cardValue == getInedx[1] || cardValue == getInedx[2] {
				CardData = append(CardData, cardData)
				room.ZhuaHuaScore++
			}
		}
		BuZhong = append(BuZhong, cardData)
	}

	return
}

//记录分饼
func (room *ZP_RoomData) RecordFollowCard(cbCenterCard int) bool {
	if room.IsFollowCard {
		return false
	}

	log.Debug("记录分饼")
	room.FollowCard = append(room.FollowCard, cbCenterCard)

	count := len(room.FollowCard)
	if count == 0 {
		begin := count - 4
		for i := begin; i < count; i++ {
			if room.FollowCard[i] != cbCenterCard {
				room.IsFollowCard = true //取消跟牌
				return false
			}
		}
	} else {
		return true
	}

	log.Debug("有分饼，牌值：%x", cbCenterCard)
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
		//抢杠胡分
		room.HuKindScore[u.ChairId][IDX_SUB_SCORE_QGH] = 0
	}

	cbTargetAction := OperateCode
	wTargetUser := u.ChairId

	//执行判断
	for i := 0; i < userCnt; i++ {
		//获取动作
		cbUserAction := room.UserAction[i]
		if room.IsResponse[i] {
			cbUserAction = room.PerformAction[i]
		} else {
			cbUserAction = room.UserAction[i]
		}

		log.Debug("userd:%d cbUserAction:%d ", i, room.UserAction[i])
		log.Debug("room.PerformAction:%d ", i, room.PerformAction)
		//优先级别
		cbUserActionRank := room.MjBase.LogicMgr.GetUserActionRank(cbUserAction)
		cbTargetActionRank := room.MjBase.LogicMgr.GetUserActionRank(cbTargetAction)
		log.Debug("cbUserActionRank:%d", cbUserActionRank)
		log.Debug("cbTargetActionRank:%d", cbTargetActionRank)

		//动作判断
		if cbUserActionRank > cbTargetActionRank {
			wTargetUser = i
			cbTargetAction = cbUserAction
		}
	}

	log.Debug("wTargetUser:%d u.ChairId：%d", wTargetUser, u.ChairId)
	if room.IsResponse[wTargetUser] == false { //最高权限的人没响应
		return -1, u.ChairId
	}

	if cbTargetAction == WIK_NULL {
		room.IsResponse = make([]bool, userCnt)
		room.UserAction = make([]int, userCnt)
		room.OperateCard = make([][]int, userCnt)
		for i := 0; i < userCnt; i++ {
			room.OperateCard[i] = make([]int, 60)
		}
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
	kind, AnalyseItem := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], pWeaveItem, room.SendCardData, room.ChiHuRight[u.ChairId], room.GetCfg().MaxCount, false)
	room.ChiHuKind[u.ChairId] = int(kind)
	room.ProvideCard = room.SendCardData

	//特殊胡牌类型
	room.SpecialCardKind(AnalyseItem, u.ChairId)
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
		chihuKind, AnalyseItem := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[wChiHuUser], pWeaveItem, room.OperateCard[wTargetUser][0], room.ChiHuRight[wChiHuUser], room.GetCfg().MaxCount, false)
		room.ChiHuKind[wChiHuUser] = chihuKind

		//特殊胡牌类型
		room.SpecialCardKind(AnalyseItem, wChiHuUser)

		//插入扑克
		if room.ChiHuKind[wChiHuUser] != WIK_NULL {
			wTargetUser = wChiHuUser
		}
	}
}

//特殊胡牌类型及算分
func (room *ZP_RoomData) SpecialCardKind(TagAnalyseItem []*TagAnalyseItem, HuUserID int) {
	winScore := room.HuKindScore[HuUserID]
	for _, v := range TagAnalyseItem {
		kind := 0

		kind = room.IsDaSanYuan(v) //大三元
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DSY] = 12
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsXiaoSanYuan(v) //小三元
		if kind > 0 {
			winScore[IDX_SUB_SCORE_XSY] = 6
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHunYiSe(v) //混一色
		if kind > 0 {
			winScore[IDX_SUB_SCORE_CYS] = 6
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind, _ = room.IsQingYiSe(v) //清一色
		if kind > 0 {
			winScore[IDX_SUB_SCORE_QYS] = 24
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHuaYiSe(v) //花一色
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HYS] = 12
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsGangKaiHua(v) //杠上开花
		if kind > 0 {
			winScore[IDX_SUB_SCORE_GSKH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHuaKaiHua(v) //花上开花
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HSKH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsBaiLiu(v) //佰六
		if kind > 0 {
			winScore[IDX_SUB_SCORE_BL] = 6
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsMenQing(v) //门清
		if kind > 0 {
			winScore[IDX_SUB_SCORE_MQQ] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsMenQingBaiLiu(v) //门清佰六
		if kind > 0 {
			winScore[IDX_SUB_SCORE_BL] = 0
			winScore[IDX_SUB_SCORE_MQQ] = 0
			winScore[IDX_SUB_SCORE_MQBL] = 9
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHuWeiZhang(v) //尾单吊
		if kind > 0 {
			winScore[IDX_SUB_SCORE_WDD] = 6
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsJieTou(v) //截头
		if kind > 0 {
			winScore[IDX_SUB_SCORE_JT] = 1
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsKongXin(v) //空心
		if kind > 0 {
			winScore[IDX_SUB_SCORE_KX] = 1
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsDuiDuiHu(v) //对对胡
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DDH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsTianHu(v) //天胡
		if kind > 0 {
			winScore[IDX_SUB_SCORE_TH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsDiHu(v) //地胡
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DH] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsKeZi(v) //字牌刻字
		if kind > 0 {
			winScore[IDX_SUB_SCORE_ZPKZ] = 1
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsHaiDiLaoYue(v) //海底捞针
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HDLZ] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsWuHuaZi(v) //无花字
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HDLZ] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsAnKe(v) //暗刻
		if kind > 0 {
			winScore[IDX_SUB_SCORE_SANAK+kind/8] = 3 * (kind / 4) //2,8,16
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsDaSiXi(v) //大四喜
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DSX] = 24
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsXiaoSiXi(v) //小四喜
		if kind > 0 {
			winScore[IDX_SUB_SCORE_XSX] = 12
			room.HuKindType = append(room.HuKindType, kind)
		}
		kind = room.IsWuHuaZi(v) //无花字
		if kind > 0 {
			winScore[IDX_SUB_SCORE_WHZ] = 3
			room.HuKindType = append(room.HuKindType, kind)
		}
		//自摸
		kind = room.IsZiMo()
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_HDLZ] == 0 && winScore[IDX_SUB_SCORE_GSKH] == 0 && winScore[IDX_SUB_SCORE_HSKH] == 0 {
				winScore[IDX_SUB_SCORE_ZM] = 2
				room.HuKindType = append(room.HuKindType, kind)
			}
		}
	}
	//单吊
	if room.TingCnt[room.CurrentUser] == 1 {
		if room.CurrentUser == room.ProvideUser {
			winScore[IDX_SUB_SCORE_DDPH] = 1
			room.HuKindType = append(room.HuKindType, IDX_SUB_SCORE_DDPH)
		} else {
			winScore[IDX_SUB_SCORE_DDZM] = 1
			room.HuKindType = append(room.HuKindType, IDX_SUB_SCORE_DDZM)
		}
	}
}

//特殊胡牌算分规则
func (room *ZP_RoomData) SpecialCardScore(HuUserID int) {
	winScore := room.HuKindScore[HuUserID]
	if room.ScoreType == GAME_TYPE_33 {
		winScore[IDX_SUB_SCORE_JT] = 0
		winScore[IDX_SUB_SCORE_KX] = 0
		winScore[IDX_SUB_SCORE_DDPH] = 0
		return
	}

	if room.ScoreType == GAME_TYPE_48 {
		for k, v := range winScore {
			if v <= 0 {
				continue
			}

			switch k {
			case IDX_SUB_SCORE_ZPKZ:
				winScore[k] = 1
			case IDX_SUB_SCORE_HDLZ:
				winScore[k] = 8
			case IDX_SUB_SCORE_GSKH:
				winScore[k] = 4
			case IDX_SUB_SCORE_HSKH:
				winScore[k] = 4
			case IDX_SUB_SCORE_QYS:
				winScore[k] = 32
			case IDX_SUB_SCORE_HYS:
				winScore[k] = 16
			case IDX_SUB_SCORE_CYS:
				winScore[k] = 8
			case IDX_SUB_SCORE_DSY:
				winScore[k] = 16
			case IDX_SUB_SCORE_XSY:
				winScore[k] = 8
			case IDX_SUB_SCORE_DDH:
				winScore[k] = 8
			case IDX_SUB_SCORE_MQQ:
				winScore[k] = 4
			case IDX_SUB_SCORE_BL:
				winScore[k] = 4
			case IDX_SUB_SCORE_DH:
				winScore[k] = 4
			case IDX_SUB_SCORE_TH:
				winScore[k] = 4
			case IDX_SUB_SCORE_DDPH:
				winScore[k] = 1
			case IDX_SUB_SCORE_WDD:
				winScore[k] = 8
			case IDX_SUB_SCORE_MQBL:
				winScore[k] = 12
			case IDX_SUB_SCORE_SANAK:
				winScore[k] = 4
			case IDX_SUB_SCORE_SIAK:
				winScore[k] = 8
			case IDX_SUB_SCORE_WUAK:
				winScore[k] = 16
			case IDX_SUB_SCORE_ZM:
				winScore[k] = 1
			case IDX_SUB_SCORE_QGH:
				winScore[k] = 4
			}
		}

	} else if room.ScoreType == GAME_TYPE_88 {
		for k, v := range winScore {
			if v <= 0 {
				continue
			}

			switch k {
			case IDX_SUB_SCORE_ZPKZ:
				winScore[k] = 1
			case IDX_SUB_SCORE_HDLZ:
				winScore[k] = 8
			case IDX_SUB_SCORE_GSKH:
				winScore[k] = 8
			case IDX_SUB_SCORE_HSKH:
				winScore[k] = 8
			case IDX_SUB_SCORE_QYS:
				winScore[k] = 32
			case IDX_SUB_SCORE_HYS:
				winScore[k] = 16
			case IDX_SUB_SCORE_CYS:
				winScore[k] = 8
			case IDX_SUB_SCORE_DSY:
				winScore[k] = 16
			case IDX_SUB_SCORE_XSY:
				winScore[k] = 8
			case IDX_SUB_SCORE_DDH:
				winScore[k] = 8
			case IDX_SUB_SCORE_MQQ:
				winScore[k] = 8
			case IDX_SUB_SCORE_BL:
				winScore[k] = 8
			case IDX_SUB_SCORE_DH:
				winScore[k] = 8
			case IDX_SUB_SCORE_TH:
				winScore[k] = 8
			case IDX_SUB_SCORE_DDPH:
				winScore[k] = 0
			case IDX_SUB_SCORE_WDD:
				winScore[k] = 8
			case IDX_SUB_SCORE_MQBL:
				winScore[k] = 12
			case IDX_SUB_SCORE_SANAK:
				winScore[k] = 8
			case IDX_SUB_SCORE_SIAK:
				winScore[k] = 16
			case IDX_SUB_SCORE_WUAK:
				winScore[k] = 32
			case IDX_SUB_SCORE_JT:
				winScore[k] = 0
			case IDX_SUB_SCORE_KX:
				winScore[k] = 0
			case IDX_SUB_SCORE_ZM:
				winScore[k] = 1
			case IDX_SUB_SCORE_QGH:
				winScore[k] = 8
			}
		}
	}
}

//总得分计算和得分类型统计
func (room *ZP_RoomData) SumGameScore(WinUser []int) {
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	for i := 0; i < UserCnt; i++ {
		playerScore := room.HuKindScore[i]

		//暗杠
		playerScore[IDX_SUB_SCORE_AG] = room.UserGangScore[i]
		room.SumScore[i] += playerScore[IDX_SUB_SCORE_AG]

		var j int
		var winCnt int
		for k := range WinUser {
			if WinUser[k] == i {
				winCnt++
				break
			}
		}
		log.Debug("winCnt:%d", winCnt)
		//黄庄
		if winCnt == 0 {
			return
		}
		if j == len(WinUser) {
			continue
		}

		//基础分
		playerScore[IDX_SUB_SCORE_JC] = 1
		room.SumScore[i] += 1
		//补花得分
		if room.FlowerCnt[i] > 1 {
			if room.FlowerCnt[i] < 8 {
				playerScore[IDX_SUB_SCORE_HUA] = room.FlowerCnt[i]
			} else { //八张花牌
				playerScore[IDX_SUB_SCORE_HUA] = 16
			}
			room.SumScore[i] += playerScore[IDX_SUB_SCORE_HUA]
		}
		//连庄
		if i == room.BankerUser {
			playerScore[IDX_SUB_SCORE_LZ] = room.LianZhuang
			log.Debug("room.ProvideUser:%d", room.ProvideUser)
			room.SumScore[room.ProvideUser] -= room.LianZhuang
			room.SumScore[room.BankerUser] += room.LianZhuang
		} else if room.ProvideUser == room.BankerUser {
			playerScore[IDX_SUB_SCORE_LZ] = room.LianZhuang
			room.SumScore[room.ProvideUser] += room.LianZhuang
			room.SumScore[room.BankerUser] -= room.LianZhuang
		}
		//胡牌类型分+加分项分总和
		for j := IDX_SUB_SCORE_HP; j < COUNT_KIND_SCORE; j++ {
			room.SumScore[i] += playerScore[j]
		}
		//插花分
		if i == room.ProvideUser { //自摸情况
			playerScore[IDX_SUB_SCORE_CH] = room.ChaHuaMap[0] + room.ChaHuaMap[1] + room.ChaHuaMap[2] + room.ChaHuaMap[3]
			room.SumScore[i] += playerScore[IDX_SUB_SCORE_CH]
			for j := 0; j < UserCnt; j++ { //其他玩家扣分
				if j == room.ProvideUser {
					continue
				}
				room.SumScore[j] -= room.ChaHuaMap[i] + room.ChaHuaMap[j]
			}
		} else {
			playerScore[IDX_SUB_SCORE_CH] = room.ChaHuaMap[i] + room.ChaHuaMap[room.ProvideUser]
			room.SumScore[i] += playerScore[IDX_SUB_SCORE_CH]
			room.SumScore[room.ProvideUser] -= room.ChaHuaMap[i] + room.ChaHuaMap[room.ProvideUser]
		}
		//抓花
		playerScore[IDX_SUB_SCORE_ZH] = room.ZhuaHuaScore
		room.SumScore[i] += room.ZhuaHuaScore
		//分饼
		if room.BankerUser == i {
			room.SumScore[i] -= room.FollowCardScore[i]
		} else {
			playerScore[IDX_SUB_SCORE_CH] = room.FollowCardScore[i]
			room.SumScore[i] += room.FollowCardScore[i]
		}
	}
}

func (room *ZP_RoomData) SendStatusPlay(u *user.User) {
	StatusPlay := &msg.G2C_StatusPlay{}
	//自定规则
	StatusPlay.TimeOutCard = room.MjBase.TimerMgr.GetTimeOutCard()
	StatusPlay.TimeOperateCard = room.MjBase.TimerMgr.GetTimeOperateCard()
	StatusPlay.CreateTime = room.MjBase.TimerMgr.GetCreatrTime()

	//规则
	StatusPlay.PlayerCount = room.MjBase.TimerMgr.GetPlayCount()
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	//游戏变量
	StatusPlay.BankerUser = room.BankerUser
	StatusPlay.CurrentUser = room.CurrentUser
	StatusPlay.CellScore = room.Source
	StatusPlay.MagicIndex = room.MjBase.LogicMgr.GetMagicIndex()
	StatusPlay.Trustee = room.MjBase.UserMgr.GetTrustees()
	StatusPlay.HuCardCount = make([]int, room.GetCfg().MaxCount)
	StatusPlay.HuCardData = make([][]int, room.GetCfg().MaxCount)
	for i := 0; i < room.GetCfg().MaxCount; i++ {
		StatusPlay.HuCardData[i] = make([]int, 28)
	}
	StatusPlay.OutCardDataEx = make([]int, room.GetCfg().MaxCount)
	StatusPlay.CardCount = make([]int, UserCnt)
	StatusPlay.TurnScore = make([]int, UserCnt)
	StatusPlay.CollectScore = make([]int, UserCnt)
	StatusPlay.BuHuaCnt = make([]int, UserCnt)
	StatusPlay.ChaHuaCnt = make([]int, UserCnt)

	StatusPlay.ZhuaHuaCnt = room.ZhuaHuaCnt
	for k, v := range room.ChaHuaMap {
		StatusPlay.ChaHuaCnt[k] = v
	}
	for i := 0; i < len(room.FlowerCnt); i++ {
		StatusPlay.BuHuaCnt[i] = room.FlowerCnt[i]
	}

	//状态变量
	StatusPlay.ActionCard = room.ProvideCard
	StatusPlay.LeftCardCount = room.LeftCardCount
	StatusPlay.ActionMask = room.UserAction[u.ChairId]

	StatusPlay.Ting = room.Ting
	//当前能胡的牌
	StatusPlay.OutCardCount = room.MjBase.LogicMgr.AnalyseTingCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId],
		StatusPlay.OutCardDataEx, StatusPlay.HuCardCount, StatusPlay.HuCardData, room.GetCfg().MaxCount)

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
	for j := 0; j < UserCnt; j++ {
		//设置变量
		if room.HistoryScores[j] != nil {
			StatusPlay.TurnScore[j] = room.HistoryScores[j].TurnScore
			StatusPlay.CollectScore[j] = room.HistoryScores[j].CollectScore
		}
	}

	u.WriteMsg(StatusPlay)
}

//算分
func (room *ZP_RoomData) CalHuPaiScore(EndScore []int) {
	//CellScore := room.Source
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	UserScore := make([]int, UserCnt) //玩家手上分
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		if u.Status != US_PLAYING {
			return
		}
		UserScore[u.ChairId] = int(u.Score)
	})

	//WinUser := make([]int, UserCnt)
	var WinUser []int
	WinCount := 0

	for i := 0; i < UserCnt; i++ {
		if WIK_CHI_HU == room.ChiHuKind[(room.BankerUser+i)%UserCnt] {
			WinUser = append(WinUser, (room.BankerUser+i)%UserCnt)
			room.SpecialCardScore(WinUser[WinCount])
			WinCount++
		}
	}

	if WinCount > 0 {
		//插花
		tempZhuaHuaCnt := room.ZhuaHuaCnt
		leftZhuaHuaCnt := room.ZhuaHuaCnt
		huaUser := mj_zp_msg.HuaUser{}
		for k, v := range WinUser {
			//一炮多响，数量随机
			if WinCount > 1 && k < WinCount-1 {
				var error error
				room.ZhuaHuaCnt, error = utils.RandInt(1, leftZhuaHuaCnt)
				if error == nil {
					return
				}
			}

			//进行抓花
			ZhongCard, BuZhong := room.OnZhuaHua(v)

			//抓花派位
			for _, cardV := range ZhongCard {
				for {
					randV, randOk := utils.RandInt(0, 15)
					if randOk == nil && room.ZhuaHuaMap[randV] == nil {
						huaUser.Card = cardV
						log.Debug("中花：%d", cardV)
						huaUser.ChairID = v
						huaUser.IsZhong = true
						room.ZhuaHuaMap[randV] = &huaUser
						break
					}
				}
			}
			for _, cardV2 := range BuZhong {
				for {
					randV, randOk := utils.RandInt(0, 15)
					if randOk == nil && room.ZhuaHuaMap[randV] == nil {
						huaUser.Card = cardV2
						huaUser.ChairID = v
						log.Debug("中花：%d", cardV2)
						huaUser.IsZhong = false
						room.ZhuaHuaMap[randV] = &huaUser
						break
					}
				}
			}

			leftZhuaHuaCnt -= room.ZhuaHuaCnt
		}
		room.ZhuaHuaCnt = tempZhuaHuaCnt

		//连庄次数
		if room.CurrentUser == room.BankerUser {
			room.LianZhuang = 1
		}

		//连庄
		if WinCount > 1 {
			//一炮多响,庄家当庄
			var Zhuang bool
			for _, v := range WinUser {
				if room.BankerUser == v {
					Zhuang = true
				}
			}
			if Zhuang == false {
				room.BankerUser = room.BankerUser + 1
			}
		} else {
			if room.CurrentUser == room.BankerUser {
				room.BankerUser = room.BankerUser
			} else {
				room.BankerUser += 1
			}
		}

	} else { //荒庄
		room.BankerUser = room.LastCatchCardUser //最后一个摸牌的人当庄
	}

	room.SumGameScore(WinUser)
}

//杠计分
func (room *ZP_RoomData) CallGangScore() {
	lcell := room.Source
	//暗杠得分
	if room.GangStatus == WIK_AN_GANG {
		room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
			if u.Status != US_PLAYING {
				return
			}
			if u.ChairId != room.CurrentUser {
				room.UserGangScore[u.ChairId] -= lcell
				room.UserGangScore[room.CurrentUser] += lcell
			}
		})
	}
}

//出牌禁忌
func (room *ZP_RoomData) RecordBanCard(OperateCode, ChairId int) {
	room.BanUser[ChairId] |= OperateCode
}

//清除出牌禁忌
func (room *ZP_RoomData) ClearBanCard(ChairId int) {
	room.BanUser[ChairId] = 0
	room.BanCardCnt[ChairId] = [9]int{}
}

//吃啥打啥
func (room *ZP_RoomData) OutOfChiCardRule(CardData, ChairId int) bool {
	if room.BanUser[ChairId]&LimitChi != 0 && room.BanCardCnt[ChairId][LimitChi] == CardData {
		return false
	}
	return true
}

/////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////
//////////////////与base逻辑一致
func (room *ZP_RoomData) NotifySendCard(u *user.User, cbCardData int, bSysOut bool) {
	//设置变量
	room.SendStatus = OutCard_Send
	room.SendCardData = 0
	room.UserAction[u.ChairId] = WIK_NULL

	//出牌记录
	room.OutCardUser = u.ChairId
	room.OutCardData = cbCardData

	//构造数据
	OutCard := &mj_zp_msg.G2C_ZPMJ_OutCard{}
	OutCard.OutCardUser = u.ChairId
	OutCard.OutCardData = cbCardData
	OutCard.SysOut = bSysOut

	//发送消息
	room.MjBase.UserMgr.SendMsgAll(OutCard)
	room.ProvideUser = u.ChairId
	room.ProvideCard = cbCardData

	//用户切换
	room.CurrentUser = (u.ChairId + 1) % room.MjBase.UserMgr.GetMaxPlayerCnt()
}

func (room *ZP_RoomData) AnGang(u *user.User, cbOperateCode int, cbOperateCard []int) int {
	log.Debug("########## cbOperateCode:%d", cbOperateCode)
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

		Wrave := &msg.WeaveItem{}
		Wrave.Param = WIK_AN_GANG
		Wrave.ProvideUser = u.ChairId
		Wrave.WeaveKind = cbOperateCode
		Wrave.CenterCard = cbOperateCard[0]
		Wrave.CardData = make([]int, 4)
		for j := 0; j < 4; j++ {
			Wrave.CardData[j] = cbOperateCard[0]
		}
		room.WeaveItemArray[u.ChairId] = append(room.WeaveItemArray[u.ChairId], Wrave)
	}

	//删除扑克
	room.CardIndex[u.ChairId][cbCardIndex] = 0
	room.GangStatus = cbGangKind
	room.ProvideGangUser = wProvideUser
	room.GangCard[u.ChairId] = true
	room.GangCount[u.ChairId]++

	//构造结果
	OperateResult := &mj_zp_msg.G2C_ZPMJ_OperateResult{}
	OperateResult.OperateUser = u.ChairId
	OperateResult.ProvideUser = wProvideUser
	OperateResult.OperateCode = cbOperateCode
	OperateResult.OperateCard[0] = cbOperateCard[0]

	//发送消息
	room.MjBase.UserMgr.SendMsgAll(OperateResult)

	return cbGangKind
}

func (room *ZP_RoomData) CallOperateResult(wTargetUser, cbTargetAction int) {
	//构造结果
	OperateResult := &mj_zp_msg.G2C_ZPMJ_OperateResult{}
	OperateResult.OperateUser = wTargetUser
	OperateResult.OperateCode = cbTargetAction
	if room.ProvideUser == INVALID_CHAIR {
		OperateResult.ProvideUser = wTargetUser
	} else {
		OperateResult.ProvideUser = room.ProvideUser
	}

	cbTargetCard := room.OperateCard[wTargetUser][0]
	OperateResult.OperateCard[0] = cbTargetCard
	if cbTargetAction&(WIK_LEFT|WIK_CENTER|WIK_RIGHT) != 0 {
		OperateResult.OperateCard[1] = room.OperateCard[wTargetUser][1]
		OperateResult.OperateCard[2] = room.OperateCard[wTargetUser][2]
	} else if cbTargetAction&WIK_PENG != 0 {
		OperateResult.OperateCard[1] = cbTargetCard
		OperateResult.OperateCard[2] = cbTargetCard
	}

	//用户状态
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	log.Debug("+++++++++++++++++++++UserAction :%v", room.UserAction)
	room.IsResponse = make([]bool, UserCnt)
	room.UserAction = make([]int, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.OperateCard[i] = make([]int, 60)
	}
	room.PerformAction = make([]int, UserCnt)
	log.Debug("状态清除")
	log.Debug("room.IsResponse %v", room.IsResponse)
	log.Debug("room.UserActionResponse %v", room.UserAction)
	log.Debug("room.OperateCard %v", room.OperateCard)
	log.Debug("len1 OperateCard: %d %d", len(room.OperateCard), len(room.OperateCard[1]))

	//如果非杠牌
	if cbTargetAction != WIK_GANG {
		room.ProvideUser = INVALID_CHAIR
		room.ProvideCard = 0

		gcr := &mj_base.TagGangCardResult{}
		room.UserAction[wTargetUser] |= room.MjBase.LogicMgr.AnalyseGangCard(room.CardIndex[wTargetUser], room.WeaveItemArray[wTargetUser], 0, gcr)

		if room.Ting[wTargetUser] == false {
			HuData := &msg.G2C_Hu_Data{OutCardData: make([]int, room.GetCfg().MaxCount), HuCardCount: make([]int, room.GetCfg().MaxCount), HuCardData: make([][]int, room.GetCfg().MaxCount), HuCardRemainingCount: make([][]int, room.GetCfg().MaxCount)}
			for k := 0; k < room.GetCfg().MaxCount; k++ {
				HuData.HuCardData[k] = make([]int, 28)
				HuData.HuCardRemainingCount[k] = make([]int, 28)
			}

			cbCount := room.MjBase.LogicMgr.AnalyseTingCard(room.CardIndex[wTargetUser], room.WeaveItemArray[wTargetUser], HuData.OutCardData, HuData.HuCardCount, HuData.HuCardData, room.GetCfg().MaxCount)
			HuData.OutCardCount = cbCount
			if cbCount > 0 {
				room.UserAction[wTargetUser] |= WIK_LISTEN
				for i := 0; i < room.GetCfg().MaxCount; i++ {
					if HuData.HuCardCount[i] > 0 {
						for j := 0; j < HuData.HuCardCount[i]; j++ {
							HuData.HuCardRemainingCount[i][j] = room.GetRemainingCount(wTargetUser, HuData.HuCardData[i][j])
						}
					} else {
						break
					}
				}
				u := room.MjBase.UserMgr.GetUserByChairId(wTargetUser)
				u.WriteMsg(HuData)
			}
		}
		OperateResult.ActionMask |= room.UserAction[wTargetUser]
	}

	//发送消息
	room.MjBase.UserMgr.SendMsgAll(OperateResult)

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

//派发扑克
func (room *ZP_RoomData) DispatchCardData(wCurrentUser int, bTail bool) int {
	//状态效验
	if room.SendStatus == Not_Send {
		log.Error("at DispatchCardData f room.SendStatus == Not_Send")
		return -1
	}

	//丢弃扑克
	if (room.OutCardUser != INVALID_CHAIR) && (room.OutCardData != 0) {
		if len(room.DiscardCard[room.OutCardUser]) < 1 {
			room.DiscardCard[room.OutCardUser] = make([]int, 60)
		}

		room.DiscardCard[room.OutCardUser] = append(room.DiscardCard[room.OutCardUser], room.OutCardData)
	}

	//荒庄结束
	if room.LeftCardCount <= room.EndLeftCount {
		log.Debug("荒庄结束,room.LeftCardCount:%d,room.EndLeftCount:%d", room.LeftCardCount, room.EndLeftCount)
		room.ProvideUser = INVALID_CHAIR
		return 1
	}

	//发送扑克
	room.ProvideCard = room.GetSendCard(bTail, room.MjBase.UserMgr.GetMaxPlayerCnt())
	room.SendCardData = room.ProvideCard
	room.LastCatchCardUser = wCurrentUser
	//清除禁止胡牌的牌

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
	}

	//加牌
	room.CardIndex[wCurrentUser][room.MjBase.LogicMgr.SwitchToCardIndex(room.ProvideCard)]++
	//room.UserCatchCardCount[wCurrentUser]++;

	if !room.MjBase.UserMgr.IsTrustee(wCurrentUser) {
		//胡牌判断
		chr := 0
		room.CardIndex[wCurrentUser][room.MjBase.LogicMgr.SwitchToCardIndex(room.SendCardData)]--
		log.Debug("befer %v ", room.UserAction[wCurrentUser])
		huKind, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser],
			room.SendCardData, chr, room.GetCfg().MaxCount, false)
		room.UserAction[wCurrentUser] |= huKind
		log.Debug("afert %v ", room.UserAction[wCurrentUser])
		room.CardIndex[wCurrentUser][room.MjBase.LogicMgr.SwitchToCardIndex(room.SendCardData)]++

		//杠牌判断
		if (room.LeftCardCount > room.EndLeftCount) && !room.Ting[wCurrentUser] {
			GangCardResult := &mj_base.TagGangCardResult{}
			room.UserAction[wCurrentUser] |= room.MjBase.LogicMgr.AnalyseGangCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser], room.ProvideCard, GangCardResult)
		}
	}

	//听牌判断
	HuData := &mj_zp_msg.G2C_ZPMJ_HuData{OutCardData: make([]int, room.GetCfg().MaxCount), HuCardCount: make([]int, room.GetCfg().MaxCount), HuCardData: make([][]int, room.GetCfg().MaxCount), HuCardRemainingCount: make([][]int, room.GetCfg().MaxCount)}
	for i := 0; i < room.GetCfg().MaxCount; i++ {
		HuData.HuCardData[i] = make([]int, 28)
		HuData.HuCardRemainingCount[i] = make([]int, 28)
	}

	if room.Ting[wCurrentUser] == false {
		cbCount := room.MjBase.LogicMgr.AnalyseTingCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser], HuData.OutCardData, HuData.HuCardCount, HuData.HuCardData, room.GetCfg().MaxCount)
		room.TingCnt[wCurrentUser] = int(cbCount)
		HuData.OutCardCount = int(cbCount)
		if cbCount > 0 {
			room.UserAction[wCurrentUser] |= WIK_LISTEN

			for i := 0; i < room.GetCfg().MaxCount; i++ {
				if HuData.HuCardCount[i] > 0 {
					for j := 0; j < HuData.HuCardCount[i]; j++ {
						HuData.HuCardRemainingCount[i] = append(HuData.HuCardRemainingCount[i], room.GetRemainingCount(wCurrentUser, HuData.HuCardData[i][j]))
					}
				} else {
					break
				}
			}

			u.WriteMsg(HuData)
		}
	}

	log.Debug("User Action === %v , %d", room.UserAction, room.UserAction[wCurrentUser])
	//构造数据
	SendCard := &mj_zp_msg.G2C_ZPMJ_SendCard{}
	SendCard.SendCardUser = wCurrentUser
	SendCard.CurrentUser = wCurrentUser
	SendCard.Tail = bTail
	SendCard.ActionMask = room.UserAction[wCurrentUser]
	SendCard.CardData = room.ProvideCard
	//发送数据
	u.WriteMsg(SendCard)
	SendCard.CardData = 0
	room.MjBase.UserMgr.SendMsgAllNoSelf(u.Id, SendCard)

	room.UserActionDone = false
	if room.MjBase.UserMgr.IsTrustee(wCurrentUser) {
		room.UserActionDone = true
	} else {
		//超时定时
		room.OutCardTimer(u)
	}
	return 0
}

//解散接触
func (room *ZP_RoomData) DismissEnd() {
	//清理变量
	room.ClearAllTimer()

	//变量定义
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	GameConclude := &mj_zp_msg.G2C_ZPMJ_GameConclude{}
	GameConclude.ChiHuKind = make([]int, UserCnt)
	GameConclude.CardCount = make([]int, UserCnt)
	GameConclude.HandCardData = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		GameConclude.HandCardData[i] = make([]int, 28)
	}
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

	//发送信息
	room.MjBase.UserMgr.SendMsgAll(GameConclude)
}

//空闲状态
func (room *ZP_RoomData) SendStatusReady(u *user.User) {
	StatusFree := &msg.G2C_StatusFree{}
	StatusFree.CellScore = room.Source                                     //基础积分
	StatusFree.TimeOutCard = room.MjBase.TimerMgr.GetTimeOutCard()         //出牌时间
	StatusFree.TimeOperateCard = room.MjBase.TimerMgr.GetTimeOperateCard() //操作时间
	StatusFree.CreateTime = room.MjBase.TimerMgr.GetCreatrTime()           //开始时间
	for _, v := range room.HistoryScores {
		StatusFree.TurnScore = append(StatusFree.TurnScore, v.TurnScore)
		StatusFree.CollectScore = append(StatusFree.TurnScore, v.CollectScore)
	}
	StatusFree.PlayerCount = room.MjBase.TimerMgr.GetPlayCount() //玩家人数
	StatusFree.MaCount = 0                                       //码数
	StatusFree.CountLimit = room.MjBase.TimerMgr.GetMaxPayCnt()  //局数限制
	StatusFree.ZhuaHuaCnt = room.ZhuaHuaCnt
	u.WriteMsg(StatusFree)
}

///////////////////////////////////////////////////////////////////////////////////
//定时器

//超时出牌定时
func (room *ZP_RoomData) OutCardTimer(u *user.User) {
	//stop
	if room.OutCardTime != nil {
		room.OutCardTime.Stop()
	}

	room.OutCardTime = room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OutCardTime)*time.Second, func() {
		log.Debug("超时出牌定时")
		card := room.SendCardData
		for j := 0; j < room.GetCfg().MaxIdx; j++ {
			if room.CardIndex[u.ChairId][j] > 0 {
				card = room.MjBase.LogicMgr.SwitchToCardData(j)
				break
			}
		}
		log.Debug("用户%d超时打牌：%x", u.ChairId, card)
		room.MjBase.OutCard([]interface{}{u, card, true})
	})
}

//操作定时器
func (room *ZP_RoomData) OperateCardTimer(u *user.User) {
	chairID := u.ChairId

	//todo,是否stop OutCardTime
	if room.OutCardTime != nil {
		room.OutCardTime.Stop()
	}
	if room.OperateTime[chairID] != nil {
		room.OperateTime[chairID].Stop()
	}

	operateTimer := room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OperateCardTime)*time.Second, func() {
		log.Debug("操作定时器")
		if room.UserAction[chairID] != WIK_LISTEN {
			operateCard := []int{0, 0, 0}
			room.MjBase.UserOperateCard([]interface{}{u, WIK_NULL, operateCard})
		} else {
			room.OnUserListenCard(u, false)
		}
	})
	room.OperateTime = append(room.OperateTime, operateTimer)
}

//清理定时器
func (room *ZP_RoomData) ClearAllTimer() {
	if room.OutCardTime != nil {
		room.OutCardTime.Stop()
	}
	for k := range room.OperateTime {
		if room.OperateTime[k] != nil {
			room.OperateTime[k].Stop()
		}
	}
	log.Debug("zpmj at ClearAllTimer")
}
