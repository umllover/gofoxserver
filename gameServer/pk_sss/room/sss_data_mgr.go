package room

import (
	"encoding/json"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	. "mj/common/cost"
	"mj/common/msg/pk_sss_msg"

	"mj/gameServer/common/pk"
	"mj/gameServer/db/model"

	//dbg "github.com/funny/debug"

	"math/rand"

	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

// 游戏状态
const (
	GAME_FREE       = 100 // 空闲
	GAME_SEND_CARD  = 101 // 发牌
	GAME_SETSEGMENT = 102 //组牌
	GAME_COMPARE    = 103 //比牌
	GAME_END        = 104 //结束

)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewDataMgr(info *model.CreateRoomInfo, uid int64, ConfigIdx int, name string, temp *base.GameServiceOption, base *SSS_Entry) *sss_data_mgr {
	d := new(sss_data_mgr)
	d.RoomData = pk_base.NewDataMgr(info.RoomId, uid, ConfigIdx, name, temp, base.Entry_base, info.OtherInfo)
	var setInfo sssOtherInfo
	if err := json.Unmarshal([]byte(info.OtherInfo), &setInfo); err == nil {
		d.wanFa = setInfo.WanFa
		d.jiaYiSe = setInfo.JiaYiSe
		d.jiaGongGong = setInfo.JiaGongGong
		d.jiaDaXiaoWan = setInfo.JiaDaXiaoWan
	}
	return d
}

type sssOtherInfo struct {
	WanFa        int
	JiaYiSe      bool
	JiaGongGong  bool
	JiaDaXiaoWan bool
}

type sss_data_mgr struct {
	*pk_base.RoomData

	//游戏变量
	wanFa        int
	jiaYiSe      bool
	jiaGongGong  bool
	jiaDaXiaoWan bool
	laiZi        []int
	publicCards  []int

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
	m_bThreeKillResult      map[*user.User]int   //全垒打加减分
	m_bToltalWinDaoShu      map[*user.User]int   //总共道数
	m_bCompareDouble        map[*user.User]int   //打枪的道数
	m_bSpecialCompareResult map[*user.User]int   //特殊牌型比较结果
	m_lGameScore            map[*user.User]int   //游戏积分
	m_nXShoot               int                  //几家打枪
	m_lCellScore            int                  //单元底分

	// 游戏状态

	GameStatus        int
	BtCardSpecialData []int
	AllResult         [][]int //每一局的结果
}

func (room *sss_data_mgr) InitRoom(UserCnt int) {
	//初始化
	log.Debug("初始化房间")

	room.cbSpecialResult = make(map[*user.User]int, UserCnt)
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
	room.m_bThreeKillResult = make(map[*user.User]int, UserCnt)
	room.m_bToltalWinDaoShu = make(map[*user.User]int, UserCnt)
	room.m_lGameScore = make(map[*user.User]int, UserCnt)
	room.m_nXShoot = 0
	room.BtCardSpecialData = make([]int, 13)
	room.LeftCardCount = room.GetCfg().MaxRepertory

	room.laiZi = make([]int, 0, 6)
	room.publicCards = make([]int, 0, 3)

	room.AllResult = make([][]int, room.PkBase.TimerMgr.GetMaxPlayCnt())
}
func (room *sss_data_mgr) ComputeChOut() {
	userMgr := room.PkBase.UserMgr
	gameLogic := room.PkBase.LogicMgr

	userMgr.ForEachUser(func(u *user.User) {
		room.CbResult[u] = make([]int, 3)
		if room.SpecialTypeTable[u] == false {
			//ResultTemp := make([]int, 3)
			bCardData := make([]int, 5)
			util.DeepCopy(&bCardData, &room.m_bSegmentCard[u][2])
			tagCardTypeHou := gameLogic.GetType(bCardData, len(bCardData))
			if tagCardTypeHou.BStraightFlush {
				//ResultTemp[2] = 5
				room.CbResult[u][2] = 5
			}
			bCardDataHouZ := make([]int, 5)
			util.DeepCopy(&bCardDataHouZ, &room.m_bSegmentCard[u][1])
			tagCardTypeHouZ := gameLogic.GetType(bCardDataHouZ, len(bCardData))
			if tagCardTypeHouZ.BStraightFlush {
				//ResultTemp[1] = 10
				room.CbResult[u][1] = 10
			}

			//后敦炸弹
			if CT_FIVE_FOUR_ONE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				log.Debug("后敦炸弹")
				//ResultTemp[2] = 4
				room.CbResult[u][2] = 4
			}
			//中敦炸弹
			if CT_FIVE_FOUR_ONE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				log.Debug("中敦炸弹")
				//ResultTemp[1] = 8
				room.CbResult[u][1] = 8
			}
			//后敦葫芦
			if CT_FIVE_THREE_DEOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				log.Debug("后敦葫芦")
				//ResultTemp[2] = 1
				room.CbResult[u][2] = 1
			}
			//中敦葫芦
			if CT_FIVE_THREE_DEOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				log.Debug("中敦葫芦")
				//ResultTemp[1] = 2
				room.CbResult[u][1] = 2
			}
			//后墩同花
			if CT_FIVE_FLUSH == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				log.Debug("后墩同花")
				//ResultTemp[2] = 1
				room.CbResult[u][2] = 1
			}
			//中墩同花
			if CT_FIVE_FLUSH == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				log.Debug("中墩同花")
				//ResultTemp[1] = 1
				room.CbResult[u][1] = 1
			}
			//后墩顺子
			if CT_FIVE_MIXED_FLUSH_NO_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) ||
				CT_FIVE_MIXED_FLUSH_FIRST_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) ||
				CT_FIVE_MIXED_FLUSH_BACK_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				log.Debug("后墩顺子")
				//ResultTemp[2] = 1
				room.CbResult[u][2] = 1
			}
			//中墩顺子
			if CT_FIVE_MIXED_FLUSH_NO_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) ||
				CT_FIVE_MIXED_FLUSH_FIRST_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) ||
				CT_FIVE_MIXED_FLUSH_BACK_A == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				log.Debug("中墩顺子")
				//ResultTemp[1] = 1
				room.CbResult[u][1] = 1
			}
			//后敦三张
			if CT_THREE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				log.Debug("后敦三张")
				//ResultTemp[2] = 1
				room.CbResult[u][2] = 1
			}
			//中敦三张
			if CT_THREE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				log.Debug("中敦三张")
				//ResultTemp[1] = 1
				room.CbResult[u][1] = 1
			}
			//前敦三张
			if CT_THREE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][0], 3, room.BtCardSpecialData) {
				log.Debug("前敦三张")
				//ResultTemp[0] = 3
				room.CbResult[u][0] = 3
			}
			//后敦两对
			if CT_FIVE_TWO_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				log.Debug("后敦两对")
				//ResultTemp[2] = 1
				room.CbResult[u][2] = 1
			}
			//中敦两对
			if CT_FIVE_TWO_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				log.Debug("中敦两对")
				//ResultTemp[1] = 1
				room.CbResult[u][1] = 1
			}
			//后敦一对
			if CT_ONE_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				log.Debug("后敦一对")
				//ResultTemp[2] = 1
				room.CbResult[u][2] = 1
			}
			//中敦一对
			if CT_ONE_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				log.Debug("中敦一对")
				//ResultTemp[1] = 1
				room.CbResult[u][1] = 1
			}
			//前敦一对
			if CT_ONE_DOUBLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][0], 3, room.BtCardSpecialData) {
				log.Debug("前敦一对")
				//ResultTemp[0] = 1
				room.CbResult[u][0] = 1
			}
			//后敦散牌
			if CT_SINGLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) {
				log.Debug("后敦散牌")
				//ResultTemp[2] = 1
				room.CbResult[u][2] = 1
			}
			//中敦散牌
			if CT_SINGLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
				log.Debug("中敦散牌")
				//ResultTemp[1] = 1
				room.CbResult[u][1] = 1
			}
			//前敦散牌
			if CT_SINGLE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][0], 3, room.BtCardSpecialData) {
				log.Debug("前敦散牌")
				//ResultTemp[0] = 1
				room.CbResult[u][0] = 1
			}
			//log.Debug("%d   zzzzzzzzzz", ResultTemp)

		} else {
			//至尊清龙
			if room.cbSpecialResult[u] == 0 && CT_THIRTEEN_FLUSH == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 104
				log.Debug("至尊清龙 %d", room.cbSpecialResult[u])
			}
			//一条龙
			if room.cbSpecialResult[u] == 0 && CT_THIRTEEN == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 52
				log.Debug("一条龙 %d", room.cbSpecialResult[u])
			}
			//十二皇族
			if room.cbSpecialResult[u] == 0 && CT_TWELVE_KING == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 24
				log.Debug("十二皇族 %d", room.cbSpecialResult[u])
			}
			//三同花顺
			if room.cbSpecialResult[u] == 0 && CT_THREE_STRAIGHTFLUSH == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 36
				log.Debug("三同花顺 %d", room.cbSpecialResult[u])
			}
			//三分天下
			if room.cbSpecialResult[u] == 0 && CT_THREE_BOMB == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 32
				log.Debug("三分天下 %d", room.cbSpecialResult[u])
			}
			//全大
			if room.cbSpecialResult[u] == 0 && CT_ALL_BIG == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 10
				log.Debug("全大 %d", room.cbSpecialResult[u])
			}
			//全小
			if room.cbSpecialResult[u] == 0 && CT_ALL_SMALL == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 10
				log.Debug("全小 %d", room.cbSpecialResult[u])
			}
			//凑一色
			if room.cbSpecialResult[u] == 0 && CT_SAME_COLOR == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 10
				log.Debug("凑一色 %d", room.cbSpecialResult[u])
			}
			//四套冲三
			if room.cbSpecialResult[u] == 0 && CT_FOUR_THREESAME == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 16
				log.Debug("四套冲三 %d", room.cbSpecialResult[u])
			}
			//五对冲三
			if room.cbSpecialResult[u] == 0 && CT_FIVEPAIR_THREE == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				room.cbSpecialResult[u] = 5
				log.Debug("五对冲三 %d", room.cbSpecialResult[u])
			}
			//六对半
			if room.cbSpecialResult[u] == 0 && CT_SIXPAIR == gameLogic.GetSSSCardType(room.m_bUserCardData[u], 13, room.BtCardSpecialData) {
				//后敦炸弹 中敦炸弹
				if CT_FIVE_FOUR_ONE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][2], 5, room.BtCardSpecialData) ||
					CT_FIVE_FOUR_ONE == gameLogic.GetSSSCardType(room.m_bSegmentCard[u][1], 5, room.BtCardSpecialData) {
					room.cbSpecialResult[u] = 8
					log.Debug("六对半 后敦炸弹 中敦炸弹 %d", room.cbSpecialResult[u])
				} else {
					room.cbSpecialResult[u] = 6
					log.Debug("六对半 %d", room.cbSpecialResult[u])
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
					log.Debug("三顺子 后敦同花顺 中敦同花顺 %d", room.cbSpecialResult[u])
				} else {
					room.cbSpecialResult[u] = 6
					log.Debug("三顺子 %d", room.cbSpecialResult[u])
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
					log.Debug("三同花 后敦同花顺 中敦同花顺 %d", room.cbSpecialResult[u])
				} else {
					room.cbSpecialResult[u] = 6
					log.Debug("三同花 %d", room.cbSpecialResult[u])
				}

			}
			log.Debug("cbSpecialResult %d", room.cbSpecialResult[u])
		}

	})

}

