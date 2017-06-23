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
)

func CreaterRoom(args []interface{}) RoomMgr.IRoom {
	info := args[0].(*model.CreateRoomInfo)

	u := args[1].(user.User)
	retCode := 0
	defer func() {
		if retCode != 0 {
			u.WriteMsg(&msg.G2C_CreateTableFailure{ErrorCode: retCode, DescribeString: "创建房间失败"})
		}
	}()

	if info.KindId != common.KIND_TYPE_HZMJ {
		retCode = ErrParamError
		return nil
	}

	cfg := &mj_base.NewMjCtlConfig{
		NUserF:  room_base.NewRoomUserMgr,
		NDataF:  mj_base.NewDataMgr,
		NBaseF:  room_base.NewRoomBase,
		NLogicF: mj_base.NewBaseLogic,
		NTimerF: room_base.NewRoomTimerMgr,
	}
	r := mj_base.NewMJBase(info, u.Id, 0, info.Num, 0, 0, 4, cfg)
	if r == nil {
		retCode = Errunlawful
		return nil
	}

	u.KindID = info.KindId
	u.RoomId = r.DataMgr.GetRoomId()
	RegisterHandler(r)
	return r
}
