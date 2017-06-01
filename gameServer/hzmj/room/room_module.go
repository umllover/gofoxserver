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
)

var (
	idLock sync.RWMutex
	IncId = 0
)

func NewRoom(mgrCh* chanrpc.Server, param *msg.C2G_CreateTable, t *tbase.GameServiceOption, rid, userCnt, uid int) *Room {
	skeleton := base.NewSkeleton()
	Room := new(Room)
	Room.Skeleton = skeleton
	Room.ChanRPC= skeleton.ChanRPCServer
	fmt.Println("@@@@@@@@@@@@@@@@ NewRoom", skeleton.ChanRPCServer)
	Room.mgrCh =mgrCh


	Room.RoomInfo = common.NewRoomInfo(userCnt)
	Room.id = rid
	Room.Kind = t.KindID
	Room.ServerId = t.ServerID
	Room.Name = fmt.Sprintf( strconv.Itoa(common.KIND_TYPE_HZMJ) +"_%v", Room.id)
	Room.CloseSig = make(chan bool, 1)
	Room.TimeLimit = param.DrawTimeLimit
	Room.CountLimit = param.DrawCountLimit
	Room.Source = param.CellScore
	Room.Password = param.Password
	Room.JoinGamePeopleCount = param.JoinGamePeopleCount
	Room.CreateUser = uid
	Room.CustomRule = new(msg.CustomRule)
	Room.Response = make([]bool, userCnt)
	Room.gameLogic = DefaultGameLogic
	Room.EendTime = time.Now().Unix() + 900
	RegisterHandler(Room)
	Room.OnInit()
	go Room.run()
	log.Debug("new room ok .... ")
	return Room
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
	id 			int   //唯一id 房间id
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
	MagicIndex int8//财神索引
	ProvideCard int8 	//供应扑克
	LeftCardCount int8 //剩下拍的数量
	Response []bool						//响应标志
	UserAction []int8			//用户动作
	CardIndex [][]int8		//用户扑克[GAME_PLAYER][MAX_INDEX]
	WeaveItemCount []int8				//组合数目
	WeaveItemArray [][]*msg.WeaveItem;		//组合扑克
	DiscardCount[]int8								//丢弃数目
	DiscardCard[][]int8				//丢弃记录
	OutCardData int8  	//出牌扑克
	OutCardUser int									//当前出牌用户
	HeapHead int									//堆立头部
	HeapTail int									//堆立尾部
	HeapCardInfo [][]int8;						//堆牌信息
	SendCardData int8 					//发牌扑克
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


////////////////// 上面run 和 Destroy 请勿随意修改 //////  下面函数自由操作
func (r *Room) OnInit() {
	r.Skeleton.AfterFunc(10 * time.Second, r.checkDestroyRoom)
}

func (r *Room) OnDestroy() {
	idGenerate.DelRoomId(r.id)
}

func (r *Room) GetRoomId() int{
	return r.id
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



