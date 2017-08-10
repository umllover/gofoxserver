package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"

	"github.com/lovelly/leaf/timer"

	"mj/common/cost"
	"mj/common/msg/nn_tb_msg"
	"mj/common/msg/pk_common_msg"
	"mj/gameServer/user"
	"time"

	"mj/common/msg"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

// 游戏状态
const (
	//GAME_STATUS_NULL = 1000 // 空
	//PLAYER_ENTER_ROOM  	= 1001 // 玩家进入房间
	GAME_STATUS_START          = 1002 // 游戏开始
	GAME_STATUS_CALL_SCORE     = 1003 // 抢庄
	GAME_STATUS_ADD_SCORE      = 1004 // 加注
	GAME_STATUS_SEND_LAST_CARD = 1005 // 发最后一张牌
	GAME_STATUS_OPEN_CARD      = 1006 // 亮牌
	GAME_STATUS_CAL_SCORE      = 1007 // 结算
)

// 定时器 -- for test
const (
	TIME_CALL_SCORE = 3000
	TIME_ADD_SCORE  = 3000
	TIME_OPEN_CARD  = 3000
)

func NewDataMgr(id int, uid int64, ConfigIdx int, name string, temp *base.GameServiceOption, base *NNTB_Entry, info *msg.L2G_CreatorRoom) *nntb_data_mgr {
	d := new(nntb_data_mgr)
	d.RoomData = pk_base.NewDataMgr(id, uid, ConfigIdx, name, temp, base.Entry_base, info)
	return d
}

// 亮牌信息
type OpenCardInfo struct {
	CardData []int // 亮牌数据
	CardType int   //亮牌牌型
}

type nntb_data_mgr struct {
	*pk_base.RoomData

	//游戏变量
	CardData       [][]int //用户扑克
	PublicCardData []int   //公共牌 两张
	RepertoryCard  []int   //库存扑克
	LeftCardCount  int     //库存剩余扑克数量

	OpenCardMap       map[int]OpenCardInfo //记录亮牌数据
	CallScoreTimesMap map[int]int          //记录叫分信息
	CalScoreMap       map[int]int          //记录算分
	AddScoreTimesMap  map[int]int          //记录用户加注信息
	UserGameStatusMap map[int]int          //记录用户游戏状态信息 用于断线重连

	BankerUser int //庄家用户

	// 游戏状态
	GameStatus     int
	CallScoreTimer *timer.Timer
	AddScoreTimer  *timer.Timer
	OpenCardTimer  *timer.Timer
}

func (room *nntb_data_mgr) SetUserAddScoreTimes(chairId int, addScoreTimes int) {
	room.AddScoreTimesMap[chairId] = addScoreTimes
}

func (room *nntb_data_mgr) SetUserCallScoreTimes(chairId int, callScoreTimes int) {
	room.CallScoreTimesMap[chairId] = callScoreTimes
}

func (room *nntb_data_mgr) SendStatusReady(u *user.User) {
	StatusFree := &pk_common_msg.G2C_PKCOMMON_StatusFree{}

	StatusFree.CellScore = room.PkBase.Temp.Source                         //基础积分
	StatusFree.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()         //出牌时间
	StatusFree.TimeOperateCard = room.PkBase.TimerMgr.GetTimeOperateCard() //操作时间
	StatusFree.TimeStartGame = room.PkBase.TimerMgr.GetCreatrTime()        //开始时间
	StatusFree.TurnScore = make([]int, room.PkBase.TimerMgr.GetMaxPlayCnt())
	StatusFree.CollectScore = make([]int, room.PlayerCount)
	StatusFree.EachRoundScore = make([][]int, room.PlayerCount, room.PkBase.TimerMgr.GetMaxPlayCnt())
	StatusFree.InitScore = make([]int, room.PlayerCount)
	log.Debug("at send status ready %v", room.InitScoreMap)

	for i := 0; i < room.PlayerCount; i++ {
		StatusFree.InitScore[i] = room.InitScoreMap[i]
	}

	StatusFree.PlayerCount = room.PlayerCount
	StatusFree.CurrentPlayCount = room.PkBase.TimerMgr.GetPlayCount()
	StatusFree.CountLimit = room.PkBase.TimerMgr.GetMaxPlayCnt() //局数限制
	StatusFree.GameRoomName = room.Name

	u.WriteMsg(StatusFree)
}

