package internal

import (
	"errors"
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	rgst "mj/common/register"
	"mj/hallServer/center"
	"mj/hallServer/common"
	"mj/hallServer/conf"
	"mj/hallServer/db/model"
	"mj/hallServer/id_generate"
	"mj/hallServer/user"
	"strconv"
	"sync"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
)

var (
	reg          = rgst.NewRegister(ChanRPC)
	gameLists    = make(map[int]*ServerInfo) //k1 NodeID,
	gameListLock sync.RWMutex
	roomList     = make(map[int]*msg.RoomInfo)    // key1 is roomId
	roomKindList = make(map[int]map[int]struct{}) //key1 is kind Id key2 incId
	MatchRpc     *chanrpc.Server
)

type ServerInfo struct {
	wight int
	list  map[int]*msg.TagGameServer //key is KindID
}

func init() {
	reg.RegisterC2S(&msg.C2L_GetRoomList{}, GetRoomList)

	reg.RegisterRpc("sendGameList", sendGameList)
	reg.RegisterRpc("updateGameInfo", updateGameInfo)
	reg.RegisterRpc("delGameList", delGameList)
	reg.RegisterRpc("CloseServerAgent", CloseServerAgent)
	reg.RegisterRpc("addyNewRoom", addyNewRoom)
	reg.RegisterRpc("NewServerAgent", NewServerAgent)
	reg.RegisterRpc("FaildServerAgent", FaildServerAgent)
	reg.RegisterRpc("SendPlayerBrief", sendPlayerBrief)
	reg.RegisterRpc("GetMatchRoomsByKind", GetMatchRoomsByKind)
	reg.RegisterRpc("GetRoomByRoomId", GetRoomByRoomId)
	reg.RegisterRpc("HaseRoom", HaseRoom)
	reg.RegisterRpc("CheckVaildIds", CheckVaildIds)

	reg.RegisterS2S(&msg.S2S_notifyDelRoom{}, notifyDelRoom)
	reg.RegisterS2S(&msg.UpdateRoomInfo{}, updateRoom)
	reg.RegisterS2S(&msg.RoomInfo{}, notifyNewRoom)

	center.SetGameListRpc(ChanRPC)
}

////// c2s

