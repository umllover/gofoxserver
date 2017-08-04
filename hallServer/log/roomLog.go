package log

import (
	"mj/hallServer/db/model/stats"

	"time"

	. "github.com/lovelly/leaf/log"
)

type RoomLog struct{}

//添加创建房间记录
func (roomLog *RoomLog) AddCreateRoomLog(roomId int, userId int64, roomName string, kindId, serverId, nodeId int, payType, retCode int) {
	createTime := time.Now()
	Info := &stats.RoomLog{
		UserId:     userId,
		PayType:    payType,
		RoomId:     roomId,
		RoomName:   roomName,
		NodeId:     nodeId,
		KindId:     kindId,
		ServiceId:  serverId,
		CreateTime: &createTime,
	}
	if retCode == 0 {
		Info.NomalOpen = 1
	} else {
		Info.NomalOpen = 0
	}
	_, err := stats.RoomLogOp.Insert(Info)
	if err != nil {
		Error("insert Data into table roomlog Error:%v", err.Error())
	}
}

//获取创建房间记录
func (roomLog *RoomLog) GetRoomLogRecode(roomId, kindId, serverId int) (roomRecord *stats.RoomLog) {
	logInfo := make(map[string]interface{})
	logInfo["room_id"] = roomId
	logInfo["kind_id"] = kindId
	logInfo["service_id"] = serverId
	logData, err := stats.RoomLogOp.GetByMap(logInfo)
	if err != nil {
		Error("Select Data from recode Error:%v", err.Error())
		return nil
	}
	return logData
}

//更新创建房间记录
func (RoomLog *RoomLog) UpdateRoomLogRecode(roomLog *stats.RoomLog, time time.Time, code int) {
	if roomLog == nil {
		Error("没有这条记录存在")
		return
	}
	myLogInfo := make(map[string]interface{})
	myLogInfo["end_time"] = &time
	myLogInfo["end_type"] = code
	err := stats.RoomLogOp.UpdateWithMap(roomLog.RecodeId, myLogInfo)
	if err != nil {
		Error("结束时间和结束状态记录更新失败：%s", err.Error())
	}
}