func (room *nntb_data_mgr) SendStatusPlay(u *user.User) {
	StatusPlay := &pk_common_msg.G2C_PKCOMMON_StatusPlay{}

	log.Debug("at send status play, player count%d", room.PlayerCount)
	//游戏变量
	StatusPlay.CellScore = room.CellScore
	StatusPlay.GameStatus = room.GameStatus
	if room.GameStatus > GAME_STATUS_CALL_SCORE {
		StatusPlay.BankerUser = room.BankerUser
	}
	StatusPlay.PlayerCount = room.PlayerCount
	StatusPlay.GameRoomName = room.Name
	StatusPlay.CurrentPlayCount = room.PkBase.TimerMgr.GetPlayCount()
	StatusPlay.LimitPlayCount = room.PkBase.TimerMgr.GetMaxPlayCnt()

	userMgr := room.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		userReLoginInfo := new(pk_common_msg.UserReLoginInfo)
		userReLoginInfo.ChairID = u.ChairId
		userReLoginInfo.UserGameStatus = room.UserGameStatusMap[u.ChairId]
		userReLoginInfo.CallScoreTimes = room.CallScoreTimesMap[u.ChairId]
		userReLoginInfo.AddScoreTimes = room.AddScoreTimesMap[u.ChairId]
		StatusPlay.UserReLoginInfos = append(StatusPlay.UserReLoginInfos, userReLoginInfo)
	})

	StatusPlay.PublicCardData = append(StatusPlay.PublicCardData, room.PublicCardData...)

	StatusPlay.HandCardData = make([][]int, room.PlayerCount)
	for i := 0; i < room.PlayerCount; i++ {
		StatusPlay.HandCardData[i] = append(StatusPlay.HandCardData[i], room.CardData[i]...)
	}

	for i := 0; i < room.PlayerCount; i++ {
		StatusPlay.InitScore = append(StatusPlay.InitScore, room.InitScoreMap[i])
	}

	u.WriteMsg(StatusPlay)
}

func (room *nntb_data_mgr) BeforeStartGame(UserCnt int) {
	log.Debug("init room")
	room.InitRoom(UserCnt)
	room.GameStatus = GAME_STATUS_START
	room.SetAllUserGameStatus(GAME_STATUS_START)
}

func (room *nntb_data_mgr) StartGameing() {
	// 发牌
	log.Debug("dispatch card")
	room.StartDispatchCard()
}

func (room *nntb_data_mgr) AfterStartGame() {
	// 叫分
	log.Debug("call score")
	room.GameStatus = GAME_STATUS_CALL_SCORE
	log.Debug("begin call score timer")

	room.CallScoreTimer = room.PkBase.AfterFunc(TIME_CALL_SCORE*time.Second, func() {
		log.Debug("end call score timer")
		if room.GameStatus == GAME_STATUS_CALL_SCORE { // 超时叫分结束
			room.CallScoreEnd()
		}
	})
}

