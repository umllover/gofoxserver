package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/pk_sss_msg"

	//dbg "github.com/funny/debug"

	"math/rand"

	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

//dbg "github.com/funny/debug"
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

func NewDataMgr(info *msg.L2G_CreatorRoom, uid int64, ConfigIdx int, name string, temp *base.GameServiceOption, base *SSS_Entry) *sss_data_mgr {
	d := new(sss_data_mgr)
	d.RoomData = pk_base.NewDataMgr(info.RoomID, uid, ConfigIdx, name, temp, base.Entry_base, info)
	//var setInfo sssOtherInfo
	//if err := json.Unmarshal([]byte(info.OtherInfo), &setInfo); err == nil {
	d.wanFa = int(info.OtherInfo["wanFa"].(float64))
	d.jiaYiSe = info.OtherInfo["jiaYiSe"].(bool)
	d.jiaGongGong = info.OtherInfo["jiaGongGong"].(bool)
	d.jiaDaXiaoWan = info.OtherInfo["jiaDaXiaoWan"].(bool)
	//}
	return d
}

type sssCardType struct {
	CT      int
	Item    *TagAnalyseItem
	isLaiZi bool
}

type sss_data_mgr struct {
	*pk_base.RoomData
	// 定时器
	ShowCardTimer *time.Timer

	//游戏变量
	wanFa        int
	jiaYiSe      bool
	jiaGongGong  bool
	jiaDaXiaoWan bool

	bCardData []int //牌的总数

	LeftCardCount int                 //剩下拍的数量
	OpenCardMap   map[*user.User]bool //摊牌数据

	// 游戏状态

	GameStatus int

	AllResult [][]int //每一局的结果

	gameEndStatus *pk_sss_msg.G2C_SSS_COMPARE
	gameRecord    *pk_sss_msg.G2C_SSS_Record

	Players               []int           //玩家
	PlayerCards           [][]int         //玩家手牌
	PlayerSpecialCardType []sssCardType   //玩家特殊牌型数据
	PlayerSegmentCards    [][][]int       //玩家组牌结果
	PlayerSegmentCardType [][]sssCardType //玩家组牌牌型数据
	Results               [][]int         //玩家每一道牌的水数
	SpecialResults        []int           //玩家特殊牌水数
	ToltalResults         []int           //玩家总共水数
	CompareResults        [][]int         //玩家每一道比较结果
	SpecialCompareResults []int           //玩家特殊牌型比较结果
	ShootState            [][]int         //打枪(0赢的玩家,1输的玩家)
	ShootResults          []int           //打枪的水数
	ShootNum              int             //几家打枪
	AddCards              []int           //加牌
	PublicCards           []int           //公共牌
	UniversalCards        []int           //万能牌

}

func (room *sss_data_mgr) InitRoom(UserCnt int) {
	//初始化
	log.Debug("初始化房间")
	room.PlayerCount = UserCnt
	room.Players = make([]int, UserCnt)
	room.cleanRoom(UserCnt)
}

func (room *sss_data_mgr) cleanRoom(UserCnt int) {
	log.Debug("清理房间")

	room.bCardData = make([]int, room.GetCfg().MaxRepertory) //牌堆
	room.OpenCardMap = make(map[*user.User]bool, UserCnt)

	room.LeftCardCount = room.GetCfg().MaxRepertory

	room.AllResult = make([][]int, room.PkBase.TimerMgr.GetMaxPlayCnt())

	room.PlayerCards = make([][]int, UserCnt)
	room.PlayerSpecialCardType = make([]sssCardType, UserCnt)
	room.PlayerSegmentCards = make([][][]int, UserCnt)
	room.PlayerSegmentCardType = make([][]sssCardType, UserCnt)
	room.Results = make([][]int, UserCnt)
	for i := range room.Results {
		room.Results[i] = make([]int, 3)
	}
	room.SpecialResults = make([]int, UserCnt)
	room.ToltalResults = make([]int, UserCnt)
	room.CompareResults = make([][]int, UserCnt)
	for i := range room.CompareResults {
		room.CompareResults[i] = make([]int, 3)
	}
	room.SpecialCompareResults = make([]int, UserCnt)
	room.ShootState = make([][]int, 6)
	room.ShootResults = make([]int, 6)

	room.AddCards = make([]int, 0)
	room.PublicCards = make([]int, 0, 3)
	room.UniversalCards = make([]int, 0, 3)

}

