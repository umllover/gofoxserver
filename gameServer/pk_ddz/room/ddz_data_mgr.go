package room

import (
	"mj/common/cost"
	"mj/common/msg/pk_ddz_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"encoding/json"

	"math/rand"
	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

func NewDDZDataMgr(info *model.CreateRoomInfo, uid, ConfigIdx int, name string, temp *base.GameServiceOption, base *DDZ_Entry) *ddz_data_mgr {
	d := new(ddz_data_mgr)
	d.RoomData = pk_base.NewDataMgr(info.RoomId, uid, ConfigIdx, name, temp, base.Entry_base)

	var setInfo pk_ddz_msg.C2G_DDZ_CreateRoomInfo
	if err := json.Unmarshal([]byte(info.OtherInfo), &setInfo); err == nil {
		d.EightKing = setInfo.King
		d.GameType = setInfo.GameType
	}
	return d
}

const (
	// 游戏状态
	GAME_STATUS_FREE = 0
	GAME_STATUS_CALL = 1
	GAME_STATUS_PLAY = 2

	// 用户叫分信息
	CALLSCORE_CALLING = 0XFFFF // 正在叫分状态
	CALLSCORE_NOCALL  = 0xFFFE // 未叫状态

	// 用户出牌状态
	OUTCARD_OUTING = 0XFFFF // 出牌中
	OUTCARD_PASS   = 0XFFFE // 不出
)

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
	GameType  int  // 游戏类型
	LiziCard  int  // 癞子牌

	// 炸弹信息
	EachBombCount []int // 炸弹个数
	KingCount     []int // 八王个数

	// 叫分信息
	BankerScore int   // 庄家叫分
	ScoreInfo   []int // 叫分信息

	// 出牌信息
	TurnCardStatus []int     // 用户出牌状态
	TurnCardData   [][][]int // 出牌数据
	RepertoryCard  []int     // 库存扑克

	// 扑克信息
	BankerCard   [3]int       // 游戏底牌
	HandCardData [][]int      // 手上扑克
	ShowCardSign map[int]bool // 用户明牌标识
}

