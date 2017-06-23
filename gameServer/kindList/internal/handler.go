package internal

import (
	"mj/common/msg"
	"mj/gameServer/common"
	"mj/gameServer/conf"
	"mj/gameServer/db/model/base"
	"reflect"

	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/log"
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.Register(id, f)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	//rpc
	handleRpc("GetKindList", GetKindList)

}

///// rpc
func GetKindList(args []interface{}) (interface{}, error) {
	ip, port := conf.GetServerAddrAndPort()

	ret := make([]*msg.TagGameServer, 0)
	for kind, v := range modules {
		templates, ok := base.GameServiceOptionCache.GetKey1(kind)
		if !ok {
			continue
		}
		for _, template := range templates {
			svrInfo := &msg.TagGameServer{}
			svrInfo.KindID = kind
			svrInfo.NodeID = conf.Server.NodeId
			svrInfo.SortID = template.SortID
			svrInfo.ServerID = template.ServerID
			svrInfo.ServerPort = port
			svrInfo.ServerType = int64(template.ServerType)
			svrInfo.OnLineCount = int64(v.GetClientCount())
			svrInfo.FullCount = common.TableFullCount
			svrInfo.RestrictScore = int64(template.RestrictScore)
			svrInfo.MinTableScore = int64(template.MinTableScore)
			svrInfo.MinEnterScore = int64(template.MinEnterScore)
			svrInfo.MaxEnterScore = int64(template.MaxEnterScore)
			svrInfo.ServerAddr = ip
			svrInfo.ServerName = template.ServerName
			svrInfo.SurportType = 0
			svrInfo.TableCount = v.GetTableCount()
			ret = append(ret, svrInfo)
		}
	}

	log.Debug("at GetKindList ==== %v", ret)
	return ret, nil
}

