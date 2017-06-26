package room_base

import (
	"time"

	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/timer"
)

func NewRoomTimerMgr() *RoomTimerMgr {
	r := new(RoomTimerMgr)
	return r
}

type RoomTimerMgr struct {
	EndTime         *timer.Timer //开局为加入的超时
	ChaHuaTime      *timer.Timer //插花超时
	TimeLimit       int          //时间限制
	CountLimit      int          //局数限制
	TimeOutCard     int          //出牌时间
	TimeOperateCard int          //操作时间
	PlayCount       int          //已玩局数
	MaxPlayCnt      int          //玩家主动设置的最大局数
	CreateTime      int64        //创建时间
}

func (room *RoomTimerMgr) GetTimeOperateCard() int {
	return room.TimeOperateCard
}

func (room *RoomTimerMgr) GetTimeOutCard() int {
	return room.MaxPlayCnt
}

func (room *RoomTimerMgr) GetMaxPayCnt() int {
	return room.MaxPlayCnt
}

func (room *RoomTimerMgr) GetCreatrTime() int64 {
	return room.CreateTime
}

func (room *RoomTimerMgr) GetPlayCount() int {
	return room.PlayCount
}

func (room *RoomTimerMgr) AddPlayCount() {
	room.PlayCount++
}

func (room *RoomTimerMgr) GetCountLimit() int {
	return room.CountLimit
}

func (room *RoomTimerMgr) GetTimeLimit() int {
	return room.TimeLimit
}

//创建房间多久没开始解散房间
func (room *RoomTimerMgr) StartCreatorTimer(Skeleton *module.Skeleton, cb func()) {
	room.EndTime.Stop()
	if room.TimeLimit != 0 {
		room.EndTime = Skeleton.AfterFunc(time.Duration(room.TimeLimit)*time.Second, cb)
	}
}

//开始多久没打完解散房间房间
func (room *RoomTimerMgr) StartPlayingTimer(Skeleton *module.Skeleton, cb func()) {
	room.EndTime.Stop()
	if room.TimeLimit != 0 {
		room.EndTime = Skeleton.AfterFunc(time.Duration(room.TimeLimit)*time.Second, cb)
	}
}
