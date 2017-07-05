package pk_base

import (
	"strconv"
	"time"

	"mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/nn_tb_msg"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"

	"github.com/lovelly/leaf/timer"
	"github.com/lovelly/leaf/util"
)

func NewDataMgr(id, uid, ConfigIdx int, name string, temp *base.GameServiceOption, base *Entry_base) *RoomData {
	r := new(RoomData)
	r.id = id
	if name == "" {
		r.Name = temp.RoomName
	} else {
		r.Name = name
	}
	r.CreateUser = uid
	r.PkBase = base
	r.ConfigIdx = ConfigIdx
	return r
}

//当一张桌子理解
type RoomData struct {
	id         int
	Name       string //房间名字
	CreateUser int    //创建房间的人
	PkBase   *Entry_base
	ConfigIdx  int //配置文件索引

	IsGoldOrGameScore int    //金币场还是积分场 0 标识 金币场 1 标识 积分场
	Password          string // 密码

	//游戏变量
	CardData       [][]int //用户扑克
	PublicCardData []int   //公共牌 两张
	RepertoryCard  []int   //库存扑克
	LeftCardCount  int     //库存剩余扑克数量
	OpenCardMap		map[*user.User][]int //亮牌数据
	Count_All         int                //是不是全部都点了 全点了直接动画
	Qiang             []bool             //1抢 0不抢
	CellScore         int                //底分
	ScoreTimes        int                //倍数
	CallScoreTimesMap map[int]*user.User //记录叫分信息
	PlayCount         int                //游戏局数
	PlayerCount       int                //指定游戏人数，2-4

	BankerUser      int     //庄家用户
	FisrtCallUser   int     //始叫用户
	CurrentUser     int     //当前用户
	ExitScore       int64   //强退分数
	EscapeUserScore []int64 //逃跑玩家分数
	DynamicScore    int64   //总分

	//用户数据
	IsOpenCard          []bool  //是否摊牌
	DynamicJoin         []int   //动态加入
	PlayStatus          []int   //游戏状态
	CallStatus          []int   //叫庄状态
	OxCard              []int   //牛牛数据
	TableScore          []int64 //下注数目
	BuckleServiceCharge []bool  //收服务费

	//下注信息
	//TurnMaxScore		[]int64			//最大下注
	//MaxScoreTimes		int				//最大倍数
	ScoreMap map[*user.User]int //记录用户加注信息

	//历史积分
	HistoryScores []*HistoryScore //历史积分

	// 游戏状态
	GameStatus     	int
	CallScoreTimer 	*timer.Timer
	AddScoreTimer  	*timer.Timer
	OpenCardTimer	*timer.Timer
}

func (room *RoomData) GetCfg() *PK_CFG {
	return GetCfg(room.ConfigIdx)
}

func (room *RoomData) CanOperatorRoom(uid int) bool {
	if uid == room.CreateUser {
		return true
	}
	return false
}

func (room *RoomData) GetCurrentUser() int {
	return room.CurrentUser
}

func (room *RoomData) GetRoomId() int {
	return room.id
}

func (room *RoomData) SendPersonalTableTip(u *user.User) {
	u.WriteMsg(&msg.G2C_PersonalTableTip{
		TableOwnerUserID:  room.CreateUser,                                                 //桌主 I D
		DrawCountLimit:    room.PkBase.TimerMgr.GetMaxPayCnt(),                           //局数限制
		DrawTimeLimit:     room.PkBase.TimerMgr.GetTimeLimit(),                           //时间限制
		PlayCount:         room.PkBase.TimerMgr.GetPlayCount(),                           //已玩局数
		PlayTime:          int(room.PkBase.TimerMgr.GetCreatrTime() - time.Now().Unix()), //已玩时间
		CellScore:         room.CellScore,                                                  //游戏底分
		IniScore:          0,                                                               //room.IniSource,                                                //初始分数
		ServerID:          strconv.Itoa(room.id),                                           //房间编号
		IsJoinGame:        0,                                                               //是否参与游戏 todo  tagPersonalTableParameter
		IsGoldOrGameScore: room.IsGoldOrGameScore,                                          //金币场还是积分场 0 标识 金币场 1 标识 积分场
	})
}

