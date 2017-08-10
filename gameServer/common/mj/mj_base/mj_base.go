package mj_base

import (
	"errors"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/mj_zp_msg"
	"mj/gameServer/RoomMgr"
	. "mj/gameServer/common"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/room_base"
	"mj/gameServer/conf"
	"mj/gameServer/db/model/base"
	datalog "mj/gameServer/log"
	"mj/gameServer/user"
	"time"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
	"github.com/lovelly/leaf/timer"
)

type Mj_base struct {
	room_base.BaseManager
	UserMgr  room_base.UserManager
	TimerMgr room_base.TimerManager
	DataMgr  DataManager
	LogicMgr LogicManager

	Temp             *base.GameServiceOption //模板
	Status           int
	IsClose          bool
	DelayCloseTimer  *timer.Timer
	RoomTrusteeTimer *timer.Timer
}

//创建的配置文件
type NewMjCtlConfig struct {
	BaseMgr  room_base.BaseManager
	UserMgr  room_base.UserManager
	TimerMgr room_base.TimerManager
	DataMgr  DataManager
	LogicMgr LogicManager
}

func NewMJBase(KindId, ServiceId int) *Mj_base {
	Temp, ok1 := base.GameServiceOptionCache.Get(KindId, ServiceId)
	if !ok1 {
		log.Error("at NewMJBase not foud template  kindid:%d  serverid:%d", KindId, ServiceId)
		return nil
	}

	mj := new(Mj_base)
	mj.Temp = Temp
	return mj
}

func (r *Mj_base) RegisterBaseFunc() {
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
	r.GetChanRPC().Register("AddPlayCnt", r.AddPlayCnt)
}

func (r *Mj_base) Init(cfg *NewMjCtlConfig) {
	r.UserMgr = cfg.UserMgr
	r.DataMgr = cfg.DataMgr
	r.BaseManager = cfg.BaseMgr
	r.LogicMgr = cfg.LogicMgr
	r.TimerMgr = cfg.TimerMgr
	r.RoomRun(r.DataMgr.GetRoomId())
	r.TimerMgr.StartCreatorTimer(func() {
		roomLogData := datalog.RoomLog{}
		logData := roomLogData.GetRoomLogRecode(r.DataMgr.GetRoomId(), r.Temp.KindID, r.Temp.ServerID)
		roomLogData.UpdateGameLogRecode(logData, 4)
		r.OnEventGameConclude(0, nil, NO_START_GER_DISMISS)
	})

	r.DataMgr.InitRoomOne()
}

func (r *Mj_base) GetRoomId() int {
	return r.DataMgr.GetRoomId()
}

