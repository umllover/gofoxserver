package internal

import (
	"mj/common/msg"
	"mj/common/register"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/gate"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	reg.RegisterRpc("createRoom", createRoom)
	reg.RegisterRpc("addRoomMember", addRoomMember)
	reg.RegisterRpc("delRoomMember", delRoomMember)

	reg.RegisterC2S(&msg.C2G_GameChart_ToAll{}, SendChatMsgToAll)

}

//发送给房间所有人
func SendChatMsgToAll(args []interface{}) {
	getData := args[0].(*msg.C2G_GameChart_ToAll)
	agent := args[1].(gate.Agent)
	user := agent.UserData().(*user.User)

	sendData := &msg.G2C_GameChart_ToAll{}
	sendData.ChatColor = getData.ChatColor
	sendData.SendUserID = user.Id
	sendData.TargetUserID = getData.SendUserID
	sendData.ChatString = getData.ChatString
	sendData.ChatType = getData.ChatType
	sendData.ChatIndex = getData.ChatIndex

	SendMsgToAll(user.ChatRoomId, sendData)
}

//发送给房间某人
func sendCharMsgToUser(args []interface{}) {

}
