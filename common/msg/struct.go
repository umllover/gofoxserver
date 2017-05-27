package msg


///// not regist,  sub obj ////////////////////
type DateTime struct {
	Year int 	`json:"wYear"`
	Month int	`json:"wMonth"`
	DayOfWeek int	`json:"wDayOfWeek"`
	Day int	`json:"wDay"`
	Hour int	`json:"wHour"`
	Minute int	`json:"wMinute"`
	Second int	`json:"wSecond"`
	Milliseconds int	`json:"wMilliseconds"`
}

//游戏房间
type TagGameServer struct {
	KindID int			`json:"wKindID"` 					//名称索引
	NodeID int			`json:"wNodeID"` 				//节点索引
	SortID int			`json:"wSortID"` 				//排序索引
	ServerID int		`json:"wServerID"` 					//房间索引
	ServerPort int		`json:"wServerPort"` 				//房间端口
	ServerType int64	`json:"wServerType"` 			//房间类型
	OnLineCount int64	`json:"dwOnLineCount"` 				//在线人数
	FullCount int		`json:"dwFullCount"` 				//满员人数

	//fdl add 手机房间显示上下限
	RestrictScore int64		`json:"lRestrictScore"` 				//限制积分
	MinTableScore int64		`json:"lMinTableScore"` 				//最低积分
	MinEnterScore int64		`json:"lMinEnterScore"` 				//最低积分
	MaxEnterScore int64		`json:"lMaxEnterScore"` 				//最高积分
	//fdl add end

	ServerAddr string			`json:"szServerAddr"` 		//房间名称
	ServerName string			`json:"szServerName"`  	//房间名称
	//私人房添加
	SurportType int64			`json:"DdwSurportType"` 			//支持类型
	TableCount int				`json:"wTableCount"` 		//桌子数目
};