func (room *ddz_data_mgr) InitRoom(UserCnt int) {
	log.Debug("初始化房间参数")
	room.RoomData.InitRoom(UserCnt)

	room.GameStatus = GAME_STATUS_FREE
	room.CurrentUser = cost.INVALID_CHAIR
	room.CallScoreUser = cost.INVALID_CHAIR
	room.BankerUser = cost.INVALID_CHAIR
	room.TurnWiner = cost.INVALID_CHAIR

	room.TimeHeadOutCard = 0
	room.OutCardCount = make([]int, room.PlayerCount)
	room.EachBombCount = make([]int, room.PlayerCount)
	room.KingCount = make([]int, room.PlayerCount)

	room.CallScoreUser = 0
	room.BankerScore = 0
	room.ScoreInfo = make([]int, room.PlayerCount)

	room.TurnCardStatus = make([]int, room.PlayerCount)
	room.TurnCardData = make([][][]int, room.PlayerCount)
	room.RepertoryCard = make([]int, room.GetCfg().MaxRepertory)

	room.HandCardData = make([][]int, room.PlayerCount)
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
	log.Debug("发送空闲状态场景消息")
	room.GameStatus = GAME_STATUS_FREE
	StatusFree := &pk_ddz_msg.G2C_DDZ_StatusFree{}

	StatusFree.CellScore = room.PkBase.Temp.CellScore // 基础积分

	StatusFree.GameType = room.GameType
	StatusFree.EightKing = room.EightKing

	StatusFree.PlayCount = room.PkBase.TimerMgr.GetMaxPayCnt()

	StatusFree.TimeOutCard = room.PkBase.TimerMgr.GetTimeOutCard()       // 出牌时间
	StatusFree.TimeCallScore = room.GetCfg().CallScoreTime               // 叫分时间
	StatusFree.TimeStartGame = room.PkBase.TimerMgr.GetTimeOperateCard() // 开始时间
	StatusFree.TimeHeadOutCard = room.TimeHeadOutCard                    // 首出时间
	StatusFree.TurnScore = append(StatusFree.TurnScore, 12)
	StatusFree.CollectScore = append(StatusFree.CollectScore, 11)
	for _, v := range room.HistoryScores {
		StatusFree.TurnScore = append(StatusFree.TurnScore, v.TurnScore)
		StatusFree.CollectScore = append(StatusFree.TurnScore, v.CollectScore)
	}

	// 发送明牌标识
	for k, v := range room.ShowCardSign {
		StatusFree.ShowCardSign[k] = v
	}

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
	//util.DeepCopy(room.BankerCard[:], &room.RepertoryCard[len(room.RepertoryCard)-3:])
	copy(room.BankerCard[:], room.RepertoryCard[len(room.RepertoryCard)-3:])
	room.RepertoryCard = room.RepertoryCard[:len(room.RepertoryCard)-3]

	log.Debug("底牌%v", room.BankerCard)
	log.Debug("剩余牌%v", room.RepertoryCard)

	cardCount := len(room.RepertoryCard) / room.PkBase.Temp.MaxPlayer

	//构造变量
	GameStart := &pk_ddz_msg.G2C_DDZ_GameStart{}

	// 初始化叫分信息
	for i := 0; i < room.PlayerCount; i++ {
		room.ScoreInfo[i] = CALLSCORE_NOCALL
	}

	// 初始化牌
	var cardData [][]int
	for i := 0; i < room.PlayerCount; i++ {
		tempCardData := room.RepertoryCard[len(room.RepertoryCard)-cardCount:]
		room.RepertoryCard = room.RepertoryCard[:len(room.RepertoryCard)-cardCount]
		cardData = append(cardData, tempCardData)
		if room.CallScoreUser == cost.INVALID_CHAIR {
			for _, v := range tempCardData {
				if v == 0x33 {
					room.CallScoreUser = i
					room.ScoreInfo[i] = CALLSCORE_CALLING
					break
				}
			}
		}
	}

	GameStart.CallScoreUser = room.CallScoreUser

	if room.ShowCardSign == nil {
		room.ShowCardSign = make(map[int]bool)
	}

	util.DeepCopy(&GameStart.ShowCard, &room.ShowCardSign)

	//发送数据
	room.PkBase.UserMgr.ForEachUser(func(u *user.User) {

		GameStart.CardData = append([][]int{})
		for i := 0; i < room.PkBase.Temp.MaxPlayer; i++ {
			if room.ShowCardSign[i] || u.ChairId == i {
				GameStart.CardData = append(GameStart.CardData, cardData[i])
				//util.DeepCopy(GameStart.CardData[i], cardData[i])
			} else {
				GameStart.CardData = append(GameStart.CardData, nil)
			}
		}

		log.Debug("需要发送的扑克牌%v", GameStart)
		u.WriteMsg(GameStart)
	})

}

// 用户叫分(抢庄)
func (r *ddz_data_mgr) CallScore(u *user.User, scoreTimes int) {

	// 判断当前叫分玩家是否正确
	if r.ScoreInfo[u.ChairId] != CALLSCORE_CALLING {
		cost.RenderErrorMessage(cost.ErrDDZCSUser)
		return
	}

	if scoreTimes <= r.BankerScore && r.BankerScore != 0 {
		cost.RenderErrorMessage(cost.ErrDDZCSValid)
		log.Debug("用户叫分%d必须大于当前分数%d", scoreTimes, r.BankerScore)
		return
	}
	r.BankerScore = scoreTimes
	r.ScoreInfo[u.ChairId] = scoreTimes

	nextCallUser := (r.CallScoreUser + 1) % r.PlayCount // 下一个叫分玩家

	if r.ScoreInfo[nextCallUser] != CALLSCORE_NOCALL {
		// 叫分结束，发庄家信息
		r.BankerInfo()
	} else {
		r.CallScoreUser = nextCallUser
		GameCallSore := &pk_ddz_msg.G2C_DDZ_CallScore{}
		util.DeepCopy(GameCallSore.ScoreInfo, r.ScoreInfo)
		r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
			u.WriteMsg(GameCallSore)
		})
	}
}