func (room *sss_data_mgr) ComputeResult() {
	shootCount := make(map[*user.User]map[*user.User]int, 10)
	room.m_nXShoot = 0
	WinNum := make(map[*user.User]int, 4)
	for i := range room.m_bShootState {
		room.m_bShootState[i] = make([]*user.User, 2)
	}
	gameLogic := room.PkBase.LogicMgr
	userMgrW := room.PkBase.UserMgr
	userMgrW.ForEachUser(func(u *user.User) {
		room.m_bCompareResult[u] = make([]int, 3)
	})
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

							room.m_bShootState[room.m_nXShoot][0] = uW ///<赢的
							room.m_bShootState[room.m_nXShoot][1] = uN ///<输的
							room.m_nXShoot++
							WinNum[uW]++
						} else if !gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][0], room.m_bSegmentCard[uW][0], 3, 3, true) &&
							!gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][1], room.m_bSegmentCard[uW][1], 5, 5, true) &&
							!gameLogic.CompareSSSCard(room.m_bSegmentCard[uN][2], room.m_bSegmentCard[uW][2], 5, 5, true) {

							room.m_bCompareDouble[uW] += lWinDaoShu
							lWinDaoShu *= 2
							shootCount[uW] = make(map[*user.User]int)
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
				room.m_lGameScore[uW] += lWinDaoShu * 2 //room.m_lCellScore
				room.m_bToltalWinDaoShu[uW] += lWinDaoShu
			}
		})

	})
	AllKillCount := 0
	///<下面判断是否全垒打在加减分
	userMgr := room.PkBase.UserMgr
	userMgrq := room.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		if WinNum[u] == 3 {
			userMgrq.ForEachUser(func(uN *user.User) {
				if u == uN {
					AllKillCount = room.m_bCompareDouble[u] * 2

					room.m_lGameScore[uN] += AllKillCount * 2 //m_lCellScore
					room.m_bToltalWinDaoShu[uN] += AllKillCount
					room.m_bThreeKillResult[uN] = AllKillCount
				} else {
					AllKillCount = 3                          //room.shootCount[j][i]
					room.m_lGameScore[uN] += AllKillCount * 2 //m_lCellScore
					room.m_bToltalWinDaoShu[uN] += AllKillCount
					room.m_bThreeKillResult[uN] = AllKillCount
				}
			})
		}
	})
}

