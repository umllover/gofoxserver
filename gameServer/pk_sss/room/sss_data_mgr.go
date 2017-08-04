package room

import (
	"encoding/json"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	. "mj/common/cost"
	"mj/common/msg/pk_sss_msg"

	"mj/gameServer/db/model"

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
	//laiZi        []int
	//publicCards  []int

	bCardData []int //牌的总数
	// m_bUserCardData         map[*user.User][]int   //玩家扑克
	// m_bSegmentCard          map[*user.User][][]int //分段扑克
	// m_bFinishSegment        []int                  //完成分段
	// m_bShowCardCount        int                    //摊牌数目
	// m_bCompleteCompareCount int                    //完成比较
	// m_bOverTime             []int                  //摊牌超时
	// m_bUserLeft             []int                  //玩家强退

	// SpecialTypeTable map[*user.User]bool  //是否特殊牌型
	// Dragon           map[*user.User]bool  //是否倒水
	// m_nPlayerCount   int                  //实际玩家人数
	// CbResult         map[*user.User][]int //每一道的道数
	// cbSpecialResult  map[*user.User]int   //特殊牌型的道数

	LeftCardCount int                 //剩下拍的数量
	OpenCardMap   map[*user.User]bool //摊牌数据
	//比较结果

	// m_bCompareResult        map[*user.User][]int //每一道比较结果
	// m_bShootState           [][]*user.User       //打枪(0赢的玩家,1输的玩家)
	// m_bThreeKillResult      map[*user.User]int   //全垒打加减分
	// m_bToltalWinDaoShu      map[*user.User]int   //总共道数
	// m_bCompareDouble        map[*user.User]int   //打枪的道数
	// m_bSpecialCompareResult map[*user.User]int   //特殊牌型比较结果
	// m_lGameScore            map[*user.User]int   //游戏积分
	// m_nXShoot               int                  //几家打枪
	// m_lCellScore            int                  //单元底分

	// 游戏状态

	GameStatus int
	// BtCardSpecialData []int
	AllResult [][]int //每一局的结果

	gameEndStatus *pk_sss_msg.G2C_SSS_COMPARE

	/////////////////////////

	//PlayerNum             int       //玩家数量
	Players               []int     //玩家
	PlayerCards           [][]int   //玩家手牌
	PlayerSegmentCards    [][][]int //玩家组牌结果
	Results               [][]int   //玩家每一道牌的水数
	SpecialResults        []int     //玩家特殊牌水数
	ToltalResults         []int     //玩家总共水数
	CompareResults        [][]int   //玩家每一道比较结果
	SpecialCompareResults []int     //玩家特殊牌型比较结果
	ShootState            [][]int   //打枪(0赢的玩家,1输的玩家)
	ShootResults          []int     //打枪的水数
	ShootNum              int       //几家打枪
	AddCards              []int     //加牌
	PublicCards           []int     //公共牌
	UniversalCards        []int     //万能牌

}

