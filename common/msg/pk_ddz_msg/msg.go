package pk_ddz_msg

import (
	"mj/common/msg"
)

var (
	Processor = msg.Processor
)

func init() {
	//DDZ msg
	Processor.Register(&G2C_DDZ_StatusFree{})
	Processor.Register(&G2C_DDZ_StatusCall{})
	Processor.Register(&G2C_DDZ_StatusPlay{})
	Processor.Register(&G2C_DDZ_GameStart{})
	Processor.Register(&G2C_DDZ_CallScore{})
	Processor.Register(&G2C_DDZ_BankerInfo{})
	Processor.Register(&G2C_DDZ_OutCard{})
	Processor.Register(&G2C_DDZ_PassCard{})
	Processor.Register(&G2C_DDZ_GameConclude{})
	Processor.Register(&G2C_DDZ_TRUSTEE{})
	Processor.Register(&C2G_DDZ_CallScore{})
	Processor.Register(&C2G_DDZ_OutCard{})
	Processor.Register(&C2G_DDZ_TRUSTEE{})
}

// 游戏场景
type G2C_DDZ_GAMESTATUS struct {
	StatusData string // 场景附带消息，json格式
}

//空闲状态
type G2C_DDZ_StatusFree struct {
	// 游戏属性
	CellScore int // 基础积分

	GameType  int  // 游戏类型(0：经典场 1：欢乐场 2：癞子场)
	EightKing bool // 是否八王模式
	PlayCount int  // 游戏局数

	// 时间信息
	TimeOutCard     int // 出牌时间
	TimeCallScore   int // 叫分时间
	TimeStartGame   int // 开始时间
	TimeHeadOutCard int // 首出时间

	// 历史积分
	TurnScore    []int //积分信息
	CollectScore []int //积分信息

	ShowCardSign map[int]bool // 用户明牌标识
	TrusteeSign  []bool       // 托管标识
}

//叫分状态
type G2C_DDZ_StatusCall struct {
	// 时间信息
	TimeOutCard     int //出牌时间
	TimeCallScore   int //叫分时间
	TimeStartGame   int //开始时间
	TimeHeadOutCard int //首出时间

	// 游戏信息
	GameType      int   // 游戏类型(0：经典场 1：欢乐场 2：癞子场)
	LaiziCard     int   // 癞子牌
	EightKing     bool  // 是否八王模式
	CellScore     int   // 单元积分
	CurrentUser   int   // 当前玩家
	BankerScore   int   // 庄家叫分
	ScoreInfo     []int // 叫分信息
	HandCardCount []int //扑克数目

	// 历史积分
	TurnScore    []int // 积分信息
	CollectScore []int // 积分信息

	// 明牌
	ShowCardSign map[int]bool // 明牌标识
	ShowCardData [][]int      // 明牌数据
}

//游戏状态
type G2C_DDZ_StatusPlay struct {
	// 时间信息
	TimeOutCard     int //出牌时间
	TimeCallScore   int //叫分时间
	TimeStartGame   int //开始时间
	TimeHeadOutCard int //首出时间

	//游戏变量
	CellScore int //单元积分

	BankerUser  int  //庄家用户
	CurrentUser int  //当前玩家
	BankerScore int  //庄家叫分
	EightKing   bool // 是否八王模式
	GameType    int  // 游戏类型(0：经典场 1：欢乐场 2：癞子场)
	LaiziCard   int  // 癞子牌

	//出牌信息
	TurnWiner    int   //出牌玩家
	TurnCardData []int //出牌数据

	//扑克信息
	BankerCard    [3]int //游戏底牌
	HandCardCount []int  //扑克数目

	EachBombCount []int // 炸弹个数
	KingCount     []int // 八王个数
	//历史积分
	TurnScore    []int //积分信息
	CollectScore []int //积分信息

	// 明牌
	ShowCardSign map[int]bool // 明牌标识
	ShowCardData [][]int      // 明牌数据
}

//发送扑克
type G2C_DDZ_GameStart struct {
	CallScoreUser int          // 叫分玩家
	LiziCard      int          // 癞子牌
	ShowCard      map[int]bool // 明牌信息
	CardData      [][]int      // 扑克列表
}

//用户叫分
type G2C_DDZ_CallScore struct {
	ScoreInfo []int // 叫分信息
}

//庄家信息
type G2C_DDZ_BankerInfo struct {
	BankerUser  int    // 庄家玩家
	CurrentUser int    // 当前玩家
	BankerScore int    // 庄家叫分
	BankerCard  [3]int // 庄家扑克
}

//用户出牌
type G2C_DDZ_OutCard struct {
	CurrentUser int   //当前玩家
	OutCardUser int   //出牌玩家
	CardData    []int //扑克列表
}

//放弃出牌
type G2C_DDZ_PassCard struct {
	TurnOver     bool //一轮结束
	CurrentUser  int  //当前玩家
	PassCardUser int  //放弃玩家
}

//游戏结束
type G2C_DDZ_GameConclude struct {
	//积分变量
	CellScore int //单元积分

	//春天标志
	SpringSign int //春天标志(0：无 1：春天 2：反春天)

	//炸弹信息
	EachBombCount []int //炸弹个数

	//游戏信息
	BankerScore  int     //叫分数目
	HandCardData [][]int //扑克列表
	GameScore    []int   //游戏积分
	// 八王信息
	KingCount []int // 八王信息
}

//托管
type G2C_DDZ_TRUSTEE struct {
	TrusteeUser int  //托管玩家
	Trustee     bool //托管标志
}

// 用户明牌
type G2C_DDZ_ShowCard struct {
	ShowCardUser int   // 明牌用户
	CardData     []int // 明牌数据
}

// 叫分失败
type G2C_DDZ_CALLScoreFail struct {
	CallScoreUser int    // 当前叫分玩家
	CallScore     int    // 当前叫分
	ErrorCode     int    // 错误代码
	ErrorStr      string // 错误信息
}

//////////////////////////////////////////////////////////////////////////////////
//C->S

//用户叫分
type C2G_DDZ_CallScore struct {
	CallScore int //叫分数目
}

//用户出牌
type C2G_DDZ_OutCard struct {
	CardType int   // 牌型
	CardData []int //扑克数据
}

//托管
type C2G_DDZ_TRUSTEE struct {
	Trustee bool //托管标志
}

// 明牌
type C2G_DDZ_SHOWCARD struct {
}

// 斗地主创建房间附带信息
type C2G_DDZ_CreateRoomInfo struct {
	GameType int
	King     bool
}
