package room_base

import (
	. "mj/common/cost"
	"mj/gameServer/common"
	"mj/gameServer/db/model/base"
	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/timer"
)

func NewRoomTimerMgr(playCount int, temp *base.GameServiceOption, Skeleton *module.Skeleton) *RoomTimerMgr {
	r := new(RoomTimerMgr)
	r.CreateTime = time.Now().Unix()

	//TODO 这里的playCount不能完全信赖客户端吧？
	if playCount < temp.PlayTurnCount {
		r.RoomPlayCnt = playCount
	} else {
		r.RoomPlayCnt = temp.PlayTurnCount
	}
	r.MaxPlayCnt = r.RoomPlayCnt

	r.TimeLimit = temp.TimeAfterBeginTime
	r.TimeOutCard = temp.OutCardTime
	r.TimeOperateCard = temp.OperateCardTime
	r.KickOut = make(map[int64]*timer.Timer)
	r.LeaveTimer = make(map[int64]*timer.Timer)
	r.OfflineKickotTime = temp.TimeOffLineCount
	r.TimeLimitNotBegin = temp.TimeNotBeginGame
	r.Skeleton = Skeleton
	return r
}

type RoomTimerMgr struct {
	Skeleton          *module.Skeleton
	EndTime           *timer.Timer           //开局为加入的超时
	ChaHuaTime        *timer.Timer           //插花超时
	KickOut           map[int64]*timer.Timer //即将被踢出的超时定时器
	LeaveTimer        map[int64]*timer.Timer //请求离开房间超时
	TimeLimit         int                    //一局玩多久
	TimeLimitNotBegin int                    //创建房间后多久没开始
	TimeOutCard       int                    //出牌时间
	TimeOperateCard   int                    //操作时间
	PlayCount         int                    //已玩局数
	MaxPlayCnt        int                    //最大可玩局数(续费后可能会变)
	RoomPlayCnt       int                    //房间局数配置(房间创建后固定不变)
	CreateTime        int64                  //创建时间
	OfflineKickotTime int                    //离线超时时间
}

func (room *RoomTimerMgr) GetTimeOperateCard() int {
	return room.TimeOperateCard
}

func (room *RoomTimerMgr) AddMaxPlayCnt(cnt int) {
	room.MaxPlayCnt += cnt
}

func (room *RoomTimerMgr) GetTimeOutCard() int {
	return room.TimeOutCard
}

func (room *RoomTimerMgr) GetMaxPlayCnt() int {
	return room.MaxPlayCnt
}

func (room *RoomTimerMgr) GetRoomPlayCnt() int {
	return room.RoomPlayCnt
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

func (room *RoomTimerMgr) ResetPlayCount() {
	room.PlayCount = 0
}

func (room *RoomTimerMgr) GetTimeLimit() int {
	return room.TimeLimit
}

//创建房间多久没开始解散房间
func (room *RoomTimerMgr) StartCreatorTimer(cb func()) {
	log.Debug("StartCreatorTimer %d", room.TimeLimitNotBegin)
	if room.EndTime != nil {
		room.EndTime.Stop()
	}

	if room.TimeLimitNotBegin != 0 {
		log.Debug("StartCreatorTimer %d not begin limit", room.TimeLimitNotBegin)
		room.EndTime = room.Skeleton.AfterFunc(time.Duration(room.TimeLimitNotBegin)*time.Second, cb)
	}
}

//停止创建房间没开始的定时器
func (room *RoomTimerMgr) StopCreatorTimer() {
	if room.EndTime != nil {
		room.EndTime.Stop()
	}
}

//玩家离线超时
func (room *RoomTimerMgr) StartKickoutTimer(uid int64, cb func()) {
	if room.OfflineKickotTime != 0 {
		log.Debug("StartKickoutTimer : %d ", room.OfflineKickotTime)
		room.KickOut[uid] = room.Skeleton.AfterFunc(time.Duration(room.OfflineKickotTime)*time.Second, cb)
	} else {
		cb()
	}
}

//关闭离线超时
func (room *RoomTimerMgr) StopOfflineTimer(uid int64) {
	timer, ok := room.KickOut[uid]
	if ok {
		timer.Stop()
	}
}

//请求离开超时
func (room *RoomTimerMgr) StartReplytIimer(uid int64, cb func()) {
	ReqLeaveTimer := common.GetGlobalVarInt(LeaveRoomTimer)
	if ReqLeaveTimer != 0 {
		log.Debug("StartKickoutTimer : %d ", ReqLeaveTimer)
		room.LeaveTimer[uid] = room.Skeleton.AfterFunc(time.Duration(ReqLeaveTimer)*time.Second, cb)
	} else {
		cb()
	}
}

func (room *RoomTimerMgr) StopReplytIimer(uid int64) {
	log.Debug("at StopReplytIimer ... ")
	t := room.LeaveTimer[uid]
	if t != nil {
		t.Stop()
	}
}
