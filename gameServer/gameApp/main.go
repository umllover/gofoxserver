package main

import (
	"mj/common"
	"mj/common/consul"
	"github.com/lovelly/leaf"
	lconf "github.com/lovelly/leaf/conf"
	"mj/gameServer/conf"
	"mj/gameServer/center"
	"mj/gameServer/kindList"
	"mj/gameServer/gate"
	"mj/gameServer/login"
	"github.com/lovelly/leaf/module"
	"mj/gameServer/db/model/base"
	"mj/gameServer/db"
)

func main() {
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

	common.Init()
	consul.SetConfig(&conf.ConsulConfig{})
	consul.SetSelfId(conf.ServerName())
	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()

	modules := []module.Module{center.Module}
	modules = append(modules, gate.Module)
	modules = append(modules, center.Module)
	modules = append(modules, consul.Module)
	modules = append(modules, login.Module)
	modules = append(modules, kindList.GetModules()...)

	leaf.Run(modules...)
}
