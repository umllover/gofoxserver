package room

import "mj/gameServer/common/mj/mj_base"

func NewHZlogic(ConfIdx int) *hz_logic {
	l := new(hz_logic)
	l.BaseLogic = mj_base.NewBaseLogic(ConfIdx)
	return l
}

type hz_logic struct {
	*mj_base.BaseLogic
}
