package room

import (
	"mj/gameServer/common"
	"mj/gameServer/common/mj_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"mj/common/msg/mj_zp_msg"

	"time"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/timer"
)

type ZP_RoomData struct {
	*mj_base.RoomData
	ChaHuaTime *timer.Timer

	ChaHuaMap map[int]int
}

func NewDataMgr(id, uid int, name string, temp *base.GameServiceOption) *ZP_RoomData {
	r := new(ZP_RoomData)
	r.ChaHuaMap = make(map[int]int)
	r.RoomData = mj_base.NewDataMgr(id, uid, name, temp)
	return r
}

func (room *ZP_RoomData) StartGameing(userMgr common.UserManager, gameLogic common.LogicManager, template *base.GameServiceOption) {

}

func (room *ZP_RoomData) AfterStartGame(userMgr common.UserManager, gameLogic common.LogicManager) {

}

func (room *ZP_RoomData) SetChaHua(arg interface{}) {
	bRoom := arg.(zpMj_base)
	bRoom.UserMgr.SendMsgAll(&mj_zp_msg.G2C_MJZP_GetChaHua{})
	bRoom.GetSkeleton().AfterFunc(time.Duration(bRoom.Temp.OutCardTime)*time.Second, func() {
		bRoom.StartPlay()
	})
}

func (room *ZP_RoomData) GetChaHua(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	getData := &mj_zp_msg.C2G_MJZP_SetChaHua{}
	room.ChaHuaMap[user.ChairId] = getData.SetCount
}
