package room

import "mj/gameServer/common/pk/pk_base"

func NewDataMgr() *nntb_data_mgr {
	d := new(nntb_data_mgr)
	d.RoomData = pk_base.NewDataMgr()
}

type nntb_data_mgr struct {
	*pk_base.RoomData
}
