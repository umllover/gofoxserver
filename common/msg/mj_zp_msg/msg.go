package mj_zp_msg

import "mj/common/msg"

func init() {
	msg.Processor.Register(&C2G_MJZP_SetChaHua{})
	msg.Processor.Register(&G2C_MJZP_FlowerCard{})
	msg.Processor.Register(&C2G_MJZP_OperateNotify{})
}

//获取插花
type G2C_MJZP_GetChaHua struct {
}

//设置插花
type C2G_MJZP_SetChaHua struct {
	SetCount int //设置插花数量
}

//补花
type G2C_MJZP_FlowerCard struct {
	ReplaceUser  int  //补牌用户
	ReplaceCard  int  //补牌扑克
	NewCard      int  //补完扑克
	IsInitFlower bool //是否开局补花，true开局补花
}

//操作提示
type C2G_MJZP_OperateNotify struct {
	ActionMask int //动作掩码
	ActionCard int //动作扑克
}
