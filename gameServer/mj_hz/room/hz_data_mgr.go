package room

import (
	"mj/common/msg"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"

	"github.com/lovelly/leaf/log"
)

type hz_data struct {
	*mj_base.RoomData
}

func NewHZDataMgr(id int, uid int64, configIdx int, name string, temp *base.GameServiceOption, base *hz_entry, info *msg.L2G_CreatorRoom) *hz_data {
	d := new(hz_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base, info)

	getData, ok := d.OtherInfo["zhaMa"].(float64)
	if !ok {
		log.Error("hzmj at NewDataMgr [zhaMa] error")
		return nil
	}

	//TODO 客户端发的个数有误，暂时强制改掉
	getData = 2

	d.ZhuaHuaCnt = int(getData)

	return d
}
