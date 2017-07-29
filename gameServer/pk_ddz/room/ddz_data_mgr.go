package room

import (
	"math/rand"
	"mj/common/cost"
	"mj/common/msg/pk_ddz_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"encoding/json"

	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

func NewDDZDataMgr(info *model.CreateRoomInfo, uid int64, ConfigIdx int, name string, temp *base.GameServiceOption, base *DDZ_Entry) *ddz_data_mgr {
	d := new(ddz_data_mgr)
	d.RoomData = pk_base.NewDataMgr(info.RoomId, uid, ConfigIdx, name, temp, base.Entry_base, info.OtherInfo)
	d.initParam()

	var setInfo pk_ddz_msg.C2G_DDZ_CreateRoomInfo
	if err := json.Unmarshal([]byte(info.OtherInfo), &setInfo); err == nil {
		d.EightKing = setInfo.King
		d.GameType = setInfo.GameType
	}
	return d
}

type ddz_data_mgr struct {
	*pk_base.RoomData
	GameStatus int // 当前游戏状态

	CurrentUser int // 当前玩家

	BankerUser int // 地主
	TurnWiner  int // 出牌玩家

	EightKing bool // 是否八王模式
	GameType  int  // 游戏类型
	LiziCard  int  // 癞子牌

	// 炸弹信息
	EachBombCount []int // 炸弹个数
	KingCount     []int // 八王个数

	// 叫分信息
	ScoreInfo []int // 叫分信息

	// 出牌信息
	TurnCardStatus []int                          // 用户出牌状态
	TurnCardData   [][]pk_ddz_msg.C2G_DDZ_OutCard // 出牌数据
	RepertoryCard  []int                          // 库存扑克

	// 扑克信息
	BankerCard   [3]int  // 游戏底牌
	HandCardData [][]int // 手上扑克
	ShowCardSign []bool  // 用户明牌标识

	// 定时器
	CardTimer *time.Timer
}

func (room *ddz_data_mgr) resetData() {
	room.GameStatus = GAME_STATUS_FREE
	room.CurrentUser = cost.INVALID_CHAIR
	room.BankerUser = cost.INVALID_CHAIR
	room.TurnWiner = cost.INVALID_CHAIR
	room.LiziCard = 0
	room.EachBombCount = make([]int, room.PlayerCount)
	room.KingCount = append([]int{})
	room.ScoreInfo = make([]int, room.PlayerCount)
	room.TurnCardStatus = make([]int, room.PlayerCount)
	room.TurnCardData = make([][]pk_ddz_msg.C2G_DDZ_OutCard, room.PlayerCount)

	nMaxCardCount := room.GetCfg().MaxRepertory
	if room.EightKing {
		nMaxCardCount += 6
	}
	room.RepertoryCard = make([]int, nMaxCardCount)
	room.BankerCard = [3]int{}
	room.HandCardData = append([][]int{})

	room.ScoreTimes = 0
	room.HistoryScores = make([]*pk_base.HistoryScore, room.PlayerCount)
}

func (room *ddz_data_mgr) InitRoom(UserCnt int) {
	log.Debug("初始化房间参数%d", UserCnt)
	room.RoomData.InitRoom(UserCnt)
	room.PlayerCount = UserCnt

	room.resetData()
}

// 初始化部分数据
func (r *ddz_data_mgr) initParam() {
	// 明牌标识
	r.ShowCardSign = append([]bool{})
	for i := 0; i < 3; i++ {
		r.ShowCardSign = append(r.ShowCardSign, false)
	}
}

// 空闲状态场景
func (room *ddz_data_mgr) SendStatusReady(u *user.User) {
	log.Debug("发送空闲状态场景消息")
	room.GameStatus = GAME_STATUS_FREE
	StatusFree := &pk_ddz_msg.G2C_DDZ_StatusFree{}

	StatusFree.CellScore = room.PkBase.Temp.Source // 基础积分

	StatusFree.GameType = room.GameType
	StatusFree.EightKing = room.EightKing

	StatusFree.PlayCount = room.PkBase.TimerMgr.GetMaxPlayCnt()

	StatusFree.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()       // 出牌时间
	StatusFree.TimeCallScore = room.GetCfg().CallScoreTime               // 叫分时间
	StatusFree.TimeStartGame = room.PkBase.TimerMgr.GetTimeOperateCard() // 开始时间 	// 首出时间
	StatusFree.TurnScore = append(StatusFree.TurnScore, 12)
	StatusFree.CollectScore = append(StatusFree.CollectScore, 11)
	//for _, v := range room.HistoryScores {
	//	StatusFree.TurnScore = append(StatusFree.TurnScore, v.TurnScore)
	//	StatusFree.CollectScore = append(StatusFree.TurnScore, v.CollectScore)
	//}

	// 发送明牌标识
	StatusFree.ShowCardSign = make([]bool, len(room.ShowCardSign))
	util.DeepCopy(&StatusFree.ShowCardSign, &room.ShowCardSign)

	// 发送托管标识
	trustees := room.PkBase.UserMgr.GetTrustees()
	StatusFree.TrusteeSign = make([]bool, len(trustees))
	util.DeepCopy(&StatusFree.TrusteeSign, &trustees)

	u.WriteMsg(StatusFree)

}

// 开始游戏前
func (room *ddz_data_mgr) BeforeStartGame(UserCnt int) {
	room.InitRoom(UserCnt)
}

// 游戏开始
func (room *ddz_data_mgr) StartGameing() {
	room.GameStatus = GAME_STATUS_CALL
	room.SendGameStart()
}

// 叫分状态
func (room *ddz_data_mgr) SendStatusCall(u *user.User) {
	StatusCall := &pk_ddz_msg.G2C_DDZ_StatusCall{}

	StatusCall.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()
	StatusCall.TimeCallScore = room.PkBase.TimerMgr.GetTimeOperateCard()
	StatusCall.TimeStartGame = int(room.PkBase.TimerMgr.GetCreatrTime())

	StatusCall.GameType = room.GameType
	StatusCall.LaiziCard = room.LiziCard
	StatusCall.EightKing = room.EightKing
	StatusCall.CellScore = room.PkBase.Temp.Source
	StatusCall.CurrentUser = room.CurrentUser
	StatusCall.BankerScore = room.ScoreTimes
	StatusCall.ScoreInfo = util.CopySlicInt(room.ScoreInfo)
	StatusCall.HandCardCount = make([]int, len(room.HandCardData))
	for i := 0; i < len(room.HandCardData); i++ {
		StatusCall.HandCardCount[i] = len(room.HandCardData[i])
	}

	UserCnt := room.PkBase.UserMgr.GetMaxPlayerCnt()
	StatusCall.ScoreInfo = util.CopySlicInt(room.ScoreInfo)

	StatusCall.TurnScore = make([]int, UserCnt)
	StatusCall.CollectScore = make([]int, UserCnt)

	//历史积分
	for j := 0; j < UserCnt; j++ {
		//设置变量
		if room.HistoryScores[j] != nil {
			StatusCall.TurnScore[j] = room.HistoryScores[j].TurnScore
			StatusCall.CollectScore[j] = room.HistoryScores[j].CollectScore
		}
	}

	StatusCall.ShowCardSign = make([]bool, len(room.ShowCardSign))
	util.DeepCopy(&StatusCall.ShowCardSign, &room.ShowCardSign)
	//发送数据
	for i := 0; i < room.PkBase.Temp.MaxPlayer; i++ {
		if room.ShowCardSign[i] || u.ChairId == i {
			StatusCall.ShowCardData = append(StatusCall.ShowCardData, room.HandCardData[i])
		} else {
			StatusCall.ShowCardData = append(StatusCall.ShowCardData, nil)
		}
	}

	log.Debug("叫分进行时%v", StatusCall)
	u.WriteMsg(StatusCall)
}

// 游戏状态
func (room *ddz_data_mgr) SendStatusPlay(u *user.User) {
	if room.GameStatus == GAME_STATUS_CALL {
		room.SendStatusCall(u)
		return
	}
	if room.GameStatus == GAME_STATUS_FREE {
		room.SendStatusReady(u)
		return
	}
	StatusPlay := &pk_ddz_msg.G2C_DDZ_StatusPlay{}
	//自定规则
	StatusPlay.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()
	StatusPlay.TimeCallScore = room.PkBase.TimerMgr.GetTimeOperateCard()
	StatusPlay.TimeStartGame = int(room.PkBase.TimerMgr.GetCreatrTime())

	//游戏变量
	StatusPlay.CellScore = room.PkBase.Temp.Source

	StatusPlay.BankerUser = room.BankerUser
	StatusPlay.CurrentUser = room.CurrentUser
	StatusPlay.BankerScore = room.ScoreTimes
	StatusPlay.EightKing = room.EightKing
	StatusPlay.GameType = room.GameType
	StatusPlay.LaiziCard = room.LiziCard

	StatusPlay.TurnWiner = room.TurnWiner
	if StatusPlay.TurnWiner != cost.INVALID_CHAIR {
		turnCardData := room.TurnCardData[room.TurnWiner][len(room.TurnCardData[room.TurnWiner])-1]
		util.DeepCopy(&StatusPlay.TurnCardData, &turnCardData)
	}
	util.DeepCopy(&StatusPlay.BankerCard, &room.BankerCard)
	StatusPlay.HandCardCount = make([]int, room.PlayerCount)
	for i := 0; i < room.PlayerCount; i++ {
		StatusPlay.HandCardCount[i] = len(room.HandCardData[i])
	}
	StatusPlay.EachBombCount = util.CopySlicInt(room.EachBombCount)
	StatusPlay.KingCount = util.CopySlicInt(room.KingCount)

	UserCnt := room.PkBase.UserMgr.GetMaxPlayerCnt()
	StatusPlay.TurnScore = make([]int, UserCnt)
	StatusPlay.CollectScore = make([]int, UserCnt)

	//历史积分
	for j := 0; j < UserCnt; j++ {
		//设置变量
		if room.HistoryScores[j] != nil {
			StatusPlay.TurnScore[j] = room.HistoryScores[j].TurnScore
			StatusPlay.CollectScore[j] = room.HistoryScores[j].CollectScore
		}
	}

	StatusPlay.ShowCardSign = make([]bool, len(room.ShowCardSign))
	util.DeepCopy(&StatusPlay.ShowCardSign, &room.ShowCardSign)

	// 发送数据
	for i := 0; i < room.PkBase.Temp.MaxPlayer; i++ {
		if room.ShowCardSign[i] || u.ChairId == i {
			StatusPlay.ShowCardData = append(StatusPlay.ShowCardData, room.HandCardData[i])
			//util.DeepCopy(GameStart.CardData[i], cardData[i])
		} else {
			StatusPlay.ShowCardData = append(StatusPlay.ShowCardData, nil)
		}
	}

	log.Debug("游戏进行时%v", StatusPlay)
	u.WriteMsg(StatusPlay)
}

// 开始游戏，发送扑克
func (room *ddz_data_mgr) SendGameStart() {

	userMgr := room.PkBase.UserMgr
	gameLogic := room.PkBase.LogicMgr

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, cost.US_PLAYING)
	})

	// 打乱牌
	gameLogic.RandCardList(room.RepertoryCard, pk_base.GetCardByIdx(room.ConfigIdx))

	// 底牌
	//util.DeepCopy(room.BankerCard[:], &room.RepertoryCard[len(room.RepertoryCard)-3:])
	log.Debug("发牌数据%v,%v", room.BankerCard, room.RepertoryCard)
	copy(room.BankerCard[:], room.RepertoryCard[len(room.RepertoryCard)-3:])
	room.RepertoryCard = room.RepertoryCard[:len(room.RepertoryCard)-3]

	log.Debug("底牌%v", room.BankerCard)
	log.Debug("剩余牌%v", room.RepertoryCard)

	cardCount := len(room.RepertoryCard) / room.PkBase.Temp.MaxPlayer

	//构造变量
	GameStart := &pk_ddz_msg.G2C_DDZ_GameStart{}

	// 发完牌就选癞子
	if room.GameType == GAME_TYPE_LZ {
		// 随机选一张牌为癞子牌
		ran := rand.New(rand.NewSource(time.Now().UnixNano()))
		room.LiziCard = ran.Intn(13)
		GameStart.LiziCard = room.LiziCard
		room.PkBase.LogicMgr.SetParamToLogic(room.LiziCard)
	}

	// 初始化叫分信息
	for i := 0; i < room.PlayerCount; i++ {
		room.ScoreInfo[i] = CALLSCORE_NOCALL
	}

	// 初始化牌
	room.HandCardData = append([][]int{})
	for i := 0; i < room.PlayerCount; i++ {
		tempCardData := util.CopySlicInt(room.RepertoryCard[len(room.RepertoryCard)-cardCount:])
		room.PkBase.LogicMgr.SortCardList(tempCardData, len(tempCardData))
		room.RepertoryCard = room.RepertoryCard[:len(room.RepertoryCard)-cardCount]
		room.HandCardData = append(room.HandCardData, tempCardData)
		if room.CurrentUser == cost.INVALID_CHAIR {
			for _, v := range tempCardData {
				if v == 0x33 {
					room.CurrentUser = i
					break
				}
			}
		}
	}

	if room.CurrentUser == cost.INVALID_CHAIR {
		room.CurrentUser = util.RandInterval(0, 2)
	}

	room.ScoreInfo[room.CurrentUser] = CALLSCORE_CALLING

	GameStart.CallScoreUser = room.CurrentUser

	GameStart.ShowCard = make([]bool, len(room.ShowCardSign))
	util.DeepCopy(&GameStart.ShowCard, &room.ShowCardSign)

	//发送数据
	room.PkBase.UserMgr.ForEachUser(func(u *user.User) {

		GameStart.CardData = append([][]int{})
		for i := 0; i < room.PkBase.Temp.MaxPlayer; i++ {
			if room.ShowCardSign[i] || u.ChairId == i {
				GameStart.CardData = append(GameStart.CardData, room.HandCardData[i])
				//util.DeepCopy(GameStart.CardData[i], cardData[i])
			} else {
				GameStart.CardData = append(GameStart.CardData, nil)
			}
		}

		log.Debug("需要发送的扑克牌%v", GameStart)
		u.WriteMsg(GameStart)
	})

	// 启动定时器
	room.startOperateCardTimer(room.GetCfg().CallScoreTime)
}

