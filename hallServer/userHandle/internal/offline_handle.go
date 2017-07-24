package internal

import (
	"encoding/json"
	"mj/common/msg"
	"mj/hallServer/center"
	"mj/hallServer/db/model"
	"mj/hallServer/user"

	"github.com/lovelly/leaf/log"
)

const (
	MailTypeDianZhan = 1
)

//后期压力这个服务改为redis 做

func loadHandles(player *user.User) {
	handler, _ := model.UserOfflineHandlerOp.QueryByMap(map[string]interface{}{
		"user_id": player.Id,
	})

	for _, v := range handler {
		handlerEventFunc(player, v)
	}
}

func handlerEventFunc(player *user.User, v *model.UserOfflineHandler) {
	switch v.HType {
	case MailTypeDianZhan:
		hanDlerDianZhan(player, v)
	}
}

func AddOfflineHandler(htype int, uid int64, data interface{}) bool {
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

	center.SendMsgToHallUser(uid, &msg.S2S_OfflineHandler{EventID: int(id)})
	return true
}

func hanDlerDianZhan(player *user.User, msg *model.UserOfflineHandler) {
	player.Star++
	model.UserattrOp.UpdateWithMap(player.UserId, map[string]interface{}{
		"star": player.Star,
	})

	//player.WriteMsg(msg.L2C_BeStar{Star:player.Star})
}
