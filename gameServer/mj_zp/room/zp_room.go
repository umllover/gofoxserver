package room

import (
	"mj/common/msg"
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common"
	"mj/gameServer/common/mj/mj_base"

	"mj/gameServer/common/room_base"
	"mj/gameServer/db/model/base"

	"github.com/lovelly/leaf/log"
)

func CreaterRoom(args []interface{}) RoomMgr.IRoom {
	log.Debug("创建漳浦麻将房间！")
	info := args[0].(*msg.L2G_CreatorRoom)
	if info.KindId != common.KIND_TYPE_ZPMJ {
		log.Error("at creator zpmj error info.KindId != common.KIND_TYPE_ZPMJ ")
		return nil
	}

	temp, ok := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	if !ok {
		log.Error("at creator zpmj error not foud temolaye kindId:%d, serverId:%d ", info.KindId, info.ServiceId)
		return nil
	}

	r := NewMJBase(info)
	if r == nil {
		log.Error("at creator zpmj error NewMJBase faild  roomID:%d,", info.RoomID)
		return nil
	}
	zpBase := room_base.NewRoomBase()
	zpData := NewDataMgr(info, info.CreatorUid, mj_base.IDX_ZPMJ, "", temp, r)
	if zpData == nil {
		log.Error("at creator zpmj error NewDataMgr faild roomID:%d,", info.RoomID)
		return nil
	}

	cfg := &mj_base.NewMjCtlConfig{
		BaseMgr:  zpBase,
		DataMgr:  zpData,
		UserMgr:  room_base.NewRoomUserMgr(info, temp),
		LogicMgr: NewBaseLogic(mj_base.IDX_ZPMJ),
		TimerMgr: room_base.NewRoomTimerMgr(info.PlayCnt, temp, zpBase.GetSkeleton()),
	}

	if cfg.BaseMgr == nil || cfg.DataMgr == nil || cfg.UserMgr == nil || cfg.LogicMgr == nil || cfg.TimerMgr == nil {
		log.Error("at CreaterRoom mermber faild kind:%d, RoomID:%d uid:%d", info.KindId, info.RoomID, info.CreatorUid)
		return nil
	}
	r.Init(cfg)

	RegisterHandler(r)
	return r
}
