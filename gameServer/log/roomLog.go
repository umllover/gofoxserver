package log

import (
	"mj/gameServer/db/model/stats"

	"time"

	. "github.com/lovelly/leaf/log"
	"github.com/name5566/leaf/log"
)

type RoomLog struct{}

//添加创建房间记录
func (roomLog *RoomLog) AddCreateRoomLog(roomId int, userId int64, roomName string, kindId, serverId, nodeId int, createTime time.Time, payType, retCode int) {

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

//游戏结束记录类型
func (RoomLog *RoomLog) UpdateGameLogRecode(recodeId int, code int) {
	if recodeId >= 0 {
		log.Error("没有这条记录存在")
		return
	}
	myLogInfo := make(map[string]interface{})
	myLogInfo["game_end_type"] = code
	err := stats.RoomLogOp.UpdateWithMap(recodeId, myLogInfo)
	if err != nil {
		Error("结束时间和结束状态记录更新失败：%s", err.Error())
	}
}

//更新解散房间记录类型
func (roomLog *RoomLog) UpdateRoomLogRecode(recodeId int, time time.Time, code int) {
	if recodeId >= 0 {
		log.Error("没有这条记录存在")
		return
	}
	myLogInfo := make(map[string]interface{})
	myLogInfo["end_time"] = &time
	myLogInfo["room_end_type"] = code
	err := stats.RoomLogOp.UpdateWithMap(recodeId, myLogInfo)
	if err != nil {
		Error("结束时间和结束状态记录更新失败：%s", err.Error())
	}
}

//更新是否为他人创建房间
func (RoomLog *RoomLog) UpdateRoomLogForOthers(recodeId int, code int) {
	if recodeId >= 0 {
		log.Error("没有这条记录存在")
		return
	}
	myLogInfo := make(map[string]interface{})
	myLogInfo["create_others"] = code
	err := stats.RoomLogOp.UpdateWithMap(recodeId, myLogInfo)
	if err != nil {
		Error("是否为他人开房状态记录更新失败：%s", err.Error())
	}
}
