package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	. "mj/common/cost"
	"mj/common/msg/pk_sss_msg"

	"mj/gameServer/common/pk"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

// 游戏状态
const (
	GAME_START = 1002 // 游戏开始
)

func NewDataMgr(id int, uid int64, ConfigIdx int, name string, temp *base.GameServiceOption, base *SSS_Entry) *sss_data_mgr {
	d := new(sss_data_mgr)
	d.RoomData = pk_base.NewDataMgr(id, uid, ConfigIdx, name, temp, base.Entry_base)
	return d
}

type sss_data_mgr struct {
	*pk_base.RoomData

	//游戏变量

	bCardData               []int                  //牌的总数
	m_bUserCardData         map[*user.User][]int   //玩家扑克
	m_bSegmentCard          map[*user.User][][]int //分段扑克
	m_bFinishSegment        []int                  //完成分段
	m_bShowCardCount        int                    //摊牌数目
	m_bCompleteCompareCount int                    //完成比较
	m_bOverTime             []int                  //摊牌超时
	m_bUserLeft             []int                  //玩家强退

	SpecialTypeTable map[*user.User]bool  //是否特殊牌型
	Dragon           map[*user.User]bool  //是否倒水
	m_nPlayerCount   int                  //实际玩家人数
	CbResult         map[*user.User][]int //每一道的道数
	cbSpecialResult  map[*user.User]int   //特殊牌型的道数

	LeftCardCount int                  //剩下拍的数量
	OpenCardMap   map[*user.User][]int //摊牌数据
	//比较结果

	m_bCompareResult        map[*user.User][]int //每一道比较结果
	m_bShootState           [][]*user.User       //打枪(0赢的玩家,1输的玩家)
	m_bThreeKillResult      []int                //全垒打加减分
	m_bToltalWinDaoShu      []int                //总共道数
	m_bCompareDouble        map[*user.User]int   //打枪的道数
	m_bSpecialCompareResult map[*user.User]int   //特殊牌型比较结果
	m_lGameScore            []int                //游戏积分
	m_nXShoot               int                  //几家打枪

	// 游戏状态

	GameStatus        int
	BtCardSpecialData []int
}

