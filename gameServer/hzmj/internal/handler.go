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
	"mj/gameServer/idGenerate"
	"github.com/name5566/leaf/log"
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

	feeTemp, ok1 := base.PersonalTableFeeCache.Get(recvMsg.ServerId, recvMsg.Kind, recvMsg.DrawCountLimit, recvMsg.DrawTimeLimit)
	if !ok1 {
		log.Error("not foud PersonalTableFeeCache")
		retCode = NoFoudTemplate
		return
	}

	if template.CardOrBean == 0 { //消耗游戏豆
		if user.RoomCard < feeTemp.TableFee {
			retCode = NotEnoughFee
			return
		}
	}else if  template.CardOrBean == 1 { //消耗房卡
		if user.RoomCard < template.FeeBeanOrRoomCard {
			retCode = NotEnoughFee
			return
		}
	}else{
		retCode = ConfigError
		return
	}

	rid, iok := idGenerate.GetRoomId(user.Id)
	if !iok {
		retCode = RandRoomIdError
		return
	}

	if recvMsg.CellScore > template.CellScore {
		retCode = MaxSoucrce
		return
	}

	r  := room.NewRoom(ChanRPC, recvMsg, template, rid)
	retMsg.TableID = r.GetRoomId()
	retMsg.DrawCountLimit = r.CountLimit
	retMsg.DrawTimeLimit = r.TimeLimit
	retMsg.Beans = feeTemp.TableFee
	retMsg.RoomCard = user.RoomCard
	addRoom(r)
}




