package room

import (
	"mj/gameServer/common/pk_base/PKBaseLogic"
	"mj/gameServer/db/model"
)

func NewSSSEntry(info *model.CreateRoomInfo) *SSS_Entry {
	e := new(SSS_Entry)
	return e
	e.Entry_base = PKBaseLogic.NewPKBase(info)
	return e
}

///主消息入口
type SSS_Entry struct {
	*PKBaseLogic.Entry_base
}
