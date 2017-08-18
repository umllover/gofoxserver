package room_base

import (
	"errors"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/mj_zp_msg"
	"mj/gameServer/RoomMgr"
	. "mj/gameServer/common"
	"mj/gameServer/conf"
	"mj/gameServer/db/model/base"
	datalog "mj/gameServer/log"
	"mj/gameServer/user"
	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
	"github.com/lovelly/leaf/timer"
)

type Entry_base struct {
	BaseManager
	UserMgr  UserManager
	TimerMgr TimerManager
	DataMgr  BData

	OnUserTrusteeCb  func(wChairID int, bTrustee bool) bool
	Temp             *base.GameServiceOption //模板
	Status           int
	IsClose          bool
	DelayCloseTimer  *timer.Timer
	RoomTrusteeTimer *timer.Timer
}

//创建的配置文件
type NewMjCtlConfig struct {
	BaseMgr  BaseManager
	UserMgr  UserManager
	TimerMgr TimerManager
	DataMgr  BData
}

func NewEntryBase(KindId, ServiceId int) *Entry_base {
	Temp, ok1 := base.GameServiceOptionCache.Get(KindId, ServiceId)
	if !ok1 {
		log.Error("at NewMJBase not foud template  kindid:%d  serverid:%d", KindId, ServiceId)
		return nil
	}

	mj := new(Entry_base)
	mj.Temp = Temp
	return mj
}

func (r *Entry_base) RegisterBaseFunc() {
	r.GetChanRPC().Register("Sitdown", r.Sitdown)
	r.GetChanRPC().Register("UserStandup", r.UserStandup)
	r.GetChanRPC().Register("GetUserChairInfo", r.GetUserChairInfo)
	r.GetChanRPC().Register("DissumeRoom", r.DissumeRoom)
	r.GetChanRPC().Register("UserReady", r.UserReady)
	r.GetChanRPC().Register("userRelogin", r.UserReLogin)
	r.GetChanRPC().Register("userOffline", r.UserOffline)
	r.GetChanRPC().Register("SetGameOption", r.SetGameOption)
	r.GetChanRPC().Register("ReqLeaveRoom", r.ReqLeaveRoom)
	r.GetChanRPC().Register("ReplyLeaveRoom", r.ReplyLeaveRoom)
	r.GetChanRPC().Register("RenewalFeesSetInfo", r.RenewalFeesSetInfo)
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
		if retcode != 0 {
			cluster.SendMsgToHallUser(u.HallNodeId, u.Id, &msg.JoinRoomFaild{RoomID: r.DataMgr.GetRoomId()})
		}
	}()

	if r.Status == RoomStatusStarting && r.Temp.DynamicJoin == 1 {
		retcode = GameIsStart
		return
	}

	retcode = r.UserMgr.Sit(u, chairID, r.Status)

}

//起立
func (r *Entry_base) UserStandup(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserStandup{})
	u := args[1].(*user.User)
	r.ReqLeaveRoom([]interface{}{u})
	return
	//retcode := 0
	//defer func() {
	//	if retcode != 0 {
	//		u.WriteMsg(RenderErrorMessage(retcode))
	//	}
	//}()
	//
	//if r.Status == RoomStatusStarting {
	//	retcode = ErrGameIsStart
	//	return
	//}
	//
	//r.UserMgr.Standup(u)
}

