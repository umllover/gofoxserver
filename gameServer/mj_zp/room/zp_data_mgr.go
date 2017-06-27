package room

import (
	"math"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/common"
	"mj/gameServer/common/mj_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"mj/common/msg/mj_zp_msg"

	"time"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/timer"
	"github.com/lovelly/leaf/util"
)

type ZP_RoomData struct {
	*mj_base.RoomData
	ChaHuaTime *timer.Timer

	FollowCard      []int //跟牌
	FollowCardScore []int //跟牌得分

	FlowerCnt map[int]int //补花数
	ChaHuaMap map[int]int //插花数
}

func NewDataMgr(id, uid, OriCardIdx int, name string, temp *base.GameServiceOption, base *mj_base.Mj_base) *ZP_RoomData {
	r := new(ZP_RoomData)
	r.ChaHuaMap = make(map[int]int)
	r.RoomData = mj_base.NewDataMgr(id, uid, OriCardIdx, name, temp, base)
	return r
}

func (room *ZP_RoomData) InitRoom(UserCnt int) {
	//初始化
	log.Debug("初始化漳浦房间")
	room.RepertoryCard = make([]int, MAX_REPERTORY)
	room.CardIndex = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.CardIndex[i] = make([]int, MAX_INDEX)
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

	room.LeftCardCount = MAX_REPERTORY
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
		room.ChaHuaTime = room.MjBase.AfterFunc(time.Duration(0)*time.Second, func() {
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
	}
}

func (room *ZP_RoomData) AfterStartGame() {

}

//获得插花
func (room *ZP_RoomData) GetChaHua(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	getData := args[0].(*mj_zp_msg.C2G_MJZP_SetChaHua)
	room.ChaHuaMap[user.ChairId] = getData.SetCount
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
func (room *ZP_RoomData) OnUserReplaceCard(args []interface{}) bool {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)
	gameLogic := room.MjBase.LogicMgr

	getData := args[0].(*mj_zp_msg.C2G_MJZP_ReplaceCard)

	if gameLogic.RemoveCard(room.CardIndex[user.ChairId], getData.CardData) == false {
		log.Debug("[用户补花] 用户：%d补花失败", user.ChairId)
		return false
	}

	//记录补花
	room.FlowerCnt[user.ChairId]++

	//是否花杠
	if room.FlowerCnt[user.ChairId] == 8 {
		room.MjBase.OnEventGameConclude(user.ChairId, user, GER_NORMAL)
	}

	//状态变量
	room.SendStatus = BuHua_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideUser = INVALID_CHAIR

	//派发扑克
	room.DispatchCardData(user.ChairId, true)

	outData := &mj_zp_msg.G2C_MJZP_ReplaceCard{}
	outData.IsInitFlower = false
	outData.ReplaceUser = user.ChairId
	outData.ReplaceCard = getData.CardData
	outData.NewCard = room.SendCardData
	room.MjBase.UserMgr.SendMsgAll(&outData)

	log.Debug("[用户补花] 用户：%d,花牌：%x 新牌：%x", user.ChairId, getData.CardData, room.SendCardData)
	return true
}

//用户听牌
func (room *ZP_RoomData) OnUserListenCard(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)
	//gameLogic := room.MjBase.LogicMgr

	getData := args[0].(*mj_zp_msg.C2G_MJZP_ListenCard)
	if getData.ListenCard { //todo,用户点击听
		//sendData := &mj_zp_msg.G2C_MJZP_ListenCard{}

		//if WIK_LISTEN == gameLogic.AnalyseTingCard(room.CardIndex[user.ChairId], room.WeaveItemArray[user.ChairId],
		//	, sendData.HuCardCount, sendData.HuCardData) {
		//
		//}
	} else {
		room.Ting[user.ChairId] = false
		sendData := &mj_zp_msg.G2C_MJZP_ListenCard{}
		sendData.ListenUser = user.ChairId
		sendData.IsListen = false
		room.MjBase.UserMgr.SendMsgAll(&sendData)
	}
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

			outData := &mj_zp_msg.G2C_MJZP_ReplaceCard{}
			outData.ReplaceUser = playerIndex
			outData.IsInitFlower = true
			for j := MAX_INDEX - MAX_HUA_INDEX; j < MAX_INDEX; j++ {
				if room.CardIndex[playerIndex][j] == 1 {
					for {
						outData.NewCard = room.GetSendCard(true, playerCNT)
						newCardIndex := SwitchToCardIndex(outData.NewCard)
						outData.ReplaceCard = SwitchToCardIndex(j)
						room.MjBase.UserMgr.SendMsgAll(&outData)

						room.FlowerCnt[playerIndex]++
						room.CardIndex[playerIndex][newCardIndex]++
						if newCardIndex >= (MAX_INDEX-MAX_HUA_INDEX) && newCardIndex <= MAX_INDEX {
							break
						}
					}
				}
			}
		}
	}
	playerIndex++
}

