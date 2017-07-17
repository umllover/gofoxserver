package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"

	"github.com/lovelly/leaf/timer"

	"mj/common/cost"
	"mj/common/msg/nn_tb_msg"
	"mj/gameServer/user"
	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

// 游戏状态
const (
	GAME_NULL = 1000 // 空
	//PLAYER_ENTER_ROOM  	= 1001 // 玩家进入房间
	GAME_START       = 1002 // 游戏开始
	CALL_SCORE_TIMES = 1003 // 抢庄
	ADD_SCORE        = 1004 // 加注
	SEND_LAST_CARD   = 1005 // 发最后一张牌
	OPEN_CARD        = 1006 // 亮牌
	CAL_SCORE        = 1007 // 结算
)

// 定时器 -- for test
const (
	CALL_SCORE_TIME = 10
	ADD_SCORE_TIME  = 10
	OPEN_CARD_TIME  = 30
)

func NewDataMgr(id int, uid int64, ConfigIdx int, name string, temp *base.GameServiceOption, base *NNTB_Entry) *nntb_data_mgr {
	d := new(nntb_data_mgr)
	d.RoomData = pk_base.NewDataMgr(id, uid, ConfigIdx, name, temp, base.Entry_base)
	return d
}

type nntb_data_mgr struct {
	*pk_base.RoomData

	//游戏变量
	CardData          [][]int              //用户扑克
	PublicCardData    []int                //公共牌 两张
	RepertoryCard     []int                //库存扑克
	LeftCardCount     int                  //库存剩余扑克数量
	OpenCardMap       map[*user.User][]int //亮牌数据
	CallScoreTimesMap map[int]*user.User   //记录叫分信息
	CalScoreMap       map[*user.User]int   //计分信息
	AddScoreMap       map[*user.User]int   //记录用户加注信息

	BankerUser *user.User //庄家用户

	// 游戏状态
	GameStatus     int
	CallScoreTimer *timer.Timer
	AddScoreTimer  *timer.Timer
	OpenCardTimer  *timer.Timer
}

func (room *nntb_data_mgr) SendStatusReady(u *user.User) {
	StatusFree := &nn_tb_msg.G2C_TBNN_StatusFree{}

	StatusFree.CellScore = room.CellScore                                  //基础积分
	StatusFree.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()         //出牌时间
	StatusFree.TimeOperateCard = room.PkBase.TimerMgr.GetTimeOperateCard() //操作时间
	StatusFree.TimeStartGame = room.PkBase.TimerMgr.GetCreatrTime()        //开始时间
	for _, v := range room.HistoryScores {
		StatusFree.TurnScore = append(StatusFree.TurnScore, v.TurnScore)
		StatusFree.CollectScore = append(StatusFree.TurnScore, v.CollectScore)
	}
	StatusFree.PlayerCount = room.PkBase.TimerMgr.GetPlayCount() //玩家人数
	StatusFree.CountLimit = room.PkBase.TimerMgr.GetMaxPayCnt()  //局数限制
	StatusFree.GameRoomName = room.Name

	u.WriteMsg(StatusFree)
}

func (room *nntb_data_mgr) SendStatusPlay(u *user.User) {
	StatusPlay := &nn_tb_msg.G2C_TBNN_StatusPlay{}

	UserCnt := room.PkBase.UserMgr.GetMaxPlayerCnt()
	//游戏变量
	StatusPlay.BankerUser = room.BankerUser.ChairId
	StatusPlay.CellScore = room.CellScore

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

	u.WriteMsg(StatusPlay)
}

func (room *nntb_data_mgr) BeforeStartGame(UserCnt int) {
	room.GameStatus = GAME_START
	room.InitRoom(UserCnt)
}

func (room *nntb_data_mgr) StartGameing() {
	room.StartDispatchCard()
}

func (room *nntb_data_mgr) AfterStartGame() {
	room.GameStatus = CALL_SCORE_TIMES
	room.CallScoreTimer = room.PkBase.AfterFunc(CALL_SCORE_TIME*time.Second, func() {
		if room.GameStatus != ADD_SCORE { // 超时叫分结束
			room.CallScoreEnd()
		}
	})

	room.CallScoreTimer.Stop()

}

func (room *nntb_data_mgr) InitRoom(UserCnt int) {
	//初始化
	room.CardData = make([][]int, UserCnt)
	room.PublicCardData = make([]int, room.GetCfg().PublicCardCount)

	room.PlayerCount = UserCnt

	room.CallScoreTimesMap = make(map[int]*user.User)
	room.AddScoreMap = make(map[*user.User]int)
	room.OpenCardMap = make(map[*user.User][]int)
	room.RepertoryCard = make([]int, room.GetCfg().MaxRepertory)

	room.ExitScore = 0
	room.DynamicScore = 0
	room.BankerUser = nil
	room.FisrtCallUser = cost.INVALID_CHAIR
	room.CurrentUser = cost.INVALID_CHAIR

}

