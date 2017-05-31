package msg

////// c 2 s
//手机登录
type C2G_GR_LogonMobile struct {
	GameID int							//游戏标识
	ProcessVersion int				//进程版本

	//桌子区域
	DeviceType int                       //设备类型
	BehaviorFlags int                     //行为标识
	PageTableCount int                    //分页桌数

	//登录信息
	UserID int					//用户 I D
	Password string;				//登录密码
	MachineID string				//机器标识
};

// 请求更换椅子消息
type C2G_GR_UserChairReq struct {

}

//请求创建房间消息
type C2G_CreateTable struct {
	CellScore int							    //底分设置
	DrawCountLimit int					//局数限制
	DrawTimeLimit int					//时间限制
	JoinGamePeopleCount int			//参与游戏的人数
	RoomTax int								//单独一个私人房间的税率，千分比
	Password string	//密码设置
    GameRule []int8				//游戏规则 弟 0 位标识 是否设置规则 0 代表未设置 1 代表设置
	Kind int 				//游戏类型
	ServerId int 			//子类型
}

//查询房间信息
type C2G_SearchServerTable struct {
	ServerID int
	KindID int
}


// 出牌
type C2G_HZOutCard struct {

}




//// s 2 c
//登录成功
type G2C_LogonSuccess struct {

}

//登录失败
type G2C_LogonFailur struct {
	ResultCode int
	DescribeString string
}

// 创建房间失败消息
type G2C_CreateTableFailure struct {
	ErrorCode int
	DescribeString string
}

//创建房间成功的消息
type G2C_CreateTableSucess struct {
	TableID int					//房间编号
	DrawCountLimit int				//局数限制
	DrawTimeLimit int				//时间限制
	Beans int						//游戏豆
	RoomCard int					//房卡数量
}

//查询房间的结果
type G2C_SearchResult struct {
	ServerID int							//房间 I D
	TableID int								//桌子 I D
}