func (r *sss_data_mgr) checkLaiZi(carData []int) (bool, []int) {
	laiZiCount := 0
	tempData := make([]int, len(carData))
	copy(tempData, carData)
	if len(r.UniversalCards) > 0 {
		for i := range carData {
			for j := range r.UniversalCards {
				if carData[i] == r.UniversalCards[j] {
					tempData[i] = 0xFF
					laiZiCount++
				}
			}
		}
	}

	if laiZiCount == len(carData) {
		bossCount := 0
		for i := range carData {
			if carData[i] == 0x4E || carData[i] == 0x4F {
				tempData[i] = 0xFF
				bossCount++
			} else {
				tempData[i] = carData[i]
			}
		}
		return bossCount != 0, tempData
	}

	return laiZiCount > 0, tempData
}

func (r *sss_data_mgr) ComputeChOut() {
	lg := r.PkBase.LogicMgr.(*sss_logic)
	//lg.UniversalCards = r.UniversalCards
	userMgr := r.PkBase.UserMgr
	userMgr.ForEachUser(func(u *user.User) {
		i := u.ChairId
		var ct int
		var item *TagAnalyseItem

		r.PlayerSegmentCardType[i] = make([]sssCardType, 3)
		//特殊牌型
		isLaiZi, tempData := r.checkLaiZi(r.PlayerCards[i])
		ct, item = lg.SSSGetCardType(tempData)
		r.PlayerSpecialCardType[i].CT = ct
		r.PlayerSpecialCardType[i].Item = item
		r.PlayerSpecialCardType[i].isLaiZi = isLaiZi
		switch ct {
		case CT_THIRTEEN_FLUSH: //至尊清龙
			log.Debug("至尊清龙")
			r.SpecialResults[i] = 104
		case CT_THIRTEEN: //一条龙
			log.Debug("一条龙")
			r.SpecialResults[i] = 52
		case CT_THREE_STRAIGHTFLUSH: //三同花顺
			log.Debug("三同花顺")
			r.SpecialResults[i] = 36
		case CT_THREE_BOMB: //三分天下
			log.Debug("三分天下")
			r.SpecialResults[i] = 32
		case CT_FOUR_THREESAME: //四套三条
			log.Debug("四套三条")
			r.SpecialResults[i] = 16
		case CT_SIXPAIR: //六对半
			log.Debug("六对半")
			r.SpecialResults[i] = 6
			//有炸弹（四条）
			ct1, _ := lg.SSSGetCardType(r.PlayerSegmentCards[i][1])
			ct2, _ := lg.SSSGetCardType(r.PlayerSegmentCards[i][2])
			if ct1 == CT_FIVE_FOUR_ONE || ct2 == CT_FIVE_FOUR_ONE {
				r.SpecialResults[i] = 10
			}
		case CT_THREE_FLUSH: //三同花
			log.Debug("三同花")
			r.SpecialResults[i] = 6
			//有同花顺
			if lg.IsLine(r.PlayerSegmentCards[i][1], len(r.PlayerSegmentCards[i][1]), true) ||
				lg.IsLine(r.PlayerSegmentCards[i][2], len(r.PlayerSegmentCards[i][2]), true) {
				r.SpecialResults[i] = 10
			}
		case CT_THREE_STRAIGHT: //三顺子
			log.Debug("三顺子")
			r.SpecialResults[i] = 6
			//有同花顺
			if lg.IsLine(r.PlayerSegmentCards[i][1], len(r.PlayerSegmentCards[i][1]), true) ||
				lg.IsLine(r.PlayerSegmentCards[i][2], len(r.PlayerSegmentCards[i][2]), true) {
				r.SpecialResults[i] = 10
			}
		default: //普通牌型
			//前敦
			isLaiZi, tempData := r.checkLaiZi(r.PlayerSegmentCards[i][0])
			ct, item = lg.SSSGetCardType(tempData)
			r.PlayerSegmentCardType[i][0].CT = ct
			r.PlayerSegmentCardType[i][0].Item = item
			r.PlayerSegmentCardType[i][0].isLaiZi = isLaiZi
			switch ct {
			case CT_SINGLE: //散牌
				log.Debug("前敦散牌")
				r.Results[i][0] = 1
			case CT_ONE_DOUBLE: //对子
				log.Debug("前敦对子")
				r.Results[i][0] = 1
			case CT_THREE: //三条
				log.Debug("前敦三条")
				r.Results[i][0] = 3
			}
			//中墩
			isLaiZi, tempData = r.checkLaiZi(r.PlayerSegmentCards[i][1])
			ct, item = lg.SSSGetCardType(tempData)
			r.PlayerSegmentCardType[i][1].CT = ct
			r.PlayerSegmentCardType[i][1].Item = item
			r.PlayerSegmentCardType[i][1].isLaiZi = isLaiZi
			switch ct {
			case CT_SINGLE: //散牌
				log.Debug("中墩散牌")
				r.Results[i][1] = 1
			case CT_ONE_DOUBLE: //对子
				log.Debug("中墩对子")
				r.Results[i][1] = 1
			case CT_FIVE_TWO_DOUBLE: //两对
				log.Debug("中墩两对")
				r.Results[i][1] = 1
			case CT_THREE: //三条
				log.Debug("中墩三条")
				r.Results[i][1] = 3
			case CT_FIVE_STRAIGHT: //顺子
				log.Debug("中墩顺子")
				r.Results[i][1] = 1
			case CT_FIVE_FLUSH: //同花
				log.Debug("中墩同花")
				r.Results[i][1] = 1
			case CT_FIVE_THREE_DEOUBLE: //葫芦
				log.Debug("中墩葫芦")
				r.Results[i][1] = 2
			case CT_FIVE_FOUR_ONE: //铁支
				log.Debug("中墩铁支")
				r.Results[i][1] = 8
			case CT_FIVE_STRAIGHT_FLUSH:
				log.Debug("中墩同花顺")
				r.Results[i][1] = 10
			case CT_FIVE_SAME:
				log.Debug("中墩五同")
				r.Results[i][1] = 14
			}

			//尾墩
			isLaiZi, tempData = r.checkLaiZi(r.PlayerSegmentCards[i][2])
			ct, item = lg.SSSGetCardType(tempData)
			r.PlayerSegmentCardType[i][2].CT = ct
			r.PlayerSegmentCardType[i][2].Item = item
			r.PlayerSegmentCardType[i][2].isLaiZi = isLaiZi
			switch ct {
			case CT_SINGLE: //散牌
				log.Debug("后墩散牌")
				r.Results[i][2] = 1
			case CT_ONE_DOUBLE: //对子
				log.Debug("后墩对子")
				r.Results[i][2] = 1
			case CT_FIVE_TWO_DOUBLE: //两对
				log.Debug("后墩两对")
				r.Results[i][2] = 1
			case CT_THREE: //三条
				log.Debug("后墩三条")
				r.Results[i][2] = 3
			case CT_FIVE_STRAIGHT: //顺子
				log.Debug("后墩顺子")
				r.Results[i][1] = 1
			case CT_FIVE_FLUSH: //同花
				log.Debug("后墩同花")
				r.Results[i][2] = 1
			case CT_FIVE_THREE_DEOUBLE: //葫芦
				log.Debug("后墩葫芦")
				r.Results[i][2] = 1
			case CT_FIVE_FOUR_ONE: //铁支
				log.Debug("后墩铁支")
				r.Results[i][2] = 4
			case CT_FIVE_STRAIGHT_FLUSH:
				log.Debug("后墩同花顺")
				r.Results[i][2] = 5
			case CT_FIVE_SAME:
				log.Debug("后墩五同")
				r.Results[i][2] = 7
			}
		}
	})
}

