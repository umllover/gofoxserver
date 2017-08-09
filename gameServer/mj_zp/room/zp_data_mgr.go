package room

import (
	"math"
	. "mj/common/cost"
	"mj/common/msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/conf"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"mj/common/msg/mj_zp_msg"

	"mj/common/utils"

	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/timer"
	"github.com/lovelly/leaf/util"
)

type ZP_RoomData struct {
	*mj_base.RoomData
	ChaHuaTime  *timer.Timer //插花定时器
	OutCardTime *timer.Timer //出牌定时器

	ZhuaHuaCnt int  //抓花个数
	WithZiCard bool //带字牌
	WithChaHua bool //是否插花
	ScoreType  int  //算分制式

	FollowCard   []int       //跟牌
	IsFollowCard bool        //是否跟牌
	LianZhuang   int         //连庄次数
	ChaHuaMap    map[int]int //插花数
	ChaHuaScore  []int       //插花分
	HuKindType   []int       //胡牌类型
	TingCnt      [4]int      //听牌个数

	ZhuaHuaMap      [16]*mj_zp_msg.HuaUser   //插花数据
	HuKindScore     [4][COUNT_KIND_SCORE]int //特殊胡牌分
	ZhuaHuaScore    [4]int                   //插花得分
	FollowCardScore []int                    //跟牌得分
	SumScore        [4]int                   //游戏总分
}

func NewZPDataMgr(info *msg.L2G_CreatorRoom, uid int64, configIdx int, name string, temp *base.GameServiceOption, base *ZP_base) *ZP_RoomData {
	r := new(ZP_RoomData)
	r.ChaHuaMap = make(map[int]int)
	r.RoomData = mj_base.NewDataMgr(info.RoomID, uid, configIdx, name, temp, base.Mj_base, info)

	r.IniSource = temp.IniScore
	getData, ok := r.OtherInfo["zhuaHua"].(float64)
	if !ok {
		log.Error("zpmj at NewDataMgr [zhuaHua] error")
		return nil
	}
	r.ZhuaHuaCnt = int(getData)

	getData2, ok := r.OtherInfo["wanFa"].(bool)
	if !ok {
		log.Error("zpmj at NewDataMgr [wanFa] error")
		return nil
	}
	r.WithZiCard = getData2

	getData3, ok := r.OtherInfo["suanFen"].(float64)
	if !ok {
		log.Error("zpmj at NewDataMgr [suanFen] error")
		return nil
	}
	r.ScoreType = int(getData3)

	getData4, ok := r.OtherInfo["chaHua"].(bool)
	if !ok {
		log.Error("zpmj at NewDataMgr [chaHua] error")
		return nil
	}
	r.WithChaHua = getData4
	return r
}

func (room *ZP_RoomData) InitRoom(UserCnt int) {
	//初始化
	log.Debug("zpmj at InitRoom version 000002")
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
	room.BanUser = make([]int, UserCnt)
	room.BanCardCnt = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.HeapCardInfo[i] = make([]int, 2)
		room.BanCardCnt[i] = make([]int, 9)
	}
	room.OperateTime = make([]*timer.Timer, UserCnt)

	room.UserActionDone = false
	room.SendStatus = Not_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideGangUser = INVALID_CHAIR
	room.MinusLastCount = 0
	room.MinusHeadCount = room.GetCfg().MaxRepertory
	room.OutCardCount = 0

	//设置漳浦麻将牌数据
	room.EndLeftCount = 16
	room.IsFollowCard = false
	room.TingCnt = [4]int{}
	room.FollowCard = room.FollowCard[0:0]
	room.HuKindType = room.HuKindType[0:0]
	room.HuKindType = append(room.HuKindType, 1)
	room.FollowCardScore = make([]int, UserCnt)
	room.SumScore = [4]int{}
	room.HuKindScore = [4][COUNT_KIND_SCORE]int{}
	room.ZhuaHuaMap = [16]*mj_zp_msg.HuaUser{}
	room.ZhuaHuaScore = [4]int{}

	room.FlowerCnt = make([]int, UserCnt)
	room.IsResponse = make([]bool, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
	room.FlowerCard = make([][]int, UserCnt)
	room.ChaHuaScore = make([]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.OperateCard[i] = make([]int, 60)
		room.MjBase.UserMgr.SetUsetTrustee(i, false)
	}
	log.Debug("len1 OperateCard: %d %d", len(room.OperateCard), len(room.OperateCard[1]))
	room.PerformAction = make([]int, UserCnt)
}

func (room *ZP_RoomData) BeforeStartGame(UserCnt int) {
	room.InitRoom(UserCnt)
}

func (room *ZP_RoomData) StartGameing() {
	log.Debug("开始漳浦游戏")
	if room.MjBase.TimerMgr.GetPlayCount() == 0 && room.WithChaHua == true {
		room.MjBase.UserMgr.SendMsgAll(&mj_zp_msg.G2C_MJZP_NotifiChaHua{})

		room.ChaHuaTime = room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OperateCardTime)*time.Second, func() {
			log.Debug("超时插花")
			for i := 0; i < 4; i++ {
				_, ok := room.ChaHuaMap[i]
				if !ok {
					room.ChaHuaMap[i] = 0
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
			room.InitOutCardTimer(u)
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
		room.InitOutCardTimer(u)
	}
}

func (room *ZP_RoomData) AfterStartGame() {

}

//获得插花
func (room *ZP_RoomData) GetChaHua(u *user.User, setCount int) {
	log.Debug("获得插花")
	_, ok := room.ChaHuaMap[u.ChairId]
	if ok {
		log.Debug("%d 已经插花", u.ChairId)
		return
	}
	room.ChaHuaMap[u.ChairId] = setCount

	sendData := &mj_zp_msg.G2C_MJZP_UserCharHua{}
	sendData.SetCount = setCount
	sendData.Chair = u.ChairId
	room.MjBase.UserMgr.SendMsgAll(sendData)
	userCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	if len(room.ChaHuaMap) == userCnt && room.MjBase.TimerMgr.GetPlayCount() == 0 && room.WithChaHua == true {
		if room.ChaHuaTime != nil {
			room.ChaHuaTime.Stop()
		}
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
		room.InitOutCardTimer(u)
	}
}

func (room *ZP_RoomData) SendReplaceCard(ReplaceUser, ReplaceCard, NewCard int, IsInitFlower bool) {
	room.MjBase.UserMgr.SendMsgAll(&mj_zp_msg.G2C_MJZP_ReplaceCard{
		ReplaceUser:  ReplaceUser,
		ReplaceCard:  ReplaceCard,
		NewCard:      NewCard,
		IsInitFlower: IsInitFlower,
		ActionMask:   room.UserAction[ReplaceUser],
	})
}

//用户听牌
func (room *ZP_RoomData) OnUserListenCard(u *user.User, bListenCard bool) bool {
	gameLogic := room.MjBase.LogicMgr

	if bListenCard {
		if WIK_LISTEN == gameLogic.AnalyseTingCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], nil, nil, nil) {
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
	bufferCount := 0
	for _, v := range OriDataArray {
		if v >= 0x31 && v <= 0x37 {
			continue
		}
		NewDataArray[bufferCount] = v
		bufferCount++
	}
}

//庄家开局动作
func (room *ZP_RoomData) InitBankerAction() {
	userMgr := room.MjBase.UserMgr
	UserCnt := userMgr.GetMaxPlayerCnt()
	gameLogic := room.MjBase.LogicMgr
	room.UserAction = make([]int, UserCnt)

	gangCardResult := &mj_base.TagGangCardResult{}
	room.UserAction[room.BankerUser] |= gameLogic.AnalyseGangCard(room.CardIndex[room.BankerUser], nil, 0, gangCardResult)

	//胡牌判断
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]--
	huKind, _ := gameLogic.AnalyseChiHuCard(room.CardIndex[room.BankerUser], []*msg.WeaveItem{}, room.SendCardData)
	if huKind {
		room.UserAction[room.BankerUser] |= WIK_CHI_HU
	}
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]++

	if room.UserAction[room.BankerUser] != 0 {
		outData := &mj_zp_msg.G2C_MJZP_OperateNotify{}
		outData.ActionCard = room.SendCardData
		outData.ActionMask = room.UserAction[room.BankerUser]
		u := userMgr.GetUserByChairId(room.BankerUser)
		u.WriteMsg(outData)
		//定时
		room.OperateCardTimer(u)
	}
}

