package room

import (
	. "mj/common/cost"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"
)

func NewSSSEntry(info *model.CreateRoomInfo) *SSS_Entry {
	e := new(SSS_Entry)
	e.Entry_base = pk_base.NewPKBase(info)
	return e
}

///主消息入口
type SSS_Entry struct {
	*pk_base.Entry_base
}

//玩家准备
func (room *SSS_Entry) UserReady(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserReady)
	u := args[1].(*user.User)
	if u.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	log.Debug("at UserReady ==== ")
	room.UserMgr.SetUsetStatus(u, US_READY)

	log.Debug("ren shu %d", room.UserMgr.GetCurPlayerCnt()) //|| room.UserMgr.GetCurPlayerCnt() >= 1
	if room.UserMgr.IsAllReady() {
		//派发初始扑克
		room.DataMgr.BeforeStartGame(room.UserMgr.GetMaxPlayerCnt())
		room.DataMgr.StartGameing()
		room.DataMgr.AfterStartGame()

		room.Status = RoomStatusStarting
		room.TimerMgr.StartPlayingTimer(room.GetSkeleton(), func() {
			room.OnEventGameConclude(0, nil, GER_DISMISS)
		})
	}
}
