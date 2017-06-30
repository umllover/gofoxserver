package room

import (
	"mj/gameServer/common/pk_base/PKBaseLogic"
	"mj/gameServer/db/model"
)

func NewNNTBEntry(info *model.CreateRoomInfo) *NNTB_Entry {
	e := new(NNTB_Entry)
	return e
	e.Entry_base = PKBaseLogic.NewPKBase(info)
	return e
}

///主消息入口
type NNTB_Entry struct {
	*PKBaseLogic.Entry_base
}
