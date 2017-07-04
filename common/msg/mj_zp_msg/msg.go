package mj_zp_msg

import "mj/common/msg"

func init() {
	msg.Processor.Register(&C2G_MJZP_SetChaHua{})
	msg.Processor.Register(&C2G_MJZP_OperateNotify{})
	msg.Processor.Register(&C2G_MJZP_ReplaceCard{})
	msg.Processor.Register(&C2G_MJZP_ListenCard{})
	msg.Processor.Register(&C2G_ZPMJ_OutCard{})
	msg.Processor.Register(&C2G_ZPMJ_OperateCard{})

	msg.Processor.Register(&G2C_MJZP_GetChaHua{})
	msg.Processor.Register(&G2C_MJZP_ReplaceCard{})
	msg.Processor.Register(&G2C_MJZP_ListenCard{})
	msg.Processor.Register(&G2C_ZPMJ_GameConclude{})
	msg.Processor.Register(&G2C_ZPMJ_ZhuaHua{})
	msg.Processor.Register(&C2G_MJZP_AllChaHua{})
	msg.Processor.Register(&G2C_ZPMG_GameStart{})
	msg.Processor.Register(&C2G_MJZP_Trustee{})

}

//获取插花
type G2C_MJZP_GetChaHua struct {
}

//设置插花
type C2G_MJZP_SetChaHua struct {
	SetCount int //设置插花数量
}

//超时插花
type C2G_MJZP_AllChaHua struct {
	ChaHuaCnt [4]int
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

//操作提示
type C2G_MJZP_OperateNotify struct {
	ActionMask int //动作掩码
	ActionCard int //动作扑克
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
}

//抓花
type G2C_ZPMJ_ZhuaHua struct {
	ZhongHua []int
	BuZhong  []int
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
	Trustee int //1：托管 0：取消托管
}
