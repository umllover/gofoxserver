package msg

////// c 2 s
//手机登录
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
	UserID    int    //用户 I D
	Password  string //登录密码
	MachineID string //机器标识
}

// 请求更换椅子消息
type C2G_GR_UserChairReq struct {
}

//请求创建房间消息
type C2G_CreateTable struct {
	CellScore           int    //底分设置
	DrawCountLimit      int    //局数限制
	DrawTimeLimit       int    //时间限制
	JoinGamePeopleCount int    //参与游戏的人数
	RoomTax             int    //单独一个私人房间的税率，千分比
	Password            string //密码设置
	GameRule            []int8 //游戏规则 弟 0 位标识 是否设置规则 0 代表未设置 1 代表设置
	Kind                int    //游戏类型
	ServerId            int    //子类型
}

//请求坐下
type C2G_UserSitdown struct {
	TableID  int    // 桌子号码
	ChairID  int    // 椅子号码
	Password string //房间密码
}

//请求玩家信息
type C2G_REQUserInfo struct {
	UserID   int
	TablePos int
}

//配置信息
type C2G_GameOption struct {
	AllowLookon   int //旁观标志
	FrameVersion  int //框架版本
	ClientVersion int //游戏版本
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
//登录成功
type G2C_LogonFinish struct {
}

//登录失败
type G2C_LogonFailur struct {
	ResultCode     int
	DescribeString string
}

// 创建房间失败消息
type G2C_CreateTableFailure struct {
	ErrorCode      int
	DescribeString string
}

//创建房间成功的消息
type G2C_CreateTableSucess struct {
	TableID        int //房间编号
	DrawCountLimit int //局数限制
	DrawTimeLimit  int //时间限制
	Beans          int //游戏豆
	RoomCard       int //房卡数量
}

//查询房间的结果
type G2C_SearchResult struct {
	ServerID int //房间 I D
	TableID  int //桌子 I D
}

//玩家状态
type G2C_UserStatus struct {
	UserID     int
	UserStatus *UserStu
}

//发送提示信息
type G2C_PersonalTableTip struct {
	TableOwnerUserID  int    //桌主 I D
	DrawCountLimit    int    //局数限制
	DrawTimeLimit     int    //时间限制
	PlayCount         int    //已玩局数
	PlayTime          int    //已玩时间
	CellScore         int    //游戏底分
	IniScore          int    //初始分数
	ServerID          string //房间编号
	IsJoinGame        int    //是否参与游戏
	IsGoldOrGameScore int    //金币场还是积分场 0 标识 金币场 1 标识 积分场
}

//游戏属性 ， 游戏未开始发送的结构
type G2C_StatusFree struct {
	CellScore       int   //基础积分
	TimeOutCard     int   //出牌时间
	TimeOperateCard int   //操作时间
	TimeStartGame   int64 //开始时间
	TurnScore       []int //积分信息
	CollectScore    []int //积分信息
	PlayerCount     int   //玩家人数
	MaCount         int   //码数
	CountLimit      int   //局数限制
}

//游戏状态 游戏已经开始了发送的结构
type G2C_StatusPlay struct {
	//时间信息
	TimeOutCard     int   //出牌时间
	TimeOperateCard int   //叫分时间
	TimeStartGame   int64 //开始时间

	//游戏变量
	CellScore   int //单元积分
	BankerUser  int //庄家用户
	CurrentUser int //当前用户
	MagicIndex  int //财神索引

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

	//扑克数据
	CardCount    []int //扑克数目
	CardData     []int //扑克列表 MAX_COUNT
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
	TurnScore    []int //积分信息
	CollectScore []int //积分信息
}

//约战类型特殊属性
type G2C_Record struct {
	Count       int
	HuCount     []int8  //胡牌次数
	MaCount     []int8  //中码个数
	AnGang      []int8  //暗杠次数
	MingGang    []int8  //明杠次数
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
	GameID int //游戏 I D
	UserID int //用户 I D

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
}

type SysMsg struct {
	ClientID int
	Type     int
	Context  string
}

//////////////////////// hzmj proto begin //////////////////////////////
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

//发送扑克
type G2C_HZMG_GameStart struct {
	BankerUser   int     //当前庄家
	ReplaceUser  int     //补花用户
	SiceCount    int     //骰子点数
	HeapHead     int     //牌堆头部
	HeapTail     int     //牌堆尾部
	MagicIndex   int     //财神索引
	HeapCardInfo [][]int //堆立信息
	UserAction   int     //用户动作
	CardData     []int   //麻将列表
	OutCardCount int
	OutCardData  [][]int
}

//游戏结束
type G2C_HZMJ_GameConclude struct {
	//积分变量
	CellScore int   //单元积分
	GameScore []int //游戏积分
	Revenue   []int //税收积分
	GangScore []int //本局杠输赢分
	//结束信息
	ProvideUser  int   //供应用户
	ProvideCard  int   //供应扑克
	SendCardData int   //最后发牌
	ChiHuKind    []int //胡牌类型
	ChiHuRight   []int //胡牌类型
	LeftUser     int   //玩家逃跑
	LianZhuang   int   //连庄

	//游戏信息
	CardCount    []int   //扑克数目
	HandCardData [][]int //扑克列表

	MaCount []int //码数
	MaData  []int //码数据
}

// 出牌
type C2G_HZMJ_HZOutCard struct {
	CardData int
}

//出操作
type C2G_HZMJ_OperateCard struct {
	OperateCode int
	OperateCard []int
}

//// s to c
//用户出牌
type G2C_HZMJ_OutCard struct {
	OutCardUser int  //出牌用户
	OutCardData int  //出牌扑克
	SysOut      bool //托管系统出牌
}

type G2C_HZMJ_OperateNotify struct {
	ActionMask int //动作掩码
	ActionCard int //动作扑克
}

//发送扑克
type G2C_HZMJ_SendCard struct {
	CardData     int  //扑克数据
	ActionMask   int  //动作掩码
	CurrentUser  int  //当前用户
	SendCardUser int  //发牌用户
	ReplaceUser  int  //补牌用户
	Tail         bool //末尾发牌
}

//操作命令
type G2C_HZMJ_OperateResult struct {
	OperateUser int    //操作用户
	ActionMask  int    //动作掩码
	ProvideUser int    //供应用户
	OperateCode int    //操作代码
	OperateCard [3]int //操作扑克
}

type G2C_HZMJ_Trustee struct { //用户托管
	Trustee bool //是否托管
	ChairID int  //托管用户
}

///////////////////////// hzmj proto end ///////////////////////////////

///////////////////////// game chart begin ///////////////////////////////
type C2G_GameChart_ToAll struct {
	ChatColor	int	//字体颜色
	SendUserID	int	//发送者id
	ChatString	string	//消息内容
	ChatIndex	int	//第几条消息
}

type G2C_GameChart_ToAll struct {
	ChatColor	int
	SendUserID	int
	TargetUserID	int
	ClientID	int
	ChatIndex	int
	ChatString	string
}


///////////////////////// game chart end ///////////////////////////////