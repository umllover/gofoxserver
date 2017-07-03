package room

import "mj/gameServer/common/pk/pk_base"

func NewNNTBZLogic() *nntb_logic {
	l := new(nntb_logic)
	l.BaseLogic = pk_base.NewBaseLogic()
	return l
}

type nntb_logic struct {
	*pk_base.BaseLogic
}