func (room *nntb_data_mgr) InitRoom(UserCnt int) {

	log.Debug("nn init room version 28001 player count %d", UserCnt)
	//初始化
	room.CardData = make([][]int, UserCnt)

	for i := 0; i < UserCnt; i++ {
		room.CardData[i] = make([]int, room.GetCfg().MaxCount)
	}
	room.PublicCardData = make([]int, room.GetCfg().PublicCardCount)
	room.LeftCardCount = room.GetCfg().MaxRepertory

	room.PlayerCount = UserCnt
	room.CellScore = room.PkBase.Temp.Source

	room.CallScoreTimesMap = make(map[int]int)
	room.AddScoreTimesMap = make(map[int]int)

	room.OpenCardMap = make(map[int]OpenCardInfo)
	room.CalScoreMap = make(map[int]int)
	room.UserGameStatusMap = make(map[int]int)
	room.RepertoryCard = make([]int, room.GetCfg().MaxRepertory)

	room.ExitScore = 0
	room.DynamicScore = 0
	room.BankerUser = 0
	room.FisrtCallUser = cost.INVALID_CHAIR
	room.CurrentUser = cost.INVALID_CHAIR

	userMgr := room.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(&nn_tb_msg.G2C_TBNN_GameStart{
			PlayerCount: room.PlayerCount,
		})
	})

}

func (r *nntb_data_mgr) GetOneCard() int { // 从牌堆取出一张
	r.LeftCardCount--
	return r.RepertoryCard[r.LeftCardCount]
}

func (room *nntb_data_mgr) StartDispatchCard() {

	userMgr := room.PkBase.UserMgr
	gameLogic := room.PkBase.LogicMgr

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, cost.US_PLAYING)
	})

	gameLogic.RandCardList(room.RepertoryCard, pk_base.GetTBNNCards())

	//分发扑克
	// 两张公共牌
	for i := 0; i < room.GetCfg().PublicCardCount; i++ {
		room.PublicCardData[i] = room.GetOneCard()
		log.Debug("public card %d", room.PublicCardData[i])
	}

	PublicCardData := &nn_tb_msg.G2C_TBNN_PublicCard{}
	util.DeepCopy(&PublicCardData.PublicCardData, &room.PublicCardData)
	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(PublicCardData)
	})

	// 再发每个用户4张牌
	userMgr.ForEachUser(func(u *user.User) {
		for i := 0; i < room.GetCfg().MaxCount-1; i++ {
			room.CardData[u.ChairId][i] = room.GetOneCard()
		}
	})

	// 把整幅牌发出去
	usersCardData := &nn_tb_msg.G2C_TBNN_SendCard{}
	usersCardData.CardData = make([][]int, room.PlayerCount)
	userMgr.ForEachUser(func(u *user.User) {
		usersCardData.CardData[u.ChairId] = make([]int, room.GetCfg().MaxCount)
	})
	util.DeepCopy(&usersCardData.CardData, &room.CardData)

	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(usersCardData)
	})

	return
}

//正常结束房间
func (room *nntb_data_mgr) NormalEnd(cbReason int) {

	userMgr := room.PkBase.UserMgr

	calScore := &nn_tb_msg.G2C_TBNN_CalScore{}
	calScore.GameTax = make([]int, room.PlayerCount)
	calScore.CardType = make([]int, room.PlayerCount)
	calScore.GameScore = make([]int, room.PlayerCount)
	calScore.CardData = make([][]int, room.PlayerCount)

	for i := 0; i < room.PlayerCount; i++ {
		calScore.CardData[i] = make([]int, pk_base.GetCfg(pk_base.IDX_TBNN).MaxCount)
	}
	calScore.InitScore = make([]int, room.PlayerCount)

	userMgr.ForEachUser(func(u *user.User) {
		calScore.GameScore[u.ChairId] = room.CalScoreMap[u.ChairId]
		openCardInfo := OpenCardInfo{
			CardType: room.OpenCardMap[u.ChairId].CardType,
			CardData: room.OpenCardMap[u.ChairId].CardData,
		}
		calScore.CardType[u.ChairId] = openCardInfo.CardType
		util.DeepCopy(&calScore.CardData[u.ChairId], &openCardInfo.CardData)
		// 更新积分
		room.InitScoreMap[u.ChairId] += room.CalScoreMap[u.ChairId]
	})

	log.Debug("normal end init score map %v", room.InitScoreMap)

	for i := 0; i < room.PlayerCount; i++ {
		calScore.InitScore[i] = room.InitScoreMap[i]
	}

	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(calScore)
	})

	roundScore := make([]int, room.PlayerCount)
	for i := 0; i < room.PlayerCount; i++ {
		roundScore[i] = room.CalScoreMap[i]
	}
	room.EachRoundScoreMap[room.PkBase.TimerMgr.GetPlayCount()] = roundScore
	log.Debug("normal end each round score map %v", room.EachRoundScoreMap)

	gameEnd := &nn_tb_msg.G2C_TBNN_GameEnd{}
	gameEnd.CurrentPlayCount = room.PkBase.TimerMgr.GetPlayCount()
	gameEnd.LimitPlayCount = room.PkBase.TimerMgr.GetMaxPlayCnt()
	gameEnd.EachRoundScore = make([][]int, gameEnd.CurrentPlayCount)
	gameEnd.Reason = cbReason
	for i := 0; i < gameEnd.CurrentPlayCount; i++ {
		gameEnd.EachRoundScore[i] = append(gameEnd.EachRoundScore[i], room.EachRoundScoreMap[i+1]...)
	}
	for i := 0; i < room.PlayerCount; i++ {
		gameEnd.InitScore = append(gameEnd.InitScore, room.InitScoreMap[i])
	}
	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(gameEnd)
	})

}

