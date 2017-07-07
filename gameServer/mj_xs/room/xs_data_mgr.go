package room

import (
	"mj/common/msg"
	"mj/common/msg/mj_xs_msg"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
)

func NewXSDataMgr(id, uid, configIdx int, name string, temp *base.GameServiceOption, base *xs_entry) *xs_data {
	d := new(xs_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base)
}

type xs_data struct {
	*mj_base.RoomData
	ZhuaHuaCnt   int  //扎花个数
	ZhuaHuaScore int  //扎花分数
	FengQaun     int  //风圈
	IsFirst      bool //是否首发
}

func (room *xs_data) AfterStartGame() {
	//检查自摸
	room.CheckZiMo()
	//通知客户端开始了
	room.SendGameStart()
}

//发送开始
func (room *xs_data) SendGameStart() {
	//构造变量
	GameStart := &mj_xs_msg.G2C_MJXS_GameStart{}
	GameStart.BankerUser = room.BankerUser
	GameStart.SiceCount = room.SiceCount
	GameStart.SunWindCount = 0
	GameStart.LeftCardCount = room.LeftCardCount
	GameStart.First = room.IsFirst
	GameStart.FengQuan = room.FengQaun
	GameStart.InitialBankerUser = room.BankerUser
	//发送数据
	room.MjBase.UserMgr.ForEachUser(func(u *user.User) {
		GameStart.UserAction = room.UserAction[u.ChairId]
		GameStart.CardData = room.MjBase.LogicMgr.GetUserCards(room.CardIndex[u.ChairId])
		u.WriteMsg(GameStart)
	})
}

//发送操作结果
func (room *xs_data) SendOperateResult(u *user.User, wrave *msg.WeaveItem) {
	OperateResult := &mj_xs_msg.G2C_MJXS_OperateResult{}
	OperateResult.ProvideUser = wrave.ProvideUser
	OperateResult.OperateCode = wrave.WeaveKind
	OperateResult.OperateCard = wrave.CenterCard
	if u != nil {
		OperateResult.OperateUser = u.ChairId
	} else {
		OperateResult.OperateUser = wrave.OperateUser
		OperateResult.ActionMask = wrave.ActionMask
	}
	room.MjBase.UserMgr.SendMsgAll(OperateResult)
}
