package mj_zp_msg

import (
	"mj/common/msg"
)

func init() {
	msg.Processor.Register(&C2G_MJZP_SetChaHua{})
	msg.Processor.Register(&C2G_MJZP_ReplaceCard{})
	msg.Processor.Register(&C2G_MJZP_ListenCard{})
	msg.Processor.Register(&C2G_ZPMJ_OutCard{})
	msg.Processor.Register(&C2G_ZPMJ_OperateCard{})
	msg.Processor.Register(&C2G_MJZP_Trustee{})

	msg.Processor.Register(&G2C_MJZP_NotifiChaHua{})
	msg.Processor.Register(&G2C_MJZP_ReplaceCard{})
	msg.Processor.Register(&G2C_MJZP_ListenCard{})
	msg.Processor.Register(&G2C_ZPMJ_GameConclude{})
	msg.Processor.Register(&G2C_ZPMG_GameStart{})
	msg.Processor.Register(&G2C_ZPMJ_Trustee{})
	msg.Processor.Register(&G2C_ZPMJ_OutCard{})
	msg.Processor.Register(&G2C_ZPMJ_OperateResult{})
	msg.Processor.Register(&G2C_ZPMJ_SendCard{})
	msg.Processor.Register(&G2C_ZPMJ_StatusPlay{})
	msg.Processor.Register(&G2C_MJZP_UserCharHua{})
	msg.Processor.Register(&G2C_MJZP_OperateNotify{})
	msg.Processor.Register(&G2C_ZPMJ_HuData{})
}

type G2C_MJZP_OperateNotify struct {
	ActionMask int //动作掩码
	ActionCard int //动作扑克
}

//通知插花
type G2C_MJZP_NotifiChaHua struct {
}

//设置插花
type C2G_MJZP_SetChaHua struct {
	SetCount int //设置插花数量
}

//用户插花通知
type G2C_MJZP_UserCharHua struct {
	Chair    int //椅子
	SetCount int //设置插花数量
}

//补花
type G2C_MJZP_ReplaceCard struct {
	ReplaceUser  int  //补牌用户
	ReplaceCard  int  //补牌扑克
	NewCard      int  //补完扑克
	IsInitFlower bool //是否开局补花，true开局补花
}

//补花
type C2G_MJZP_ReplaceCard struct {
	CardData int //扑克数据
}

//听牌
type C2G_MJZP_ListenCard struct {
	ListenCard bool //是否听牌
}

//听牌
type G2C_MJZP_ListenCard struct {
	ListenUser  int   //听牌用户
	IsListen    bool  //是否听牌
	HuCardCount int   //胡几张牌
	HuCardData  []int //胡牌数据
}

// 出牌
type C2G_ZPMJ_OutCard struct {
	CardData int
}

//出操作
type C2G_ZPMJ_OperateCard struct {
	OperateCode int
	OperateCard []int
}

//游戏结束
type G2C_ZPMJ_GameConclude struct {
	//积分变量
	CellScore int   //单元积分
	GameScore []int //游戏积分
	Revenue   []int //税收积分
	GangScore []int //本局杠输赢分
	//结束信息
	ProvideUser  int          //供应用户
	ProvideCard  int          //供应扑克
	SendCardData int          //最后发牌
	ChiHuKind    []int        //胡牌类型
	ChiHuRight   []int        //胡牌类型
	LeftUser     int          //玩家逃跑
	LianZhuang   int          //连庄
	ScoreKind    [4][35]int   //得分类型
	ZhuaHua      [16]*HuaUser //用户抓花
	//type HuaUser struct {
	//	chairID int
	//	card    int
	//}

	//游戏信息
	CardCount    []int   //扑克数目
	HandCardData [][]int //扑克列表

	MaCount []int //码数
	MaData  []int //码数据
}

//发送扑克
type G2C_ZPMG_GameStart struct {
	BankerUser   int     //当前庄家
	ReplaceUser  int     //补花用户
	SiceCount    int     //骰子点数
	HeapHead     int     //牌堆头部
	HeapTail     int     //牌堆尾部
	HeapCardInfo [][]int //堆立信息
	UserAction   int     //用户动作
	CardData     []int   //麻将列表
}

//托管
type C2G_MJZP_Trustee struct {
	Trustee bool //1：托管 0：取消托管
}

//用户出牌
type G2C_ZPMJ_OutCard struct {
	OutCardUser int  //出牌用户
	OutCardData int  //出牌扑克
	SysOut      bool //托管系统出牌
}

//操作命令
type G2C_ZPMJ_OperateResult struct {
	OperateUser int    //操作用户
	ActionMask  int    //动作掩码
	ProvideUser int    //供应用户
	OperateCode int    //操作代码
	OperateCard [3]int //操作扑克
}

//发送扑克
type G2C_ZPMJ_SendCard struct {
	CardData     int  //扑克数据
	ActionMask   int  //动作掩码
	CurrentUser  int  //当前用户
	SendCardUser int  //发牌用户
	ReplaceUser  int  //补牌用户
	Tail         bool //末尾发牌
}

type G2C_ZPMJ_Trustee struct { //用户托管
	Trustee bool //是否托管
	ChairID int  //托管用户
}

//游戏状态 游戏已经开始了发送的结构
type G2C_ZPMJ_StatusPlay struct {
	//时间信息
	TimeOutCard     int   //出牌时间
	TimeOperateCard int   //叫分时间
	CreateTime      int64 //开始时间

	//游戏变量
	CellScore   int   //单元积分
	BankerUser  int   //庄家用户
	CurrentUser int   //当前用户
	MagicIndex  int   //财神索引
	ChaHuaCnt   []int //插花数
	BuHuaCnt    []int //补花数
	ZhuaHuaCnt  int   //抓花数

	//规则
	PlayerCount int //玩家人数
	MaCount     int //码数

	//状态变量
	ActionCard    int    //动作扑克
	ActionMask    int    //动作掩码
	LeftCardCount int    //剩余数目
	Trustee       []bool //是否托管 index 就是椅子id
	Ting          []bool //是否听牌  index chairId

	//出牌信息
	OutCardUser  int     //出牌用户
	OutCardData  int     //出牌扑克
	DiscardCount []int   //丢弃数目
	DiscardCard  [][]int //丢弃记录

	//扑克数据
	CardCount    []int //扑克数目
	CardData     []int //扑克列表 room.GetCfg().MaxCount
	SendCardData int   //发送扑克

	//组合扑克
	WeaveItemCount []int              //组合数目
	WeaveItemArray [][]*msg.WeaveItem //组合扑克 [GAME_PLAYER][MAX_WEAVE]

	//堆立信息
	HeapHead     int     //堆立头部
	HeapTail     int     //堆立尾部
	HeapCardInfo [][]int //堆牌信息

	HuCardCount   []int
	HuCardData    [][]int
	OutCardCount  int
	OutCardDataEx []int
	//历史积分
	TurnScore    []int //积分信息
	CollectScore []int //积分信息
}

//听牌
type G2C_ZPMJ_HuData struct {
	//出哪几张能听
	OutCardCount int
	OutCardData  []int
	//听后能胡哪几张牌
	HuCardCount []int
	HuCardData  [][]int
	//胡牌剩余数
	HuCardRemainingCount [][]int
}

//抓花结构体子项
type HuaUser struct {
	ChairID int  //椅子号
	Card    int  //牌值
	IsZhong bool //是否中花
}
