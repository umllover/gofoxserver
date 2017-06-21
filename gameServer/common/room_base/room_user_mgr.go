package room_base

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/Chat"
	"mj/gameServer/common"
	"mj/gameServer/conf"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"mj/gameServer/idGenerate"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/timer"
)

func NewRoomUserMgr(roomId, UserCnt int, Temp *base.GameServiceOption) common.UserManager {
	r := new(RoomUserMgr)
	r.UserCnt = UserCnt
	r.id = roomId
	r.Users = make([]*user.User, r.UserCnt)
	r.Trustee = make([]bool, r.UserCnt)
	r.Onlookers = make(map[int]*user.User)
	return r
}

type RoomUserMgr struct {
	id          int                  //唯一id 房间id
	Kind        int                  //模板表第一类型
	ServerId    int                  //模板表第二类型 注意 非房间id
	EendTime    int64                //结束时间
	UserCnt     int                  //可以容纳的用户数量
	PlayerCount int                  //当前用户人数
	JoinCount   int                  //房主设置的游戏人数
	Users       []*user.User         /// index is chairId
	Onlookers   map[int]*user.User   /// 旁观的玩家
	ChatRoomId  int                  //聊天房间id
	KickOut     map[int]*timer.Timer //即将被踢出的超时定时器
	Trustee     []bool               //是否托管 index 就是椅子id
}

func (r *RoomUserMgr) GetTrustees() []bool {
	return r.Trustee
}

func (r *RoomUserMgr) SetUsetTrustee(chairId int, isTruste bool) {
	r.Trustee[chairId] = isTruste
}

func (r *RoomUserMgr) IsTrustee(chairId int) bool {
	return r.Trustee[chairId]
}

func (r *RoomUserMgr) GetCurPlayerCnt() int {
	r.PlayerCount = 0
	for _, u := range r.Users {
		if u != nil {
			r.PlayerCount++
		}
	}
	return r.PlayerCount
}

func (r *RoomUserMgr) GetMaxPlayerCnt() int {
	return r.PlayerCount
}

func (r *RoomUserMgr) IsInRoom(userId int) bool {
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

func (r *RoomUserMgr) GetUserByChairId(chairId int) *user.User {
	if len(r.Users) <= chairId {
		return nil
	}
	return r.Users[chairId]
}

func (r *RoomUserMgr) GetUserByUid(userId int) (*user.User, int) {
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

func (r *RoomUserMgr) EnterRoom(chairId int, u *user.User) bool {
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

func (r *RoomUserMgr) GetChairId() int {
	for i, u := range r.Users {
		if u == nil {
			return i
		}
	}
	return -1
}

func (r *RoomUserMgr) LeaveRoom(u *user.User) bool {
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

func (r *RoomUserMgr) SendMsg(chairId int, data interface{}) bool {
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

func (r *RoomUserMgr) SendMsgAll(data interface{}) {
	for _, u := range r.Users {
		if u != nil {
			u.WriteMsg(data)
		}
	}
}

func (r *RoomUserMgr) SendOnlookers(data interface{}) {
	for _, u := range r.Onlookers {
		if u != nil {
			u.WriteMsg(data)
		}
	}
}

func (r *RoomUserMgr) SendMsgAllNoSelf(selfid int, data interface{}) {
	for _, u := range r.Users {
		log.Debug("SendMsgAllNoSelf %v ", (u != nil && u.Id != selfid))
		if u != nil && u.Id != selfid {
			u.WriteMsg(data)
		}
	}
}

func (r *RoomUserMgr) ForEachUser(fn func(u *user.User)) {
	for _, u := range r.Users {
		if u != nil {
			fn(u)
		}
	}
}

func (r *RoomUserMgr) WriteTableScore(source []*msg.TagScoreInfo, usercnt, Type int) {
	for _, u := range r.Users {
		u.ChanRPC().Go("WriteUserScore", source[u.ChairId], Type)
	}
}

//获取对方信息
func (room *RoomUserMgr) GetUserInfoByChairId(ChairID int) interface{} {
	tagUser := room.GetUserByChairId(ChairID)
	if tagUser == nil {
		log.Error("at GetUserChairInfo no foud tagUser %v", ChairID)
		return nil
	}

	return &msg.G2C_UserEnter{
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
	}
}

//坐下
func (room *RoomUserMgr) Sit(u *user.User, ChairID int) int {
	oldUser := room.GetUserByChairId(ChairID)
	if oldUser != nil {
		return ChairHasUser
	}

	if room.ChatRoomId == 0 {
		id, err := Chat.ChanRPC.Call1("createRoom", u.Agent)
		if err != nil {
			log.Error("create Chat Room faild")
			return ErrCreateRoomFaild
		}
		room.ChatRoomId = id.(int)
	}

	_, chairId := room.GetUserByUid(u.Id)
	if chairId > 0 {
		room.LeaveRoom(u)
	}

	room.EnterRoom(ChairID, u)
	//把自己的信息推送给所有玩家
	room.NotifyUserInfo(u)
	//把所有玩家信息推送给自己
	room.GetAllUsetInfo(u)

	Chat.ChanRPC.Go("addRoomMember", room.ChatRoomId, u.Agent)
	room.SetUsetStatus(u, US_SIT)
	return 0
}

//广播某个玩家的信息
func (room *RoomUserMgr) NotifyUserInfo(u *user.User) {
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
}

func (room *RoomUserMgr) GetAllUsetInfo(u *user.User) {
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
}

//起立
func (room *RoomUserMgr) Standup(u *user.User) bool {
	room.SetUsetStatus(u, US_FREE)
	room.LeaveRoom(u)
	return true
}

func (room *RoomUserMgr) SetUsetStatus(u *user.User, stu int) {
	u.Status = stu
	room.SendMsgAll(&msg.G2C_UserStatus{
		UserID: u.Id,
		UserStatus: &msg.UserStu{
			TableID:    room.id,
			ChairID:    u.ChairId,
			UserStatus: u.Status,
		},
	})
}

//通知房间解散
func (room *RoomUserMgr) RoomDissume() {
	Cance := &msg.G2C_CancelTable{}
	room.ForEachUser(func(u *user.User) {
		u.WriteMsg(Cance)
	})

	Diis := &msg.G2C_PersonalTableEnd{}
	room.ForEachUser(func(u *user.User) {
		u.WriteMsg(Diis)
	})

	for _, u := range room.Users {
		if u != nil {
			u.ChanRPC().Go("LeaveRoom")
		}
	}

	idGenerate.DelRoomId(room.id)
}

func (room *RoomUserMgr) IsAllReady() bool {
	for _, u := range room.Users {
		if u == nil || u.Status != US_READY {
			return false
		}
	}
	return true
}

func (room *RoomUserMgr) ReLogin(u *user.User, Status int) {
	tm, ok := room.KickOut[u.Id]
	if ok {
		tm.Stop()
		delete(room.KickOut, u.Id)
	}

	if Status == RoomStatusStarting {
		room.SetUsetStatus(u, US_PLAYING)
	} else {
		room.SetUsetStatus(u, US_SIT)
	}
}

func (room *RoomUserMgr) AddKickOutTimer(uid int, t *timer.Timer) {
	room.KickOut[uid] = t
}

func (room *RoomUserMgr) SendUserInfoToSelf(u *user.User) {
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
}
