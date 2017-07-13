package room

import "mj/gameServer/common/mj/mj_base"

func NewXSlogic(ConfIdx int) *xs_logic {
	l := new(xs_logic)
	l.BaseLogic = mj_base.NewBaseLogic(ConfIdx)
	return l
}

type xs_logic struct {
	*mj_base.BaseLogic
}
