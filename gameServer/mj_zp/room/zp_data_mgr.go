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

func NewDataMgr(id, uid int, name string, temp *base.GameServiceOption, base *Mj_base) *ZP_RoomData {
	r := new(ZP_RoomData)
	r.ChaHuaMap = make(map[int]int)
	r.RoomData = mj_base.NewDataMgr(id, uid, name, temp, base)
	return r
}

func (room *ZP_RoomData) StartGameing() {
	if room.MjBase.TimerMgr.GetPlayCount() == 0 {
		room.MjBase.UserMgr.SendMsgAll(&mj_zp_msg.G2C_MJZP_GetChaHua{})
		room.ChaHuaTime = room.MjBase.AfterFunc(time.Duration(room.MjBase.Temp.OutCardTime)*time.Second, func() {
			room.StartDispatchCard(room.MjBase.UserMgr, room.MjBase.LogicMgr, room.MjBase.Temp)
			//检查自摸
			room.CheckZiMo(room.MjBase.LogicMgr, room.MjBase.UserMgr)
			//通知客户端开始了
			room.SendGameStart()
		})
	} else {
		room.StartDispatchCard(room.MjBase.UserMgr, room.MjBase.LogicMgr, room.MjBase.Temp)
		//检查自摸
		room.CheckZiMo(room.MjBase.LogicMgr, room.MjBase.UserMgr)
		//通知客户端开始了
		room.SendGameStart(room.MjBase.LogicMgr, room.MjBase.UserMgr)
	}
}

func (room *ZP_RoomData) AfterStartGame(userMgr common.UserManager, gameLogic common.LogicManager) {

}

func (room *ZP_RoomData) GetChaHua(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	getData := &mj_zp_msg.C2G_MJZP_SetChaHua{}
	room.ChaHuaMap[user.ChairId] = getData.SetCount
	if len(room.ChaHuaMap) == 4 {
		room.StartDispatchCard(room.MjBase.UserMgr, room.MjBase.LogicMgr, room.MjBase.Temp)
		//检查自摸
		room.CheckZiMo(room.MjBase.LogicMgr, room.MjBase.UserMgr)
		//通知客户端开始了
		room.SendGameStart(room.MjBase.LogicMgr, room.MjBase.UserMgr)
	}
}
