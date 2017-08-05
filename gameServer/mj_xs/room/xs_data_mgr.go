package room

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/mj_xs_msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/timer"
	"github.com/lovelly/leaf/util"
)

func NewXSDataMgr(id int, uid int64, configIdx int, name string, temp *base.GameServiceOption, base *xs_entry, set map[string]interface{}) *xs_data {
	d := new(xs_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base, set)
	return d
}

type xs_data struct {
	*mj_base.RoomData
	ZhuaHuaCnt      int                      //扎花个数
	ZhuaHuaScore    []int                    //扎花分数
	FengQaun        int                      //风圈
	IsFirst         bool                     //是否首发
	HuKindScore     [4][COUNT_KIND_SCORE]int //特殊胡牌分
	LianZhuang      int                      //连庄次数
	SumScore        [4]int                   //游戏总分
	FollowCardScore []int                    //跟牌得分
	HuKindType      []int                    //胡牌类型
}

func (room *xs_data) InitRoom(UserCnt int) {
	//初始化
	log.Debug("mjxs at InitRoom")
	room.RepertoryCard = make([]int, room.GetCfg().MaxRepertory)
	room.CardIndex = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.CardIndex[i] = make([]int, room.GetCfg().MaxIdx)
	}
	room.ChiHuKind = make([]int, UserCnt)
	room.ChiPengCount = make([]int, UserCnt)
	room.GangCard = make([]bool, UserCnt) //杠牌状态
	room.GangCount = make([]int, UserCnt)
	room.Ting = make([]bool, UserCnt)
	room.UserAction = make([]int, UserCnt)
	room.DiscardCard = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.DiscardCard[i] = make([]int, 60)
	}
	room.UserGangScore = make([]int, UserCnt)
	room.WeaveItemArray = make([][]*msg.WeaveItem, UserCnt)
	room.ChiHuRight = make([]int, UserCnt)
	room.HeapCardInfo = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.HeapCardInfo[i] = make([]int, 2)
	}
	room.OperateTime = make([]*timer.Timer, UserCnt)

	room.UserActionDone = false
	room.SendStatus = Not_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideGangUser = INVALID_CHAIR
	room.MinusLastCount = 0
	room.MinusHeadCount = room.GetCfg().MaxRepertory
	room.OutCardCount = 0

	//设置xs麻将牌数据
	room.EndLeftCount = 16
	room.ZhuaHuaScore = make([]int, room.MjBase.UserMgr.GetMaxPlayerCnt())
	room.FlowerCnt = [4]int{}
	room.BanCardCnt = [4][9]int{}
	room.BanUser = [4]int{}

	room.IsResponse = make([]bool, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.OperateCard[i] = make([]int, 60)
	}
	log.Debug("len1 OperateCard: %d %d", len(room.OperateCard), len(room.OperateCard[1]))
	room.PerformAction = make([]int, UserCnt)
}

func (room *xs_data) BeforeStartGame(UserCnt int) {
	log.Debug("###################### BeforeStartGame")
	room.InitRoom(UserCnt)
}

func (room *xs_data) AfterStartGame() {
	//检查自摸
	room.CheckZiMo()
	//通知客户端开始了
	room.SendGameStart()
}

//发送开始
func (room *xs_data) SendGameStart() {
	//构造变量
	GameStart := &mj_xs_msg.G2C_GameStart{}
	GameStart.BankerUser = room.BankerUser
	GameStart.SiceCount = room.SiceCount
	GameStart.SunWindCount = 0
	GameStart.LeftCardCount = room.GetLeftCard()
	GameStart.First = room.IsFirst
	GameStart.FengQuan = room.FengQaun
	GameStart.InitialBankerUser = room.BankerUser
	//发送数据
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		GameStart.UserAction = room.UserAction[u.ChairId]
		GameStart.CardData = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[u.ChairId])
		u.WriteMsg(GameStart)
	})
}

//发送操作结果
func (room *xs_data) SendOperateResult(u *user.User, wrave *msg.WeaveItem) {
	OperateResult := &mj_xs_msg.G2C_OperateResult{}
	OperateResult.ProvideUser = wrave.ProvideUser
	OperateResult.OperateCode = wrave.WeaveKind
	OperateResult.OperateCard = wrave.CenterCard
	if u != nil {
		OperateResult.OperateUser = u.ChairId
	} else {
		OperateResult.OperateUser = wrave.OperateUser
		OperateResult.ActionMask = wrave.ActionMask
	}
	room.MjBase.UserMgr.SendMsgAll(OperateResult)
}