func (room *nntb_data_mgr) SetAllUserGameStatus(gameStatus int) {
	userMgr := room.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		room.UserGameStatusMap[u.ChairId] = gameStatus
	})
	log.Debug("set all user game status %v", room.UserGameStatusMap)
}

func (room *nntb_data_mgr) SetUserGameStatus(chairId int, gameStatus int) {
	room.UserGameStatusMap[chairId] = gameStatus
	log.Debug("set user game status chair id %d, user game status %d", chairId, room.UserGameStatusMap[chairId])
}

//解散接触
func (room *nntb_data_mgr) DismissEnd(cbReason int) {

}

// 用户叫分(抢庄)
func (r *nntb_data_mgr) CallScore(u *user.User, scoreTimes int) {
	if r.GameStatus != GAME_STATUS_CALL_SCORE {
		return
	}
	log.Debug("call score times userChairId:%d, scoretimes:%d", u.ChairId, scoreTimes)

	r.SetUserCallScoreTimes(u.ChairId, scoreTimes)
	r.SetUserGameStatus(u.ChairId, GAME_STATUS_CALL_SCORE)

	maxScoreTimes := 0
	for chairId, s := range r.CallScoreTimesMap {
		if s > maxScoreTimes {
			maxScoreTimes = s
			r.BankerUser = chairId
		}
	}
	r.ScoreTimes = maxScoreTimes

	// 广播叫分
	callScore := &nn_tb_msg.G2C_TBNN_CallScore{}
	callScore.ChairID = u.ChairId
	callScore.CallScore = scoreTimes
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u1 *user.User) {
		u1.WriteMsg(callScore)
	})

	if len(r.CallScoreTimesMap) == r.PlayerCount {
		//叫分结束
		r.CallScoreEnd()
	}
}

// 判定是否有人叫分

func (r *nntb_data_mgr) IsAnyOneCallScore() bool {
	for _, s := range r.CallScoreTimesMap {
		if s > 0 {
			return true
		}
	}
	return false
}

// 叫分结束