// 检查是否能叫分
func (r *ddz_data_mgr) checkCallScore(u *user.User, scoreTimes int) int {

	return 0
}

// 庄家信息
func (r *ddz_data_mgr) BankerInfo() {
	GameBankerInfo := &pk_ddz_msg.G2C_DDZ_BankerInfo{}
	GameBankerInfo.BankerUser = r.BankerUser
	GameBankerInfo.CurrentUser = r.CurrentUser
	GameBankerInfo.BankerScore = r.BankerScore

	// 随机选一张牌为癞子牌
	ran := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.LiziCard = ran.Intn(13)

	GameBankerInfo.LiziCard = r.LiziCard

	r.PkBase.LogicMgr.SetParamToLogic(r.LiziCard)

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
func (r *ddz_data_mgr) OpenCard(u *user.User, cardType int, cardData []int) {
	// 检查当前是否该用户出牌
	if r.CurrentUser != u.ChairId {

		return
	}

	// 检查所出牌是否完整在手上
	if len(cardData) > len(r.HandCardData[u.ChairId]) {

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
			return
		}
	}

	// 对比牌型
	r.PkBase.LogicMgr.SortCardList(cardData, len(cardData)) // 排序

	r.TurnWiner = u.ChairId
	r.CurrentUser = r.nextUser(r.CurrentUser)
	r.TurnCardStatus[r.CurrentUser] = OUTCARD_OUTING
	// 把所出的牌存到出牌数据里
	r.TurnCardData[u.ChairId] = append(r.TurnCardData[u.ChairId], cardData)
	r.TurnCardStatus[u.ChairId] = len(r.TurnCardData[u.ChairId]) - 1

	// 从手牌删除数据
	r.PkBase.LogicMgr.RemoveCardList(cardData, r.HandCardData[u.ChairId])
}

// 下一个玩家
func (r *ddz_data_mgr) nextUser(u int) int {
	return (u + 1) % r.PlayCount
}

func (r *ddz_data_mgr) OutCard() {
	DataOutCard := &pk_ddz_msg.G2C_DDZ_OutCard{}

	DataOutCard.CurrentUser = r.CurrentUser
	DataOutCard.OutCardUser = r.TurnWiner
	util.DeepCopy(DataOutCard.CardData, r.TurnCardData[r.TurnWiner][len(r.TurnCardData[r.TurnWiner])])

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

	r.TurnCardStatus[r.CurrentUser] = OUTCARD_PASS
	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataPassCard)
	})
}

// 游戏正常结束
func (r *ddz_data_mgr) NormalEnd() {
	DataGameConclude := &pk_ddz_msg.G2C_DDZ_GameConclude{}
	DataGameConclude.CellScore = r.CellScore
	DataGameConclude.GameScore = make([]int, r.PlayCount)

	// 算分数
	nMultiple := r.ScoreTimes

	// 春天标识
	if r.OutCardCount[r.BankerUser] <= 1 {
		DataGameConclude.SpringSign = 2 // 地主只出了一次牌
	} else {
		DataGameConclude.SpringSign = 1
		for i := 1; i < r.PlayCount; i++ {
			if r.OutCardCount[(i+r.BankerUser)%r.PlayCount] > 0 {
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

	for i := 0; i < r.PlayCount; i++ {
		if i == r.BankerUser {
			DataGameConclude.GameScore[i] = (0 - gameScore) * (r.PlayCount - 1)
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
	DataGameConclude.GameScore = make([]int, r.PlayCount)

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
	room.PkBase.UserMgr.SetUsetTrustee(u.ChairId, t)
	DataTrustee := &pk_ddz_msg.G2C_DDZ_TRUSTEE{}
	DataTrustee.TrusteeUser = u.ChairId
	DataTrustee.Trustee = t

	room.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(DataTrustee)
	})
}
