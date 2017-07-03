package room

import (
	"mj/gameServer/common/pk/pk_base"
)

func NewDDZLogic(ConfigIdx int) *ddz_logic {
	l := new(ddz_logic)
	l.BaseLogic = pk_base.NewBaseLogic(ConfigIdx)
	return l
}

type ddz_logic struct {
	*pk_base.BaseLogic
}
