package room

import (
	"mj/common/msg"
	"mj/gameServer/common/pk/pk_base"
)

func NewNNTBEntry(info *msg.L2G_CreatorRoom) *NNTB_Entry {
	e := new(NNTB_Entry)
	e.Entry_base = pk_base.NewPKBase(info)
	return e
}

///主消息入口
type NNTB_Entry struct {
	*pk_base.Entry_base
}
