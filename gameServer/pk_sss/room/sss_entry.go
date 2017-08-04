package room

import (
	"mj/common/msg/pk_sss_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"
	"mj/gameServer/user"
)

///主消息入口
type SSS_Entry struct {
	*pk_base.Entry_base
}

func NewSSSEntry(info *model.CreateRoomInfo) *SSS_Entry {
	e := new(SSS_Entry)
	e.Entry_base = pk_base.NewPKBase(info)
	return e
}

// 十三水摊牌
func (r *SSS_Entry) ShowSSSCard(args []interface{}) {
	recvMsg := args[0].(*pk_sss_msg.C2G_SSS_Open_Card)
	u := args[1].(*user.User)

	r.DataMgr.(*sss_data_mgr).ShowSSSCard(u, recvMsg.Dragon, recvMsg.SpecialType, recvMsg.SpecialData, recvMsg.FrontCard, recvMsg.MidCard, recvMsg.BackCard)
	return
}
