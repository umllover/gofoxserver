package room

import (
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"
)

func NewHZDataMgr(id, uid, configIdx int, name string, temp *base.GameServiceOption, base *hz_data) *hz_data {
	d := new(hz_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.MjBase)
}

type hz_data struct {
	*mj_base.RoomData
}
