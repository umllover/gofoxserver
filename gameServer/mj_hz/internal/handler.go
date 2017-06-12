package internal

import (
	"mj/common/msg"
	"mj/gameServer/mj_hz/room"
	"reflect"

	"github.com/lovelly/leaf/gate"
	//"mj/gameServer/db/model/base"
	. "mj/common/cost"
	"mj/gameServer/common"
	"mj/gameServer/db/model/base"
	"mj/gameServer/idGenerate"
	"mj/gameServer/user"

	"mj/common/msg/mj_hz_msg"

	"github.com/lovelly/leaf/log"
)

const (
	UserCount = 4
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}) {
	ChanRPC.Register(id, f)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	// c 2 s
	handlerC2S(&mj_hz_msg.C2G_HZMJ_HZOutCard{}, HZOutCard)
	handlerC2S(&mj_hz_msg.C2G_HZMJ_OperateCard{}, OperateCard)
	// rpc
	handleRpc("DelRoom", DelRoom)
	handleRpc("CreateRoom", CreaterRoom)
	handleRpc("Sitdown", Sitdown)
	handleRpc("SetGameOption", SetGameOption)
	handleRpc("UserStandup", UserStandup)
	handleRpc("GetUserChairInfo", GetUserChairInfo)
	handleRpc("UserReady", UserReady)
}

func HZOutCard(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.ChanRPC.Go("OutCard", args[0], user)
	}
}

func OperateCard(args []interface{}) {
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	r := getRoom(user.RoomId)
	if r != nil {
		r.ChanRPC.Go("OperateCard", args[0], user)
	}
}

//////////////// rcp ///////////////////
func DelRoom(args []interface{}) {
	id := args[0].(int)
	delRoom(id)
}

func Sitdown(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_UserSitdown)
	user := args[1].(*user.User)
	r := getRoom(user.RoomId)
	if r == nil {
		r = getRoom(recvMsg.TableID)
	}
	if r != nil {
		r.ChanRPC.Go("Sitdown", args...)
	} else {
		log.Error("at Sitdown no foud room %v", args[0])
	}
}

//只读信息 不涉及竞争， 在这里就处理了， 不投递进房间
func GetUserChairInfo(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_REQUserChairInfo)
	agent := args[1].(gate.Agent)
	user, ok := agent.UserData().(*user.User)
	if !ok {
		log.Error("at GerUserInfo user not logon")
		return
	}

	r := getRoom(user.RoomId)
	if r == nil {
		log.Error("at GetUserChairInfo no foud room %v, userId:%d", args[0], user.Id)
		return
	}

	tagUser := r.GetUserByChairId(recvMsg.ChairID)
	if tagUser == nil {
		log.Error("at GetUserChairInfo no foud tagUser %v, userId:%d", args[0], user.Id)
		return
	}

	agent.WriteMsg(&msg.G2C_UserEnter{
		GameID:      tagUser.GameID,                   //游戏 I D
		UserID:      tagUser.Id,                       //用户 I D
		FaceID:      tagUser.FaceID,                   //头像索引
		CustomID:    tagUser.CustomID,                 //自定标识
		Gender:      tagUser.Gender,                   //用户性别
		MemberOrder: tagUser.Accountsinfo.MemberOrder, //会员等级
		TableID:     tagUser.RoomId,                   //桌子索引
		ChairID:     tagUser.ChairId,                  //椅子索引
		UserStatus:  tagUser.Status,                   //用户状态
		Score:       tagUser.Score,                    //用户分数
		WinCount:    tagUser.WinCount,                 //胜利盘数
		LostCount:   tagUser.LostCount,                //失败盘数
		DrawCount:   tagUser.DrawCount,                //和局盘数
		FleeCount:   tagUser.FleeCount,                //逃跑盘数
		Experience:  tagUser.Experience,               //用户经验
		NickName:    tagUser.NickName,                 //昵称
		HeaderUrl:   tagUser.HeadImgUrl,               //头像
	})
}

func SetGameOption(args []interface{}) {

	user := args[1].(*user.User)
	r := getRoom(user.RoomId)
	if r != nil {
		r.ChanRPC.Go("SetGameOption", args...)
	} else {
		log.Error("at SetGameOption no foud room %v", args[0])
	}
}

func UserStandup(args []interface{}) {
	user := args[1].(*user.User)
	r := getRoom(user.RoomId)
	if r != nil {
		r.ChanRPC.Go("UserStandup", args...)
	} else {
		log.Error("at UserStandup no foud room %v", args[0])
	}
}

func UserReady(args []interface{}) {
	user := args[1].(*user.User)
	r := getRoom(user.RoomId)
	if r != nil {
		r.ChanRPC.Go("UserReady", args...)
	} else {
		log.Error("at UserReady no foud room %v", args[0])
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
		} else {
			agent.WriteMsg(&msg.G2C_CreateTableFailure{ErrorCode: retCode, DescribeString: "创建房间失败"})
		}
	}()

	user := agent.UserData().(*user.User)
	if wTableCount > 10000 {
		retCode = RoomFull
		return
	}

	if recvMsg.Kind != common.KIND_TYPE_HZMJ {
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
	} else if template.CardOrBean == 1 { //消耗房卡
		if user.RoomCard < template.FeeBeanOrRoomCard {
			retCode = NotEnoughFee
			return
		}
	} else {
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

	r := room.NewRoom(ChanRPC, recvMsg, template, rid, UserCount, user.Id)
	if recvMsg.DrawTimeLimit == 0 {
		r.TimeLimit = feeTemp.DrawTimeLimit
		r.CountLimit = feeTemp.DrawCountLimit
		r.Source = feeTemp.IniScore
	}
	r.TimeOutCard = template.OutCardTime
	r.TimeOperateCard = template.OperateCardTime
	retMsg.TableID = r.GetRoomId()
	retMsg.DrawCountLimit = r.CountLimit
	retMsg.DrawTimeLimit = r.TimeLimit
	retMsg.Beans = feeTemp.TableFee
	retMsg.RoomCard = user.RoomCard
	user.KindID = recvMsg.Kind
	user.RoomId = r.GetRoomId()
	addRoom(r)
}
