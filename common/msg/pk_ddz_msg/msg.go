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
	Processor.Register(&G2C_DDZ_AndroidCard{})
	Processor.Register(&G2C_DDZ_CheatCard{})
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

//空闲状态
type G2C_DDZ_StatusFree struct {
	// 游戏属性
	CellScore int // 基础积分

	// 时间信息
	TimeOutCard     int // 出牌时间
	TimeCallScore   int // 叫分时间
	TimeStartGame   int // 开始时间
	TimeHeadOutCard int // 首出时间

	// 历史积分
	TurnScore    []int //积分信息
	CollectScore []int //积分信息
}

//叫分状态
type G2C_DDZ_StatusCall struct {
	// 时间信息
	TimeOutCard     int //出牌时间
	TimeCallScore   int //叫分时间
	TimeStartGame   int //开始时间
	TimeHeadOutCard int //首出时间

	// 游戏信息
	CellScore    int   // 单元积分
	CurrentUser  int   // 当前玩家
	BankerScore  int   // 庄家叫分
	ScoreInfo    []int // 叫分信息
	HandCardData []int // 手上扑克

	// 历史积分
	TurnScore    []int // 积分信息
	CollectScore []int // 积分信息

	// 明牌
	ShowCardSign map[int]bool  // 明牌标识
	ShowCardData map[int][]int // 明牌数据
}

//游戏状态
type G2C_DDZ_StatusPlay struct {
	// 时间信息
	TimeOutCard     int //出牌时间
	TimeCallScore   int //叫分时间
	TimeStartGame   int //开始时间
	TimeHeadOutCard int //首出时间

	//游戏变量
	CellScore   int //单元积分
	BombCount   int //炸弹次数
	BankerUser  int //庄家用户
	CurrentUser int //当前玩家
	BankerScore int //庄家叫分

	//出牌信息
	TurnWiner     int   //出牌玩家
	TurnCardCount int   //出牌数目
	TurnCardData  []int //出牌数据

	//扑克信息
	BankerCard    [3]int //游戏底牌
	HandCardData  []int  //手上扑克
	HandCardCount []int  //扑克数目

	//历史积分
	TurnScore    []int //积分信息
	CollectScore []int //积分信息

	// 明牌
	ShowCardSign map[int]bool  // 明牌标识
	ShowCardData map[int][]int // 明牌数据
}

//发送扑克
type G2C_DDZ_GameStart struct {
	StartUser      int   //开始玩家
	CurrentUser    int   //当前玩家
	ValidCardData  int   //明牌扑克
	ValidCardIndex int   //明牌位置
	CardData       []int //扑克列表
}

//机器人扑克
type G2C_DDZ_AndroidCard struct {
	HandCard    [][]int //手上扑克
	CurrentUser int     //当前玩家
}

//作弊扑克
type G2C_DDZ_CheatCard struct {
	CardUser  []int   //作弊玩家
	UserCount int     //作弊数量
	CardData  [][]int //扑克列表
	CardCount []int   //扑克数量

}

//用户叫分
type G2C_DDZ_CallScore struct {
	CurrentUser   int //当前玩家
	CallScoreUser int //叫分玩家
	CurrentScore  int //当前叫分
	UserCallScore int //上次叫分
}

//庄家信息
type G2C_DDZ_BankerInfo struct {
	BankerUser  int    //庄家玩家
	CurrentUser int    //当前玩家
	BankerScore int    //庄家叫分
	BankerCard  [3]int //庄家扑克
}

//用户出牌
type G2C_DDZ_OutCard struct {
	CardCount   int   //出牌数目
	CurrentUser int   //当前玩家
	OutCardUser int   //出牌玩家
	CardData    []int //扑克列表
}

//放弃出牌
type G2C_DDZ_PassCard struct {
	TurnOver     int //一轮结束
	CurrentUser  int //当前玩家
	PassCardUser int //放弃玩家
}

//游戏结束
type G2C_DDZ_GameConclude struct {
	//积分变量
	CellScore int   //单元积分
	GameScore []int //游戏积分

	//春天标志
	ChunTian    bool //春天标志
	FanChunTian bool //春天标志

	//炸弹信息
	BombCount     int   //炸弹个数
	EachBombCount []int //炸弹个数

	//游戏信息
	BankerScore  int   //叫分数目
	CardCount    []int //扑克数目
	HandCardData []int //扑克列表
}

//托管
type G2C_DDZ_TRUSTEE struct {
	TrusteeUser int //托管玩家
	Trustee     int //托管标志
}

// 用户明牌
type G2C_DDZ_ShowCard struct {
	ShowCardUser int   // 明牌用户
	CardData     []int // 明牌数据
}

//////////////////////////////////////////////////////////////////////////////////
//C->S

//用户叫分
type C2G_DDZ_CallScore struct {
	CallScore int //叫分数目
}

//用户出牌
type C2G_DDZ_OutCard struct {
	CardCount int   //出牌数目
	CardData  []int //扑克数据
}

//托管
type C2G_DDZ_TRUSTEE struct {
	Trustee bool //托管标志
}
