package idGenerate

import (
	"mj/common/utils"

	"mj/gameServer/db/model"
)

var (
	ids = make(map[int]*model.RoomId)
)

func GetRoomId(uid int) (int, bool) {
	for i := 0; i < 100; i++ {
		r, rerr := utils.RandInt(100000, 1000000)
		if rerr != nil {
			continue
		}
		if _, ok := ids[r]; ok {
			continue
		}
		rid := &model.RoomId{Id: r, UserId: uid}
		_, err := model.RoomIdOp.Insert(rid)
		if err == nil {
			ids[r] = rid
			return r, true
		}
		continue
	}
	return 0, false
}

func DelRoomId(rid int) {
	delete(ids, rid)
	model.RoomIdOp.Delete(rid)
}
