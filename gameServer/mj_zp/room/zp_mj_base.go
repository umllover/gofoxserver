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
			room.OnEventGameConclude(GER_NORMAL)
		}
	}

	return
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

	u := room.UserMgr.GetUserByChairId(wChairID)
	if u == nil {
		return false
	}

	if bTrustee {
		if wChairID == room.GetDataMgr().GetCurrentUser() && !room.GetDataMgr().IsActionDone() {
			cardindex := room.GetDataMgr().GetTrusteeOutCard(wChairID)
			if cardindex == INVALID_BYTE {
				return false
			}
			card := room.LogicMgr.SwitchToCardData(cardindex)
			room.OutCard([]interface{}{u, card, true})
		} else if room.GetDataMgr().GetCurrentUser() == INVALID_CHAIR && !room.GetDataMgr().IsActionDone() {
			operateCard := []int{0, 0, 0}
			room.UserOperateCard([]interface{}{u, WIK_NULL, operateCard})
		}
		//TODO 测试 启动机器人
		//u.RunRobot()
	} else {
		//TODO 测试 关闭机器人
		//u.StopRobot()
	}
	return true
}
