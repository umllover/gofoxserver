package internal

import (
	"fmt"
	"math"
	"mj/common/cost"
	"mj/common/msg"
	"mj/hallServer/common"
	"mj/hallServer/conf"
	"mj/hallServer/db/model"
	"mj/hallServer/id_generate"
	"reflect"
	"strconv"
	"strings"

	"mj/hallServer/center"
	"mj/hallServer/user"

	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

var (
	gameLists    = make(map[int]*ServerInfo)   //k1 NodeID,
	roomList     = make(map[int]*msg.RoomInfo) // key1 is roomId
	roomKindList = make(map[int]map[int]int)   //key1 is kind Id key2 incId
	KindListInc  = 0
	Test         = false
)

type ServerInfo struct {
	wight int
	list  map[int]*msg.TagGameServer //key is KindID
}

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

	handleRpc("sendGameList", sendGameList)
	handleRpc("updateGameInfo", updateGameInfo)
	handleRpc("delGameList", delGameList)
	handleRpc("NewServerAgent", NewServerAgent)
	handleRpc("CloseServerAgent", CloseServerAgent)

	handleRpc("notifyNewRoom", NotifyNewRoom)
	handleRpc("notifyDelRoom", NotifyDelRoom)
	handleRpc("updateRoomInfo", UpdateRoom)

	handleRpc("SvrverFaild", ServerFaild)
	handleRpc("SendPlayerBrief", SendPlayerBrief)

	handleRpc("GetMatchRooms", GetMatchRooms)
}

////// c2s
//玩家请求查找房间
func SrarchTable(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_SearchServerTable)
	agent := args[1].(gate.Agent)
	retcode := 0
	defer func() {
		if retcode != 0 {
			agent.WriteMsg(cost.RenderErrorMessage(retcode))
		}
	}()

	roomInfo := getRoomInfo(recvMsg.TableID)
	if roomInfo == nil {
		log.Error("at SrarchTable not foud room, %v", recvMsg)
		retcode = cost.ErrNoFoudRoom
		return
	}

	agent.ChanRPC().Go("SrarchTableResult", roomInfo)
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

//////////////////// rpc
func sendGameList(args []interface{}) {
	agent := args[0].(gate.Agent)
	list := make(msg.L2C_ServerList, 0)
	for _, v := range gameLists {
		for _, v1 := range v.list {
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
	for _, v := range args {
		log.Debug("at NotifyNewRoom === %v", v)
	}

	roomInfo := args[0].(*msg.RoomInfo)
	roomInfo.Players = make(map[int]*msg.PlayerBrief)
	roomInfo.MachPlayer = make(map[int]struct{})
	addRoom(roomInfo)
}

func addRoom(recvMsg *msg.RoomInfo) {
	roomList[recvMsg.RoomID] = recvMsg
	m, ok := roomKindList[recvMsg.KindID]
	if !ok {
		m = make(map[int]int)
		roomKindList[recvMsg.KindID] = m
	}

	if int32(KindListInc) >= math.MaxInt32 {
		KindListInc = 0
	}
	KindListInc++
	recvMsg.Idx = KindListInc
	m[KindListInc] = recvMsg.RoomID
}

func NotifyDelRoom(args []interface{}) {
	log.Debug("at NotifyDelRoom === %v", args)
	roomId := args[0].(int)
	delRoom(roomId)
}

func delRoom(roomId int) {
	ri := roomList[roomId]
	delete(roomList, roomId)
	id_generate.DelRoomId(roomId)
	model.CreateRoomInfoOp.Delete(roomId)
	m, ok := roomKindList[ri.KindID]
	if ok && ri != nil {
		delete(m, ri.Idx)
	} else {
		log.Error("at NotifyDelRoom not foud kind id %v", ri.KindID)
	}
}

func UpdateRoom(args []interface{}) {
	info := args[0].(*msg.UpdateRoomInfo)
	room, ok := roomList[info.RoomId]
	if !ok {
		log.Debug("at  UpdateRoom not foud kindid:%d", info.RoomId)
		return
	}

	switch info.OpName {
	case "CurPayCnt":
		room.CurPayCnt = info.Data["CurPayCnt"].(int)
	case "AddPlayerId":
		pinfo := info.Data["info"].(*msg.PlayerBrief)
		room.Players[pinfo.UID] = pinfo
		room.CurCnt = len(room.Players)
		center.SendMsgToThisNodeUser(pinfo.UID, "JoinRoom", room)
	case "DelPlayerId":
		id := info.Data["UID"].(int)
		status := info.Data["Status"].(int)
		delete(room.Players, id)
		room.CurCnt = len(room.Players)
		if status == 0 { //返回钱
			center.SendMsgToThisNodeUser(id, "restoreToken", info.RoomId)
		}
		center.SendMsgToThisNodeUser(id, "LeaveRoom", info.RoomId)
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
			log.Error("GetKindList error:%s", err.Error())
			return
		}

		ret := data.([]*msg.TagGameServer)

		for _, v := range ret {
			if Test {
				if v.NodeID != conf.Server.NodeId {
					continue
				}
			}
			addGameList(v)
			log.Debug("add sverInfo %v", v)
		}
	})

	cluster.AsynCall(serverName, skeleton.GetChanAsynRet(), "GetRooms", func(data interface{}, err error) {
		if err != nil {
			log.Error("GetKindList error:%s", err.Error())
			return
		}

		ret := data.([]*msg.RoomInfo)

		for _, v := range ret {
			if Test {
				if v.NodeID != conf.Server.NodeId {
					continue
				}
			}
			addRoom(v)
			log.Debug("add room %v", v)
		}
	})
}

