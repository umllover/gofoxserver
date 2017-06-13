package mj_ctl_base

//麻将的通用后台配置管理类

type Ctl_base struct {
	Kind            int   //第一类型
	ServerId        int   //第二类型 注意 非房间id
	TimeLimit       int   //时间显示
	CountLimit      int   //局数限制
	TimeOutCard     int   //出牌时间
	TimeOperateCard int   //操作时间
	TimeStartGame   int64 //开始时间

}

func NewCtlBase() *Ctl_base {
	return new(Ctl_base)
}

//多久没玩家加入超时
func (c *Ctl_base) JoinTimeOut() bool {
	return false
}

//出牌超时
func (c *Ctl_base) isTimeOutOutCardTime() bool {

	return false
}

//吃胡操作
func (c *Ctl_base) isTimeOutOperateTime() bool {

	return false
}

//房间超时
func (c *Ctl_base) isRoomTimeOout() bool {

	return false
}

// 单局超时时间
func (c *Ctl_base) OneRoundTimeOut() bool {
	return false
}

//某个玩家离线后的超时解散时间
func (c *Ctl_base) OffLineTimeOut() bool {
	return false
}
