package room

import (
	"mj/common/msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/common/msg/nn_tb_msg"
	"mj/gameServer/user"
)

func NewNNTBEntry(info *msg.L2G_CreatorRoom) *NNTB_Entry {
	e := new(NNTB_Entry)
	e.Entry_base = pk_base.NewPKBase(info)
	return e
}

///主消息入口
type NNTB_Entry struct {
	*pk_base.Entry_base
}




//叫分(倍数)
func (room *NNTB_Entry) CallScore(args []interface{}) {
	recvMsg := args[0].(*nn_tb_msg.C2G_TBNN_CallScore)
	u := args[1].(*user.User)

	room.DataMgr.CallScore(u, recvMsg.CallScore)
	return
}

//加注
func (r *NNTB_Entry) AddScore(args []interface{}) {
	recvMsg := args[0].(*nn_tb_msg.C2G_TBNN_AddScore)
	u := args[1].(*user.User)

	r.DataMgr.AddScore(u, recvMsg.Score)
	return
}

// 亮牌
func (r *NNTB_Entry) OpenCard(args []interface{}) {
	recvMsg := args[0].(*nn_tb_msg.C2G_TBNN_OpenCard)
	u := args[1].(*user.User)

	r.DataMgr.OpenCard(u, recvMsg.CardType, recvMsg.CardData)
	return
}