func (r *nntb_data_mgr) GetOneCard() int { // 从牌堆取出一张
	r.LeftCardCount -= 1
	return r.RepertoryCard[r.LeftCardCount]
}

func (room *nntb_data_mgr) StartDispatchCard() {
	log.Debug("start dispatch card")

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
		log.Debug("public card %d", room.CardData[i])
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

	userIndex := 0
	userMgr.ForEachUser(func(u *user.User) {
		UserCardData := &nn_tb_msg.G2C_TBNN_SendCard{}
		util.DeepCopy(&UserCardData.CardData, &room.CardData[userIndex])
		u.WriteMsg(UserCardData)
		userIndex++
	})

	return
}

func (room *nntb_data_mgr) SendGameStart() {
	//构造变量
	/*GameStart := &nn_tb_msg.G2C_TBNN_GameStart{}
	GameStart.BankerUser = room.BankerUser
	/*
	GameStart.SiceCount = room.SiceCount
	GameStart.HeapHead = room.HeapHead
	GameStart.HeapTail = room.HeapTail
	GameStart.MagicIndex = room.MjBase.LogicMgr.GetMagicIndex()
	GameStart.HeapCardInfo = room.HeapCardInfo
	GameStart.CardData = make([]int, MAX_COUNT)
	//发送数据
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		GameStart.UserAction = room.UserAction[u.ChairId]
		GameStart.CardData = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[u.ChairId])
		u.WriteMsg(GameStart)
	})*/
	log.Debug("startgame ... ")

}

//正常结束房间
func (room *nntb_data_mgr) NormalEnd() {

	userMgr := room.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		calScore := &nn_tb_msg.G2C_TBNN_CalScore{}
		calScore.GameScore = room.CalScoreMap[u]
		calScore.CardData = make([]int, pk_base.GetCfg(pk_base.IDX_TBNN).MaxCount)
		util.DeepCopy(calScore.CardData, room.OpenCardMap[u])

		u.WriteMsg(calScore)

		//历史积分
		if room.HistoryScores[u.ChairId] == nil {
			room.HistoryScores[u.ChairId] = &pk_base.HistoryScore{}
		}
		room.HistoryScores[u.ChairId].TurnScore = room.CalScoreMap[u]
		room.HistoryScores[u.ChairId].CollectScore += room.CalScoreMap[u]
	})

	room.GameStatus = GAME_NULL

	/*
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
			GameConclude.HandCardData[i] = make([]int, MAX_COUNT)
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

		nCount := 0
		if nCount > 1 {
			nCount++
		}

		for i := 0; i < nCount; i++ {
			GameConclude.MaData[i] = room.RepertoryCard[room.MinusLastCount+i]
		}

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
				room.HistoryScores[u.ChairId] = &HistoryScore{}
			}
			room.HistoryScores[u.ChairId].TurnScore = GameConclude.GameScore[u.ChairId]
			room.HistoryScores[u.ChairId].CollectScore += GameConclude.GameScore[u.ChairId]

		})

		//发送数据
		room.MjBase.UserMgr.SendMsgAll(GameConclude)

		//写入积分 todo
		room.MjBase.UserMgr.WriteTableScore(ScoreInfoArray, room.MjBase.UserMgr.GetMaxPlayerCnt(), HZMJ_CHANGE_SOURCE)
	*/
}

//解散接触
func (room *nntb_data_mgr) DismissEnd() {
	/*
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
			GameConclude.HandCardData[i] = make([]int, MAX_COUNT)
		}

		room.BankerUser = INVALID_CHAIR

		GameConclude.SendCardData = room.SendCardData

		//用户扑克
		for i := 0; i < UserCnt; i++ {
			if len(room.CardIndex[i]) > 0 {
				GameConclude.HandCardData[i] = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[i])
				GameConclude.CardCount[i] = len(GameConclude.HandCardData[i])
			}
		}

		//发送信息
		room.MjBase.UserMgr.SendMsgAll(GameConclude)
	*/
}

// 用户叫分(抢庄)
func (r *nntb_data_mgr) CallScore(u *user.User, scoreTimes int) {
	log.Debug("call score times userChairId:%d, scoretimes:%d", u.ChairId, scoreTimes)

	r.CallScoreTimesMap[scoreTimes] = u
	maxScoreTimes := 0
	for s, _ := range r.CallScoreTimesMap {
		if s > maxScoreTimes {
			maxScoreTimes = s
		}
	}
	r.BankerUser = r.CallScoreTimesMap[maxScoreTimes]
	r.ScoreTimes = maxScoreTimes

	// 广播叫分
	callScore := &nn_tb_msg.G2C_TBNN_CallScore{}
	callScore.ChairID = u.ChairId
	callScore.CallScore = scoreTimes
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u1 *user.User) {
		if u != u1 {
			u1.WriteMsg(callScore)
		}
	})

	if len(r.CallScoreTimesMap) == r.PlayerCount {
		//叫分结束
		r.CallScoreEnd()
	}
}

