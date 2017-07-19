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
	OUTCARD_OUTING   = 0XFFFF // 出牌中
	OUTCARD_PASS     = 0XFFFE // 不出
	OUTCARD_MAXCOUNT = 10000  // 最大出牌次数

	// 游戏类型
	GAME_TYPE_INVALID = 255 // 无效类型
	GAME_TYPE_CLASSIC = 0   // 经典场
	GAME_TYPE_HAPPY   = 1   // 欢乐场
	GAME_TYPE_LZ      = 2   // 癞子场
	GAME_TYPE_MONEY   = 3   // 金币场
)

const (
	// 牌类型
	CT_ERROR           = 0     // 错误类型
	CT_SINGLE          = 0x100 // 单张牌（散牌）(结尾两位16进制代表牌的逻辑数值，直接可以拿来比大小)
	CT_DOUBLE          = 0x200 // 对子牌(结尾两位16进制代表牌的逻辑数值，直接可以拿来比大小)
	CT_THREE           = 0x300 // 三张牌(结尾两位16进制代表牌的逻辑数值，直接可以拿来比大小)
	CT_THREE_TAKE_ONE  = 0x400 // 三带一(结尾两位16进制代表三根主牌的逻辑数值，直接可以拿来比大小)
	CT_THREE_TAKE_TWO  = 0x500 // 三带二(结尾两位16进制代表三根主牌的逻辑数值，直接可以拿来比大小)
	CT_SINGLE_LINE     = 0x600 // 单顺子（第二位16进制代表顺子的张数，第一位16进制代表最大牌的逻辑值）
	CT_DOUBLE_LINE     = 0x700 // 双顺子（第二位16进制代表顺子的对子数，第一位16进制代表最大牌的逻辑值）
	CT_THREE_LINE      = 0x800 // 三顺子（第二位16进制代表顺子的对子数，第一位16进制代表最大牌的逻辑值）
	CT_THREE_LINE_TAKE = 0X900 // 飞机带翅膀(第二位16进制代表最大飞机主牌逻辑数值，第二位16进制中，第一位代表翅膀是对子还是单根，其它三位代表多少对)
	CT_FOUR_TAKE_TWO   = 0xA00 // 四带二(第一位16进制代表带的是单还是对子，第二位代表四根主牌的逻辑值)
	CT_BOMB_CARD       = 0xB00 // 炸弹类型(第一位16进制代表癞子数量，只要不为0，则任何为0的炸弹都比它大，第二位代表炸弹代表的逻辑值)
	CT_KING            = 0xC00 // 火箭(第一位16进制代表小王数量，第二位代表大王数量，两位相加多的大，一样多时，大王多的大)
)
