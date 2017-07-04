package room

import (
	"mj/common/msg/nn_tb_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"
	"mj/gameServer/user"
)

func NewDDZEntry(info *model.CreateRoomInfo) *DDZ_Entry {
	e := new(DDZ_Entry)
	return e
	e.Entry_base = pk_base.NewPKBase(info)
	return e
}

///主消息入口
type DDZ_Entry struct {
	*pk_base.Entry_base
}

//叫分(倍数)
func (room *DDZ_Entry) CallScore(args []interface{}) {
	recvMsg := args[0].(*nn_tb_msg.C2G_TBNN_CallScore)
	u := args[1].(*user.User)

	room.DataMgr.CallScore(u, recvMsg.CallScore)
	return
}
