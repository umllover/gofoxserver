package common

import (
	"github.com/lovelly/leaf/chanrpc"
	"mj/gameServer/user"
	"github.com/lovelly/leaf/log"
)


type Module interface {
	GetChanRPC() *chanrpc.Server
	OnInit()
	OnDestroy()
	Run(closeSig chan bool)
	GetClientCount() int
	GetTableCount() int
}


/// 房间里面的玩家管理
type RoomInfo struct {
	EendTime int64
	SiceCount int									//骰子点数
	BankerUser int									//庄家用户
	Ting []bool;									//是否听牌
	/*
		todo 桌子属性 基础牌操作
	 */
	Users []*user.User /// index is chairId
}

func (r *RoomInfo) CheckDestroy(curTime int64) bool {
	if len(r.Users) < 1 {
		return true //没人关闭房间 todo
	}

	if r.EendTime <  curTime {
		return true //时间到了关闭房间 todo
	}
	return false
}

func (r *RoomInfo) GetUserCount() int {
	return len(r.Users)
}

func (r *RoomInfo) IsInRoom(userId int) bool {
	for _, u := range r.Users {
		if u.Id == userId {
			return true
		}
	}
	return false
}

func (r *RoomInfo) GetUserByChairId(chairId int)*user.User {
	if len(r.Users) <= chairId  {
		return nil
	}
	return r.Users[chairId]
}


func (r *RoomInfo) GetUserByUid(userId int)*user.User {
	for _, u := range r.Users {
		if u.Id == userId {
			return u
		}
	}
	return nil
}

func (r *RoomInfo) EnterRoom(chairId int, u *user.User) bool {
	if len(r.Users) <= chairId  {
		return false
	}
	ok := !r.IsInRoom(u.Id)
	if ok {
		log.Debug("%v user is inroom,", u.Id)
		return true
	}

	old := r.Users[chairId]
	if old != nil {
		log.Error("at chair %d have user",chairId)
		return false
	}
	r.Users[chairId] = u
	return ok
}

func (r *RoomInfo) LeaveRoom(chairId int) bool {
	if len(r.Users) <= chairId  {
		return false
	}

	r.Users[chairId] = nil
	log.Debug("%v user leave room,  left %v count", chairId, r.GetUserCount())
	return true
}


func (r *RoomInfo) SendMsg(chairId int,data interface{}) bool {
	if len(r.Users) <= chairId  {
		return false
	}

	u := r.Users[chairId]
	if u == nil {
		return false
	}

	u.WriteMsg(data)
	return true
}

func (r *RoomInfo) SendMsgAll(data interface{})  {
	for _, u := range r.Users {
		if u != nil {
			u.WriteMsg(data)
		}
	}
}


