package internal

import (
	"mj/hallServer/conf"
	"mj/hallServer/db"
	"testing"

	"sync"

	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/log"
)

var Wg sync.WaitGroup

func TestSendMsgToAll(t *testing.T) {
	//startTimer(0)
	ReciveGMMsg(5, 5, "dfsd")
	Wg.Wait()
}

func init() {
	Wg.Add(1)
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
	log.Debug("数据库初始化完毕")
}

func InitLog() {
	logger, err := log.New(conf.Server.LogLevel, "", conf.LogFlag)
	if err != nil {
		panic(err)
	}
	log.Export(logger)
}