func (room *sss_data_mgr) InitRoom(UserCnt int) {
	//初始化
	log.Debug("初始化房间")

	//room.cbSpecialResult = make(map[*user.User]int, UserCnt)
	//room.CbResult = make(map[*user.User][]int, UserCnt)
	room.PlayerCount = UserCnt
	//room.m_bSegmentCard = make(map[*user.User][][]int, UserCnt)
	//room.bCardData = make([]int, room.GetCfg().MaxRepertory) //牌堆
	room.OpenCardMap = make(map[*user.User]bool, UserCnt)
	//room.Dragon = make(map[*user.User]bool, UserCnt)
	//room.SpecialTypeTable = make(map[*user.User]bool, UserCnt)
	//room.m_bUserCardData = make(map[*user.User][]int, UserCnt)
	//room.m_bCompareDouble = make(map[*user.User]int, UserCnt)
	//room.m_bCompareResult = make(map[*user.User][]int, UserCnt)
	//room.m_bShootState = make([][]*user.User, UserCnt)
	//room.m_bSpecialCompareResult = make(map[*user.User]int, UserCnt)
	//room.m_bThreeKillResult = make(map[*user.User]int, UserCnt)
	//room.m_bToltalWinDaoShu = make(map[*user.User]int, UserCnt)
	//room.m_lGameScore = make(map[*user.User]int, UserCnt)
	//room.m_nXShoot = 0
	//room.BtCardSpecialData = make([]int, 13)
	room.LeftCardCount = room.GetCfg().MaxRepertory

	//room.laiZi = make([]int, 0, 6)

	room.AllResult = make([][]int, room.PkBase.TimerMgr.GetMaxPayCnt())

	room.gameEndStatus = &pk_sss_msg.G2C_SSS_COMPARE{}

	/////////////////////////////
	//room.PlayerNum = UserCnt
	room.Players = make([]int, UserCnt)
	room.PlayerCards = make([][]int, UserCnt)
	room.PlayerSegmentCards = make([][][]int, UserCnt)
	room.Results = make([][]int, UserCnt)
	for i := range room.Results {
		room.Results[i] = make([]int, 3)
	}
	room.SpecialResults = make([]int, UserCnt)
	room.ToltalResults = make([]int, UserCnt)
	room.CompareResults = make([][]int, UserCnt)
	room.SpecialCompareResults = make([]int, UserCnt)
	room.ShootState = make([][]int, 6)
	room.ShootResults = make([]int, 6)
	room.AddCards = make([]int, 0)
	room.PublicCards = make([]int, 0, 3)
	room.UniversalCards = make([]int, 0, 3)

}

