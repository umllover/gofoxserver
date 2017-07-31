package room

import (
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/common/room_base"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"
)

func CreaterRoom(args []interface{}) RoomMgr.IRoom {
	info := args[0].(*model.CreateRoomInfo)
	u := args[1].(*user.User)
	if info.KindId != common.KIND_TYPE_DDZ {
		log.Debug("at CreaterRoom info.KindId != common.KIND_TYPE_DDZ uid:%d", u.Id)
		return nil
	}

	temp, ok := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	if !ok {
		log.Debug("at CreaterRoom not foud template kind:%d, serverId:%d, uid:%d", info.KindId, info.ServiceId, u.Id)
		return nil
	}
	r := NewDDZEntry(info)
	rbase := room_base.NewRoomBase()
	cfg := &pk_base.NewPKCtlConfig{
		BaseMgr:  rbase,
		DataMgr:  NewDDZDataMgr(info, u.Id, pk_base.IDX_DDZ, temp.RoomName, temp, r),
		UserMgr:  room_base.NewRoomUserMgr(info, temp),
		LogicMgr: NewDDZLogic(pk_base.IDX_DDZ, info),
		TimerMgr: room_base.NewRoomTimerMgr(info.Num, temp, rbase.GetSkeleton()),
	}
	r.Init(cfg)
	if r == nil {
		log.Debug("at CreaterRoom NewPKBase error, uid:%d", u.Id)
		return nil
	}

	u.KindID = info.KindId
	u.RoomId = r.DataMgr.GetRoomId()
	RegisterHandler(r)
	return r
}
