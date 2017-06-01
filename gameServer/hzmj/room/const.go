package room

const (
	//用户状态
	US_NULL		=				0x00								//没有状态
	US_FREE		=				0x01								//站立状态
	US_SIT		=				0x02								//坐下状态
	US_READY	=				0x03								//同意状态
	US_LOOKON	=				0x04								//旁观状态
	US_PLAYING	=				0x05								//游戏状态
	US_OFFLINE	=				0x06								//断线状态
)

const (
	//房间状态
	RoomStatusReady = 0
	RoomStatusStarting = 1
	RoomStatusEnd = 2
)

const (
	//动作标志
	WIK_NULL		=			0x00								//没有类型
	WIK_LEFT		=			0x01								//左吃类型
	WIK_CENTER		=			0x02								//中吃类型
	WIK_RIGHT		=			0x04								//右吃类型
	WIK_PENG		=			0x08								//碰牌类型
	WIK_GANG		=			0x10								//杠牌类型
	WIK_LISTEN		=			0x20								//听牌类型
	WIK_CHI_HU		=			0x40								//吃胡类型
	WIK_FANG_PAO	=			0x80								//放炮
)

type HistoryScore struct {
	TurnScore int
	CollectScore int
}