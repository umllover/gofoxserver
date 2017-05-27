package room

import (
	"github.com/lovelly/leaf/chanrpc"
)

func RegisterHandler(r *Room) {
	r.ChanRPC.RegisterFromType("EnterRoom", OutCard, chanrpc.FuncThis, r)
}

func OutCard(args []interface{}) (error) {
	card := args[0].(int)
	room := args[len(args) - 1].(*Room)
	room.SendMsgAll(card )
	return nil
}

