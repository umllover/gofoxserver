package room

import (
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/mj_xs_msg"
	"mj/gameServer/common"
	"mj/gameServer/common/room_base"
	tbase "mj/gameServer/db/model/base"
	"mj/gameServer/idGenerate"
	"strconv"
	"time"

	"mj/gameServer/common/mj_ctl_base"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
)

func NewRoom(mgrCh *chanrpc.Server, param *msg.C2G_CreateTable, t *tbase.GameServiceOption, rid, userCnt, uid int) *Room {
	room := new(Room)
	room.RoomBase = room_base.NewRoomBase(userCnt, rid, mgrCh, fmt.Sprintf(strconv.Itoa(common.KIND_TYPE_XSMJ)+"_%v", rid))
	room.Kind = t.KindID
	room.ServerId = t.ServerID
	room.CloseSig = make(chan bool, 1)
	room.TimeLimit = param.DrawTimeLimit
	room.CountLimit = param.DrawCountLimit
	room.Source = param.CellScore
	room.Password = param.Password
	room.CreateUser = uid
	room.Response = make([]bool, userCnt)
	room.gameLogic = NewGameLogic()
	room.Owner = uid
	room.BankerUser = INVALID_CHAIR
	now := time.Now().Unix()
	room.TimeStartGame = now
	room.EendTime = now + 900
	room.CardIndex = make([][]int, room.UserCnt)
	room.HistoryScores = make([]*HistoryScore, room.UserCnt)
	RegisterHandler(room)

	log.Debug("new room ok .... ")
	return room
}

//吧room 当一张桌子理解
type Room struct {
	// 游戏字段
	*room_base.RoomBase
	*mj_ctl_base.Ctl_base
	ChatRoomId        int                         //聊天房间id
	Name              string                      //房间名字
	Kind              int                         //第一类型
	ServerId          int                         //第二类型 注意 非房间id
	Source            int                         //底分
	IniSource         int                         //初始分数
	TimeLimit         int                         //时间显示
	CountLimit        int                         //局数限制
	IsGoldOrGameScore int                         //金币场还是积分场 0 标识 金币场 1 标识 积分场
	Password          string                      // 密码
	ProvideCard       int                         //供应扑克
	ResumeUser        int                         //还原用户
	ProvideUser       int                         //供应用户
	LeftCardCount     int                         //剩下拍的数量
	Owner             int                         //房主id
	OutCardCount      int                         //出牌数目
	ChiHuCard         int                         //吃胡扑克
	SiceCount         int                         //色子大小
	SendCardCount     int                         //发牌数目
	RepertoryCard     []int                       //库存扑克
	Response          []bool                      //响应标志
	ChiHuKind         []int                       //吃胡结果
	ChiHuRight        []int                       //胡牌类型
	UserAction        []int                       //用户动作
	ChiPengCount      []int                       //吃碰杠次数
	PerformAction     []int                       //执行动作
	CardIndex         [][]int                     //用户扑克[GAME_PLAYER][MAX_INDEX]
	WeaveItemCount    []int                       //组合数目
	WeaveItemArray    [][]*mj_xs_msg.TagWeaveItem //组合扑克
	DiscardCount      []int                       //丢弃数目
	DiscardCard       [][]int                     //丢弃记录
	OutCardData       int                         //出牌扑克
	OutCardUser       int                         //当前出牌用户
	SendCardData      int                         //发牌扑克

	GangStatus        bool              //杠牌状态
	SendStatus        bool              //发牌状态
	OperateCard       []int             //操作扑克
	FengQuan          int               //圈风
	Zfb               []int             //中发白的碰刻杠
	Dnxb              []int             //东南西北的碰刻杠
	UserWindCount     []int             //花牌个数
	UserWindData      [][]int           //花牌数据	 已出
	AllWindCount      int               //花牌计数
	WindCount         []int             //临时花牌个数
	TempWinCount      []int             //临时花牌个数
	WindData          [][]int           ///临时花牌数据
	SumWindCount      int               //开始总花牌
	RemainCardCount   int               //预留扑克
	GangCount         int               //杠牌次数
	ChiHuResult       []*TagChiHuResult //吃胡结果
	RealDnxb          []int             //手上的东南西北
	RealZfb           []int             //手上中发白
	Change            bool              //是否换圈
	EnjoinCardData    [][]int           //禁止出牌 GAME_PLAYER
	EnjoinCardCount   []int             //禁止出牌GAME_PLAYER
	EnjoinPengCard    []int             //禁止碰牌
	EnjoinHuCard      []int             //禁止胡牌
	InitialBankerUser int               //初始庄家
	gameLogic         *GameLogic        //逻辑类
	HistoryScores     []*HistoryScore   //历史积分
	Ting              []bool            //是否听牌
	BankerUser        int               //庄家用户
	CreateUser        int               //创建房间的人
	Status            int               //当前状态
	PlayCount         int               //已玩局数
	AllowLookon       map[int]int       //旁观玩家
	TurnScore         []int             //积分信息
	CollectScore      []int             //积分信息
	Trustee           []bool            //是否托管 index 就是椅子id
	CurrentUser       int               //当前玩家
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
	r.Skeleton.AfterFunc(10*time.Second, r.checkDestroyRoom)
}

func (r *Room) OnDestroy() {
	idGenerate.DelRoomId(r.GetRoomId())
	r.MgrCh.Go("DelRoom", r.GetRoomId())
}

//这里添加定时操作
func (r *Room) checkDestroyRoom() {
	nowTime := time.Now().Unix()
	if r.CheckDestroy(nowTime) {
		r.Destroy()
		return
	}

	r.Skeleton.AfterFunc(10*time.Second, r.checkDestroyRoom)
}
