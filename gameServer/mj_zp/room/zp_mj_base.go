package room

import (
	. "mj/common/cost"
	"mj/gameServer/common/mj_base"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"
)

type zpMj_base struct {
	*mj_base.Mj_base
}

func (room *zpMj_base) UserReady(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserReady)
	u := args[1].(*user.User)
	if u.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	log.Debug("at UserReady ==== ")
	room.UserMgr.SetUsetStatus(u, US_READY)
	if room.UserMgr.IsAllReady() {
		//初始房间
		room.DataMgr.InitRoom(room.UserMgr.GetMaxPlayerCnt())
		//第一局抓花
		if room.TimerMgr.GetPlayCount() == 0 {
			room.DataMgr.SetChaHua(room)
		} else {
			room.StartPlay()
		}
	}
}

func (room *zpMj_base) StartPlay() {
	//派发初始扑克
	room.DataMgr.StartDispatchCard(room.UserMgr, room.LogicMgr, room.Temp)
	room.Status = RoomStatusStarting
	//检查自摸
	room.DataMgr.CheckZiMo(room.LogicMgr, room.UserMgr)
	//通知客户端开始了
	room.DataMgr.SendGameStart(room.LogicMgr, room.UserMgr, room.TimerMgr)
}