func (room *sss_data_mgr) InitRoom(UserCnt int) {
	//初始化
	log.Debug("初始化房间")

	room.CbResult = make(map[*user.User][]int, UserCnt)
	room.PlayerCount = UserCnt
	room.m_bSegmentCard = make(map[*user.User][][]int, UserCnt)
	room.bCardData = make([]int, room.GetCfg().MaxRepertory) //牌堆
	room.OpenCardMap = make(map[*user.User][]int, UserCnt)
	room.Dragon = make(map[*user.User]bool, UserCnt)
	room.SpecialTypeTable = make(map[*user.User]bool, UserCnt)
	room.m_bUserCardData = make(map[*user.User][]int, UserCnt)
	room.m_bCompareDouble = make(map[*user.User]int, UserCnt)
	room.m_bCompareResult = make(map[*user.User][]int, UserCnt)
	room.m_bShootState = make([][]*user.User, UserCnt)
	room.m_bSpecialCompareResult = make(map[*user.User]int, UserCnt)
	room.m_bThreeKillResult = make([]int, UserCnt)
	room.m_bToltalWinDaoShu = make([]int, UserCnt)
	room.m_lGameScore = make([]int, UserCnt)
	room.m_nXShoot = 0

	room.LeftCardCount = room.GetCfg().MaxRepertory
}
func (room *sss_data_mgr) ComputeChOut() {
	userMgr := room.PkBase.UserMgr
	gameLogic := room.PkBase.LogicMgr

	userMgr.ForEachUser(func(u *user.User) {
		if room.SpecialTypeTable[u] == false {
			ResultTemp := make([]int, 3)
			bCardData := make([]int, 5)
			util.DeepCopy(&bCardData, &room.m_bSegmentCard[u][2])
			tagCardTypeHou := gameLogic.GetType(bCardData, len(bCardData))
			if tagCardTypeHou.BStraightFlush {
				ResultTemp[2] = 5
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[2])
			}
			bCardDataHouZ := make([]int, 5)
			util.DeepCopy(&bCardDataHouZ, &room.m_bSegmentCard[u][1])
			tagCardTypeHouZ := gameLogic.GetType(bCardDataHouZ, len(bCardData))
			if tagCardTypeHouZ.BStraightFlush {
				ResultTemp[1] = 10
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[1])
			}

			log.Debug("%d", room.CbResult[u])
			//后敦炸弹
			if CT_FIVE_FOUR_ONE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				ResultTemp[2] = 4
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[2])
			}
			//中敦炸弹
			if CT_FIVE_FOUR_ONE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				ResultTemp[1] = 8
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[1])
			}
			//后敦葫芦
			if CT_FIVE_THREE_DEOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				ResultTemp[2] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[2])
			}
			//中敦葫芦
			if CT_FIVE_THREE_DEOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				ResultTemp[1] = 2
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[1])
			}
			//后墩同花
			if CT_FIVE_FLUSH == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				ResultTemp[2] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[2])
			}
			//中墩同花
			if CT_FIVE_FLUSH == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				ResultTemp[1] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[1])
			}
			//后墩顺子
			if CT_FIVE_MIXED_FLUSH_NO_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) ||
				CT_FIVE_MIXED_FLUSH_FIRST_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) ||
				CT_FIVE_MIXED_FLUSH_BACK_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				ResultTemp[2] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[2])
			}
			//中墩顺子
			if CT_FIVE_MIXED_FLUSH_NO_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) ||
				CT_FIVE_MIXED_FLUSH_FIRST_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) ||
				CT_FIVE_MIXED_FLUSH_BACK_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				ResultTemp[1] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[1])
			}
			//后敦三张
			if CT_THREE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				ResultTemp[2] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[2])
			}
			//中敦三张
			if CT_THREE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				ResultTemp[1] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[1])
			}
			//前敦三张
			if CT_THREE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][0], 3, room.BtCardSpecialData) {
				ResultTemp[0] = 3
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[0])
			}
			//后敦两对
			if CT_FIVE_TWO_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				ResultTemp[2] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[2])
			}
			//中敦两对
			if CT_FIVE_TWO_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				ResultTemp[1] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[1])
			}
			//后敦一对
			if CT_ONE_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				ResultTemp[2] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[2])
			}
			//中敦一对
			if CT_ONE_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				ResultTemp[1] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[1])
			}
			//前敦一对
			if CT_ONE_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][0], 3, room.BtCardSpecialData) {
				ResultTemp[0] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[0])
			}
			//后敦散牌
			if CT_SINGLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				ResultTemp[2] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[2])
			}
			//中敦散牌
			if CT_SINGLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				ResultTemp[1] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[1])
			}
			//前敦散牌
			if CT_SINGLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][0], 3, room.BtCardSpecialData) {
				ResultTemp[0] = 1
				room.CbResult[u] = append(room.CbResult[u], ResultTemp[0])
			}
		} else {
			//至尊清龙
			if room.cbSpecialResult[u] == 0 && CT_THIRTEEN_FLUSH == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 104
			}
			//一条龙
			if room.cbSpecialResult[u] == 0 && CT_THIRTEEN == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 52
			}
			//十二皇族
			if room.cbSpecialResult[u] == 0 && CT_TWELVE_KING == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 24
			}
			//三同花顺
			if room.cbSpecialResult[u] == 0 && CT_THREE_STRAIGHTFLUSH == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 36
			}
			//三分天下
			if room.cbSpecialResult[u] == 0 && CT_THREE_BOMB == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 32
			}
			//全大
			if room.cbSpecialResult[u] == 0 && CT_ALL_BIG == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 10
			}
			//全小
			if room.cbSpecialResult[u] == 0 && CT_ALL_SMALL == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 10
			}
			//凑一色
			if room.cbSpecialResult[u] == 0 && CT_SAME_COLOR == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 10
			}
			//四套冲三
			if room.cbSpecialResult[u] == 0 && CT_FOUR_THREESAME == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 16
			}
			//五对冲三
			if room.cbSpecialResult[u] == 0 && CT_FIVEPAIR_THREE == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 5
			}
			//六对半
			if room.cbSpecialResult[u] == 0 && CT_SIXPAIR == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				//后敦炸弹 中敦炸弹
				if CT_FIVE_FOUR_ONE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) ||
					CT_FIVE_FOUR_ONE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
					room.cbSpecialResult[u] = 8
				} else {
					room.cbSpecialResult[u] = 6
				}
			}
			//三顺子
			if room.cbSpecialResult[u] == 0 && CT_THREE_STRAIGHT == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				//后敦同花顺 中敦同花顺

				tagCardTypeHou := new(pk.TagAnalyseType) //后敦同花顺

				bCardData := make([]int, 5)
				util.DeepCopy(&bCardData, &room.m_bSegmentCard[u][2])

				tagCardTypeHou = gameLogic.GetType(bCardData, 5)

				tagCardTypezhong := new(pk.TagAnalyseType) //中敦同花顺
				bCardDatazhong := make([]int, 5)
				util.DeepCopy(&bCardDatazhong, &room.m_bSegmentCard[u][1])

				tagCardTypezhong = gameLogic.GetType(bCardDatazhong, 5)

				if tagCardTypeHou.BStraightFlush || tagCardTypezhong.BStraightFlush {
					room.cbSpecialResult[u] = 10
				} else {
					room.cbSpecialResult[u] = 6
				}
			}
			//三同花
			if room.cbSpecialResult[u] == 0 && CT_THREE_FLUSH == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				//后敦同花顺 中敦同花顺

				tagCardTypeHou := new(pk.TagAnalyseType) //后敦同花顺

				bCardData := make([]int, 5)
				util.DeepCopy(&bCardData, &room.m_bSegmentCard[u][2])

				tagCardTypeHou = gameLogic.GetType(bCardData, 5)

				tagCardTypezhong := new(pk.TagAnalyseType) //中敦同花顺
				bCardDatazhong := make([]int, 5)
				util.DeepCopy(&bCardDatazhong, &room.m_bSegmentCard[u][1])

				tagCardTypezhong = gameLogic.GetType(bCardDatazhong, 5)

				if tagCardTypeHou.BStraightFlush || tagCardTypezhong.BStraightFlush {
					room.cbSpecialResult[u] = 10
				} else {
					room.cbSpecialResult[u] = 6
				}
			}
		}
	})

}

