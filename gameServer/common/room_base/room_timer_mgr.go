package room_base

import (
	. "mj/common/cost"
	"mj/gameServer/common"
	"mj/gameServer/db/model/base"
	"time"

	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/timer"
)

func NewRoomTimerMgr(roomId, UserCnt int, Temp *base.GameServiceOption) common.TimerManager {
	r := new(RoomTimerMgr)
	return r
}

type RoomTimerMgr struct {
	EndTime         *timer.Timer //开局为加入的超时
	TimeLimit       int          //时间限制
	CountLimit      int          //局数限制
	TimeOutCard     int          //出牌时间
	TimeOperateCard int          //操作时间
	PlayCount       int          //已玩局数
	MaxPlayCnt      int          //玩家主动设置的最大局数
	CreateTime      int64        //创建时间
}

func (room *RoomTimerMgr) GetPlayCount() int {
	return room.PlayCount
}

func (room *RoomTimerMgr) AddPlayCount() {
	room.PlayCount++
}

//创建房间多久没开始解散房间
func (room *RoomTimerMgr) StartCreatorTimer(skeleton *module.Skeleton, cb func()) {
	room.EndTime.Stop()
	if room.TimeLimit != 0 {
		room.EndTime = skeleton.AfterFunc(time.Duration(room.TimeLimit)*time.Second,cb)
	}
}

//开始多久没打完解散房间房间
func (room *RoomTimerMgr) StartPlayingTimer(skeleton *module.Skeleton, cb func()) {
	room.EndTime.Stop()
	if room.TimeLimit != 0 {
		room.EndTime = skeleton.AfterFunc(time.Duration(room.TimeLimit)*time.Second,cb)
	}
}
