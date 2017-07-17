package room

import (
	. "mj/common/cost"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model"
	"mj/gameServer/user"

	"mj/common/msg/mj_zp_msg"

	"github.com/lovelly/leaf/log"
)

type ZP_base struct {
	*mj_base.Mj_base
}

func NewMJBase(info *model.CreateRoomInfo) *ZP_base {
	mj := new(ZP_base)
	mj.Mj_base = mj_base.NewMJBase(info)

	return mj
}

//出牌
func (room *ZP_base) OutCard(args []interface{}) {
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
		log.Error("zpmj at OnUserOutCard game status != RoomStatusStarting ")
		retcode = ErrGameNotStart
		return
	}

	//效验参数
	if u.ChairId != room.DataMgr.GetCurrentUser() {
		log.Error("u.ChairId:%d  room.DataMgr.GetCurrentUser():%d", u.ChairId, room.DataMgr.GetCurrentUser())
		log.Error("zpmj at OnUserOutCard not self out ")
		retcode = ErrNotSelfOut
		return
	}

	if !room.LogicMgr.IsValidCard(CardData) {
		log.Error("zpmj at OnUserOutCard IsValidCard card ")
		retcode = NotValidCard
	}

	//吃啥打啥
	if !room.DataMgr.OutOfChiCardRule(CardData, u.ChairId) {
		log.Error("zpmj at OutOfChiCardRule IsValidCard card ")
		retcode = NotValidCard
	}

	//清除出牌禁忌
	room.DataMgr.ClearBanCard(u.ChairId)

	//删除扑克
	if !room.LogicMgr.RemoveCard(room.DataMgr.GetUserCardIndex(u.ChairId), CardData) {
		log.Error("zpmj at OnUserOutCard not have card ")
		log.Error("user:%d card:%d", u.ChairId, CardData)
		retcode = ErrNotFoudCard
		return
	}

	//todo.记录出牌用户

	//记录出牌数
	room.DataMgr.RecordOutCarCnt()

	//记录跟牌
	room.DataMgr.RecordFollowCard(CardData)

	u.UserLimit &= ^LimitChiHu
	u.UserLimit &= ^LimitPeng
	u.UserLimit &= ^LimitGang

	var bSysOut bool
	if len(args) == 3 {
		bSysOut = args[2].(bool)
	}
	room.DataMgr.NotifySendCard(u, CardData, bSysOut)

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

// 吃碰杠胡各种操作
func (room *ZP_base) UserOperateCard(args []interface{}) {
	u := args[0].(*user.User)
	OperateCode := args[1].(int)
	OperateCard := args[2].([]int)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		}
	}()
	log.Debug("吃碰杠胡各种操作 user:%d", room.DataMgr.GetCurrentUser())
	if room.DataMgr.GetCurrentUser() == INVALID_CHAIR {

		//效验状态
		if !room.DataMgr.HasOperator(u.ChairId, OperateCode) {
			log.Debug("user:%d,OperateCode:%d,OperateCard：%v", u.ChairId, OperateCode, OperateCard)
			retcode = ErrNoOperator
			return
		}

		//变量定义
		cbTargetAction, wTargetUser := room.DataMgr.CheckUserOperator(u, room.UserMgr.GetMaxPlayerCnt(), OperateCode, OperateCard)
		if cbTargetAction < 0 {
			log.Debug("wait other user")
			return
		}

		//放弃操作
		if cbTargetAction == WIK_NULL {
			//用户状态
			if room.DataMgr.DispatchCardData(room.DataMgr.GetResumeUser(), room.DataMgr.GetGangStatus() != WIK_GANERAL) > 0 {
				room.OnEventGameConclude(room.DataMgr.GetProvideUser(), nil, GER_NORMAL)
			}
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
				room.OnEventGameConclude(room.DataMgr.GetProvideUser(), nil, GER_NORMAL)
			}
		}
	} else {

		//扑克效验
		if (OperateCode != WIK_NULL) && (OperateCode != WIK_CHI_HU) && (!room.LogicMgr.IsValidCard(OperateCard[0])) {
			log.Error("OperateCode != WIK_NULL) && (OperateCode != WIK_CHI_HU) && (!room.LogicMgr.IsValidCard(OperateCard[0])")
			return
		}

		//设置变量
		// room.UserAction[room.CurrentUser] = WIK_NULL

		//执行动作
		switch OperateCode {
		case WIK_GANG: //杠牌操作
			cbGangKind := room.DataMgr.AnGang(u, OperateCode, OperateCard)
			//效验动作
			bAroseAction := false
			if cbGangKind == WIK_MING_GANG {
				bAroseAction = room.DataMgr.EstimateUserRespond(u.ChairId, OperateCard[0], EstimatKind_GangCard)
			}

			//清除操作定时
			room.DataMgr.StopOperateCardTimer(u)

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

//抓花
func (room *ZP_base) ZhaMa(args []interface{}) {
	return
}

//托管
func (room *ZP_base) OnUserTrustee(wChairID int, bTrustee bool) bool {
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
