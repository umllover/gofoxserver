package room

import (
	"mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/pk_ddz_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
	"github.com/mitchellh/mapstructure"
)

func NewDDZDataMgr(info *msg.L2G_CreatorRoom, uid int64, ConfigIdx int, name string, temp *base.GameServiceOption, base *DDZ_Entry) *ddz_data_mgr {
	d := new(ddz_data_mgr)
	d.RoomData = pk_base.NewDataMgr(info.RoomID, uid, ConfigIdx, name, temp, base.Entry_base, info)
	d.initParam()

	var setInfo pk_ddz_msg.C2G_DDZ_CreateRoomInfo
	err := mapstructure.Decode(info.OtherInfo, &setInfo)
	if err != nil {
		log.Error(" mapstructure.Decode error")
	}
	d.EightKing = setInfo.King
	d.GameType = setInfo.GameType

	return d
}

type ddz_data_mgr struct {
	*pk_base.RoomData
	GameStatus int // 当前游戏状态

	CurrentUser int // 当前玩家

	BankerUser int // 地主
	TurnWiner  int // 出牌玩家
	WinnerUser int // 上一个赢的玩家

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
	RecordOutCards []int                          // 出牌历史记录

	// 扑克信息
	MaxCardCount int     // 最大扑克数
	BankerCard   [3]int  // 游戏底牌
	HandCardData [][]int // 手上扑克
	ShowCardSign []bool  // 用户明牌标识
	RecordInfo   [][]int // 历史积分

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

	room.MaxCardCount = room.GetCfg().MaxRepertory
	if room.EightKing {
		room.MaxCardCount += 6
	}

	room.BankerCard = [3]int{}
	room.HandCardData = append([][]int{})
	room.RecordOutCards = []int{}

	room.ScoreTimes = 0
}

