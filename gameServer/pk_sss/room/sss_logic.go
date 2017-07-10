package room

import "mj/gameServer/common/pk/pk_base"

func NewSssZLogic(ConfigIdx int) *sss_logic {
	l := new(sss_logic)
	l.BaseLogic = pk_base.NewBaseLogic(ConfigIdx)
	return l
}

type sss_logic struct {
	*pk_base.BaseLogic
}