// 用户叫分(抢庄)
func (r *ddz_data_mgr) CallScore(u *user.User, scoreTimes int) {
	// 判断当前是否为叫分状态
	if r.GameStatus != GAME_STATUS_CALL {
		log.Debug("叫分错误，当前游戏状态为%d", r.GameStatus)
		return
	}

	// 判断当前叫分玩家是否正确
	if r.ScoreInfo[u.ChairId] != CALLSCORE_CALLING {
		log.Debug("叫分玩家%d叫分失败，当前玩家叫分状态%v", u.ChairId, r.ScoreInfo)
		cost.RenderErrorMessage(cost.ErrDDZCSUser)
		return
	}

	if scoreTimes <= r.ScoreTimes && scoreTimes != 0 {
		cost.RenderErrorMessage(cost.ErrDDZCSValid)
		log.Debug("用户叫分%d必须大于当前分数%d", scoreTimes, r.ScoreTimes)
		return
	}

	if scoreTimes == 0 {

	} else {
		r.ScoreTimes = scoreTimes
		r.BankerUser = u.ChairId
	}

	r.ScoreInfo[u.ChairId] = scoreTimes

	nextCallUser := (r.CurrentUser + 1) % r.PlayerCount // 下一个叫分玩家
	r.CurrentUser = nextCallUser

	isEnd := (r.ScoreTimes == CALLSCORE_MAX) || (r.ScoreInfo[nextCallUser] != CALLSCORE_NOCALL)

	if !isEnd {
		r.ScoreInfo[nextCallUser] = CALLSCORE_CALLING
		r.resetOperateCardTimer(r.GetCfg().CallScoreTime)
	} else {
		// 叫分结束，看谁叫的分数大就是地主
		var score int
		for i, v := range r.ScoreInfo {
			log.Debug("遍历叫分%d,%d", i, v)
			if v > score && v >= 0 && v <= CALLSCORE_MAX {
				score = v
				r.ScoreTimes = v
				r.BankerUser = i
			}
		}
		// 如果都未叫，则随机选一个作为地主，并且倍数默认为1
		if score == 0 {
			r.BankerUser = util.RandInterval(0, 2)
			r.ScoreTimes = 1
		}
		r.CurrentUser = r.BankerUser
		r.resetOperateCardTimer(r.PkBase.Temp.OutCardTime)
	}

	// 发送叫分信息
	GameCallSore := &pk_ddz_msg.G2C_DDZ_CallScore{}
	util.DeepCopy(&GameCallSore.ScoreInfo, &r.ScoreInfo)
	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		log.Debug("发送叫分信息%v", GameCallSore)
		u.WriteMsg(GameCallSore)
	})

	if isEnd {
		// 叫分结束，发庄家信息
		r.BankerInfo()
	}

	r.checkNextUserTrustee()
}