func (r *nntb_data_mgr) CallScoreEnd() {
	r.CallScoreTimer.Stop()
	r.SetAllUserGameStatus(GAME_STATUS_CALL_SCORE)
	log.Debug("call score end")
	// 发回叫分结果
	userMgr := r.PkBase.UserMgr
	//如果没有任何人叫分
	if !r.IsAnyOneCallScore() {
		r.BankerUser = 0
		r.ScoreTimes = 1
	}

	callScoreEnd := &nn_tb_msg.G2C_TBNN_CallScoreEnd{
		Banker:     r.BankerUser,
		ScoreTimes: r.ScoreTimes,
	}

	userMgr.ForEachUser(func(u *user.User) {
		if r.CallScoreTimesMap[u.ChairId] == r.ScoreTimes {
			callScoreEnd.ScoreTimesUser = append(callScoreEnd.ScoreTimesUser, u.ChairId)
		}
	})

	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(callScoreEnd)
	})

	// 进入加注
	r.GameStatus = GAME_STATUS_ADD_SCORE
	log.Debug("begin add score timer ")
	r.AddScoreTimer = r.PkBase.AfterFunc(TIME_ADD_SCORE*time.Second, func() { // 超时加注结束
		log.Debug("end add score timer")
		if r.GameStatus == GAME_STATUS_ADD_SCORE {
			r.AddScoreEnd()
		}
	})
}

// 用户加注
func (r *nntb_data_mgr) AddScore(u *user.User, score int) {

	if r.GameStatus != GAME_STATUS_ADD_SCORE {
		return
	}

	log.Debug("add score userChairId:%d, score:%d", u.ChairId, score)
	r.SetUserAddScoreTimes(u.ChairId, score)
	r.SetUserGameStatus(u.ChairId, GAME_STATUS_ADD_SCORE)

	// 广播加注
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(uFunc *user.User) {

		addScore := &nn_tb_msg.G2C_TBNN_AddScore{}
		addScore.ChairID = u.ChairId
		addScore.AddScoreCount = score
		uFunc.WriteMsg(addScore)
	})

	if len(r.AddScoreTimesMap) == r.PlayerCount-1 { //全加过加注结束 庄家不能加注
		r.AddScoreEnd()
	}
}

// 加注结束
func (r *nntb_data_mgr) AddScoreEnd() {
	r.AddScoreTimer.Stop()
	r.SetAllUserGameStatus(GAME_STATUS_ADD_SCORE)
	log.Debug("add score end")

	userMgr := r.PkBase.UserMgr

	userMgr.ForEachUser(func(u *user.User) {
		// 没有加过的默认一倍
		if r.AddScoreTimesMap[u.ChairId] <= 0 {
			r.SetUserAddScoreTimes(u.ChairId, 1)
		}
	})

	log.Debug("add score end  map %v", r.AddScoreTimesMap)

	// 进入最后一张牌
	log.Debug("enter send last card")
	r.GameStatus = GAME_STATUS_SEND_LAST_CARD
	r.SetAllUserGameStatus(GAME_STATUS_SEND_LAST_CARD)

	// 发最后一张牌
	userMgr.ForEachUser(func(u *user.User) {
		lastCard := r.GetOneCard()
		r.CardData[u.ChairId][r.GetCfg().MaxCount-1] = lastCard
	})

	lastCardData := &nn_tb_msg.G2C_TBNN_LastCard{}
	lastCardData.LastCard = make([][]int, r.PlayerCount)
	userMgr.ForEachUser(func(u *user.User) {
		lastCardData.LastCard[u.ChairId] = make([]int, r.GetCfg().MaxCount)
	})
	util.DeepCopy(&lastCardData.LastCard, &r.CardData)

	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(lastCardData)
	})

	// 进入亮牌
	r.GameStatus = GAME_STATUS_OPEN_CARD
	r.EnterOpenCard()

}

// 进入亮牌