func (r *sss_data_mgr) ComputeResult() {
	lg := r.PkBase.LogicMgr.(*sss_logic)
	//lg.UniversalCards = r.UniversalCards
	userMgr := r.PkBase.UserMgr
	//打枪次数
	shootPlayerNum := make([]int, r.PlayerCount)
	userMgr.ForEachUser(func(u *user.User) {
		i := u.ChairId
		winPoint := 0
		userMgr.ForEachUser(func(u *user.User) {
			j := u.ChairId
			if i != j {

				//都是普通牌型
				if r.SpecialResults[i] == 0 && r.SpecialResults[j] == 0 {
					firstResult := lg.SSSCompareCard(r.PlayerSegmentCardType[j][0], r.PlayerSegmentCardType[i][0])
					switch firstResult {
					case 1:
						winPoint += r.Results[i][0]
						r.CompareResults[i][0] += r.Results[i][0]
					case -1:
						winPoint -= r.Results[j][0]
						r.CompareResults[i][0] -= r.Results[j][0]
					}
					midResult := lg.SSSCompareCard(r.PlayerSegmentCardType[j][1], r.PlayerSegmentCardType[i][1])
					switch midResult {
					case 1:
						winPoint += r.Results[i][1]
						r.CompareResults[i][1] += r.Results[i][1]
					case -1:
						winPoint -= r.Results[j][1]
						r.CompareResults[i][1] -= r.Results[j][1]
					}
					backResult := lg.SSSCompareCard(r.PlayerSegmentCardType[j][2], r.PlayerSegmentCardType[i][2])
					switch backResult {
					case 1:
						winPoint += r.Results[i][2]
						r.CompareResults[i][2] += r.Results[i][2]
					case -1:
						winPoint -= r.Results[j][2]
						r.CompareResults[i][2] -= r.Results[j][2]
					}
					// 打枪
					if firstResult >= 0 && midResult >= 0 && backResult >= 0 && firstResult+midResult+backResult > 0 {
						r.ToltalResults[j] -= winPoint
						winPoint *= 2
						r.ShootResults[j] = winPoint
						shootPlayerNum[i]++
					}
				}
				//特殊对普通
				if r.SpecialResults[i] > 0 && r.SpecialResults[j] == 0 {
					winPoint += r.SpecialResults[i]
					r.SpecialCompareResults[i] += r.SpecialResults[i]
				}
				//普通对特殊
				if r.SpecialResults[i] == 0 && r.SpecialResults[j] > 0 {
					winPoint -= r.SpecialResults[j]
					r.SpecialCompareResults[i] -= r.SpecialResults[j]
				}
				//都是特殊牌型
				if r.SpecialResults[i] > 0 && r.SpecialResults[j] > 0 {
					switch lg.SSSCompareCard(r.PlayerSpecialCardType[j], r.PlayerSpecialCardType[i]) {
					case 1:
						winPoint += r.SpecialResults[i]
						r.SpecialCompareResults[i] += r.SpecialResults[i]
					case -1:
						winPoint -= r.SpecialResults[j]
						r.SpecialCompareResults[i] -= r.SpecialResults[j]
					}
				}
			}
		})
		r.ToltalResults[i] += winPoint
	})

	//全垒打加分
	userMgr.ForEachUser(func(u *user.User) {
		i := u.ChairId
		if (r.PlayerCount >= 4) && (shootPlayerNum[i] == r.PlayerCount-1) {
			r.ToltalResults[i] *= 2
			for j, v := range r.ShootResults {
				if j == i {
					continue
				}
				r.ToltalResults[j] -= v * 2
			}
			return
		}
	})
}