//发牌
func (room *ZP_RoomData) StartDispatchCard() {
	log.Debug("开始发牌")
	userMgr := room.MjBase.UserMgr
	gameLogic := room.MjBase.LogicMgr

	//初始化变量
	gameLogic.SetMagicIndex(room.GetCfg().MaxIdx)

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, US_PLAYING)
	})

	var minSice int
	UserCnt := userMgr.GetMaxPlayerCnt()
	room.SiceCount, minSice = room.GetSice()

	gameLogic.RandCardList(room.RepertoryCard, mj_base.GetCardByIdx(room.ConfigIdx))

	//剔除大字
	log.Debug("剔除大字before:%v", room.RepertoryCard)
	if room.WithZiCard == false {
		tempCard := make([]int, room.GetCfg().MaxRepertory-7*4)
		room.RemoveAllZiCar(tempCard, room.RepertoryCard)
		room.RepertoryCard = tempCard
		log.Debug("剔除大字1:%v", room.RepertoryCard)
		room.MinusHeadCount = len(room.RepertoryCard)
	}

	log.Debug("======房间Id：%d", room.ID)

	//选取庄家
	if room.BankerUser == INVALID_CHAIR {
		_, room.BankerUser = room.MjBase.UserMgr.GetUserByUid(room.CreatorUid)
	}

	//分发扑克
	userMgr.ForEachUser(func(u *user.User) {
		for i := 0; i < room.GetCfg().MaxCount-1; i++ {
			setIndex := SwitchToCardIndex(room.GetHeadCard())
			room.CardIndex[u.ChairId][setIndex]++
		}
		log.Debug("用户%d手牌：%v", u.ChairId, room.CardIndex[u.ChairId])
	})

	room.SendCardData = room.GetHeadCard()
	room.CardIndex[room.BankerUser][SwitchToCardIndex(room.SendCardData)]++
	room.ProvideCard = room.SendCardData
	room.ProvideUser = room.BankerUser
	room.CurrentUser = room.BankerUser

	if conf.Test {
		room.RepalceCard()
	}

	//todo,测试手牌
	var temp []int
	temp = make([]int, 42)
	temp[0] = 3 //三张一同
	temp[1] = 3 //三张二同
	temp[2] = 3 //三张三同
	temp[3] = 3 //三张四同
	temp[4] = 3 //三张五同
	temp[5] = 2

	//room.FlowerCnt[0] = 1 //花牌
	room.SendCardData = 0x04
	room.CardIndex[0] = temp
	GetCardWordArray(room.CardIndex[0])
	log.Debug("@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	log.Debug("room.CardIndex:%v", room.CardIndex[0])
	//
	//for k := range room.RepertoryCard {
	//	room.RepertoryCard[k] = 0x01
	//}

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
		if wCenterUser == u.ChairId {
			return
		}

		//托管不响应
		if room.MjBase.UserMgr.IsTrustee(u.ChairId) {
			log.Debug("at EstimateUserRespond ======== IsTrustee ChairId:%v", u.ChairId)
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
			eatUser := (wCenterUser + UserCnt + 1) % UserCnt //4==GAME_PLAYER
			if eatUser == u.ChairId {
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimateEatCard(room.CardIndex[u.ChairId], cbCenterCard)
				log.Debug("吃牌用户：%d 动作：%d,wCenterUser:%d", u.ChairId, room.UserAction[u.ChairId], wCenterUser)
				room.BanCardCnt[u.ChairId][LimitChi] = cbCenterCard
			}

			//杠牌判断
			if room.IsEnoughCard() && u.UserLimit&LimitGang == 0 {
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimateGangCard(room.CardIndex[u.ChairId], cbCenterCard)
			}

			//吃胡判断
			if u.ChairId != wCenterUser {
				if u.UserLimit&LimitChiHu == 0 {
					//有禁止吃胡的牌
					if !(room.BanUser[u.ChairId]&LimitChiHu != 0 && room.BanCardCnt[u.ChairId][LimitChiHu] == cbCenterCard) {
						log.Debug("有吃胡2")
						hu, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], cbCenterCard)
						if hu {
							log.Debug("有吃胡3")
							room.UserAction[u.ChairId] |= WIK_CHI_HU
							room.BanCardCnt[u.ChairId][LimitChiHu] = cbCenterCard
						}
					}
				}
			}
		}

		//检查抢杠胡
		if EstimatKind == EstimatKind_GangCard {
			//只有庄家和闲家之间才能放炮
			if u.UserLimit&LimitChiHu == 0 {
				//吃胡判断
				hu, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], cbCenterCard)
				if hu {
					room.UserAction[u.ChairId] |= WIK_CHI_HU
					//抢杠胡特殊分
					room.HuKindScore[u.ChairId][IDX_SUB_SCORE_QGH] = 3
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
				u.WriteMsg(&mj_zp_msg.G2C_MJZP_OperateNotify{
					ActionMask: room.UserAction[u.ChairId],
					ActionCard: room.ProvideCard,
				})
				//定时
				if room.MjBase.UserMgr.IsTrustee(u.ChairId) {
					u := room.MjBase.UserMgr.GetUserByChairId(u.ChairId)
					operateCard := []int{0, 0, 0}
					room.MjBase.UserOperateCard([]interface{}{u, WIK_NULL, operateCard})
				} else {
					room.OperateCardTimer(u)
				}
			}
		})
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

//向客户端发牌
func (room *ZP_RoomData) SendGameStart() {
	//构造变量
	GameStart := &mj_zp_msg.G2C_ZPMG_GameStart{}
	GameStart.BankerUser = room.BankerUser
	GameStart.SiceCount = room.SiceCount
	GameStart.HeapHead = room.HeapHead
	GameStart.HeapTail = room.HeapTail
	GameStart.HeapCardInfo = room.HeapCardInfo
	//发送数据
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		GameStart.UserAction = room.UserAction[u.ChairId]
		GameStart.CardData = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[u.ChairId])
		u.WriteMsg(GameStart)
	})
}

