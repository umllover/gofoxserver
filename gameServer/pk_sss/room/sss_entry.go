package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"
)

func NewSSSEntry(info *model.CreateRoomInfo) *SSS_Entry {
	e := new(SSS_Entry)
	e.Entry_base = pk_base.NewPKBase(info)
	return e
}

///主消息入口
type SSS_Entry struct {
	*pk_base.Entry_base
}
