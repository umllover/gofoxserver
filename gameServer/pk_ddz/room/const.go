package room

const (
	// 游戏状态
	GAME_STATUS_FREE = 0
	GAME_STATUS_CALL = 1
	GAME_STATUS_PLAY = 2

	// 用户叫分信息
	CALLSCORE_CALLING = 0XFFFF // 正在叫分状态
	CALLSCORE_NOCALL  = 0xFFFE // 未叫状态
	CALLSCORE_MAX     = 3      // 允许叫分的最大值

	// 用户出牌状态
	OUTCARD_OUTING = 0XFFFF // 出牌中
	OUTCARD_PASS   = 0XFFFE // 不出

	// 游戏类型
	GAME_TYPE_INVALID = 255 // 无效类型
	GAME_TYPE_CLASSIC = 0   // 经典场
	GAME_TYPE_HAPPY   = 1   // 欢乐场
	GAME_TYPE_LZ      = 2   // 癞子场
	GAME_TYPE_MONEY   = 3   // 金币场
)
