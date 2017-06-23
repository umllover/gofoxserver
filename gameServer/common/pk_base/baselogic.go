package pk_base

import (
	"mj/common/msg"
	"mj/gameServer/common"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

// 扑克通用逻辑
const (
	LOGIC_MASK_COLOR	=			0xF0								//花色掩码
	LOGIC_MASK_VALUE	=			0x0F								//数值掩码
)
//获取数值
func GetCardValue(CardData int) int {
	return CardData&LOGIC_MASK_VALUE
}
//获取花色
func GetCardColor(CardData int) int {
	return CardData&LOGIC_MASK_COLOR
}