func (room *sss_data_mgr) ComputeResult() {
	shootCount := make(map[*user.User]map[*user.User]int, 10)
	m_nXShoot := 0
	WinNum := make(map[*user.User]int, 4)

	gameLogic := room.PkBase.LogicMgr
	userMgrW := room.PkBase.UserMgr
	userMgrW.ForEachUser(func(uW *user.User) {
		lWinDaoShu := 0
		userMgrN := room.PkBase.UserMgr
		userMgrN.ForEachUser(func(uN *user.User) {
			if uW != uN {
				if room.Dragon[uW] && room.Dragon[uN] == false { ///<一家倒水一家不倒水
					if room.SpecialTypeTable[uN] == false { ///<不等于特殊牌型

						lWinDaoShu -= room.CbResult[uN][0]
						room.m_bCompareResult[uW][0] -= room.CbResult[uN][0]

						lWinDaoShu -= room.CbResult[uN][1]
						room.m_bCompareResult[uW][1] -= room.CbResult[uN][1]

						lWinDaoShu -= room.CbResult[uN][2]
						room.m_bCompareResult[uW][2] -= room.CbResult[uN][2]

					} else {
						lWinDaoShu -= room.cbSpecialResult[uN]
						room.m_bSpecialCompareResult[uW] -= room.cbSpecialResult[uN]
					}
				} else if room.Dragon[uW] == false && room.Dragon[uN] { ///<一家不倒水，一家倒水
					if room.SpecialTypeTable[uW] == false { ///<不等于特殊牌型
						lWinDaoShu += room.CbResult[uN][0]
						room.m_bCompareResult[uW][0] += room.CbResult[uN][0]

						lWinDaoShu += room.CbResult[uN][1]
						room.m_bCompareResult[uW][1] += room.CbResult[uN][1]

						lWinDaoShu += room.CbResult[uN][2]
						room.m_bCompareResult[uW][2] += room.CbResult[uN][2]

						WinNum[uW]++
					} else {
						lWinDaoShu += room.cbSpecialResult[uW]
						room.m_bSpecialCompareResult[uW] += room.cbSpecialResult[uW]
					}
				} else if room.Dragon[uW] == false && room.Dragon[uN] == false {

					if room.SpecialTypeTable[uW] == false && room.SpecialTypeTable[uN] == false {

						if gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][0], room.m_bSegmentCard[uW][0], 3, 3, true) {
							lWinDaoShu += room.CbResult[uW][0]
							room.m_bCompareResult[uW][0] += room.CbResult[uW][0]
						} else {
							lWinDaoShu -= room.CbResult[uN][0]
							room.m_bCompareResult[uW][0] -= room.CbResult[uN][0]
						}

						if gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][1], room.m_bSegmentCard[uW][1], 5, 5, true) {
							lWinDaoShu += room.CbResult[uW][1]
							room.m_bCompareResult[uW][1] += room.CbResult[uW][1]
						} else {
							lWinDaoShu -= room.CbResult[uN][1]
							room.m_bCompareResult[uW][1] -= room.CbResult[uN][1]
						}
						if gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][2], room.m_bSegmentCard[uW][2], 5, 5, true) {
							lWinDaoShu += room.CbResult[uW][2]
							room.m_bCompareResult[uW][2] += room.CbResult[uW][2]
						} else {
							lWinDaoShu -= room.CbResult[uN][2]
							room.m_bCompareResult[uW][2] -= room.CbResult[uN][2]
						}

						if gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][0], room.m_bSegmentCard[uW][0], 3, 3, true) &&
							gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][1], room.m_bSegmentCard[uW][1], 5, 5, true) &&
							gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][2], room.m_bSegmentCard[uW][2], 5, 5, true) {
							room.m_bCompareDouble[uW] += lWinDaoShu
							lWinDaoShu *= 2

							room.m_bShootState[m_nXShoot][0] = uW ///<赢的
							room.m_bShootState[m_nXShoot][1] = uN ///<输的
							m_nXShoot++
							WinNum[uW]++
						} else if !gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][0], room.m_bSegmentCard[uW][0], 3, 3, true) &&
							!gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][1], room.m_bSegmentCard[uW][1], 5, 5, true) &&
							!gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][2], room.m_bSegmentCard[uW][2], 5, 5, true) {

							room.m_bCompareDouble[uW] += lWinDaoShu
							lWinDaoShu *= 2

							shootCount[uW][uN] = lWinDaoShu
						}
					} else if room.SpecialTypeTable[uW] == true && room.SpecialTypeTable[uN] == false {
						WinNum[uW]++ //add
						lWinDaoShu += room.cbSpecialResult[uW]
						room.m_bSpecialCompareResult[uW] += room.cbSpecialResult[uW]
					} else if room.SpecialTypeTable[uW] == true && room.SpecialTypeTable[uN] == true {
						if gameLogic.GetSSSCardType(room.m_bUserCardData[uW], 13, room.BtCardSpecialData) > gameLogic.GetSSSCardType(room.m_bUserCardData[uN], 13, room.BtCardSpecialData) {
							WinNum[uW]++ //add
							lWinDaoShu += room.cbSpecialResult[uW]
							room.m_bSpecialCompareResult[uW] += room.cbSpecialResult[uW]
						} else if gameLogic.GetSSSCardType(room.m_bUserCardData[uW], 13, room.BtCardSpecialData) < gameLogic.GetSSSCardType(room.m_bUserCardData[uN], 13, room.BtCardSpecialData) {
							lWinDaoShu -= room.cbSpecialResult[uN]
							room.m_bSpecialCompareResult[uW] -= room.cbSpecialResult[uN]
						}
					} else if room.SpecialTypeTable[uW] == false && room.SpecialTypeTable[uN] == true {
						lWinDaoShu -= room.cbSpecialResult[uN]
						room.m_bSpecialCompareResult[uW] -= room.cbSpecialResult[uN]
					}
				}
			}
		})
	})
}