//坐下
func (r *Mj_base) Sitdown(args []interface{}) {
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
func (r *Mj_base) UserStandup(args []interface{}) {
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

func (r *Mj_base) AddPlayCnt(args []interface{}) (interface{}, error) {
	log.Debug("at AddPlayCnt .... ")
	if r.IsClose {
		return 1, errors.New("room is close ")
	}

	addCnt := args[0].(int)
	//不需要续费或者已经有人续过费了
	if r.TimerMgr.GetPlayCount() < r.TimerMgr.GetMaxPlayCnt() {
		return 2, errors.New("room playCnt >= maxPlayCnt")
	}

	r.TimerMgr.AddMaxPlayCnt(addCnt)
	//r.TimerMgr.ResetPlayCount()

	if r.DelayCloseTimer != nil {
		r.DelayCloseTimer.Stop()
		r.DelayCloseTimer = nil
	}
	log.Debug("at AddPlayCnt ...1111 . ")
	return 0, nil
}

//获取对方信息
func (room *Mj_base) GetUserChairInfo(args []interface{}) {
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
func (room *Mj_base) DissumeRoom(args []interface{}) {
	room.UserMgr.ForEachUser(func(u *user.User) {
		room.UserMgr.LeaveRoom(u, room.Status)
	})

	room.OnEventGameConclude(0, nil, GER_DISMISS)

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
func (room *Mj_base) UserReady(args []interface{}) {
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

	log.Debug("at UserReady ==== ")
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
		room.DataMgr.BeforeStartGame(room.UserMgr.GetMaxPlayerCnt())
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
func (room *Mj_base) UserReLogin(args []interface{}) error {
	u := args[0].(*user.User)
	roomUser := room.getRoomUser(u.Id)
	if roomUser == nil {
		return errors.New(" UserReLogin not old user ")
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
		room.OnUserTrustee(u.ChairId, false)
	}
	return nil
}

//玩家离线
func (room *Mj_base) UserOffline(args []interface{}) {
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
func (room *Mj_base) OffLineTimeOut(u *user.User) {
	room.UserMgr.LeaveRoom(u, room.Status)
	if room.UserMgr.GetCurPlayerCnt() == 0 { //没人了直接销毁
		log.Debug("at OffLineTimeOut ======= ")
		room.AfterEnd(true, GER_DISMISS)
	}
}

//获取房间基础信息
func (room *Mj_base) GetBirefInfo() *msg.RoomInfo {
	BirefInf := &msg.RoomInfo{}
	BirefInf.ServerID = room.Temp.ServerID
	BirefInf.KindID = room.Temp.KindID
	BirefInf.NodeID = conf.Server.NodeId
	BirefInf.SvrHost = conf.Server.WSAddr
	BirefInf.PayType = room.UserMgr.GetPayType()
	BirefInf.RoomID = room.DataMgr.GetRoomId()
	BirefInf.CurCnt = room.UserMgr.GetCurPlayerCnt()
	BirefInf.MaxPlayerCnt = room.UserMgr.GetMaxPlayerCnt() //最多多人数
	BirefInf.PayCnt = room.TimerMgr.GetMaxPlayCnt()        //可玩局数
	BirefInf.CurPayCnt = room.TimerMgr.GetPlayCount()      //已玩局数
	BirefInf.CreateTime = room.TimerMgr.GetCreatrTime()    //创建时间
	BirefInf.CreateUserId = room.DataMgr.GetCreator()
	BirefInf.IsPublic = room.UserMgr.IsPublic()
	BirefInf.Players = make(map[int64]*msg.PlayerBrief)
	BirefInf.MachPlayer = make(map[int64]int64) //todo
	return BirefInf

}

//游戏配置
func (room *Mj_base) SetGameOption(args []interface{}) {
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

//出牌
func (room *Mj_base) OutCard(args []interface{}) {
	u := args[0].(*user.User)
	CardData := args[1].(int)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	//效验状态
	if room.Status != RoomStatusStarting {
		log.Error("at OnUserOutCard game status != RoomStatusStarting ")
		retcode = ErrGameNotStart
		return
	}

	//效验参数
	if u.ChairId != room.DataMgr.GetCurrentUser() {
		log.Error("at OnUserOutCard not self out ")
		log.Error("u.ChairId:%d, room.DataMgr.GetCurrentUser():%d", u.ChairId, room.DataMgr.GetCurrentUser())
		retcode = ErrNotSelfOut
		return
	}

	if !room.LogicMgr.IsValidCard(CardData) {
		log.Error("at OnUserOutCard IsValidCard card ")
		retcode = NotValidCard
	}

	//吃啥打啥
	if !room.DataMgr.OutOfChiCardRule(CardData, u.ChairId) {
		log.Error(" at OutOfChiCardRule IsValidCard card ")
		retcode = NotValidCard
	}

	//删除扑克
	if !room.LogicMgr.RemoveCard(room.DataMgr.GetUserCardIndex(u.ChairId), CardData) {
		log.Error("at OnUserOutCard not have card ：%d chairid:%d", CardData, u.ChairId)
		retcode = ErrNotFoudCard
		return
	}

	//记录出牌数
	room.DataMgr.RecordOutCarCnt()

	u.UserLimit &= ^LimitChiHu
	u.UserLimit &= ^LimitPeng
	u.UserLimit &= ^LimitGang

	log.Debug("AAAAAAAAAAAAAAAAA ")
	room.DataMgr.NotifySendCard(u, CardData, false)

	//响应判断
	bAroseAction := room.DataMgr.EstimateUserRespond(u.ChairId, CardData, EstimatKind_OutCard)

	//派发扑克
	if !bAroseAction {
		if room.DataMgr.DispatchCardData(room.DataMgr.GetCurrentUser(), false) > 0 {
			room.OnEventGameConclude(room.DataMgr.GetProvideUser(), nil, GER_NORMAL)
		}
	}

	return
}

//插花
func (room *Mj_base) ChaHuaMsg(args []interface{}) {
	u := args[1].(*user.User)
	getData := args[0].(*mj_zp_msg.C2G_MJZP_SetChaHua)
	room.DataMgr.GetChaHua(u, getData.SetCount)
}

//补花
func (room *Mj_base) OnUserReplaceCardMsg(args []interface{}) {
	u := args[0].(*user.User)
	CardData := args[1].(int)
	room.DataMgr.OnUserReplaceCard(u, CardData)
}

//用户听牌
func (room *Mj_base) OnUserListenCardMsg(args []interface{}) {
	u := args[1].(*user.User)
	getData := args[0].(*mj_zp_msg.C2G_MJZP_ListenCard)
	room.DataMgr.OnUserListenCard(u, getData.ListenCard)
}

//用户托管
func (room *Mj_base) OnRecUserTrustee(args []interface{}) {
	u := args[1].(*user.User)
	getData := args[0].(*mj_zp_msg.C2G_MJZP_Trustee)
	ok := room.OnUserTrustee(u.ChairId, getData.Trustee)
	if !ok {
		log.Error("at OnRecUserTrustee user.chairid:", u.ChairId)
	}
}

// 吃碰杠胡各种操作
func (room *Mj_base) UserOperateCard(args []interface{}) {
	u := args[0].(*user.User)
	OperateCode := args[1].(int)
	OperateCard := args[2].([]int)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	if room.DataMgr.GetCurrentUser() == INVALID_CHAIR {
		//效验状态
		if !room.DataMgr.HasOperator(u.ChairId, OperateCode) {
			retcode = ErrNoOperator
			return
		}

		//变量定义
		cbTargetAction, wTargetUser := room.DataMgr.CheckUserOperator(u, room.UserMgr.GetMaxPlayerCnt(), OperateCode, OperateCard)
		if cbTargetAction < 0 {
			log.Debug("wait other user, OperateCode=%d, OperateCard=%d, cbTargetAction=%v, wTargetUser=%v", OperateCode, OperateCard, cbTargetAction, wTargetUser)
			return
		}

		//放弃操作
		if cbTargetAction == WIK_NULL {
			//用户状态
			if room.DataMgr.DispatchCardData(room.DataMgr.GetResumeUser(), room.DataMgr.GetGangStatus() != WIK_GANERAL) > 0 {
				room.OnEventGameConclude(room.DataMgr.GetProvideUser(), nil, GER_NORMAL)
			}
			room.DataMgr.ResetUserOperate()
			return
		}

		//胡牌操作
		if cbTargetAction == WIK_CHI_HU {
			room.DataMgr.UserChiHu(wTargetUser, room.UserMgr.GetMaxPlayerCnt())
			room.OnEventGameConclude(room.DataMgr.GetProvideUser(), nil, GER_NORMAL)
			return
		}

		//收集牌
		room.DataMgr.WeaveCard(cbTargetAction, wTargetUser)

		//删除扑克
		if room.DataMgr.RemoveCardByOP(wTargetUser, cbTargetAction) == false {
			log.Error("at UserOperateCard RemoveCardByOP error")
			return
		}

		room.DataMgr.CallOperateResult(wTargetUser, cbTargetAction)
		if cbTargetAction == WIK_GANG {
			if room.DataMgr.DispatchCardData(wTargetUser, true) > 0 {
				room.DataMgr.GetRoomId()
				room.OnEventGameConclude(room.DataMgr.GetProvideUser(), nil, GER_NORMAL)
			}
		}
	} else {
		//扑克效验
		if (OperateCode != WIK_NULL) && (OperateCode != WIK_CHI_HU) && (!room.LogicMgr.IsValidCard(OperateCard[0])) {
			return
		}

		//设置变量
		room.DataMgr.ResetUserOperateEx(u)

		//执行动作
		switch OperateCode {
		case WIK_GANG: //杠牌操作
			cbGangKind := room.DataMgr.AnGang(u, OperateCode, OperateCard)
			//效验动作
			bAroseAction := false
			if cbGangKind == WIK_MING_GANG {
				bAroseAction = room.DataMgr.EstimateUserRespond(u.ChairId, OperateCard[0], EstimatKind_GangCard)
			}

			//发送扑克
			if !bAroseAction {
				if room.DataMgr.DispatchCardData(u.ChairId, true) > 0 {
					room.OnEventGameConclude(room.DataMgr.GetProvideUser(), nil, GER_NORMAL)
				}
			}
		case WIK_CHI_HU: //自摸
			//结束游戏
			room.DataMgr.ZiMo(u)
			room.OnEventGameConclude(room.DataMgr.GetProvideUser(), nil, GER_NORMAL)
		}
	}
}

//玩家离开房间
func (room *Mj_base) ReqLeaveRoom(args []interface{}) {
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
			room.OnEventGameConclude(player.ChairId, player, USER_LEAVE)
		})
	}
}

//其他玩家响应玩家离开房间的请求
func (room *Mj_base) ReplyLeaveRoom(args []interface{}) {
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

		room.OnEventGameConclude(player.ChairId, player, USER_LEAVE)
	} else if ret == -1 { //有人拒绝
		room.TimerMgr.StopReplytIimer(ReplyUid)
		reqPlayer, _ := room.UserMgr.GetUserByUid(ReplyUid)
		if reqPlayer != nil {
			reqPlayer.WriteMsg(&msg.G2C_LeaveRoomRsp{Status: room.Status, Code: ErrRefuseLeave})
		}
	}
}

//游戏结束
func (room *Mj_base) OnEventGameConclude(ChairId int, user *user.User, cbReason int) {
	switch cbReason {
	case GER_NORMAL: //常规结束
		room.DataMgr.NormalEnd(cbReason)
		room.AfterEnd(false, cbReason)
	case GER_DISMISS: //游戏解散
		room.DataMgr.DismissEnd(cbReason)
		room.AfterEnd(true, cbReason)
	case USER_LEAVE: //用户请求解散
		room.DataMgr.NormalEnd(cbReason)
		room.AfterEnd(true, cbReason)
	case NO_START_GER_DISMISS: //没开始就解散
		room.DataMgr.DismissEnd(cbReason)
		room.AfterEnd(true, cbReason)
	}
	room.Status = RoomStatusEnd
	log.Debug("at OnEventGameConclude cbReason:%d ", cbReason)
	return
}

// 如果这里不能满足 afertEnd 请重构这个到个个组件里面
func (room *Mj_base) AfterEnd(Forced bool, cbReason int) {
	roomStatus := room.Status
	room.TimerMgr.AddPlayCount()
	//room.OnRoomTrustee()
	if Forced || room.TimerMgr.GetPlayCount() >= room.TimerMgr.GetMaxPlayCnt() {
		if room.DelayCloseTimer != nil {
			room.DelayCloseTimer.Stop()
		}
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
			room.UserMgr.RoomDissume()
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

////todo,房间托管
//func (room *Mj_base) OnRoomTrustee() {
//	TrusteeCnt := 0
//	room.UserMgr.ForEachUser(func(u *user.User) {
//		if room.UserMgr.IsTrustee(u.ChairId) {
//			TrusteeCnt++
//		}
//	})
//
//	var AddPlayCount func()
//	AddPlayCount = func() {
//		if room.TimerMgr.GetPlayCount() <= room.TimerMgr.GetMaxPlayCnt() {
//			room.TimerMgr.AddPlayCount()
//			room.RoomTrusteeTimer = room.AfterFunc(time.Duration(room.Temp.TimeRoomTrustee)*time.Second, AddPlayCount)
//			log.Debug("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 局数+1 总局数;%d", room.TimerMgr.GetPlayCount())
//		} else {
//			room.OnEventGameConclude(0, nil, GER_NORMAL)
//			log.Debug("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@ 游戏结束 总局数;%d", room.TimerMgr.GetPlayCount())
//		}
//	}
//
//	if TrusteeCnt == room.UserMgr.GetMaxPlayerCnt() && room.TimerMgr.GetPlayCount() <= room.TimerMgr.GetMaxPlayCnt() {
//		log.Debug("进入房间托管")
//		room.RoomTrusteeTimer = room.AfterFunc(time.Duration(room.Temp.TimeRoomTrustee)*time.Second, AddPlayCount)
//	}
//}

//托管
func (room *Mj_base) OnUserTrustee(wChairID int, bTrustee bool) bool {

	//效验状态
	if wChairID >= room.UserMgr.GetMaxPlayerCnt() {
		return false
	}

	room.UserMgr.SetUsetTrustee(wChairID, bTrustee)

	room.UserMgr.SendMsgAll(&mj_zp_msg.G2C_ZPMJ_Trustee{
		Trustee: bTrustee,
		ChairID: wChairID,
	})

	if bTrustee {
		if wChairID == room.DataMgr.GetCurrentUser() && !room.DataMgr.IsActionDone() {
			cardindex := room.DataMgr.GetTrusteeOutCard(wChairID)
			if cardindex == INVALID_BYTE {
				return false
			}
			u := room.UserMgr.GetUserByChairId(wChairID)
			card := room.LogicMgr.SwitchToCardData(cardindex)
			room.OutCard([]interface{}{u, card, true})
		} else if room.DataMgr.GetCurrentUser() == INVALID_CHAIR && !room.DataMgr.IsActionDone() {
			u := room.UserMgr.GetUserByChairId(wChairID)
			if u == nil {
				return false
			}
			operateCard := []int{0, 0, 0}
			room.UserOperateCard([]interface{}{u, WIK_NULL, operateCard})
		}
	}
	return true
}

func (room *Mj_base) getRoomUser(uid int64) *user.User {
	u, _ := room.UserMgr.GetUserByUid(uid)
	return u
}
