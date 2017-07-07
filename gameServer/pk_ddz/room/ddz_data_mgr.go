package room

import (
	//"mj/common/cost"
	//"mj/common/msg/pk_ddz_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"
	//"mj/gameServer/user"

	//"github.com/lovelly/leaf/log"
	//"github.com/lovelly/leaf/util"
)

func NewDataMgr(id, uid, ConfigIdx int, name string, temp *base.GameServiceOption, base *DDZ_Entry) *ddz_data_mgr {
	d := new(ddz_data_mgr)
	d.RoomData = pk_base.NewDataMgr(id, uid, ConfigIdx, name, temp, base.Entry_base)

	return d
}

type ddz_data_mgr struct {
	*pk_base.RoomData
	CurrentUser   int   // 当前玩家
	OutCardCount  []int // 出牌次数
	CallScoreUser int   // 叫分玩家

	TimeHeadOutCard int // 首出时间

	// 托管信息
	OffLineTrustee bool // 离线托管

	// 炸弹信息
	BombCount     int   // 炸弹个数
	EachBombCount []int // 炸弹个数

	// 叫分信息
	CallScoreCount int   // 叫分次数
	BankerScore    int   // 庄家叫分
	ScoreInfo      []int // 叫分信息

	// 出牌信息
	TurnWiner     int   // 出牌玩家
	TurnCardCount int   // 出牌数目
	TurnCardData  []int // 出牌数据

	// 扑克信息
	BankerCard    [3]int       // 游戏底牌
	HandCardCount []int        // 扑克数目
	HandCardData  [][]int      // 手上扑克
	ShowCardSign  map[int]bool // 用户明牌标识
}

