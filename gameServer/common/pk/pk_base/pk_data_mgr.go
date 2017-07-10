package pk_base

import (
	"strconv"
	"time"

	"mj/common/msg"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
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
	PkBase     *Entry_base
	ConfigIdx  int //配置文件索引

	IsGoldOrGameScore int    //金币场还是积分场 0 标识 金币场 1 标识 积分场
	Password          string // 密码

	CellScore  int //底分
	ScoreTimes int //倍数

	PlayCount   int //游戏局数
	PlayerCount int //指定游戏人数，2-4

	FisrtCallUser   int     //始叫用户
	CurrentUser     int     //当前用户
	ExitScore       int64   //强退分数
	EscapeUserScore []int64 //逃跑玩家分数
	DynamicScore    int64   //总分

	//历史积分
	HistoryScores []*HistoryScore //历史积分
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
	log.Debug("pk data mgr init")
	room.CellScore = room.PkBase.Temp.CellScore
}

// 游戏开始
func (room *RoomData) BeforeStartGame(UserCnt int) {

}
func (room *RoomData) StartGameing() {

}
func (room *RoomData) AfterStartGame() {

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

// 其它操作，各个游戏自己有自己的游戏指令
func (room *RoomData) OtherOperation(args []interface{}) {

}
func (room *RoomData) ShowSSSCard(u *user.User, bDragon bool, bSpecialType bool, btSpecialData []int, bFrontCard []int, bMidCard []int, bBackCard []int) {

}
