package room

import (
	"mj/common/msg"
	"mj/gameServer/common/mj/mj_base"
)

func NewXSEntry(info *msg.L2G_CreatorRoom) *xs_entry {
	e := new(xs_entry)
	e.Mj_base = mj_base.NewMJBase(info.KindId, info.ServiceId)
	return e
}

type xs_entry struct {
	*mj_base.Mj_base
}
