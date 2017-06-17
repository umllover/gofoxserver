package room_base

import (
	"mj/common/base"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/Chat"
	"mj/gameServer/conf"
	"mj/gameServer/db/model"
	tbase "mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"sync"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/timer"
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
	Kind                int                //模板表第一类型
	ServerId            int                //模板表第二类型 注意 非房间id
	Name                string             //房间名字
	EendTime            int64              //结束时间
	CreateTime          int64              //创建时间
	UserCnt             int                //可以容纳的用户数量
	PlayerCount         int                //当前用户人数
	JoinGamePeopleCount int                //房主设置的游戏人数
	Users               []*user.User       /// index is chairId
	Onlookers           map[int]*user.User /// 旁观的玩家
	CreateUser          int                //创建房间的人
	Owner               int                //房主id
	Status              int                //当前状态
	ChatRoomId          int                //聊天房间id
	TimeStartGame       int64              //开始时间
	MaxPayCnt           int                //最大局数
	PlayCount           int                //已玩局数
	KickOut             map[int]*timer.Timer
	Temp                *tbase.GameServiceOption
	//cb
	StartGameCb    func() //开始函数
	GameConcludeCb func(ChairId int, user *user.User, cbReason int) bool
	OnDestroyCb    func() //销毁回调函数
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
		defer func() {
			if r := recover(); r != nil {
				log.Recover(r)
			}
		}()

		log.Debug("room Room start run Name:%s", r.Name)
		r.Run(r.CloseSig)
		r.End()
		log.Debug("room Room End run Name:%s", r.Name)
	}()
}

func (r *RoomBase) GetCurlPlayerCount() int {
	cnt := 0
	for _, u := range r.Users {
		if u != nil {
			cnt++
		}
	}

	return cnt
}

func (r *RoomBase) Destroy() {
	defer func() {
		if r := recover(); r != nil {
			log.Recover(r)
		}
	}()
	r.OnDestroyCb()
	r.CloseSig <- true
	log.Debug("room Room Destroy ok,  Name:%s", r.Name)
}

func (r *RoomBase) End() {
	for _, u := range r.Users {
		if u != nil {
			u.ChanRPC().Go("LeaveRoom")
		}
	}
}

func (r *RoomBase) GetBirefInfo() *msg.RoomInfo {
	msg := &msg.RoomInfo{}
	msg.ServerID = r.ServerId
	msg.KindID = r.Kind
	msg.NodeID = conf.Server.NodeId
	msg.RoomID = r.GetRoomId()
	msg.CurCnt = r.PlayerCount
	msg.MaxCnt = r.UserCnt           //最多多人数
	msg.PayCnt = r.MaxPayCnt         //可玩局数
	msg.CurPayCnt = r.PlayCount      //已玩局数
	msg.CreateTime = r.TimeStartGame //创建时间
	return msg
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
	if chairId == INVALID_CHAIR {
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

	locak := &model.Gamescorelocker{}
	locak.UserID = u.Id
	locak.KindID = u.KindID
	locak.ServerID = u.ServerID
	locak.NodeID = conf.Server.NodeId
	_, err := model.GamescorelockerOp.Insert(locak)
	if err != nil {
		log.Error("at EnterRoom  updaye .Gamescorelocker error:%s", err.Error())
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
	err := model.GamescorelockerOp.Delete(u.Id)
	if err != nil {
		log.Error("at EnterRoom  updaye .Gamescorelocker error:%s", err.Error())
	}
	u.ChanRPC().Go("LeaveRoom")
	r.Users[u.ChairId] = nil
	u.ChairId = INVALID_CHAIR
	u.RoomId = 0
	log.Debug("%v user leave room,  left %v count", u.ChairId, r.PlayerCount)

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

func (r *RoomBase) CheckPlayerCnt() bool {
	r.PlayerCount = 0
	for _, u := range r.Users {
		if u != nil {
			r.PlayerCount++
		}
	}
	return r.PlayerCount == 0
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

//// base msg =======================

func (room *RoomBase) DissumeRoom(args []interface{}) {
	u := args[0].(*user.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode, "解散房间失败."))
		}
	}()
	if u.Id != room.Owner {
		retcode = NotOwner
		return
	}

	Cance := &msg.G2C_CancelTable{}
	room.ForEachUser(func(u *user.User) {
		u.WriteMsg(Cance)
	})

	Diis := &msg.G2C_PersonalTableEnd{}
	room.ForEachUser(func(u *user.User) {
		u.WriteMsg(Diis)
	})

	room.Destroy()
}

//获取对方信息
func (room *RoomBase) GetUserChairInfo(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_REQUserChairInfo)
	u := args[1].(*user.User)
	tagUser := room.GetUserByChairId(recvMsg.ChairID)
	if tagUser == nil {
		log.Error("at GetUserChairInfo no foud tagUser %v, userId:%d", args[0], u.Id)
		return
	}

	u.WriteMsg(&msg.G2C_UserEnter{
		UserID:      tagUser.Id,          //用户 I D
		FaceID:      tagUser.FaceID,      //头像索引
		CustomID:    tagUser.CustomID,    //自定标识
		Gender:      tagUser.Gender,      //用户性别
		MemberOrder: tagUser.MemberOrder, //会员等级
		TableID:     tagUser.RoomId,      //桌子索引
		ChairID:     tagUser.ChairId,     //椅子索引
		UserStatus:  tagUser.Status,      //用户状态
		Score:       tagUser.Score,       //用户分数
		WinCount:    tagUser.WinCount,    //胜利盘数
		LostCount:   tagUser.LostCount,   //失败盘数
		DrawCount:   tagUser.DrawCount,   //和局盘数
		FleeCount:   tagUser.FleeCount,   //逃跑盘数
		Experience:  tagUser.Experience,  //用户经验
		NickName:    tagUser.NickName,    //昵称
		HeaderUrl:   tagUser.HeadImgUrl,  //头像
	})
}