func (room *ddz_data_mgr) InitRoomOne() {

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

	// 发送明牌标识
	StatusFree.ShowCardSign = make([]bool, len(room.ShowCardSign))
	util.DeepCopy(&StatusFree.ShowCardSign, &room.ShowCardSign)

	u.WriteMsg(StatusFree)

	// 以下用扑克公告协议，到时候再打开
	/*
		StatusFree := &pk_common_msg.G2C_PKCOMMON_StatusFree{}

		StatusFree.CellScore = room.PkBase.Temp.Source // 基础积分
		StatusFree.GameRoomName = room.Name

		StatusFree.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()              // 出牌时间
		StatusFree.TimeOperateCard = room.GetCfg().CallScoreTime                    // 叫分时间
		StatusFree.TimeStartGame = int64(room.PkBase.TimerMgr.GetTimeOperateCard()) // 开始时间 	// 首出时间

		StatusFree.PlayMode = room.GameType
		StatusFree.CountLimit = room.PkBase.TimerMgr.GetMaxPlayCnt()

		StatusFree.PlayerCount = room.PkBase.TimerMgr.GetMaxPlayCnt()
		StatusFree.CurrentPlayCount = room.PkBase.TimerMgr.GetPlayCount()

		StatusFree.EightKing = room.EightKing
		// 发送明牌标识
		StatusFree.ShowCardSign = make([]bool, len(room.ShowCardSign))
		util.DeepCopy(&StatusFree.ShowCardSign, &room.ShowCardSign)

		u.WriteMsg(StatusFree)
	*/
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

	StatusCall.ScoreInfo = util.CopySlicInt(room.ScoreInfo)

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

	// 发送托管标识
	trustees := room.PkBase.UserMgr.GetTrustees()
	StatusCall.TrusteeSign = make([]bool, len(trustees))
	util.DeepCopy(&StatusCall.TrusteeSign, &trustees)

	log.Debug("叫分进行时%v", StatusCall)

	u.WriteMsg(StatusCall)

	// 以下用扑克公告协议，到时候再打开
	/*
		StatusCall := &pk_common_msg.G2C_PKCOMMON_StatusCall{}

		StatusCall.CallBanker = room.CurrentUser
		StatusCall.CellScore = room.PkBase.Temp.Source
		StatusCall.GameRoomName = room.Name

		StatusCall.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()
		StatusCall.TimeCallScore = room.PkBase.TimerMgr.GetTimeOperateCard()

		StatusCall.PlayMode = room.GameType
		StatusCall.WildCard = room.LiziCard
		StatusCall.EightKing = room.EightKing
		StatusCall.BankerScore = room.ScoreTimes
		StatusCall.ScoreInfo = util.CopySlicInt(room.ScoreInfo)
		StatusCall.HandCardCount = make([]int, len(room.HandCardData))
		for i := 0; i < len(room.HandCardData); i++ {
			StatusCall.HandCardCount[i] = len(room.HandCardData[i])
		}

		StatusCall.ScoreInfo = util.CopySlicInt(room.ScoreInfo)

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

		// 发送托管标识
		trustees := room.PkBase.UserMgr.GetTrustees()
		StatusCall.TrusteeSign = make([]bool, len(trustees))
		util.DeepCopy(&StatusCall.TrusteeSign, &trustees)

		log.Debug("叫分进行时%v", StatusCall)

		u.WriteMsg(StatusCall)
	*/
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
	// 发送托管标识
	trustees := room.PkBase.UserMgr.GetTrustees()
	StatusPlay.TrusteeSign = make([]bool, len(trustees))
	util.DeepCopy(&StatusPlay.TrusteeSign, &trustees)

	log.Debug("游戏进行时%v", StatusPlay)
	u.WriteMsg(StatusPlay)

	// 以下用扑克公告协议，到时候再打开
	/*
	   StatusPlay := &pk_common_msg.G2C_PKCOMMON_StatusPlay{}

	   StatusPlay.CellScore = room.PkBase.Temp.Source
	   StatusPlay.PlayerCount = room.PlayerCount
	   StatusPlay.BankerUser = room.BankerUser
	   StatusPlay.PublicCardData = util.CopySlicInt(room.BankerCard[:])
	   // 发送数据
	   for i := 0; i < room.PkBase.Temp.MaxPlayer; i++ {
	   	if room.ShowCardSign[i] || u.ChairId == i {
	   		StatusPlay.HandCardData = append(StatusPlay.HandCardData, room.HandCardData[i])
	   	} else {
	   		StatusPlay.HandCardData = append(StatusPlay.HandCardData, nil)
	   	}
	   }

	   StatusPlay.GameRoomName = room.Name

	   StatusPlay.CurrentPlayCount = room.PkBase.TimerMgr.GetPlayCount()
	   StatusPlay.LimitPlayCount = room.PkBase.TimerMgr.GetMaxPlayCnt()
	   StatusPlay.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()
	   StatusPlay.CurrentUser = room.CurrentUser
	   StatusPlay.EightKing = room.EightKing
	   StatusPlay.PlayMode = room.GameType
	   StatusPlay.WildCard = room.LiziCard
	   StatusPlay.BankerScore = room.ScoreTimes
	   StatusPlay.TurnUser = room.TurnWiner
	   if StatusPlay.TurnUser != cost.INVALID_CHAIR {
	   	turnCardData := room.TurnCardData[room.TurnWiner][len(room.TurnCardData[room.TurnWiner])-1]
	   	util.DeepCopy(&StatusPlay.TurnCardData, &turnCardData)
	   }

	   StatusPlay.HandCardCount = make([]int, room.PlayerCount)
	   for i := 0; i < room.PlayerCount; i++ {
	   	StatusPlay.HandCardCount[i] = len(room.HandCardData[i])
	   }
	   StatusPlay.EachBombCount = util.CopySlicInt(room.EachBombCount)
	   StatusPlay.KingCount = util.CopySlicInt(room.KingCount)

	   StatusPlay.ShowCardSign = make([]bool, len(room.ShowCardSign))
	   util.DeepCopy(&StatusPlay.ShowCardSign, &room.ShowCardSign)

	   // 发送托管标识
	   trustees := room.PkBase.UserMgr.GetTrustees()
	   StatusPlay.TrusteeSign = make([]bool, len(trustees))
	   util.DeepCopy(&StatusPlay.TrusteeSign, &trustees)

	   log.Debug("游戏进行时%v", StatusPlay)
	   u.WriteMsg(StatusPlay)
	*/
}

