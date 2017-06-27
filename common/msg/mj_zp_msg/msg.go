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

}

//获取插花
type G2C_MJZP_GetChaHua struct {
}

//设置插花
type C2G_MJZP_SetChaHua struct {
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
	ListenUser  int     //听牌用户
	IsListen    bool    //是否听牌
	HuCardCount int     //胡几张牌
	HuCardData  [42]int //胡牌数据
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
