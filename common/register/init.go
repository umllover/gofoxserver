package register

import (
	"encoding/gob"
	"mj/common/msg"

	"github.com/lovelly/leaf/nsq/cluster"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	gob.Register(bson.NewObjectId())
	gob.Register([]bson.ObjectId{})
	gob.Register(map[string]string{})
	gob.Register([]*msg.TagGameServer{})
	gob.Register(&msg.RoomInfo{})
	gob.Register([]*msg.RoomInfo{})
	gob.Register(&msg.RoomEndInfo{})
	gob.Register(&msg.UpdateRoomInfo{})
	gob.Register(&msg.PlayerBrief{})

	cluster.Processor.Register(&msg.S2S_KindListResult{})
	cluster.Processor.Register(&msg.S2S_GetKindList{})
	cluster.Processor.Register(&msg.RoomInfo{})
	cluster.Processor.Register(&msg.S2S_GetRoomsResult{})
	cluster.Processor.Register(&msg.UpdateRoomInfo{})
	cluster.Processor.Register(&msg.PlayerBrief{})
	cluster.Processor.Register(&msg.S2S_GetRooms{})
	cluster.Processor.Register(&msg.S2S_notifyDelRoom{})
	cluster.Processor.Register(&msg.S2S_GetPlayerInfo{})
	cluster.Processor.Register(&msg.S2S_NotifyOtherNodelogout{})
	cluster.Processor.Register(&msg.S2S_NotifyOtherNodeLogin{})
	cluster.Processor.Register(&msg.S2S_GetPlayerInfoResult{})
	cluster.Processor.Register(&msg.S2S_HanldeFromUserMsg{})
	cluster.Processor.Register(&msg.S2S_RenewalFee{})
	cluster.Processor.Register(&msg.L2G_CreatorRoom{})
	cluster.Processor.Register(&msg.JoinRoom{})
	cluster.Processor.Register(&msg.LeaveRoom{})
	cluster.Processor.Register(&msg.StartRoom{})
	cluster.Processor.Register(&msg.JoinRoomFaild{})
	cluster.Processor.Register(&msg.RoomEndInfo{})
	//cluster.Processor.Register(&msg.L2L_NewRoomInfo{})
	cluster.Processor.Register(&msg.S2S_OfflineHandler{})
	cluster.Processor.Register(&msg.RoomReturnMoney{})
	cluster.Processor.Register(&msg.S2S_CloseRoom{})
	cluster.Processor.Register(&msg.S2S_RenewalFeeFaild{})
	cluster.Processor.Register(&msg.OfflineAddElectId{})
}
