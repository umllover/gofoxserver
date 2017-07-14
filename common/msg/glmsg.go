package msg

//逻辑服和大厅服消息
type PlayerBrief struct {
	UID     int64
	Name    string
	HeadUrl string
	Icon    int
}

//房间简要信息
type RoomInfo struct {
	ServerID     int                    //第二类型
	KindID       int                    //第一类型
	RoomID       int                    //6位房号
	NodeID       int                    //在哪个节点上
	CurCnt       int                    //当前人数
	MaxPlayerCnt int                    //最多多人数
	PayCnt       int                    //可玩局数
	PayType      int                    //支付类型
	CurPayCnt    int                    //已玩局数
	CreateTime   int64                  //创建时间
	CreateUserId int64                  //房间房间的人
	Idx          int                    //服务器标识用的字段
	IsPublic     bool                   //是否公开匹配
	Players      map[int64]*PlayerBrief //玩家id
	SvrHost      string                 //哪个ip上的房间
	MachPlayer   map[int64]struct{}
}

///通知大厅房间结束
type RoomEndInfo struct {
	RoomId int //房间id
	Status int //0是没开始， 1是开始了
}

type UpdateRoomInfo struct {
	RoomId int
	OpName string
	Data   map[string]interface{}
}

//通知玩家在其他服登录
type S2S_NotifyOtherNodeLogin struct {
	Uid        int64
	ServerName string
}

//通知玩家在别的福登出
type S2S_NotifyOtherNodelogout struct {
	Uid int64
}

//获取游戏服的可玩游戏列表
type S2S_GetKindList struct {
}

//获取游戏服的所有房间
type S2S_GetRooms struct {
}

//通知大厅游戏解散
type S2S_notifyDelRoom struct {
	RoomID int //房间id
}