// 庄家信息
func (r *ddz_data_mgr) BankerInfo() {
	for _, v := range r.BankerCard {
		r.HandCardData[r.BankerUser] = append(r.HandCardData[r.BankerUser], v)
	}
	r.PkBase.LogicMgr.SortCardList(r.HandCardData[r.BankerUser], len(r.HandCardData[r.BankerUser]))
	log.Debug("地主的牌%v", r.HandCardData[r.BankerUser])

	r.GameStatus = GAME_STATUS_PLAY // 确定地主了，进入游戏状态

	for i := 1; i < r.PlayerCount; i++ {
		r.TurnCardStatus[(r.BankerUser+i)%r.PlayerCount] = OUTCARD_PASS
	}
	r.TurnCardStatus[r.BankerUser] = OUTCARD_OUTING

	GameBankerInfo := &pk_ddz_msg.G2C_DDZ_BankerInfo{}
	GameBankerInfo.BankerUser = r.BankerUser
	GameBankerInfo.CurrentUser = r.CurrentUser
	GameBankerInfo.BankerScore = r.ScoreTimes

	util.DeepCopy(&GameBankerInfo.BankerCard, &r.BankerCard)

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		log.Debug("庄家信息%v", GameBankerInfo)
		u.WriteMsg(GameBankerInfo)
	})
}

