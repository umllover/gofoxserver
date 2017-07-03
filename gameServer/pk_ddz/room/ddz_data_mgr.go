package room

import (
	"auth_server/db/model/base"
	"mj/gameServer/common/pk/pk_base"
)

func NewDataMgr(id, uid, ConfigIdx int, name string, temp *base.GameServiceOption, base *DDZ_Entry) *ddz_data_mgr {
	d := new(ddz_data_mgr)
	d.RoomData = pk_base.NewDataMgr(id, uid, ConfigIdx, name, temp, base.Entry_base)
	return d
}

type ddz_data_mgr struct {
	*pk_base.RoomData
}
