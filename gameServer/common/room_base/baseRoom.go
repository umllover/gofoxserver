package room_base

import (
	"mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/user"
	"time"

	"github.com/lovelly/leaf/chanrpc"
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
	id          int    //唯一id 房间id
	EendTime    int64  //结束时间
	SiceCount   int    //骰子点数
	BankerUser  int    //庄家用户
	CurrentUser int    //当前操作用户
	Ting        []bool //是否听牌
	CreateUser  int    //创建房间的人
	Status      int    //当前状态
	PlayCount   int    //已玩局数
	CreateTime  int64  //创建时间
	UserCnt     int    //可以容纳的用户数量

	Users        []*user.User /// index is chairId
	AllowLookon  map[int]int  //旁观玩家
	TurnScore    []int        //积分信息
	CollectScore []int        //积分信息
	Trustee      []bool       //是否托管 index 就是椅子id
}

func NewRoomInfo(userCnt, rid int) *RoomInfo {
	r := new(RoomInfo)
	r.id = rid
	r.Users = make([]*user.User, userCnt)
	r.AllowLookon = make(map[int]int)
	r.CreateTime = time.Now().Unix()
	r.TurnScore = make([]int, userCnt)
	r.CollectScore = make([]int, userCnt)
	r.Trustee = make([]bool, userCnt)
	r.UserCnt = userCnt
	return r
}

func (r *RoomInfo) GetRoomId() int {
	return r.id
}

func (r *RoomInfo) SetRoomStatus(su int) {
	r.Status = su
}

func (r *RoomInfo) GetRoomStatus() int {
	return r.Status
}

func (r *RoomInfo) CheckDestroy(curTime int64) bool {
	if len(r.Users) < 1 {
		return true //没人关闭房间 todo
	}

	if r.EendTime < curTime {
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

func (r *RoomInfo) GetUserByChairId(chairId int) *user.User {
	if len(r.Users) <= chairId {
		return nil
	}
	return r.Users[chairId]
}

func (r *RoomInfo) GetUserByUid(userId int) (*user.User, int) {
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
		log.Error("at chair %d have user", chairId)
		return false
	}
	r.Users[chairId] = u
	u.ChairId = chairId
	u.RoomId = r.id
	return true
}

func (r *RoomInfo) GetChairId() int {
	for i, u := range r.Users {
		if u == nil {
			return i
		}
	}
	return -1
}

func (r *RoomInfo) LeaveRoom(u *user.User) bool {
	if len(r.Users) <= u.ChairId {
		return false
	}

	r.Users[u.ChairId] = nil
	u.ChairId = cost.INVALID_CHAIR
	u.RoomId = cost.INVALID_TABLE
	log.Debug("%v user leave room,  left %v count", u.ChairId, r.GetUserCount())
	return true
}

func (r *RoomInfo) SendMsg(chairId int, data interface{}) bool {
	if len(r.Users) <= chairId {
		return false
	}

	u := r.Users[chairId]
	if u == nil {
		return false
	}

	u.WriteMsg(data)
	return true
}

func (r *RoomInfo) SendMsgAll(data interface{}) {
	for _, u := range r.Users {
		if u != nil {
			u.WriteMsg(data)
		}
	}
}

func (r *RoomInfo) SendMsgAllNoSelf(selfid int, data interface{}) {
	for _, u := range r.Users {
		log.Debug("SendMsgAllNoSelf %v ", (u != nil && u.Id != selfid))
		if u != nil && u.Id != selfid {
			u.WriteMsg(data)
		}
	}
}

func (r *RoomInfo) ForEachUser(fn func(u *user.User)) {
	for _, u := range r.Users {
		if u != nil {
			fn(u)
		}
	}
}

func (r *RoomInfo) WriteTableScore(source []*msg.TagScoreInfo, usercnt, Type int) {
	for _, u := range r.Users {
		u.ChanRPC().Go("WriteUserScore", source[u.ChairId], Type)
	}
}
