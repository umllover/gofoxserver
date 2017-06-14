package common

const (
	KIND_TYPE_HZMJ = 389
	KIND_TYPE_XSMJ = 390
	KIND_TYPE_TBNN = 28
)

//踢出玩家原因
const (
	UserOffline    = 0 //玩家主动底线
	NetWorkError   = 1 //网络错误
	ReadErrorMsg   = 2 //读到非法消息
	ServerShutdown = 3 //服务器关闭
)

const (
	TableFullCount = 1
)
