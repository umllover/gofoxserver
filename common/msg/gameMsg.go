package msg

////// c 2 s
//登录游戏服
type C2G_GR_LogonMobile struct {
	GameID         int //游戏标识
	KindID         int
	ServerID       int
	ProcessVersion int //进程版本

	//桌子区域
	DeviceType     int //设备类型
	BehaviorFlags  int //行为标识
	PageTableCount int //分页桌数

	//登录信息
	UserID     int64  //用户 I D
	Password   string //登录密码
	MachineID  string //机器标识
	HallNodeID int
}

//登录成功
type G2C_LogonFinish struct {
}

//登录失败
type G2C_LogonFailure struct {
	ResultCode     int
	DescribeString string
}

//重连游戏服
type C2G_Reconnect struct {
	KindID   int
	ServerID int
	UserID   int64  //用户 I D
	Password string //登录密码
}

//重连游戏服结果
type G2C_ReconnectRsp struct {
	Code int //非0位失败
}

// 请求更换椅子消息
type C2G_GR_UserChairReq struct {
}

//请求创建房间消息
type C2G_LoadTable struct {
	RoomID int //房间id
}

//请求退出房间
type C2G_LeaveRoom struct {
}

//请求退出房间结果
type G2C_LeaveRoomRsp struct {
	Code   int //非0为失败
	Status int // 房间状态 0是没开始， 其他都是开始了
}

//别人退出房间的广播
type G2C_LeaveRoomBradcast struct {
	UserID int64 //用户id
}

//别人同意或拒绝的结果通知
type G2C_ReplyRsp struct {
	UserID int64 //谁同意或者拒绝你了
	Agree  bool  //ture 是同意你了， false 是拒绝你了
}

//同意还是拒绝解散房间
type C2G_ReplyLeaveRoom struct {
	Agree  bool  //true是同意玩家退出， false 是拒绝
	UserID int64 //同意或者拒绝谁
}

type G2C_CancelTable struct{}
type G2C_PersonalTableEnd struct{}

//请求坐下
type C2G_UserSitdown struct {
	TableID  int    // 桌子号码
	ChairID  int    // 椅子号码
	Password string //房间密码
}

//坐下结果
type G2C_UserSitDownRst struct {
	Code int //非0为失败
}

//请求玩家信息
type C2G_REQUserInfo struct {
	UserID   int
	TablePos int
}

//请求房间的基础信息
type C2G_GameOption struct {
	AllowLookon int //旁观标志
}

//房间信息
type G2C_PersonalTableTip struct {
	TableOwnerUserID  int64                  //桌主 I D
	PlayerCnt         int                    //玩家数量
	DrawCountLimit    int                    //局数限制
	DrawTimeLimit     int                    //时间限制
	PlayCount         int                    //已玩局数
	PlayTime          int                    //已玩时间
	CellScore         int                    //游戏底分
	IniScore          int                    //初始分数
	ServerID          string                 //房间编号
	PayType           int                    //1是自己付钱， 2是AA
	IsJoinGame        int                    //是否参与游戏
	IsGoldOrGameScore int                    //金币场还是积分场 0 标识 金币场 1 标识 积分场
	OtherInfo         map[string]interface{} //客户端的配置信息
	LeaveInfo         *LeaveReq              //key 是谁申请退出了，value 是同意的玩家的数组
}

//请求用户信息
type C2G_REQUserChairInfo struct {
	TableID int
	ChairID int
}

//用户起立
type C2G_UserStandup struct {
	TableID    int
	ChairID    int
	ForceLeave int8
}

//用户准备
type C2G_UserReady struct {
	TableID int
	ChairID int
}

//// s 2 c ////////////////////////////

//玩家状态
type G2C_UserStatus struct {
	UserID     int64
	UserStatus *UserStu
}

//请求退出房间的信息
type LeaveReq struct {
	LeftTimes int64
	AgreeInfo []int64
}

//游戏属性 ， 游戏未开始发送的结构
type G2C_StatusFree struct {
	CellScore       int     //基础积分
	TimeOutCard     int     //出牌时间
	TimeOperateCard int     //操作时间
	CreateTime      int64   //开始时间
	TurnScore       []int   //总积分信息 index 是chairId
	CollectScore    [][]int //积分信息 index1 是局数 2是chairID
	PlayerCount     int     //玩家人数
	MaCount         int     //码数
	CountLimit      int     //局数限制
	ZhuaHuaCnt      int     //抓花数
}

