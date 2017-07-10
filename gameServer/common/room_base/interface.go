package room_base

import (
	"mj/common/msg"
	"mj/gameServer/user"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/timer"
)

type Module interface {
	GetChanRPC() *chanrpc.Server
	GetClientCount() int
	GetTableCount() int
	OnDestroy()
	OnInit()
	Run(chan bool)
	CreateRoom(args ...interface{}) bool
}

type TimerManager interface {
	StartCreatorTimer(Skeleton *module.Skeleton, cb func())
	StartPlayingTimer(Skeleton *module.Skeleton, cb func())
	StartKickoutTimer(Skeleton *module.Skeleton, uid int, cb func())
	StopOfflineTimer(uid int)

	GetTimeLimit() int
	GetPlayCount() int
	AddPlayCount()
	GetMaxPayCnt() int
	GetCreatrTime() int64
	GetTimeOutCard() int
	GetTimeOperateCard() int
}

type UserManager interface {
	Sit(*user.User, int) int
	Standup(*user.User) bool
	ForEachUser(fn func(*user.User))
	LeaveRoom(*user.User, int) bool
	SetUsetStatus(*user.User, int)
	ReLogin(*user.User, int)
	IsAllReady() bool
	RoomDissume()
	SendUserInfoToSelf(*user.User)
	SendMsgAll(data interface{})
	SendMsgAllNoSelf(selfid int, data interface{})
	WriteTableScore(source []*msg.TagScoreInfo, usercnt, Type int)
	SendDataToHallUser(chiairID int, funcName string, data interface{})
	SendMsgToHallServerAll(funcName string, data interface{})
	SendCloseRoomToHall(data interface{})

	GetCurPlayerCnt() int
	GetMaxPlayerCnt() int
	GetUserInfoByChairId(int) interface{}
	GetUserByChairId(int) *user.User
	GetUserByUid(userId int) (*user.User, int)
	SetUsetTrustee(chairId int, isTruste bool)
	IsTrustee(chairId int) bool
	GetTrustees() []bool
}

type BaseManager interface {
	Destroy(int)
	RoomRun(int)
	GetSkeleton() *module.Skeleton
	AfterFunc(d time.Duration, cb func()) *timer.Timer
	GetChanRPC() *chanrpc.Server
}
