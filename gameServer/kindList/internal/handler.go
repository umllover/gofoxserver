package internal

import (
	"mj/common/msg"
	"mj/gameServer/common"
	"mj/gameServer/conf"
	"mj/gameServer/db/model/base"

	"mj/gameServer/RoomMgr"

	"mj/common/register"

	"mj/gameServer/center"

	"github.com/lovelly/leaf/log"
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
			continue
		}
		for _, template := range templates {
			svrInfo := &msg.TagGameServer{}
			svrInfo.KindID = kind
			svrInfo.NodeID = conf.Server.NodeId
			svrInfo.SortID = template.SortID
			svrInfo.ServerID = template.ServerID
			svrInfo.ServerPort = port
			svrInfo.ServerType = int64(template.ServerType)
			svrInfo.OnLineCount = int64(v.GetClientCount())
			svrInfo.FullCount = common.TableFullCount
			svrInfo.MinTableScore = int64(template.MinTableScore)
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
	retCode := 0
	recvMsg := args[0].(*msg.S2S_RenewalFee)
	defer func() {
		if retCode != 0 { //通知大厅续费失败
			center.SendDataToHallUser(recvMsg.HallName, recvMsg.UserId, &msg.S2S_RenewalFeeFaild{RoomId: recvMsg.RoomID})
		}
	}()
	room := RoomMgr.GetRoom(recvMsg.RoomID)
	if room == nil {
		retCode = 1
		return
	}

	_, err := room.GetChanRPC().Call1("AddPayCnt", recvMsg.AddCnt)
	if err != nil {
		retCode = 2
		return
	}
	return
}
