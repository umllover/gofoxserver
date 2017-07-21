package msg

/////////////////////// hall server msg ////////////////////////////

////////////////////// c 2 l ////////////////////
//登录消息
type C2L_Login struct {
	ModuleID     int
	PlazaVersion int
	DeviceType   int
	LogonPass    string
	Accounts     string
	MachineID    string
	MobilePhone  string
	SessionKey   string
}

//注册消息
type C2L_Regist struct {
	ModuleID     int //模块标识
	PlazaVersion int //广场版本
	DeviceType   int //设备类型
	//密码变量
	LogonPass  string //登录密码
	InsurePass string //银行密码

	//注册信息
	FaceID   int8   //头像标识
	Gender   int8   //用户性别
	Accounts string //登录帐号
	NickName string //用户昵称

	//连接信息
	MachineID   string //机器标识
	MobilePhone string //电话号码  //默认不获取本机号码
}

//注册消息
type L2C_RegistResult struct {
	Code int //非0位失败
}

//请求创建房间消息
type C2L_CreateTable struct {
	DrawCountLimit int                    //局数限制
	Password       string                 //密码设置
	Kind           int                    //游戏类型
	ServerId       int                    //子类型
	PayType        int                    //1是自己付钱， 2是AA
	Public         bool                   //是否公开
	RoomName       string                 //房间名字
	OtherInfo      map[string]interface{} //其他配置， 对应 key v 结构 客户端 {k1:v1,k2:v2}即可
}

//删除自己创建的房间
type C2L_DeleteRoom struct {
	RoomId int //房间id
}

//删除房间的结果
type L2C_DeleteRoomResult struct {
	Code int //非0 为删除失败的错误码
}

//查询房间信息
type C2L_SearchServerTable struct {
	TableID int
	KindID  int
}

//查询房间的结果
type L2C_SearchResult struct {
	TableID  int //桌子 I D 返回0 是没匹配到
	ServerIP string
}

//获取玩家显示信息
type C2L_User_Individual struct {
	UserId int
}

//请求房间列表
type C2L_GetRoomList struct {
	KindID int //要查看哪个游戏类型
	Num    int //获取第几页
}

//请求匹配一个房间
type C2L_QuickMatch struct {
	KindID int //要匹配的游戏类型1
}

//服务器接收到了请求匹配的结果
type L2C_QuickMatchOk struct {
	MatchTime int //多少秒没收到结果后退出匹配, 间隔时间
}

//请求查看开放记录
type C2L_ReqCreatorRoomRecord struct {
}

//请求查看房间玩家简要信息
type C2L_ReqRoomPlayerBrief struct {
	RoomId int
}

//领取奖励
type C2L_DrawSahreAward struct {
	DrawId int //领取奖励的key
}

//设置推荐人
type C2L_SetElect struct {
	ElectUid int64 //推荐人id
}

//设置推荐人结果
type L2C_SetElectResult struct {
	RetCode int // 0带表成功， 其他则是错误码
}

//请求获取验证号码
type C2L_ReqBindMaskCode struct {
}

//请求获取验证号码返回
type L2C_ReqBindMaskCodeRsp struct {
	Code int //非0位失败
}

//绑定手机号码
type C2L_SetPhoneNumber struct {
	PhoneNumber string
	MaskCode    int
}

//绑定手机结果
type L2C_SetPhoneNumberRsp struct {
	Code int //非0 位设置失败
}

//点赞
type C2L_DianZhan struct {
	UserID int //用户id
}

//点赞结果
type C2L_DianZhanRsp struct {
	Star int //当前赞数
}

//续费 再来一局
type C2L_RenewalFees struct {
}

//续费 再来一局 结果
type L2C_RenewalFeesRsp struct {
	Code int //非0位失败
}

//修改名字
type C2L_ChangeUserName struct {
	NewName string
}

//修改名字结果
type L2C_ChangeUserNameRsp struct {
	Code    int //非0w位失败
	NewName string
}

//修改个性签名
type C2L_ChangeSign struct {
	Sign string
}

//修改个性签名结果
type L2C_ChangeSignRsp struct {
	Code    int //非0w位失败
	NewSign string
}

