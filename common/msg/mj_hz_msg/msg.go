package mj_hz_msg

import (
	"mj/common/msg"
)

var (
	Processor = msg.Processor
)

func init() {
	//HZMJ msg
	Processor.Register(&G2C_HZMG_GameStart{})
	Processor.Register(&C2G_HZMJ_HZOutCard{})
	Processor.Register(&G2C_HZMJ_OutCard{})
	Processor.Register(&G2C_HZMJ_OperateNotify{})
	Processor.Register(&G2C_HZMJ_SendCard{})
	Processor.Register(&C2G_HZMJ_OperateCard{})
	Processor.Register(&G2C_HZMJ_OperateResult{})
	Processor.Register(&G2C_HZMJ_Trustee{})
}

//发送扑克
type G2C_HZMG_GameStart struct {
	BankerUser   int     //当前庄家
	ReplaceUser  int     //补花用户
	SiceCount    int     //骰子点数
	HeapHead     int     //牌堆头部
	HeapTail     int     //牌堆尾部
	MagicIndex   int     //财神索引
	HeapCardInfo [][]int //堆立信息
	UserAction   int     //用户动作
	CardData     []int   //麻将列表
}

//游戏结束
type G2C_GameConclude struct {
	//积分变量
	CellScore int   //单元积分
	GameScore []int //游戏积分
	Revenue   []int //税收积分
	GangScore []int //本局杠输赢分
	//结束信息
	ProvideUser  int   //供应用户
	ProvideCard  int   //供应扑克
	SendCardData int   //最后发牌
	ChiHuKind    []int //胡牌类型
	ChiHuRight   []int //胡牌类型
	LeftUser     int   //玩家逃跑
	LianZhuang   int   //连庄

	//游戏信息
	CardCount    []int   //扑克数目
	HandCardData [][]int //扑克列表

	MaCount []int //码数
	MaData  []int //码数据
	Reason int //结算原因

	AllScore    []int   //总结算分
	DetailScore [][]int //单局结算分
}

// 出牌
type C2G_HZMJ_HZOutCard struct {
	CardData int
}

//出操作
type C2G_HZMJ_OperateCard struct {
	OperateCode int
	OperateCard []int
}

//请求扎码
type C2G_HZMJ_ZhaMa struct {
}

//// s to c
//用户出牌
type G2C_HZMJ_OutCard struct {
	OutCardUser int  //出牌用户
	OutCardData int  //出牌扑克
	SysOut      bool //托管系统出牌
}

type G2C_HZMJ_OperateNotify struct {
	ActionMask int //动作掩码
	ActionCard int //动作扑克
}

//发送扑克
type G2C_HZMJ_SendCard struct {
	CardData     int  //扑克数据
	ActionMask   int  //动作掩码
	CurrentUser  int  //当前用户
	SendCardUser int  //发牌用户
	ReplaceUser  int  //补牌用户
	Tail         bool //末尾发牌
}

//操作命令
type G2C_HZMJ_OperateResult struct {
	OperateUser int    //操作用户
	ActionMask  int    //动作掩码
	ProvideUser int    //供应用户
	OperateCode int    //操作代码
	OperateCard [3]int //操作扑克
}

type G2C_HZMJ_Trustee struct { //用户托管
	Trustee bool //是否托管
	ChairID int  //托管用户
}

//抓花
type G2C_HZMJ_ZhuaHua struct {
	ZhongHua []int
	BuZhong  []int
}