// 明牌
func (r *ddz_data_mgr) ShowCard(u *user.User) {
	log.Debug("当前明牌信息%v,玩家id", r.ShowCardSign, u.ChairId)

	r.ShowCardSign[u.ChairId] = true

	DataShowCard := &pk_ddz_msg.G2C_DDZ_ShowCard{}
	DataShowCard.ShowCardUser = u.ChairId
	if len(r.HandCardData) > u.ChairId {
		log.Debug("明牌前%v", r.HandCardData[u.ChairId])
		DataShowCard.CardData = util.CopySlicInt(r.HandCardData[u.ChairId])
	}

	log.Debug("%v", r.HandCardData)

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		log.Debug("明牌数据%v", DataShowCard)
		u.WriteMsg(DataShowCard)
	})
}

// 用户出牌
func (r *ddz_data_mgr) OpenCard(u *user.User, cardType int, cardData []int) {
	if cardType == 0 {
		log.Debug("用户%d要不起", u.ChairId)
		r.PassCard(u)
		return
	}
	log.Debug("用户%d出牌%v", u.ChairId, cardData)
	// 检查当前是否是游戏中
	if r.GameStatus != GAME_STATUS_PLAY {
		log.Debug("出牌错误，当前游戏状态为%d", r.GameStatus)
		return
	}
	// 检查当前是否该用户出牌
	if r.TurnCardStatus[u.ChairId] != OUTCARD_OUTING {
		log.Debug("出牌错误，当前出牌人为%d", r.CurrentUser)
		return
	}

	// 检查所出牌是否完整在手上
	if len(cardData) > len(r.HandCardData[u.ChairId]) {
		log.Debug("所出的牌%v，手上的牌%v，超过数量了", cardData, r.HandCardData[u.ChairId])
		return
	}

	var b bool
	for _, vOut := range cardData {
		b = false
		for _, vHand := range r.HandCardData[u.ChairId] {
			if vHand == vOut {
				b = true
				break
			}
		}
		if !b {
			// 所出的牌没在该人手上
			log.Debug("所出的牌%v，手上的牌%v", cardData, r.HandCardData[u.ChairId])
			return
		}
	}

	nowCard := pk_ddz_msg.C2G_DDZ_OutCard{}
	nowCard.CardData = util.CopySlicInt(cardData)

	// 对比牌型
	r.PkBase.LogicMgr.SortCardList(nowCard.CardData, len(nowCard.CardData)) // 排序

	var lastTurnCard pk_ddz_msg.C2G_DDZ_OutCard
	for i := 1; i < r.PlayerCount; i++ {
		lastUser := r.lastUser(u.ChairId)
		uStatus := r.TurnCardStatus[lastUser]
		if uStatus < OUTCARD_MAXCOUNT {
			log.Debug("手牌%d,%d", lastUser, uStatus)
			log.Debug("手牌信息%v", r.TurnCardData)
			lastTurnCard = r.TurnCardData[lastUser][uStatus]
			break
		}
	}

	log.Debug("上一次出牌的信息%v", lastTurnCard)
	if lastTurnCard.CardType == CT_ERROR {
		nowCard.CardType = r.PkBase.LogicMgr.GetCardType(nowCard.CardData)
		if nowCard.CardType == CT_ERROR {
			// 出牌无效
			return
		}
	} else {
		var isType bool
		nowCard.CardType, isType = r.PkBase.LogicMgr.CompareCardWithParam(lastTurnCard.CardData, nowCard.CardData, []interface{}{lastTurnCard.CardType})
		if !isType {
			log.Debug("出牌数据有错")
			return
		}
	}

	// 判断是否火箭
	if nowCard.CardType >= CT_KING {
		r.KingCount = append(r.KingCount, len(nowCard.CardData))
	} else if nowCard.CardType >= CT_BOMB_CARD {
		// 炸弹
		r.EachBombCount[u.ChairId]++
	}

	r.TurnWiner = u.ChairId
	r.CurrentUser = r.nextUser(u.ChairId)
	r.TurnCardStatus[r.CurrentUser] = OUTCARD_OUTING
	// 把所出的牌存到出牌数据里

	r.TurnCardData[u.ChairId] = append(r.TurnCardData[u.ChairId], nowCard)
	r.TurnCardStatus[u.ChairId] = len(r.TurnCardData[u.ChairId]) - 1

	// 从手牌删除数据
	log.Debug("出牌玩家%d删除前的手牌%v", r.TurnWiner, r.HandCardData)
	r.HandCardData[u.ChairId], _ = r.PkBase.LogicMgr.RemoveCardList(cardData, r.HandCardData[u.ChairId])
	log.Debug("删除后的手牌%v", r.HandCardData)

	// 发送给所有玩家
	DataOutCard := &pk_ddz_msg.G2C_DDZ_OutCard{}
	if len(r.HandCardData[u.ChairId]) == 0 {
		DataOutCard.CurrentUser = cost.INVALID_CHAIR
	} else {
		DataOutCard.CurrentUser = r.CurrentUser
	}
	DataOutCard.OutCardUser = u.ChairId
	util.DeepCopy(&DataOutCard.CardData, &nowCard)
	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		log.Debug("出牌数据%v", DataOutCard)
		u.WriteMsg(DataOutCard)
	})

	if len(r.HandCardData[u.ChairId]) == 0 {
		log.Debug("游戏结束")
		r.PkBase.OnEventGameConclude(0, nil, cost.GER_NORMAL)
		return
	}
	r.resetOperateCardTimer(r.PkBase.Temp.OutCardTime)
	r.checkNextUserTrustee()
}

