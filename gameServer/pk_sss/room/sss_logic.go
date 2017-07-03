package room

import "mj/gameServer/common/pk/pk_base"

func NewNNTBZLogic() *sss_logic {
	l := new(sss_logic)
	l.BaseLogic = pk_base.NewBaseLogic()
	return l
}

type sss_logic struct {
	*pk_base.BaseLogic
}