//正常结束房间
func (room *sss_data_mgr) NormalEnd(a int) {
	log.Debug("关闭房间")

}

//解散结束
func (room *sss_data_mgr) DismissEnd(a int) {

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
	//清理上一局数据
	room.cleanRoom(room.PlayerCount)

	userMgr := room.PkBase.UserMgr
	gameLogic := room.PkBase.LogicMgr
	defaultCards := pk_base.GetCardByIdx(room.ConfigIdx)
	if room.wanFa == 1 {
		randNum := rand.Intn(13)
		room.UniversalCards = append(room.UniversalCards, defaultCards[randNum])
		room.UniversalCards = append(room.UniversalCards, defaultCards[randNum+13])
		room.UniversalCards = append(room.UniversalCards, defaultCards[randNum+13])
		room.UniversalCards = append(room.UniversalCards, defaultCards[randNum+13])
	}

	if room.jiaGongGong {
		len := len(defaultCards)
		tempCards := make([]int, len)
		copy(tempCards, defaultCards)
		for i := 0; i < 3; i++ {
			randNum := rand.Intn(len)
			room.PublicCards = append(room.PublicCards, tempCards[randNum])
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
	room.bCardData = make([]int, room.LeftCardCount)
	gameLogic.RandCardList(room.bCardData, defaultCards)

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, US_PLAYING)
	})

	userMgr.ForEachUser(func(u *user.User) {
		room.PlayerCards[u.ChairId] = make([]int, 0, 13)
		for i := 0; i < 13; i++ {
			room.PlayerCards[u.ChairId] = append(room.PlayerCards[u.ChairId], room.GetOneCard())
		}
	})

	userMgr.ForEachUser(func(u *user.User) {
		SendCard := &pk_sss_msg.G2C_SSS_SendCard{}
		SendCard.CardData = room.PlayerCards[u.ChairId]
		SendCard.CellScore = room.CellScore
		SendCard.Laizi = room.UniversalCards
		SendCard.PublicCards = room.PublicCards
		u.WriteMsg(SendCard)
	})

	// 启动定时器
	room.startShowCardTimer(40)
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
func (room *sss_data_mgr) ShowSSSCard(u *user.User, bDragon bool, bSpecialType bool, btSpecialData []int, FrontCard []int, MidCard []int, BackCard []int) {
	userMgr := room.PkBase.UserMgr

	// //解除托管
	// trustees := userMgr.GetTrustees()
	// for i := range trustees {
	// 	if trustees[i] == true {
	// 		if u == userMgr.GetUserByChairId(i) {
	// 			room.Trustee(u, false)
	// 			break
	// 		}
	// 	}
	// }

	room.PlayerSegmentCards[u.ChairId] = append(room.PlayerSegmentCards[u.ChairId], FrontCard, MidCard, BackCard)
	room.PlayerCards[u.ChairId] = make([]int, 0, 13)
	room.PlayerCards[u.ChairId] = append(room.PlayerCards[u.ChairId], FrontCard...)
	room.PlayerCards[u.ChairId] = append(room.PlayerCards[u.ChairId], MidCard...)
	room.PlayerCards[u.ChairId] = append(room.PlayerCards[u.ChairId], BackCard...)

	userMgr.ForEachUser(func(user *user.User) {
		user.WriteMsg(&pk_sss_msg.G2C_SSS_Open_Card{CurrentUser: u.ChairId})
	})

	room.OpenCardMap[u] = true
	if len(room.OpenCardMap) == room.PlayerCount { //已全摊
		room.stopShowCardTimer()
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
		gameEnd.M_nXShoot = 0
		//CbThreeKillResult      []int      //全垒打加减分
		gameEnd.CbThreeKillResult = make([]int, room.PlayerCount)
		//BEnterExit             bool       //是否一进入就离开
		gameEnd.BEnterExit = false
		//WAllUser               int        //全垒打用户
		gameEnd.WAllUser = 0
		//copy(room.m_lGameScore,room.m_lLeftScore)

		gameEnd.BAllSpecialCard = false

		userMgr.ForEachUser(func(u *user.User) {
			gameEnd.CbCardData[u.ChairId] = make([]int, 13)
			copy(gameEnd.CbCardData[u.ChairId], room.PlayerCards[u.ChairId])
			gameEnd.CbCompareResult[u.ChairId] = make([]int, 3)
			copy(gameEnd.CbCompareResult[u.ChairId], room.CompareResults[u.ChairId])
			gameEnd.CbCompareDouble[u.ChairId] = 0
			gameEnd.BToltalWinDaoShu[u.ChairId] = room.ToltalResults[u.ChairId]
			gameEnd.LGameScore[u.ChairId] = room.ToltalResults[u.ChairId]
			gameEnd.CbSpecialCompareResult[u.ChairId] = room.SpecialCompareResults[u.ChairId]
			gameEnd.BSpecialCard[u.ChairId] = false

		})

		userMgr.ForEachUser(func(u *user.User) {
			u.WriteMsg(gameEnd)
		})
		room.gameEndStatus = gameEnd

		room.AllResult[room.PkBase.TimerMgr.GetPlayCount()-1] = gameEnd.LGameScore
		//room.PkBase.TimerMgr.AddPlayCount()
		//最后一局
		if room.PkBase.TimerMgr.GetPlayCount() >= room.PkBase.TimerMgr.GetMaxPlayCnt() {
			gameRecord := &pk_sss_msg.G2C_SSS_Record{}
			util.DeepCopy(&gameRecord.AllResult, &room.AllResult)
			allScore := make([]int, room.PlayerCount)

			for i := 0; i < room.PkBase.TimerMgr.GetPlayCount(); i++ {
				for j := range allScore {
					allScore[j] += room.AllResult[i][j]
				}
			}
			gameRecord.AllScore = allScore

			gameRecord.Reason = GER_NORMAL
			userMgr.ForEachUser(func(u *user.User) {
				u.WriteMsg(gameRecord)
			})
			room.gameRecord = gameRecord
		}

	}

}

