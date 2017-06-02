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

const(
	UserCount = 4
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
	handleRpc("Sitdown", Sitdown, chanrpc.FuncCommon)
	handleRpc("SetGameOption", SetGameOption, chanrpc.FuncCommon)
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

func Sitdown(args []interface{}){
	user := args[1].(*user.User)
	r := getRoom(user.RoomId)
	if r != nil {
		r.ChanRPC.Go("Sitdown", args...)
	}else {
		log.Error("at Sitdown no foud room %v", args[0])
	}
}

func SetGameOption(args []interface{}){

	user := args[1].(*user.User)
	r := getRoom(user.RoomId)
	if r != nil {
		r.ChanRPC.Go("SetGameOption", args...)
	}else {
		log.Error("at SetGameOption no foud room %v", args[0])
	}
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

	agent.WriteMsg(&msg.G2C_ConfigServer{
		TableCount: template.TableCount,
		ChairCount: 4,
		ServerType: template.ServerType,
		ServerRule: template.ServerRule,
	})

	agent.WriteMsg(&msg.G2C_ConfigFinish{})

	r  := room.NewRoom(ChanRPC, recvMsg, template, rid, UserCount, user.Id)
	retMsg.TableID = r.GetRoomId()
	retMsg.DrawCountLimit = r.CountLimit
	retMsg.DrawTimeLimit = r.TimeLimit
	retMsg.Beans = feeTemp.TableFee
	retMsg.RoomCard = user.RoomCard
	user.KindID =  recvMsg.Kind
	user.RoomId = r.GetRoomId()
	addRoom(r)

	agent.WriteMsg(&msg.G2C_UserEnter{
		GameID : user.GameID,						//游戏 I D
		UserID : user.Id,							//用户 I D
		FaceID : user.FaceID,							//头像索引
		CustomID :user.CustomID,						//自定标识
		Gender :user.Gender,							//用户性别
		MemberOrder :user.Accountsinfo.MemberOrder,					//会员等级
		TableID : user.RoomId,							//桌子索引
		ChairID : user.ChairId,							//椅子索引
		UserStatus :user.Status,						//用户状态
		Score :user.Score,								//用户分数
		WinCount : user.WinCount,							//胜利盘数
		LostCount : user.LostCount,						//失败盘数
		DrawCount : user.DrawCount,						//和局盘数
		FleeCount : user.FleeCount,						//逃跑盘数
		Experience : user.Experience,						//用户经验
		NickName: user.NickName,				//昵称
		HeaderUrl :user.HeadImgUrl, 				//头像
	})
}




