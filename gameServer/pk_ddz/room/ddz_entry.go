package room

import (
	"mj/common/msg"
	"mj/common/msg/pk_ddz_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"
)

func NewDDZEntry(info *msg.L2G_CreatorRoom) *DDZ_Entry {
	e := new(DDZ_Entry)
	e.Entry_base = pk_base.NewPKBase(info)
	return e
}

///主消息入口
type DDZ_Entry struct {
	*pk_base.Entry_base
}

// 叫分(倍数)
func (room *DDZ_Entry) CallScore(args []interface{}) {
	recvMsg := args[0].(*pk_ddz_msg.C2G_DDZ_CallScore)
	u := args[1].(*user.User)

	room.DataMgr.CallScore(u, recvMsg.CallScore)
	return
}

// 用户出牌
func (room *DDZ_Entry) OutCard(args []interface{}) {

	recvMsg := args[0].(*pk_ddz_msg.C2G_DDZ_OutCard)
	u := args[1].(*user.User)

	log.Debug("用户%d出牌%v", u.ChairId, recvMsg)
	room.DataMgr.OpenCard(u, recvMsg.CardType, recvMsg.CardData)
}

// 托管
func (room *DDZ_Entry) CTrustee(args []interface{}) {
	recvMsg := args[0].(*pk_ddz_msg.C2G_DDZ_TRUSTEE)
	u := args[1].(*user.User)

	room.DataMgr.OtherOperation([]interface{}{"Trustee", u, recvMsg})
}

// 明牌
func (r *DDZ_Entry) ShowCard(args []interface{}) {
	u := args[0].(*user.User)
	r.DataMgr.OtherOperation([]interface{}{"ShowCard", u})
}

// 放弃出牌
func (room *DDZ_Entry) PassCard(args []interface{}) {
	u := args[1].(*user.User)
	room.DataMgr.OtherOperation([]interface{}{"PassCard", u})
}
