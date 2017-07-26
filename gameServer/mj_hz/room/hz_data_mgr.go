package room

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/mj_hz_msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
)

func NewHZDataMgr(id int, uid int64, configIdx int, name string, temp *base.GameServiceOption, base *hz_entry) *hz_data {
	d := new(hz_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base)
	return d
}

type hz_data struct {
	*mj_base.RoomData
	ZhuaHuaCnt   int //扎花个数
	ZhuaHuaScore int //扎花分数
}

func (room *hz_data) AfterStartGame() {
	//检查自摸
	room.CheckZiMo()
	//检测起手杠牌
	room.CheckGameStartGang()
	//通知客户端开始了
	room.SendGameStart()
}

func (room *hz_data) SendGameStart() {
	//构造变量
	GameStart := &mj_hz_msg.G2C_HZMG_GameStart{}
	GameStart.BankerUser = room.BankerUser
	GameStart.SiceCount = room.SiceCount
	GameStart.HeapHead = room.HeapHead
	GameStart.HeapTail = room.HeapTail
	GameStart.HeapCardInfo = room.HeapCardInfo
	//发送数据
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		GameStart.UserAction = room.UserAction[u.ChairId]
		GameStart.CardData = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[u.ChairId])
		u.WriteMsg(GameStart)
	})
}

func (room *hz_data) OnZhuaHua(CenterUser int) (CardData []int, BuZhong []int) {
	count := room.ZhuaHuaCnt
	if count == 0 {
		return
	}

	isWin := false
	for chairId, v := range room.UserAction {
		if v&WIK_CHI_HU != 0 && chairId == CenterUser {
			isWin = true
		}
	}

	if !isWin {
		return
	}

	//抓花规则
	getInedx := [3]int{1, 5, 9}
	for i := 0; i < count; i++ {
		cardData := room.GetHeadCard()
		cardColor := room.MjBase.LogicMgr.GetCardColor(cardData)
		cardValue := room.MjBase.LogicMgr.GetCardValue(cardData)
		if cardColor == 3 {
			if cardValue >= 5 {
				//中发白
				temp := cardValue - 4
				if temp == getInedx[0] || temp == getInedx[1] || temp == getInedx[2] {
					CardData = append(CardData, cardData)
					room.ZhuaHuaScore++
				}
			}
		} else if cardColor >= 0 && cardColor <= 2 {
			if cardValue == getInedx[0] || cardValue == getInedx[1] || cardValue == getInedx[2] {
				CardData = append(CardData, cardData)
				room.ZhuaHuaScore++
			}
		}
		BuZhong = append(BuZhong, cardData)
	}
	return
}

//正常结束房间
func (room *hz_data) NormalEnd() {
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

	for i := range GameConclude.HandCardData {
		GameConclude.HandCardData[i] = make([]int, room.GetCfg().MaxCount)
	}

	GameConclude.SendCardData = room.SendCardData
	GameConclude.LeftUser = INVALID_CHAIR
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
}
