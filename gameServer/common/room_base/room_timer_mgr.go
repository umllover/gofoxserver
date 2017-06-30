package room_base

import (
	"mj/gameServer/db/model/base"
	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/timer"
)

func NewRoomTimerMgr(payCnt int, temp *base.GameServiceOption) *RoomTimerMgr {
	r := new(RoomTimerMgr)
	r.CreateTime = time.Now().Unix()
	if payCnt < temp.PlayTurnCount {
		r.MaxPlayCnt = payCnt
	} else {
		r.MaxPlayCnt = temp.PlayTurnCount
	}

	r.TimeLimit = temp.TimeAfterBeginTime
	r.TimeOutCard = temp.OutCardTime
	r.TimeOperateCard = temp.OperateCardTime
	r.KickOut = make(map[int]*timer.Timer)
	r.OfflineKickotTime = temp.TimeOffLineCount
	r.TimeLimitNotBegin = temp.TimeNotBeginGame
	return r
}

type RoomTimerMgr struct {
	EndTime           *timer.Timer         //开局为加入的超时
	ChaHuaTime        *timer.Timer         //插花超时
	KickOut           map[int]*timer.Timer //即将被踢出的超时定时器
	TimeLimit         int                  //一局玩多久
	TimeLimitNotBegin int                  //创建房间后多久没开始
	TimeOutCard       int                  //出牌时间
	TimeOperateCard   int                  //操作时间
	PlayCount         int                  //已玩局数
	MaxPlayCnt        int                  //玩家主动设置的最大局数
	CreateTime        int64                //创建时间
	OfflineKickotTime int                  //离线超时时间
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

func (room *RoomTimerMgr) GetTimeLimit() int {
	return room.TimeLimit
}

//创建房间多久没开始解散房间
func (room *RoomTimerMgr) StartCreatorTimer(Skeleton *module.Skeleton, cb func()) {
	log.Debug("StartCreatorTimer 111111111 %d", room.TimeLimitNotBegin)
	if room.EndTime != nil {
		room.EndTime.Stop()
	}

	if room.TimeLimitNotBegin != 0 {
		log.Debug("StartCreatorTimer %d", room.TimeLimitNotBegin)
		room.EndTime = Skeleton.AfterFunc(time.Duration(room.TimeLimitNotBegin)*time.Second, cb)
	}
}

//开始多久没打完解散房间房间
func (room *RoomTimerMgr) StartPlayingTimer(Skeleton *module.Skeleton, cb func()) {
	if room.EndTime != nil {
		room.EndTime.Stop()
	}
	if room.TimeLimit != 0 {
		room.EndTime = Skeleton.AfterFunc(time.Duration(room.TimeLimit)*time.Second, cb)
	}
}

//玩家离线超时
func (room *RoomTimerMgr) StartKickoutTimer(Skeleton *module.Skeleton, uid int, cb func()) {
	if room.OfflineKickotTime != 0 {
		log.Debug("StartKickoutTimer : %d ", room.OfflineKickotTime)
		room.KickOut[uid] = Skeleton.AfterFunc(time.Duration(room.OfflineKickotTime)*time.Second, cb)
	} else {
		cb()
	}
}

//关闭离线超时
func (room *RoomTimerMgr) StopOfflineTimer(uid int) {
	timer, ok := room.KickOut[uid]
	if ok {
		timer.Stop()
	}
}