// 叫分结束
func (r *nntb_data_mgr) CallScoreEnd() {
	// 发回叫分结果
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(&nn_tb_msg.G2C_TBNN_CallScoreEnd{
			Banker:     r.BankerUser.ChairId,
			ScoreTimes: r.ScoreTimes,
		})
	})

	// 进入加注
	log.Debug("enter add score")
	r.GameStatus = ADD_SCORE

	r.AddScoreTimer = r.PkBase.AfterFunc(ADD_SCORE_TIME*time.Second, func() { // 超时加注结束
		if r.GameStatus != SEND_LAST_CARD {
			r.AddScoreEnd()
		}
	})
	r.AddScoreTimer.Stop()
}

// 用户加注
func (r *nntb_data_mgr) AddScore(u *user.User, score int) {
	log.Debug("add score userChairId:%d, score:%d", u.ChairId, score)
	r.AddScoreMap[u] = score

	// 广播加注
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		addScore := &nn_tb_msg.G2C_TBNN_AddScore{}
		addScore.ChairID = u.ChairId
		addScore.AddScoreCount = score
		u.WriteMsg(addScore)
	})

	if len(r.AddScoreMap) == r.PlayerCount { //全加过加注结束
		r.AddScoreEnd()
	}
}

// 加注结束
func (r *nntb_data_mgr) AddScoreEnd() {

	// 进入最后一张牌
	log.Debug("enter last card")
	r.GameStatus = SEND_LAST_CARD

	// 发最后一张牌
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		lastCard := r.GetOneCard()
		r.CardData[u.ChairId][r.GetCfg().MaxCount-1] = lastCard
		u.WriteMsg(&nn_tb_msg.G2C_TBNN_LastCard{
			LastCard: lastCard,
		})
	})

	// 进入亮牌
	r.EnterOpenCard()

}

// 进入亮牌
func (r *nntb_data_mgr) EnterOpenCard() {
	log.Debug("enter open card")
	r.GameStatus = OPEN_CARD
	// 亮牌超时
	r.OpenCardTimer = r.PkBase.AfterFunc(OPEN_CARD_TIME, func() { // 超时亮牌结束
		if r.GameStatus != CAL_SCORE {
			// 没有亮牌的用户自动亮牌
			userMgr := r.PkBase.UserMgr
			userMgr.ForEachUser(func(u *user.User) {
				if r.OpenCardMap[u] == nil {
					// 需要改进
					r.OpenCard(u, 0, r.CardData[u.ChairId])
				}
			})
		}
	})
	r.OpenCardTimer.Stop()
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
	for i := 0; i < pk_base.GetCfg(pk_base.IDX_TBNN).MaxCount; i++ {
		if !r.IsValidCard(chairID, cardData[i]) {
			return false
		}
	}
	return true
}

// 亮牌
func (r *nntb_data_mgr) OpenCard(u *user.User, cardType int, cardData []int) {
	// 验证牌数据
	if !r.IsValidCardData(u.ChairId, cardData) {
		log.Debug("user open card failed at %d", u.ChairId)
		return
	}

	// 验证牌型
	if r.PkBase.LogicMgr.GetCardType(cardData) != cardType {
		return
	}

	r.OpenCardMap[u] = cardData
	// 广播亮牌
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(&nn_tb_msg.G2C_TBNN_Open_Card{
			ChairID:  u.ChairId,
			CardType: cardType,
			CardData: cardData,
		})
	})

	if len(r.OpenCardMap) == r.PlayerCount { // 全亮过
		r.OpenCardEnd() // 亮牌结束
	}
}

// 亮牌结束
func (r *nntb_data_mgr) OpenCardEnd() {
	r.GameStatus = CAL_SCORE
	// 比牌
	logicMgr := r.PkBase.LogicMgr
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		if u != r.BankerUser { // 闲家与庄家比
			if logicMgr.CompareCard(r.OpenCardMap[r.BankerUser], r.OpenCardMap[u]) { // 庄家比闲家大
				r.CalScoreMap[r.BankerUser] += r.CellScore * r.ScoreTimes * r.AddScoreMap[u]
				r.CalScoreMap[u] -= r.CellScore * r.ScoreTimes * r.AddScoreMap[u]
			} else {
				r.CalScoreMap[r.BankerUser] -= r.CellScore * r.ScoreTimes * r.AddScoreMap[u]
				r.CalScoreMap[u] += r.CellScore * r.ScoreTimes * r.AddScoreMap[u]
			}
		}
	})

	// 游戏结束
	userMgr.ForEachUser(func(u *user.User) {
		r.PkBase.OnEventGameConclude(u.ChairId, u, cost.GER_NORMAL)
	})
}