func (r *Entry_base) RenewalFeesSetInfo(args []interface{}) (interface{}, error) {
	//addCnt := args[0].(int)
	rUserId := args[1].(int64)
	rNodeId := args[2].(int)
	if r.IsClose {
		return 2, errors.New("room is close ")
	}

	//不需要续费或者已经有人续过费了
	if r.TimerMgr.GetPlayCount() < r.TimerMgr.GetMaxPlayCnt() {
		return 3, errors.New("room playCnt >= maxPlayCnt")
	}

	if r.DelayCloseTimer != nil {
		r.DelayCloseTimer.Stop()
		r.DelayCloseTimer = nil
	}

	//旧房主uid
	oldCreator := r.DataMgr.GetCreator()
	//续费的人成为新房主
	r.DataMgr.ResetRoomCreator(rUserId, rNodeId)
	//未开始游戏定时器
	r.TimerMgr.StartCreatorTimer(func() {
		roomLogData := datalog.RoomLog{}
		logData := roomLogData.GetRoomLogRecode(r.DataMgr.GetRoomId(), r.Temp.KindID, r.Temp.ServerID)
		roomLogData.UpdateGameLogRecode(logData, 4)
		r.OnEventGameConclude(NO_START_GER_DISMISS)
	})
	//重置已玩次数
	r.TimerMgr.ResetPlayCount()
	//更新游戏服房间状态
	r.Status = RoomStatusReady
	//更新大厅房间信息
	RoomMgr.UpdateRoomToHall(&msg.UpdateRoomInfo{
		RoomId: r.DataMgr.GetRoomId(),
		OpName: "SetRoomInfo",
		Data: map[string]interface{}{
			"RoomStatus": RoomStatusReady,
			"NewCreator": rUserId,    //续费玩家id
			"oldCreator": oldCreator, //旧房主
		},
	})
	//重置其他(与玩法相关联的东西)
	r.DataMgr.ResetGameAfterRenewal()

	//更新玩家状态，并下发续费成功
	r.UserMgr.ForEachUser(func(u *user.User) {
		r.UserMgr.SetUsetStatus(u, US_SIT)
		u.WriteMsg(&msg.G2C_RenewalFeesSuccess{UserID: rUserId})
	})

	return 0, nil
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

//大厅服发来的解散房间
func (room *Entry_base) DissumeRoom(args []interface{}) {
	room.OnEventGameConclude(GER_DISMISS)

	room.UserMgr.ForEachUser(func(u *user.User) {
		room.UserMgr.LeaveRoom(u, room.Status)
	})

	roomLogData := datalog.RoomLog{}
	logData := roomLogData.GetRoomLogRecode(room.DataMgr.GetRoomId(), room.Temp.KindID, room.Temp.ServerID)
	if logData != nil {
		user, _ := room.UserMgr.GetUserByUid(logData.UserId)
		if user == nil {
			roomLogData.UpdateRoomLogForOthers(logData, CreateRoomForOthers)
		}
	}
}

//玩家准备
func (room *Entry_base) UserReady(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserReady)
	u := args[1].(*user.User)
	retCode := 0
	defer func() {
		if retCode != 0 {
			u.WriteMsg(RenderErrorMessage(retCode))
		}
	}()

	if u.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		retCode = ErrPlayerIsReady
		return
	}

	if room.DelayCloseTimer != nil {
		if room.TimerMgr.GetMaxPlayCnt() <= room.TimerMgr.GetPlayCount() {
			log.Debug("Max Play count limit, curCount=%d, maxCount=%d", room.TimerMgr.GetPlayCount(), room.TimerMgr.GetMaxPlayCnt())
			retCode = ErrRenewalFee
			return
		} else {
			log.Debug("ErrRoomIsClose")
			retCode = ErrRoomIsClose
			return
		}
	}

	log.Debug("at Entry_base UserReady ==== ")
	if u.Status != US_PLAYING {
		room.UserMgr.SetUsetStatus(u, US_READY)
	}

	if room.UserMgr.IsAllReady() {
		room.UserMgr.ResetBeginPlayer()
		RoomMgr.UpdateRoomToHall(&msg.UpdateRoomInfo{ //通知大厅服这个房间加局数
			RoomId: room.DataMgr.GetRoomId(),
			OpName: "AddPlayCnt",
			Data: map[string]interface{}{
				"Status": RoomStatusStarting,
				"Cnt":    1,
			},
		})
		//房间托管
		if room.RoomTrusteeTimer != nil {
			room.RoomTrusteeTimer.Stop()
			log.Debug("@@@@@@@@@@@@@@@@取消房间托管 游戏局数:%d", room.TimerMgr.GetPlayCount())
		}

		if room.TimerMgr.GetPlayCount() < room.TimerMgr.GetMaxPlayCnt() {
			room.TimerMgr.AddPlayCount()
		}
		room.DataMgr.BeforeStartGame(room.UserMgr.GetCurPlayerCnt())
		room.DataMgr.StartGameing()
		room.DataMgr.AfterStartGame()
		//派发初始扑克
		room.Status = RoomStatusStarting
		room.TimerMgr.StopCreatorTimer()
	} else {
		log.Debug(" not all ready")
	}
}

//玩家重登
func (room *Entry_base) UserReLogin(args []interface{}) (interface{}, error) {
	u := args[0].(*user.User)
	roomUser := room.getRoomUser(u.Id)
	if roomUser == nil {
		log.Debug("UserReLogin not old user")
		return false, nil
	}
	log.Debug("at ReLogin have old user ")
	u.ChairId = roomUser.ChairId
	u.Status = roomUser.Status
	u.ChatRoomId = roomUser.ChatRoomId
	u.RoomId = room.DataMgr.GetRoomId()
	u.Score = int64(room.DataMgr.GetUserScore(u.ChairId))
	room.UserMgr.ReLogin(u, room.Status)
	room.TimerMgr.StopOfflineTimer(u.Id)
	//重入取消托管
	if room.Temp.OffLineTrustee == 1 {
		room.OnUserTrusteeCb(u.ChairId, false)
	}
	return true, nil
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
	if room.Temp.OffLineTrustee == 0 {
		room.TimerMgr.StartKickoutTimer(u.Id, func() {
			room.OffLineTimeOut(u)
		})
	}
}