// 判断下一个玩家托管状态
func (r *ddz_data_mgr) checkNextUserTrustee() {
	log.Debug("当前玩家%d,上次出牌玩家%d,托管状态%v", r.CurrentUser, r.TurnWiner, r.PkBase.UserMgr.GetTrustees())
	if r.PkBase.UserMgr.IsTrustee(r.CurrentUser) {
		// 出牌玩家为托管状态
		if r.GameStatus == GAME_STATUS_CALL {
			// 叫分状态，托管则不叫
			r.CallScore(r.PkBase.UserMgr.GetUserByChairId(r.CurrentUser), 0)
		} else if r.GameStatus == GAME_STATUS_PLAY {
			// 出牌状态
			if r.CurrentUser == r.TurnWiner || r.TurnWiner == cost.INVALID_CHAIR {
				// 上一个出牌玩家是自己，则选最小牌
				var cardData []int
				cardData = append(cardData, r.HandCardData[r.CurrentUser][len(r.HandCardData[r.CurrentUser])-1])
				r.OpenCard(r.PkBase.UserMgr.GetUserByChairId(r.CurrentUser), 1, cardData)
			} else {
				// 上一个出牌玩家不是自己，则不出
				log.Debug("玩家不出%d,%d,%d", r, r.PkBase, r.PkBase.UserMgr)
				r.OpenCard(r.PkBase.UserMgr.GetUserByChairId(r.CurrentUser), 0, nil)
			}
		}
	}
}

