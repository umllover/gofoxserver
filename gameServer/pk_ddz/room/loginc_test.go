package room

import (
	"testing"

	"github.com/lovelly/leaf/log"
)

func TestCardType(t *testing.T) {
	lg := new(ddz_logic)

	var c0 []int
	log.Debug("空牌-%d", lg.GetCardType(c0))
	c1 := [...]int{0x01}
	log.Debug("单牌%d", lg.GetCardType(c1[:]))
	c2 := [...]int{0x03, 0x33}
	log.Debug("对子-%d", lg.GetCardType(c2[:]))
	c21 := [...]int{0x03, 0x31}
	log.Debug("无效两根-%d", lg.GetCardType(c21[:]))
	c3 := [...]int{0x03, 0x23, 0x33}
	log.Debug("三根-%d", lg.GetCardType(c3[:]))
	c31 := [...]int{0x04, 0x34, 0x24, 0x08}
	log.Debug("三代一-%d", lg.GetCardType(c31[:]))
	c32 := [...]int{0x04, 0x34, 0x24, 0x08, 0x18}
	log.Debug("三代二-%d", lg.GetCardType(c32[:]))
	c5 := [...]int{0x03, 0x34, 0x25, 0x06, 0x17}
	log.Debug("顺子%d", lg.GetCardType(c5[:]))
	c51 := [...]int{0x03, 0x34, 0x25, 0x06, 0x17, 0x02}
	log.Debug("带2的顺子-%d", lg.GetCardType(c51[:]))
	c4 := [...]int{0x03, 0x33, 0x24, 0x04}
	log.Debug("两个连续对子%d", lg.GetCardType(c4[:]))
	c6 := [...]int{0x03, 0x33, 0x22, 0x02, 0x14, 0x04}
	log.Debug("带2连对%d", lg.GetCardType(c6[:]))
	c61 := [...]int{0x03, 0x33, 0x25, 0x05, 0x14, 0x04}
	log.Debug("连对%d", lg.GetCardType(c61[:]))
	c62 := [...]int{0x03, 0x33, 0x23, 0x02, 0x12, 0x32}
	log.Debug("带2三顺子%d", lg.GetCardType(c62[:]))
	c63 := [...]int{0x03, 0x33, 0x23, 0x04, 0x14, 0x24}
	log.Debug("三顺子%d", lg.GetCardType(c63[:]))
	c64 := [...]int{0x03, 0x33, 0x23, 0x04, 0x14, 0x24, 0x01, 0x02}
	log.Debug("飞机带两单%d", lg.GetCardType(c64[:]))
	c65 := [...]int{0x03, 0x33, 0x23, 0x04, 0x14, 0x24, 0x01, 0x11, 0x02, 0x12}
	log.Debug("飞机带两对%d", lg.GetCardType(c65[:]))
	c41 := [...]int{0x03, 0x33, 0x23, 0x13, 0x14, 0x25}
	log.Debug("四带两单%d", lg.GetCardType(c41[:]))
	c42 := [...]int{0x03, 0x33, 0x23, 0x13, 0x14, 0x24, 0x15, 0x25}
	log.Debug("四带两对%d", lg.GetCardType(c42[:]))
	c40 := [...]int{0x03, 0x33, 0x23, 0x13}
	log.Debug("炸弹%d", lg.GetCardType(c40[:]))

	var ck []int
	for i := 0; i < 8; i++ {
		ck = append(ck, 0x4E+i%2)
		log.Debug("八王类型%d", lg.GetCardType(ck[:]))
	}
}
