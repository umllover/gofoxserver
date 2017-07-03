package room

import (
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model"
)

func NewXSEntry(info *model.CreateRoomInfo) *xs_entry {
	e := new(xs_entry)
	e.Mj_base = mj_base.NewMJBase(info)
	return e
}

type xs_entry struct {
	*mj_base.Mj_base
}
