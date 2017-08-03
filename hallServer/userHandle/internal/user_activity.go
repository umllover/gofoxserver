package internal

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/db/model/stats"
	"mj/hallServer/db/model/base"
	"mj/hallServer/user"
	"time"
)

//领取奖励
func (m *UserModule) DrawSahreAward(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_DrawSahreAward)
	retMsg := &msg.L2C_DrawSahreAwardResult{}
	player := m.a.UserData().(*user.User)
	defer func() {
		player.WriteMsg(retMsg)
	}()

	template, ok := base.ActivityCache.Get(recvMsg.DrawId)
	if !ok {
		retMsg.RetCode = ErrNotFoudTemplate
		return
	}

	times := player.GetTimesByType(template.Id, template.DrawType)
	if times >= template.DrawTimes {
		retMsg.RetCode = ErrMaxDrawTimes
		return
	}

	player.IncreaseTimesByType(template.Id, 1, template.DrawType)

	switch template.ItemType {
	case 1:
		player.AddCurrency(template.Amount)
	}
	now := time.Now()
	stats.DrawAwardLogOp.Insert(&stats.DrawAwardLog{
		Id:          template.Id,
		DrawId:      recvMsg.DrawId,
		Description: template.Description,
		DrawCount:   template.DrawTimes,
		DrawType:    template.DrawType,
		Amount:      template.Amount,
		ItemType:    template.ItemType,
		DrawTime:    &now,
	})

}

//玩家请求次数信息
func (m *UserModule) ReqTimesInfo(args []interface{}) {
	player := m.a.UserData().(*user.User)
	player.WriteMsg(player.GetTimeInfo())
}
