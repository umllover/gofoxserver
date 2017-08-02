package internal

import (
	"github.com/lovelly/leaf/log"
	"mj/common/consul"
	"mj/gameServer/kindList"

	"mj/common/msg"

	"github.com/lovelly/leaf/nsq/cluster"
)

//新的节点启动了
func serverStart(args []interface{}) {
	svr := args[0].(*consul.CacheInfo)
	log.Debug("%s on line", svr.Csid)
	cluster.AddClient(&cluster.NsqClient{Addr: svr.Host, ServerName: svr.Csid})
}

//节点关闭了
func serverFaild(args []interface{}) {
	svr := args[0].(*consul.CacheInfo)
	log.Debug("%s off line", svr.Csid)
	cluster.RemoveClient(svr.Csid)
}

func LoadRoom(info interface{}) bool {
	retMsg := info.(*msg.L2G_CreatorRoom)
	mod, ok := kindList.GetModByKind(retMsg.KindId)
	if !ok {
		log.Error("at L2G_CreatorRoom not foud kind %d; roomid:%d", retMsg.KindId, retMsg.RoomID)
		return false
	}

	log.Debug("begin CreateRoom.....")
	ok1 := mod.CreateRoom(retMsg)
	if !ok1 {
		log.Debug("at L2G_CreatorRoom mod.CreateRoom falild roomid:%d,", retMsg.RoomID)
		return false
	}

	return true
}
