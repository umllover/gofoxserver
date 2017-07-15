package internal

import (
	"math"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/hallServer/center"
	"mj/hallServer/common"
	"mj/hallServer/conf"
	"mj/hallServer/db/model"
	"mj/hallServer/id_generate"
	"mj/hallServer/user"
	"strconv"

	rgst "mj/common/register"

	"errors"

	"time"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
)

var (
	reg          = rgst.NewRegister(ChanRPC)
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

func init() {
	reg.RegisterC2S(&msg.C2L_SearchServerTable{}, SrarchTable)
	reg.RegisterC2S(&msg.C2L_GetRoomList{}, GetRoomList)

	reg.RegisterRpc("sendGameList", sendGameList)
	reg.RegisterRpc("updateGameInfo", updateGameInfo)
	reg.RegisterRpc("delGameList", delGameList)
	reg.RegisterRpc("CloseServerAgent", CloseServerAgent)
	reg.RegisterRpc("addyNewRoom", addyNewRoom)
	reg.RegisterRpc("notifyDelRoom", notifyDelRoom)
	reg.RegisterRpc("NewServerAgent", NewServerAgent)
	reg.RegisterRpc("FaildServerAgent", FaildServerAgent)
	reg.RegisterRpc("SendPlayerBrief", sendPlayerBrief)
	reg.RegisterRpc("GetMatchRooms", getMatchRooms)
	reg.RegisterRpc("HaseRoom", HaseRoom)

	reg.RegisterS2S(&msg.UpdateRoomInfo{}, updateRoom)
	reg.RegisterS2S(&msg.RoomInfo{}, notifyNewRoom)

	center.SetGameListRpc(ChanRPC)
}

////// c2s
//玩家请求查找房间
func SrarchTable(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_SearchServerTable)
	agent := args[1].(gate.Agent)
	retcode := 0
	defer func() {
		if retcode != 0 {
			agent.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	roomInfo := getRoomInfo(recvMsg.TableID)
	if roomInfo == nil {
		log.Error("at SrarchTable not foud room, %v", recvMsg)
		retcode = ErrNoFoudRoom
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
	skeleton.AfterFunc(1*time.Second, func() {
		agent.WriteMsg(&msg.L2C_ServerListFinish{})
	})

}

func updateGameInfo(args []interface{}) {

}

//别的服通知的增加的房间
func notifyNewRoom(args []interface{}) {
	for _, v := range args {
		log.Debug("at NotifyNewRoom === %v", v)
	}

	roomInfo := args[0].(*msg.RoomInfo)
	roomInfo.Players = make(map[int64]*msg.PlayerBrief)
	roomInfo.MachPlayer = make(map[int64]struct{})
	addRoom(roomInfo)
}

//本服增加创建的房间
func addyNewRoom(args []interface{}) {
	for _, v := range args {
		log.Debug("at NotifyNewRoom === %v", v)
	}

	roomInfo := args[0].(*msg.RoomInfo)
	roomInfo.Players = make(map[int64]*msg.PlayerBrief)
	roomInfo.MachPlayer = make(map[int64]struct{})
	addRoom(roomInfo)
	center.BroadcastToHall(roomInfo)
}

func addRoom(recvMsg *msg.RoomInfo) {
	recvMsg.Players = make(map[int64]*msg.PlayerBrief)
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

func notifyDelRoom(args []interface{}) {
	log.Debug("at NotifyDelRoom === %v", args)
	roomId := args[0].(int)
	delRoom(roomId)
}

func delRoom(roomId int) {
	ri := roomList[roomId]
	delete(roomList, roomId)
	id_generate.DelRoomId(roomId)
	model.CreateRoomInfoOp.Delete(roomId)
	if ri != nil {
		m, ok := roomKindList[ri.KindID]
		if ok {
			delete(m, ri.Idx)
		} else {
			log.Error("at NotifyDelRoom not foud kind id %v", ri.KindID)
		}
	}
}

func updateRoom(args []interface{}) {
	info := args[0].(*msg.UpdateRoomInfo)
	room, ok := roomList[info.RoomId]
	if !ok {
		log.Debug("at  UpdateRoom not foud kindid:%d", info.RoomId)
		return
	}

	switch info.OpName {
	case "CurPayCnt":
		room.CurPayCnt = int(info.Data["CurPayCnt"].(float64))
	case "AddPlayerId":
		pinfo := &msg.PlayerBrief{
			UID:     int64(info.Data["UID"].(float64)),
			Name:    info.Data["Name"].(string),
			HeadUrl: info.Data["HeadUrl"].(string),
			Icon:    int(info.Data["Icon"].(float64)),
		}
		room.Players[pinfo.UID] = pinfo
		room.CurCnt = len(room.Players)
		center.SendToThisNodeUser(pinfo.UID, "JoinRoom", room)
	case "DelPlayerId":
		id := int64(info.Data["UID"].(float64))
		status := int(info.Data["Status"].(float64))
		delete(room.Players, id)
		room.CurCnt = len(room.Players)
		if status == 0 { //返回钱
			center.SendToThisNodeUser(id, "restoreToken", info.RoomId)
		}
		center.SendToThisNodeUser(id, "LeaveRoom", info.RoomId)
	}

}

func getRoomInfo(tableId int) *msg.RoomInfo {
	return roomList[tableId]
}

func NewServerAgent(args []interface{}) {
	serverName := args[0].(string)
	log.Debug("at NewServerAgent :%s", serverName)

	cluster.AsynCall(serverName, skeleton.GetChanAsynRet(), &msg.S2S_GetKindList{}, func(data interface{}, err error) {
		if err == nil {
			log.Debug("data === %v", data)
			ret := data.(*msg.S2S_KindListResult)
			for _, v := range ret.Data {
				if Test {
					if v.NodeID != conf.Server.NodeId {
						continue
					}
				}
				addGameList(v)
				log.Debug("add sverInfo %v", v)
			}
		} else {
			log.Debug("S2S_GetKindList error:%s", err.Error())
		}
	})

	cluster.AsynCall(serverName, skeleton.GetChanAsynRet(), &msg.S2S_GetRooms{}, func(data interface{}, err error) {
		if err == nil {
			log.Debug("data ======= %v", data)
			ret := data.(*msg.S2S_GetRoomsResult)
			for _, v := range ret.Data {
				if Test {
					if v.NodeID != conf.Server.NodeId {
						continue
					}
				}
				addRoom(v)
				log.Debug("add room %v", v)
			}
		} else {
			log.Debug("S2S_GetRooms error:%s", err.Error())
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
			log.Debug("node id =%d,  self node id =%d", v.list[kindId].NodeID, conf.Server.NodeId)
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

func FaildServerAgent(args []interface{}) {
	id := args[0].(int)
	for roomId, v := range roomList {
		if v.NodeID == id {
			delete(roomList, roomId)
		}
	}
}

func sendPlayerBrief(args []interface{}) {
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

func getMatchRooms(args []interface{}) (interface{}, error) {
	ret := make(map[int][]*msg.RoomInfo)
	for _, v := range roomList {
		if !v.IsPublic {
			continue
		}
		if v.MaxPlayerCnt >= len(v.MachPlayer) {
			continue
		}
		ret[v.KindID] = append(ret[v.KindID], v)
	}
	return ret, nil
}

func HaseRoom(args []interface{}) (interface{}, error) {
	id := args[0].(int)
	_, ok := roomList[id]
	if ok {
		return nil, nil
	}
	return nil, errors.New("no room")
}