// 下一个玩家
func (r *ddz_data_mgr) nextUser(u int) int {
	return (u + 1) % r.PlayerCount
}

// 上一个玩家
func (r *ddz_data_mgr) lastUser(u int) int {
	return (u + r.PlayerCount - 1) % r.PlayerCount
}

// 放弃出牌
func (r *ddz_data_mgr) PassCard(u *user.User) {
	// 检查当前是否是游戏中
	if r.GameStatus != GAME_STATUS_PLAY {
		log.Debug("出牌错误，当前游戏状态为%d", r.GameStatus)
		return
	}
	// 检查当前是否该用户出牌
	if r.TurnCardStatus[u.ChairId] != OUTCARD_OUTING {
		log.Debug("出牌错误，当前出牌人为%d", r.CurrentUser)
		return
	}
	// 如果上一个出牌人是自己，则不能放弃
	if r.TurnWiner == u.ChairId || r.TurnWiner == cost.INVALID_CHAIR {
		log.Debug("不允许放弃")
		return
	}

	r.CurrentUser = r.nextUser(u.ChairId)
	r.TurnCardStatus[u.ChairId] = OUTCARD_PASS
	r.TurnCardStatus[r.CurrentUser] = OUTCARD_OUTING

	DataPassCard := &pk_ddz_msg.G2C_DDZ_PassCard{}

	DataPassCard.TurnOver = r.CurrentUser == r.TurnWiner
	DataPassCard.CurrentUser = r.CurrentUser
	DataPassCard.PassCardUser = u.ChairId

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		//log.Debug("用户放弃出牌%v", DataPassCard)
		u.WriteMsg(DataPassCard)
	})

	r.checkNextUserTrustee()
}

