package internal

import (
	"encoding/json"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/hallServer/center"
	"mj/hallServer/db/model"
	"mj/hallServer/user"

	"github.com/lovelly/leaf/log"
)

//后期压力这个服务改为redis 做

func loadHandles(player *user.User) {
	log.Debug("at loadHandles .................... ")
	handler, _ := model.UserOfflineHandlerOp.QueryByMap(map[string]interface{}{
		"user_id": player.Id,
	})
	log.Debug("at loadHandles, handler=%v", handler)
	for _, v := range handler {
		log.Debug("at loadHandles .................... 1111111111111111111 ")
		handlerEventFunc(player, v)
	}
}

func handlerEventFunc(player *user.User, v *model.UserOfflineHandler) {
	ret := false
	switch v.HType {
	case OfflineTypeDianZhan:
		ret = handlerDianZhan(player, v)
	case OfflineRoomEndInfo:
		ReturnMoney := &msg.RoomReturnMoney{}
		err := json.Unmarshal([]byte(v.Context), ReturnMoney)
		if err == nil {
			ret = handlerOfflineRoomEndInfo(player, ReturnMoney)
		}
	case OfflineReturnMoney:
		ret = handlerOfflineRoomReturnMone(player, v)
	}

	if ret {
		model.UserOfflineHandlerOp.Delete(v.Id)
	}
}

func AddOfflineHandler(htype string, uid int64, data interface{}, Notify bool) bool {
	log.Debug("#################### AddOfflineHandler htype=%d, uid=%d, data=%v", htype, uid, data)
	h := &model.UserOfflineHandler{
		UserId: uid,
		HType:  htype,
	}

	if data != nil {
		text, err := json.Marshal(data)
		if err != nil {
			log.Debug("add AddOfflineHandler error:%s", err.Error())
			return false
		}
		h.Context = string(text)
	}

	id, ierr := model.UserOfflineHandlerOp.Insert(h)
	if ierr != nil {
		log.Debug("add AddOfflineHandler UserOfflineHandlerOp insert error:%s", ierr.Error())
		return false
	}

	if Notify {
		center.SendMsgToHallUser(uid, &msg.S2S_OfflineHandler{EventID: int(id)})
	}

	return true
}

func handlerDianZhan(player *user.User, msg *model.UserOfflineHandler) bool {
	player.Star++
	model.UserattrOp.UpdateWithMap(player.UserId, map[string]interface{}{
		"star": player.Star,
	})
	//player.WriteMsg(msg.L2C_BeStar{Star:player.Star})
	return true
}

//返还钱给玩家
func handlerOfflineRoomEndInfo(player *user.User, ReturnMoney *msg.RoomReturnMoney) bool {
	log.Debug("at handlerOfflineRoomEndInfo ...............")
	record := player.GetRecord(ReturnMoney.RoomId)
	log.Debug("############## handlerReturnMoney RoomId=%d, record=%v", ReturnMoney.RoomId, record)
	if record != nil {
		log.Debug("############## record.RoomId=%d, record.Amount=%d", record.RoomId, record.Amount)
		player.DelRecord(record.RoomId)
		player.AddCurrency(record.Amount)
	}
	player.DelRooms(ReturnMoney.RoomId)
	return true

}

//返还钱给玩家
func handlerOfflineRoomReturnMone(player *user.User, data *model.UserOfflineHandler) bool {
	ReturnMoney := &msg.RoomReturnMoney{}
	err := json.Unmarshal([]byte(data.Context), ReturnMoney)
	if err != nil {
		log.Error("at handlerReturnMoney Unmarshal error ")
		return false
	}
	record := player.GetRecord(ReturnMoney.RoomId)
	log.Debug("############## handlerReturnMoney RoomId=%d, record=%v", ReturnMoney.RoomId, record)
	if record != nil {
		log.Debug("############## record.RoomId=%d, record.Amount=%d", record.RoomId, record.Amount)
		player.DelRecord(record.RoomId)
		player.AddCurrency(record.Amount)
	}
	return true
}
