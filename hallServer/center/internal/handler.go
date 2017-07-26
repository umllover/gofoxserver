package internal

import (
	"errors"
	"mj/common/consul"
	"mj/common/cost"
	"mj/common/msg"
	"mj/common/register"
	"mj/hallServer/conf"
	"mj/hallServer/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
)

var (
	GamelistRpc *chanrpc.Server
)

func init() {
	reg := register.NewRegister(ChanRPC)
	reg.RegisterRpc("SelfNodeAddPlayer", SelfNodeAddPlayer)
	reg.RegisterRpc("SelfNodeDelPlayer", SelfNodeDelPlayer)
	reg.RegisterRpc("SendMsgToSelfNotdeUser", SendMsgToSelfNotdeUser)
	reg.RegisterRpc("SendMsgToHallUser", SendMsgToHallUser)

	reg.RegisterRpc("ServerFaild", serverFaild)
	reg.RegisterRpc("ServerStart", serverStart)

	reg.RegisterS2S(&msg.S2S_GetPlayerInfo{}, GetPlayerInfo)
	reg.RegisterS2S(&msg.S2S_NotifyOtherNodeLogin{}, NotifyOtherNodeLogin)
	reg.RegisterS2S(&msg.S2S_NotifyOtherNodelogout{}, NotifyOtherNodelogout)
	reg.RegisterS2S(&msg.S2S_HanldeFromUserMsg{}, HanldeFromGameMsg)

	consul.SetHookRpc(ChanRPC)
}

//玩家在本服节点登录
func SelfNodeAddPlayer(args []interface{}) {
	uid := args[0].(int64)
	ch := args[1].(*chanrpc.Server)
	Users[uid] = ch
	cluster.Broadcast(cost.HallPrefix, &msg.S2S_NotifyOtherNodeLogin{
		Uid:        uid,
		ServerName: conf.ServerName(),
	})
}

//本服玩家登出
func SelfNodeDelPlayer(args []interface{}) {
	uid := args[0].(int64)
	delete(Users, uid)
	cluster.Broadcast(cost.HallPrefix, &msg.S2S_NotifyOtherNodelogout{
		Uid: uid,
	})
}

//玩家在别的节点登录了
func NotifyOtherNodeLogin(args []interface{}) {
	recvMsg := args[0].(*msg.S2S_NotifyOtherNodeLogin)
	OtherUsers[recvMsg.Uid] = recvMsg.ServerName
	log.Debug("user %d login on %s", recvMsg.Uid, recvMsg.ServerName)
}

//玩家在别的节点登出了
func NotifyOtherNodelogout(args []interface{}) {
	recvMsg := args[0].(*msg.S2S_NotifyOtherNodelogout)
	log.Debug("user %d logout on %s", recvMsg.Uid, OtherUsers[recvMsg.Uid])
	delete(OtherUsers, recvMsg.Uid)
}

func SendMsgToSelfNotdeUser(args []interface{}) {
	uid := args[0].(int64)
	FuncName := args[1].(string)
	ch, ok := Users[uid]
	if ok {
		ch.Go(FuncName, args[2:]...)
		return
	} else {

	}
	log.Debug("at SendMsgToSelfNotdeUser player not in node")
	return
}

func SendMsgToHallUser(args []interface{}) {
	sendMsg := &msg.S2S_HanldeFromUserMsg{}
	sendMsg.Uid = args[0].(int64)
	data, err := cluster.Processor.Marshal(args[1])
	if err != nil {
		log.Debug("SendMsgToHallUser error: %s", err.Error())
		return
	}
	sendMsg.Data = data[0]

	u := Users[sendMsg.Uid]
	if u != nil {
		HanldeFromGameMsg([]interface{}{sendMsg})
		return
	}

	serverName := OtherUsers[sendMsg.Uid]
	if serverName == "" {
		return
	}

	cluster.Go(serverName, sendMsg)
}

//处理来自游戏服的消息
func HanldeFromGameMsg(args []interface{}) {
	recvMsg := args[0].(*msg.S2S_HanldeFromUserMsg)
	log.Debug("at HanldeFromGameMsg == ", msg.S2S_HanldeFromUserMsg{})
	data, err := cluster.Processor.Unmarshal(recvMsg.Data)
	if err != nil {
		log.Error("at HanldeFromGameMsg Unmarshal error:%s", err.Error())
	}
	msgId, err1 := cluster.Processor.GetMsgId(data)
	if err1 != nil {
		log.Error("at HanldeFromGameMsg error:%s", err1.Error())
	}
	SendMsgToSelfNotdeUser([]interface{}{recvMsg.Uid, msgId, data})
}

func GetPlayerInfo(args []interface{}) (interface{}, error) {
	recvMsg := args[0].(*msg.S2S_GetPlayerInfo)
	log.Debug("at GetPlayerInfo uid:%d", recvMsg.Uid)
	ch, chok := Users[recvMsg.Uid]
	if !chok {
		return nil, errors.New("not foud user ch")
	}
	us, err := ch.TimeOutCall1("GetUser", 5)
	if err != nil {
		return nil, err
	}

	u, ok := us.(*user.User)
	if !ok {
		return nil, errors.New("user data error")
	}

	gu := &msg.S2S_GetPlayerInfoResult{
		Id:          u.Id,
		NickName:    u.NickName,
		Currency:    u.Currency,
		RoomCard:    u.RoomCard,
		FaceID:      u.FaceID,
		CustomID:    u.CustomID,
		HeadImgUrl:  u.HeadImgUrl,
		Experience:  u.Experience,
		Gender:      u.Gender,
		WinCount:    u.WinCount,
		LostCount:   u.LostCount,
		DrawCount:   u.DrawCount,
		FleeCount:   u.FleeCount,
		UserRight:   u.Accountsmember.UserRight,
		Score:       u.Score,
		Revenue:     u.Revenue,
		InsureScore: u.InsureScore,
		MemberOrder: u.MemberOrder,
		RoomId:      u.Roomid,
	}
	return gu, nil
}

//新的节点启动了
func serverStart(args []interface{}) {
	svr := args[0].(*consul.CacheInfo)
	log.Debug("%s on line", svr.Csid)
	cluster.AddClient(&cluster.NsqClient{Addr: svr.Host, ServerName: svr.Csid})
	if ok, _ := regexp.Match(cost.GamePrefix, []byte(svr.Csid)); ok { //如果是游戏服启动
		GamelistRpc.Go("NewServerAgent", svr.Csid)
	}
}

//节点关闭了
func serverFaild(args []interface{}) {
	svr := args[0].(*consul.CacheInfo)
	log.Debug("%s off line", svr.Csid)
	list := strings.Split(svr.Csid, "_")
	if len(list) < 2 {
		log.Error("at ServerFaild param error ")
		return
	}

	id, err := strconv.Atoi(list[1])
	if err != nil {
		log.Error("at ServerFaild param error : %s", err.Error())
		return
	}
	cluster.RemoveClient(svr.Csid)
	if ok, _ := regexp.Match(cost.GamePrefix, []byte(svr.Csid)); ok { //如果是游戏服关闭
		GamelistRpc.Go("FaildServerAgent", id)
	}
}