func (room *RoomData) SendStatusReady(u *user.User) {
	StatusFree := &nn_tb_msg.G2C_TBNN_StatusFree{}

	StatusFree.CellScore = room.CellScore                                    //基础积分
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

func (room *RoomData) SendStatusPlay(u *user.User) {
	StatusPlay := &nn_tb_msg.G2C_TBNN_StatusPlay{}

	UserCnt := room.PkBase.UserMgr.GetMaxPlayerCnt()
	//游戏变量
	StatusPlay.BankerUser = room.BankerUser
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

//派发扑克
/*
func (room *RoomData) DispatchCardData(wCurrentUser int) int {

	return 0
}*/

func (room *RoomData) BeforeStartGame(UserCnt int) {
	room.GameStatus = GAME_NULL
	room.InitRoom(UserCnt)
}

func (room *RoomData) StartGameing() {
	room.GameStatus = GAME_START
	room.StartDispatchCard()
}


func (room *RoomData) AfterStartGame() {
	room.GameStatus = CALL_SCORE_TIMES
	room.CallScoreTimer = room.PkBase.AfterFunc(CALL_SCORE_TIME * time.Second, func() {
		if room.GameStatus != ADD_SCORE { // 超时叫分结束
			room.CallScoreEnd()
		}
	})

	room.CallScoreTimer.Stop()
	
	//检查自摸
	//room.CheckZiMo()
	//通知客户端开始了
	//room.SendGameStart()
}

func (room *RoomData) InitRoom(UserCnt int) {
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

}

func (r *RoomData) GetOneCard() int  { // 从牌堆取出一张
	r.LeftCardCount -= 1
	return r.RepertoryCard[r.LeftCardCount]
}

func (room *RoomData) StartDispatchCard() {
	log.Debug("start dispatch card")

	userMgr := room.PkBase.UserMgr
	gameLogic := room.PkBase.LogicMgr

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, cost.US_PLAYING)
	})

	gameLogic.RandCardList(room.RepertoryCard, GetCardByIdx(room.ConfigIdx))

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

func (room *RoomData) SendGameStart() {
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
func (room *RoomData) NormalEnd() {
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
func (room *RoomData) DismissEnd() {
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

// 设置底分
func (room *RoomData) SetCellScore(cellScore int) {
	room.CellScore = cellScore
}

// 设置倍数
func (r *RoomData) SetScoreTimes(scoreTimes int) {
	r.ScoreTimes = scoreTimes
}

// 用户叫分(抢庄)
func (r *RoomData) CallScore(u *user.User, scoreTimes int) {
	log.Debug("add score times userChairId:%d, scoretimes:%d", u.ChairId, scoreTimes)
	r.CallScoreTimesMap[scoreTimes] = u
	maxScoreTimes := 0
	for s, _ := range r.CallScoreTimesMap {
		if s > maxScoreTimes {
			maxScoreTimes = s
		}
	}
	r.BankerUser = r.CallScoreTimesMap[maxScoreTimes].ChairId
	r.ScoreTimes = maxScoreTimes

	if len(r.CallScoreTimesMap) == r.PlayerCount {
		//叫分结束
		r.CallScoreEnd()
	}
}

// 叫分结束
func (r * RoomData) CallScoreEnd()  {
	// 发回叫分结果
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(&nn_tb_msg.G2C_TBNN_CallBanker{
			CallBanker: r.BankerUser,
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
func (r *RoomData) AddScore(u *user.User, score int) {
	log.Debug("add score userChairId:%d, score:%d", u.ChairId, score)
	r.ScoreMap[u] += score

	if len(r.ScoreMap) == r.PlayerCount { //全加过加注结束
		r.AddScoreEnd()
	}
}

// 加注结束
func (r * RoomData) AddScoreEnd() {
	// 发回加注结果
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(&nn_tb_msg.G2C_TBNN_AddScore{
			AddScoreCount: r.ScoreMap[u],
		})
	})

	// 进入最后一张牌
	log.Debug("enter last card")
	r.GameStatus = SEND_LAST_CARD

	// 发最后一张牌
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
func (r *RoomData) EnterOpenCard()  {
	log.Debug("enter open card")
	r.GameStatus = OPEN_CARD
	// 亮牌超时
	r.OpenCardTimer = r.PkBase.AfterFunc(OPEN_CARD_TIME, func() { // 超时亮牌结束
		if r.GameStatus != CAL_SCORE {
			//r.OpenCardEnd()
			log.Debug("open card time out")
		}
	})
	r.OpenCardTimer.Stop()
}

// 亮牌
func (r *RoomData) OpenCard(u *user.User, cardData []int)  {
	r.OpenCardMap[u] = cardData
	if len(r.OpenCardMap) == r.PlayerCount { // 全亮过
		r.OpenCardEnd() // 亮牌结束
	}
}

// 亮牌结束
func (r *RoomData) OpenCardEnd()  {
	// 比牌 发回比牌结果
}

