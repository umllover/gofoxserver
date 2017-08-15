package pk_common_msg

import (
	"mj/common/msg"
	"mj/common/msg/pk_ddz_msg"
)

func init() {
	// ----- G2C --------
	msg.Processor.Register(&G2C_PKCOMMON_StatusFree{})
	msg.Processor.Register(&G2C_PKCOMMON_StatusCall{})
	msg.Processor.Register(&G2C_PKCOMMON_StatusScore{})
	msg.Processor.Register(&G2C_PKCOMMON_StatusPlay{})
}

//---------- 游戏状态-------
type G2C_PKCOMMON_StatusFree struct {
	CellScore int //基础积分

	TurnScore    []int  //积分信息
	CollectScore []int  //积分信息
	GameRoomName string //房间名称

	TimeOutCard     int   //出牌时间
	TimeOperateCard int   //操作时间
	TimeStartGame   int64 //开始时间

	TimesCount int //倍数
	PlayMode   int //游戏模式 ddz{0：经典场 1：欢乐场 2：癞子场}
	CountLimit int //局数限制

	PlayerCount      int     // 游戏人数
	CurrentPlayCount int     //房间已玩局数
	EachRoundScore   [][]int //房间每局游戏比分
	InitScore        []int   //积分信息

	EightKing    bool   // 是否八王模式
	ShowCardSign []bool // 用户明牌标识

}

type G2C_PKCOMMON_StatusCall struct {
	CallBanker  int   //叫庄用户
	DynamicJoin int   //动态加入
	PlayStatus  []int //用户状态

	CellScore int //基础积分
	//历史积分
	TurnScore    []int64 //积分信息
	CollectScore []int64 //积分信息
	GameRoomName string  //房间名称

	TimeOutCard   int   // 出牌时间
	TimeCallScore int   // 叫分时间
	PlayMode      int   // 游戏模式 ddz{0：经典场 1：欢乐场 2：癞子场}
	WildCard      int   // 万能牌
	EightKing     bool  // 是否八王模式
	BankerScore   int   // 庄家叫分
	ScoreInfo     []int // 叫分信息
	HandCardCount []int // 扑克数目

	// 明牌
	ShowCardSign []bool  // 明牌标识
	ShowCardData [][]int // 明牌数据

	// 托管状态
	TrusteeSign []bool // 托管标识
}

type G2C_PKCOMMON_StatusScore struct {
	//下注信息
	PlayStatusi  []int   //用户状态
	DynamicJoin  int     //动态加入
	TurnMaxScore int64   //最大下注
	TableScore   []int64 //下注数目
	BankerUser   int     //庄家用户
	TurnScore    []int64 //积分信息
	CollectScore []int64 //积分信息
	GameRoomName string  //房间名称
}

type UserReLoginInfo struct {
	ChairID        int
	UserGameStatus int
	CallScoreTimes int
	AddScoreTimes  int
	OpenCardData []int
}
type G2C_PKCOMMON_StatusPlay struct {
	CellScore int //基础积分

	UserReLoginInfos []*UserReLoginInfo
	GameStatus       int     // 游戏状态
	PlayerCount      int     // 玩家人数
	BankerUser       int     // 庄家用户
	PublicCardData   []int   // 公共牌(ddz的底牌)
	HandCardData     [][]int // 桌面扑克

	InitScore    []int  //积分信息
	GameRoomName string //房间名称

	CurrentPlayCount int  // 房间已玩局数
	LimitPlayCount   int  // 总局数
	TimeOutCard      int  // 出牌时间
	CurrentUser      int  // 当前玩家
	EightKing        bool // 是否八王模式
	PlayMode         int  // 游戏类型(0：经典场 1：欢乐场 2：癞子场)
	WildCard         int  // 万能牌
	BankerScore      int  // 庄家叫分
	//出牌信息
	TurnUser     int                        // 出牌玩家
	TurnCardData pk_ddz_msg.C2G_DDZ_OutCard // 出牌数据

	//扑克信息
	HandCardCount []int // 扑克数目

	EachBombCount []int // 炸弹个数
	KingCount     []int // 八王个数

	// 明牌
	ShowCardSign []bool // 明牌标识

	// 托管
	TrusteeSign []bool // 托管标识
}
