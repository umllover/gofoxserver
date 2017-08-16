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
	if info.KindId != common.KIND_TYPE_SSS {
		log.Debug("at CreateRoom info.KindId != common.KIND_TYPE_SSS uid:%d", info.CreatorUid)
		return nil
	}

	temp, ok := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	if !ok {
		log.Debug("at CreateRoom not foud template kind:%d, serverId:%d, uid:%d", info.KindId, info.ServiceId, info.CreatorUid)
		return nil
	}
	r := NewSSSEntry(info)
	rbase := room_base.NewRoomBase()
	cfg := &pk_base.NewPKCtlConfig{
		BaseMgr:  rbase,
		DataMgr:  NewDataMgr(info, info.CreatorUid, pk_base.IDX_SSS, temp.RoomName, temp, r),
		UserMgr:  room_base.NewRoomUserMgr(info, temp),
		LogicMgr: NewSssZLogic(pk_base.IDX_SSS),
		TimerMgr: room_base.NewRoomTimerMgr(info.PlayCnt, temp, rbase.GetSkeleton()),
	}
	r.Init(cfg)
	if r == nil {
		log.Debug("at CreateRoom NewSSS error, uid:%d", info.CreatorUid)
		return nil
	}

	RegisterHandler(r)
	return r
}
