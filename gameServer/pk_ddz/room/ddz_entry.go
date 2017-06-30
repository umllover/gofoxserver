package room

import (
	"mj/gameServer/common/pk_base/PKBaseLogic"
	"mj/gameServer/db/model"
)

func NewDDZEntry(info *model.CreateRoomInfo) *DDZ_Entry {
	e := new(DDZ_Entry)
	return e
	e.Entry_base = PKBaseLogic.NewPKBase(info)
	return e
}

///主消息入口
type DDZ_Entry struct {
	*PKBaseLogic.Entry_base
}


