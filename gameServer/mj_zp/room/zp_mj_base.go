package room

import (
	. "mj/common/cost"
	"mj/gameServer/common/mj_base"
	"mj/gameServer/db/model"
	"mj/gameServer/user"

	"mj/common/msg/mj_zp_msg"

	"mj/common/msg/mj_hz_msg"

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
	recvMsg := args[0].(*mj_zp_msg.C2G_ZPMJ_OutCard)
	u := args[1].(*user.User)
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
		retcode = ErrNotSelfOut
		return
	}

	if !room.LogicMgr.IsValidCard(recvMsg.CardData) {
		log.Error("at OnUserOutCard IsValidCard card ")
		retcode = NotValidCard
	}

	//删除扑克
	if !room.LogicMgr.RemoveCard(room.DataMgr.GetUserCardIndex(u.ChairId), recvMsg.CardData) {
		log.Error("at OnUserOutCard not have card ")
		retcode = ErrNotFoudCard
		return
	}

	u.UserLimit |= ^LimitChiHu
	u.UserLimit |= ^LimitPeng
	u.UserLimit |= ^LimitGang

	room.DataMgr.NotifySendCard(u, recvMsg.CardData, false)

	//响应判断
	bAroseAction := room.DataMgr.EstimateUserRespond(u.ChairId, recvMsg.CardData, EstimatKind_OutCard)

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
	recvMsg := args[0].(*mj_hz_msg.C2G_HZMJ_OperateCard)
	u := args[1].(*user.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	if room.DataMgr.GetCurrentUser() == INVALID_CHAIR {
		//效验状态
		if !room.DataMgr.HasOperator(u.ChairId, recvMsg.OperateCode) {
			retcode = ErrNoOperator
			return
		}

		//变量定义
		cbTargetAction, wTargetUser := room.DataMgr.CheckUserOperator(u, room.UserMgr.GetMaxPlayerCnt(), recvMsg)
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
		if room.DataMgr.RemoveCardByOP(wTargetUser, cbTargetAction) {
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
		if (recvMsg.OperateCode != WIK_NULL) && (recvMsg.OperateCode != WIK_CHI_HU) && (!room.LogicMgr.IsValidCard(recvMsg.OperateCard[0])) {
			return
		}

		//设置变量
		// room.UserAction[room.CurrentUser] = WIK_NULL

		//执行动作
		switch recvMsg.OperateCode {
		case WIK_GANG: //杠牌操作
			cbGangKind := room.DataMgr.AnGang(u, recvMsg.OperateCode, recvMsg.OperateCard)
			//效验动作
			bAroseAction := false
			if cbGangKind == WIK_MING_GANG {
				bAroseAction = room.DataMgr.EstimateUserRespond(u.ChairId, recvMsg.OperateCard[0], EstimatKind_GangCard)
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
