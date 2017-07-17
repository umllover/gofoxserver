package room

import (
	"encoding/json"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/mj_xs_msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"

	"github.com/lovelly/leaf/log"

	"github.com/lovelly/leaf/timer"
)

func NewXSDataMgr(id int, uid int64, configIdx int, name string, temp *base.GameServiceOption, base *xs_entry, set string) *xs_data {
	d := new(xs_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base)

	//房间游戏设置
	info := make(map[string]interface{})
	err := json.Unmarshal([]byte(set), &info)
	if err != nil {
		log.Error("zpmj at NewDataMgr error:%s", err.Error())
		return nil
	}

	return d
}

type xs_data struct {
	*mj_base.RoomData
	ZhuaHuaCnt   int  //扎花个数
	ZhuaHuaScore int  //扎花分数
	FengQaun     int  //风圈
	IsFirst      bool //是否首发
}

func (room *xs_data) InitRoom(UserCnt int) {
	//初始化
	log.Debug("mjxs at InitRoom")
	room.RepertoryCard = make([]int, room.GetCfg().MaxRepertory)
	room.CardIndex = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.CardIndex[i] = make([]int, room.GetCfg().MaxIdx)
	}
	room.ChiHuKind = make([]int, UserCnt)
	room.ChiPengCount = make([]int, UserCnt)
	room.GangCard = make([]bool, UserCnt) //杠牌状态
	room.GangCount = make([]int, UserCnt)
	room.Ting = make([]bool, UserCnt)
	room.UserAction = make([]int, UserCnt)
	room.DiscardCard = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.DiscardCard[i] = make([]int, 60)
	}
	room.UserGangScore = make([]int, UserCnt)
	room.WeaveItemArray = make([][]*msg.WeaveItem, UserCnt)
	room.ChiHuRight = make([]int, UserCnt)
	room.HeapCardInfo = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.HeapCardInfo[i] = make([]int, 2)
	}
	room.OperateTime = make([]*timer.Timer, UserCnt)

	room.LeftCardCount = room.GetCfg().MaxRepertory
	room.UserActionDone = false
	room.SendStatus = Not_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideGangUser = INVALID_CHAIR
	room.HistoryScores = make([]*HistoryScore, UserCnt)
	room.MinusLastCount = 0
	room.MinusHeadCount = 0
	room.OutCardCount = 0

	//设置xs麻将牌数据
	room.EndLeftCount = 16
	room.ZhuaHuaScore = 0
	room.FlowerCnt = [4]int{}
	room.BanCardCnt = [4][9]int{}
	room.BanUser = [4]int{}

	room.IsResponse = make([]bool, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.OperateCard[i] = make([]int, 60)
	}
	log.Debug("len1 OperateCard: %d %d", len(room.OperateCard), len(room.OperateCard[1]))
	room.PerformAction = make([]int, UserCnt)
}

func (room *xs_data) BeforeStartGame(UserCnt int) {
	log.Debug("###################### BeforeStartGame")
	room.InitRoom(UserCnt)
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
	GameStart := &mj_xs_msg.G2C_GameStart{}
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
	OperateResult := &mj_xs_msg.G2C_OperateResult{}
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