//游戏状态 游戏已经开始了发送的结构
type G2C_StatusPlay struct {
	//时间信息
	TimeOutCard     int   //出牌时间
	TimeOperateCard int   //叫分时间
	CreateTime      int64 //开始时间
	PlayCount       int   //已玩局数

	//游戏变量
	CellScore   int   //单元积分
	BankerUser  int   //庄家用户
	CurrentUser int   //当前用户
	MagicIndex  int   //财神索引
	ChaHuaCnt   []int //插花数
	BuHuaCnt    []int //补花数
	BuHuaCard   []int //最新补花卡牌
	ZhuaHuaCnt  int   //抓花数

	//规则
	PlayerCount int //玩家人数
	MaCount     int //码数

	//状态变量
	ActionCard    int    //动作扑克
	ActionMask    int    //动作掩码
	LeftCardCount int    //剩余数目
	Trustee       []bool //是否托管 index 就是椅子id
	Ting          []bool //是否听牌  index chairId

	//出牌信息
	OutCardUser  int     //出牌用户
	OutCardData  int     //出牌扑克
	DiscardCount []int   //丢弃数目
	DiscardCard  [][]int //丢弃记录
	BanOutCard   []int   //禁出卡牌

	//扑克数据
	CardCount    []int //扑克数目
	CardData     []int //扑克列表 room.GetCfg().MaxCount
	SendCardData int   //发送扑克

	//组合扑克
	WeaveItemCount []int          //组合数目
	WeaveItemArray [][]*WeaveItem //组合扑克 [GAME_PLAYER][MAX_WEAVE]

	//堆立信息
	HeapHead     int     //堆立头部
	HeapTail     int     //堆立尾部
	HeapCardInfo [][]int //堆牌信息

	HuCardCount   []int
	HuCardData    [][]int
	OutCardCount  int
	OutCardDataEx []int
	//历史积分
	TurnScore    []int   //总积分信息 index 是chairId
	CollectScore [][]int //积分信息 index1 是局数 2是chairID

}

//约战类型特殊属性
type G2C_Record struct {
	Count       int
	HuCount     []int   //胡牌次数
	MaCount     []int   //中码个数
	AnGang      []int   //暗杠次数
	MingGang    []int   //明杠次数
	AllScore    []int   //总结算分
	DetailScore [][]int //单局结算分
}

// 游戏状态
type G2C_GameStatus struct {
	GameStatus  int //游戏状态
	AllowLookon int //旁观标志
}

//房间配置
type G2C_ConfigServer struct {
	//房间属性
	TableCount int //桌子数目
	ChairCount int //椅子数目

	//房间配置
	ServerType int //房间类型
	ServerRule int //房间规则
}

//发送配置完成
type G2C_ConfigFinish struct {
}

//用户信息
type G2C_UserEnter struct {
	KindID int   //游戏 I D
	UserID int64 //用户 I D

	//头像信息
	FaceID   int8 //头像索引
	CustomID int  //自定标识

	//用户属性
	Gender      int8 //用户性别
	MemberOrder int8 //会员等级

	//用户状态
	TableID    int //桌子索引
	ChairID    int //椅子索引
	UserStatus int //用户状态

	//积分信息
	Score int64 //用户分数

	//游戏信息
	WinCount   int    //胜利盘数
	LostCount  int    //失败盘数
	DrawCount  int    //和局盘数
	FleeCount  int    //逃跑盘数
	Experience int    //用户经验
	NickName   string //昵称
	HeaderUrl  string //头像
	Sign       string //签名
	Star       int    //点赞数
}

type SysMsg struct {
	ClientID int64
	Type     int
	Context  string
}

type G2C_Hu_Data struct {
	//出哪几张能听
	OutCardCount int
	OutCardData  []int
	//听后能胡哪几张牌
	HuCardCount []int
	HuCardData  [][]int
	//胡牌剩余数
	HuCardRemainingCount [][]int
}

///////////////////////// game chart begin ///////////////////////////////
type C2G_GameChart_ToAll struct {
	ChatColor  int    //字体颜色
	SendUserID int    //发送者id
	ChatString string //消息内容
	ChatIndex  int    //第几条消息
	ChatType   int    //1是语音 0 是普通聊天
}

type G2C_GameChart_ToAll struct {
	ChatColor    int    //颜色 无效
	SendUserID   int64  //谁发的消息
	TargetUserID int    //发给谁的消息  无效
	ClientID     int    //无效
	ChatIndex    int    //消息的下标， 如 1 我的等的花都谢了  2快点吗
	ChatString   string //消息内容 ，
	ChatType     int    //1是语音 0 是普通聊天
}

//结束消息， 各个游戏自己实现
type G2C_GameConclude struct {
}

//踢出玩家

type G2C_KickOut struct {
	Reason int //踢出原因 1是服务器主动踢出， 2是踢号
}

//补花
type C2G_ReplaceCard struct {
	CardData int //扑克数据
}

//补花
type G2C_ReplaceCard struct {
	ReplaceUser  int //补牌用户
	ReplaceCard  int //补牌扑克
	NewCard      int //补完扑克
	IsInitFlower bool //是否开局补花，true开局补花
}

///////////////////////// game chart end ///////////////////////////////