//响应判断
func (room *xs_data) EstimateUserRespond(wCenterUser int, cbCenterCard int, EstimatKind int) bool {
	//变量定义
	bAroseAction := false
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	room.ClearStatus()
	//动作判断
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		//用户过滤
		if wCenterUser == u.ChairId || room.MjBase.UserMgr.IsTrustee(u.ChairId) {
			return
		}

		//出牌类型
		if EstimatKind == EstimatKind_OutCard {
			//吃碰判断
			if u.UserLimit&LimitPeng == 0 {
				//碰牌判断
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimatePengCard(room.CardIndex[u.ChairId], cbCenterCard)
			}

			//吃牌判断
			wEatUser := (wCenterUser + UserCnt - 1) % UserCnt
			if wEatUser == u.ChairId {
				room.UserAction[wEatUser] |= room.MjBase.LogicMgr.EstimateEatCard(room.CardIndex[u.ChairId], cbCenterCard)
			}

			//杠牌判断
			log.Debug(".room.LeftCardCount > room.EndLeftCount %v, %v", room.IsEnoughCard(), u.UserLimit&LimitGang)
			if room.IsEnoughCard() && u.UserLimit&LimitGang == 0 {
				room.UserAction[u.ChairId] |= room.MjBase.LogicMgr.EstimateGangCard(room.CardIndex[u.ChairId], cbCenterCard)
			}
		}

		if u.UserLimit|LimitChiHu == 0 {
			//吃胡判断
			hu, _ := room.MjBase.LogicMgr.AnalyseChiHuCard(room.CardIndex[u.ChairId], room.WeaveItemArray[u.ChairId], cbCenterCard)
			if hu {
				room.UserAction[u.ChairId] |= WIK_CHI_HU
			}
		}

		//结果判断
		if room.UserAction[u.ChairId] != WIK_NULL {
			bAroseAction = true
		}
	})

	//结果处理
	if bAroseAction {
		//设置变量
		room.ProvideUser = wCenterUser
		room.ProvideCard = cbCenterCard
		room.ResumeUser = room.CurrentUser
		room.CurrentUser = INVALID_CHAIR

		//发送提示
		room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
			log.Debug("########### EstimateUserRespond ActionMask %v ###########", room.UserAction[u.ChairId])
			if room.UserAction[u.ChairId] != WIK_NULL {
				u.WriteMsg(&mj_xs_msg.G2C_OperateNotify{
					ActionMask: room.UserAction[u.ChairId],
					ActionCard: room.ProvideCard,
					ResumeUser: room.ResumeUser,
				})
			}
		})
		return true
	}

	if room.GangStatus != WIK_GANERAL {
		room.GangOutCard = true
		room.GangStatus = WIK_GANERAL
		room.ProvideGangUser = INVALID_CHAIR
	} else {
		room.GangOutCard = false
	}

	return false
}