func (r *sss_data_mgr) ComputeChOut() {
	lg := r.PkBase.LogicMgr
	for i := 0; i < r.PlayerCount; i++ {
		//特殊牌型
		switch lg.GetCardType(r.PlayerCards[i]) {
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
			if lg.GetCardType(r.PlayerSegmentCards[i][1]) == CT_FIVE_FOUR_ONE ||
				lg.GetCardType(r.PlayerSegmentCards[i][2]) == CT_FIVE_FOUR_ONE {
				r.SpecialResults[i] = 10
			}
		case CT_THREE_STRAIGHT: //三顺子
			log.Debug("三顺子")
			r.SpecialResults[i] = 6
			//todo 有同花顺
		case CT_THREE_FLUSH: //三同花
			log.Debug("三同花")
			r.SpecialResults[i] = 6
			//todo 有同花顺
		default: //普通牌型
			//前敦
			switch lg.GetCardType(r.PlayerSegmentCards[i][0]) {
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
			switch lg.GetCardType(r.PlayerSegmentCards[i][1]) {
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
			case CT_FIVE_MIXED_FLUSH_NO_A, CT_FIVE_MIXED_FLUSH_FIRST_A, CT_FIVE_MIXED_FLUSH_BACK_A: //顺子
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
			case CT_FIVE_STRAIGHT_FLUSH_NO_A, CT_FIVE_STRAIGHT_FLUSH_FIRST_A, CT_FIVE_STRAIGHT_FLUSH_BACK_A:
				log.Debug("中墩同花顺")
				r.Results[i][1] = 10
			}
			//尾墩
			switch lg.GetCardType(r.PlayerSegmentCards[i][2]) {
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
			case CT_FIVE_MIXED_FLUSH_NO_A, CT_FIVE_MIXED_FLUSH_FIRST_A, CT_FIVE_MIXED_FLUSH_BACK_A: //顺子
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
			case CT_FIVE_STRAIGHT_FLUSH_NO_A, CT_FIVE_STRAIGHT_FLUSH_FIRST_A, CT_FIVE_STRAIGHT_FLUSH_BACK_A:
				log.Debug("后墩同花顺")
				r.Results[i][2] = 5
			}
		}
	}
}

func (r *sss_data_mgr) ComputeResult() {
	lg := r.PkBase.LogicMgr.(*sss_logic)
	//打枪次数
	shootPlayerNum := make([]int, r.PlayerCount)
	for i := 0; i < r.PlayerCount; i++ {
		winPoint := 0
		for j := 0; j < r.PlayerCount; j++ {
			if i == j {
				continue
			}
			//都是普通牌型
			if r.SpecialResults[i] == 0 && r.SpecialResults[j] == 0 {
				firstResult := lg.SSSCompareCard(r.PlayerSegmentCards[j][0], r.PlayerSegmentCards[i][0])
				switch firstResult {
				case 1:
					winPoint += r.Results[i][0]
					r.CompareResults[i][0] += r.Results[i][0]
				case -1:
					winPoint -= r.Results[j][0]
					r.CompareResults[i][0] -= r.Results[j][0]
				}
				midResult := lg.SSSCompareCard(r.PlayerSegmentCards[j][1], r.PlayerSegmentCards[i][1])
				switch midResult {
				case 1:
					winPoint += r.Results[i][1]
					r.CompareResults[i][1] += r.Results[i][1]
				case -1:
					winPoint -= r.Results[j][1]
					r.CompareResults[i][1] -= r.Results[j][1]
				}
				backResult := lg.SSSCompareCard(r.PlayerSegmentCards[j][2], r.PlayerSegmentCards[i][2])
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
				switch lg.SSSCompareCard(r.PlayerCards[j], r.PlayerCards[i]) {
				case 1:
					winPoint += r.SpecialResults[i]
					r.SpecialCompareResults[i] += r.SpecialResults[i]
				case -1:
					winPoint -= r.SpecialResults[j]
					r.SpecialCompareResults[i] -= r.SpecialResults[j]
				}
			}

		}
		r.ToltalResults[i] += winPoint
	}

	//全垒打加分
	for i := 0; i < r.PlayerCount; i++ {
		if (r.PlayerCount >= 4) && (shootPlayerNum[i] == r.PlayerCount-1) {
			r.ToltalResults[i] *= 2
			for j, v := range r.ShootResults {
				if j == i {
					continue
				}
				r.ToltalResults[j] -= v * 2
			}
			break
		}
	}
}

//正常结束房间
func (room *sss_data_mgr) NormalEnd() {
	log.Debug("关闭房间")

}

func (room *sss_data_mgr) AfterEnd(a bool) {
	log.Debug("SSS AfterEnd")
}

//解散结束
func (room *sss_data_mgr) DismissEnd() {

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
	gameLogic.RandCardList(room.bCardData, defaultCards)

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, US_PLAYING)
	})

	userMgr.ForEachUser(func(u *user.User) {
		for i := 0; i < pk_base.GetCfg(pk_base.IDX_SSS).MaxCount; i++ {
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

	//room.SpecialTypeTable[u] = bSpecialType
	//room.Dragon[u] = bDragon

	// room.m_bSegmentCard[u] = append(room.m_bSegmentCard[u], bFrontCard, bMidCard, bBackCard)

	// room.m_bUserCardData[u] = make([]int, 0, 13)
	// room.m_bUserCardData[u] = append(room.m_bUserCardData[u], bFrontCard...)
	// room.m_bUserCardData[u] = append(room.m_bUserCardData[u], bMidCard...)
	// room.m_bUserCardData[u] = append(room.m_bUserCardData[u], bBackCard...)

	room.PlayerSegmentCards[u.ChairId] = append(room.PlayerSegmentCards[u.ChairId], FrontCard, MidCard, BackCard)
	room.PlayerCards[u.ChairId] = make([]int, 0, 13)
	room.PlayerCards[u.ChairId] = append(room.PlayerCards[u.ChairId], FrontCard...)
	room.PlayerCards[u.ChairId] = append(room.PlayerCards[u.ChairId], MidCard...)
	room.PlayerCards[u.ChairId] = append(room.PlayerCards[u.ChairId], BackCard...)

	// btSpecialDataTemp := make([]int, 13)

	// if bSpecialType {
	// 	util.DeepCopy(&btSpecialDataTemp, &btSpecialData)
	// }

	userMgr.ForEachUser(func(user *user.User) {
		user.WriteMsg(&pk_sss_msg.G2C_SSS_Open_Card{CurrentUser: u.ChairId})
	})

	room.OpenCardMap[u] = true
	//log.Debug("%d cccccc", len(room.OpenCardMap))
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
		gameEnd.M_nXShoot = 0
		//CbThreeKillResult      []int      //全垒打加减分
		gameEnd.CbThreeKillResult = make([]int, room.PlayerCount)
		//BEnterExit             bool       //是否一进入就离开
		gameEnd.BEnterExit = false
		//WAllUser               int        //全垒打用户
		gameEnd.WAllUser = 0
		//copy(room.m_lGameScore,room.m_lLeftScore)

		gameEnd.BAllSpecialCard = false

		// nSpecialCard := 0
		// nDragon := 0

		// userMgr.ForEachUser(func(u *user.User) {
		// 	if room.SpecialTypeTable[u] {
		// 		nSpecialCard++
		// 	}
		// 	if room.Dragon[u] {
		// 		nDragon++
		// 	}
		// })

		// if room.PlayerCount == nSpecialCard+nDragon || room.PlayerCount <= nSpecialCard+1 {
		// 	gameEnd.BAllSpecialCard = true
		// } else {
		// 	gameEnd.BAllSpecialCard = false
		// }

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
			// for i := range room.m_bShootState {
			// 	if room.m_bShootState[i][0] != nil {
			// 		gameEnd.ShootState[i][0] = room.m_bShootState[i][0].ChairId

			// 	}
			// 	if room.m_bShootState[i][1] != nil {
			// 		gameEnd.ShootState[i][1] = room.m_bShootState[i][1].ChairId

			// 	}
			// }
		})

		userMgr.ForEachUser(func(u *user.User) {
			u.WriteMsg(gameEnd)
		})
		room.gameEndStatus = gameEnd

		room.AllResult[room.PkBase.TimerMgr.GetPlayCount()] = gameEnd.LGameScore
		room.PkBase.TimerMgr.AddPlayCount()
		//最后一局
		//if room.PkBase.TimerMgr.GetPlayCount() >= room.PkBase.TimerMgr.GetMaxPayCnt() {

		gameRecord := &pk_sss_msg.G2C_SSS_Record{}
		util.DeepCopy(&gameRecord.AllResult, &room.AllResult)
		allScore := make([]int, room.PlayerCount)

		for i := 0; i < room.PkBase.TimerMgr.GetPlayCount(); i++ {
			for j := range allScore {

				allScore[j] += room.AllResult[i][j]
			}
		}
		gameRecord.AllScore = allScore

		userMgr.ForEachUser(func(u *user.User) {
			u.WriteMsg(gameRecord)
		})
		//}

	}

}

// 空闲状态场景
func (room *sss_data_mgr) SendStatusReady(u *user.User) {
	log.Debug("发送空闲状态场景消息")
	StatusFree := &pk_sss_msg.G2C_SSS_StatusFree{
		PlayerCount:      room.PkBase.UserMgr.GetCurPlayerCnt(),
		SubCmd:           room.GameStatus,
		CurrentPlayCount: room.PkBase.TimerMgr.GetPlayCount(),
		MaxPlayCount:     room.PkBase.TimerMgr.GetMaxPayCnt(),
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

	statusPlay.PlayerCount = room.PkBase.UserMgr.GetCurPlayerCnt()
	statusPlay.CurrentPlayCount = room.PkBase.TimerMgr.GetPlayCount()
	statusPlay.MaxPlayCount = room.PkBase.TimerMgr.GetMaxPayCnt()
	statusPlay.Laizi = room.UniversalCards
	statusPlay.PublicCards = room.PublicCards

	u.WriteMsg(statusPlay)
}
