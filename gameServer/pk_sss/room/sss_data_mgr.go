package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	. "mj/common/cost"
	"mj/common/msg/pk_sss_msg"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

// 游戏状态
const (
	GAME_START = 1002 // 游戏开始
)

func NewDataMgr(id, uid, ConfigIdx int, name string, temp *base.GameServiceOption, base *SSS_Entry) *sss_data_mgr {
	d := new(sss_data_mgr)
	d.RoomData = pk_base.NewDataMgr(id, uid, ConfigIdx, name, temp, base.Entry_base)
	return d
}

type sss_data_mgr struct {
	*pk_base.RoomData

	//游戏变量
	bCardData               []int     //牌的总数
	m_bUserCardData         [][]int   //玩家扑克
	m_bSegmentCard          [][][]int //分段扑克
	m_bFinishSegment        []int     //完成分段
	m_bShowCardCount        int       //摊牌数目
	m_bCompleteCompareCount int       //完成比较
	m_bOverTime             []int     //摊牌超时
	m_bUserLeft             []int     //玩家强退
	m_bSpecialTypeTable     []int     //是否特殊牌型
	m_nPlayerCount          int       //实际玩家人数

	LeftCardCount int //剩下拍的数量
	//比较结果
	m_bCompareResult        [][]int //每一道比较结果
	m_bShootState           [][]int //打枪(0赢的玩家,1输的玩家)
	m_bThreeKillResult      []int   //全垒打加减分
	m_bToltalWinDaoShu      []int   //总共道数
	m_bCompareDouble        []int   //打枪的道数
	m_bSpecialCompareResult []int   //特殊牌型比较结果
	m_lGameScore            []int   //游戏积分
	m_nXShoot               int     //几家打枪

	// 游戏状态
	GameStatus int
}

func (room *sss_data_mgr) InitRoom(UserCnt int) {
	//初始化
	log.Debug("初始化房间")

	room.m_bCompareDouble = make([]int, UserCnt)
	room.m_bCompareResult = make([][]int, UserCnt)
	room.m_bShootState = make([][]int, UserCnt)
	room.m_bSpecialCompareResult = make([]int, UserCnt)
	room.m_bThreeKillResult = make([]int, UserCnt)
	room.m_bToltalWinDaoShu = make([]int, UserCnt)
	room.m_lGameScore = make([]int, UserCnt)
	room.m_nXShoot = 0

	room.LeftCardCount = room.GetCfg().MaxRepertory
}

func (room *sss_data_mgr) BeforeStartGame(UserCnt int) {
	room.GameStatus = GAME_START
	room.InitRoom(UserCnt)
}

func (room *sss_data_mgr) StartGameing() {
	room.StartDispatchCard()
}
func (r *sss_data_mgr) GetOneCard() int { // 从牌堆取出一张
	r.LeftCardCount -= 1
	return r.bCardData[r.LeftCardCount]
}
func (room *sss_data_mgr) StartDispatchCard() {
	log.Debug("begin start game sss")
	userMgr := room.PkBase.UserMgr
	gameLogic := room.PkBase.LogicMgr

	gameLogic.RandCardList(room.bCardData, pk_base.GetSSSCards())

	userMgr.ForEachUser(func(u *user.User) {
		userMgr.SetUsetStatus(u, US_PLAYING)
	})

	userMgr.ForEachUser(func(u *user.User) {
		for i := 0; i < 13; i++ {
			room.m_bUserCardData[u.ChairId][i] = room.GetOneCard()
		}
	})

	userMgr.ForEachUser(func(u *user.User) {
		SendCard := &pk_sss_msg.G2C_SSS_SendCard{}
		util.DeepCopy(&SendCard.AllHandCardData, &room.m_bUserCardData)
		SendCard.CellScore = room.CellScore
		u.WriteMsg(SendCard)
	})
}
