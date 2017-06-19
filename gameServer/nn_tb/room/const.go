package room

const (
	//用户状态
	US_NULL    = 0x00 //没有状态
	US_FREE    = 0x01 //站立状态
	US_SIT     = 0x02 //坐下状态
	US_READY   = 0x03 //同意状态
	US_LOOKON  = 0x04 //旁观状态
	US_PLAYING = 0x05 //游戏状态
	US_OFFLINE = 0x06 //断线状态
)

const (
	//房间状态
	RoomStatusReady    = 0
	RoomStatusStarting = 1
	RoomStatusEnd      = 2
)

const (
	//发牌状态

	Not_Send     = iota //无
	OutCard_Send        //出牌后发牌
	Gang_Send           //杠牌后发牌
	BuHua_Send          //补花后发牌

)

const (
	//结束原因
	GER_NORMAL        = 0x00 //常规结束
	GER_DISMISS       = 0x01 //游戏解散
	GER_USER_LEAVE    = 0x02 //用户离开
	GER_NETWORK_ERROR = 0x03 //网络错误
)

const (
	//分数模式
	SCORE_GENRE_NORMAL   = 0x0100 //普通模式
	SCORE_GENRE_POSITIVE = 0x0200 //非负模式
)

const (
	//积分类型
	SCORE_TYPE_NULL    = 0x00 //无效积分
	SCORE_TYPE_WIN     = 0x01 //胜局积分
	SCORE_TYPE_LOSE    = 0x02 //输局积分
	SCORE_TYPE_DRAW    = 0x03 //和局积分
	SCORE_TYPE_FLEE    = 0x04 //逃局积分
	SCORE_TYPE_PRESENT = 0x10 //赠送积分
	SCORE_TYPE_SERVICE = 0x11 //服务积分
)

const (
	//税收定义
	REVENUE_BENCHMARK   = 0    //税收起点
	REVENUE_DENOMINATOR = 1000 //税收分母
	PERSONAL_ROOM_CHAIR = 8    //私人房间座子上椅子的最大数目
)


type HistoryScore struct {
	TurnScore    int
	CollectScore int
}

const (
	GAME_PLAYER 		= 4		//房间人数
	MAX_COUNT			= 5		//牌数
)

//扑克数据
var  CardDataArray  = []int {
	0x01,0x02,0x03,0x04,0x05,0x06,0x07,0x08,0x09,0x0A,0x0B,0x0C,0x0D,	//方块 A - K
	0x11,0x12,0x13,0x14,0x15,0x16,0x17,0x18,0x19,0x1A,0x1B,0x1C,0x1D,	//梅花 A - K
	0x21,0x22,0x23,0x24,0x25,0x26,0x27,0x28,0x29,0x2A,0x2B,0x2C,0x2D,	//红桃 A - K
	0x31,0x32,0x33,0x34,0x35,0x36,0x37,0x38,0x39,0x3A,0x3B,0x3C,0x3D,	//黑桃 A - K 3 13 54
	0x4E,0x4F,//14 15
}

//花色
const (
	LOGIC_MASK_COLOR = 0xF0
	LOGIC_MASK_VALUE = 0x0F
)

// 扑克类型
const (

)
