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
	StartKickoutTimer(Skeleton *module.Skeleton, uid int64, cb func())
	StopOfflineTimer(uid int64)
	StartReplytIimer(uid int64, cb func())
	StopReplytIimer(uid int64)

	GetTimeLimit() int
	GetPlayCount() int
	AddPlayCount()
	GetMaxPayCnt() int
	AddMaxPayCnt(int)
	GetCreatrTime() int64
	GetTimeOutCard() int
	GetTimeOperateCard() int
}

type UserManager interface {
	Sit(*user.User, int) int
	Standup(*user.User) bool
	ForEachUser(fn func(*user.User))
	GetLeaveInfo()map[int64][]int64
	LeaveRoom(*user.User, int) bool
	SetUsetStatus(*user.User, int)
	ReLogin(*user.User, int)
	IsAllReady() bool
	RoomDissume()
	SendUserInfoToSelf(*user.User)
	SendMsgAll(data interface{})
	SendMsgAllNoSelf(selfid int64, data interface{})
	WriteTableScore(source []*msg.TagScoreInfo, usercnt, Type int)
	SendDataToHallUser(chiairID int, data interface{})
	SendMsgToHallServerAll(data interface{})
	ReplyLeave(*user.User, bool, int64, int) bool
	DeleteReply(uid int64)

	GetCurPlayerCnt() int
	GetPayType() int
	IsPublic() bool
	GetMaxPlayerCnt() int
	GetUserInfoByChairId(int) interface{}
	GetUserByChairId(int) *user.User
	GetUserByUid(userId int64) (*user.User, int)
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
