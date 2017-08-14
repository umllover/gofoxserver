package internal

import (
	"mj/hallServer/common"
	"mj/hallServer/conf"
	"mj/hallServer/db"
	"mj/hallServer/db/model/base"
	"mj/hallServer/user"
	"testing"

	lconf "github.com/lovelly/leaf/conf"

	"github.com/lovelly/leaf/log"
)

func TestUserTime(t *testing.T) {
	player := user.NewUser(110)
	player.SetTimes(common.ActivityBindPhome, 0)
}

func init() {
	conf.Init()
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ServerName = conf.ServerName()
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath
	lconf.ListenAddr = conf.Server.ListenAddr
	lconf.ConnAddrs = conf.Server.ConnAddrs
	lconf.PendingWriteNum = conf.Server.PendingWriteNum
	InitLog()

	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()
}

func InitLog() {
	logger, err := log.New(conf.Server.LogLevel, "", conf.LogFlag)
	if err != nil {
		panic(err)
	}
	log.Export(logger)
}
