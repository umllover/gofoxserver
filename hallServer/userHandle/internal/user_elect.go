package internal

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/hallServer/common"
	"mj/hallServer/db/model"
	"mj/hallServer/db/model/stats"
	"mj/hallServer/user"
	"time"
)

//填写推荐人信息
func (m *UserModule) SetElect(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_SetElect)
	retMsg := &msg.L2C_SetElectResult{}
	player := m.a.UserData().(*user.User)
	defer func() {
		player.WriteMsg(retMsg)
	}()
	_, ok := model.UserattrOp.Get(recvMsg.ElectUid)
	if !ok {
		retMsg.RetCode = ErrNotFoudPlayer
	}
	player.ElectUid = recvMsg.ElectUid
	model.UserattrOp.UpdateWithMap(player.Id, map[string]interface{}{
		"elect_uid": player.ElectUid,
	})

	model.UserSpreadOp.Insert(&model.UserSpread{
		UserId:    recvMsg.ElectUid,
		SpreadUid: player.Id,
	})
	now := time.Now()

	stats.RecommendLogOp.Insert(&stats.RecommendLog{
		SubElectUid: player.Id,
		ElectUid:    recvMsg.ElectUid,
		ElectTime:   &now,
	})

}

//领取推举人奖励
func (m *UserModule) DrawElectAward(player *user.User, cnt int) {
	award := common.GetGlobalVarInt(MAX_ELECT_AWARD)
	list, err := model.UserSpreadOp.QueryByMap(map[string]interface{}{
		"user_id": player.Id,
	})
	if err != nil {
		return
	}

	for _, v := range list {
		model.UserSpreadOp.UpdateWithMap(v.UserId, v.SpreadUid, map[string]interface{}{
			"status": 1,
		})
		player.AddCurrency(award)
	}
}