// 空闲状态场景
func (room *sss_data_mgr) SendStatusReady(u *user.User) {
	log.Debug("发送空闲状态场景消息")
	StatusFree := &pk_sss_msg.G2C_SSS_StatusFree{
		PlayerCount:      room.PkBase.UserMgr.GetCurPlayerCnt(),
		SubCmd:           room.GameStatus,
		CurrentPlayCount: room.PkBase.TimerMgr.GetPlayCount(),
		MaxPlayCount:     room.PkBase.TimerMgr.GetMaxPlayCnt(),
	}

	room.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(StatusFree)
	})

}

func (room *sss_data_mgr) SendStatusPlay(u *user.User) {
	statusPlay := &pk_sss_msg.G2C_SSS_StatusPlay{}
	//WCurrentUser       int             `json:"wCurrentUser"`       //当前玩家
	statusPlay.WCurrentUser = u.ChairId
	//LCellScore         int             `json:"lCellScore"`         //单元底分
	statusPlay.LCellScore = 0
	//NChip              []int           `json:"nChip"`              //下注大小
	statusPlay.NChip = make([]int, 0)
	//BHandCardData      []int           `json:"bHandCardData"`      //扑克数据

	// for ucd := range room.m_bUserCardData {
	// 	if ucd.ChairId == u.ChairId {
	// 		statusPlay.BHandCardData = room.m_bUserCardData[ucd]
	// 	}
	// }

	statusPlay.BHandCardData = room.PlayerCards[u.ChairId]

	//BSegmentCard       [][]int       `json:"bSegmentCard"`         //分段扑克
	statusPlay.BSegmentCard = room.PlayerSegmentCards[u.ChairId]
	//BFinishSegment     []bool          `json:"bFinishSegment"`     //完成分段
	statusPlay.BFinishSegment = make([]bool, room.PlayerCount)
	for user := range room.OpenCardMap {
		statusPlay.BFinishSegment[user.ChairId] = room.OpenCardMap[user]
	}
	//WUserToltalChip    int             `json:"wUserToltalChip"`    //总共金币
	statusPlay.WUserToltalChip = 0
	//BOverTime          []bool          `json:"bOverTime"`          //超时状态
	statusPlay.BOverTime = make([]bool, 0)
	//BSpecialTypeTable1 []bool          `json:"bSpecialTypeTable1"` //是否特殊牌型
	statusPlay.BSpecialTypeTable1 = make([]bool, 0)
	//BDragon1           []bool          `json:"bDragon1"`           //是否倒水
	statusPlay.BDragon1 = make([]bool, 0)
	//BAllHandCardData   [][]int         `json:"bAllHandCardData"`   //所有玩家的扑克数据
	statusPlay.BAllHandCardData = make([][]int, 0)
	//SGameEnd           G2C_SSS_COMPARE `json:"sGameEnd"`           //游戏结束数据
	statusPlay.SGameEnd = *room.gameEndStatus
	statusPlay.Record = *room.gameRecord
	statusPlay.PlayerCount = room.PkBase.UserMgr.GetCurPlayerCnt()
	statusPlay.CurrentPlayCount = room.PkBase.TimerMgr.GetPlayCount()
	statusPlay.MaxPlayCount = room.PkBase.TimerMgr.GetMaxPlayCnt()
	statusPlay.Laizi = room.UniversalCards
	statusPlay.PublicCards = room.PublicCards

	u.WriteMsg(statusPlay)
}

