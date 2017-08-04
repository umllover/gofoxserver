package main

import (
	"mj/common/consul"
	"mj/gameServer/Chat"
	"mj/gameServer/center"
	"mj/gameServer/conf"
	"mj/gameServer/db"
	"mj/gameServer/db/model/base"
	"mj/gameServer/gate"
	"mj/gameServer/kindList"

	"flag"

	"fmt"
	"os"

	"mj/gameServer/http_service"

	"mj/gameServer/userHandle"

	"github.com/lovelly/leaf"
	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/module"
)

var version = 0

var printVersion = flag.Bool("version", false, "print version")
var Test = flag.Bool("Test", false, "reload base db")

func main() {
	flag.Parse()
	if *printVersion {
		fmt.Println(" version: ", version)
		os.Exit(0)
	}
	Init()
	http_service.StartHttpServer()
	http_service.StartPrivateServer()
	consul.SetConfig(&conf.ConsulConfig{})
	consul.SetSelfId(conf.ServerName())
	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()
	kindList.Init()

	modules := []module.Module{}
	modules = append(modules, userHandle.UserMgr)
	modules = append(modules, gate.Module)
	modules = append(modules, center.Module)
	modules = append(modules, Chat.Module)
	modules = append(modules, kindList.GetModules()...)
	modules = append(modules, consul.Module)
	leaf.Run(modules...)
}

func Init() {
	conf.Init()
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
	conf.Test = *Test
	leaf.InitLog()
	leaf.OnDestroy = userHandle.KickOutUser
}
