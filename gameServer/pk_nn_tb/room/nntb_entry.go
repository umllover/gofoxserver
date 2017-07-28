package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"
)

func NewNNTBEntry(info *model.CreateRoomInfo) *NNTB_Entry {
	e := new(NNTB_Entry)
	e.Entry_base = pk_base.NewPKBase(info)
	return e
}

///主消息入口
type NNTB_Entry struct {
	*pk_base.Entry_base
}