// 托管
func (room *sss_data_mgr) Trustee(u *user.User, t bool) {
	room.PkBase.UserMgr.SetUsetTrustee(u.ChairId, t)
	DataTrustee := &pk_sss_msg.G2C_SSS_TRUSTEE{}
	DataTrustee.TrusteeUser = u.ChairId
	DataTrustee.Trustee = t

	room.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		log.Debug("托管状态%v", DataTrustee)
		u.WriteMsg(DataTrustee)
	})
}

// 托管操作
func (room *sss_data_mgr) trusteeOperate() {
	trustees := room.PkBase.UserMgr.GetTrustees()
	for i := range trustees {
		u := room.PkBase.UserMgr.GetUserByChairId(i)
		if u != nil {
			if trustees[i] == true {
				segmentCard1, segmentCard2, segmentCard3 := room.getSegmentCard(i)
				room.ShowSSSCard(u, false, false, []int{}, segmentCard1, segmentCard2, segmentCard3)
			} else {
				if !room.OpenCardMap[u] {
					room.Trustee(u, true)
					segmentCard1, segmentCard2, segmentCard3 := room.getSegmentCard(i)
					room.ShowSSSCard(u, false, false, []int{}, segmentCard1, segmentCard2, segmentCard3)
				}
			}
		}
	}
}

