package room

import (
	"github.com/lovelly/leaf/module"
	"mj/gameServer/base"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"time"
	"mj/gameServer/common"
	"fmt"
	"strconv"
	"sync"
	"mj/common/msg"
	tbase "mj/gameServer/db/model/base"
	"mj/gameServer/idGenerate"
	. "mj/common/cost"
)

var (
	idLock sync.RWMutex
	IncId = 0
)

func NewRoom(mgrCh* chanrpc.Server, param *msg.C2G_CreateTable, t *tbase.GameServiceOption, rid, userCnt, uid int) *Room {
	skeleton := base.NewSkeleton()
	room := new(Room)
	room.Skeleton = skeleton
	room.ChanRPC= skeleton.ChanRPCServer
	room.mgrCh =mgrCh
	room.RoomInfo = common.NewRoomInfo(userCnt, rid)
	room.Kind = t.KindID
	room.ServerId = t.ServerID
	room.Name = fmt.Sprintf( strconv.Itoa(common.KIND_TYPE_HZMJ) +"_%v", room.GetRoomId())
	room.CloseSig = make(chan bool, 1)
	room.TimeLimit = param.DrawTimeLimit
	room.CountLimit = param.DrawCountLimit
	room.Source = param.CellScore
	room.Password = param.Password
	room.JoinGamePeopleCount = param.JoinGamePeopleCount
	room.CreateUser = uid
	room.CustomRule = new(msg.CustomRule)
	room.Response = make([]bool, userCnt)
	room.gameLogic = DefaultGameLogic
	room.EendTime = time.Now().Unix() + 900
	room.Owner = uid
	room.BankerUser = INVALID_CHAIR
	room.Record = &msg.G2C_Record{}

	room.CardIndex = make([][]uint8, room.UserCnt)
	room.HeapCardInfo  = make([][]uint8,room.UserCnt)			//堆牌信息
	room.HistoryScores  = make([]*HistoryScore,room.UserCnt)
	RegisterHandler(room)
	room.OnInit()
	go room.run()
	log.Debug("new room ok .... ")
	return room
}

//吧room 当一张桌子理解
type Room struct {
	// module 必须字段
	*module.Skeleton
	ChanRPC *chanrpc.Server //接受客户端消息的chan
	mgrCh* chanrpc.Server  //管理类的chan 例如红中麻将 就是红中麻将module的 ChanRPC
	CloseSig  chan bool
	wg       sync.WaitGroup //

	// 游戏字段
	*common.RoomInfo
	Name          string  //房间名字
	Kind 		int  //第一类型
	ServerId    int  //第二类型 注意 非房间id
	Source int //底分
	IniSource int //初始分数
	TimeLimit int //时间显示
	CountLimit int //局数限制
	IsGoldOrGameScore int //金币场还是积分场 0 标识 金币场 1 标识 积分场
	Password string 	// 密码
	JoinGamePeopleCount int	 //参与游戏的人数
	*msg.CustomRule  //自定义规则
	Record *msg.G2C_Record //约战类型特殊记录
	IsDissumGame bool					//是否强制解散游戏
	MagicIndex uint8//财神索引
	ProvideCard uint8 	//供应扑克
	ResumeUser int									//还原用户
	ProvideUser int		//供应用户
	LeftCardCount uint8 //剩下拍的数量
	EndLeftCount uint8 								//荒庄牌数
	LastCatchCardUser int						//最后一个摸牌的用户
	Owner int 			//房主id
	OutCardCount uint8								//出牌数目
	ChiHuCard uint8									//吃胡扑克
	MinusHeadCount uint8								//头部空缺
	MinusLastCount uint8								//尾部空缺
	SiceCount int 										//色子大小
	SendCardCount uint8									//发牌数目
	UserActionDone bool
	SendStatus uint8									//发牌状态
	GangStatus uint8									//杠牌状态
	GangOutCard bool									//杠后出牌
	ProvideGangUser int									//供杠用户
	GangCard []bool						//杠牌状态
	GangCount []uint8						//杠牌次数
	RepertoryCard []uint8								//库存扑克
	UserGangScore []int									//游戏中杠的输赢
	Response []bool										//响应标志
	ChiHuKind []int									//吃胡结果
	ChiHuRight []int								//胡牌类型
	UserMaCount []uint8
	UserAction []uint8								//用户动作
	OperateCard	[][]uint8				//操作扑克
	ChiPengCount []uint8		//吃碰杠次数
	PerformAction []uint8							//执行动作
	HandCardCount []uint8							//扑克数目
	CardIndex [][]uint8								//用户扑克[GAME_PLAYER][MAX_INDEX]
	WeaveItemCount []uint8							//组合数目
	WeaveItemArray [][]*msg.WeaveItem;				//组合扑克
	DiscardCount[]uint8								//丢弃数目
	DiscardCard[][]uint8							//丢弃记录
	OutCardData uint8  								//出牌扑克
	OutCardUser int									//当前出牌用户
	HeapHead int									//堆立头部
	HeapTail int									//堆立尾部
	HeapCardInfo [][]uint8;							//堆牌信息
	SendCardData uint8 								//发牌扑克
	gameLogic *GameLogic
	HistoryScores  []*HistoryScore
}

func (r *Room)run(){
	log.Debug("room Room start run Name:%s", r.Name)
	r.Run(r.CloseSig)
	log.Debug("room Room End run Name:%s", r.Name)
}

func  (r *Room) Destroy(){
	defer func() {
		if r := recover(); r != nil {
			log.Recover(r)
		}
	}()

	r.CloseSig <- true
	r.OnDestroy()
	log.Debug("room Room Destroy ok,  Name:%s", r.Name)
}

func (r *Room)  GetCurlPlayerCount() int {
	cnt := 0
	for _, u := range r.Users {
		if u != nil {
			cnt ++
		}
	}

	return cnt
}


////////////////// 上面run 和 Destroy 请勿随意修改 //////  下面函数自由操作
func (r *Room) OnInit() {
	r.Skeleton.AfterFunc(10 * time.Second, r.checkDestroyRoom)
}

func (r *Room) OnDestroy() {
	idGenerate.DelRoomId(r.GetRoomId())
}


//这里添加定时操作
func (r *Room) checkDestroyRoom() {
	nowTime := time.Now().Unix()
	if r.CheckDestroy(nowTime) {
		r.Destroy()
		return
	}

	r.Skeleton.AfterFunc(10 * time.Second, r.checkDestroyRoom)
}

func (r *Room) GetChanRPC() *chanrpc.Server {
	return r.ChanRPC
}