func CloseServerAgent(args []interface{}) {
	log.Debug("at CloseServerAgent")
}

///////////////// help
func delGameList(args []interface{}) {
	NodeId := args[0].(int)
	delete(gameLists, NodeId)
}

func addGameList(v *msg.TagGameServer) {
	gminfo, ok := gameLists[v.NodeID]
	if !ok {
		gminfo = new(ServerInfo)
		gminfo.list = make(map[int]*msg.TagGameServer)
		gameLists[v.NodeID] = gminfo
	}

	gminfo.list[v.KindID] = v
}

func GetSvrByKind(kindId int) (string, int) {
	var minv *ServerInfo
	for _, v := range gameLists {
		if _, ok := v.list[kindId]; !ok {
			continue
		}

		if minv == nil {
			minv = v
		}

		if v.wight < minv.wight {
			minv = v
		}

		if Test {
			fmt.Println(v.list[kindId].NodeID, conf.Server.NodeId)
			if v.list[kindId].NodeID == conf.Server.NodeId {
				minv = v
				break
			}
		}

	}

	if minv == nil || len(minv.list) < 1 {
		return "", 0
	}
	minv.wight++
	return minv.list[kindId].ServerAddr + ":" + strconv.Itoa(minv.list[kindId].ServerPort), minv.list[kindId].NodeID
}

func GetSvrByNodeID(nodeid int) string {
	for _, v := range gameLists {
		for _, v1 := range v.list {
			if v1.NodeID != nodeid {
				break
			}
			return v1.ServerAddr + ":" + strconv.Itoa(v1.ServerPort)
		}
	}
	return ""
}

func ServerFaild(args []interface{}) {
	svrId := args[0].(string)
	list := strings.Split(svrId, "_")
	if len(list) < 2 {
		log.Error("at ServerFaild param error ")
		return
	}

	id, err := strconv.Atoi(list[1])
	if err != nil {
		log.Error("at ServerFaild param error : %s", err.Error())
		return
	}

	for roomId, v := range roomList {
		if v.NodeID == id {
			delete(roomList, roomId)
		}
	}
}

func SendPlayerBrief(args []interface{}) {
	roomId := args[0].(int)
	u := args[1].(*user.User)
	retMsg := &msg.L2C_RoomPlayerBrief{}
	r := roomList[roomId]
	if r != nil {
		for _, v := range r.Players {
			retMsg.Players = append(retMsg.Players, v)
		}
	}
	u.WriteMsg(retMsg)
}

func GetMatchRooms(args []interface{}) (interface{}, error) {
	ret := make(map[int][]*msg.RoomInfo)
	for _, v := range roomList {
		if !v.IsPublic {
			continue
		}
		if v.MaxCnt >= len(v.MachPlayer) {
			continue
		}
		ret[v.KindID] = append(ret[v.KindID], v)
	}
	return ret, nil
}
