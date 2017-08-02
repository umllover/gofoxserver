package log

import (
	"mj/gameServer/db/model/stats"
	"time"

	. "github.com/lovelly/leaf/log"
)

type RoomLog struct{}

//添加创建房间记录
func (roomLog *RoomLog) AddCreateRoomLog(roomId int, userId int64, roomName string, kindId, serverId, nodeId int, createTime time.Time, payType, retCode int) {

	logInfo := &stats.RoomLog{
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
		logInfo.NomalOpen = 1
	} else {
		logInfo.NomalOpen = 0
	}
	logInfo.CreateOthers = 1
	_, err := stats.RoomLogOp.Insert(logInfo)
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
	}
	return logData
}

//更新创建房间记录
func (roomLog *RoomLog) UpdateRoomLogRecode(recodeId int, time time.Time, code int) {
	myLogInfo := make(map[string]interface{})
	myLogInfo["end_time"] = &time
	myLogInfo["end_type"] = code
	err := stats.RoomLogOp.UpdateWithMap(recodeId, myLogInfo)
	if err != nil {
		Error("结束时间和结束状态记录更新失败：%s", err.Error())
	}
}