func (room *sss_data_mgr) getSegmentCard(chairId int) (segmentCard1, segmentCard2, segmentCard3 []int) {

	cardData := room.PlayerCards[chairId]
	newCardData := []int{}
	//后墩
	segmentCard3, newCardData = room.get5card(cardData)
	//中墩
	segmentCard2, newCardData = room.get5card(newCardData)
	//前墩
	segmentCard1 = newCardData

	return
}

func (room *sss_data_mgr) get5card(cardData []int) (segmentCard []int, newCardData []int) {
	lg := room.PkBase.LogicMgr.(*sss_logic)

	segmentCard = make([]int, 0, 5)
	index := 0

	TagAnalyseItemArray := lg.AnalyseCard(cardData)
	//五同
	if len(segmentCard) == 0 && TagAnalyseItemArray.bFiveCount > 0 {
		index = TagAnalyseItemArray.bFiveFirst[0]
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[index:index+5]...)
	}
	//同花顺
	if len(segmentCard) == 0 {
		for i := 3; i >= 0; i-- {
			uniqueColorCard := lg.GetUniqueColorCard(TagAnalyseItemArray.cardData, i)
			l := len(uniqueColorCard)
			if l < 5 {
				continue
			}
			for j := 0; j <= l-5; j++ {
				if lg.IsLine(uniqueColorCard[j:j+5], 5, true) {
					segmentCard = append(segmentCard, uniqueColorCard[j:j+5]...)
					break
				}
			}
			if len(segmentCard) == 0 {
				logicCards := make([]int, 6)
				for _, v := range uniqueColorCard {
					for k := range logicCards {
						if lg.GetCardValue(v) == k && logicCards[k] == 0 {
							logicCards[k] = v
						}
					}
				}
				isSmallStraight := true
				for _, v := range logicCards[1:6] {
					if v == 0 {
						isSmallStraight = false
					}
				}
				if isSmallStraight {
					segmentCard = append(segmentCard, logicCards[1:6]...)
					break
				}
			}
		}
	}

	//铁支
	if len(segmentCard) == 0 && TagAnalyseItemArray.bFourCount > 0 && TagAnalyseItemArray.bOneCount > 0 {
		index = TagAnalyseItemArray.bFourFirst[0]
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[index:index+4]...)
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[TagAnalyseItemArray.bOneFirst[0]])
	}
	//葫芦
	if len(segmentCard) == 0 && TagAnalyseItemArray.bThreeCount > 0 && TagAnalyseItemArray.bTwoCount > 0 {
		index = TagAnalyseItemArray.bThreeFirst[0]
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[index:index+3]...)
		index = TagAnalyseItemArray.bTwoFirst[0]
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[index:index+2]...)
	}
	//同花
	if len(segmentCard) == 0 {
		for i := 3; i >= 0; i-- {
			colorCard := lg.GetColorCard(TagAnalyseItemArray.cardData, i)
			l := len(colorCard)
			if l < 5 {
				continue
			}
			segmentCard = append(segmentCard, colorCard[:5]...)
			break
		}
	}
	//顺子
	if len(segmentCard) == 0 {
		l := len(TagAnalyseItemArray.cardData)
		for j := 0; j <= l-5; j++ {
			if lg.IsLine(TagAnalyseItemArray.cardData[j:j+5], 5, false) {
				segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[j:j+5]...)
				break
			}
		}
		if len(segmentCard) == 0 {
			logicCards := make([]int, 6)
			for _, v := range TagAnalyseItemArray.cardData {
				for k := range logicCards {
					if lg.GetCardValue(v) == k && logicCards[k] == 0 {
						logicCards[k] = v
					}
				}
			}
			isSmallStraight := true
			for _, v := range logicCards[1:6] {
				if v == 0 {
					isSmallStraight = false
				}
			}
			if isSmallStraight {
				segmentCard = append(segmentCard, logicCards[1:6]...)
			}
		}
	}
	//三条
	if len(segmentCard) == 0 && TagAnalyseItemArray.bThreeCount > 0 && TagAnalyseItemArray.bOneCount >= 2 {
		index = TagAnalyseItemArray.bThreeFirst[0]
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[index:index+3]...)
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[TagAnalyseItemArray.bOneFirst[0]])
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[TagAnalyseItemArray.bOneFirst[1]])
	}
	//两对
	if len(segmentCard) == 0 && TagAnalyseItemArray.bTwoCount >= 2 && TagAnalyseItemArray.bOneCount >= 1 {
		index = TagAnalyseItemArray.bTwoFirst[0]
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[index:index+2]...)
		index = TagAnalyseItemArray.bTwoFirst[1]
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[index:index+2]...)
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[TagAnalyseItemArray.bOneFirst[0]])
	}
	//对子
	if len(segmentCard) == 0 && TagAnalyseItemArray.bTwoCount >= 1 && TagAnalyseItemArray.bOneCount >= 3 {
		index = TagAnalyseItemArray.bTwoFirst[0]
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[index:index+2]...)
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[TagAnalyseItemArray.bOneFirst[0]])
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[TagAnalyseItemArray.bOneFirst[1]])
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[TagAnalyseItemArray.bOneFirst[2]])
	}
	//散牌
	if len(segmentCard) == 0 {
		segmentCard = append(segmentCard, TagAnalyseItemArray.cardData[:5]...)
	}

	newCardData = lg.getUnUsedCard(TagAnalyseItemArray.cardData, segmentCard)

	return

}

func (r *sss_data_mgr) startShowCardTimer(nTime int) {
	if r.ShowCardTimer != nil {
		r.ShowCardTimer.Stop()
		r.ShowCardTimer = nil
	}

	f := func() {
		r.trusteeOperate()
	}

	r.ShowCardTimer = time.AfterFunc(time.Duration(nTime+5)*time.Second, f)
}

func (r *sss_data_mgr) resetShowCardTimer(nTime int) {
	log.Debug("重置定时器时间%d", nTime)
	if r.ShowCardTimer != nil {
		r.ShowCardTimer.Reset(time.Duration(nTime+5) * time.Second)
	}
}

func (r *sss_data_mgr) stopShowCardTimer() {
	if r.ShowCardTimer != nil {
		log.Debug("停止定时器")
		r.ShowCardTimer.Stop()
		r.ShowCardTimer = nil
	}
}
