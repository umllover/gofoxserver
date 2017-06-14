package internal

import (
	"mj/common/msg"
	"reflect"

	"mj/hallServer/conf"

	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

var (
	gameLists  = make(map[int]map[int]*msg.TagGameServer) //k1 is kind k2 is server Id
	SvrvetType = make(map[int]map[int]struct{})           //key sverId v KingId
	roomList   = make(map[int]map[int]*msg.RoomInfo)      //key1 is kindid key2 is roomId
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
	handlerC2S(&msg.C2L_SearchServerTable{}, SrarchTable)

	handleRpc("sendGameList", sendGameList)
	handleRpc("updateGameInfo", updateGameInfo)
	handleRpc("delGameList", delGameList)
	handleRpc("NewServerAgent", NewServerAgent)
	handleRpc("CloseServerAgent", CloseServerAgent)

	handleRpc("notifyNewRoom", AddRoom)
	handleRpc("notifyDelRoom", DelRoom)
	handleRpc("updateRoomInfo", UpdateRoom)
}

////// c2s
func SrarchTable(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_SearchServerTable)
	retMsg := &msg.G2C_SearchResult{}
	agent := args[1].(gate.Agent)
	defer func() {
		agent.WriteMsg(retMsg)
	}()

	roomInfo := getRoomInfo(recvMsg.KindID, recvMsg.TableID)
	if roomInfo == nil {
		log.Error("at SrarchTable not foud room, %v", recvMsg)
		return
	}
	retMsg.TableID = roomInfo.RoomID
	retMsg.ServerID = roomInfo.ServerID
	return
}

//////////////////// rpc
func sendGameList(args []interface{}) {
	agent := args[0].(gate.Agent)
	list := make(msg.L2C_ServerList, 0)
	for _, v := range gameLists {
		for _, v1 := range v {
			list = append(list, v1)
		}
	}
	agent.WriteMsg(&list)
	finish := &msg.L2C_ServerListFinish{}
	agent.WriteMsg(finish)
}

func updateGameInfo(args []interface{}) {

}

func AddRoom(args []interface{}) {
	log.Debug("at AddRoom === %v", args)
	recvMsg := args[0].(*msg.RoomInfo)

	m, ok := roomList[recvMsg.KindID]
	if !ok {
		m = make(map[int]*msg.RoomInfo)
		roomList[recvMsg.KindID] = m
	}
	m[recvMsg.RoomID] = recvMsg
}

func DelRoom(args []interface{}) {
	log.Debug("at DelRoom === %v", args)
	id := args[0].(int)
	delete(roomList, id)
}

func UpdateRoom(args []interface{}) {
	info := args[0].(map[string]interface{})
	kindID := info["KindID"].(int)
	roomID := info["RoomID"].(int)
	m, ok := roomList[kindID]
	if !ok {
		log.Debug("at  UpdateRoom not foud kindid:%d", kindID)
		return
	}

	room, ok2 := m[roomID]
	if !ok2 {
		log.Debug("at  UpdateRoom not foud roomID:%d", roomID)
		return
	}

	for k, v := range info {
		switch k {
		case "CurCnt":
			room.CurCnt = v.(int)
		case "CurPayCnt":
			room.CurPayCnt = v.(int)
		}
	}
}

func getRoomInfo(kindId, tableId int) *msg.RoomInfo {
	m, ok := roomList[kindId]
	if !ok {
		log.Debug("at getRoomInfo not foud kindID:%v", kindId)
		return nil
	}
	return m[tableId]
}

func NewServerAgent(args []interface{}) {
	serverName := args[0].(string)
	log.Debug("at NewServerAgent :%s", serverName)
	cluster.AsynCall(serverName, skeleton.GetChanAsynRet(), "GetKindList", func(data interface{}, err error) {
		if err != nil {
			log.Debug("GetKindList error:%s", err.Error())
			return
		}

		ret := data.([]*msg.TagGameServer)

		for _, v := range ret {
			if conf.Server.TestNode {
				if v.NodeID != conf.Server.NodeId {
					continue
				}
			}
			addGameList(v)
			log.Debug("data %v", v)
		}
	})
}

func CloseServerAgent(args []interface{}) {
	log.Debug("at CloseServerAgent")
}

///////////////// help
func delGameList(args []interface{}) {
	svrId := args[0].(int)
	typeInfo := SvrvetType[svrId]
	for kind, _ := range typeInfo {
		gminfo, ok := gameLists[kind]
		if ok {
			delete(gminfo, svrId)
		}
	}
}

func addGameList(v *msg.TagGameServer) {
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
	typeInfo[v.KindID] = struct{}{}
}
