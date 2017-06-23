package mj_zp_msg

import "mj/common/msg"

func init() {
	msg.Processor.Register(&C2G_MJZP_SetChaHua{})

}

//获取插花
type G2C_MJZP_GetChaHua struct {
}

//设置插花
type C2G_MJZP_SetChaHua struct {
	SetCount int //设置插花数量
}
