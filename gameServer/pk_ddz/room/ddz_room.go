package room

import (
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/common/room_base"
	"mj/gameServer/db/model/base"

	"mj/common/msg"

	"github.com/lovelly/leaf/log"
)

func CreaterRoom(args []interface{}) RoomMgr.IRoom {
	info := args[0].(*msg.L2G_CreatorRoom)
	if info.KindId != common.KIND_TYPE_DDZ {
		log.Error("at CreaterRoom info.KindId != common.KIND_TYPE_DDZ uid:%d", info.CreatorUid)
		return nil
	}

	temp, ok := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	if !ok {
		log.Error("at CreaterRoom not foud template kind:%d, serverId:%d, uid:%d", info.KindId, info.ServiceId, info.CreatorUid)
		return nil
	}
	r := NewDDZEntry(info)
	if r == nil {
		log.Error("at CreaterRoom NewPKBase error, uid:%d", info.CreatorUid)
		return nil
	}

	rbase := room_base.NewRoomBase()
	cfg := &pk_base.NewPKCtlConfig{
		BaseMgr:  rbase,
		DataMgr:  NewDDZDataMgr(info, info.CreatorUid, pk_base.IDX_DDZ, temp.RoomName, temp, r),
		UserMgr:  room_base.NewRoomUserMgr(info, temp),
		LogicMgr: NewDDZLogic(pk_base.IDX_DDZ, info),
		TimerMgr: room_base.NewRoomTimerMgr(info.PlayCnt, temp, rbase.GetSkeleton()),
	}
	if cfg.BaseMgr == nil || cfg.DataMgr == nil || cfg.UserMgr == nil || cfg.LogicMgr == nil || cfg.TimerMgr == nil {
		log.Error("at CreaterRoom mermber faild kind:%d, RoomID:%d uid:%d", info.KindId, info.RoomID, info.CreatorUid)
		return nil
	}
	r.Init(cfg)

	RegisterHandler(r)
	return r
}
