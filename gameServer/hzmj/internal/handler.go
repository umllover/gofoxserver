package internal

import (

	"mj/common/msg"
	"mj/gameServer/hzmj/room"
	"reflect"
	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/gate"
	//"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	. "mj/common/cost"
	"mj/gameServer/common"
	"mj/gameServer/db/model/base"
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
	// c 2 s
	handlerC2S(&msg.C2G_HZOutCard{}, HZOutCard)

	// rpc
	handleRpc("DelRoom", DelRoom, chanrpc.FuncCommon)
	handleRpc("CreateRoom", CreaterRoom, chanrpc.FuncCommon)
	handleRpc("SrarchTableInfo", SrarchTableInfo, chanrpc.FuncCommon)
}



func HZOutCard (args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.ChanRPC.Go("OutCard", args[0])
	}
}


//////////////// rcp ///////////////////
func DelRoom(args []interface{}){
	id := args[0].(int)
	delRoom(id)
}

func CreaterRoom(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_CreateTable)
	retMsg := &msg.G2C_CreateTableSucess{}
	agent := args[1].(gate.Agent)
	retCode := 0
	defer func() {
		if retCode == 0 {
			agent.WriteMsg(retMsg)
		}else {
			agent.WriteMsg(&msg.G2C_CreateTableFailure{ErrorCode:retCode, DescribeString:"创建房间失败"})
		}
	}()

	user := agent.UserData().(*user.User)
	_ = user
	if wTableCount > 10000 {
		retCode = RoomFull
		return
	}

	if recvMsg.Kind !=  common.KIND_TYPE_HZMJ {
		retCode = CreateParamError
		return
	}

	template, ok := base.GameServiceOptionCache.Get(recvMsg.Kind, recvMsg.ServerId)
	if !ok {
		retCode = NoFoudTemplate
		return
	}

	if template.CardOrBean == 0 { //消耗游戏豆

	}else if  template.CardOrBean == 1 { //消耗房卡

	}else{
		retCode = ConfigError
		return
	}


	r  := room.NewRoom(ChanRPC, recvMsg, template)
	addRoom(r)
}

func SrarchTableInfo(args []interface{}) {
	retMsg := &msg.G2C_SearchResult{}
	agent := args[1].(gate.Agent)
	defer func(){
		agent.WriteMsg(retMsg)
	}()

	user := agent.UserData().(*user.User)
	r := getRoom(user.RoomId)
	if r == nil {
		return
	}

	retMsg.TableID = r.GetRoomId()
	retMsg.ServerID = r.ServerId
	return
}


