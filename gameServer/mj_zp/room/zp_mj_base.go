package room

import (
	. "mj/common/cost"
	"mj/common/msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/user"

	"mj/common/msg/mj_zp_msg"

	"github.com/lovelly/leaf/log"
)

type ZP_base struct {
	*mj_base.Mj_base
}

func NewMJBase(info *msg.L2G_CreatorRoom) *ZP_base {
	mj := new(ZP_base)
	mj.Mj_base = mj_base.NewMJBase(info.KindId, info.ServiceId)
	if mj.Mj_base == nil {
		return nil
	}
	return mj
}
func (room *ZP_base) GetDataMgr() *ZP_RoomData {
	return room.DataMgr.(*ZP_RoomData)
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
	if u.ChairId != room.GetDataMgr().GetCurrentUser() {
		log.Error("u.ChairId:%d  room.GetDataMgr().GetCurrentUser():%d", u.ChairId, room.GetDataMgr().GetCurrentUser())
		log.Error("zpmj at OnUserOutCard not self out ")
		retcode = ErrNotSelfOut
		return
	}

	if !room.LogicMgr.IsValidCard(CardData) {
		log.Error("zpmj at OnUserOutCard IsValidCard card ")
		retcode = NotValidCard
	}

	//吃啥打啥
	if !room.GetDataMgr().OutOfChiCardRule(CardData, u.ChairId) {
		log.Error("zpmj at OutOfChiCardRule IsValidCard card ")
		retcode = NotValidCard
	}

	//删除扑克
	if !room.LogicMgr.RemoveCard(room.GetDataMgr().GetUserCardIndex(u.ChairId), CardData) {
		log.Error("zpmj at OnUserOutCard not have card ")
		log.Error("user:%d card:%d", u.ChairId, CardData)
		retcode = ErrNotFoudCard
		return
	}

	//记录出牌数
	room.GetDataMgr().RecordOutCarCnt()

	//记录跟牌
	room.GetDataMgr().RecordFollowCard(u.ChairId, CardData)

	u.UserLimit &= ^LimitChiHu
	u.UserLimit &= ^LimitPeng
	u.UserLimit &= ^LimitGang

	var bSysOut bool
	if len(args) == 3 {
		bSysOut = args[2].(bool)
	}
	room.GetDataMgr().NotifySendCard(u, CardData, bSysOut)

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

// 吃碰杠胡各种操作
func (room *ZP_base) UserOperateCard(args []interface{}) {
	log.Debug("???????????????????????????????1111111111111111111")
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
			log.Debug("user:%d,OperateCode:%d,OperateCard：%v", u.ChairId, OperateCode, OperateCard)
			retcode = ErrNoOperator
			return
		}

		//变量定义
		cbTargetAction, wTargetUser := room.GetDataMgr().CheckUserOperator(u, room.UserMgr.GetMaxPlayerCnt(), OperateCode, OperateCard)
		if cbTargetAction < 0 {
			log.Debug("wait other user")
			return
		}

		//放弃操作
		if cbTargetAction == WIK_NULL {
			//用户状态
			if room.GetDataMgr().DispatchCardData(room.GetDataMgr().GetResumeUser(), room.GetDataMgr().GetGangStatus() != WIK_GANERAL) > 0 {
				room.OnEventGameConclude(room.GetDataMgr().GetProvideUser(), nil, GER_NORMAL)
			}
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
				room.OnEventGameConclude(room.GetDataMgr().GetProvideUser(), nil, GER_NORMAL)
			}
		}
	} else {

		//扑克效验
		if (OperateCode != WIK_NULL) && (OperateCode != WIK_CHI_HU) && (!room.LogicMgr.IsValidCard(OperateCard[0])) {
			log.Error("OperateCode != WIK_NULL) && (OperateCode != WIK_CHI_HU) && (!room.LogicMgr.IsValidCard(OperateCard[0])")
			return
		}

		//设置变量
		room.GetDataMgr().ResetUserOperateEx(u)

		//执行动作
		switch OperateCode {
		case WIK_GANG: //杠牌操作
			cbGangKind := room.GetDataMgr().AnGang(u, OperateCode, OperateCard)
			if cbGangKind == 0 {
				return
			}

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