//正常结束房间
func (room *ZP_RoomData) NormalEnd(cbReason int) {
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
	GameConclude.ChaHuaScore = make([]int, UserCnt)

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
	room.CalHuPaiScore(UserGameScore)

	//拷贝码数据
	GameConclude.MaCount = make([]int, 0)
	util.DeepCopy(&GameConclude.ZhuaHua, &room.ZhuaHuaMap) //抓花数据
	util.DeepCopy(&GameConclude.ChaHuaScore, &room.ChaHuaScore)

	//积分变量
	ScoreInfoArray := make([]*msg.TagScoreInfo, UserCnt)

	GameConclude.ProvideUser = room.ProvideUser
	GameConclude.ProvideCard = room.ProvideCard

	//统计积分
	GameConclude.CellScore = room.Source
	DetailScore := make([]int, room.MjBase.UserMgr.GetMaxPlayerCnt())
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		if u.Status != US_PLAYING {
			return
		}
		util.DeepCopy(&GameConclude.ScoreKind[u.ChairId], &room.HuKindScore[u.ChairId]) //游戏得分类型
		GameConclude.GameScore[u.ChairId] = UserGameScore[u.ChairId]
		//胡牌分
		GameConclude.GameScore[u.ChairId] = room.SumScore[u.ChairId]

		//收税
		if GameConclude.GameScore[u.ChairId] > 0 && room.MjBase.Temp.GameType == GAME_GENRE_ZhuanShi {
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
		room.HistorySe.AllScore[u.ChairId] += GameConclude.GameScore[u.ChairId]
		DetailScore[u.ChairId] = GameConclude.GameScore[u.ChairId]
	})

	room.HistorySe.DetailScore = append(room.HistorySe.DetailScore, DetailScore)
	GameConclude.Reason = cbReason

	GameConclude.AllScore = room.HistorySe.AllScore
	GameConclude.DetailScore = room.HistorySe.DetailScore
	//发送数据
	room.MjBase.UserMgr.SendMsgAll(GameConclude)

	//写入积分 todo
	//room.MjBase.UserMgr.WriteTableScore(ScoreInfoArray, room.MjBase.UserMgr.GetMaxPlayerCnt(), ZPMJ_CHANGE_SOURCE)
}

//进行抓花
func (room *ZP_RoomData) OnZhuaHua(winUser []int) (CardData [][]int, BuZhong []int) {

	count := room.ZhuaHuaCnt
	if count == 0 {
		log.Debug("抓花0")
		return
	}

	CardData = make([][]int, len(winUser))
	for i := 0; i < count; i++ {
		cardData := room.GetHeadCard()

		var isZhong bool
		for k, userIndex := range winUser {
			//抓花规则
			var getInedx [3]int
			userCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
			index := [4][3]int{{1, 5, 9}, {0, 2, 6}, {0, 3, 7}, {0, 4, 8}}
			if room.BankerUser == userIndex {
				getInedx = index[0]
			} else {
				v := int(math.Abs(float64(userIndex-room.BankerUser+userCnt))) % userCnt
				getInedx = index[int(v)]
			}

			cardColor := room.MjBase.LogicMgr.GetCardColor(cardData)
			cardValue := room.MjBase.LogicMgr.GetCardValue(cardData)
			if cardColor == 0x30 {
				//东南西北
				if cardValue < 5 {
					if cardValue == getInedx[0] || cardValue == getInedx[1] || cardValue == getInedx[2] {
						isZhong = true
						CardData[k] = append(CardData[k], cardData)
						break
					}
				} else {
					//中发白
					temp := cardValue - 4
					if temp == getInedx[0] || temp == getInedx[1] || temp == getInedx[2] {
						isZhong = true
						CardData[k] = append(CardData[k], cardData)
						break
					}
				}
			} else if cardColor >= 0x00 && cardColor <= 0x20 {
				if cardValue == getInedx[0] || cardValue == getInedx[1] || cardValue == getInedx[2] {
					isZhong = true
					CardData[k] = append(CardData[k], cardData)
					break
				}
			}
		}
		if !isZhong {
			BuZhong = append(BuZhong, cardData)
		}
	}

	log.Debug("================================== zhonghua:%v buzhong:%v", CardData, BuZhong)
	return
}

//记录分饼
func (room *ZP_RoomData) RecordFollowCard(wTargetUser, cbCenterCard int) bool {
	if room.IsFollowCard {
		return false
	}

	if len(room.WeaveItemArray[wTargetUser]) > 0 {
		room.IsFollowCard = true //取消跟牌
		return false
	}

	log.Debug("记录分饼")
	room.FollowCard = append(room.FollowCard, cbCenterCard)

	count := len(room.FollowCard) % 4
	if count == 0 {
		begin := 0
		if len(room.FollowCard) > 8 {
			begin = count - 4
		}
		for i := begin; i < len(room.FollowCard); i++ {
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
	room.OperateCard[u.ChairId] = make([]int, 4)
	if len(OperateCard) < 2 {
		room.BuildOpCard(u.ChairId, OperateCode, OperateCard[0])
	} else {
		for i, card := range OperateCard {
			room.OperateCard[u.ChairId][i] = card
		}
	}
	room.StopOperateCardTimer(u) //清理定时

	u.UserLimit = 0
	//放弃操作
	if OperateCode == WIK_NULL {
		log.Debug("放弃操作")
		//////禁止这轮吃胡
		//if room.HasOperator(u.ChairId, WIK_CHI_HU) {
		//	u.UserLimit |= LimitChiHu
		//}
		//抢杠胡分
		room.HuKindScore[u.ChairId][IDX_SUB_SCORE_QGH] = 0
		//记录放弃操作
		room.RecordBanCard(OperateCode, u.ChairId)
		room.StopOperateCardTimer(u)
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
		//优先级别
		cbUserActionRank := room.MjBase.LogicMgr.GetUserActionRank(cbUserAction)
		cbTargetActionRank := room.MjBase.LogicMgr.GetUserActionRank(cbTargetAction)

		//动作判断
		if cbUserActionRank > cbTargetActionRank {
			wTargetUser = i
			cbTargetAction = cbUserAction
		}
	}
	for i := 0; i < userCnt; i++ {
		if i == wTargetUser {
			continue
		}
		clearUser := room.MjBase.UserMgr.GetUserByChairId(i)
		room.StopOperateCardTimer(clearUser)
	}

	if room.IsResponse[wTargetUser] == false { //最高权限的人没响应
		return -1, u.ChairId
	}

	//吃胡等待
	if cbTargetAction == WIK_CHI_HU {
		for i := 0; i < userCnt; i++ {
			if room.IsResponse[i] == false && room.UserAction[i]&WIK_CHI_HU != 0 {
				return -1, u.ChairId
			}
		}
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
	kind, AnalyseItem := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], pWeaveItem, room.SendCardData)
	if kind {
		room.ChiHuKind[u.ChairId] = WIK_CHI_HU
	}

	if room.FlowerCnt[u.ChairId] == 8 {
		room.ChiHuKind[u.ChairId] = WIK_CHI_HU
	}
	room.ProvideCard = room.SendCardData

	//特殊胡牌类型
	room.CurrentUser = u.ChairId
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
		chihuKind, AnalyseItem := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[wChiHuUser], pWeaveItem, room.OperateCard[wTargetUser][0])
		if chihuKind {
			room.ChiHuKind[wChiHuUser] = WIK_CHI_HU
		}

		//特殊胡牌类型
		room.CurrentUser = wChiHuUser
		room.SpecialCardKind(AnalyseItem, wChiHuUser)

		//插入扑克
		if room.ChiHuKind[wChiHuUser] != WIK_NULL {
			wTargetUser = wChiHuUser
		}
	}
}