//正常结束房间
func (room *xs_data) NormalEnd(cbReason int) {
	//变量定义
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	GameConclude := &mj_xs_msg.G2C_GameConclude{}
	GameConclude.ChiHuKind = make([]int, UserCnt)
	GameConclude.CardCount = make([]int, UserCnt)
	GameConclude.CardData = make([][]int, UserCnt)
	GameConclude.GameScore = make([]int, UserCnt)
	GameConclude.ChiHuRight = make([]int, UserCnt)

	for i := range GameConclude.CardData {
		GameConclude.CardData[i] = make([]int, room.GetCfg().MaxCount)
	}

	//结束信息
	for i := 0; i < UserCnt; i++ {
		GameConclude.ChiHuKind[i] = room.ChiHuKind[i]
		//权位过滤
		if room.ChiHuKind[i] == WIK_CHI_HU {
			room.FiltrateRight(i, &room.ChiHuRight[i])
			GameConclude.ChiHuRight[i] = room.ChiHuRight[i]
		}
		GameConclude.CardData[i] = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[i])
		GameConclude.CardCount[i] = len(GameConclude.CardData[i])
	}

	//计算胡牌输赢分
	UserGameScore := make([]int, UserCnt)
	room.CalHuPaiScore(UserGameScore)

	//积分变量
	ScoreInfoArray := make([]*msg.TagScoreInfo, UserCnt)

	GameConclude.ProvideUser = room.ProvideUser
	GameConclude.ProvideCard = room.ProvideCard

	//统计积分
	DetailScore := make([]int, room.MjBase.UserMgr.GetMaxPlayerCnt())
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		if u.Status != US_PLAYING {
			return
		}
		GameConclude.GameScore[u.ChairId] = UserGameScore[u.ChairId]
		//胡牌分算完后再加上杠的输赢分就是玩家本轮最终输赢分
		GameConclude.GameScore[u.ChairId] += room.UserGangScore[u.ChairId]
		GameConclude.GangScore[u.ChairId] = room.UserGangScore[u.ChairId]

		ScoreInfoArray[u.ChairId] = &msg.TagScoreInfo{}
		ScoreInfoArray[u.ChairId].Score = GameConclude.GameScore[u.ChairId]
		if ScoreInfoArray[u.ChairId].Score > 0 {
			ScoreInfoArray[u.ChairId].Type = SCORE_TYPE_WIN
		} else {
			ScoreInfoArray[u.ChairId].Type = SCORE_TYPE_LOSE
		}

		//历史积分
		room.HistorySe.AllScore[u.ChairId] += GameConclude.GameScore[u.ChairId]
		DetailScore[u.ChairId] = GameConclude.GameScore[u.ChairId]
	})

	room.HistorySe.DetailScore = append(room.HistorySe.DetailScore, DetailScore)
	GameConclude.Reason = cbReason
	GameConclude.AllScore = room.HistorySe.AllScore
	GameConclude.DetailScore = room.HistorySe.DetailScore
	GameConclude.AllScore = room.HistorySe.AllScore
	GameConclude.DetailScore = room.HistorySe.DetailScore
	//发送数据
	room.MjBase.UserMgr.SendMsgAll(GameConclude)

	//写入积分 todo
	room.MjBase.UserMgr.WriteTableScore(ScoreInfoArray, room.MjBase.UserMgr.GetMaxPlayerCnt(), HZMJ_CHANGE_SOURCE)
}

//算分
func (room *xs_data) CalHuPaiScore(EndScore []int) {
	//CellScore := room.Source
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	UserScore := make([]int, UserCnt) //玩家手上分
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		if u.Status != US_PLAYING {
			return
		}
		UserScore[u.ChairId] = int(u.Score)
	})

	var WinUser []int
	WinCount := 0
	for i := 0; i < UserCnt; i++ {
		if WIK_CHI_HU == room.ChiHuKind[i] {
			WinUser = append(WinUser, i)
			room.CurrentUser = i
			room.SpecialCardScore(i)
			WinCount++
		}
	}

	//进行抓花
	ZhongCard, BuZhong := room.OnZhuaHua(WinUser)
	log.Debug("========================= ZhongCard:%d,BuZhong:%d", len(ZhongCard), len(BuZhong))
	if WinCount > 0 {

		//总分
		room.SumGameScore(WinUser)

		//连庄
		if WinCount > 1 {
			//一炮多响,庄家当庄
			var Zhuang bool
			for _, v := range WinUser {
				if room.BankerUser == v {
					Zhuang = true
				}
			}
			if Zhuang == false {
				room.BankerUser = room.BankerUser + 1
				room.LianZhuang = 1 //连庄次数
			} else {
				room.LianZhuang += 1 //连庄次数
			}
		} else {
			if WinUser[0] == room.BankerUser {
				room.BankerUser = room.BankerUser
				room.LianZhuang += 1 //连庄次数
			} else {
				room.BankerUser += 1
				room.LianZhuang = 1 //连庄次数
			}
		}

		if room.BankerUser > 3 {
			room.BankerUser = 0
		}
	} else { //荒庄
		room.BankerUser = room.BankerUser
	}
}

//杠计分
func (room *xs_data) CallGangScore() {
	lcell := room.Source
	//暗杠得分
	if room.GangStatus == WIK_AN_GANG {
		room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
			if u.Status != US_PLAYING {
				return
			}
			if u.ChairId != room.CurrentUser {
				room.UserGangScore[u.ChairId] -= lcell
				room.UserGangScore[room.CurrentUser] += lcell
			}
		})
	}
}

