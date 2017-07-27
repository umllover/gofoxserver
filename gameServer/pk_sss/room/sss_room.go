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
	if info.KindId != common.KIND_TYPE_SSS {
		log.Debug("at CreaterRoom info.KindId != common.KIND_TYPE_SSS uid:%d", u.Id)
		return nil
	}

	temp, ok := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	if !ok {
		log.Debug("at CreaterRoom not foud template kind:%d, serverId:%d, uid:%d", info.KindId, info.ServiceId, u.Id)
		return nil
	}
	r := NewSSSEntry(info)
	cfg := &pk_base.NewPKCtlConfig{
		BaseMgr:  room_base.NewRoomBase(),
		DataMgr:  NewDataMgr(info, u.Id, pk_base.IDX_SSS, temp.RoomName, temp, r, info.OtherInfo),
		UserMgr:  room_base.NewRoomUserMgr(info, temp),
		LogicMgr: NewSssZLogic(pk_base.IDX_SSS),
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
