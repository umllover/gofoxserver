package room

import (
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"
)

func NewDDZEntry(info *model.CreateRoomInfo) *DDZ_Entry {
	e := new(DDZ_Entry)
	return e
	e.Entry_base = pk_base.NewPKBase(info)
	return e
}

///主消息入口
type DDZ_Entry struct {
	*pk_base.Entry_base
}