//特殊胡牌算分规则
func (room *xs_data) SpecialCardScore(HuUserID int) {
	score := room.Source
	winScore := &room.HuKindScore[HuUserID]

	for k, v := range winScore {
		if v <= 0 {
			continue
		}

		switch k {
		case IDX_SUB_SCORE_ZPKZ:
			winScore[k] = 1 * score
		case IDX_SUB_SCORE_HDLZ:
			winScore[k] = 8 * score
		case IDX_SUB_SCORE_GSKH:
			winScore[k] = 4 * score
		case IDX_SUB_SCORE_HSKH:
			winScore[k] = 4 * score
		case IDX_SUB_SCORE_QYS:
			winScore[k] = 32 * score
		case IDX_SUB_SCORE_HYS:
			winScore[k] = 16 * score
		case IDX_SUB_SCORE_CYS:
			winScore[k] = 8 * score
		case IDX_SUB_SCORE_DSY:
			winScore[k] = 16 * score
		case IDX_SUB_SCORE_XSY:
			winScore[k] = 8 * score
		case IDX_SUB_SCORE_DDH:
			winScore[k] = 8 * score
		case IDX_SUB_SCORE_MQQ:
			winScore[k] = 4 * score
		case IDX_SUB_SCORE_BL:
			winScore[k] = 4 * score
		case IDX_SUB_SCORE_DH:
			winScore[k] = 4 * score
		case IDX_SUB_SCORE_TH:
			winScore[k] = 4 * score
		case IDX_SUB_SCORE_DDPH:
			winScore[k] = 1 * score
		case IDX_SUB_SCORE_WDD:
			winScore[k] = 8 * score
		case IDX_SUB_SCORE_MQBL:
			winScore[k] = 12 * score
		case IDX_SUB_SCORE_SANAK:
			winScore[k] = 4 * score
		case IDX_SUB_SCORE_SIAK:
			winScore[k] = 8 * score
		case IDX_SUB_SCORE_WUAK:
			winScore[k] = 16 * score
		case IDX_SUB_SCORE_ZM:
			winScore[k] = 1 * score
		case IDX_SUB_SCORE_QGH:
			winScore[k] = 4 * score
		case IDX_SUB_SCORE_WHZ:
			winScore[k] = 4 * score
		case IDX_SUB_SCORE_ZYS:
			winScore[k] = 16 * score
		case IDX_SUB_SCORE_ZPG:
			winScore[k] = 1 * score
		}
	}

}

