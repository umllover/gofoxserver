package main

import (
	"flag"
	"fmt"
	"mj/common"
	"mj/common/consul"
	. "mj/common/cost"
	"mj/hallServer/center"
	"mj/hallServer/conf"
	"mj/hallServer/db"
	"mj/hallServer/db/model/base"
	"mj/hallServer/game_list"
	"mj/hallServer/gate"
	"mj/hallServer/http_service"
	"mj/hallServer/race_msg"
	"mj/hallServer/shop"
	"mj/hallServer/userHandle"
	"os"

	"mj/hallServer/match_room"

	"mj/hallServer/times_mgr"

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
	game_list.SetTest(*Test)
	if *reloadDB {
		db.NeedReloadBaseDB = true
		log.Debug("need reload base db")
		db.RefreshInTime()
	}

	common.Init()
	http_service.StartHttpServer()
	http_service.StartPrivateServer()
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
		game_list.Module,
		race_msg.Module,
		match_room.Module,
		times_mgr.Module,
		shop.Module,
	)
}

func Init() {
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
	leaf.InitLog()
}
