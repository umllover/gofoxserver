package room

import (
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/common"
	"mj/gameServer/common/room_base"
	"mj/gameServer/conf"
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
	room.CellScore = param.CellScore
	room.Password = param.Password
	room.JoinGamePeopleCount = param.JoinGamePeopleCount
	room.CreateUser = uid
	//room.Response = make([]bool, userCnt)
	room.gameLogic = NewGameLogic()
	room.Owner = uid
	room.BankerUser = INVALID_CHAIR
	//room.Record = &msg.G2C_Record{HuCount: make([]int, room.UserCnt), MaCount: make([]int, room.UserCnt), AnGang: make([]int, room.UserCnt), MingGang: make([]int, room.UserCnt), AllScore: make([]int, room.UserCnt), DetailScore: make([][]int, room.UserCnt)}
	now := time.Now().Unix()
	room.TimeStartGame = now
	/*room.CardIndex = make([][]int, room.UserCnt)
	room.HeapCardInfo = make([][]int, room.UserCnt) //堆牌信息
	room.HistoryScores = make([]*HistoryScore, room.UserCnt)
	room.AllowLookon = make(map[int]int)
	room.TurnScore = make([]int, userCnt)
	room.CollectScore = make([]int, userCnt)
	room.Trustee = make([]bool, userCnt)*/
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
	MaxPayCnt       int   //最大局数
	Temp            *tbase.GameServiceOption
	EndTime         *timer.Timer
	KickOut         map[int]*timer.Timer

	ChatRoomId        int                //聊天房间id
	Name              string             //房间名字
	CellScore            int                //底分
	InitCellScore        int                //初始分数
	IsGoldOrGameScore int                //金币场还是积分场 0 标识 金币场 1 标识 积分场

	Password          string             // 密码
	Owner             int                //房主id

	HistoryScores     []*HistoryScore    //历史积分
	AllowLookon       map[int]int        //旁观标志

	/*MaCount           int                //码数，1：一码全中，2-6：对应码数
	Record            *msg.G2C_Record    //约战类型特殊记录
	IsDissumGame      bool               //是否强制解散游戏
	MagicIndex        int                //财神索引
	ProvideCard       int                //供应扑克
	ResumeUser        int                //还原用户
	ProvideUser       int                //供应用户
	LeftCardCount     int                //剩下拍的数量
	EndLeftCount      int                //荒庄牌数
	LastCatchCardUser int                //最后一个摸牌的用户
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
	CardIndex         [][]int            //用户扑克[][MAX_INDEX]
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
	CurrentUser       int                //当前操作用户
	Ting              []bool             //是否听牌
	BankerUser        int                //庄家用户
	TurnScore         []int              //积分信息
	CollectScore      []int              //积分信息
	Trustee           []bool             //是否托管 index 就是椅子id
	PlayCount         int                //已玩局数 */

	CardData[][]		int				//用户扑克
	Count_All			bool			//是不是全部都点了 全点了直接动画
	Qiang				[]bool			//1抢 0不抢
	PlayCount			int						//游戏局数
	PlayerCount			int					//指定游戏人数，2-4

	ChongXuan			bool			//是不是需要重选庄家
	RefreshCfg			bool							//每盘刷新
	RefreshStorageCfg	bool					//是否刷新库存
	StorageOff			int64							//库存减少值
	StorageMax			int64							//库存减少值

	SpecialClient		[]bool          //特殊终端
	BankerUser			int							//庄家用户
	FisrtCallUser		int						//始叫用户
	CurrentUser			int							//当前用户
	ExitScore			int64							//强退分数
	EscapeUserScore		[]int64        //逃跑玩家分数
	DynamicScore		int64                        //总分

	GameRule			bool							//规则标志
	DrawCellScore		int64						//底注积分
	SetCellScoreUser	int 					//设置底分

	IsOpenCard			[]bool				//是否摊牌
	DynamicJoin			[]int           //动态加入
	PlayStatus			[]int			//游戏状态
	CallStatus			[]bool				//叫庄状态
	OxCard				[]int					//牛牛数据
	TableScore			[]int64				//下注数目
	BuckleServiceCharge	[]bool			//收服务费

	HandCardData		[][]int 		//桌面扑克

	TurnMaxScore		[]int64			//最大下注
	MaxScoreTimes		int						//最大倍数

	//HINSTANCE						m_hInst
	//IServerControl*					m_pServerContro
///////////////
	//DWORD							m_dwCheatGameID					//作弊帐号
	//DWORD							m_dwCheatCount						//作弊次数
	//BYTE							m_cbCheatType

	//BYTE							m_cbControl						//是否控制

	/*static bool						m_bAllocConsole
	static LONGLONG					m_lStockScore							//总输赢分
	static LONGLONG					m_lStorageDeduct						//回扣变量
	static LONGLONG					m_lStorageScore						//总抽水分
	LONGLONG						m_lStockLimit							//总输赢分


	TCHAR							m_szConfigFileName[MAX_PATH]		//配置文件
	TCHAR							m_szRoomName[32]					//配置房间

	INT								m_nRobotWinRate					//赢牌几率
	//组件变量
protected:
	CGameLogic						m_GameLogic							//游戏逻辑
	ITableFrame						* m_pITableFrame						//框架接口
	CHistoryScore					m_HistoryScore							//历史成绩
	tagCustomRule *					m_pGameCustomRule						//自定规则

	tagGameServiceOption		    *m_pGameServiceOption					//配置参数
	tagGameServiceAttrib			*m_pGameServiceAttrib					//游戏属性
*/


	/*//不同人数比例下出牌概率控制
protected:
	int								m_nTwCountOneUser						//总人数2人，1个真实玩家，玩家赢牌概率
	int								m_nThCountOneUser
	int								m_nThCountTwoUser
	int								m_nFoCountOneUser
	int								m_nFoCountTwoUser
	int								m_nFoCountThrUser
	int								m_nFvCountOneUser
	int								m_nFvCountTwoUser
	int								m_nFvCountThrUser
	int								m_nFvCountForUser
	int								m_nSxCountOneUser
	int								m_nSxCountTwoUser
	int								m_nSxCountThrUser
	int								m_nSxCountForUser
	int								m_nSxCountFivUser

	//玩家赢分控制
protected:
	LONGLONG						m_lWinScoreLevel[WinScoreLevel]		//赢分库存级别
	int								m_nDecreasePro[WinScoreLevel]			//赢分库存级别，对应概率调整
	*/
	gameLogic         *GameLogic
}

func (r *Room) Destroy() {
	defer func() {
		if r := recover() r != nil {
			log.Recover(r)
		}
	}()
	r.OnDestroy()
	r.CloseSig <- true
	log.Debug("room Room Destroy ok,  Name:%s", r.Name)
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

func (r *Room) GetBirefInfo() *msg.RoomInfo {
	msg := &msg.RoomInfo{}
	msg.ServerID = r.ServerId
	msg.KindID = r.Kind
	msg.NodeID = conf.Server.NodeId
	msg.RoomID = r.GetRoomId()
	msg.CurCnt = r.PlayerCount
	msg.MaxCnt = r.UserCnt           //最多多人数
	msg.PayCnt = r.MaxPayCnt         //可玩局数
	msg.CurPayCnt = r.PlayCount      //已玩局数
	msg.CreateTime = r.TimeStartGame //创建时间
	return msg
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
	} else {
		if room.CheckPlayerCnt() {
			room.Destroy()
		}
	}
}
