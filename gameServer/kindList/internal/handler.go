package internal

import (
	"mj/common/msg"
	"github.com/lovelly/leaf/cluster"
	"reflect"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"mj/gameServer/db/model/base"
	"mj/gameServer/conf"
	"github.com/lovelly/leaf/gate"
	. "mj/common/cost"
	"mj/gameServer/user"
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
	//rpc
	handleRpc("GetKindList", GetKindList, chanrpc.FuncCommon)

	//c2s
	handlerC2S(&msg.C2G_GR_UserChairReq{}, UserChairReq)
	handlerC2S(&msg.C2G_CreateTable{}, CreateTable)
	handlerC2S(&msg.C2G_UserSitdown{}, UserSitdown)
	handlerC2S(&msg.C2G_GameOption{}, SetGameOption)
	handlerC2S(&msg.C2G_UserStandup{}, UserStandup)
	handlerC2S(&msg.C2G_REQUserChairInfo{}, GetUserChairInfo)
	handlerC2S(&msg.C2G_UserReady{}, UserReady)
}

//客户端请求更换椅子
func UserChairReq(args []interface{}) {


}

func GetUserChairInfo (args []interface{}) {
	agent := args[1].(gate.Agent)
	user  := agent.UserData().(*user.User)
	mod, ok := GetModByKind(user.KindID)
	if !ok {
		log.Error("at GetUserChairInfo not foud module")
		return
	}

	mod.GetChanRPC().Go("GetUserChairInfo",  args[0], user)
}

//起立
func UserStandup(args []interface{}) {
	agent := args[1].(gate.Agent)
	user  := agent.UserData().(*user.User)
	mod, ok := GetModByKind(user.KindID)
	if !ok {
		log.Error("at UserStandup not foud module")
		return
	}

	mod.GetChanRPC().Go("UserStandup",  args[0], user)

}

//创建桌子
func CreateTable(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_CreateTable)
	agent := args[1].(gate.Agent)
	retCode := 0

	defer func() {
		if retCode != 0 {
			agent.WriteMsg(&msg.G2C_CreateTableFailure{ErrorCode:retCode, DescribeString:"创建房间失败"})
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

func UserSitdown(args []interface{}) {
	agent := args[1].(gate.Agent)
	user  := agent.UserData().(*user.User)
	mod, ok := GetModByKind(user.KindID)
	if !ok {
		log.Error("at UserSitdown not foud module")
		return
	}

	mod.GetChanRPC().Go("Sitdown",  args[0], user)
}


func SetGameOption(args []interface{}) {
	agent := args[1].(gate.Agent)
	user  := agent.UserData().(*user.User)
	mod, ok := GetModByKind(user.KindID)
	if !ok {
		log.Error("at UserSitdown not foud module")
		return
	}

	mod.GetChanRPC().Go("SetGameOption",  args[0], user)
}

func UserReady(args []interface{}) {
	agent := args[1].(gate.Agent)
	user  := agent.UserData().(*user.User)
	mod, ok := GetModByKind(user.KindID)
	if !ok {
		log.Error("at UserReady not foud module")
		return
	}

	mod.GetChanRPC().Go("UserReady",  args[0], user)

}


///// rpc
func GetKindList(args []interface{})(interface{}, error){
	ip, port := conf.GetServerAddrAndPort()

	ret := make([]*msg.TagGameServer, 0)
	for kind, v := range modules {
		templates, ok := base.GameServiceOptionCache.GetKey1(kind)
		if !ok {
			continue
		}
		for _, template := range templates{
			svrInfo := &msg.TagGameServer{}
			svrInfo.KindID = kind
			svrInfo.NodeID = template.NodeID
			svrInfo.SortID  = template.SortID
			svrInfo.ServerID = template.ServerID
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
	}

	log.Debug("at GetKindList ==== %v", ret)
	return  ret, nil
}
