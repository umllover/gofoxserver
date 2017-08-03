package internal

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/hallServer/db/model/base"
	dataLog "mj/hallServer/log"
	"mj/hallServer/user"
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

	drawAward := dataLog.DrawAwardLog{}
	drawAward.AddDrawAdardLog(template.Id, recvMsg.DrawId, template.Description, template.DrawTimes, template.DrawType, template.Amount, template.ItemType)

}
