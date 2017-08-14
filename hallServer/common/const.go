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

const (
	ActivityShareDay    = 1 //每日分享
	ActivityRechangeDay = 2 //每日充值
	ActivityBindPhome   = 3 //绑定手机
	ActivitySetSetElect = 4 //设置推荐人
)

//次数索引
const (
	dayShareTimes = 1
	shareTimes    = 2
)
