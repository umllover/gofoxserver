package pk_base

import (
	"strconv"
	"time"

	"mj/common/cost"
	"mj/common/msg"
	dbase "mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"
)

func NewDataMgr(id int, uid int64, ConfigIdx int, name string, temp *dbase.GameServiceOption, base *Entry_base) *RoomData {
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
	r.PlayerCount = temp.MaxPlayer

	r.KindID = temp.KindID
	r.ServerID = temp.ServerID

	return r
}

//当一张桌子理解
type RoomData struct {
	id         int
	KindID     int
	ServerID   int
	Name       string //房间名字
	CreateUser int64  //创建房间的人
	PkBase     *Entry_base
	ConfigIdx  int //配置文件索引

	IsGoldOrGameScore int    //金币场还是积分场 0 标识 金币场 1 标识 积分场
	Password          string // 密码

	CellScore  int //底分
	ScoreTimes int //倍数

	InitScoreMap map[int]int // 初始积分

	PlayerCount int //游戏人数，

	FisrtCallUser   int     //始叫用户
	CurrentUser     int     //当前用户
	ExitScore       int64   //强退分数
	EscapeUserScore []int64 //逃跑玩家分数
	DynamicScore    int64   //总分

	EachRoundScoreMap map[int][]int // 每局比分

	HistoryScores    []*HistoryScore //历史积分
	CurrentPlayCount int
}

func (r *RoomData) OnCreateRoom() {
	log.Debug("at pk data mgr create room")
	// 初始化积分
	log.Debug("at new data mgr %d %d %d ", r.KindID, r.ServerID, r.PkBase.TimerMgr.GetMaxPayCnt())

	r.InitScoreMap = make(map[int]int)
	persionalTalbleFeeCache := dbase.PersonalTableFeeCache
	persionalTableFee, ok := persionalTalbleFeeCache.Get(r.KindID, r.ServerID, r.PkBase.TimerMgr.GetMaxPayCnt())
	if ok {
		log.Debug("get persional table fee ok")
		initScore := persionalTableFee.IniScore
		for i := 0; i < r.PlayerCount; i++ { //每个玩家初始积分1000
			r.InitScoreMap[i] = initScore
		}
	} else {
		for i := 0; i < r.PlayerCount; i++ {
			r.InitScoreMap[i] = 1000
		}
	}

	//  每局积分
	r.EachRoundScoreMap = make(map[int][]int)

}

func (room *RoomData) GetCfg() *PK_CFG {
	return GetCfg(room.ConfigIdx)
}

func (room *RoomData) GetCreater() int64 {
	return room.CreateUser
}

// 其它操作，各个游戏自己有自己的游戏指令
func (room *RoomData) OtherOperation(args []interface{}) {

}
func (room *RoomData) CanOperatorRoom(uid int64) bool {
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
		DrawCountLimit:    room.PkBase.TimerMgr.GetMaxPayCnt(),                           //局数限制
		DrawTimeLimit:     room.PkBase.TimerMgr.GetTimeLimit(),                           //时间限制
		PlayCount:         room.PkBase.TimerMgr.GetPlayCount(),                           //已玩局数
		PlayTime:          int(room.PkBase.TimerMgr.GetCreatrTime() - time.Now().Unix()), //已玩时间
		CellScore:         room.CellScore,                                                //游戏底分
		IniScore:          0,                                                             //room.IniSource,                                                //初始分数
		ServerID:          strconv.Itoa(room.id),                                         //房间编号
		IsJoinGame:        0,                                                             //是否参与游戏 todo  tagPersonalTableParameter
		IsGoldOrGameScore: room.IsGoldOrGameScore,                                        //金币场还是积分场 0 标识 金币场 1 标识 积分场
	})
}

// 设置底分
func (room *RoomData) SetCellScore(cellScore int) {
	room.CellScore = cellScore
}

// 设置倍数
func (room *RoomData) SetScoreTimes(scoreTimes int) {
	room.ScoreTimes = scoreTimes
}

func (room *RoomData) InitRoom(UserCnt int) {
	room.PlayerCount = UserCnt
	room.CellScore = room.PkBase.Temp.CellScore
}

// 游戏开始
func (room *RoomData) BeforeStartGame(UserCnt int) {

}
func (room *RoomData) StartGameing() {

}
func (room *RoomData) AfterStartGame() {

}
func (room *RoomData) AfertEnd(Forced bool) {
	log.Debug("ggggggggggggg")
}

// 游戏结束
func (room *RoomData) NormalEnd() {

}
func (room *RoomData) DismissEnd() {

}

func (room *RoomData) SendStatusPlay(u *user.User) {

}
func (room *RoomData) SendStatusReady(u *user.User) {

}

// 叫分 加注 亮牌
func (room *RoomData) CallScore(u *user.User, scoreTimes int) {

}
func (room *RoomData) AddScore(u *user.User, score int) {

}
func (room *RoomData) OpenCard(u *user.User, cardType int, cardData []int) {

}

func (r *RoomData) ShowCard(u *user.User) {
}

func (r *RoomData) Trustee(u *user.User) {}

func (r *RoomData) AfterEnd(Forced bool) {
	log.Debug("at pk data mgr after end")
	r.PkBase.TimerMgr.AddPlayCount()
	if Forced || r.PkBase.TimerMgr.GetPlayCount() >= r.PkBase.TimerMgr.GetMaxPayCnt() {
		log.Debug("Forced :%v, PlayTurnCount:%v, temp PlayTurnCount:%d", Forced, r.PkBase.TimerMgr.GetPlayCount(), r.PkBase.TimerMgr.GetMaxPayCnt())
		r.PkBase.UserMgr.SendCloseRoomToHall(&msg.RoomEndInfo{
			RoomId: r.PkBase.DataMgr.GetRoomId(),
			Status: r.PkBase.Status,
		})
		r.PkBase.Destroy(r.PkBase.DataMgr.GetRoomId())
		r.PkBase.UserMgr.RoomDissume()

		return
	}

	r.PkBase.UserMgr.ForEachUser(func(u *user.User) {
		r.PkBase.UserMgr.SetUsetStatus(u, cost.US_FREE)
	})

}
func (room *RoomData) ShowSSSCard(u *user.User, bDragon bool, bSpecialType bool, btSpecialData []int, bFrontCard []int, bMidCard []int, bBackCard []int) {

}