func GetRoomList(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_GetRoomList)
	retMsg := &msg.L2C_GetRoomList{}
	retMsg.Lists = make([]*msg.RoomInfo, 0)
	agent := args[1].(gate.Agent)
	defer func() {
		agent.WriteMsg(retMsg)
	}()

	if recvMsg.Num > common.GetGlobalVarInt(MAX_SHOW_ENTRY) {
		return
	}

	m, ok := roomKindList[recvMsg.KindID]

	if ok {
		for roomID, _ := range m {
			if retMsg.Count >= recvMsg.Num {
				break
			}

			r := roomList[roomID]
			if r == nil {
				log.Debug("error at GetRoomList")
				continue
			}

			if !r.IsPublic {
				continue
			}
			retMsg.Lists = append(retMsg.Lists, r)
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
	roomInfo.MachPlayer = make(map[int64]int64)
	addRoom(roomInfo)
}

//本服增加创建的房间
func addyNewRoom(args []interface{}) {
	for _, v := range args {
		log.Debug("at NotifyNewRoom === %v", v)
	}

	roomInfo := args[0].(*msg.RoomInfo)
	roomInfo.Players = make(map[int64]*msg.PlayerBrief)
	roomInfo.MachPlayer = make(map[int64]int64)
	addRoom(roomInfo)
	cluster.BroadcastToHall(roomInfo)
}

func addRoom(recvMsg *msg.RoomInfo) {
	recvMsg.Players = make(map[int64]*msg.PlayerBrief)
	roomList[recvMsg.RoomID] = recvMsg
	m, ok := roomKindList[recvMsg.KindID]
	if !ok {
		m = make(map[int]struct{})
		roomKindList[recvMsg.KindID] = m
	}
	m[recvMsg.RoomID] = struct{}{}
	log.Debug("addRoom ok, RoomID = %d", recvMsg.RoomID)
}

func notifyDelRoom(args []interface{}) {
	log.Debug("at NotifyDelRoom === %v", args)
	recvMsg := args[0].(*msg.S2S_notifyDelRoom)
	delRoom(recvMsg.RoomID)
}

func delRoom(roomId int) {
	ri := roomList[roomId]
	delete(roomList, roomId)
	id_generate.DelRoomId(roomId)
	model.CreateRoomInfoOp.Delete(roomId) //todo
	if ri != nil {
		m, ok := roomKindList[ri.KindID]
		if ok {
			delete(m, roomId)
		} else {
			log.Error("at NotifyDelRoom not foud kind id %v", ri.KindID)
		}
	}
}

func updateRoom(args []interface{}) {
	info := args[0].(*msg.UpdateRoomInfo)
	room, ok := roomList[info.RoomId]
	if !ok {
		log.Debug("at  UpdateRoom not foud RoomId:%d", info.RoomId)
		return
	}
	log.Debug("=============================info.OpName=%v, info.Data[HallNodeName]=%v", info.OpName, info.Data["HallNodeName"])

	switch info.OpName {
	case "AddPlayCnt":
		if room.CurPayCnt == 0 {
			room.Status = RoomStatusStarting
		}
		room.CurPayCnt += 1
		for _, ply := range room.Players {
			if ply.HallNodeName == conf.ServerName() {
				center.SendToThisNodeUser(ply.UID, "GameStart", &msg.StartRoom{RoomId: room.RoomID})
			}
		}

	case "AddPlayerId":
		pinfo := &msg.PlayerBrief{
			UID:          int64(info.Data["UID"].(float64)),
			Name:         info.Data["Name"].(string),
			HeadUrl:      info.Data["HeadUrl"].(string),
			Icon:         int(info.Data["Icon"].(float64)),
			HallNodeName: info.Data["HallNodeName"].(string),
		}
		room.Players[pinfo.UID] = pinfo
		room.CurCnt = len(room.Players)

		status := int(info.Data["Status"].(float64))
		if pinfo.HallNodeName == conf.ServerName() {
			center.SendToThisNodeUser(pinfo.UID, "JoinRoom", &msg.JoinRoom{Rinfo: room, Status: status})
		}

	case "DelPlayerId":
		log.Debug("at hall room DelPlayerId ........................")
		id := int64(info.Data["UID"].(float64))
		status := int(info.Data["Status"].(float64))
		ply, ok := room.Players[id]
		if ok {
			delete(room.Players, id)
			room.CurCnt = len(room.Players)
			if ply.HallNodeName == conf.ServerName() {
				center.SendToThisNodeUser(id, "LeaveRoom", &msg.LeaveRoom{RoomId: room.RoomID, Status: status})
			}

			if MatchRpc != nil {
				MatchRpc.Go("delMatchPlayer", id, room)
			}
		}
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
				if conf.Test {
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
				if conf.Test {
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
	gameListLock.Lock()
	defer gameListLock.Unlock()
	NodeId := args[0].(int)
	delete(gameLists, NodeId)
}

func addGameList(v *msg.TagGameServer) {
	gameListLock.Lock()
	defer gameListLock.Unlock()
	gminfo, ok := gameLists[v.NodeID]
	if !ok {
		gminfo = new(ServerInfo)
		gminfo.list = make(map[int]*msg.TagGameServer)
		gameLists[v.NodeID] = gminfo
	}

	gminfo.list[v.KindID] = v
}

func GetSvrByKind(kindId int) (string, int) {
	gameListLock.RLock()
	defer gameListLock.RUnlock()
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

		if conf.Test {
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

	if conf.Test {
		if minv.list[kindId].NodeID != conf.Server.NodeId {
			return "", 0
		}
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
	log.Debug("AT 			FaildServerAgent  ")
	id := args[0].(int)
	ids := make(map[int]bool)
	for roomId, v := range roomList {
		if v.NodeID == id {
			for uid, _ := range v.MachPlayer { //房间因为服务器宕机关闭
				///通知大厅房间结束
				center.SendToThisNodeUser(uid, "RoomEndInfo", &msg.RoomEndInfo{RoomId: roomId, Status: v.Status, CreateUid: v.CreateUserId})
			}
			delete(roomList, roomId)
			ids[roomId] = true
		}
	}

	for _, m := range roomKindList {
		for roomid, _ := range m {
			if ids[roomid] {
				delete(m, roomid)
			}
		}
	}

	delete(gameLists, id)

	log.Debug("AT 			FaildServerAgent  ok")
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

//不公开的房间 和人数满的 无法被获取到
func GetMatchRoomsByKind(args []interface{}) (interface{}, error) {
	kind := args[0].(int)
	ret := make([]*msg.RoomInfo, 0)

	now := time.Now().Unix()
	log.Debug("at GetMatchRoomsByKind === %v", roomKindList)
	log.Debug("at GetMatchRoomsByKind rooms  === %v", roomList)
	for roomID, _ := range roomKindList[kind] {
		v := roomList[roomID]
		if v == nil {
			fmt.Println("111111111111")
			continue
		}
		CheckTimeOut(v, now)
		if !v.IsPublic {
			fmt.Println("222222222222")
			continue
		}
		if v.MaxPlayerCnt <= v.MachCnt {
			fmt.Println("3333333333333 %d", v.MaxPlayerCnt)
			continue
		}
		ret = append(ret, v)
	}
	return ret, nil
}

func GetRoomByRoomId(args []interface{}) (interface{}, error) {
	roomid := args[0].(int)
	ret, ok := roomList[roomid]
	if !ok {
		return nil, errors.New("not foud")
	}
	CheckTimeOut(ret, time.Now().Unix())
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

func CheckVaildIds(args []interface{}) {
	ids := args[0].(map[int]struct{})
	ch := args[1].(*chanrpc.Server)
	var invalidIds []int
	for id, _ := range ids {
		_, ok := roomList[id]
		if !ok {
			invalidIds = append(invalidIds, id)
		}
	}
	if len(invalidIds) > 0 {
		ch.Go("DeleteVaildIds", invalidIds)
	}
}

func CheckTimeOut(r *msg.RoomInfo, now int64) {
	for uid, t := range r.MachPlayer {
		if t < now {
			if _, ok := r.Players[uid]; !ok {
				log.Error("at CheckTimeOut player :%d not join room ")
				delete(r.MachPlayer, uid)
			}
		}
	}
}
