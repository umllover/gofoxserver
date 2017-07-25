package pk_base

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/nn_tb_msg"
	"mj/gameServer/common/pk"

	"mj/gameServer/common/room_base"
	"mj/gameServer/conf"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"
)

//创建的配置文件
type NewPKCtlConfig struct {
	BaseMgr  room_base.BaseManager
	TimerMgr room_base.TimerManager
	UserMgr  room_base.UserManager
	DataMgr  pk.DataManager
	LogicMgr pk.LogicManager
}

//消息入口文件
type Entry_base struct {
	room_base.BaseManager
	UserMgr  room_base.UserManager
	TimerMgr room_base.TimerManager

	DataMgr  pk.DataManager
	LogicMgr pk.LogicManager

	Temp   *base.GameServiceOption //模板
	Status int

	BtCardSpecialData []int
}

func NewPKBase(info *model.CreateRoomInfo) *Entry_base {
	Temp, ok1 := base.GameServiceOptionCache.Get(info.KindId, info.ServiceId)
	log.Debug("new pk base %d %d", info.KindId, info.ServiceId)
	if !ok1 {
		log.Error("at NewPKBase not foud config .... ")
		return nil
	}

	pk := new(Entry_base)
	pk.Temp = Temp
	return pk
}

func (r *Entry_base) Init(cfg *NewPKCtlConfig) {
	r.UserMgr = cfg.UserMgr
	r.DataMgr = cfg.DataMgr
	r.BaseManager = cfg.BaseMgr
	r.LogicMgr = cfg.LogicMgr
	r.TimerMgr = cfg.TimerMgr
	r.RoomRun(r.DataMgr.GetRoomId())

	r.DataMgr.OnCreateRoom()

	r.TimerMgr.StartCreatorTimer(r.GetSkeleton(), func() {
		r.OnEventGameConclude(0, nil, GER_DISMISS)
	})

}

func (r *Entry_base) GetRoomId() int {
	return r.DataMgr.GetRoomId()
}

//坐下
func (r *Entry_base) Sitdown(args []interface{}) {
	chairID := args[0].(int)
	u := args[1].(*user.User)

	retcode := 0
	defer func() {
		u.WriteMsg(&msg.G2C_UserSitDownRst{Code: retcode})
	}()
	if r.Status == RoomStatusStarting && r.Temp.DynamicJoin == 1 {
		retcode = GameIsStart
		return
	}

	retcode = r.UserMgr.Sit(u, chairID)

}

//起立
func (r *Entry_base) UserStandup(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserStandup{})
	u := args[1].(*user.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	if r.Status == RoomStatusStarting {
		retcode = ErrGameIsStart
		return
	}

	r.UserMgr.Standup(u)
}

//获取对方信息
func (room *Entry_base) GetUserChairInfo(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_REQUserChairInfo)
	u := args[1].(*user.User)
	info := room.UserMgr.GetUserInfoByChairId(recvMsg.ChairID).(*msg.G2C_UserEnter)
	if info == nil {
		log.Error("at GetUserChairInfo no foud tagUser %v, userId:%d", args[0], u.Id)
		return
	}
	u.WriteMsg(info)
}

//解散房间
func (room *Entry_base) DissumeRoom(args []interface{}) {
	u := args[0].(*user.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode, "解散房间失败."))
		}
	}()

	if !room.DataMgr.CanOperatorRoom(u.Id) {
		retcode = NotOwner
		return
	}

	room.UserMgr.ForEachUser(func(u *user.User) {
		room.UserMgr.LeaveRoom(u, room.Status)
	})

	room.OnEventGameConclude(0, nil, GER_DISMISS)
	room.Destroy(room.DataMgr.GetRoomId())
}

//玩家准备
func (room *Entry_base) UserReady(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserReady)
	u := args[1].(*user.User)
	if u.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	log.Debug("at UserReady")
	room.UserMgr.SetUsetStatus(u, US_READY)

	if room.UserMgr.IsAllReady() {
		log.Debug("all user are ready start game")
		//派发初始扑克
		room.DataMgr.BeforeStartGame(room.UserMgr.GetMaxPlayerCnt())
		room.DataMgr.StartGameing()
		room.DataMgr.AfterStartGame()

		room.Status = RoomStatusStarting
		room.TimerMgr.StartPlayingTimer(room.GetSkeleton(), func() {
			room.OnEventGameConclude(0, nil, GER_DISMISS)
		})
	}
}

//玩家重登
func (room *Entry_base) UserReLogin(args []interface{}) {
	u := args[0].(*user.User)
	roomUser := room.getRoomUser(u.Id)
	if roomUser == nil {
		return
	}
	log.Debug("at ReLogin have old user ")
	u.ChairId = roomUser.ChairId
	u.RoomId = roomUser.RoomId
	room.UserMgr.ReLogin(u, room.Status)
}

//玩家离线
func (room *Entry_base) UserOffline(args []interface{}) {
	u := args[0].(*user.User)
	log.Debug("at UserOffline .... uid:%d", u.Id)
	if u.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	room.UserMgr.SetUsetStatus(u, US_OFFLINE)
	room.TimerMgr.StartKickoutTimer(room.GetSkeleton(), u.Id, func() {
		room.OffLineTimeOut(u)
	})
}

//离线超时踢出
func (room *Entry_base) OffLineTimeOut(u *user.User) {
	room.UserMgr.LeaveRoom(u, room.Status)
	if room.Status != RoomStatusReady {
		room.OnEventGameConclude(0, nil, GER_DISMISS)
	} else {
		if room.UserMgr.GetCurPlayerCnt() == 0 { //没人了直接销毁
			room.Destroy(room.DataMgr.GetRoomId())
		}
	}
}