//正常结束房间
func (room *sss_data_mgr) NormalEnd() {

	room.ComputeChOut()
	//room.ComputeResult()

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

func (room *sss_data_mgr) AfertEnd(a bool) {

}

//解散结束
func (room *sss_data_mgr) DismissEnd() {
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
func (room *sss_data_mgr) BeforeStartGame(UserCnt int) {
	room.GameStatus = GAME_START
	room.InitRoom(UserCnt)
}

func (room *sss_data_mgr) StartGameing() {
	room.StartDispatchCard()
}

func (room *sss_data_mgr) GetOneCard() int { // 从牌堆取出一张
	room.LeftCardCount -= 1
	return room.bCardData[room.LeftCardCount]
}
func (room *sss_data_mgr) StartDispatchCard() {
	log.Debug("begin start game sss")
	userMgr := room.PkBase.UserMgr
	gameLogic := room.PkBase.LogicMgr

	gameLogic.RandCardList(room.bCardData, pk_base.GetCardByIdx(room.ConfigIdx))

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, US_PLAYING)
	})

	userMgr.ForEachUser(func(u *user.User) {

		for i := 0; i < pk_base.GetCfg(pk_base.IDX_SSS).MaxCount; i++ {
			room.m_bUserCardData[u] = append(room.m_bUserCardData[u], room.GetOneCard())

		}
	})

	userMgr.ForEachUser(func(u *user.User) {
		SendCard := &pk_sss_msg.G2C_SSS_SendCard{}
		SendCard.CardData = room.m_bUserCardData[u]
		SendCard.CellScore = room.CellScore
		u.WriteMsg(SendCard)
	})
}
func (room *sss_data_mgr) AfterStartGame() {

}

