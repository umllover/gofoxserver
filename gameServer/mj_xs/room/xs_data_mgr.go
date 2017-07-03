package room

import (
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"
)

func NewXSDataMgr(id, uid, configIdx int, name string, temp *base.GameServiceOption, base *xs_entry) *xs_data {
	d := new(xs_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base)
}

type xs_data struct {
	*mj_base.RoomData
	ZhuaHuaCnt   int //扎花个数
	ZhuaHuaScore int //扎花分数
}
