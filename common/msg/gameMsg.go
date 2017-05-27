package msg

////// c 2 s
//手机登录
type C2G_GR_LogonMobile struct {
//版本信息
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





//// s 2 c
//登录成功
type G2C_LogonSuccess struct {

}

//登录失败
type G2C_LogonFailur struct {
	ResultCode int
	DescribeString string
}

////////////// game sverver proto
//创建房间
type C2G_CreateRoom struct {

}

// 出牌
type C2G_HZOutCard struct {

}


