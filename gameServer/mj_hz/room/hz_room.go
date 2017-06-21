package room

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/common"
	"mj/gameServer/common/mj_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/idGenerate"
	"mj/gameServer/user"

	"mj/gameServer/common/room_base"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

func CreaterRoom(args []interface{}) room_base.Module {
	recvMsg := args[0].(*msg.C2G_CreateTable)
	retMsg := &msg.G2C_CreateTableSucess{}
	agent := args[1].(gate.Agent)
	retCode := 0
	defer func() {
		if retCode == 0 {
			agent.WriteMsg(retMsg)
		} else {
			agent.WriteMsg(&msg.G2C_CreateTableFailure{ErrorCode: retCode, DescribeString: "创建房间失败"})
		}
	}()

	u := agent.UserData().(*user.User)
	if recvMsg.Kind != common.KIND_TYPE_HZMJ {
		retCode = CreateParamError
		return nil
	}

	template, ok := base.GameServiceOptionCache.Get(recvMsg.Kind, recvMsg.ServerId)
	if !ok {
		retCode = NoFoudTemplate
		return nil
	}

	feeTemp, ok1 := base.PersonalTableFeeCache.Get(recvMsg.ServerId, recvMsg.Kind, recvMsg.DrawCountLimit, recvMsg.DrawTimeLimit)
	if !ok1 {
		log.Error("not foud PersonalTableFeeCache")
		retCode = NoFoudTemplate
		return nil
	}

	rid, iok := idGenerate.GetRoomId(u.Id)
	if !iok {
		retCode = RandRoomIdError
		return nil
	}

	if recvMsg.CellScore > template.CellScore {
		retCode = MaxSoucrce
		return nil
	}

	r := mj_base.NewMJBase(recvMsg.RoomID, u.Id, room_base.NewRoomUserMgr)
	if r == nil {
		retCode = Errunlawful
		return nil
	}

	retMsg.TableID = r.DataMgr.GetRoomId()
	retMsg.DrawCountLimit = r.DataMgr.GetCountLimit()
	retMsg.DrawTimeLimit = r.DataMgr.GetTimeLimit()
	retMsg.Beans = feeTemp.TableFee
	retMsg.RoomCard = u.RoomCard
	u.KindID = recvMsg.Kind
	u.RoomId = r.DataMgr.GetRoomId()
	RegisterHandler(r)
	return r
}