//离线超时踢出
func (room *Entry_base) OffLineTimeOut(u *user.User) {
	room.UserMgr.LeaveRoom(u, room.Status)
	if room.UserMgr.GetCurPlayerCnt() == 0 { //没人了直接销毁
		log.Debug("at OffLineTimeOut ======= ")
		room.AfterEnd(true, GER_DISMISS)
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
	BirefInf.CurPayCnt = room.TimerMgr.GetPlayCount()      //已玩局数
	BirefInf.PayCnt = room.TimerMgr.GetMaxPlayCnt()        //可玩局数
	BirefInf.RoomPlayCnt = room.TimerMgr.GetRoomPlayCnt()  //房间局数配置
	BirefInf.CreateTime = room.TimerMgr.GetCreatrTime()    //创建时间
	BirefInf.CreateUserId = room.DataMgr.GetCreator()
	BirefInf.IsPublic = room.UserMgr.IsPublic()
	BirefInf.Players = make(map[int64]*msg.PlayerBrief)
	BirefInf.MachPlayer = make(map[int64]int64) //todo
	return BirefInf

}

//游戏配置
func (room *Entry_base) SetGameOption(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_GameOption)
	u := args[1].(*user.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	log.Debug("at SetGameOption room id :%d", u.RoomId)

	//if u.ChairId == INVALID_CHAIR {
	//	retcode = ErrNoSitdowm
	//	return
	//}

	AllowLookon := 0
	if u.Status == US_LOOKON {
		AllowLookon = 1
	}
	u.WriteMsg(&msg.G2C_GameStatus{
		GameStatus:  room.Status,
		AllowLookon: AllowLookon,
	})

	room.DataMgr.SendPersonalTableTip(u)

	if room.Status == RoomStatusReady || room.Status == RoomStatusEnd { // 没开始
		room.DataMgr.SendStatusReady(u)
	} else { //开始了
		room.DataMgr.SendStatusPlay(u)
	}

	//把所有玩家信息推送给自己
	room.UserMgr.SendUserInfoToSelf(u)

}

//用户托管
func (room *Entry_base) OnRecUserTrustee(args []interface{}) {
	u := args[1].(*user.User)
	getData := args[0].(*mj_zp_msg.C2G_MJZP_Trustee)
	ok := room.OnUserTrusteeCb(u.ChairId, getData.Trustee)
	if !ok {
		log.Error("at OnRecUserTrustee user.chairid:", u.ChairId)
	}
}

//玩家离开房间
func (room *Entry_base) ReqLeaveRoom(args []interface{}) (interface{}, error) {
	player := args[0].(*user.User)
	if room.Status == RoomStatusReady {
		if room.UserMgr.LeaveRoom(player, room.Status) {
			player.WriteMsg(&msg.G2C_LeaveRoomRsp{Status: room.Status})
		} else {
			player.WriteMsg(&msg.G2C_LeaveRoomRsp{Status: room.Status, Code: ErrLoveRoomFaild})
		}
	} else {
		room.UserMgr.AddLeavePly(player.Id)
		room.UserMgr.SendMsgAllNoSelf(player.Id, &msg.G2C_LeaveRoomBradcast{UserID: player.Id})
		room.TimerMgr.StartReplytIimer(player.Id, func() {
			player.WriteMsg(&msg.G2C_LeaveRoomRsp{Status: room.Status})
			room.OnEventGameConclude(USER_LEAVE)
		})
	}
	return nil, nil
}

//其他玩家响应玩家离开房间的请求
func (room *Entry_base) ReplyLeaveRoom(args []interface{}) {
	log.Debug("at ReplyLeaveRoom ")
	player := args[0].(*user.User)
	Agree := args[1].(bool)
	ReplyUid := args[2].(int64)
	ret := room.UserMgr.ReplyLeave(player, Agree, ReplyUid, room.Status)
	if ret == 1 {
		reqPlayer, _ := room.UserMgr.GetUserByUid(ReplyUid)
		if reqPlayer != nil {
			reqPlayer.WriteMsg(&msg.G2C_LeaveRoomRsp{Status: room.Status})
		}

		room.OnEventGameConclude(USER_LEAVE)
	} else if ret == -1 { //有人拒绝
		room.TimerMgr.StopReplytIimer(ReplyUid)
		reqPlayer, _ := room.UserMgr.GetUserByUid(ReplyUid)
		if reqPlayer != nil {
			reqPlayer.WriteMsg(&msg.G2C_LeaveRoomRsp{Status: room.Status, Code: ErrRefuseLeave})
		}
	}
}

//游戏结束
func (room *Entry_base) OnEventGameConclude(cbReason int) {
	if room.Status == RoomStatusClose {
		log.Debug("double close room")
		return
	}
	switch cbReason {
	case GER_NORMAL: //常规结束
		room.DataMgr.NormalEnd(cbReason)
		room.AfterEnd(false, cbReason)
		room.OnRoomTrustee()
	case GER_DISMISS: //游戏解散
		room.DataMgr.DismissEnd(cbReason)
		room.AfterEnd(true, cbReason)
	case USER_LEAVE: //用户请求解散
		room.DataMgr.NormalEnd(cbReason)
		room.AfterEnd(true, cbReason)
	case NO_START_GER_DISMISS: //没开始就解散
		room.DataMgr.DismissEnd(cbReason)
		room.AfterEnd(true, cbReason)
	case ROOM_TRUSTEE: //房间托管结束
		room.DataMgr.NormalEnd(cbReason)
		room.AfterEnd(false, cbReason)
	}
	log.Debug("at OnEventGameConclude cbReason:%d ", cbReason)
	return
}

// 如果这里不能满足 afertEnd 请重构这个到个个组件里面
func (room *Entry_base) AfterEnd(Forced bool, cbReason int) {
	roomStatus := room.Status
	room.Status = RoomStatusEnd //一局结束状态
	if Forced || room.TimerMgr.GetPlayCount() >= room.TimerMgr.GetMaxPlayCnt() {
		if room.DelayCloseTimer != nil {
			room.DelayCloseTimer.Stop()
		}
		room.Status = RoomStatusClose
		log.Debug("Forced :%v, room.Status:%d, PlayTurnCount:%d, temp PlayTurnCount:%d", Forced, roomStatus, room.TimerMgr.GetPlayCount(), room.TimerMgr.GetMaxPlayCnt())
		closeFunc := func() {
			room.IsClose = true
			room.UserMgr.SendMsgToHallServerAll(&msg.RoomEndInfo{
				RoomId: room.DataMgr.GetRoomId(),
				Status: roomStatus,
			})

			//全付的房间，若没开始过并且创建的房主没在，则返还给他钱
			room.UserMgr.CheckRoomReturnMoney(roomStatus, room.DataMgr.GetCreatorNodeId(), room.DataMgr.GetRoomId(), room.DataMgr.GetCreator())

			room.Destroy(room.DataMgr.GetRoomId())
			room.UserMgr.RoomDissume(cbReason)
		}

		if GER_NORMAL != cbReason {
			room.DelayCloseTimer = room.AfterFunc(2*time.Second, closeFunc)
		} else { //常规结束延迟
			room.DelayCloseTimer = room.AfterFunc(time.Duration(GetGlobalVarInt(DelayDestroyRoom))*time.Second, closeFunc)
		}

		return
	}

	room.UserMgr.ForEachUser(func(u *user.User) {
		room.UserMgr.SetUsetStatus(u, US_SIT)
	})
}

//房间托管
func (room *Entry_base) OnRoomTrustee() bool {
	if room.Status == RoomStatusStarting {
		return false
	}

	var AddPlayCount func()
	AddPlayCount = func() {
		if room.TimerMgr.GetPlayCount() < room.TimerMgr.GetMaxPlayCnt() {
			room.TimerMgr.AddPlayCount()

			RoomMgr.UpdateRoomToHall(&msg.UpdateRoomInfo{ //通知大厅服这个房间加局数
				RoomId: room.DataMgr.GetRoomId(),
				OpName: "AddPlayCnt",
				Data: map[string]interface{}{
					"Status": RoomStatusStarting,
					"Cnt":    1,
				},
			})
			room.RoomTrusteeTimer = room.AfterFunc(time.Duration(room.Temp.TimeRoomTrustee)*time.Second, AddPlayCount)
			log.Debug("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 局数+1 总局数;%d", room.TimerMgr.GetPlayCount())
		} else { //最后一局
			room.OnEventGameConclude(ROOM_TRUSTEE)
			log.Debug("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 游戏结束 总局数;%d", room.TimerMgr.GetPlayCount())
		}
	}

	if room.TimerMgr.GetPlayCount() < room.TimerMgr.GetMaxPlayCnt() {
		log.Debug("进入房间托管")
		room.RoomTrusteeTimer = room.AfterFunc(time.Duration(room.Temp.TimeRoomTrustee)*time.Second, AddPlayCount)
	}
	return true
}

func (room *Entry_base) getRoomUser(uid int64) *user.User {
	u, _ := room.UserMgr.GetUserByUid(uid)
	return u
}