//玩家摊牌
func (room *sss_data_mgr) ShowSSSCard(u *user.User, bDragon bool, bSpecialType bool, btSpecialData []int, bFrontCard []int, bMidCard []int, bBackCard []int) {

	room.SpecialTypeTable[u] = bSpecialType
	room.Dragon[u] = bDragon
	room.m_bSegmentCard[u] = append(room.m_bSegmentCard[u], bFrontCard, bMidCard, bBackCard)

	btSpecialDataTemp := make([]int, 13)

	if bSpecialType {
		util.DeepCopy(&btSpecialDataTemp, &btSpecialData)
	}

	// 广播摊牌
	userMgr := room.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(&pk_sss_msg.G2C_SSS_Open_Card{
			CurrentUser:    u.ChairId,
			FrontCard:      bFrontCard,
			MidCard:        bMidCard,
			BackCard:       bBackCard,
			CanSeeShowCard: false,
			SpecialType:    bSpecialType,
			SpecialData:    btSpecialDataTemp,
			ShowUser:       u.ChairId,
			Dragon:         bDragon,
		})
	})

	room.OpenCardMap[u] = bFrontCard
	log.Debug("%d cccccc", len(room.OpenCardMap))
	if len(room.OpenCardMap) == room.PlayerCount { //已全摊
		// 游戏结束
		//userMgr.ForEachUser(func(u *user.User) {
		room.PkBase.OnEventGameConclude(u.ChairId, u, GER_NORMAL)
		//})
	}

}
