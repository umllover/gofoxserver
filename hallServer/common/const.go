package common

//踢出玩家原因
const (
	UserOffline    = 0 //玩家主动底线
	NetWorkError   = 1 //网络错误
	ReadErrorMsg   = 2 //读到非法消息
	ServerShutdown = 3 //服务器关闭
)

const (
	ListsMaxCnt = 30 //房间列表一次多去的最多条数
	PackCount   = 5  //一个多少条
)
