package user

import (
	"mj/common/msg"
	"mj/common/msg/mj_zp_msg"
	"mj/gameServer/RoomMgr"

	"github.com/lovelly/leaf/log"
)

func NewRobot(uid int64, player *User, room RoomMgr.IRoom) *Robot {
	robot := new(Robot)
	robot.UserID = uid
	robot.Player = player
	robot.ch = make(chan interface{}, 10)
	robot.closeCh = make(chan bool, 1)
	robot.room = room
	robot.Run()
	return robot
}

type Robot struct {
	Player  *User
	UserID  int64
	ch      chan interface{}
	closeCh chan bool
	room    RoomMgr.IRoom
}

func (r *Robot) WriteMsg(msg interface{}) {
	r.ch <- msg
}

func (r *Robot) Run() {
	go r.run()
}

func (r *Robot) Close() {
	r.closeCh <- true
}

func (r *Robot) run() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Robot run error:%v", err)
			r.Run()
		}
	}()

	for {
		select {
		case _ = <-r.closeCh:
			return
		case data := <-r.ch:
			r.HandleMsg(data)
		}

	}
}

func (r *Robot) HandleMsg(data interface{}) {
	msgName, err := msg.Processor.GetMsgId(data)
	if err != nil {
		log.Error("at HandleMsg error: %s", err.Error())
		return
	}
	r.Player.WriteMsg(data)
	switch msgName {
	case "G2C_ZPMJ_SendCard":
		recvMsg := data.(*mj_zp_msg.G2C_ZPMJ_SendCard)
		r.room.GetChanRPC().Go("OutCard", r.Player, recvMsg.CardData)
	default:

	}
}
