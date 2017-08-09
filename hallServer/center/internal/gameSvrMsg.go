package internal

import (
	. "mj/common/cost"
	"mj/common/msg"
)

func RoomReturnMoney(args []interface{}) {
	recvMsg := args[0].(*msg.RoomReturnMoney)
	AddOfflineHandler(MailTypeReturnMoney, recvMsg.CreatorUid, recvMsg, true)
}
