package main

import (
	"mj/common"
	"mj/common/consul"
	. "mj/common/cost"
	"mj/hallServer/center"
	"mj/hallServer/conf"
	"mj/hallServer/db"
	"mj/hallServer/db/model/base"
	"mj/hallServer/gameList"
	"mj/hallServer/gate"
	"mj/hallServer/userHandle"

	"github.com/lovelly/leaf"
	lconf "github.com/lovelly/leaf/conf"
)

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath
	lconf.ListenAddr = conf.Server.ListenAddr
	lconf.ServerName = conf.ServerName()
	lconf.ConnAddrs = conf.Server.ConnAddrs
	lconf.PendingWriteNum = conf.Server.PendingWriteNum
	lconf.HeartBeatInterval = conf.HeartBeatInterval

	common.Init()
	consul.SetConfig(&conf.ConsulConfig{})
	consul.SetSelfId(lconf.ServerName)
	consul.AddinitiativeSvr(GamePrefix)
	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()
	leaf.Run(
		gate.Module,
		center.Module,
		consul.Module,
		userHandle.UserMgr,
		gameList.Module,
	)
}
