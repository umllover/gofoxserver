package log

import (
	. "mj/common/cost"
	"mj/gameServer/db/model/stats"
	"time"

	"github.com/lovelly/leaf/log"
)

type GetinRoomLog struct{}

func (getinRoomLog *GetinRoomLog) AddGetinRoomLogInfo(RoomId int, userId int64, KindId, ServerId int, RoomName string, NodeId, Num, Status, MaxPlayerCnt, PayType, public int) {
	now := time.Now()
	getInLog := &stats.GetinRoomLog{}
	if public == GIRPublic {
		getInLog.Public = GIRPublic
	} else {
		getInLog.Public = GIRPrivate
	}
	getInLog.RoomId = RoomId
	getInLog.UserId = userId
	getInLog.KindId = KindId
	getInLog.ServiceId = ServerId
	getInLog.RoomName = RoomName
	getInLog.NodeId = NodeId
	getInLog.Num = Num
	getInLog.Status = Status
	getInLog.MaxPlayerCnt = MaxPlayerCnt
	getInLog.PayType = PayType
	getInLog.GetInTime = &now

	_, err := stats.GetinRoomLogOp.Insert(getInLog)
	if err != nil {
		log.Error("添加进入房间信息到数据库失败：%s", err.Error())
	}
}
