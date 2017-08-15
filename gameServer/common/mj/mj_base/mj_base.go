package mj_base

import (
	. "mj/common/cost"
	"mj/common/msg/mj_zp_msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/room_base"
	"mj/gameServer/db/model/base"
	datalog "mj/gameServer/log"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"
)

type Mj_base struct {
	*room_base.Entry_base
	LogicMgr LogicManager
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
	mj.Entry_base = room_base.NewEntryBase(KindId, ServiceId)
	mj.Temp = Temp
	return mj
}

func (r *Mj_base) GetDataMgr() DataManager {
	return r.DataMgr.(DataManager)
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
	if u.ChairId != room.GetDataMgr().GetCurrentUser() {
		log.Error("at OnUserOutCard not self out ")
		log.Error("u.ChairId:%d, room.GetDataMgr().GetCurrentUser():%d", u.ChairId, room.GetDataMgr().GetCurrentUser())
		retcode = ErrNotSelfOut
		return
	}

	if !room.LogicMgr.IsValidCard(CardData) {
		log.Error("at OnUserOutCard IsValidCard card ")
		retcode = NotValidCard
	}

	//吃啥打啥
	if !room.GetDataMgr().OutOfChiCardRule(CardData, u.ChairId) {
		log.Error(" at OutOfChiCardRule IsValidCard card ")
		retcode = NotValidCard
	}

	//删除扑克
	if !room.LogicMgr.RemoveCard(room.GetDataMgr().GetUserCardIndex(u.ChairId), CardData) {
		log.Error("at OnUserOutCard not have card ：%d chairid:%d", CardData, u.ChairId)
		retcode = ErrNotFoudCard
		return
	}

	//记录出牌数
	room.GetDataMgr().RecordOutCarCnt()

	u.UserLimit &= ^LimitChiHu
	u.UserLimit &= ^LimitPeng
	u.UserLimit &= ^LimitGang

	log.Debug("AAAAAAAAAAAAAAAAA ")
	room.GetDataMgr().NotifySendCard(u, CardData, false)

	//响应判断
	bAroseAction := room.GetDataMgr().EstimateUserRespond(u.ChairId, CardData, EstimatKind_OutCard)

	//派发扑克
	if !bAroseAction {
		if room.GetDataMgr().DispatchCardData(room.GetDataMgr().GetCurrentUser(), false) > 0 {
			room.OnEventGameConclude(room.GetDataMgr().GetProvideUser(), nil, GER_NORMAL)
		}
	}

	return
}

//插花
func (room *Mj_base) ChaHuaMsg(args []interface{}) {
	u := args[1].(*user.User)
	getData := args[0].(*mj_zp_msg.C2G_MJZP_SetChaHua)
	room.GetDataMgr().GetChaHua(u, getData.SetCount)
}

//补花
func (room *Mj_base) OnUserReplaceCardMsg(args []interface{}) {
	u := args[0].(*user.User)
	CardData := args[1].(int)
	room.GetDataMgr().OnUserReplaceCard(u, CardData)
}

//用户听牌
func (room *Mj_base) OnUserListenCardMsg(args []interface{}) {
	u := args[1].(*user.User)
	getData := args[0].(*mj_zp_msg.C2G_MJZP_ListenCard)
	room.GetDataMgr().OnUserListenCard(u, getData.ListenCard)
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

	if room.GetDataMgr().GetCurrentUser() == INVALID_CHAIR {
		//效验状态
		if !room.GetDataMgr().HasOperator(u.ChairId, OperateCode) {
			retcode = ErrNoOperator
			return
		}

		//变量定义
		cbTargetAction, wTargetUser := room.GetDataMgr().CheckUserOperator(u, room.UserMgr.GetMaxPlayerCnt(), OperateCode, OperateCard)
		if cbTargetAction < 0 {
			log.Debug("wait other user, OperateCode=%d, OperateCard=%d, cbTargetAction=%v, wTargetUser=%v", OperateCode, OperateCard, cbTargetAction, wTargetUser)
			return
		}

		//放弃操作
		if cbTargetAction == WIK_NULL {
			//用户状态
			if room.GetDataMgr().DispatchCardData(room.GetDataMgr().GetResumeUser(), room.GetDataMgr().GetGangStatus() != WIK_GANERAL) > 0 {
				room.OnEventGameConclude(room.GetDataMgr().GetProvideUser(), nil, GER_NORMAL)
			}
			room.GetDataMgr().ResetUserOperate()
			return
		}

		//胡牌操作
		if cbTargetAction == WIK_CHI_HU {
			room.GetDataMgr().UserChiHu(wTargetUser, room.UserMgr.GetMaxPlayerCnt())
			room.OnEventGameConclude(room.GetDataMgr().GetProvideUser(), nil, GER_NORMAL)
			return
		}

		//收集牌
		room.GetDataMgr().WeaveCard(cbTargetAction, wTargetUser)

		//删除扑克
		if room.GetDataMgr().RemoveCardByOP(wTargetUser, cbTargetAction) == false {
			log.Error("at UserOperateCard RemoveCardByOP error")
			return
		}

		room.GetDataMgr().CallOperateResult(wTargetUser, cbTargetAction)
		if cbTargetAction == WIK_GANG {
			if room.GetDataMgr().DispatchCardData(wTargetUser, true) > 0 {
				room.GetDataMgr().GetRoomId()
				room.OnEventGameConclude(room.GetDataMgr().GetProvideUser(), nil, GER_NORMAL)
			}
		}
	} else {
		//扑克效验
		if (OperateCode != WIK_NULL) && (OperateCode != WIK_CHI_HU) && (!room.LogicMgr.IsValidCard(OperateCard[0])) {
			return
		}

		//设置变量
		room.GetDataMgr().ResetUserOperateEx(u)

		//执行动作
		switch OperateCode {
		case WIK_GANG: //杠牌操作
			cbGangKind := room.GetDataMgr().AnGang(u, OperateCode, OperateCard)
			//效验动作
			bAroseAction := false
			if cbGangKind == WIK_MING_GANG {
				bAroseAction = room.GetDataMgr().EstimateUserRespond(u.ChairId, OperateCard[0], EstimatKind_GangCard)
			}

			//发送扑克
			if !bAroseAction {
				if room.GetDataMgr().DispatchCardData(u.ChairId, true) > 0 {
					room.OnEventGameConclude(room.GetDataMgr().GetProvideUser(), nil, GER_NORMAL)
				}
			}
		case WIK_CHI_HU: //自摸
			//结束游戏
			room.GetDataMgr().ZiMo(u)
			room.OnEventGameConclude(room.GetDataMgr().GetProvideUser(), nil, GER_NORMAL)
		}
	}
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
		if wChairID == room.GetDataMgr().GetCurrentUser() && !room.GetDataMgr().IsActionDone() {
			cardindex := room.GetDataMgr().GetTrusteeOutCard(wChairID)
			if cardindex == INVALID_BYTE {
				return false
			}
			u := room.UserMgr.GetUserByChairId(wChairID)
			card := room.LogicMgr.SwitchToCardData(cardindex)
			room.OutCard([]interface{}{u, card, true})
		} else if room.GetDataMgr().GetCurrentUser() == INVALID_CHAIR && !room.GetDataMgr().IsActionDone() {
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
