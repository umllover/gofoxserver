package NNBaseLogic

import (
	"strconv"
	"time"

	"mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/user"
	"mj/gameServer/db/model/base"
	"mj/common/msg/nn_tb_msg"
	"mj/gameServer/common/pk_base"

	"github.com/lovelly/leaf/log"

)

func NewDataMgr(id, uid, OriCardIdx int, name string, temp *base.GameServiceOption, base *NN_PK_base) *RoomData {
	r := new(RoomData)
	r.id = id
	if name == "" {
		r.Name = temp.GameName
	} else {
		r.Name = name
	}
	r.CreateUser = uid
	r.NNPkBase = base
	r.OriCardIdx = OriCardIdx
	return r
}

//当一张桌子理解
type RoomData struct {
	id         			int
	Name       			string //房间名字
	CreateUser 			int    //创建房间的人
	NNPkBase     			*NN_PK_base
	OriCardIdx 			int

	IsGoldOrGameScore 	int                //金币场还是积分场 0 标识 金币场 1 标识 积分场
	Password          	string             // 密码

	//游戏变量
	CardData			[][]int	//用户扑克
	HandCardData		[][]int	//桌面扑克
	//Banker				int		//庄家
	Count_All			int 	//是不是全部都点了 全点了直接动画
	Qiang				[]bool	//1抢 0不抢
	CellScore			int			//底分
	PlayCount			int			//游戏局数
	PlayerCount			int		//指定游戏人数，2-4


	BankerUser			int				//庄家用户
	FisrtCallUser		int				//始叫用户
	CurrentUser			int				//当前用户
	ExitScore			int64			//强退分数
	EscapeUserScore		[]int64        //逃跑玩家分数
	DynamicScore        int64              //总分


	//用户数据
	IsOpenCard			[]bool			//是否摊牌
	DynamicJoin			[]int           //动态加入
	PlayStatus			[]int			//游戏状态
	CallStatus			[]int			//叫庄状态
	OxCard				[]int			//牛牛数据
	TableScore			[]int64			//下注数目
	BuckleServiceCharge	[]bool			//收服务费

	//下注信息
	TurnMaxScore		[]int64			//最大下注
	MaxScoreTimes		int				//最大倍数

	//历史积分
	HistoryScores     []*pk_base.HistoryScore    //历史积分
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
		TableOwnerUserID:  room.CreateUser,                                               //桌主 I D
		DrawCountLimit:    room.NNPkBase.TimerMgr.GetCountLimit(),                          //局数限制
		DrawTimeLimit:     room.NNPkBase.TimerMgr.GetTimeLimit(),                           //时间限制
		PlayCount:         room.NNPkBase.TimerMgr.GetPlayCount(),                           //已玩局数
		PlayTime:          int(room.NNPkBase.TimerMgr.GetCreatrTime() - time.Now().Unix()), //已玩时间
		CellScore:         room.CellScore,                                                   //游戏底分
		IniScore:          0,//room.IniSource,                                                //初始分数
		ServerID:          strconv.Itoa(room.id),                                         //房间编号
		IsJoinGame:        0,                                                             //是否参与游戏 todo  tagPersonalTableParameter
		IsGoldOrGameScore: room.IsGoldOrGameScore,                                        //金币场还是积分场 0 标识 金币场 1 标识 积分场
	})
}


func (room *RoomData) SendStatusReady(u *user.User) {
	StatusFree := &nn_tb_msg.G2C_TBNN_StatusFree{}

	StatusFree.CellScore = room.CellScore                                     //基础积分
	StatusFree.TimeOutCard = room.NNPkBase.TimerMgr.GetTimeOutCard()         //出牌时间
	StatusFree.TimeOperateCard = room.NNPkBase.TimerMgr.GetTimeOperateCard() //操作时间
	StatusFree.TimeStartGame = room.NNPkBase.TimerMgr.GetCreatrTime()           //开始时间
	for _, v := range room.HistoryScores {
		StatusFree.TurnScore = append(StatusFree.TurnScore, v.TurnScore)
		StatusFree.CollectScore = append(StatusFree.TurnScore, v.CollectScore)
	}
	StatusFree.PlayerCount = room.NNPkBase.TimerMgr.GetPlayCount() //玩家人数
	StatusFree.CountLimit = room.NNPkBase.TimerMgr.GetCountLimit() //局数限制
	StatusFree.GameRoomName = room.Name

	u.WriteMsg(StatusFree)
}