//总得分计算和得分类型统计
func (room *xs_data) SumGameScore(WinUser []int) {
	log.Debug("总得分计算和得分类型统计 赢人：%d", len(WinUser))
	log.Debug("补花数：%v", room.FlowerCnt)

	score := room.Source
	UserCnt := room.MjBase.UserMgr.GetMaxPlayerCnt()
	for i := 0; i < UserCnt; i++ {
		playerScore := &room.HuKindScore[i]

		//暗杠
		if room.UserGangScore[i] > 0 {
			playerScore[IDX_SUB_SCORE_AG] = room.UserGangScore[i]
		}
		room.SumScore[i] += room.UserGangScore[i]

		//胜者
		winCnt := 0
		for k := range WinUser {
			if WinUser[k] == i {
				winCnt++
				break
			}
		}
		if winCnt == 0 {
			continue
		}

		nowCnt := 0
		tempScore := [COUNT_KIND_SCORE]int{}
		util.DeepCopy(&tempScore, playerScore)
		if i == room.ProvideUser && winCnt == 1 { //自摸情况
			for index := 0; index < UserCnt; index++ {
				if index == i {
					continue
				}
				nowCnt++

				//基础分
				playerScore[IDX_SUB_SCORE_JC] += 1 * score
				room.SumScore[index] -= 1 * score
				room.SumScore[i] += 1 * score
				log.Debug("基础分:%d,SumScore:%d", playerScore[IDX_SUB_SCORE_JC], room.SumScore[i])

				//胡牌
				testScore := 0 //todo,测试
				for j := IDX_SUB_SCORE_ZPKZ; j < COUNT_KIND_SCORE; j++ {
					if nowCnt > 1 {
						playerScore[j] += tempScore[j]
					}
					testScore += tempScore[j] //todo,测试
					room.SumScore[i] += tempScore[j]
					room.SumScore[index] -= tempScore[j]
				}
				log.Debug("胡牌分：%d", testScore)

				//连庄
				if index == 0 {
					if i == room.BankerUser { //庄W
						room.SumScore[index] -= room.LianZhuang * score
						playerScore[IDX_SUB_SCORE_LZ] += room.LianZhuang * score
						room.SumScore[room.BankerUser] += room.LianZhuang * score
					} else { // 边W
						playerScore[IDX_SUB_SCORE_LZ] = room.LianZhuang * score
						room.SumScore[room.ProvideUser] += room.LianZhuang * score
						room.SumScore[room.BankerUser] -= room.LianZhuang * score
					}
					log.Debug("连庄得分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_LZ], room.SumScore[i])
				}

				//补花得分
				if room.FlowerCnt[i] < 8 {
					playerScore[IDX_SUB_SCORE_HUA] += room.FlowerCnt[i] * score
					room.SumScore[index] -= room.FlowerCnt[i] * score
				} else { //八张花牌
					playerScore[IDX_SUB_SCORE_HUA] += 16 * score
					room.SumScore[index] -= 16 * score
				}

				//抓花分
				playerScore[IDX_SUB_SCORE_ZH] = room.ZhuaHuaScore[i] * score
				room.SumScore[index] -= room.ZhuaHuaScore[i] * score
				room.SumScore[i] += room.ZhuaHuaScore[i] * score
				log.Debug("抓花分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_ZH], room.SumScore[i])
			}
			room.SumScore[i] += playerScore[IDX_SUB_SCORE_HUA]

			log.Debug("自摸i:%d ,庄家：%d", i, room.BankerUser)
			log.Debug("补分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_HUA], room.SumScore[i])
		} else {
			//基础分
			playerScore[IDX_SUB_SCORE_JC] += 1 * score
			room.SumScore[room.ProvideUser] -= 1 * score
			room.SumScore[i] += 1 * score
			log.Debug("基础分:%d,SumScore:%d", playerScore[IDX_SUB_SCORE_JC], room.SumScore[i])

			//胡牌
			testScore := 0 //todo,测试
			for j := IDX_SUB_SCORE_ZPKZ; j < COUNT_KIND_SCORE; j++ {
				testScore += tempScore[j] //todo,测试
				room.SumScore[i] += playerScore[j]
				room.SumScore[room.ProvideUser] -= tempScore[j]
			}
			log.Debug("胡牌分：%d", testScore)

			//补花分
			if room.FlowerCnt[i] < 8 {
				playerScore[IDX_SUB_SCORE_HUA] = room.FlowerCnt[i] * score
			} else { //八张花牌
				playerScore[IDX_SUB_SCORE_HUA] = 16 * score
			}
			room.SumScore[i] += playerScore[IDX_SUB_SCORE_HUA]
			room.SumScore[room.ProvideUser] -= playerScore[IDX_SUB_SCORE_HUA]
			log.Debug("补花得分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_HUA], room.SumScore[i])

			//连庄
			if i == room.BankerUser { //庄W
				room.SumScore[room.ProvideUser] -= room.LianZhuang * score
				playerScore[IDX_SUB_SCORE_LZ] = room.LianZhuang * score
				room.SumScore[room.BankerUser] += room.LianZhuang * score
			} else if room.ProvideUser == room.BankerUser { // 边W
				playerScore[IDX_SUB_SCORE_LZ] = room.LianZhuang * score
				room.SumScore[room.ProvideUser] += room.LianZhuang * score
				room.SumScore[room.BankerUser] -= room.LianZhuang * score
			}
			log.Debug("i:%d ,庄家：%d", i, room.BankerUser)
			log.Debug("连庄得分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_LZ], room.SumScore[i])

			//抓花分
			playerScore[IDX_SUB_SCORE_ZH] = room.ZhuaHuaScore[i] * score
			room.SumScore[room.ProvideUser] -= room.ZhuaHuaScore[i] * score
			room.SumScore[i] += room.ZhuaHuaScore[i] * score
			log.Debug("抓花分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_ZH], room.SumScore[i])
		}

		//分饼
		if room.BankerUser == i {
			room.SumScore[i] += room.FollowCardScore[i] * score
		} else {
			playerScore[IDX_SUB_SCORE_FB] = room.FollowCardScore[i] * score
			room.SumScore[i] += room.FollowCardScore[i] * score
		}
		log.Debug("分饼分：%d SumScore:%d", playerScore[IDX_SUB_SCORE_FB], room.SumScore[i])
	}
	log.Debug("游戏总分：%d", room.SumScore)
}

//特殊胡牌类型及算分
func (room *xs_data) SpecialCardKind(TagAnalyseItem []*TagAnalyseItem, HuUserID int) {

	type1Cnt := 0
	type2Cnt := 0
	score := room.Source
	winScore := &room.HuKindScore[HuUserID]
	for _, v := range TagAnalyseItem {
		kind := 0
		kind = room.IsDaSanYuan(v) //大三元
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DSY] = 12 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("大三元 %d", winScore[IDX_SUB_SCORE_DSY])
		}
		kind = room.IsXiaoSanYuan(v) //小三元
		if kind > 0 {
			winScore[IDX_SUB_SCORE_XSY] = 6 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("小三元 %d", winScore[IDX_SUB_SCORE_XSY])
		}
		kind = room.IsHunYiSe(v) //混一色
		if kind > 0 {
			winScore[IDX_SUB_SCORE_CYS] = 6 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("混一色 %d", winScore[IDX_SUB_SCORE_CYS])
		}
		kind = room.IsQingYiSe(v, room.FlowerCnt) //清一色
		if kind > 0 {
			winScore[IDX_SUB_SCORE_QYS] = 24 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("清一色 %d", winScore[IDX_SUB_SCORE_QYS])
		}
		kind = room.IsHuaYiSe(v, room.FlowerCnt) //花一色
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HYS] = 12 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("花一色 %d", winScore[IDX_SUB_SCORE_HYS])
		}
		kind = room.IsGangKaiHua(v, room.WeaveItemArray[HuUserID]) //杠上开花
		if kind > 0 {
			winScore[IDX_SUB_SCORE_GSKH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("杠上开花 %d", winScore[IDX_SUB_SCORE_GSKH])
		}
		kind = room.IsHuaKaiHua(v) //花上开花
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HSKH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("花上开花 %d", winScore[IDX_SUB_SCORE_HSKH])
		}
		kind = room.IsBaiLiu(v, room.FlowerCnt) //佰六
		if kind > 0 {
			winScore[IDX_SUB_SCORE_BL] = 6 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("佰六 %d", winScore[IDX_SUB_SCORE_BL])
		}
		kind = room.IsMenQing(v) //门清
		if kind > 0 {
			winScore[IDX_SUB_SCORE_MQQ] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("门清 %d", winScore[IDX_SUB_SCORE_MQQ])
		}
		kind = room.IsMenQingBaiLiu(v, room.FlowerCnt) //门清佰六
		if kind > 0 {
			winScore[IDX_SUB_SCORE_BL] = 0
			winScore[IDX_SUB_SCORE_MQQ] = 0
			winScore[IDX_SUB_SCORE_MQBL] = 9 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("门清佰六 %d", winScore[IDX_SUB_SCORE_MQBL])
		}
		kind = room.IsHuWeiZhang(v) //尾单吊
		if kind > 0 {
			winScore[IDX_SUB_SCORE_WDD] = 6 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("尾单吊 %d", winScore[IDX_SUB_SCORE_WDD])
		}
		kind = room.IsJieTou(v) //截头
		if kind > 0 {
			winScore[IDX_SUB_SCORE_JT] = 1 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("截头 %d", winScore[IDX_SUB_SCORE_JT])
		}
		kind = room.IsKongXin(v) //空心
		if kind > 0 {
			winScore[IDX_SUB_SCORE_KX] = 1 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("空心 %d", winScore[IDX_SUB_SCORE_KX])
		}
		kind = room.IsDuiDuiHu(v) //对对胡
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DDH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("对对胡 %d", winScore[IDX_SUB_SCORE_DDH])
		}
		kind = room.IsTianHu(v) //天胡
		if kind > 0 {
			winScore[IDX_SUB_SCORE_TH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("天胡 %d", winScore[IDX_SUB_SCORE_TH])
		}
		kind = room.IsDiHu(v) //地胡
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DH] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("地胡 %d", winScore[IDX_SUB_SCORE_DH])
		}
		kind = room.IsHaiDiLaoYue(v) //海底捞针
		if kind > 0 {
			winScore[IDX_SUB_SCORE_HDLZ] = 3 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("海底捞针 %d", winScore[IDX_SUB_SCORE_HDLZ])
		}
		kind = room.IsAnKe(v) //暗刻
		if kind > 0 {
			winScore[IDX_SUB_SCORE_SANAK+kind/8] = 3 * (kind / 4) * score //2,8,16
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("%d暗刻(32,33,34) %d", IDX_SUB_SCORE_SANAK+kind/8, winScore[IDX_SUB_SCORE_SANAK+kind/8])
		}
		kind = room.IsDaSiXi(v) //大四喜
		if kind > 0 {
			winScore[IDX_SUB_SCORE_DSX] = 24 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("大四喜 %d", winScore[IDX_SUB_SCORE_DSX])
		}
		kind = room.IsXiaoSiXi(v) //小四喜
		if kind > 0 {
			winScore[IDX_SUB_SCORE_XSX] = 12 * score
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("小四喜 %d", winScore[IDX_SUB_SCORE_XSX])
		}
		//自摸
		kind = room.IsZiMo()
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_HDLZ] == 0 && winScore[IDX_SUB_SCORE_GSKH] == 0 && winScore[IDX_SUB_SCORE_HSKH] == 0 {
				winScore[IDX_SUB_SCORE_ZM] = 1 * score
				room.HuKindType = append(room.HuKindType, kind)
				log.Debug("自摸,%d", winScore[IDX_SUB_SCORE_ZM])
			}
		}
		//无花字
		kind = room.IsWuHuaZi(v, room.FlowerCnt)
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_BL] > 0 || winScore[IDX_SUB_SCORE_MQBL] > 0 {
				continue
			}
			winScore[IDX_SUB_SCORE_WHZ] = 3 * score
			log.Debug("无花字，%d", winScore[IDX_SUB_SCORE_WHZ])
		}
		//字一色
		kind = room.IsZiYiSe(v, room.FlowerCnt)
		if kind > 0 {
			winScore[IDX_SUB_SCORE_ZYS] = 12 * score
			log.Debug("字一色，%d", winScore[IDX_SUB_SCORE_ZYS])
		}
		kind, type1Cnt, type2Cnt = room.IsZiPaiGang(v) //字牌杠
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_DSX] > 0 || winScore[IDX_SUB_SCORE_XSX] > 0 {
				winScore[IDX_SUB_SCORE_ZPG] = type1Cnt * score
			} else if winScore[IDX_SUB_SCORE_DSY] > 0 || winScore[IDX_SUB_SCORE_XSY] > 0 {
				winScore[IDX_SUB_SCORE_ZPG] = type2Cnt * score
			} else {
				winScore[IDX_SUB_SCORE_ZPG] = (type2Cnt + type1Cnt) * score
			}
			type1Cnt = 0
			type2Cnt = 0
			log.Debug("字牌杠，%d", winScore[IDX_SUB_SCORE_ZPG])
		}
		kind, type1Cnt, type2Cnt = room.IsKeZi(v) //字牌刻字
		if kind > 0 {
			if winScore[IDX_SUB_SCORE_DSX] > 0 || winScore[IDX_SUB_SCORE_XSX] > 0 {
				winScore[IDX_SUB_SCORE_ZPKZ] = type1Cnt * score
			} else if winScore[IDX_SUB_SCORE_DSY] > 0 || winScore[IDX_SUB_SCORE_XSY] > 0 {
				winScore[IDX_SUB_SCORE_ZPKZ] = type2Cnt * score
			} else {
				winScore[IDX_SUB_SCORE_ZPKZ] = (type2Cnt + type1Cnt) * score
			}
			type1Cnt = 0
			type2Cnt = 0
			room.HuKindType = append(room.HuKindType, kind)
			log.Debug("字牌刻字 %d", winScore[IDX_SUB_SCORE_ZPKZ])
		}
	}
	//单吊
	if room.TingCnt[room.CurrentUser] == 1 {
		if room.CurrentUser == room.ProvideUser {
			winScore[IDX_SUB_SCORE_DDPH] = 1 * score
			room.HuKindType = append(room.HuKindType, IDX_SUB_SCORE_DDPH)
			log.Debug("单吊平胡,%d", winScore[IDX_SUB_SCORE_DDPH])
		} else {
			winScore[IDX_SUB_SCORE_DDZM] = 1 * score
			room.HuKindType = append(room.HuKindType, IDX_SUB_SCORE_DDZM)
			log.Debug("单吊自摸,%d", winScore[IDX_SUB_SCORE_DDZM])
		}
	}

}