//特殊胡牌类型及算分
func (room *ZP_RoomData) SpecialCardKind(TagAnalyseItem []*TagAnalyseItem, HuUserID int) {

	type1Cnt := 0
	type2Cnt := 0
	score := room.Source
	winScore := &room.HuKindScore[HuUserID]
	for _, v := range TagAnalyseItem {
		kind := 0
		kind = room.IsDaSanYuan(v) //大三元
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DSY] = 12 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("大三元 %d", winScore[IDX_SUB_SCORE_DSY])
		}
		kind = room.IsXiaoSanYuan(v) //小三元
		if kind > 0 {
			winScore[IDX_SUB_SCORE_XSY] = 6 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("小三元 %d", winScore[IDX_SUB_SCORE_XSY])
		}
		kind = room.IsHunYiSe(v) //混一色
		if kind > 0 {
			winScore[IDX_SUB_SCORE_CYS] = 6 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("混一色 %d", winScore[IDX_SUB_SCORE_CYS])
		}
		kind = room.IsQingYiSe(v, room.FlowerCnt) //清一色
		if kind > 0 {
			winScore[IDX_SUB_SCORE_QYS] = 24 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("清一色 %d", winScore[IDX_SUB_SCORE_QYS])
		}
		kind = room.IsHuaYiSe(v, room.FlowerCnt) //花一色
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HYS] = 12 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("花一色 %d", winScore[IDX_SUB_SCORE_HYS])
		}
		kind = room.IsGangKaiHua(v, room.WeaveItemArray[HuUserID]) //杠上开花
		if kind > 0 {
			winScore[IDX_SUB_SCORE_GSKH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("杠上开花 %d", winScore[IDX_SUB_SCORE_GSKH])
		}
		kind = room.IsHuaKaiHua(v) //花上开花
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HSKH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("花上开花 %d", winScore[IDX_SUB_SCORE_HSKH])
		}
		kind = room.IsBaiLiu(v, room.FlowerCnt) //佰六
		if kind > 0 {
			winScore[IDX_SUB_SCORE_BL] = 6 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("佰六 %d", winScore[IDX_SUB_SCORE_BL])
		}
		kind = room.IsMenQing(v) //门清
		if kind > 0 {
			winScore[IDX_SUB_SCORE_MQQ] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("门清 %d", winScore[IDX_SUB_SCORE_MQQ])
		}
		kind = room.IsMenQingBaiLiu(v, room.FlowerCnt) //门清佰六
		if kind > 0 {
			winScore[IDX_SUB_SCORE_BL] = 0
			winScore[IDX_SUB_SCORE_MQQ] = 0
			winScore[IDX_SUB_SCORE_MQBL] = 9 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("门清佰六 %d", winScore[IDX_SUB_SCORE_MQBL])
		}
		kind = room.IsHuWeiZhang(v) //尾单吊
		if kind > 0 {
			winScore[IDX_SUB_SCORE_WDD] = 6 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("尾单吊 %d", winScore[IDX_SUB_SCORE_WDD])
		}
		kind = room.IsJieTou(v) //截头
		if kind > 0 {
			winScore[IDX_SUB_SCORE_JT] = 1 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("截头 %d", winScore[IDX_SUB_SCORE_JT])
		}
		kind = room.IsKongXin(v) //空心
		if kind > 0 {
			winScore[IDX_SUB_SCORE_KX] = 1 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("空心 %d", winScore[IDX_SUB_SCORE_KX])
		}
		kind = room.IsDuiDuiHu(v) //对对胡
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DDH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("对对胡 %d", winScore[IDX_SUB_SCORE_DDH])
		}
		kind = room.IsTianHu(v) //天胡
		if kind > 0 {
			winScore[IDX_SUB_SCORE_TH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("天胡 %d", winScore[IDX_SUB_SCORE_TH])
		}
		kind = room.IsDiHu(v) //地胡
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("地胡 %d", winScore[IDX_SUB_SCORE_DH])
		}
		kind = room.IsHaiDiLaoYue(v) //海底捞针
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HDLZ] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("海底捞针 %d", winScore[IDX_SUB_SCORE_HDLZ])
		}
		kind = room.IsAnKe(v) //暗刻
		if kind > 0 {
			winScore[IDX_SUB_SCORE_SANAK+kind/8] = 3 * (kind / 4) * score //2,8,16
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("%d暗刻(32,33,34) %d", IDX_SUB_SCORE_SANAK+kind/8, winScore[IDX_SUB_SCORE_SANAK+kind/8])
		}
		kind = room.IsDaSiXi(v) //大四喜
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DSX] = 24 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("大四喜 %d", winScore[IDX_SUB_SCORE_DSX])
		}
		kind = room.IsXiaoSiXi(v) //小四喜
		if kind > 0 {
			winScore[IDX_SUB_SCORE_XSX] = 12 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("小四喜 %d", winScore[IDX_SUB_SCORE_XSX])
		}
		//自摸
		kind = room.IsZiMo()
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_HDLZ] > 0 || winScore[IDX_SUB_SCORE_GSKH] > 0 || winScore[IDX_SUB_SCORE_HSKH] > 0 ||
				winScore[IDX_SUB_SCORE_TH] > 0 || winScore[IDX_SUB_SCORE_TH] > 0 {
				continue
			}
			winScore[IDX_SUB_SCORE_ZM] = 1 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("自摸,%d", winScore[IDX_SUB_SCORE_ZM])
		}
		//无花字
		kind = room.IsWuHuaZi(v, room.FlowerCnt)
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_BL] > 0 || winScore[IDX_SUB_SCORE_MQBL] > 0 || winScore[IDX_SUB_SCORE_QYS] > 0 {
				continue
			}
			winScore[IDX_SUB_SCORE_WHZ] = 3 * score
			log.Debug("无花字，%d", winScore[IDX_SUB_SCORE_WHZ])
		}
		//字一色
		kind = room.IsZiYiSe(v, room.FlowerCnt)
		if kind > 0 {
			winScore[IDX_SUB_SCORE_ZYS] = 12 * score
			log.Debug("字一色，%d", winScore[IDX_SUB_SCORE_ZYS])
		}
		kind, type1Cnt, type2Cnt = room.IsZiPaiGang(v) //字牌杠
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_DSX] > 0 || winScore[IDX_SUB_SCORE_XSX] > 0 {
				winScore[IDX_SUB_SCORE_ZPG] = type1Cnt * score
			} else if winScore[IDX_SUB_SCORE_DSY] > 0 || winScore[IDX_SUB_SCORE_XSY] > 0 {
				winScore[IDX_SUB_SCORE_ZPG] = type2Cnt * score
			} else {
				winScore[IDX_SUB_SCORE_ZPG] = (type2Cnt + type1Cnt) * score
			}
			type1Cnt = 0
			type2Cnt = 0
			log.Debug("字牌杠，%d", winScore[IDX_SUB_SCORE_ZPG])
		}
		kind, type1Cnt, type2Cnt = room.IsKeZi(v) //字牌刻字
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_DSX] > 0 || winScore[IDX_SUB_SCORE_XSX] > 0 {
				winScore[IDX_SUB_SCORE_ZPKZ] = type1Cnt * score
			} else if winScore[IDX_SUB_SCORE_DSY] > 0 || winScore[IDX_SUB_SCORE_XSY] > 0 {
				winScore[IDX_SUB_SCORE_ZPKZ] = type2Cnt * score
			} else {
				winScore[IDX_SUB_SCORE_ZPKZ] = (type2Cnt + type1Cnt) * score
			}
			type1Cnt = 0
			type2Cnt = 0
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("字牌刻字 %d", winScore[IDX_SUB_SCORE_ZPKZ])
		}
		kind = room.IsDanDiao(v) //单吊
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_KX] > 0 || winScore[IDX_SUB_SCORE_JT] > 0 || winScore[IDX_SUB_SCORE_WDD] > 0 {
				continue
			}
			if room.CurrentUser == room.ProvideUser {
				winScore[IDX_SUB_SCORE_DDZM] = 1 * score
				room.HuKindType = append(room.HuKindType, IDX_SUB_SCORE_DDZM)
				log.Debug("单吊自摸,%d", winScore[IDX_SUB_SCORE_DDZM])
			} else {
				winScore[IDX_SUB_SCORE_DDPH] = 1 * score
				room.HuKindType = append(room.HuKindType, IDX_SUB_SCORE_DDPH)
				log.Debug("单吊平胡,%d", winScore[IDX_SUB_SCORE_DDPH])
			}
		}
	}
}

