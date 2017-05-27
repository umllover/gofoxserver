package internal

import (
	"mj/common/msg"
	"reflect"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/gate"
)

var (
	gameLists = make(map[int]map[int]*msg.TagGameServer) //k1 is kind k2 is server Id
	SvrvetType = make(map[int]map[int]struct{}) //key sverId v KingId
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

func init() {

	handleRpc("sendGameList", sendGameList, chanrpc.FuncCommon)
	handleRpc("updateGameInfo", updateGameInfo, chanrpc.FuncCommon)
	handleRpc("delGameList", delGameList, chanrpc.FuncCommon)

	handleRpc("NewServerAgent", NewServerAgent, chanrpc.FuncCommon)
	handleRpc("CloseServerAgent", CloseServerAgent, chanrpc.FuncCommon)
}

func sendGameList(args []interface{}){
	agent := args[0].(gate.Agent)
	list := make(msg.L2C_ServerList, 0)
	for _, v := range  gameLists {
		for _, v1 := range  v {
			list = append(list, v1)
		}
	}
	agent.WriteMsg(&list)
	finish := &msg.L2C_ServerListFinish{}
	agent.WriteMsg(finish)
}

func updateGameInfo(args []interface{}){

}

func NewServerAgent(args []interface{}){
	serverName := args[0].(string)
 	log.Debug("at NewServerAgent :%s", serverName)
	cluster.AsynCall(serverName,skeleton.GetChanAsynRet(), "GetKindList", func(data interface{}, err error) {
		if err != nil {
			log.Debug("GetKindList error:%s",err.Error() )
			return
		}

		ret := data.([]*msg.TagGameServer)

		for _, v := range ret {
			addGameList(v)
			log.Debug("data %v", v)
		}
	} )
}

func CloseServerAgent(args []interface{}) {
	log.Debug("at CloseServerAgent")
}


///////////////// help
func delGameList(args []interface{}){
	svrId := args[0].(int)
	typeInfo := SvrvetType[svrId]
	for kind, _ := range typeInfo {
		gminfo, ok := gameLists[kind]
		if ok {
			delete(gminfo, svrId)
		}
	}
}

func addGameList(v *msg.TagGameServer){
	gminfo, ok := gameLists[v.KindID]
	if !ok {
		gminfo = make(map[int]*msg.TagGameServer)
		gameLists[v.KindID] = gminfo
	}

	gminfo[v.ServerID] = v

	typeInfo, ok1 := SvrvetType[v.ServerID]
	if !ok1 {
		typeInfo = make(map[int]struct{})
		SvrvetType[v.ServerID] = typeInfo
	}
	typeInfo[v.KindID] = struct {}{}
}