func (room *RoomData) SendStatusPlay(u *user.User) {
	StatusPlay := &nn_tb_msg.G2C_TBNN_StatusPlay{}

	UserCnt := room.NNPkBase.UserMgr.GetMaxPlayerCnt()
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
func (room *RoomData) DispatchCardData(wCurrentUser int) int {

	return 0
}

func (room *RoomData) BeforeStartGame(UserCnt int) {
	room.InitRoom(UserCnt)
}

func (room *RoomData) StartGameing() {
	room.StartDispatchCard()
}

func (room *RoomData) AfterStartGame() {
	//检查自摸
	//room.CheckZiMo()
	//通知客户端开始了
	//room.SendGameStart()
}



func (room *RoomData) InitRoom(UserCnt int) {
	//初始化
	room.CardData = make([][]int, UserCnt)
	room.HandCardData = make([][]int, UserCnt)

	room.ExitScore = 0
	room.DynamicScore = 0
	room.BankerUser = cost.INVALID_CHAIR
	room.FisrtCallUser = cost.INVALID_CHAIR
	room.CurrentUser = cost.INVALID_CHAIR

	room.MaxScoreTimes = 0
	room.Count_All = 0
	room.Qiang = make([]bool, UserCnt)
	room.IsOpenCard = make([]bool, UserCnt)
	room.DynamicJoin = make([]int, UserCnt)
	room.TableScore = make([]int64, UserCnt)
	room.PlayStatus = make([]int, UserCnt)
	room.CallStatus = make([]int, UserCnt)
	room.EscapeUserScore = make([]int64, UserCnt)
	for i:=0;i<UserCnt;i++ {
		room.OxCard[i] = 0xFF
	}
	room.TurnMaxScore = make([]int64, UserCnt)

}

func (room *RoomData) StartDispatchCard() {
	log.Debug("begin start game tbnn")
	/*
	userMgr := room.NNPkBase.UserMgr
	gameLogic := room.NNPkBase.LogicMgr

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, US_PLAYING)
	})

	var minSice int
	UserCnt := userMgr.GetMaxPlayerCnt()
	room.SiceCount, minSice = room.GetSice()

	gameLogic.RandCardList(room.RepertoryCard, getCardByIdx(room.OriCardIdx))

	//红中可以当财神
	gameLogic.SetMagicIndex(gameLogic.SwitchToCardIndex(0x35))

	//分发扑克
	userMgr.ForEachUser(func(u *user.User) {
		room.LeftCardCount -= 1
		room.MinusHeadCount += 1
		for i := 0; i < MAX_COUNT-1; i++ {
			room.CardIndex[u.ChairId][i] = gameLogic.SwitchToCardIndex(room.RepertoryCard[room.LeftCardCount])
		}
	})

	OwnerUser, _ := userMgr.GetUserByUid(room.CreateUser)
	if room.BankerUser == INVALID_CHAIR && (room.MjBase.Temp.ServerType&GAME_GENRE_PERSONAL) != 0 { //房卡模式下先把庄家给房主
		if OwnerUser != nil {
			room.BankerUser = OwnerUser.ChairId
		} else {
			log.Error("get bamkerUser error at StartGame")
		}
	}

	if room.BankerUser == INVALID_CHAIR {
		room.BankerUser = util.RandInterval(0, UserCnt-1)
	}

	if room.BankerUser >= UserCnt {
		log.Error(" room.BankerUser >=UserCnt %d,  %d", room.BankerUser, UserCnt)
	}

	room.MinusHeadCount++
	room.SendCardData = room.RepertoryCard[room.LeftCardCount]
	room.LeftCardCount--

	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]++
	room.ProvideCard = room.SendCardData
	room.ProvideUser = room.BankerUser
	room.CurrentUser = room.BankerUser

	//堆立信息
	SiceCount := LOBYTE(room.SiceCount) + HIBYTE(room.SiceCount)
	TakeChairID := (room.BankerUser + SiceCount - 1) % UserCnt
	TakeCount := MAX_REPERTORY - room.LeftCardCount
	for i := 0; i < UserCnt; i++ {
		//计算数目
		var ValidCount int
		if i == 0 {
			ValidCount = HEAP_FULL_COUNT - room.HeapCardInfo[TakeChairID][1] - (minSice)*2
		} else {
			ValidCount = HEAP_FULL_COUNT - room.HeapCardInfo[TakeChairID][1]
		}

		RemoveCount := int(math.Min(float64(ValidCount), float64(TakeCount)))

		//提取扑克
		TakeCount -= RemoveCount
		if i == 0 {
			room.HeapCardInfo[TakeChairID][1] += RemoveCount
		} else {
			room.HeapCardInfo[TakeChairID][0] += RemoveCount
		}

		//完成判断
		if TakeCount == 0 {
			room.HeapHead = TakeChairID
			room.HeapTail = (room.BankerUser + SiceCount - 1) % UserCnt
			break
		}
		//切换索引
		TakeChairID = (TakeChairID + UserCnt - 1) % UserCnt
	}

	room.UserAction = make([]int, UserCnt)

	gangCardResult := &common.TagGangCardResult{}
	room.UserAction[room.BankerUser] |= gameLogic.AnalyseGangCard(room.CardIndex[room.BankerUser], nil, 0, gangCardResult)

	//胡牌判断
	chr := 0
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]--
	room.UserAction[room.BankerUser] |= gameLogic.AnalyseChiHuCard(room.CardIndex[room.BankerUser], []*msg.WeaveItem{}, room.SendCardData, chr, true)
	room.CardIndex[room.BankerUser][gameLogic.SwitchToCardIndex(room.SendCardData)]++
*/
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


