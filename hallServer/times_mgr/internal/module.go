package internal

import (
	"mj/hallServer/base"
	"mj/hallServer/user"
	"mj/hallServer/userHandle"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/timer"
)

type MachPlayer struct {
	ch      *chanrpc.Server
	EndTime int64
	Uid     int
}

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
)

type TimesModule struct {
	*module.Skeleton
}

func (m *TimesModule) OnInit() {
	m.Skeleton = skeleton
	corn, err := timer.NewCronExpr("0 0 12 * * *")
	if err != nil {
		log.Fatal("at TimesModule OnInit NewCronExpr error:%s", err.Error())
	}
	m.Skeleton.CronFunc(corn, m.ClearDayTimes)

	wcorn, werr := timer.NewCronExpr("0 0 12 * * 0")
	if werr != nil {
		log.Fatal("at TimesModule OnInit NewCronExpr error:%s", err.Error())
	}
	m.Skeleton.CronFunc(wcorn, m.CliearWeekTimes)
}

func (m *TimesModule) OnDestroy() {

}

func (m *TimesModule) ClearDayTimes() {
	defer func() {
		corn, _ := timer.NewCronExpr("0 0 12 * * *")
		m.Skeleton.CronFunc(corn, m.ClearDayTimes)
	}()
	userHandle.ForEachUser(func(u *user.User) {
		u.ClearDayTimes()
	})

}

func (m *TimesModule) CliearWeekTimes() {
	defer func() {
		wcorn, _ := timer.NewCronExpr("0 0 12 * * 0")
		m.Skeleton.CronFunc(wcorn, m.CliearWeekTimes)
	}()

	userHandle.ForEachUser(func(u *user.User) {
		u.ClearWeekTimes()
	})
}
