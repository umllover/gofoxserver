package internal

import (
	"mj/common/msg"
	"mj/hallServer/user"
)

//重登的时候删除已经不存在的房间， 后期这些房间放在redis
func (m *UserModule) RoomReturnMoney(args []interface{}) {
	player := m.a.UserData().(*user.User)
	data := args[0].(*msg.RoomReturnMoney)
	handlerOfflineRoomEndInfo(player, data)
}
