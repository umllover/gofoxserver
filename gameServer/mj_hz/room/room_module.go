package room

import (
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/common"
	"mj/gameServer/common/room_base"
	tbase "mj/gameServer/db/model/base"
	"mj/gameServer/idGenerate"
	"mj/gameServer/user"
	"strconv"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/timer"
)

func NewRoom(mgrCh *chanrpc.Server, param *msg.C2G_CreateTable, t *tbase.GameServiceOption, rid, userCnt, uid int) *Room {
	room := new(Room)
	room.RoomBase = room_base.NewRoomBase(userCnt, rid, mgrCh, fmt.Sprintf(strconv.Itoa(common.KIND_TYPE_HZMJ)+"_%v", rid))
	room.Kind = t.KindID
	room.ServerId = t.ServerID
	room.Temp = t
	room.CloseSig = make(chan bool, 1)
	room.TimeLimit = param.DrawTimeLimit
	room.CountLimit = param.DrawCountLimit
	room.Source = param.CellScore
	room.Password = param.Password
	room.JoinGamePeopleCount = param.JoinGamePeopleCount
	room.CreateUser = uid
	room.Response = make([]bool, userCnt)
	room.gameLogic = NewGameLogic()
	room.Owner = uid
	room.BankerUser = INVALID_CHAIR
	room.Record = &msg.G2C_Record{HuCount: make([]int, room.UserCnt), MaCount: make([]int, room.UserCnt), AnGang: make([]int, room.UserCnt), MingGang: make([]int, room.UserCnt), AllScore: make([]int, room.UserCnt), DetailScore: make([][]int, room.UserCnt)}
	now := time.Now().Unix()
	room.TimeStartGame = now
	room.CardIndex = make([][]int, room.UserCnt)
	room.HeapCardInfo = make([][]int, room.UserCnt) //堆牌信息
	room.HistoryScores = make([]*HistoryScore, room.UserCnt)
	room.AllowLookon = make(map[int]int)
	room.TurnScore = make([]int, userCnt)
	room.CollectScore = make([]int, userCnt)
	room.Trustee = make([]bool, userCnt)
	room.KickOut = make(map[int]*timer.Timer)
	RegisterHandler(room)
	room.OnInit()
	room.RoomRun()
	if room.Temp.TimeNotBeginGame != 0 {
		room.EndTime = room.Skeleton.AfterFunc(time.Duration(room.Temp.TimeNotBeginGame)*time.Second, room.AfterNotBegin)
	}

	log.Debug("new room ok .... ")
	return room
}

//吧room 当一张桌子理解
type Room struct {
	*room_base.RoomBase
	Kind            int   //第一类型
	ServerId        int   //第二类型 注意 非房间id
	TimeLimit       int   //时间显示
	CountLimit      int   //局数限制
	TimeOutCard     int   //出牌时间
	TimeOperateCard int   //操作时间
	TimeStartGame   int64 //开始时间
	Status          int   //当前状态
	Temp            *tbase.GameServiceOption
	EndTime         *timer.Timer
	KickOut         map[int]*timer.Timer

	ChatRoomId        int                //聊天房间id
	Name              string             //房间名字
	Source            int                //底分
	IniSource         int                //初始分数
	IsGoldOrGameScore int                //金币场还是积分场 0 标识 金币场 1 标识 积分场
	Password          string             // 密码
	MaCount           int                //码数，1：一码全中，2-6：对应码数
	Record            *msg.G2C_Record    //约战类型特殊记录
	IsDissumGame      bool               //是否强制解散游戏
	MagicIndex        int                //财神索引
	ProvideCard       int                //供应扑克
	ResumeUser        int                //还原用户
	ProvideUser       int                //供应用户
	LeftCardCount     int                //剩下拍的数量
	EndLeftCount      int                //荒庄牌数
	LastCatchCardUser int                //最后一个摸牌的用户
	Owner             int                //房主id
	OutCardCount      int                //出牌数目
	ChiHuCard         int                //吃胡扑克
	MinusHeadCount    int                //头部空缺
	MinusLastCount    int                //尾部空缺
	SiceCount         int                //色子大小
	SendCardCount     int                //发牌数目
	UserActionDone    bool               //操作完成
	SendStatus        int                //发牌状态
	GangStatus        int                //杠牌状态
	GangOutCard       bool               //杠后出牌
	ProvideGangUser   int                //供杠用户
	GangCard          []bool             //杠牌状态
	GangCount         []int              //杠牌次数
	RepertoryCard     []int              //库存扑克
	UserGangScore     []int              //游戏中杠的输赢
	Response          []bool             //响应标志
	ChiHuKind         []int              //吃胡结果
	ChiHuRight        []int              //胡牌类型
	UserMaCount       []int              //下注用户数
	UserAction        []int              //用户动作
	OperateCard       [][]int            //操作扑克
	ChiPengCount      []int              //吃碰杠次数
	PerformAction     []int              //执行动作
	HandCardCount     []int              //扑克数目
	CardIndex         [][]int            //用户扑克[GAME_PLAYER][MAX_INDEX]
	WeaveItemCount    []int              //组合数目
	WeaveItemArray    [][]*msg.WeaveItem //组合扑克
	DiscardCount      []int              //丢弃数目
	DiscardCard       [][]int            //丢弃记录
	OutCardData       int                //出牌扑克
	OutCardUser       int                //当前出牌用户
	HeapHead          int                //堆立头部
	HeapTail          int                //堆立尾部
	HeapCardInfo      [][]int            //堆牌信息
	SendCardData      int                //发牌扑克
	HistoryScores     []*HistoryScore    //历史积分
	CurrentUser       int                //当前操作用户
	Ting              []bool             //是否听牌
	BankerUser        int                //庄家用户
	AllowLookon       map[int]int        //旁观标志
	TurnScore         []int              //积分信息
	CollectScore      []int              //积分信息
	Trustee           []bool             //是否托管 index 就是椅子id
	PlayCount         int                //已玩局数
	gameLogic         *GameLogic
}

func (r *Room) GetCurlPlayerCount() int {
	cnt := 0
	for _, u := range r.Users {
		if u != nil {
			cnt++
		}
	}

	return cnt
}

func (r *Room) OnInit() {
	r.Skeleton.AfterFunc(10*time.Second, r.Update)
}

func (r *Room) OnDestroy() {
	idGenerate.DelRoomId(r.GetRoomId())
	r.MgrCh.Go("DelRoom", r.GetRoomId())
}

func (r *Room) Update() {
	r.Skeleton.AfterFunc(10*time.Second, r.Update)
}

/////////////////// 超时处理函数
//多久没开始解散房间
func (room *Room) AfterNotBegin() {
	if room.Status == RoomStatusReady {
		room.OnEventGameConclude(0, nil, GER_DISMISS)
	}
}

//开始多久没打完界山房间
func (room *Room) AfterGameTimeOut() {
	room.OnEventGameConclude(0, nil, GER_DISMISS)
}

//玩家离线超时踢出
func (room *Room) OfflineKickOut(user *user.User) {
	room.LeaveRoom(user)
	if room.Status != RoomStatusReady {
		room.OnEventGameConclude(0, nil, GER_DISMISS)
	}
}
