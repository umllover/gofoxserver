package room

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common"
	"mj/gameServer/common/mj_base"
	"mj/gameServer/db/model"
	"mj/gameServer/user"

	"mj/gameServer/common/room_base"
	"mj/gameServer/db/model/base"
)

func CreaterRoom(args []interface{}) RoomMgr.IRoom {
	info := args[0].(*model.CreateRoomInfo)

	u := args[1].(*user.User)
	retCode := 0
	defer func() {
		if retCode != 0 {
			u.WriteMsg(&msg.L2C_CreateTableFailure{ErrorCode: retCode, DescribeString: "创建房间失败"})
		}
	}()

	if info.KindId != common.KIND_TYPE_HZMJ {
		retCode = ErrParamError
		return nil
	}

	temp, ok := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	if !ok {
		retCode = NoFoudTemplate
		return nil
	}
	r := mj_base.NewMJBase(info)
	cfg := &mj_base.NewMjCtlConfig{
		BaseMgr:  room_base.NewRoomBase(),
		DataMgr:  NewDataMgr(info.RoomId, u.Id, mj_base.IDX_ZPMJ, temp.GameName, temp, r),
		UserMgr:  room_base.NewRoomUserMgr(info.RoomId, info.MaxPlayerCnt, temp),
		LogicMgr: NewBaseLogic(),
		TimerMgr: room_base.NewRoomTimerMgr(info.Num, temp),
	}
	r.Init(cfg)
	if r == nil {
		retCode = Errunlawful
		return nil
	}

	u.KindID = info.KindId
	u.RoomId = r.DataMgr.GetRoomId()
	RegisterHandler(r)
	return r
}