/////////// l 2 c /////////////////////////
//登录失败
type L2C_LogonFailure struct {
	ResultCode     int    //错误代码
	DescribeString string //描述消息
}

//登录成功
type L2C_LogonSuccess struct {
	FaceID     int8   `json:"wFaceID"`      //头像标识
	Gender     int8   `json:"cbGender"`     //用户性别
	UserID     int64  `json:"dwUserID"`     //用户 I D
	Spreader   int    `json:"szSpreader"`   //推荐人用户标识
	GameID     int    `json:"dwGameID"`     //游戏 I D
	Experience int    `json:"dwExperience"` //经验数值
	LoveLiness int    `json:"dwLoveLiness"` //用户魅力
	NickName   string `json:"szNickName"`   //用户昵称

	//用户成绩
	UserScore    int64     `json:"lUserScore"`   //用户欢乐豆
	UserInsure   int64     `json:"lUserInsure"`  //用户银行
	Medal        int       `json:"dwMedal"`      //用户钻石
	UnderWrite   string    `json:"szUnderWrite"` //个性签名
	WinCount     int       `json:"dwWinCount"`   //赢局数
	LostCount    int       `json:"dwLostCount"`  //输局数
	DrawCount    int       `json:"dwDrawCount"`  //和局数
	FleeCount    int       `json:"dwFleeCount"`  //跑局数
	RegisterDate *DateTime `json:"RegisterDate"` //注册时间

	//额外信息
	MbTicket        int `json:"dwMbTicket"`        //手机兑换券数量
	MbPayTotal      int `json:"dwMbPayTotal"`      //手机充值累计
	MbVipLevel      int `json:"dwMbVipLevel"`      //手机VIP等级
	PayMbVipUpgrade int `json:"dwPayMbVipUpgrade"` //手机VIP升级，所需充值数（vip最高级时该值为0）

	//约战房相关
	RoomCard     int    `json:"lRoomCard"`      //用户房卡
	LockServerID int    `json:"dwLockServerID"` //锁定房间
	KindID       int    `json:"dwKindID"`       //游戏类型
	HallNodeID   int    `json:"HallNodeID"`
	ServerIP     string `json:"ServerIP"`
}

//房间列表
type L2C_ServerList []*TagGameServer

//房间列表发送成功
type L2C_ServerListFinish struct{}

//个人资料
type L2C_UserIndividual struct {
	//用户信息
	UserID      int64  //用户 I D
	NickName    string //昵称
	Accounts    string //账号
	WinCount    int    //赢数
	LostCount   int    //输数
	DrawCount   int    //平数
	Medal       int
	RoomCard    int  //房卡
	MemberOrder int8 //会员等级
	Score       int64
	HeadImgUrl  string
	PhomeNumber string //电话号码
	Sign        string //个性签名
	Star        int    //赞数
}

//返回房间列表
type L2C_GetRoomList struct {
	Lists []*RoomInfo //房间信息
	Count int         //有多少条C2L_ReqCreatorRoomRecord
}

// 创建房间失败消息
type L2C_CreateTableFailure struct {
	ErrorCode      int
	DescribeString string
}

//创建房间成功的消息
type L2C_CreateTableSucess struct {
	TableID        int    //房间编号
	DrawCountLimit int    //局数限制
	DrawTimeLimit  int    //时间限制
	Beans          int    //游戏豆
	RoomCard       int    //房卡数量
	ServerIP       string //去哪个逻辑服玩
}

//返回开房记录
type L2C_CreatorRoomRecord struct {
	Records []*CreatorRoomInfo //创建的房间的简要信息
}

//请求房间内玩家简要信息的返回值
type L2C_RoomPlayerBrief struct {
	Players []*PlayerBrief //房间内玩家的简要信息
}

//登录时下发已领取过的奖励信息
type L2C_ActivityInfo struct {
	DayTimes  map[int]int64 //每日次数信息
	Times     map[int]int64 //永久次数信息
	WeekTimes map[int]int64 //周次数信息
}

//领取奖励结果
type L2C_DrawSahreAwardResult struct {
	DrawId  int //领取奖励的key
	Times   int //已经领取的次数
	RetCode int //领取成功还是失败的结果
}
