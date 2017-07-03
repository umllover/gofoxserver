package room

import "mj/gameServer/common/pk/pk_base"

func NewNNTBZLogic(ConfigIdx int) *nntb_logic {
	l := new(nntb_logic)
	l.BaseLogic = pk_base.NewBaseLogic(ConfigIdx)
	return l
}

type nntb_logic struct {
	*pk_base.BaseLogic
}
