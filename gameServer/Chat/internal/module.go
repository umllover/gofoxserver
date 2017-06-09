package internal

import (
	"mj/gameServer/base"
	"mj/gameServer/common"
	"mj/gameServer/conf"
	"mj/gameServer/hzmj"

	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/gate"
	"mj/gameServer/user"
	"github.com/lovelly/leaf/log"
	"mj/gameServer/hzmj/room"
	"go/importer"
)

type ChatRoom struct {
	members	map[int]gate.Agent
	memChatIndex map[int]int			//成员聊天索引
}

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	roomList = make(map[int]*ChatRoom)
	roomID int
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton

}

func (m *Module) OnDestroy() {

}

//创建聊天房间
func createRoom(ag gate.Agent)  {
	user := ag.UserData().(*user.User)
	room := new(ChatRoom)
	room.members[user.Id] = ag
	roomList[roomID] = room
	room.memChatIndex[user.Id]=0
	user.ChatRoomId=roomID
	roomID++
}

//增加聊天房间成员
func addRoomMember(roomID int, ag gate.Agent)  {
	room,ok:= roomList[roomID]
	if !ok {
		log.Error("聊天房间：%s不存在",roomID)
		return
	}
	user := ag.UserData().(*user.User)
	room.members[user.Id] = ag
	room.memChatIndex[user.Id]=0
}

//删除聊天房间成员
func delRoomMember(GetRoomID int,ag gate.Agent)  {
	room,ok := roomList[GetRoomID]
	if !ok {
		log.Error("聊天房间：%s不存在",GetRoomID)
		return
	}
	user :=ag.UserData().(*user.User)

	size := len(room.members)
	if size>1 {
		delete(room.members,user.Id)
		delete(room.memChatIndex,user.Id)
	}else {
		delete(room.memChatIndex,user.Id)
		delete(room.members,user.Id)
		delete(roomList,GetRoomID)
	}
}

func SendMsgToUser(getRoomID int,userID int,data interface{})  {
	room, ok := roomList[getRoomID]
	if !ok {
		log.Error("聊天房间：%s不存在",getRoomID)
		return
	}
	for id, ag := range room.members {
		if id==userID {
			ag.WriteMsg(data)
			return
		}
	}
}

func SendMsgToAll(getRoomID int, data interface{}){
	room, ok := roomList[getRoomID]
	if !ok {
		log.Error("聊天房间：%s不存在",getRoomID)
		return
	}
	for id, ag := range room.members {
		ag.WriteMsg(data)
		user :=ag.UserData().(*user.User)
		room.memChatIndex[user.Id]++
	}
}
