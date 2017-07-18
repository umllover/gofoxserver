package room

import (
	"mj/common/cost"
	"mj/common/msg/pk_ddz_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"encoding/json"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

func NewDDZDataMgr(info *model.CreateRoomInfo, uid int64, ConfigIdx int, name string, temp *base.GameServiceOption, base *DDZ_Entry) *ddz_data_mgr {
	d := new(ddz_data_mgr)
	d.RoomData = pk_base.NewDataMgr(info.RoomId, uid, ConfigIdx, name, temp, base.Entry_base)

	var setInfo pk_ddz_msg.C2G_DDZ_CreateRoomInfo
	if err := json.Unmarshal([]byte(info.OtherInfo), &setInfo); err == nil {
		d.EightKing = setInfo.King
	}
	return d
}

type ddz_data_mgr struct {
	*pk_base.RoomData
	GameStatus int // 当前游戏状态

	CurrentUser   int // 当前玩家
	CallScoreUser int // 叫分玩家
	BankerUser    int // 地主
	TurnWiner     int // 出牌玩家

	TimeHeadOutCard int // 首出时间

	OutCardCount []int // 出牌次数

	EightKing bool // 是否八王模式

	// 炸弹信息
	EachBombCount []int // 炸弹个数
	KingCount     []int // 八王个数

	// 叫分信息
	CallScoreCount int   // 叫分次数
	BankerScore    int   // 庄家叫分
	ScoreInfo      []int // 叫分信息

	// 出牌信息
	TurnCardData  []int // 出牌数据
	RepertoryCard []int // 库存扑克

	// 扑克信息
	BankerCard   [3]int  // 游戏底牌
	HandCardData [][]int // 手上扑克
	ShowCardSign []bool  // 用户明牌标识
	TrusteeSign  []bool  // 托管标识
}

func (room *ddz_data_mgr) InitRoom(UserCnt int) {
	room.RoomData.InitRoom(UserCnt)

	room.GameStatus = GAME_STATUS_FREE
	room.CurrentUser = cost.INVALID_CHAIR
	room.CallScoreUser = cost.INVALID_CHAIR
	room.BankerUser = cost.INVALID_CHAIR
	room.TurnWiner = cost.INVALID_CHAIR

	room.TimeHeadOutCard = 0
	room.OutCardCount = make([]int, room.CurrentPlayCount)
	room.EachBombCount = make([]int, room.CurrentPlayCount)
	room.KingCount = make([]int, room.CurrentPlayCount)

	room.CallScoreUser = 0
	room.BankerScore = 0
	room.ScoreInfo = make([]int, room.CurrentPlayCount)

	room.TurnCardData = make([]int, room.CurrentPlayCount)
	room.RepertoryCard = make([]int, room.CurrentPlayCount)

	room.HandCardData = make([][]int, room.CurrentPlayCount)
	room.ShowCardSign = make([]bool, room.CurrentPlayCount)
	room.TrusteeSign = make([]bool, room.CurrentPlayCount)
}

// 游戏开始
func (room *ddz_data_mgr) StartGameing() {
	room.GameStatus = GAME_STATUS_PLAY
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
	/*for _, v := range room.HistoryScores {
		StatusFree.TurnScore = append(StatusFree.TurnScore, v.TurnScore)
		StatusFree.CollectScore = append(StatusFree.TurnScore, v.CollectScore)
	}*/

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

	/*//历史积分
	for j := 0; j < UserCnt; j++ {
		//设置变量
		if room.HistoryScores[j] != nil {
			StatusCall.TurnScore[j] = room.HistoryScores[j].TurnScore
			StatusCall.CollectScore[j] = room.HistoryScores[j].CollectScore
		}
	}*/

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

	/*//历史积分
	for j := 0; j < UserCnt; j++ {
		//设置变量
		if room.HistoryScores[j] != nil {
			StatusPlay.TurnScore[j] = room.HistoryScores[j].TurnScore
			StatusPlay.CollectScore[j] = room.HistoryScores[j].CollectScore
		}
	}*/

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

	cardCount := int(len(room.RepertoryCard) / room.CurrentPlayCount)

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
	r.ScoreInfo[u.ChairId] = scoreTimes
	r.CallScoreCount++

	if r.CallScoreCount >= r.PlayerCount {
		//叫分结束，发庄家信息
		r.BankerInfo()
	} else {
		r.CurrentUser = (r.CallScoreUser + 1) % r.CurrentPlayCount
		GameCallSore := &pk_ddz_msg.G2C_DDZ_CallScore{}
		GameCallSore.CurrentUser = r.CurrentUser
		GameCallSore.CallScoreUser = r.CallScoreUser
		GameCallSore.CurrentScore = r.BankerScore
		GameCallSore.UserCallScore = r.ScoreInfo[u.ChairId]
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
func (r *ddz_data_mgr) ShowCard(u *user.User) {
	r.ShowCardSign[u.ChairId] = true

	DataShowCard := &pk_ddz_msg.G2C_DDZ_ShowCard{}
	DataShowCard.ShowCardUser = u.ChairId
	util.DeepCopy(DataShowCard.CardData, r.HandCardData[u.ChairId])

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataShowCard)
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
	DataGameConclude := &pk_ddz_msg.G2C_DDZ_GameConclude{}
	DataGameConclude.CellScore = r.CellScore
	DataGameConclude.GameScore = make([]int, r.CurrentPlayCount)

	// 算分数
	nMultiple := r.ScoreTimes

	// 春天标识
	if r.OutCardCount[r.BankerUser] <= 1 {
		DataGameConclude.SpringSign = 2 // 地主只出了一次牌
	} else {
		DataGameConclude.SpringSign = 1
		for i := 1; i < r.CurrentPlayCount; i++ {
			if r.OutCardCount[(i+r.BankerUser)%r.CurrentPlayCount] > 0 {
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

	// 明牌翻倍
	for _, v := range r.ShowCardSign {
		if v {
			nMultiple <<= 1
		}
	}

	// 八王
	util.DeepCopy(DataGameConclude.KingCount, r.KingCount)
	for _, v := range r.KingCount {
		if v == 8 {
			nMultiple *= 8 * 2
		} else if v >= 2 {
			nMultiple *= v
		}
	}

	// 炸弹
	util.DeepCopy(DataGameConclude.EachBombCount, r.EachBombCount)

	DataGameConclude.BankerScore = r.BankerScore
	util.DeepCopy(DataGameConclude.HandCardData, r.HandCardData)

	// 计算积分
	gameScore := r.BankerScore * nMultiple

	if len(r.HandCardData[r.BankerUser]) <= 0 {
		gameScore = 0 - gameScore
	}

	for i := 0; i < r.CurrentPlayCount; i++ {
		if i == r.BankerUser {
			DataGameConclude.GameScore[i] = (0 - gameScore) * (r.CurrentPlayCount - 1)
		} else {
			DataGameConclude.GameScore[i] = gameScore
		}
	}

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataGameConclude)
	})
}

//解散房间结束
func (r *ddz_data_mgr) DismissEnd() {
	DataGameConclude := &pk_ddz_msg.G2C_DDZ_GameConclude{}
	DataGameConclude.CellScore = r.CellScore
	DataGameConclude.GameScore = make([]int, r.CurrentPlayCount)

	// 炸弹
	util.DeepCopy(DataGameConclude.EachBombCount, r.EachBombCount)

	DataGameConclude.BankerScore = r.BankerScore
	util.DeepCopy(DataGameConclude.HandCardData, r.HandCardData)

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataGameConclude)
	})
}

// 托管
func (room *ddz_data_mgr) Trustee(u *user.User, t bool) {
	room.TrusteeSign[u.ChairId] = t
	DataTrustee := &pk_ddz_msg.G2C_DDZ_TRUSTEE{}
	DataTrustee.TrusteeUser = u.ChairId
	DataTrustee.Trustee = t

	room.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataTrustee)
	})
}
