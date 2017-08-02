package internal

import (
	"mj/hallServer/conf"
	"mj/hallServer/db"
	"testing"

	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/log"
)

func TestGameRoomID(t *testing.T) {
	//var wg sync.WaitGroup
	//runtime.GOMAXPROCS(4)
	//m := make(map[int]bool)
	////wg.Add(100)
	//for i := 0; i < 100; i++ {
	//	i, err := IncRoomCnt(772954)
	//	if err != nil {
	//		log.Debug("err :%s", err.Error())
	//		return
	//	}
	//	if m[i] {
	//		log.Debug("aaaaaaaaaaaa %d", i)
	//	}
	//	m[i] = true
	//}
	//wg.Wait()
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
	lconf.HeartBeatInterval = conf.HeartBeatInterval
	InitLog()

	db.InitDB(&conf.DBConfig{})
}

func InitLog() {
	logger, err := log.New(conf.Server.LogLevel, "", conf.LogFlag)
	if err != nil {
		panic(err)
	}
	log.Export(logger)
}