//庄家开局动作
func (room *ZP_RoomData) InitBankerAction() {
	userMgr := room.MjBase.UserMgr
	UserCnt := userMgr.GetMaxPlayerCnt()
	gameLogic := room.MjBase.LogicMgr
	room.UserAction = make([]int, UserCnt)

	gangCardResult := &common.TagGangCardResult{}
	room.UserAction[room.BankerUser] |= gameLogic.AnalyseGangCard(room.CardIndex[room.BankerUser], nil, 0, gangCardResult)

	//胡牌判断
	chr := 0
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]--
	room.UserAction[room.BankerUser] |= gameLogic.AnalyseChiHuCard(room.CardIndex[room.BankerUser], []*msg.WeaveItem{}, room.SendCardData, chr, true)
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

	gameLogic.RandCardList(room.RepertoryCard, mj_base.GetZpmjCards())

	//todo,剔除大字

	//分发扑克
	userMgr.ForEachUser(func(u *user.User) {
		for i := 0; i < MAX_COUNT-1; i++ {
			room.LeftCardCount -= 1
			room.MinusHeadCount += 1
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
	TakeCount := MAX_REPERTORY - room.LeftCardCount
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
					if room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], cbCenterCard, chr, false) == WIK_CHI_HU {
						room.UserAction[u.ChairId] |= WIK_CHI_HU
					}
				}
			}
		}

		//检查抢杠胡
		if EstimatKind == EstimatKind_GangCard {
			//只有庄家和闲家之间才能放炮
			MogicCard := room.MjBase.LogicMgr.SwitchToCardData(room.MjBase.LogicMgr.GetMagicIndex())
			if room.MjBase.LogicMgr.GetMagicIndex() == MAX_INDEX || (room.MjBase.LogicMgr.GetMagicIndex() != MAX_INDEX && cbCenterCard != MogicCard) {
				if u.UserLimit|LimitChiHu == 0 {
					//吃胡判断
					chr := 0
					room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], cbCenterCard, chr, false)
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
	//todo,
}

//进行抓花
func (room *ZP_RoomData) OnZhuaHua(CenterUser int) {
	//todo,进行抓花

	//room.RepertoryCard[room.l]
}

//记录分饼
func (room *ZP_RoomData) RecordFollowCard(cbCenterCard int) bool {
	room.FollowCard = append(room.FollowCard, cbCenterCard)

	count := len(room.FollowCard) % 4
	if count == 0 {
		begin := count - 4
		for i := begin; i < count; i++ {
			if room.FollowCard[i] != cbCenterCard {
				return false
			}
		}
	}

	userCNT := room.MjBase.UserMgr.GetMaxPlayerCnt()
	for i := 0; i < userCNT; i++ {
		if i == room.BankerUser {
			room.FollowCardScore[room.BankerUser] -= 3
			continue
		} else {
			room.FollowCardScore[i] += 1
		}
	}

	return true
}
