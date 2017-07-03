package room

import (
	"mj/gameServer/common/pk/pk_base"
)

func NewDDZLogic() *ddz_logic {
	l := new(ddz_logic)
	l.BaseLogic = pk_base.NewBaseLogic()
	return l
}

type ddz_logic struct {
	*pk_base.BaseLogic
}
