package common

import (
	"github.com/lovelly/leaf/chanrpc"
	"mj/gameServer/user"
	"github.com/lovelly/leaf/log"
	"time"
	"mj/common/cost"
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
	CurrentUser int									//当前操作用户
	Ting []bool;									//是否听牌
	CreateUser  int										//创建房间的人
	Status int								//当前状态
	PlayCount int 					//已玩局数
	CreateTime int64
	UserCnt int8
	/*
		todo 桌子属性 基础牌操作
	 */
	Users []*user.User /// index is chairId
	AllowLookon map[int]int //旁观玩家
	TurnScore []int  //积分信息
	CollectScore []int  //积分信息
	Trustee []bool								//是否托管 index 就是椅子id
}

func NewRoomInfo(userCnt int)*RoomInfo{
	r := new(RoomInfo)
	r.Users =make([]*user.User, userCnt)
	r.AllowLookon = make(map[int]int)
	r.CreateTime = time.Now().Unix()
	r.TurnScore = make([]int, userCnt)
	r.CollectScore = make([]int, userCnt)
	r.Trustee = make([]bool, userCnt)
	r.UserCnt = int8(userCnt)
	return r
}

func (r *RoomInfo)SetRoomStatus(su int){
	r.Status = su
}

func(r *RoomInfo) GetRoomStatus()int{
	return r.Status
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
		if u == nil {
			continue
		}
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

func (r *RoomInfo) GetUserByUid(userId int)(*user.User, int) {
	for i, u := range r.Users {
		if u == nil {
			continue
		}
		if u.Id == userId {
			return u, i
		}
	}
	return nil, -1
}

func (r *RoomInfo) EnterRoom(chairId int, u *user.User) bool {
	if chairId == cost.INVALID_CHAIR {
		chairId = r.GetChairId()
	}
	if len(r.Users) <= chairId || chairId == -1 {
		log.Error("at EnterRoom max chairId, user len :%d, chairId:%d", len(r.Users), chairId)
		return false
	}

	if r.IsInRoom(u.Id) {
		log.Debug("%v user is inroom,", u.Id)
		return true
	}

	old := r.Users[chairId]
	if old != nil {
		log.Error("at chair %d have user",chairId)
		return false
	}
	r.Users[chairId] = u
	u.ChairId = chairId
	return true
}

func  (r *RoomInfo) GetChairId() int {
	for i, u := range r.Users {
		if u == nil {
			return i
		}
	}
	return -1
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


