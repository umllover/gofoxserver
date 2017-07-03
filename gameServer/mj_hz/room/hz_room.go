package room

import (
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model"
	"mj/gameServer/user"

	"mj/gameServer/common/room_base"
	"mj/gameServer/db/model/base"

	"github.com/lovelly/leaf/log"
)

func CreaterRoom(args []interface{}) RoomMgr.IRoom {
	info := args[0].(*model.CreateRoomInfo)
	u := args[1].(*user.User)

	if info.KindId != common.KIND_TYPE_HZMJ {
		log.Debug("at CreaterRoom info.KindId != common.KIND_TYPE_HZMJ uid:%d", u.Id)
		return nil
	}

	temp, ok := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	if !ok {
		log.Debug("at CreaterRoom not foud template kind:%d, serverId:%d, uid:%d", info.KindId, info.ServiceId, u.Id)
		return nil
	}
	r := NewHZEntry(info)
	cfg := &mj_base.NewMjCtlConfig{
		BaseMgr:  room_base.NewRoomBase(),
		DataMgr:  NewHZDataMgr(info.RoomId, u.Id, mj_base.IDX_HZMJ, "", temp, r),
		UserMgr:  room_base.NewRoomUserMgr(info.RoomId, info.MaxPlayerCnt, temp),
		LogicMgr: NewHZlogic(mj_base.IDX_HZMJ),
		TimerMgr: room_base.NewRoomTimerMgr(info.Num, temp),
	}
	r.Init(cfg)
	if r == nil {
		log.Debug("at CreaterRoom NewMJBase error, uid:%d", u.Id)
		return nil
	}

	u.KindID = info.KindId
	u.RoomId = r.DataMgr.GetRoomId()
	RegisterHandler(r)
	return r
}
