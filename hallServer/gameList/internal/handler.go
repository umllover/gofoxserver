package internal

import (
	"mj/common/msg"
	"reflect"

	//"math"
	"mj/hallServer/common"
	"mj/hallServer/conf"
	"sort"

	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

var (
	gameLists    = make(map[int]map[int]*msg.TagGameServer) //k1 is kind k2 is server Id
	SvrvetType   = make(map[int]map[int]struct{})           //key sverId v KingId
	roomList     = make(map[int]*msg.RoomInfo)              // key1 is roomId
	roomKindList = make(map[int]map[int]int)                //key1 is kind Id key2 incId
	KindListInc  = 0
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
	handlerC2S(&msg.C2L_GetRoomList{}, GetRoomList)
	handlerC2S(&msg.C2L_QuickMatch{}, QuickMatch)

	handleRpc("sendGameList", sendGameList)
	handleRpc("updateGameInfo", updateGameInfo)
	handleRpc("delGameList", delGameList)
	handleRpc("NewServerAgent", NewServerAgent)
	handleRpc("CloseServerAgent", CloseServerAgent)

	handleRpc("notifyNewRoom", NotifyNewRoom)
	handleRpc("notifyDelRoom", NotifyDelRoom)
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

	roomInfo := getRoomInfo(recvMsg.TableID)
	if roomInfo == nil {
		log.Error("at SrarchTable not foud room, %v", recvMsg)
		return
	}
	retMsg.TableID = roomInfo.RoomID
	retMsg.ServerID = roomInfo.ServerID
	return
}

func GetRoomList(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_GetRoomList)
	retMsg := msg.L2C_GetRoomList{}
	retMsg.Lists = make([]*msg.RoomInfo, common.ListsMaxCnt)
	agent := args[1].(gate.Agent)
	defer func() {
		agent.WriteMsg(retMsg)
	}()

	curIdx := recvMsg.PageId * common.PackCount
	if curIdx > KindListInc {
		curIdx = 0
	}
	m, ok := roomKindList[recvMsg.KindID]

	if ok {
		for idx, roomID := range m {
			if idx <= curIdx {
				continue
			}
			retMsg.Lists[retMsg.Count] = roomList[roomID]
			retMsg.Count++
		}
	}
}

func QuickMatch(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_QuickMatch)
	retMsg := msg.G2C_SearchResult{}
	agent := args[1].(gate.Agent)
	defer func() {
		agent.WriteMsg(retMsg)
	}()

	m, ok := roomKindList[recvMsg.KindID]
	if !ok {
		log.Debug("not found KindID:%v", recvMsg.KindID)
		return
	}

	maxLen := len(m)
	if maxLen < 2 {
		for _, id := range m {
			v := roomList[id]
			retMsg.ServerID = v.ServerID
			retMsg.TableID = v.RoomID
			return
		}
	}

	arr := make([]*msg.RoomInfo, maxLen)
	i := 0
	for _, roomid := range m {
		arr[i] = roomList[roomid]
	}

	sort.Slice(arr, func(i, j int) bool {
		if arr[i].CreateTime < arr[j].CreateTime {
			return true
		}

		if arr[i].CurCnt < arr[j].CurCnt {
			return true
		}
		return false
	})

	v := arr[0]
	retMsg.ServerID = v.ServerID
	retMsg.TableID = v.RoomID
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

func NotifyNewRoom(args []interface{}) {
	log.Debug("at NotifyNewRoom === %v", args)
	recvMsg := args[0].(*msg.RoomInfo)
	roomList[recvMsg.RoomID] = recvMsg
	m, ok := roomKindList[recvMsg.KindID]
	if !ok {
		m = make(map[int]int)
		roomKindList[recvMsg.KindID] = m
	}

	/*if KindListInc >= math.MaxInt {
		KindListInc = 0
	}*/
	KindListInc++
	recvMsg.Idx = KindListInc
	m[KindListInc] = recvMsg.RoomID
}

func NotifyDelRoom(args []interface{}) {
	log.Debug("at NotifyDelRoom === %v", args)
	kindId := args[0].(int)
	roomId := args[1].(int)
	ri := roomList[roomId]
	delete(roomList, roomId)
	m, ok := roomKindList[kindId]
	if ok && ri != nil {
		delete(m, ri.Idx)
	} else {
		log.Error("at NotifyDelRoom not foud kind id %v", kindId)
	}
}

func UpdateRoom(args []interface{}) {
	info := args[0].(map[string]interface{})
	roomID := info["RoomID"].(int)
	room, ok := roomList[roomID]
	if !ok {
		log.Debug("at  UpdateRoom not foud kindid:%d", roomID)
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

func getRoomInfo(tableId int) *msg.RoomInfo {
	return roomList[tableId]
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
