package mj_logic_base

const (
	MAX_INDEX = 42 //最大索引

	//动作标志
	WIK_NULL     = 0x00 //没有类型
	WIK_LEFT     = 0x01 //左吃类型
	WIK_CENTER   = 0x02 //中吃类型
	WIK_RIGHT    = 0x04 //右吃类型
	WIK_PENG     = 0x08 //碰牌类型
	WIK_GANG     = 0x10 //杠牌类型
	WIK_LISTEN   = 0x20 //听牌类型
	WIK_CHI_HU   = 0x40 //吃胡类型
	WIK_FANG_PAO = 0x80 //放炮
)

//逻辑掩码
const (
	MASK_COLOR = 0xF0 //花色掩码
	MASK_VALUE = 0x0F //数值掩码
)