// 游戏正常结束
func (r *ddz_data_mgr) NormalEnd() {
	log.Debug("游戏正常结束了")
	r.GameStatus = GAME_STATUS_FREE
	DataGameConclude := &pk_ddz_msg.G2C_DDZ_GameConclude{}
	DataGameConclude.CellScore = r.PkBase.Temp.Source

	// 算分数
	nMultiple := r.ScoreTimes

	// 春天标识
	if len(r.TurnCardData[r.BankerUser]) <= 1 {
		DataGameConclude.SpringSign = 2 // 地主只出了一次牌
	} else {
		DataGameConclude.SpringSign = 1
		for i := 1; i < r.PlayerCount; i++ {
			if len(r.TurnCardData[(i+r.BankerUser)%r.PlayerCount]) > 0 {
				DataGameConclude.SpringSign = 0
				break
			}
		}
	}

	if DataGameConclude.SpringSign > 0 {
		nMultiple <<= 1 // 春天反春天，倍数翻倍
	}

	// 炸弹翻倍
	for _, v := range r.EachBombCount {
		if v > 0 {
			nMultiple <<= uint(v) // 炸弹个数翻倍
		}
	}

	// 地主明牌翻倍
	if r.ShowCardSign[r.BankerUser] == true {
		nMultiple <<= 1
	}

	// 八王
	util.DeepCopy(&DataGameConclude.KingCount, &r.KingCount)
	for _, v := range r.KingCount {
		if v == 8 {
			nMultiple *= 8 * 2
		} else if v >= 2 {
			nMultiple *= v
		}
	}

	// 炸弹
	util.DeepCopy(&DataGameConclude.EachBombCount, &r.EachBombCount)

	DataGameConclude.BankerScore = r.ScoreTimes
	util.DeepCopy(&DataGameConclude.HandCardData, &r.HandCardData)

	// 计算积分
	gameScore := r.PkBase.Temp.Source * nMultiple

	if len(r.HandCardData[r.BankerUser]) <= 0 {
		gameScore = 0 - gameScore
	}

	var score int
	DataGameConclude.GameScore = make([]int, r.PlayerCount)
	for i := 0; i < r.PlayerCount; i++ {
		if i != r.BankerUser {
			if r.ShowCardSign[i] {
				DataGameConclude.GameScore[i] = gameScore * 2
			} else {
				DataGameConclude.GameScore[i] = gameScore
			}
			score += DataGameConclude.GameScore[i]
		}
	}

	DataGameConclude.GameScore[r.BankerUser] = 0 - score

	r.sendGameEndMsg(DataGameConclude)
}