//特殊胡牌算分规则
func (room *ZP_RoomData) SpecialCardScore(HuUserID int) {
	score := room.Source
	winScore := &room.HuKindScore[HuUserID]
	if room.ScoreType == GAME_TYPE_33 {
		winScore[IDX_SUB_SCORE_DDPH] = 0
		return
	}

	if room.ScoreType == GAME_TYPE_48 {
		for k, v := range winScore {
			if v <= 0 {
				continue
			}

			switch k {
			case IDX_SUB_SCORE_HDLZ:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_GSKH:
				winScore[k] = 4 * score
			case IDX_SUB_SCORE_HSKH:
				winScore[k] = 4 * score
			case IDX_SUB_SCORE_QYS:
				winScore[k] = 32 * score
			case IDX_SUB_SCORE_HYS:
				winScore[k] = 16 * score
			case IDX_SUB_SCORE_CYS:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_DSY:
				winScore[k] = 16 * score
			case IDX_SUB_SCORE_XSY:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_DDH:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_MQQ:
				winScore[k] = 4 * score
			case IDX_SUB_SCORE_BL:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_DH:
				winScore[k] = 4 * score
			case IDX_SUB_SCORE_TH:
				winScore[k] = 4 * score
			case IDX_SUB_SCORE_DDPH:
				winScore[k] = 1 * score
			case IDX_SUB_SCORE_WDD:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_MQBL:
				winScore[k] = 12 * score
			case IDX_SUB_SCORE_SANAK:
				winScore[k] = 4 * score
			case IDX_SUB_SCORE_SIAK:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_WUAK:
				winScore[k] = 16 * score
			case IDX_SUB_SCORE_ZM:
				winScore[k] = 1 * score
			case IDX_SUB_SCORE_QGH:
				winScore[k] = 4 * score
			case IDX_SUB_SCORE_WHZ:
				winScore[k] = 4 * score
			case IDX_SUB_SCORE_ZYS:
				winScore[k] = 16 * score
			}
		}

	} else if room.ScoreType == GAME_TYPE_88 {
		for k, v := range winScore {
			if v <= 0 {
				continue
			}

			switch k {
			case IDX_SUB_SCORE_HDLZ:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_GSKH:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_HSKH:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_QYS:
				winScore[k] = 32 * score
			case IDX_SUB_SCORE_HYS:
				winScore[k] = 16 * score
			case IDX_SUB_SCORE_CYS:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_DSY:
				winScore[k] = 16 * score
			case IDX_SUB_SCORE_XSY:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_DDH:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_MQQ:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_BL:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_DH:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_TH:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_DDPH:
				winScore[k] = 0
			case IDX_SUB_SCORE_WDD:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_MQBL:
				winScore[k] = 0
			case IDX_SUB_SCORE_SANAK:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_SIAK:
				winScore[k] = 16 * score
			case IDX_SUB_SCORE_WUAK:
				winScore[k] = 32 * score
			case IDX_SUB_SCORE_JT:
				winScore[k] = 0
			case IDX_SUB_SCORE_KX:
				winScore[k] = 0
			case IDX_SUB_SCORE_ZM:
				winScore[k] = 1 * score
			case IDX_SUB_SCORE_QGH:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_WHZ:
				winScore[k] = 8 * score
			case IDX_SUB_SCORE_ZYS:
				winScore[k] = 16 * score
			}
		}
	}
}

