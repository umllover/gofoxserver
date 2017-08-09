package pk_common_msg

import (
	"mj/common/msg"
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
	PlayMode   int //游戏模式
	CountLimit int //局数限制

	PlayerCount      int     // 游戏人数
	CurrentPlayCount int     //房间已玩局数
	EachRoundScore   [][]int //房间每局游戏比分
	InitScore        []int   //积分信息
}

type G2C_PKCOMMON_StatusCall struct {
	CallBanker  int   //叫庄用户
	DynamicJoin int   //动态加入
	PlayStatus  []int //用户状态

	//历史积分
	TurnScore    []int64 //积分信息
	CollectScore []int64 //积分信息
	GameRoomName string  //房间名称
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
	//OpenCardData []
}
type G2C_PKCOMMON_StatusPlay struct {
	CellScore int //基础积分

	UserReLoginInfos []*UserReLoginInfo
	GameStatus       int //游戏状态
	PlayerCount      int //玩家人数
	BankerUser       int //庄家用户
	PublicCardData   []int
	HandCardData     [][]int //桌面扑克

	InitScore    []int  //积分信息
	GameRoomName string //房间名称

	CurrentPlayCount int //房间已玩局数
	LimitPlayCount   int //总局数
}