//解散房间结束
func (r *ddz_data_mgr) DismissEnd() {
	DataGameConclude := &pk_ddz_msg.G2C_DDZ_GameConclude{}
	DataGameConclude.CellScore = r.PkBase.Temp.Source
	DataGameConclude.GameScore = make([]int, r.PlayerCount)

	// 炸弹
	util.DeepCopy(&DataGameConclude.EachBombCount, &r.EachBombCount)

	DataGameConclude.BankerScore = r.ScoreTimes
	util.DeepCopy(&DataGameConclude.HandCardData, &r.HandCardData)

	r.sendGameEndMsg(DataGameConclude)
}

// 发送游戏结束消息
func (r *ddz_data_mgr) sendGameEndMsg(DataGameConclude *pk_ddz_msg.G2C_DDZ_GameConclude) {
	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		log.Debug("游戏结算信息%v", DataGameConclude)
		u.WriteMsg(DataGameConclude)
		// 取消所有人的托管状态
		r.PkBase.UserMgr.SetUsetTrustee(u.ChairId, false)
	})
	// 取消定时器
	r.stopOperateCardTimer()
	// 明牌重置
	r.initParam()
	// 重置数据
	r.resetData()
}

// 托管
func (room *ddz_data_mgr) Trustee(u *user.User, t bool) {
	room.PkBase.UserMgr.SetUsetTrustee(u.ChairId, t)
	DataTrustee := &pk_ddz_msg.G2C_DDZ_TRUSTEE{}
	DataTrustee.TrusteeUser = u.ChairId
	DataTrustee.Trustee = t

	room.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		log.Debug("托管状态%v", DataTrustee)
		u.WriteMsg(DataTrustee)
	})

	if u.ChairId == room.CurrentUser {
		// 托管者是当前操作者
		room.checkNextUserTrustee()
	}
}

// 托管、明牌、放弃出牌
func (r *ddz_data_mgr) OtherOperation(args []interface{}) {
	nType := args[0].(string)
	u := args[1].(*user.User)
	switch nType {
	case "Trustee":
		t := args[2].(*pk_ddz_msg.C2G_DDZ_TRUSTEE)
		r.Trustee(u, t.Trustee)
		break
	case "ShowCard":
		r.ShowCard(u)
		break
	case "PassCard":
		r.PassCard(u)
		break
	}
}

// 启动操作定时器
func (r *ddz_data_mgr) startOperateCardTimer(nTime int) {
	if r.CardTimer != nil {
		r.CardTimer.Stop()
		r.CardTimer = nil
	}

	f := func() {
		u := r.PkBase.UserMgr.GetUserByChairId(r.CurrentUser)
		log.Debug("当前操作玩家ID%v", u)
		r.Trustee(u, true)
		r.resetOperateCardTimer(nTime)
	}

	r.CardTimer = time.AfterFunc(time.Duration(nTime+5)*time.Second, f)
}

// 重置定时器
func (r *ddz_data_mgr) resetOperateCardTimer(nTime int) {
	log.Debug("重置定时器时间%d", nTime)
	if r.CardTimer != nil {
		r.CardTimer.Reset(time.Duration(nTime+5) * time.Second)
	}
}

// 停止定时器
func (r *ddz_data_mgr) stopOperateCardTimer() {
	if r.CardTimer != nil {
		log.Debug("停止定时器")
		r.CardTimer.Stop()
		r.CardTimer = nil
	}
}