//获取房间基础信息
func (room *Entry_base) GetBirefInfo() *msg.RoomInfo {
	BirefInf := &msg.RoomInfo{}
	BirefInf.ServerID = room.Temp.ServerID
	BirefInf.KindID = room.Temp.KindID
	BirefInf.NodeID = conf.Server.NodeId
	BirefInf.SvrHost = conf.Server.WSAddr
	BirefInf.PayType = room.UserMgr.GetPayType()
	BirefInf.RoomID = room.DataMgr.GetRoomId()
	BirefInf.CurCnt = room.UserMgr.GetCurPlayerCnt()
	BirefInf.MaxPlayerCnt = room.UserMgr.GetMaxPlayerCnt() //最多多人数
	BirefInf.PayCnt = room.TimerMgr.GetMaxPayCnt()         //可玩局数
	BirefInf.CurPayCnt = room.TimerMgr.GetPlayCount()      //已玩局数
	BirefInf.CreateTime = room.TimerMgr.GetCreatrTime()    //创建时间
	//BirefInf.CreateUserId = room.DataMgr.GetCreater()
	BirefInf.IsPublic = room.UserMgr.IsPublic()
	BirefInf.Players = make(map[int64]*msg.PlayerBrief)
	BirefInf.MachPlayer = make(map[int64]struct{})
	return BirefInf
}

//游戏配置
func (room *Entry_base) SetGameOption(args []interface{}) {
	u := args[1].(*user.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	if u.ChairId == INVALID_CHAIR {
		retcode = ErrNoSitdowm
		return
	}

	AllowLookon := 0
	if u.Status == US_LOOKON {
		AllowLookon = 1
	}
	u.WriteMsg(&msg.G2C_GameStatus{
		GameStatus:  room.Status,
		AllowLookon: AllowLookon,
	})

	room.DataMgr.SendPersonalTableTip(u)

	if room.Status == RoomStatusReady { // 没开始
		room.DataMgr.SendStatusReady(u)
	} else { //开始了
		//把所有玩家信息推送给自己
		room.UserMgr.SendUserInfoToSelf(u)
		room.DataMgr.SendStatusPlay(u)
	}
}

//游戏结束
func (room *Entry_base) OnEventGameConclude(ChairId int, user *user.User, cbReason int) {
	switch cbReason {
	case GER_NORMAL: //常规结束
		room.DataMgr.NormalEnd()
		//room.AfertEnd(false)// 这里需要重构 不同房间结束不一样
		room.DataMgr.AfterEnd(false)
		return
	case GER_DISMISS: //游戏解散
		room.DataMgr.DismissEnd()
		room.AfertEnd(true)
	}
	log.Error("at OnEventGameConclude error  ")
	return
}

// 如果这里不能满足 afertEnd 请重构这个到个个组件里面
func (room *Entry_base) AfertEnd(Forced bool) {
	room.TimerMgr.AddPlayCount()
	if Forced || room.TimerMgr.GetPlayCount() >= room.TimerMgr.GetMaxPayCnt() {
		log.Debug("Forced :%v, PlayTurnCount:%v, temp PlayTurnCount:%d", Forced, room.TimerMgr.GetPlayCount(), room.TimerMgr.GetMaxPayCnt())
		room.UserMgr.SendMsgToHallServerAll(&msg.RoomEndInfo{
			RoomId: room.DataMgr.GetRoomId(),
			Status: room.Status,
		})
		room.Destroy(room.DataMgr.GetRoomId())
		room.UserMgr.RoomDissume()
		return
	}

	room.UserMgr.ForEachUser(func(u *user.User) {
		room.UserMgr.SetUsetStatus(u, US_FREE)
	})
}

//计算税收 暂时未实现
func (room *Entry_base) CalculateRevenue(ChairId, lScore int) int {
	//效验参数

	UserCnt := room.UserMgr.GetMaxPlayerCnt()
	if ChairId >= UserCnt {
		return 0
	}

	return 0
}

//叫分(倍数)
func (room *Entry_base) CallScore(args []interface{}) {
	recvMsg := args[0].(*nn_tb_msg.C2G_TBNN_CallScore)
	u := args[1].(*user.User)

	room.DataMgr.CallScore(u, recvMsg.CallScore)
	return
}

//加注
func (r *Entry_base) AddScore(args []interface{}) {
	recvMsg := args[0].(*nn_tb_msg.C2G_TBNN_AddScore)
	u := args[1].(*user.User)

	r.DataMgr.AddScore(u, recvMsg.Score)
	return
}

// 亮牌
func (r *Entry_base) OpenCard(args []interface{}) {
	recvMsg := args[0].(*nn_tb_msg.C2G_TBNN_OpenCard)
	u := args[1].(*user.User)

	r.DataMgr.OpenCard(u, recvMsg.CardType, recvMsg.CardData)
	return
}

// 十三水摊牌
func (r *Entry_base) ShowSSsCard(args []interface{}) {
	//recvMsg := args[0].(*pk_sss_msg.C2G_SSS_Open_Card)
	//u := args[1].(*user.User)

	//r.DataMgr.ShowSSSCard(u, recvMsg.Dragon, recvMsg.SpecialType, recvMsg.SpecialData, recvMsg.FrontCard, recvMsg.MidCard, recvMsg.BackCard)
	return
}

func (r *Entry_base) getRoomUser(uid int64) *user.User {
	u, _ := r.UserMgr.GetUserByUid(uid)
	return u
}