//总得分计算和得分类型统计
func (room *ZP_RoomData) SumGameScore(WinUser []int) {
	log.Debug("总得分计算和得分类型统计 赢人：%d", len(WinUser))
	log.Debug("补花数：%v", room.FlowerCnt)

	score := room.Source
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	for i := 0; i < UserCnt; i++ {
		playerScore := &room.HuKindScore[i]

		//暗杠
		if room.UserGangScore[i] > 0 {
			playerScore[IDX_SUB_SCORE_AG] = room.UserGangScore[i]
		}
		room.SumScore[i] += room.UserGangScore[i]

		//胜者
		winCnt := 0
		for k := range WinUser {
			if WinUser[k] == i {
				winCnt++
				break
			}
		}
		if winCnt == 0 {
			continue
		}

		nowCnt := 0
		if i == room.ProvideUser && winCnt == 1 { //自摸情况
			for index := 0; index < UserCnt; index++ {
				if index == i {
					continue
				}
				nowCnt++

				//基础分
				playerScore[IDX_SUB_SCORE_JC] = 1 * score
				room.SumScore[index] -= 1 * score
				room.SumScore[i] += 1 * score
				log.Debug("基础分:%d,SumScore:%d", playerScore[IDX_SUB_SCORE_JC], room.SumScore[i])

				//胡牌
				testScore := 0 //todo,测试
				for j := IDX_SUB_SCORE_ZPKZ; j < COUNT_KIND_SCORE; j++ {
					testScore += playerScore[j] //todo,测试
					room.SumScore[i] += playerScore[j]
					room.SumScore[index] -= playerScore[j]
				}
				log.Debug("胡牌分：%d", testScore)

				//庄家分
				if i == room.BankerUser {
					playerScore[IDX_SUB_SCORE_ZF] = 1 * score
					room.SumScore[i] += 1 * score
					room.SumScore[index] -= 1 * score
				} else if index == room.BankerUser {
					playerScore[IDX_SUB_SCORE_ZF] = 1 * score
					room.SumScore[i] += 1 * score
					room.SumScore[index] -= 1 * score
				}
				log.Debug("庄家分:%d SumScore:%d", playerScore[IDX_SUB_SCORE_ZF], room.SumScore[i])

				//连庄
				if i == room.BankerUser { //庄W
					room.SumScore[index] -= room.LianZhuang * score
					playerScore[IDX_SUB_SCORE_LZ] = room.LianZhuang * score
					room.SumScore[room.BankerUser] += room.LianZhuang * score
				} else if room.ProvideUser == room.BankerUser && nowCnt == 0 { // 边W
					playerScore[IDX_SUB_SCORE_LZ] = room.LianZhuang * score
					room.SumScore[i] += room.LianZhuang * score
					room.SumScore[room.BankerUser] -= room.LianZhuang * score
				}
				log.Debug("连庄得分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_LZ], room.SumScore[i])

				//补花得分
				if room.FlowerCnt[i] < 8 {
					playerScore[IDX_SUB_SCORE_HUA] = room.FlowerCnt[i] * score
					room.SumScore[index] -= room.FlowerCnt[i] * score
				} else { //八张花牌
					playerScore[IDX_SUB_SCORE_HUA] = 16 * score
					room.SumScore[index] -= 16 * score
				}

				//插花分
				room.ChaHuaScore[i] += (room.ChaHuaMap[i] + room.ChaHuaMap[index]) * score
				room.ChaHuaScore[index] = -1 * (room.ChaHuaMap[i] + room.ChaHuaMap[index]) * score
				room.SumScore[index] -= (room.ChaHuaMap[i] + room.ChaHuaMap[index]) * score
				room.SumScore[i] += (room.ChaHuaMap[i] + room.ChaHuaMap[index]) * score
				log.Debug("插花分：%d SumScore:%d", room.ChaHuaMap[i]+room.ChaHuaMap[index], room.SumScore[i])

				//抓花分
				playerScore[IDX_SUB_SCORE_ZH] = room.ZhuaHuaScore[i] * score
				room.SumScore[index] -= room.ZhuaHuaScore[i] * score
				room.SumScore[i] += room.ZhuaHuaScore[i] * score
				log.Debug("抓花分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_ZH], room.SumScore[i])
			}
			room.SumScore[i] += playerScore[IDX_SUB_SCORE_HUA]

			log.Debug("自摸i:%d ,庄家：%d", i, room.BankerUser)
			log.Debug("补分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_HUA], room.SumScore[i])
		} else {
			//基础分
			playerScore[IDX_SUB_SCORE_JC] = 1 * score
			room.SumScore[room.ProvideUser] -= 1 * score
			room.SumScore[i] += 1 * score
			log.Debug("基础分:%d,SumScore:%d", playerScore[IDX_SUB_SCORE_JC], room.SumScore[i])

			//胡牌
			testScore := 0 //todo,测试
			for j := IDX_SUB_SCORE_ZPKZ; j < COUNT_KIND_SCORE; j++ {
				testScore += playerScore[j] //todo,测试
				room.SumScore[i] += playerScore[j]
				room.SumScore[room.ProvideUser] -= playerScore[j]
			}
			log.Debug("胡牌分：%d", testScore)

			//补花分
			if room.FlowerCnt[i] < 8 {
				playerScore[IDX_SUB_SCORE_HUA] = room.FlowerCnt[i] * score
			} else { //八张花牌
				playerScore[IDX_SUB_SCORE_HUA] = 16 * score
			}
			room.SumScore[i] += playerScore[IDX_SUB_SCORE_HUA]
			room.SumScore[room.ProvideUser] -= playerScore[IDX_SUB_SCORE_HUA]
			log.Debug("补花得分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_HUA], room.SumScore[i])

			//庄家分
			if i == room.BankerUser || room.BankerUser == room.ProvideUser {
				playerScore[IDX_SUB_SCORE_ZF] = 1 * score
				room.SumScore[i] += 1 * score
				room.SumScore[room.ProvideUser] -= 1 * score
			}
			log.Debug("庄家分:%d SumScore:%d", playerScore[IDX_SUB_SCORE_ZF], room.SumScore[i])

			//连庄
			if i == room.BankerUser { //庄W
				room.SumScore[room.ProvideUser] -= room.LianZhuang * score
				playerScore[IDX_SUB_SCORE_LZ] = room.LianZhuang * score
				room.SumScore[room.BankerUser] += room.LianZhuang * score
			} else if room.ProvideUser == room.BankerUser { // 边W
				playerScore[IDX_SUB_SCORE_LZ] = room.LianZhuang * score
				room.SumScore[i] += room.LianZhuang * score
				room.SumScore[room.BankerUser] -= room.LianZhuang * score
			}
			log.Debug("i:%d ,庄家：%d", i, room.BankerUser)
			log.Debug("连庄得分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_LZ], room.SumScore[i])

			//插花分
			room.ChaHuaScore[i] = (room.ChaHuaMap[i] + room.ChaHuaMap[room.ProvideUser]) * score
			room.ChaHuaScore[room.ProvideUser] += -1 * (room.ChaHuaMap[i] + room.ChaHuaMap[room.ProvideUser]) * score
			room.SumScore[room.ProvideUser] -= (room.ChaHuaMap[i] + room.ChaHuaMap[room.ProvideUser]) * score
			room.SumScore[i] += room.ChaHuaScore[i]
			log.Debug("插花分：%d SumScore:%d", room.ChaHuaMap[i]+room.ChaHuaMap[room.ProvideUser], room.SumScore[i])

			//抓花分
			playerScore[IDX_SUB_SCORE_ZH] = room.ZhuaHuaScore[i] * score
			room.SumScore[room.ProvideUser] -= room.ZhuaHuaScore[i] * score
			room.SumScore[i] += room.ZhuaHuaScore[i] * score
			log.Debug("抓花分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_ZH], room.SumScore[i])
		}

		//分饼
		if room.BankerUser == i {
			room.SumScore[i] += room.FollowCardScore[i] * score
		} else {
			playerScore[IDX_SUB_SCORE_FB] = room.FollowCardScore[i] * score
			room.SumScore[i] += room.FollowCardScore[i] * score
		}
		log.Debug("分饼分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_FB], room.SumScore[i])
	}
	log.Debug("游戏总分：%d", room.SumScore)
}

func (room *ZP_RoomData) SendStatusPlay(u *user.User) {
	StatusPlay := &msg.G2C_StatusPlay{}
	//自定规则
	StatusPlay.TimeOutCard = room.MjBase.TimerMgr.GetTimeOutCard()
	StatusPlay.TimeOperateCard = room.MjBase.TimerMgr.GetTimeOperateCard()
	StatusPlay.CreateTime = room.MjBase.TimerMgr.GetCreatrTime()
	StatusPlay.PlayerCount = room.MjBase.TimerMgr.GetPlayCount()

	//重入取消托管
	room.MjBase.OnUserTrustee(u.ChairId, false)

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
	StatusPlay.BuHuaCnt = make([]int, UserCnt)
	StatusPlay.ChaHuaCnt = make([]int, UserCnt)
	StatusPlay.BanOutCard = make([]int, UserCnt)
	StatusPlay.BuHuaCard = make([]int, UserCnt)

	StatusPlay.ZhuaHuaCnt = room.ZhuaHuaCnt
	for k, v := range room.ChaHuaMap {
		StatusPlay.ChaHuaCnt[k] = v
	}
	for i := 0; i < len(room.FlowerCnt); i++ {
		StatusPlay.BuHuaCnt[i] = room.FlowerCnt[i]
		if room.FlowerCnt[i] > 0 {
			index := room.FlowerCnt[i] - 1
			StatusPlay.BuHuaCard[i] = room.FlowerCard[i][index]
		}
	}

	//状态变量
	StatusPlay.ActionCard = room.ProvideCard
	StatusPlay.LeftCardCount = room.GetLeftCard()
	StatusPlay.ActionMask = room.UserAction[u.ChairId]

	StatusPlay.Ting = room.Ting
	//当前能胡的牌
	StatusPlay.OutCardCount = room.MjBase.LogicMgr.AnalyseTingCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId],
		StatusPlay.OutCardDataEx, StatusPlay.HuCardCount, StatusPlay.HuCardData)

	//历史记录
	StatusPlay.OutCardUser = room.OutCardUser
	StatusPlay.OutCardData = room.OutCardData
	StatusPlay.DiscardCard = room.DiscardCard
	for _, v := range room.DiscardCard {
		StatusPlay.DiscardCount = append(StatusPlay.DiscardCount, len(v))
	}
	for i := 0; i < len(room.BanCardCnt); i++ {
		if room.BanUser[i]&LimitChi != 0 {
			StatusPlay.BanOutCard[i] = room.BanCardCnt[i][LimitChi]
		}
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

	StatusPlay.TurnScore = room.HistorySe.AllScore
	StatusPlay.CollectScore = room.HistorySe.DetailScore
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

	var WinUser []int
	WinCount := 0
	for i := 0; i < UserCnt; i++ {
		if WIK_CHI_HU == room.ChiHuKind[i] {
			WinUser = append(WinUser, i)
			room.CurrentUser = i
			room.SpecialCardScore(i)
			WinCount++
		}
	}

	//进行抓花
	ZhongCard, BuZhong := room.OnZhuaHua(WinUser)
	log.Debug("========================= ZhongCard:%d,BuZhong:%d", len(ZhongCard), len(BuZhong))
	if WinCount > 0 {
		//抓花派位
		if room.ZhuaHuaCnt > 0 {
			for k, v := range WinUser {
				for _, cardV := range ZhongCard[k] {
					for {
						randV, randOk := utils.RandInt(0, 16)
						if randOk == nil && room.ZhuaHuaMap[randV] == nil {
							room.ZhuaHuaScore[v]++

							huaUser := mj_zp_msg.HuaUser{}
							huaUser.Card = cardV
							log.Debug("中花：%d", cardV)
							huaUser.ChairID = v
							huaUser.IsZhong = true
							room.ZhuaHuaMap[randV] = &huaUser
							break
						}
					}
				}
			}
			for _, cardV2 := range BuZhong {
				for {
					randV, randOk := utils.RandInt(0, 16)
					if randOk == nil && room.ZhuaHuaMap[randV] == nil {
						huaUser := mj_zp_msg.HuaUser{}
						huaUser.Card = cardV2
						//huaUser.ChairID = v
						log.Debug("不中花：%d", cardV2)
						huaUser.IsZhong = false
						room.ZhuaHuaMap[randV] = &huaUser
						break
					}
				}
			}
		}

		//总分
		room.SumGameScore(WinUser)

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
				room.LianZhuang = 0 //连庄次数
			} else {
				room.LianZhuang += 1 //连庄次数
			}
		} else {
			if WinUser[0] == room.BankerUser {
				room.BankerUser = room.BankerUser
				room.LianZhuang += 1 //连庄次数
			} else {
				room.BankerUser += 1
				room.LianZhuang = 0 //连庄次数
			}
		}

		if room.BankerUser > 3 {
			room.BankerUser = 0
		}
	} else { //荒庄
		room.BankerUser = room.BankerUser
		room.LianZhuang += 1 //连庄次数
	}
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

//出牌禁忌--吃啥打啥
func (room *ZP_RoomData) OutOfChiCardRule(CardData, ChairId int) bool {
	defer func() {
		room.BanUser[ChairId] = 0
		room.BanCardCnt[ChairId] = make([]int, 9)
	}()
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

//暗杠
func (room *ZP_RoomData) AnGang(u *user.User, cbOperateCode int, cbOperateCard []int) int {
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

	//清除操作定时
	room.StopOperateCardTimer(u)

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
		room.RecordBanCard(LimitChi, wTargetUser)
	} else if cbTargetAction&WIK_PENG != 0 {
		OperateResult.OperateCard[1] = cbTargetCard
		OperateResult.OperateCard[2] = cbTargetCard
	}

	//用户状态
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	room.IsResponse = make([]bool, UserCnt)
	room.UserAction = make([]int, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.OperateCard[i] = make([]int, 60)
	}
	room.PerformAction = make([]int, UserCnt)
	if cbTargetAction != WIK_GANG {
		nowUser := room.MjBase.UserMgr.GetUserByChairId(wTargetUser)
		room.OutCardTimerEx(nowUser)
	}

	//如果非杠牌
	if cbTargetAction != WIK_GANG {
		room.ProvideUser = INVALID_CHAIR
		room.ProvideCard = 0

		gcr := &mj_base.TagGangCardResult{}
		room.UserAction[wTargetUser] |= room.MjBase.LogicMgr.AnalyseGangCard(room.CardIndex[wTargetUser], room.WeaveItemArray[wTargetUser], 0, gcr)

		//听牌判断
		//if room.Ting[wTargetUser] == false {
		//	HuData := &mj_zp_msg.G2C_ZPMJ_HuData{OutCardData: make([]int, room.GetCfg().MaxCount), HuCardCount: make([]int, room.GetCfg().MaxCount), HuCardData: make([][]int, room.GetCfg().MaxCount), HuCardRemainingCount: make([][]int, room.GetCfg().MaxCount)}
		//	for k := 0; k < room.GetCfg().MaxCount; k++ {
		//		HuData.HuCardData[k] = make([]int, 28)
		//		HuData.HuCardRemainingCount[k] = make([]int, 28)
		//	}
		//
		//	cbCount := room.MjBase.LogicMgr.AnalyseTingCard(room.CardIndex[wTargetUser], room.WeaveItemArray[wTargetUser], HuData.OutCardData, HuData.HuCardCount, HuData.HuCardData, room.GetCfg().MaxCount)
		//	HuData.OutCardCount = cbCount
		//	if cbCount > 0 {
		//		room.UserAction[wTargetUser] |= WIK_LISTEN
		//		for i := 0; i < room.GetCfg().MaxCount; i++ {
		//			if HuData.HuCardCount[i] > 0 {
		//				for j := 0; j < HuData.HuCardCount[i]; j++ {
		//					HuData.HuCardRemainingCount[i][j] = room.GetRemainingCount(wTargetUser, HuData.HuCardData[i][j])
		//				}
		//			} else {
		//				break
		//			}
		//		}
		//		u := room.MjBase.UserMgr.GetUserByChairId(wTargetUser)
		//		u.WriteMsg(HuData)
		//	}
		//}
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
	if !room.IsEnoughCard() {
		log.Debug("荒庄结束,room.LeftCardCount:%d,room.EndLeftCount:%d", room.GetLeftCard(), room.EndLeftCount)
		room.ProvideUser = INVALID_CHAIR
		return 1
	}

	//发送扑克
	room.ProvideCard = room.GetSendCard(bTail, room.MjBase.UserMgr.GetMaxPlayerCnt())
	if room.MjBase.UserMgr.IsTrustee(wCurrentUser) {
		for {
			if room.ProvideCard >= 0x41 && room.ProvideCard <= 0x48 {
				oldCard := room.ProvideCard
				room.ProvideCard = room.GetSendCard(true, room.MjBase.UserMgr.GetMaxPlayerCnt())
				room.SendReplaceCard(wCurrentUser, oldCard, room.ProvideCard, false)
				room.FlowerCnt[wCurrentUser]++
				newCardIndex := SwitchToCardIndex(room.ProvideCard)
				oldCardIndex := SwitchToCardIndex(oldCard)
				room.CardIndex[wCurrentUser][newCardIndex]++
				room.CardIndex[wCurrentUser][oldCardIndex]--
				room.FlowerCard[wCurrentUser] = append(room.FlowerCard[wCurrentUser], oldCard)
				log.Debug("用户%d补花数：%d %d", wCurrentUser, room.FlowerCnt[wCurrentUser], oldCard)
			} else {
				break
			}
		}
	}
	room.SendCardData = room.ProvideCard
	room.LastCatchCardUser = wCurrentUser
	//清除禁止胡牌的牌

	u := room.MjBase.UserMgr.GetUserByChairId(wCurrentUser)
	if u == nil {
		log.Error("at DispatchCardData not foud user ")
		return -1
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
		hu, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser], room.SendCardData)
		if hu {
			room.UserAction[wCurrentUser] |= WIK_CHI_HU
		}
		room.CardIndex[wCurrentUser][room.MjBase.LogicMgr.SwitchToCardIndex(room.SendCardData)]++

		//杠牌判断
		if room.IsEnoughCard() && !room.Ting[wCurrentUser] {
			log.Debug("###############before 杠：%d", room.UserAction[wCurrentUser])
			GangCardResult := &mj_base.TagGangCardResult{}
			room.UserAction[wCurrentUser] |= room.MjBase.LogicMgr.AnalyseGangCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser], room.ProvideCard, GangCardResult)
			log.Debug("###############after 杠：%d", room.UserAction[wCurrentUser])
		}

		if room.FlowerCnt[wCurrentUser] == 8 {
			room.UserAction[wCurrentUser] |= WIK_CHI_HU
		}
	}

	//听牌判断
	//HuData := &mj_zp_msg.G2C_ZPMJ_HuData{OutCardData: make([]int, room.GetCfg().MaxCount), HuCardCount: make([]int, room.GetCfg().MaxCount), HuCardData: make([][]int, room.GetCfg().MaxCount), HuCardRemainingCount: make([][]int, room.GetCfg().MaxCount)}
	//for i := 0; i < room.GetCfg().MaxCount; i++ {
	//	HuData.HuCardData[i] = make([]int, 28)
	//	HuData.HuCardRemainingCount[i] = make([]int, 28)
	//}
	//
	//if room.Ting[wCurrentUser] == false {
	//	cbCount := room.MjBase.LogicMgr.AnalyseTingCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser], HuData.OutCardData, HuData.HuCardCount, HuData.HuCardData, room.GetCfg().MaxCount)
	//	room.TingCnt[wCurrentUser] = int(cbCount)
	//	HuData.OutCardCount = int(cbCount)
	//	if cbCount > 0 {
	//		room.UserAction[wCurrentUser] |= WIK_LISTEN
	//
	//		for i := 0; i < room.GetCfg().MaxCount; i++ {
	//			if HuData.HuCardCount[i] > 0 {
	//				for j := 0; j < HuData.HuCardCount[i]; j++ {
	//					HuData.HuCardRemainingCount[i] = append(HuData.HuCardRemainingCount[i], room.GetRemainingCount(wCurrentUser, HuData.HuCardData[i][j]))
	//				}
	//			} else {
	//				break
	//			}
	//		}
	//
	//		u.WriteMsg(HuData)
	//	}
	//}

	if room.SendStatus != BuHua_Send {
		//构造数据
		SendCardToMe := &mj_zp_msg.G2C_ZPMJ_SendCard{}
		SendCardToMe.SendCardUser = wCurrentUser
		SendCardToMe.CurrentUser = wCurrentUser
		SendCardToMe.Tail = bTail
		SendCardToMe.ActionMask = room.UserAction[wCurrentUser]
		SendCardToMe.CardData = room.ProvideCard
		u.WriteMsg(SendCardToMe)

		SendCard := &mj_zp_msg.G2C_ZPMJ_SendCard{}
		SendCard.SendCardUser = wCurrentUser
		SendCard.CurrentUser = wCurrentUser
		SendCard.Tail = bTail
		SendCard.ActionMask = room.UserAction[wCurrentUser]
		SendCard.CardData = 0
		room.MjBase.UserMgr.SendMsgAllNoSelf(u.Id, SendCard)
	}

	//超时定时
	room.UserActionDone = false
	if room.MjBase.UserMgr.IsTrustee(wCurrentUser) {
		room.UserActionDone = true
		cardindex := room.GetTrusteeOutCard(u.ChairId)
		if cardindex == INVALID_BYTE {
			return 0
		}
		card := room.MjBase.LogicMgr.SwitchToCardData(cardindex)
		room.MjBase.OutCard([]interface{}{u, card, true})
	} else {
		room.OutCardTimer(u)
	}
	return 0
}

//解散接触
func (room *ZP_RoomData) DismissEnd(cbReason int) {
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
	GameConclude.Reason = cbReason
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

	StatusFree.TurnScore = room.HistorySe.AllScore
	StatusFree.CollectScore = room.HistorySe.DetailScore

	StatusFree.PlayerCount = room.MjBase.TimerMgr.GetPlayCount() //玩家人数
	StatusFree.MaCount = 0                                       //码数
	StatusFree.CountLimit = room.MjBase.TimerMgr.GetMaxPlayCnt() //局数限制
	StatusFree.ZhuaHuaCnt = room.ZhuaHuaCnt
	u.WriteMsg(StatusFree)
}

//重置用户状态
func (room *ZP_RoomData) ResetUserOperateEx(u *user.User) {
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	room.UserAction = make([]int, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
	room.StopOperateCardTimer(u)
}

func (room *ZP_RoomData) GetTrusteeOutCard(wChairID int) int {
	cardindex := INVALID_BYTE
	if room.SendCardData != 0 && room.SendCardData < 0x41 {
		cardindex = room.MjBase.LogicMgr.SwitchToCardIndex(room.SendCardData)
	} else {
		for i := room.GetCfg().MaxIdx - 1 - room.GetCfg().HuaCount; i > 0; i-- {
			cardindex = i
			if room.CardIndex[wChairID][cardindex] > 0 {
				card := room.MjBase.LogicMgr.SwitchToCardData(cardindex)
				if !(card == room.BanCardCnt[wChairID][LimitChi] && room.BanUser[wChairID]&LimitChi != 0) {
					break
				} else {
					log.Debug("超时吃啥打啥")
				}
			}
		}
	}
	return cardindex
}

func (room *ZP_RoomData) CheckZiMo() {

}

///////////////////////////////////////////////////////////////////////////////////
//定时器

//出牌定时
func (room *ZP_RoomData) OutCardTimer(u *user.User) {
	//stop
	if room.OutCardTime != nil {
		log.Debug("停出牌定时 %d", u.ChairId)
		room.OutCardTime.Stop()
	}

	room.OutCardTime = room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OutCardTime)*time.Second, func() {
		log.Debug("超时---出牌用户： %d", u.ChairId)
		room.MjBase.OnUserTrustee(u.ChairId, true)
	})
}

//出牌定时2
func (room *ZP_RoomData) OutCardTimerEx(u *user.User) {
	//stop
	if room.OutCardTime != nil {
		log.Debug("停出牌定时 %d", u.ChairId)
		room.OutCardTime.Stop()
	}

	room.OutCardTime = room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OperateCardTime)*time.Second, func() {
		log.Debug("超时---出牌 %d", u.ChairId)
		room.MjBase.OnUserTrustee(u.ChairId, true)
	})
}