//正常结束房间
func (room *sss_data_mgr) NormalEnd(cbReason int) {
	log.Debug("关闭房间")

}

func (room *sss_data_mgr) AfterEnd(a bool) {
	log.Debug("SSS AfterEnd")
}

//解散结束
func (room *sss_data_mgr) DismissEnd(cbReason int) {

}

func (room *sss_data_mgr) BeforeStartGame(UserCnt int) {
	room.GameStatus = GAME_FREE
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
	defaultCards := pk_base.GetCardByIdx(room.ConfigIdx)
	if room.wanFa == 1 {
		randNum := rand.Intn(13)
		room.laiZi = append(room.laiZi, defaultCards[randNum])
		room.laiZi = append(room.laiZi, defaultCards[randNum+13])
		room.laiZi = append(room.laiZi, defaultCards[randNum+13])
		room.laiZi = append(room.laiZi, defaultCards[randNum+13])
	}

	if room.jiaGongGong {
		len := len(defaultCards)
		tempCards := make([]int, len)
		copy(tempCards, defaultCards)
		for i := 0; i < 3; i++ {
			randNum := rand.Intn(len)
			room.publicCards = append(room.publicCards, tempCards[randNum])
			if randNum != len-1 {
				tempCards[randNum], tempCards[len-1] = tempCards[len-1], tempCards[randNum]
			}
			len--
		}
	}
	addCardNum := 0
	if room.jiaYiSe {
		addCardNum++
	}
	curPlayerCnt := room.PkBase.UserMgr.GetCurPlayerCnt()
	if curPlayerCnt > 4 {
		for i := 0; i < (curPlayerCnt - 4); i++ {
			addCardNum++
		}
	}
	if addCardNum > 0 {
		defaultCards = append(defaultCards, getColorCards(addCardNum)...)
	}
	if room.jiaDaXiaoWan {
		defaultCards = append(defaultCards, 0x4E, 0x4F)
	}

	room.LeftCardCount = len(defaultCards)
	gameLogic.RandCardList(room.bCardData, defaultCards)

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
		SendCard.Laizi = room.laiZi
		SendCard.PublicCards = room.publicCards
		u.WriteMsg(SendCard)
	})
}

