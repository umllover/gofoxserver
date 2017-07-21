package internal

import (
	"mj/common/cost"
	"mj/hallServer/base"
	"mj/hallServer/conf"
	"mj/hallServer/user"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/nsq/cluster"
)

var (
	skeleton   = base.NewSkeleton()
	ChanRPC    = skeleton.ChanRPCServer
	Users      = make(map[int64]*chanrpc.Server) //本服玩家
	OtherUsers = make(map[int64]string)          //其他服登录的玩家  key is uid， values is NodeId
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	cfg := &cluster.Cluster_config{
		LogLv:              "Error",
		Channel:            conf.ServerNsqCahnnel(),
		Csmtopics:          []string{cost.HallPrefix, conf.ServerName()}, //需要订阅的主题
		CsmNsqdAddrs:       conf.Server.NsqdAddrs,
		CsmNsqLookupdAddrs: conf.Server.NsqLookupdAddrs,
		PdrNsqdAddr:        conf.Server.PdrNsqdAddr, //生产者需要连接的nsqd地址
		SelfName:           conf.ServerName(),
	}

	cluster.Start(cfg)

	user.CenterRpc = ChanRPC
}

func (m *Module) OnDestroy() {
	cluster.Stop()
}
