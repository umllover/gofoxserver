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

const (
	GAME_STATUS_FREE = 0
	GAME_STATUS_CALL = 1
	GAME_STATUS_PLAY = 2
)

type ddz_data_mgr struct {
	*pk_base.RoomData
	GameStatus int // 当前游戏状态

	CurrentUser   int // 当前玩家
	CallScoreUser int // 叫分玩家
	BankerUser    int // 地主
	TurnWiner     int // 出牌玩家

	TimeHeadOutCard int // 首出时间

	EightKing bool // 是否八王模式

	// 炸弹信息
	EachBombCount []int // 炸弹个数
	KingCount     []int // 八王个数

	// 叫分信息
	ScoreInfo []int // 叫分信息

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
	room.EachBombCount = make([]int, room.PlayerCount)

	room.CallScoreUser = 0
	room.ScoreInfo = make([]int, room.PlayCount)

	room.TurnCardData = make([]int, room.PlayCount)
	room.RepertoryCard = make([]int, room.PlayCount)

	room.HandCardData = make([][]int, room.PlayCount)
	room.ShowCardSign = make([]bool, room.PlayCount)
	room.TrusteeSign = make([]bool, room.PlayCount)
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
	StatusCall.BankerScore = room.ScoreTimes

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
	StatusPlay.BankerScore = room.ScoreTimes

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
		room.PkBase.LogicMgr.SortCardList(tempCardData, len(tempCardData))
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

	if scoreTimes <= r.ScoreTimes && r.ScoreTimes != 0 {
		log.Debug("用户叫分%d必须大于当前分数%d", scoreTimes, r.ScoreTimes)
		return
	}
	r.ScoreTimes = scoreTimes
	r.ScoreInfo[u.ChairId] = scoreTimes
	r.CallScoreCount++

	nextCallUser := (r.CallScoreUser + 1) % r.PlayerCount // 下一个叫分玩家

	isEnd := (r.ScoreTimes == CALLSCORE_MAX) || (r.ScoreInfo[nextCallUser] != CALLSCORE_NOCALL)

	if !isEnd {
		r.ScoreInfo[nextCallUser] = CALLSCORE_CALLING
		r.CallScoreUser = nextCallUser
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
	}
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

	var lastTurnCard []int
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

	if !r.PkBase.LogicMgr.CompareCard(lastTurnCard, cardData) {
		log.Debug("出牌数据错误")
		return
	}

	r.TurnWiner = u.ChairId
	r.CurrentUser = r.nextUser(r.CurrentUser)
	r.TurnCardStatus[r.CurrentUser] = OUTCARD_OUTING
	// 把所出的牌存到出牌数据里
	r.TurnCardData[u.ChairId] = append(r.TurnCardData[u.ChairId], cardData)
	r.TurnCardStatus[u.ChairId] = len(r.TurnCardData[u.ChairId]) - 1

	// 从手牌删除数据
	//r.HandCardData[u.ChairId], _ = r.PkBase.LogicMgr.RemoveCardList(cardData, r.HandCardData[u.ChairId])

	// 发送给所有玩家
	DataOutCard := pk_ddz_msg.G2C_DDZ_OutCard{}
	DataOutCard.CurrentUser = r.CurrentUser
	DataOutCard.OutCardUser = u.ChairId
	DataOutCard.CardData = make([]int, len(cardData))
	util.DeepCopy(&DataOutCard.CardData, &cardData)
	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		log.Debug("出牌数据%v", DataOutCard)
		u.WriteMsg(DataOutCard)
	})

	if len(r.HandCardData[u.ChairId]) == 0 {
		log.Debug("游戏结束")
		r.PkBase.OnEventGameConclude(0, nil, cost.GER_NORMAL)
		return
	}
	r.checkNextUserTrustee()
}

// 判断下一个玩家托管状态
func (r *ddz_data_mgr) checkNextUserTrustee() {
	log.Debug("当前玩家%d,托管状态%v", r.CurrentUser, r.PkBase.UserMgr.GetTrustees())
	if r.PkBase.UserMgr.IsTrustee(r.CurrentUser) {
		// 出牌玩家为托管状态
		if r.CurrentUser == r.TurnWiner {
			// 上一个出牌玩家是自己，则选最小牌
			var cardData []int
			cardData = append(cardData, r.HandCardData[r.CurrentUser][len(r.HandCardData[r.CurrentUser])-1])
			r.OpenCard(r.PkBase.UserMgr.GetUserByChairId(r.CurrentUser), 0, cardData)
		} else {
			// 上一个出牌玩家不是自己，则不出
			r.PassCard(r.PkBase.UserMgr.GetUserByChairId(r.CurrentUser))
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
	if r.TurnWiner == u.ChairId {
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
	DataGameConclude := &pk_ddz_msg.G2C_DDZ_GameConclude{}
	DataGameConclude.CellScore = r.CellScore
	DataGameConclude.GameScore = make([]int, r.PlayCount)

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

	// 明牌翻倍
	for _, v := range r.ShowCardSign {
		if v {
			nMultiple <<= 1
		}
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
	gameScore := r.CellScore * nMultiple

	if len(r.HandCardData[r.BankerUser]) <= 0 {
		gameScore = 0 - gameScore
	}

	for i := 0; i < r.PlayerCount; i++ {
		if i == r.BankerUser {
			DataGameConclude.GameScore = append(DataGameConclude.GameScore, (0-gameScore)*(r.PlayerCount-1))
		} else {
			DataGameConclude.GameScore = append(DataGameConclude.GameScore, gameScore)
		}
	}

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		log.Debug("游戏结算信息%v", DataGameConclude)
		u.WriteMsg(DataGameConclude)
	})
}

//解散房间结束
func (r *ddz_data_mgr) DismissEnd() {
	DataGameConclude := &pk_ddz_msg.G2C_DDZ_GameConclude{}
	DataGameConclude.CellScore = r.CellScore
	DataGameConclude.GameScore = make([]int, r.PlayCount)

	// 炸弹
	util.DeepCopy(&DataGameConclude.EachBombCount, &r.EachBombCount)

	DataGameConclude.BankerScore = r.ScoreTimes
	util.DeepCopy(&DataGameConclude.HandCardData, &r.HandCardData)

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
		log.Debug("托管状态%v", DataTrustee)
		u.WriteMsg(DataTrustee)
	})
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
