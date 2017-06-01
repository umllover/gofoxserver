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

//请求坐下
type C2G_UserSitdown struct {
	TableI int 	// 桌子号码
	ChairID int // 椅子号码
	Password string //房间密码
}

//请求玩家信息
type C2G_REQUserInfo struct {

}

//配置信息
type C2G_GameOption struct {
	AllowLookon int						//旁观标志
	FrameVersion int					//框架版本
	ClientVersion int					//游戏版本
}

// 出牌
type C2G_HZOutCard struct {

}




//// s 2 c ////////////////////////////
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

//玩家状态
type G2C_UserStatus struct {
	UserID int
	UserStatus	*UserStu
}

//发送提示信息
type G2C_PersonalTableTip struct {
	TableOwnerUserID int			//桌主 I D
	DrawCountLimit int				//局数限制
	DrawTimeLimit int				//时间限制
	PlayCount int					//已玩局数
	PlayTime int					//已玩时间
	CellScore int					//游戏底分
	IniScore int					//初始分数
	ServerID string					//房间编号
	IsJoinGame int					//是否参与游戏
	IsGoldOrGameScore int			//金币场还是积分场 0 标识 金币场 1 标识 积分场
}

//游戏属性 ， 游戏未开始发送的结构
type G2C_StatusFree struct {
	CellScore int					//基础积分
	TimeOutCard int8					//出牌时间
	TimeOperateCard int8				//操作时间
	TimeStartGame int8				//开始时间
	TurnScore []int					//积分信息
	CollectScore []int				//积分信息
	PlayerCount int					//玩家人数
	MaCount int8						//码数
	CountLimit int               	//局数限制
}

//游戏状态 游戏已经开始了发送的结构
type G2C_StatusPlay struct {
	//时间信息
	TimeOutCard int8							//出牌时间
	TimeOperateCard int8						//叫分时间
	TimeStartGame int8							//开始时间

	//游戏变量
	CellScore int								//单元积分
	BankerUser int								//庄家用户
	CurrentUser int								//当前用户
	MagicIndex int8								//财神索引

	//规则
	PlayerCount int8				//玩家人数
	MaCount int8					//码数

	//状态变量
	ActionCard int8								//动作扑克
	ActionMask int8								//动作掩码
	LeftCardCount int8							//剩余数目
	Trustee []bool								//是否托管 index 就是椅子id
	Ting []bool								//是否听牌  index chairId

	//出牌信息
	OutCardUser int									//出牌用户
	OutCardData int8								//出牌扑克
	DiscardCount[]int8								//丢弃数目
	DiscardCard[][]int8				//丢弃记录

	//扑克数据
	CardCount []int8					//扑克数目
	CardData []int8						//扑克列表 MAX_COUNT
	SendCardData int8								//发送扑克

	//组合扑克
	WeaveItemCount	[]int8				//组合数目
	WeaveItemArray	[][]*WeaveItem		//组合扑克 [GAME_PLAYER][MAX_WEAVE]

	//堆立信息
	HeapHead int									//堆立头部
	HeapTail int									//堆立尾部
	HeapCardInfo [][]int8;						//堆牌信息

	HuCardCount	[]int8
	HuCardData	[][]int8
	OutCardCount int8
	OutCardDataEx []int8
	//历史积分
	TurnScore []int						//积分信息
	CollectScore []int					//积分信息
};

//约战类型特殊属性
type G2C_Record struct {
	Count int
	HuCount []int8//胡牌次数
	MaCount []int8 //中码个数
	AnGang []int8 //暗杠次数
	MingGang []int8 //明杠次数
	AllScore []int8	//总结算分
	DetailScore [][]int;	//单局结算分
}


