package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"
	"github.com/lovelly/leaf/timer"

	"mj/gameServer/user"
	"time"
	"mj/common/msg/nn_tb_msg"
	"mj/common/cost"
	"github.com/lovelly/leaf/util"
	"github.com/lovelly/leaf/log"
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
	CAL_SCORE		= 1007// 结算
)

// 定时器 -- for test
const (
	CALL_SCORE_TIME = 10
	ADD_SCORE_TIME  = 10
	OPEN_CARD_TIME	= 30
)

func NewDataMgr(id, uid, ConfigIdx int, name string, temp *base.GameServiceOption, base *NNTB_Entry) *nntb_data_mgr {
	d := new(nntb_data_mgr)
	d.RoomData = pk_base.NewDataMgr(id, uid, ConfigIdx, name, temp, base.Entry_base)
	return d
}

type nntb_data_mgr struct {
	*pk_base.RoomData

	//游戏变量
	CardData       [][]int //用户扑克
	PublicCardData []int   //公共牌 两张
	RepertoryCard  []int   //库存扑克
	LeftCardCount  int     //库存剩余扑克数量
	OpenCardMap			map[*user.User][]int 	//亮牌数据
	CallScoreTimesMap 	map[int]*user.User 		//记录叫分信息
	CalScoreMap			map[*user.User]int		//计分信息
	AddScoreMap 		map[*user.User]int //记录用户加注信息



	BankerUser      *user.User     //庄家用户



	// 游戏状态
	GameStatus     	int
	CallScoreTimer 	*timer.Timer
	AddScoreTimer  	*timer.Timer
	OpenCardTimer	*timer.Timer

}


func (room *nntb_data_mgr) SendStatusReady(u *user.User) {
	StatusFree := &nn_tb_msg.G2C_TBNN_StatusFree{}

	StatusFree.CellScore = room.PkBase.Temp.CellScore                                    //基础积分
	StatusFree.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()         //出牌时间
	StatusFree.TimeOperateCard = room.PkBase.TimerMgr.GetTimeOperateCard() //操作时间
	StatusFree.TimeStartGame = room.PkBase.TimerMgr.GetCreatrTime()        //开始时间
	StatusFree.TurnScore = make([]int, room.PkBase.TimerMgr.GetMaxPayCnt())
	StatusFree.CollectScore = make([]int, room.PlayerCount)
	StatusFree.EachRoundScore = make([][]int, room.PlayerCount, room.PkBase.TimerMgr.GetMaxPayCnt())

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
	log.Debug("init room")
	room.InitRoom(UserCnt)
}

func (room *nntb_data_mgr) StartGameing() {
	// 发牌
	log.Debug("dispatch card")
	room.StartDispatchCard()
}


func (room *nntb_data_mgr) AfterStartGame() {
	// 叫分
	log.Debug("call score")
	room.GameStatus = CALL_SCORE_TIMES
	log.Debug("begin call score timer")
	room.CallScoreTimer = room.PkBase.AfterFunc(CALL_SCORE_TIME * time.Second, func() {
		log.Debug("end call score timer")
		if room.GameStatus != ADD_SCORE { // 超时叫分结束
			room.CallScoreEnd()
		}
	})

	room.CallScoreTimer.Stop()

}

func (room *nntb_data_mgr) InitRoom(UserCnt int) {
	//初始化
	room.CardData = make([][]int, UserCnt)
	for i:=0;i<UserCnt;i++ {
		room.CardData[i] = make([]int, room.GetCfg().MaxCount)
	}
	room.PublicCardData = make([]int, room.GetCfg().PublicCardCount)
	room.LeftCardCount = room.GetCfg().MaxRepertory

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

func (r *nntb_data_mgr) GetOneCard() int  { // 从牌堆取出一张
	r.LeftCardCount --
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

}

//解散接触
func (room *nntb_data_mgr) DismissEnd() {
	/*
		//变量定义
		UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
		GameConclude := &mj_hz_msg.G2C_HZMJ_GameConclude{}
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
	if r.GameStatus != CALL_SCORE_TIMES {
		return
	}
	log.Debug("call score times userChairId:%d, scoretimes:%d", u.ChairId, scoreTimes)

	r.CallScoreTimesMap[scoreTimes] = u

	maxScoreTimes := 1
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
func (r * nntb_data_mgr) CallScoreEnd()  {
	log.Debug("call score end")
	// 发回叫分结果
	userMgr := r.PkBase.UserMgr
	//如果没有任何人叫分
	if len(r.CallScoreTimesMap) == 0 {
		r.BankerUser = userMgr.GetUserByChairId(0)
		r.ScoreTimes = 1
	}
	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(&nn_tb_msg.G2C_TBNN_CallScoreEnd{
			Banker:     r.BankerUser.ChairId,
			ScoreTimes: r.ScoreTimes,
		})
	})

	// 进入加注
	log.Debug("enter add score")
	r.GameStatus = ADD_SCORE

	r.AddScoreTimer = r.PkBase.AfterFunc(ADD_SCORE_TIME * time.Second, func() { // 超时加注结束
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
func (r * nntb_data_mgr) AddScoreEnd() {

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
func (r *nntb_data_mgr) EnterOpenCard()  {
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
					r.OpenCard(u,0, r.CardData[u.ChairId] )
				}
			})
		}
	})
	r.OpenCardTimer.Stop()
}

// 验证
func (r *nntb_data_mgr) IsValidCard(chairID int, card int) bool  {
	// 先验证是不是在公共牌中
	for i:=0; i<pk_base.GetCfg(pk_base.IDX_TBNN).PublicCardCount; i++ {
		if card == r.PublicCardData[i] {
			return true
		}
	}
	// 是不是在用户手牌
	for i:=0; i<pk_base.GetCfg(pk_base.IDX_TBNN).MaxCount; i++ {
		if card == r.CardData[chairID][i] {
			return true
		}
	}
	return false
}

func (r *nntb_data_mgr)IsValidCardData(chairID int, cardData []int) bool {
	for i:=0; i<pk_base.GetCfg(pk_base.IDX_TBNN).MaxCount; i++ {
		if !r.IsValidCard(chairID,cardData[i]) {
			return false
		}
	}
	return true
}

// 亮牌
func (r *nntb_data_mgr) OpenCard(u *user.User, cardType int, cardData []int)  {
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
			ChairID:u.ChairId,
			CardType:cardType,
			CardData:cardData,
		})
	})

	if len(r.OpenCardMap) == r.PlayerCount { // 全亮过
		r.OpenCardEnd() // 亮牌结束
	}
}

// 亮牌结束
func (r *nntb_data_mgr) OpenCardEnd()  {
	r.GameStatus = CAL_SCORE
	// 比牌
	logicMgr := r.PkBase.LogicMgr
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		if u != r.BankerUser { // 闲家与庄家比
			if logicMgr.CompareCard(r.OpenCardMap[r.BankerUser], r.OpenCardMap[u]) { // 庄家比闲家大
				r.CalScoreMap[r.BankerUser] += r.CellScore * r.ScoreTimes * r.AddScoreMap[u]
				r.CalScoreMap[u] -= r.CellScore * r.ScoreTimes * r.AddScoreMap[u]
			}else {
				r.CalScoreMap[r.BankerUser] -= r.CellScore * r.ScoreTimes * r.AddScoreMap[u]
				r.CalScoreMap[u] += r.CellScore * r.ScoreTimes * r.AddScoreMap[u]
			}
		}
	})

	// 游戏结束
	userMgr.ForEachUser(func(u *user.User) {
		r.PkBase.OnEventGameConclude(u.ChairId, u, cost.GER_NORMAL )
	})
}