//坐下
func (room *RoomBase) Sitdown(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_UserSitdown)
	u := args[1].(*user.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	oldUser := room.GetUserByChairId(recvMsg.ChairID)
	if oldUser != nil {
		retcode = ChairHasUser
		return
	}

	template, ok := tbase.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	if !ok {
		retcode = ConfigError
		return
	}

	if room.Status == RoomStatusStarting && template.DynamicJoin == 1 {
		retcode = GameIsStart
		return
	}

	if room.ChatRoomId == 0 {
		id, err := Chat.ChanRPC.Call1("createRoom", u.Agent)
		if err != nil {
			log.Error("create Chat Room faild")
			retcode = ErrCreateRoomFaild
		}

		room.ChatRoomId = id.(int)
	}

	_, chairId := room.GetUserByUid(u.Id)
	if chairId > 0 {
		room.LeaveRoom(u)
	}

	room.EnterRoom(recvMsg.ChairID, u)
	//把自己的信息推送给所有玩家
	room.SendMsgAllNoSelf(u.Id, &msg.G2C_UserEnter{
		UserID:      u.Id,          //用户 I D
		FaceID:      u.FaceID,      //头像索引
		CustomID:    u.CustomID,    //自定标识
		Gender:      u.Gender,      //用户性别
		MemberOrder: u.MemberOrder, //会员等级
		TableID:     u.RoomId,      //桌子索引
		ChairID:     u.ChairId,     //椅子索引
		UserStatus:  u.Status,      //用户状态
		Score:       u.Score,       //用户分数
		WinCount:    u.WinCount,    //胜利盘数
		LostCount:   u.LostCount,   //失败盘数
		DrawCount:   u.DrawCount,   //和局盘数
		FleeCount:   u.FleeCount,   //逃跑盘数
		Experience:  u.Experience,  //用户经验
		NickName:    u.NickName,    //昵称
		HeaderUrl:   u.HeadImgUrl,  //头像
	})

	//把所有玩家信息推送给自己
	room.ForEachUser(func(eachuser *user.User) {
		if eachuser.Id == u.Id {
			return
		}
		u.WriteMsg(&msg.G2C_UserEnter{
			UserID:      eachuser.Id,          //用户 I D
			FaceID:      eachuser.FaceID,      //头像索引
			CustomID:    eachuser.CustomID,    //自定标识
			Gender:      eachuser.Gender,      //用户性别
			MemberOrder: eachuser.MemberOrder, //会员等级
			TableID:     eachuser.RoomId,      //桌子索引
			ChairID:     eachuser.ChairId,     //椅子索引
			UserStatus:  eachuser.Status,      //用户状态
			Score:       eachuser.Score,       //用户分数
			WinCount:    eachuser.WinCount,    //胜利盘数
			LostCount:   eachuser.LostCount,   //失败盘数
			DrawCount:   eachuser.DrawCount,   //和局盘数
			FleeCount:   eachuser.FleeCount,   //逃跑盘数
			Experience:  eachuser.Experience,  //用户经验
			NickName:    eachuser.NickName,    //昵称
			HeaderUrl:   eachuser.HeadImgUrl,  //头像
		})
	})

	Chat.ChanRPC.Go("addRoomMember", room.ChatRoomId, u.Agent)
	room.SetUsetStatus(u, US_SIT)
}

//起立
func (room *RoomBase) UserStandup(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserStandup{})
	u := args[1].(*user.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	if room.Status == RoomStatusStarting {
		retcode = ErrGameIsStart
		return
	}

	room.SetUsetStatus(u, US_FREE)
	room.LeaveRoom(u)
}

func (room *RoomBase) UserReady(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserReady)
	u := args[1].(*user.User)
	if u.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	room.SetUsetStatus(u, US_READY)
	if room.IsAllReady() {
		room.StartGameCb()
	}
}

func (room *RoomBase) UserReLogin(args []interface{}) {
	u := args[0].(*user.User)
	if u.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	tm, ok := room.KickOut[u.Id]
	if ok {
		tm.Stop()
		delete(room.KickOut, u.Id)
	}

	if room.Status == RoomStatusStarting {
		room.SetUsetStatus(u, US_PLAYING)
	} else {
		room.SetUsetStatus(u, US_SIT)
	}
}

func (room *RoomBase) UserOffline(args []interface{}) {
	u := args[0].(*user.User)
	if u.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	room.SetUsetStatus(u, US_OFFLINE)
	if room.Temp.TimeOffLineCount != 0 {
		room.KickOut[u.Id] = room.Skeleton.AfterFunc(time.Duration(room.Temp.TimeOffLineCount)*time.Second, func() {
			room.OfflineKickOut(u)
		})
	} else {
		room.OfflineKickOut(u)
	}
}

// help
func (room *RoomBase) SetUsetStatus(u *user.User, stu int) {
	u.Status = stu
	room.SendMsgAll(&msg.G2C_UserStatus{
		UserID: u.Id,
		UserStatus: &msg.UserStu{
			TableID:    room.GetRoomId(),
			ChairID:    u.ChairId,
			UserStatus: u.Status,
		},
	})
}

func (room *RoomBase) IsAllReady() bool {
	for _, u := range room.Users {
		if u == nil || u.Status != US_READY {
			return false
		}
	}
	return true
}

//玩家离线超时踢出
func (room *RoomBase) OfflineKickOut(user *user.User) {
	room.LeaveRoom(user)
	if room.Status != RoomStatusReady {
		room.GameConcludeCb(0, nil, GER_DISMISS)
	} else {
		if room.CheckPlayerCnt() {
			room.Destroy()
		}
	}
}
