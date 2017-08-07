package main

import (
	"flag"
	"fmt"
	"mj/common/consul"
	"mj/hallServer/center"
	"mj/hallServer/conf"
	"mj/hallServer/db"
	"mj/hallServer/db/model/base"
	"mj/hallServer/game_list"
	"mj/hallServer/gate"
	"mj/hallServer/http_service"
	"mj/hallServer/match_room"
	"mj/hallServer/race_msg"
	"mj/hallServer/times_mgr"
	"mj/hallServer/user"
	"mj/hallServer/userHandle"
	"os"

	"github.com/lovelly/leaf"
	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/log"
)

var version = 0

var printVersion = flag.Bool("version", false, "print version")
var reloadDB = flag.Bool("reload", false, "reload base db")
var Test = flag.Bool("Test", false, "reload base db")

func main() {
	flag.Parse()
	if *printVersion {
		fmt.Println(" version: ", version)
		os.Exit(0)
	}
	Init()
	log.Debug("enter hallApp main")
	if *reloadDB {
		db.NeedReloadBaseDB = true
		log.Debug("need reload base db")
		db.RefreshInTime()
	}

	http_service.StartHttpServer()
	http_service.StartPrivateServer()
	consul.SetConfig(&conf.ConsulConfig{})
	consul.SetSelfId(lconf.ServerName)
	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()
	user.LoadIncId()
	leaf.Run(
		userHandle.UserMgr,
		center.Module,
		consul.Module,
		game_list.Module,
		race_msg.Module,
		match_room.Module,
		times_mgr.Module,
		gate.Module,
	)
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
	conf.Test = *Test
	leaf.InitLog()
	leaf.OnDestroy = func() {
		conf.Shutdown = true
		lconf.Shutdown = true
		userHandle.KickOutUser()
	}
}
