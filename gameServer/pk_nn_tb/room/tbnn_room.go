package room

import (
	"mj/common/msg"
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/common/room_base"
	"mj/gameServer/db/model/base"

	"github.com/lovelly/leaf/log"
)

func CreaterRoom(args []interface{}) RoomMgr.IRoom {
	info := args[0].(*msg.L2G_CreatorRoom)
	if info.KindId != common.KIND_TYPE_TBNN {
		log.Error("at CreaterRoom info.KindId != common.KIND_TYPE_HZMJ uid:%d", info.CreatorUid)
		return nil
	}

	temp, ok := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	if !ok {
		log.Error("at CreaterRoom not foud template kind:%d, serverId:%d, uid:%d", info.KindId, info.ServiceId, info.CreatorUid)
		return nil
	}
	r := NewNNTBEntry(info)
	if r == nil {
		log.Error("at CreaterRoom NewMJBase error, uid:%d", info.CreatorUid)
		return nil
	}

	if r == nil {
		log.Error("at create room create entry failed")
		return nil
	}
	rbase := room_base.NewRoomBase()
	cfg := &pk_base.NewPKCtlConfig{
		BaseMgr:  rbase,
		DataMgr:  NewDataMgr(info.RoomID, info.CreatorUid, pk_base.IDX_TBNN, temp.RoomName, temp, r, info),
		UserMgr:  room_base.NewRoomUserMgr(info, temp),
		LogicMgr: NewNNTBZLogic(pk_base.IDX_TBNN),
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
