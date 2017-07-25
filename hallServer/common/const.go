package common

const (
	ListsMaxCnt = 30 //房间列表一次多去的最多条数
	PackCount   = 6  //一个多少条
)

const (
	RechargeStatusCreater = 0 //未支付
	RechargeStatusPaying  = 1 //支付完成
	RechargeStatusOk      = 2 //订单完成
)

//活动次数类型
const (
	ActivityTypeForever = 1 //永久类型
	ActivityTypeDay     = 2 //天类型
	ActivityTypeWeek    = 3 //周类型
)

//次数索引
const (
	dayShareTimes = 1
	shareTimes    = 2
)