/*
func (room *ddz_data_mgr) InitRoom(UserCnt int) {
	//初始化
	room.CardData = make([][]int, UserCnt)
	room.PublicCardData = make([]int, room.GetCfg().PublicCardCount)

	room.CallScoreTimesMap = make(map[int]*user.User)
	room.ScoreMap = make(map[*user.User]int)
	room.OpenCardMap = make(map[*user.User][]int)
	room.RepertoryCard = make([]int, room.GetCfg().MaxRepertory)

	room.ExitScore = 0
	room.DynamicScore = 0
	room.BankerUser = cost.INVALID_CHAIR
	room.FisrtCallUser = cost.INVALID_CHAIR
	room.CurrentUser = cost.INVALID_CHAIR

	//room.MaxScoreTimes = 0
	room.Count_All = 0
	room.Qiang = make([]bool, UserCnt)
	room.IsOpenCard = make([]bool, UserCnt)
	room.DynamicJoin = make([]int, UserCnt)
	room.TableScore = make([]int64, UserCnt)
	room.PlayStatus = make([]int, UserCnt)
	room.CallStatus = make([]int, UserCnt)
	room.EscapeUserScore = make([]int64, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.OxCard[i] = 0xFF
	}

	room.CallScoreUser = cost.INVALID_CHAIR
}

// 游戏开始
func (room *ddz_data_mgr) StartGameing() {
	room.GameStatus = pk_base.GAME_START
	room.SendGameStart()
}

//// 叫分
//func (room *ddz_data_mgr) AfterStartGame() {
//	room.GameStatus = pk_base.CALL_SCORE_TIMES
//	if room.CallScoreUser == cost.INVALID_CHAIR {
//		room.CallScore(room.PkBase.UserMgr.GetUserByChairId(0), 0)
//	} else {
//		room.CallScore(room.PkBase.UserMgr.GetUserByChairId(room.CallScoreUser), 0)
//	}
//}

// 空闲状态场景
func (room *ddz_data_mgr) SendStatusReady(u *user.User) {
	StatusFree := &pk_ddz_msg.G2C_DDZ_StatusFree{}

	StatusFree.CellScore = room.CellScore                                // 基础积分
	StatusFree.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()       // 出牌时间
	StatusFree.TimeCallScore = room.GetCfg().CallScoreTime               // 叫分时间
	StatusFree.TimeStartGame = room.PkBase.TimerMgr.GetTimeOperateCard() // 开始时间
	StatusFree.TimeHeadOutCard = room.TimeHeadOutCard                    // 首出时间
	for _, v := range room.HistoryScores {
		StatusFree.TurnScore = append(StatusFree.TurnScore, v.TurnScore)
		StatusFree.CollectScore = append(StatusFree.TurnScore, v.CollectScore)
	}

	u.WriteMsg(StatusFree)
}

// 叫分状态
func (room *ddz_data_mgr) SendStatusCall(u *user.User) {
	StatusCall := &pk_ddz_msg.G2C_DDZ_StatusCall{}

	StatusCall.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()
	StatusCall.TimeCallScore = 0
	StatusCall.TimeStartGame = int(room.PkBase.TimerMgr.GetCreatrTime())

	StatusCall.CellScore = room.CellScore
	StatusCall.CurrentUser = room.CurrentUser
	StatusCall.BankerScore = room.BankerScore

	UserCnt := room.PkBase.UserMgr.GetMaxPlayerCnt()
	StatusCall.ScoreInfo = make([]int, UserCnt)

	StatusCall.HandCardData = nil
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

	// 明牌数据
	for k, v := range room.ShowCardSign {
		if v == true {
			util.DeepCopy(StatusCall.ShowCardData, room.HandCardData[k])
		}

	}

	u.WriteMsg(StatusCall)
}

// 游戏状态
func (room *ddz_data_mgr) SendStatusPlay(u *user.User) {
	StatusPlay := &pk_ddz_msg.G2C_DDZ_StatusPlay{}
	//自定规则
	StatusPlay.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()
	StatusPlay.TimeCallScore = room.PkBase.TimerMgr.GetTimeOperateCard()
	StatusPlay.TimeStartGame = int(room.PkBase.TimerMgr.GetCreatrTime())

	//游戏变量
	StatusPlay.CellScore = room.CellScore
	StatusPlay.BombCount = 0
	StatusPlay.BankerUser = room.BankerUser
	StatusPlay.CurrentUser = room.CurrentUser
	StatusPlay.BankerScore = room.BankerScore

	StatusPlay.TurnWiner = 0
	StatusPlay.TurnCardCount = 0
	StatusPlay.TurnCardData = nil
	util.DeepCopy(StatusPlay.BankerCard, room.BankerCard)
	StatusPlay.HandCardData = nil
	StatusPlay.HandCardCount = nil

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
	util.DeepCopy(room.BankerCard[:], room.RepertoryCard[len(room.RepertoryCard)-3:])
	room.RepertoryCard = room.RepertoryCard[:len(room.RepertoryCard)-3]

	cardCount := int(len(room.RepertoryCard) / room.PlayCount)

	//构造变量
	GameStart := &pk_ddz_msg.G2C_DDZ_GameStart{}
	GameStart.StartUser = 0
	GameStart.CurrentUser = room.CurrentUser
	GameStart.ValidCardData = 0
	GameStart.ValidCardIndex = 0

	//发送数据
	room.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		GameStart.CardData = room.RepertoryCard[len(room.RepertoryCard)-cardCount:]
		room.RepertoryCard = room.RepertoryCard[:len(room.RepertoryCard)-cardCount]
		if room.CallScoreUser == cost.INVALID_CHAIR {
			for _, v := range GameStart.CardData {
				if v == 0x33 {
					room.CallScoreUser = u.ChairId
					break
				}
			}
		}
		u.WriteMsg(GameStart)
	})

}

// 用户叫分(抢庄)
func (r *ddz_data_mgr) CallScore(u *user.User, scoreTimes int) {

	if scoreTimes <= r.BankerScore && r.BankerScore != 0 {
		log.Debug("用户叫分%d必须大于当前分数%d", scoreTimes, r.BankerScore)
		return
	}
	r.BankerScore = scoreTimes
	r.CallScoreTimesMap[scoreTimes] = u

	if len(r.CallScoreTimesMap) == r.PlayerCount {
		//叫分结束，发庄家信息
		r.BankerInfo()
	} else {
		r.CurrentUser = (r.CallScoreUser + 1) % r.PlayCount
		GameCallSore := &pk_ddz_msg.G2C_DDZ_CallScore{}
		GameCallSore.CurrentUser = r.CurrentUser
		GameCallSore.CallScoreUser = r.CallScoreUser
		GameCallSore.CurrentScore = r.BankerScore
		GameCallSore.UserCallScore = r.CallScoreTimesMap[scoreTimes].ChairId
		r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
			u.WriteMsg(GameCallSore)
		})
	}
}

// 庄家信息
func (r *ddz_data_mgr) BankerInfo() {
	GameBankerInfo := &pk_ddz_msg.G2C_DDZ_BankerInfo{}
	GameBankerInfo.BankerUser = r.BankerUser
	GameBankerInfo.CurrentUser = r.CurrentUser
	GameBankerInfo.BankerScore = r.BankerScore

	util.DeepCopy(GameBankerInfo.BankerCard, r.BankerCard)

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(GameBankerInfo)
	})
}

// 明牌
func (r *ddz_data_mgr) ShowCard(u *user.User, cardData []int) {
	r.ShowCardSign[u.ChairId] = true

	DataShowCard := &pk_ddz_msg.G2C_DDZ_ShowCard{}
	DataShowCard.ShowCardUser = u.ChairId
	util.DeepCopy(DataShowCard.CardData, cardData)

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataShowCard)
	})
}

// 亮牌
func (r *ddz_data_mgr) OpenCard(u *user.User, cardData []int) {
	DataOutCard := &pk_ddz_msg.G2C_DDZ_OutCard{}

	DataOutCard.CardCount = len(cardData)
	DataOutCard.CurrentUser = r.CurrentUser
	DataOutCard.OutCardUser = u.ChairId

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataOutCard)
	})
}

// 用户出牌
func (r *ddz_data_mgr) OutCard() {
	DataOutCard := &pk_ddz_msg.G2C_DDZ_OutCard{}

	DataOutCard.CardCount = 0
	DataOutCard.CurrentUser = r.CurrentUser
	DataOutCard.OutCardUser = 0

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataOutCard)
	})
}

// 放弃出牌
func (r *ddz_data_mgr) PassCard() {
	DataPassCard := &pk_ddz_msg.G2C_DDZ_PassCard{}

	DataPassCard.TurnOver = 0
	DataPassCard.CurrentUser = r.CurrentUser
	DataPassCard.PassCardUser = 0

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataPassCard)
	})
}

// 游戏正常结束
func (r *ddz_data_mgr) NormalEnd() {

}

//解散房间结束
func (r *ddz_data_mgr) DismissEnd() {

}
*/