//开局定时器
func (room *ZP_RoomData) InitOutCardTimer(u *user.User) {
	//stop
	if room.OutCardTime != nil {
		room.OutCardTime.Stop()
	}

	room.OutCardTime = room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OutCardTime)*time.Second, func() {
		log.Debug("开局超时---出牌 %d", u.ChairId)
		card := 0
		for j := room.GetCfg().MaxIdx - 1; j > 0; j-- {
			if room.CardIndex[u.ChairId][j] > 0 {
				card = room.MjBase.LogicMgr.SwitchToCardData(j)
				break
			}
		}
		log.Debug("用户%d开局超时：%x", u.ChairId, card)
		room.MjBase.OutCard([]interface{}{u, card, true})
	})
}

//操作定时
func (room *ZP_RoomData) OperateCardTimer(u *user.User) {
	chairID := u.ChairId

	if room.OutCardTime != nil {
		log.Debug("OperateCardTimer停出牌定时 %d", u.ChairId)
		room.OutCardTime.Stop()
	}
	if room.OperateTime[chairID] != nil {
		log.Debug("停吃碰杠定时器 %d", u.ChairId)
		room.OperateTime[chairID].Stop()
	}

	operateTimer := room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OperateCardTime)*time.Second, func() {
		log.Debug("超时---吃碰杠定时器 %d", u.ChairId)
		if room.UserAction[chairID] != WIK_LISTEN {
			operateCard := []int{0, 0, 0}
			room.MjBase.UserOperateCard([]interface{}{u, WIK_NULL, operateCard})
		} else {
			room.OnUserListenCard(u, false)
		}
		//room.OnUserTrustee(chairID, true)
	})
	room.OperateTime[chairID] = operateTimer
}

//清理定时器
func (room *ZP_RoomData) StopOperateCardTimer(u *user.User) {
	chairID := u.ChairId

	if room.OperateTime[chairID] != nil {
		log.Debug("清除操作定时 user:%d", chairID)
		log.Debug("zpmj at StopOperateCardTimer user:%d", chairID)
		room.OperateTime[chairID].Stop()
	}
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