func (r *nntb_data_mgr) EnterOpenCard() {
	log.Debug("enter open card")
	// 亮牌超时
	log.Debug("begin open card timer")

	r.OpenCardTimer = r.PkBase.AfterFunc(TIME_OPEN_CARD*time.Second, func() { // 超时亮牌结束
		log.Debug("end open card timer")
		if r.GameStatus == GAME_STATUS_OPEN_CARD {
			// 没有亮牌的用户自动亮牌
			userMgr := r.PkBase.UserMgr
			userMgr.ForEachUser(func(u *user.User) {
				if r.OpenCardMap[u.ChairId].CardData == nil {
					log.Debug("user : %d has not open card ", u.ChairId)
					// 需要改进
					//r.OpenCard(u,0, r.CardData[u.ChairId] )
					cardData := make([]int, 5)
					util.DeepCopy(&cardData, &r.CardData[u.ChairId])
					cardData = append(cardData, r.PublicCardData...)
					log.Debug("7cards:%v", cardData)
					dstCardData, cardType := r.SelectCard(cardData)

					if dstCardData != nil {
						r.OpenCard(u, cardType, dstCardData)
					}
				}
			})
		}
	})
}

// 验证

func (r *nntb_data_mgr) IsValidCard(chairID int, card int) bool {
	// 先验证是不是在公共牌中

	for i := 0; i < pk_base.GetCfg(pk_base.IDX_TBNN).PublicCardCount; i++ {
		if card == r.PublicCardData[i] {
			return true
		}
	}
	// 是不是在用户手牌

	for i := 0; i < pk_base.GetCfg(pk_base.IDX_TBNN).MaxCount; i++ {
		if card == r.CardData[chairID][i] {
			return true
		}
	}
	return false
}

func (r *nntb_data_mgr) IsValidCardData(chairID int, cardData []int) bool {
	for i := 0; i < len(cardData); i++ {
		if !r.IsValidCard(chairID, cardData[i]) {
			return false
		}
	}
	return true
}

// 亮牌

func (r *nntb_data_mgr) OpenCard(u *user.User, cardType int, cardData []int) {
	if r.GameStatus != GAME_STATUS_OPEN_CARD {
		return
	}
	log.Debug("user: %d open card type: %d card data : %v", u.ChairId, cardType, cardData)
	// 验证牌数据
	if !r.IsValidCardData(u.ChairId, cardData) {
		log.Debug("user: %d open card  data invalid ", u.ChairId)
		return
	}

	// 验证牌型
	if r.PkBase.LogicMgr.GetCardType(cardData) != cardType {
		cardType = r.PkBase.LogicMgr.GetCardType(cardData)
		log.Debug("user: %d open card type invalid , correct type is :%d",
			u.ChairId, r.PkBase.LogicMgr.GetCardType(cardData))
	}

	openCardInfo := OpenCardInfo{
		CardData: cardData,
		CardType: cardType,
	}
	log.Debug("open card info %v", openCardInfo)

	r.OpenCardMap[u.ChairId] = openCardInfo
	r.SetUserGameStatus(u.ChairId, GAME_STATUS_OPEN_CARD)

	// 广播亮牌
	userMgr := r.PkBase.UserMgr
	openCard := &nn_tb_msg.G2C_TBNN_Open_Card{}
	openCard.CardData = make([]int, r.GetCfg().MaxCount)
	util.DeepCopy(&openCard.CardData, &cardData)
	openCard.CardType = cardType
	openCard.ChairID = u.ChairId

	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(openCard)
	})

	if len(r.OpenCardMap) == r.PlayerCount { // 全亮过
		r.OpenCardEnd() // 亮牌结束
	}
}

