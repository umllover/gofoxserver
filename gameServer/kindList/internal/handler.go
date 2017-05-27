package internal

import (
	"mj/common/msg"
	"github.com/lovelly/leaf/cluster"
	"reflect"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"mj/gameServer/db/model/base"
	"mj/gameServer/conf"
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}, fType int) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.RegisterFromType(id, f, fType)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init(){
	handleRpc("GetKindList", GetKindList, chanrpc.FuncCommon)
}

func GetKindList(args []interface{})(interface{}, error){
	log.Debug("at GetKindList ==== ")
	ip, port := conf.GetServerAddrAndPort()

	ret := make([]*msg.TagGameServer, 0)
	for kind, v := range modules {
		template, ok := base.GameServiceOptionCache.Get(kind)
		if !ok {
			continue
		}
		svrInfo := &msg.TagGameServer{}
		svrInfo.KindID = kind
		svrInfo.NodeID = template.NodeID
		svrInfo.SortID  = template.SortID
		svrInfo.ServerID = conf.Server.ServerId
		svrInfo.ServerPort =port
		svrInfo.ServerType = int64(template.ServerType)
		svrInfo.OnLineCount = int64(v.GetClientCount())
		svrInfo.FullCount = template.MaxDistributeUser
		svrInfo.RestrictScore = int64(template.RestrictScore)
		svrInfo.MinTableScore =  int64(template.MinTableScore)
		svrInfo.MinEnterScore =  int64(template.MinEnterScore)
		svrInfo.MaxEnterScore =  int64(template.MaxEnterScore)
		svrInfo.ServerAddr = ip
		svrInfo.ServerName = template.ServerName
		svrInfo.SurportType = 0
		svrInfo.TableCount = v.GetTableCount()
		ret = append(ret, svrInfo)
	}

	log.Debug("at GetKindList ==== %v", ret)
	return  ret, nil
}