func getColorCards(num int) (cards []int) {
	var colorCards = [][]int{
		[]int{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D},
		[]int{0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D},
		[]int{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D},
		[]int{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D},
	}
	len := len(colorCards)
	if num > len {
		num = len
	}
	for i := 0; i < num; i++ {
		randNum := rand.Intn(len)
		cards = append(cards, colorCards[randNum]...)
		if randNum != len-1 {
			colorCards[randNum], colorCards[len-1] = colorCards[len-1], colorCards[randNum]
		}
		len--
	}
	return
}
func (room *sss_data_mgr) AfterStartGame() {

}

//玩家摊牌
func (room *sss_data_mgr) ShowSSSCard(u *user.User, bDragon bool, bSpecialType bool, btSpecialData []int, bFrontCard []int, bMidCard []int, bBackCard []int) {
	userMgr := room.PkBase.UserMgr

	room.SpecialTypeTable[u] = bSpecialType
	room.Dragon[u] = bDragon

	room.m_bSegmentCard[u] = append(room.m_bSegmentCard[u], bFrontCard, bMidCard, bBackCard)

	room.m_bUserCardData[u] = make([]int, 0, 13)
	room.m_bUserCardData[u] = append(room.m_bUserCardData[u], bFrontCard...)
	room.m_bUserCardData[u] = append(room.m_bUserCardData[u], bMidCard...)
	room.m_bUserCardData[u] = append(room.m_bUserCardData[u], bBackCard...)

	btSpecialDataTemp := make([]int, 13)

	if bSpecialType {
		util.DeepCopy(&btSpecialDataTemp, &btSpecialData)
	}

	userMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(&pk_sss_msg.G2C_SSS_Open_Card{CurrentUser: u.ChairId})
	})

	room.OpenCardMap[u] = bFrontCard
	log.Debug("%d cccccc", len(room.OpenCardMap))
	if len(room.OpenCardMap) == room.PlayerCount { //已全摊
		// 游戏结束
		//userMgr.ForEachUser(func(u *user.User) {
		//room.PkBase.OnEventGameConclude(u.ChairId, u, GER_NORMAL)
		//})

		room.ComputeChOut()
		room.ComputeResult()

		gameEnd := &pk_sss_msg.G2C_SSS_COMPARE{}

		//LGameTax               int        //游戏税收
		gameEnd.LGameTax = 0
		//LGameEveryTax          []int      //每个玩家的税收
		gameEnd.LGameEveryTax = make([]int, room.PlayerCount)
		//LGameScore             []int      //游戏积分
		gameEnd.LGameScore = make([]int, room.PlayerCount)
		//BEndMode               int        //结束方式
		gameEnd.BEndMode = GER_NORMAL
		//CbCompareResult        [][]int    //每一道比较结果
		gameEnd.CbCompareResult = make([][]int, room.PlayerCount)
		//CbSpecialCompareResult []int      //特殊牌型比较结果
		gameEnd.CbSpecialCompareResult = make([]int, room.PlayerCount)
		//CbCompareDouble        []int      //翻倍的道数
		gameEnd.CbCompareDouble = make([]int, room.PlayerCount)
		//CbUserOverTime         []int      //玩家超时得到的道数
		gameEnd.CbUserOverTime = make([]int, room.PlayerCount)
		//CbCardData             [][]int    //扑克数据
		gameEnd.CbCardData = make([][]int, room.PlayerCount)
		//BUnderScoreDescribe    [][]int    //底分描述
		gameEnd.BUnderScoreDescribe = make([]string, room.PlayerCount)
		//BCompCardDescribe      [][][]int  //牌比描述
		gameEnd.BCompCardDescribe = make([][]string, room.PlayerCount)
		for i := 0; i < room.PlayerCount; i++ {
			gameEnd.BCompCardDescribe[i] = make([]string, 3)
		}
		//BToltalWinDaoShu       []int      //总共道数
		gameEnd.BToltalWinDaoShu = make([]int, room.PlayerCount)
		//LUnderScore            int        //底注分数
		gameEnd.LUnderScore = 0
		//BAllDisperse           []bool     //所有散牌
		gameEnd.BAllDisperse = make([]bool, room.PlayerCount)
		//BOverTime              []bool     //超时状态
		gameEnd.BOverTime = make([]bool, room.PlayerCount)
		//copy(gameEnd.BOverTime, room.m_bOverTime)
		//BUserLeft              []bool     //玩家逃跑
		gameEnd.BUserLeft = make([]bool, room.PlayerCount)
		//copy(gameEnd.BUserLeft, room.m_bUserLeft)
		//BLeft                  bool       //
		gameEnd.BLeft = false
		//LeftszName             [][]string //
		gameEnd.LeftszName = make([]string, room.PlayerCount)
		//copy(gameEnd.LeftszName,room.)
		//LeftChairID            []int      //
		gameEnd.LeftChairID = make([]int, room.PlayerCount)
		//BAllLeft               bool       //
		gameEnd.BAllLeft = false
		//LeftScore              []int      //
		gameEnd.LeftScore = make([]int, room.PlayerCount)
		//BSpecialCard           []bool     //是否为特殊牌
		gameEnd.BSpecialCard = make([]bool, room.PlayerCount)
		//BAllSpecialCard        bool       //全是特殊牌
		gameEnd.BAllSpecialCard = false
		//NTimer                 int        //结束后比牌、打枪时间
		gameEnd.NTimer = 0
		//ShootState             [][]int    //赢的玩家,输的玩家 2为赢的玩家，1为全输的玩家，0为没输没赢的玩家
		gameEnd.ShootState = make([][]int, room.PlayerCount)
		for i := range gameEnd.ShootState {
			gameEnd.ShootState[i] = make([]int, 2)

		}
		//M_nXShoot              int        //几家打枪
		gameEnd.M_nXShoot = room.m_nXShoot
		//CbThreeKillResult      []int      //全垒打加减分
		gameEnd.CbThreeKillResult = make([]int, room.PlayerCount)
		//BEnterExit             bool       //是否一进入就离开
		gameEnd.BEnterExit = false
		//WAllUser               int        //全垒打用户
		gameEnd.WAllUser = 0
		//copy(room.m_lGameScore,room.m_lLeftScore)

		nSpecialCard := 0
		nDragon := 0

		userMgr.ForEachUser(func(u *user.User) {
			if room.SpecialTypeTable[u] {
				nSpecialCard++
			}
			if room.Dragon[u] {
				nDragon++
			}
		})

		if room.PlayerCount == nSpecialCard+nDragon || room.PlayerCount <= nSpecialCard+1 {
			gameEnd.BAllSpecialCard = true
		} else {
			gameEnd.BAllSpecialCard = false
		}

		userMgr.ForEachUser(func(u *user.User) {
			gameEnd.CbCardData[u.ChairId] = make([]int, 13)
			copy(gameEnd.CbCardData[u.ChairId], room.m_bUserCardData[u])
			gameEnd.CbCompareResult[u.ChairId] = make([]int, 3)
			copy(gameEnd.CbCompareResult[u.ChairId], room.m_bCompareResult[u])
			gameEnd.CbCompareDouble[u.ChairId] = room.m_bCompareDouble[u]
			gameEnd.BToltalWinDaoShu[u.ChairId] = room.m_lGameScore[u]
			gameEnd.LGameScore[u.ChairId] = room.m_lGameScore[u]
			gameEnd.CbSpecialCompareResult[u.ChairId] = room.m_bSpecialCompareResult[u]
			gameEnd.BSpecialCard[u.ChairId] = room.SpecialTypeTable[u]
			for i := range room.m_bShootState {
				if room.m_bShootState[i][0] != nil {
					gameEnd.ShootState[i][0] = room.m_bShootState[i][0].ChairId

				}
				if room.m_bShootState[i][1] != nil {
					gameEnd.ShootState[i][1] = room.m_bShootState[i][1].ChairId

				}
			}
		})
		//dbg.Print(gameEnd)
		userMgr.ForEachUser(func(u *user.User) {
			u.WriteMsg(gameEnd)
		})
		room.AllResult[room.PkBase.TimerMgr.GetPlayCount()] = gameEnd.LGameScore
		room.PkBase.TimerMgr.AddPlayCount()
		//最后一局
		if room.PkBase.TimerMgr.GetPlayCount() >= room.PkBase.TimerMgr.GetMaxPlayCnt() {
			gameRecord := &pk_sss_msg.G2C_SSS_Record{}
			util.DeepCopy(&gameRecord, &room.AllResult)
			userMgr.ForEachUser(func(u *user.User) {
				u.WriteMsg(gameRecord)
			})
		}

	}

}

// 空闲状态场景
func (room *sss_data_mgr) SendStatusReady(u *user.User) {
	log.Debug("发送空闲状态场景消息")
	StatusFree := &pk_sss_msg.G2C_SSS_StatusFree{
		PlayerCount: room.PkBase.UserMgr.GetCurPlayerCnt(),
		SubCmd:      room.GameStatus,
	}

	room.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(StatusFree)
	})

}
