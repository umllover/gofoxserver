package room

import (
	. "mj/common/cost"
	"mj/common/msg/pk_sss_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/user"

	"mj/common/msg"

	"github.com/lovelly/leaf/log"
)

func NewSSSEntry(info *msg.L2G_CreatorRoom) *SSS_Entry {
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
		room.TimerMgr.StopCreatorTimer()
	}
}

// 十三水摊牌
func (r *SSS_Entry) ShowSSsCard(args []interface{}) {
	recvMsg := args[0].(*pk_sss_msg.C2G_SSS_Open_Card)
	u := args[1].(*user.User)

	r.DataMgr.ShowSSSCard(u, recvMsg.Dragon, recvMsg.SpecialType, recvMsg.SpecialData, recvMsg.FrontCard, recvMsg.MidCard, recvMsg.BackCard)
	return
}
