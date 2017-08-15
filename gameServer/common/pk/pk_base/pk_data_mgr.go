package pk_base

import (
	"strconv"
	"time"

	"mj/common/msg"
	dbase "mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"
)

func NewDataMgr(id int, uid int64, ConfigIdx int, name string, temp *dbase.GameServiceOption, base *Entry_base, info *msg.L2G_CreatorRoom) *RoomData {
	r := new(RoomData)
	r.id = id
	if name == "" {
		r.Name = temp.RoomName
	} else {
		r.Name = name
	}
	r.CreatorUid = uid
	r.CreatorNodeId = info.CreatorNodeId
	r.PkBase = base
	r.ConfigIdx = ConfigIdx

	r.MinPlayerCount = temp.MinPlayer
	r.MaxPlayerCount = temp.MaxPlayer

	log.Debug("new data min player count %d, max %d",
		r.MinPlayerCount, r.MaxPlayerCount)

	r.KindID = temp.KindID
	r.ServerID = temp.ServerID
	r.OtherInfo = info.OtherInfo

	return r
}

//当一张桌子理解
type RoomData struct {
	id            int
	KindID        int
	ServerID      int
	Name          string //房间名字
	CreatorUid    int64  //创建房间的人
	CreatorNodeId int    //创建房间者的NodeId
	PkBase        *Entry_base
	ConfigIdx     int //配置文件索引

	IsGoldOrGameScore int    //金币场还是积分场 0 标识 金币场 1 标识 积分场
	Password          string // 密码

	CellScore  int //底分
	ScoreTimes int //倍数

	InitScoreMap map[int]int // 初始积分

	PlayerCount    int //游戏人数，
	MinPlayerCount int // 最少游戏人数
	MaxPlayerCount int // 最大游戏人数

	FisrtCallUser   int     //始叫用户
	CurrentUser     int     //当前用户
	ExitScore       int64   //强退分数
	EscapeUserScore []int64 //逃跑玩家分数
	DynamicScore    int64   //总分

	EachRoundScoreMap map[int][]int // 每局比分

	CurrentPlayCount int

	OtherInfo map[string]interface{} //其他配置信息
}

func (room *RoomData) GetUserScore(chairid int) int {
	if chairid > room.PkBase.UserMgr.GetMaxPlayerCnt() {
		return 0
	}
	source := 0
	for _, v := range room.EachRoundScoreMap {
		if len(v) < chairid {
			continue
		}
		source += v[chairid]
	}
	return source
}

func (r *RoomData) GetCreatorNodeId() int {
	return r.CreatorNodeId
}

func (r *RoomData) GetCreator() int64 {
	return r.CreatorUid
}

func (r *RoomData) OnCreateRoom() {
	log.Debug("at pk data mgr create room")
	// 初始化积分
	log.Debug("at new data mgr %d %d %d ", r.KindID, r.ServerID, r.PkBase.TimerMgr.GetMaxPlayCnt())

	r.InitScoreMap = make(map[int]int)
	template, ok := dbase.GameServiceOptionCache.Get(r.KindID, r.ServerID)
	if ok {
		log.Debug("get persional table fee ok")
		initScore := template.IniScore
		for i := 0; i < r.MaxPlayerCount; i++ { //初始6个玩家积分1000
			r.InitScoreMap[i] = initScore
		}
	} else {
		for i := 0; i < r.MaxPlayerCount; i++ {
			r.InitScoreMap[i] = 0
		}
	}
	log.Debug("on create room init score map %v", r.InitScoreMap)

	//  每局积分
	r.EachRoundScoreMap = make(map[int][]int)

}

func (room *RoomData) GetCfg() *PK_CFG {
	return GetCfg(room.ConfigIdx)
}

func (room *RoomData) CanOperatorRoom(uid int64) bool {
	if uid == room.CreatorUid {
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
		TableOwnerUserID:  room.CreatorUid,                                               //桌主 I D
		PlayerCnt:         room.PkBase.UserMgr.GetMaxPlayerCnt(),                         //玩家数量
		DrawCountLimit:    room.PkBase.TimerMgr.GetMaxPlayCnt(),                          //局数限制
		DrawTimeLimit:     room.PkBase.TimerMgr.GetTimeLimit(),                           //时间限制
		PlayCount:         room.PkBase.TimerMgr.GetPlayCount(),                           //已玩局数
		PlayTime:          int(room.PkBase.TimerMgr.GetCreatrTime() - time.Now().Unix()), //已玩时间
		CellScore:         room.CellScore,                                                //游戏底分
		IniScore:          0,                                                             //room.IniSource,                                                //初始分数
		ServerID:          strconv.Itoa(room.id),                                         //房间编号
		PayType:           room.PkBase.UserMgr.GetPayType(),                              //支付类型
		IsJoinGame:        0,                                                             //是否参与游戏 todo  tagPersonalTableParameter
		IsGoldOrGameScore: room.IsGoldOrGameScore,                                        //金币场还是积分场 0 标识 金币场 1 标识 积分场
		OtherInfo:         room.OtherInfo,
		LeaveInfo:         room.PkBase.UserMgr.GetLeaveInfo(u.Id), //请求离家的玩家的信息
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
}

func (room *RoomData) InitRoomOne() {

}


//续费后的处理
func (room *RoomData) ResetGameAfterRenewal() {

}

// 游戏开始
func (room *RoomData) BeforeStartGame(UserCnt int) {

}
func (room *RoomData) StartGameing() {

}
func (room *RoomData) AfterStartGame() {

}

// 游戏结束
func (room *RoomData) NormalEnd(cbReason int) {

}
func (room *RoomData) DismissEnd(cbReason int) {

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

func (r *RoomData) ShowCard(u *user.User) {
}

func (r *RoomData) Trustee(u *user.User) {}