// 亮牌结束
func (r *nntb_data_mgr) OpenCardEnd() {
	r.OpenCardTimer.Stop()
	r.SetAllUserGameStatus(GAME_STATUS_OPEN_CARD)
	// 结算
	// 比牌
	log.Debug("enter cal score")
	r.GameStatus = GAME_STATUS_CAL_SCORE
	r.SetAllUserGameStatus(GAME_STATUS_CAL_SCORE)
	logicMgr := r.PkBase.LogicMgr
	userMgr := r.PkBase.UserMgr

	userMgr.ForEachUser(func(u *user.User) {
		if u.ChairId != r.BankerUser { // 闲家与庄家比
			if logicMgr.CompareCard(r.OpenCardMap[r.BankerUser].CardData, r.OpenCardMap[u.ChairId].CardData) { // 庄家比闲家大
				log.Debug("at open card end %d %d %d %d ",
					r.CellScore, r.ScoreTimes, r.AddScoreTimesMap[u.ChairId], r.PkBase.LogicMgr.GetCardTimes(r.OpenCardMap[r.BankerUser].CardType))
				r.CalScoreMap[r.BankerUser] += r.CellScore * r.ScoreTimes * r.AddScoreTimesMap[u.ChairId] *
					r.PkBase.LogicMgr.GetCardTimes(r.OpenCardMap[r.BankerUser].CardType)

				r.CalScoreMap[u.ChairId] -= r.CellScore * r.ScoreTimes * r.AddScoreTimesMap[u.ChairId] *
					r.PkBase.LogicMgr.GetCardTimes(r.OpenCardMap[r.BankerUser].CardType)
				log.Debug("banker win  : banker card: %v banker score:%d, player card: %v player score:%d",
					r.OpenCardMap[r.BankerUser], r.CalScoreMap[r.BankerUser], r.OpenCardMap[u.ChairId], r.CalScoreMap[u.ChairId])

			} else {

				log.Debug("at open card end %d %d %d %d ",
					r.CellScore, r.ScoreTimes, r.AddScoreTimesMap[u.ChairId], r.PkBase.LogicMgr.GetCardTimes(r.OpenCardMap[r.BankerUser].CardType))
				r.CalScoreMap[r.BankerUser] -= r.CellScore * r.ScoreTimes * r.AddScoreTimesMap[u.ChairId] *
					r.PkBase.LogicMgr.GetCardTimes(r.OpenCardMap[u.ChairId].CardType)

				r.CalScoreMap[u.ChairId] += r.CellScore * r.ScoreTimes * r.AddScoreTimesMap[u.ChairId] *
					r.PkBase.LogicMgr.GetCardTimes(r.OpenCardMap[u.ChairId].CardType)
				log.Debug("banker lost  : banker card: %v banker score:%d, player card: %v player score:%d",
					r.OpenCardMap[r.BankerUser], r.CalScoreMap[r.BankerUser], r.OpenCardMap[u.ChairId], r.CalScoreMap[u.ChairId])
			}
		}
	})

	log.Debug("cal score map %v", r.CalScoreMap)

	// 游戏结束

	r.PkBase.OnEventGameConclude(0, userMgr.GetUserByChairId(0), cost.GER_NORMAL)

	/*r.PkBase.AfterFunc( 15 * time.Second, func() {
		log.Debug("game end timer")
		//退出房间
		userMgr.ForEachUser(func(u *user.User) {
			userMgr.LeaveRoom(u, r.PkBase.Status)
		})
		r.PkBase.Destroy(r.PkBase.DataMgr.GetRoomId())
	})*/

}

// 7选5
func (r *nntb_data_mgr) SelectCard(cardData []int) ([]int, int) {
	cardCount := len(cardData)

	if cardCount < 5 {
		return nil, 0
	}
	r.PkBase.LogicMgr.SortCardList(cardData, cardCount)
	var cardsMap = make(map[int][]int)

	index := 0

	for i := 0; i < cardCount-4; i++ {
		for j := i + 1; j < cardCount-3; j++ {
			for k := j + 1; k < cardCount-2; k++ {
				for m := k + 1; m < cardCount-1; m++ {
					for n := m + 1; n < cardCount; n++ {

						temp := []int{cardData[i], cardData[j], cardData[k], cardData[m], cardData[n]}
						cardsMap[index] = temp
						index++
					}
				}
			}
		}
	}

	// 按照牌型来选

	for cardType := 18; cardType >= 0; cardType-- {
		for i := 0; i < len(cardsMap); i++ {
			if r.PkBase.LogicMgr.GetCardType(cardsMap[i]) == cardType {
				return cardsMap[i], cardType
			}
		}
	}

	return nil, 0
}
