package internal

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/conf"
	"mj/gameServer/db/model/base"
	"reflect"

	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/gate"
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

	handlerC2S(&msg.C2G_CreateTable{}, CreateTable)
}

//创建桌子
func CreateTable(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_CreateTable)
	agent := args[1].(gate.Agent)
	retCode := 0

	defer func() {
		if retCode != 0 {
			agent.WriteMsg(&msg.G2C_CreateTableFailure{ErrorCode: retCode, DescribeString: "创建房间失败"})
		}
	}()

	mod, ok := GetModByKind(recvMsg.Kind)
	if !ok {
		retCode = NotFoudGameType
		return
	}

	log.Debug("begin CreateRoom.....")
	mod.GetChanRPC().Go("CreateRoom", recvMsg, agent)
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
			svrInfo.FullCount = template.MaxDistributeUser
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
