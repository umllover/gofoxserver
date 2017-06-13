package room_base

import (
	"mj/common/base"
	"mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/user"
	"sync"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
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
type RoomBase struct {
	// module 必须字段
	*module.Skeleton
	ChanRPC  *chanrpc.Server //接受客户端消息的chan
	MgrCh    *chanrpc.Server //管理类的chan 例如红中麻将 就是红中麻将module的 ChanRPC
	CloseSig chan bool
	wg       sync.WaitGroup //

	id                  int                //唯一id 房间id
	Name                string             //房间名字
	EendTime            int64              //结束时间
	CreateTime          int64              //创建时间
	UserCnt             int                //可以容纳的用户数量
	PlayerCount         int                //当前用户人数
	JoinGamePeopleCount int                //房主设置的游戏人数
	Users               []*user.User       /// index is chairId
	Onlookers           map[int]*user.User /// 旁观的玩家
	CreateUser          int                //创建房间的人
}

func NewRoomBase(userCnt, rid int, mgrCh *chanrpc.Server, name string) *RoomBase {
	r := new(RoomBase)
	skeleton := base.NewSkeleton()
	r.Skeleton = skeleton
	r.ChanRPC = skeleton.ChanRPCServer
	r.MgrCh = mgrCh
	r.Name = name
	r.id = rid
	r.Users = make([]*user.User, userCnt)
	r.CreateTime = time.Now().Unix()
	r.UserCnt = userCnt
	r.Onlookers = make(map[int]*user.User)
	return r
}

func (r *RoomBase) RoomRun() {
	go func() {
		log.Debug("room Room start run Name:%s", r.id)
		r.Run(r.CloseSig)
		r.End()
		log.Debug("room Room End run Name:%s", r.id)
	}()
}

func (r *RoomBase) End() {
	for _, u := range r.Users {
		u.ChanRPC().Go("RoomClose")
	}
}

func (r *RoomBase) Destroy() {
	defer func() {
		if r := recover(); r != nil {
			log.Recover(r)
		}
	}()

	r.CloseSig <- true
	r.OnDestroy()
	log.Debug("room Room Destroy ok,  Name:%s", r.id)
}

func (r *RoomBase) OnDestroy() { // 基类实现

}

func (r *RoomBase) GetChanRPC() *chanrpc.Server {
	return r.ChanRPC
}

func (r *RoomBase) GetRoomId() int {
	return r.id
}

func (r *RoomBase) CheckDestroy(curTime int64) bool {
	if len(r.Users) < 1 {
		return true //没人关闭房间 todo
	}

	if r.EendTime < curTime {
		return true //时间到了关闭房间 todo
	}
	return false
}

func (r *RoomBase) GetUserCount() int {
	return len(r.Users)
}

func (r *RoomBase) IsInRoom(userId int) bool {
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

func (r *RoomBase) GetUserByChairId(chairId int) *user.User {
	if len(r.Users) <= chairId {
		return nil
	}
	return r.Users[chairId]
}

func (r *RoomBase) GetUserByUid(userId int) (*user.User, int) {
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

func (r *RoomBase) EnterRoom(chairId int, u *user.User) bool {
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

func (r *RoomBase) GetChairId() int {
	for i, u := range r.Users {
		if u == nil {
			return i
		}
	}
	return -1
}

func (r *RoomBase) LeaveRoom(u *user.User) bool {
	if len(r.Users) <= u.ChairId {
		return false
	}

	r.Users[u.ChairId] = nil
	u.ChairId = cost.INVALID_CHAIR
	u.RoomId = cost.INVALID_TABLE
	log.Debug("%v user leave room,  left %v count", u.ChairId, r.GetUserCount())
	return true
}

func (r *RoomBase) SendMsg(chairId int, data interface{}) bool {
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

func (r *RoomBase) SendMsgAll(data interface{}) {
	for _, u := range r.Users {
		if u != nil {
			u.WriteMsg(data)
		}
	}
}

func (r *RoomBase) SendOnlookers(data interface{}) {
	for _, u := range r.Onlookers {
		if u != nil {
			u.WriteMsg(data)
		}
	}
}

func (r *RoomBase) SendMsgAllNoSelf(selfid int, data interface{}) {
	for _, u := range r.Users {
		log.Debug("SendMsgAllNoSelf %v ", (u != nil && u.Id != selfid))
		if u != nil && u.Id != selfid {
			u.WriteMsg(data)
		}
	}
}

func (r *RoomBase) ForEachUser(fn func(u *user.User)) {
	for _, u := range r.Users {
		if u != nil {
			fn(u)
		}
	}
}

func (r *RoomBase) WriteTableScore(source []*msg.TagScoreInfo, usercnt, Type int) {
	for _, u := range r.Users {
		u.ChanRPC().Go("WriteUserScore", source[u.ChairId], Type)
	}
}
