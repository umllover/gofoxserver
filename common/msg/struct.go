package msg

///// not regist,  sub obj ////////////////////
type DateTime struct {
	Year         int `json:"wYear"`
	Month        int `json:"wMonth"`
	DayOfWeek    int `json:"wDayOfWeek"`
	Day          int `json:"wDay"`
	Hour         int `json:"wHour"`
	Minute       int `json:"wMinute"`
	Second       int `json:"wSecond"`
	Milliseconds int `json:"wMilliseconds"`
}

//游戏房间
type TagGameServer struct {
	KindID      int   `json:"wKindID"`       //名称索引
	NodeID      int   `json:"wNodeID"`       //节点索引
	SortID      int   `json:"wSortID"`       //排序索引
	ServerID    int   `json:"wServerID"`     //房间索引
	ServerPort  int   `json:"wServerPort"`   //房间端口
	ServerType  int64 `json:"wServerType"`   //房间类型
	OnLineCount int64 `json:"dwOnLineCount"` //在线人数
	FullCount   int   `json:"dwFullCount"`   //满员人数

	//fdl add 手机房间显示上下限
	RestrictScore int64 `json:"lRestrictScore"` //限制积分
	MinTableScore int64 `json:"lMinTableScore"` //最低积分
	MinEnterScore int64 `json:"lMinEnterScore"` //最低积分
	MaxEnterScore int64 `json:"lMaxEnterScore"` //最高积分
	//fdl add end

	ServerAddr string `json:"szServerAddr"` //ip地址
	ServerName string `json:"szServerName"` //游戏名称
	//私人房添加
	SurportType int64 `json:"DdwSurportType"` //支持类型
	TableCount  int   `json:"wTableCount"`    //桌子数目
}

//用户状态信息
type UserStu struct {
	TableID    int // 桌子索引
	ChairID    int // 椅子位置
	UserStatus int // 用户状态
}

//组合子项
type WeaveItem struct {
	WeaveKind   int   //组合类型
	CenterCard  int   //中心扑克
	Param       int   //类型标志
	ProvideUser int   //供应用户
	CardData    []int //麻将数据
	PublicCard  bool  //公开标志
}

//积分信息
type TagScoreInfo struct {
	Type     int //积分类型
	Score    int //用户分数
	Grade    int //用户成绩
	Revenue  int //游戏税收
	SubMonry int //扣除代币
	IsWin    int
}

//自己创建的房间的信息
type CreatorRoomInfo struct {
	CreatorTime int64  //创建时间
	RoomID      int    //房号
	RoomName    string //房间名称
	Status      int    //状态
}
