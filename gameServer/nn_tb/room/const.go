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

/*
const (
	//发牌状态

	Not_Send     = iota //无
	OutCard_Send        //出牌后发牌
	Gang_Send           //杠牌后发牌
	BuHua_Send          //补花后发牌

)
*/

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

/*
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
*/

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

