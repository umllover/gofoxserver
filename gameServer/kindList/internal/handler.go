package internal

import (
	"mj/common/msg"
	"mj/common/register"
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common"
	"mj/gameServer/conf"
	"mj/gameServer/db/model/base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	//rpc
	reg.RegisterS2S(&msg.S2S_GetKindList{}, GetKindList)
	reg.RegisterS2S(&msg.S2S_GetRooms{}, GetRooms)
	reg.RegisterS2S(&msg.S2S_RenewalFee{}, RenewalFee)
}

///// rpc
func GetKindList(args []interface{}) (interface{}, error) {
	ip, port := conf.GetServerAddrAndPort()

	ret := &msg.S2S_KindListResult{}
	for kind, v := range modules {
		templates, ok := base.GameServiceOptionCache.GetKey1(kind)
		if !ok {
			log.Error("at get kind list not foud kind %d", kind)
			continue
		}
		for _, template := range templates {
			svrInfo := &msg.TagGameServer{}
			svrInfo.KindID = kind
			svrInfo.NodeID = conf.Server.NodeId
			svrInfo.SortID = template.SortID
			svrInfo.ServerID = template.ServerID
			svrInfo.ServerPort = port
			svrInfo.ServerType = int64(template.GameType)
			svrInfo.OnLineCount = int64(v.GetClientCount())
			svrInfo.FullCount = common.TableFullCount
			svrInfo.MinTableScore = 0 //暂时无效
			svrInfo.MinEnterScore = int64(template.MinEnterScore)
			svrInfo.MaxEnterScore = int64(template.MaxEnterScore)
			svrInfo.ServerAddr = ip
			svrInfo.ServerName = template.RoomName
			svrInfo.SurportType = 0
			svrInfo.TableCount = v.GetTableCount()
			ret.Data = append(ret.Data, svrInfo)
		}
	}

	//log.Debug("at S2S_GetKindList ==== %v", ret)
	return ret, nil
}

func GetRooms(args []interface{}) (interface{}, error) {
	rooms := &msg.S2S_GetRoomsResult{}
	RoomMgr.ForEachRoom(func(r RoomMgr.IRoom) {
		rooms.Data = append(rooms.Data, r.GetBirefInfo())
	})
	log.Debug("at S2S_GetRooms ==== %v", rooms)
	return rooms, nil
}

//大厅服服通知续费
func RenewalFee(args []interface{}) {
	log.Debug("at game RenewalFee start")
	retCode := 0
	var retErr error = nil
	recvMsg := args[0].(*msg.S2S_RenewalFee)
	defer func() {
		//通知大厅续费结果
		cluster.SendMsgToHallUser(recvMsg.HallNodeID, recvMsg.UserId, &msg.S2S_RenewalFeeResult{RoomId: recvMsg.RoomID, ResultId: retCode, AddCount: recvMsg.AddCnt})
		log.Debug("at game RenewalFee end, retCode=%d, retErr=%v", retCode, retErr)
	}()

	room := RoomMgr.GetRoom(recvMsg.RoomID)
	if room == nil {
		retCode = 1
		return
	}

	ret, err := room.GetChanRPC().Call1("RenewalFeesSetInfo", recvMsg.AddCnt, recvMsg.UserId, recvMsg.HallNodeID)
	if err != nil {
		retCode = ret.(int)
		retErr = err
		return
	}

	return
}