// 开始游戏，发送扑克
func (room *ddz_data_mgr) SendGameStart() {

	userMgr := room.PkBase.UserMgr

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, cost.US_PLAYING)
	})

	// 打乱牌
	hasCard := false
	if room.GameType == GAME_TYPE_HAPPY {
		hasCard = room.sendCardRuleOfHappyType()
	}

	if !hasCard {
		room.sendCardRuleOfNormal()
	}

	//构造变量
	GameStart := &pk_ddz_msg.G2C_DDZ_GameStart{}

	// 发完牌就选癞子
	if room.GameType == GAME_TYPE_LZ {
		// 随机选一张牌为癞子牌
		room.LiziCard = util.RandInterval(1, 13)
		GameStart.LiziCard = room.LiziCard
		room.PkBase.LogicMgr.SetParamToLogic(room.LiziCard)
	}

	// 取上一次赢的玩家
	if len(room.RecordInfo) > 0 {
		room.CurrentUser = room.WinnerUser
	}

	// 如果是第一把，取黑桃三作为叫分者
	for i := 0; i < room.PlayerCount; i++ {
		if room.CurrentUser != cost.INVALID_CHAIR {
			break
		}

		tmpCardData := room.HandCardData[i]
		for _, v := range tmpCardData {
			if v == 0x33 {
				room.CurrentUser = i
				break
			}
		}
	}

	// 黑桃三可能在底牌，随机选一个
	if room.CurrentUser == cost.INVALID_CHAIR {
		room.CurrentUser = util.RandInterval(0, 2)
	}

	// 初始化叫分信息
	for i := 0; i < room.PlayerCount; i++ {
		room.ScoreInfo[i] = CALLSCORE_NOCALL
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

// 欢乐场发牌规则
func (r *ddz_data_mgr) sendCardRuleOfHappyType() bool {
	// 欢乐场，从数据库里取
	var dataCard []int
	hasCard := false
	if r.EightKing {
		if HappyCardsKing.count > 0 {
			nIndex := util.RandInterval(0, HappyCardsKing.count-1)
			dataCard = util.CopySlicInt(HappyCardsKing.Cards[nIndex][:])
			hasCard = true
		}
	} else {
		if HappyCards.count > 0 {
			nIndex := util.RandInterval(0, HappyCards.count-1)
			dataCard = util.CopySlicInt(HappyCards.Cards[nIndex][:])
			hasCard = true
		}
	}
	if hasCard {
		log.Debug("当前牌%v", dataCard)
		nIndex := 0                                          // 当前索引
		var nCount int                                       // 每次随机取的条数
		var nMaxCount = (r.MaxCardCount - 3) / r.PlayerCount // 每个人牌的最大数
		// 把结尾三张当成底牌
		for i := 0; i < 3; i++ {
			r.BankerCard[i] = dataCard[r.MaxCardCount-i-1]
		}

		for nIndex < len(dataCard)-3 {
			for i := 0; i < r.PlayerCount; i++ {
				if i >= len(r.HandCardData) {
					r.HandCardData = append(r.HandCardData, []int{})
				}

				nSurplus := nMaxCount - len(r.HandCardData[i]) // 当前扑克还差多少
				if nSurplus <= 0 {
					continue
				}
				if nSurplus <= 6 {
					// 当前剩余牌数小于6，则直接取6张
					nCount = nSurplus
				} else if nSurplus <= 12 {
					// 6~12，随机取6~nSurplus
					nCount = util.RandInterval(6, nSurplus)
				} else {
					// 大于12，则随机取6~12
					nCount = util.RandInterval(6, 12)
				}

				for j := 0; j < nCount; j++ {
					r.HandCardData[i] = append(r.HandCardData[i], dataCard[nIndex+j])
				}
				nIndex += nCount
			}
		}
		for i := 0; i < r.PlayerCount; i++ {
			log.Debug("已分配完的扑克%v", r.HandCardData[i])
		}
	}
	return hasCard
}

// 普通发牌规则
func (r *ddz_data_mgr) sendCardRuleOfNormal() {
	var RepertoryCard []int = make([]int, r.MaxCardCount)
	r.PkBase.LogicMgr.RandCardList(RepertoryCard, pk_base.GetCardByIdx(r.ConfigIdx))
	log.Debug("随机打乱扑克后的牌%v", RepertoryCard)

	// 底牌
	for i := 0; i < 3; i++ {
		r.BankerCard[i] = RepertoryCard[r.MaxCardCount-i-1]
	}
	RepertoryCard = RepertoryCard[:len(RepertoryCard)-3]

	log.Debug("底牌%v", r.BankerCard)

	cardCount := (r.MaxCardCount - 3) / r.PlayerCount

	// 初始化牌
	r.HandCardData = append([][]int{})
	for i := 0; i < r.PlayerCount; i++ {
		tempCardData := util.CopySlicInt(RepertoryCard[len(RepertoryCard)-cardCount:])
		r.PkBase.LogicMgr.SortCardList(tempCardData, len(tempCardData))
		RepertoryCard = RepertoryCard[:len(RepertoryCard)-cardCount]
		r.HandCardData = append(r.HandCardData, tempCardData)
	}
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

	r.resetOperateCardTimer(r.GetCfg().CallScoreTime)

	if scoreTimes == 0 {

	} else {
		r.ScoreTimes = scoreTimes
		r.BankerUser = u.ChairId
	}

	r.ScoreInfo[u.ChairId] = scoreTimes

	nextCallUser := r.nextUser(r.CurrentUser) // 下一个叫分玩家

	isEnd := (r.ScoreTimes == CALLSCORE_MAX) || (r.ScoreInfo[nextCallUser] != CALLSCORE_NOCALL)

	if !isEnd {
		r.ScoreInfo[nextCallUser] = CALLSCORE_CALLING
		r.CurrentUser = nextCallUser
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
			r.BankerUser = nextCallUser
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

	// 用户出牌后重置定时器
	r.resetOperateCardTimer(r.PkBase.Temp.OutCardTime)
	if r.GameType == GAME_TYPE_CLASSIC {
		// 经典场，把出牌数据收集起来
		for i := 0; i < len(nowCard.CardData); i++ {
			r.RecordOutCards = append(r.RecordOutCards, nowCard.CardData[i])
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
		r.WinnerUser = u.ChairId
		r.PkBase.OnEventGameConclude(cost.GER_NORMAL)
		return
	}

	r.checkNextUserTrustee()
}

// 判断下一个玩家托管状态
func (r *ddz_data_mgr) checkNextUserTrustee() {
	log.Debug("当前游戏状态%d,玩家%d,上次出牌玩家%d,托管状态%v", r.GameStatus, r.CurrentUser, r.TurnWiner, r.PkBase.UserMgr.GetTrustees())
	if r.PkBase.UserMgr.IsTrustee(r.CurrentUser) {
		// 出牌玩家为托管状态
		if r.GameStatus == GAME_STATUS_CALL {
			// 叫分状态，托管则不叫
			r.CallScore(r.PkBase.UserMgr.GetUserByChairId(r.CurrentUser), 0)
		} else if r.GameStatus == GAME_STATUS_PLAY {
			// 出牌状态
			if r.CurrentUser == r.TurnWiner || r.TurnWiner == cost.INVALID_CHAIR {
				// 上一个出牌玩家是自己，则选最小牌
				nLen := len(r.HandCardData[r.CurrentUser])
				if nLen > 0 {
					var cardData []int
					cardData = append(cardData, r.HandCardData[r.CurrentUser][nLen-1])
					r.OpenCard(r.PkBase.UserMgr.GetUserByChairId(r.CurrentUser), 1, cardData)
				} else {
					// 作保护，如果手牌为0了，则游戏结束
					r.WinnerUser = r.CurrentUser
					r.PkBase.OnEventGameConclude(cost.GER_NORMAL)
				}

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

	r.resetOperateCardTimer(r.PkBase.Temp.OutCardTime)
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
func (r *ddz_data_mgr) NormalEnd(cbReason int) {
	log.Debug("游戏正常结束了")
	r.stopOperateCardTimer()
	if cbReason == 2 {
		r.WinnerUser = cost.INVALID_CHAIR
	}

	DataGameConclude := &pk_ddz_msg.G2C_DDZ_GameConclude{}
	DataGameConclude.CellScore = r.PkBase.Temp.Source

	DataGameConclude.Reason = cbReason

	if r.GameStatus == GAME_STATUS_FREE {

	} else if r.GameStatus == GAME_STATUS_CALL {
		DataGameConclude.BankerScore = r.ScoreTimes
		util.DeepCopy(&DataGameConclude.HandCardData, &r.HandCardData)
		DataGameConclude.GameScore = []int{0, 0, 0}
		r.RecordInfo = append(r.RecordInfo, util.CopySlicInt(DataGameConclude.GameScore))
	} else {
		DataGameConclude.BankerScore = r.ScoreTimes

		// 常规结束才需要算积分
		if cbReason == 0 {
			// 算分数
			nMultiple := r.ScoreTimes

			if r.BankerUser != cost.INVALID_CHAIR {
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

				// 地主明牌翻倍
				if r.ShowCardSign[r.BankerUser] == true {
					nMultiple <<= 1
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

			// 八王
			for _, v := range r.KingCount {
				if v == 8 {
					nMultiple *= 8 * 2
				} else if v >= 2 {
					nMultiple *= v
				}
			}

			// 计算积分
			gameScore := r.PkBase.Temp.Source * nMultiple

			if r.BankerUser != cost.INVALID_CHAIR {
				if len(r.HandCardData[r.BankerUser]) <= 0 {
					gameScore = 0 - gameScore
				}
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
			if r.BankerUser != cost.INVALID_CHAIR {
				DataGameConclude.GameScore[r.BankerUser] = 0 - score
			}

			// 经典场把历史出牌存数据库
			log.Debug("当前类型%d", r.GameType)
			if r.GameType == GAME_TYPE_CLASSIC {
				log.Debug("是经典场")
				r.saveOutCardToDB()
			}
		} else {
			DataGameConclude.GameScore = []int{0, 0, 0}
		}

		// 炸弹
		util.DeepCopy(&DataGameConclude.KingCount, &r.KingCount)
		util.DeepCopy(&DataGameConclude.EachBombCount, &r.EachBombCount)
		util.DeepCopy(&DataGameConclude.HandCardData, &r.HandCardData)

		// 服务端收集历史积分
		r.RecordInfo = append(r.RecordInfo, util.CopySlicInt(DataGameConclude.GameScore))

		if cbReason == 0 {
			//设置玩家积分
			r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
				var nScore int
				for _, arr := range r.RecordInfo {
					if arr != nil && len(arr) > 0 {
						nScore += arr[u.ChairId]
					}
				}
				u.Score = int64(nScore)
				log.Debug("当前玩家%d积分%d", u.ChairId, nScore)
			})
		}
	}
	r.GameStatus = GAME_STATUS_FREE

	log.Debug("历史积分%v", r.RecordInfo)
	if r.PkBase.TimerMgr.GetPlayCount() >= r.PkBase.TimerMgr.GetMaxPlayCnt() || cbReason > 0 {
		util.DeepCopy(&DataGameConclude.RecordInfo, &r.RecordInfo)
	}

	r.sendGameEndMsg(DataGameConclude)
}

// 保存出牌记录到数据库里
func (r *ddz_data_mgr) saveOutCardToDB() {
	// 检查出牌是否完整
	nMaxCardCount := r.GetCfg().MaxRepertory
	if r.EightKing {
		nMaxCardCount += 6
	}
	// 把未打完的牌插入
	log.Debug("剩余牌%v", r.HandCardData)
	for i := 0; i < r.PlayerCount; i++ {
		log.Debug("剩余牌%d", i)
		cardData := r.HandCardData[i]
		for _, v := range cardData {
			r.RecordOutCards = append(r.RecordOutCards, v)
		}
	}
	log.Debug("当前最大牌数%d-------%d", nMaxCardCount, r.RecordOutCards)
	if nMaxCardCount == len(r.RecordOutCards) {
		// 数量相等才能存数据库
		// ----存数据库----
		if r.EightKing {
			UpdatehappyKingCardList(r.RecordOutCards)
		} else {
			UpdateHappyCardList(r.RecordOutCards)
		}
	}
}

//解散房间结束
func (r *ddz_data_mgr) DismissEnd(cbReason int) {
	r.stopOperateCardTimer()
	r.WinnerUser = cost.INVALID_CHAIR
	DataGameConclude := &pk_ddz_msg.G2C_DDZ_GameConclude{}
	DataGameConclude.CellScore = r.PkBase.Temp.Source
	DataGameConclude.GameScore = make([]int, r.PlayerCount)
	DataGameConclude.Reason = cbReason
	// 炸弹
	util.DeepCopy(&DataGameConclude.EachBombCount, &r.EachBombCount)

	DataGameConclude.BankerScore = r.ScoreTimes
	util.DeepCopy(&DataGameConclude.HandCardData, &r.HandCardData)
	util.DeepCopy(&DataGameConclude.RecordInfo, &r.RecordInfo)
	log.Debug("解散房间历史积分%v", r.RecordInfo)

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

	r.CardTimer = time.AfterFunc(time.Duration(nTime)*time.Second, f)
}

// 重置定时器
func (r *ddz_data_mgr) resetOperateCardTimer(nTime int) {
	log.Debug("重置定时器时间%d", nTime)
	if r.CardTimer != nil {
		r.CardTimer.Reset(time.Duration(nTime) * time.Second)
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
