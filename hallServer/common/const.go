package common

//踢出玩家原因
const (
	UserOffline    = 0 //玩家主动底线
	NetWorkError   = 1 //网络错误
	ReadErrorMsg   = 2 //读到非法消息
	ServerShutdown = 3 //服务器关闭
)
