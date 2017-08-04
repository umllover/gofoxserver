package room

import (
	. "mj/common/cost"
	"mj/common/msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"

	"github.com/lovelly/leaf/log"
)

func NewHZDataMgr(id int, uid int64, configIdx int, name string, temp *base.GameServiceOption, base *hz_entry, info *msg.L2G_CreatorRoom) *hz_data {
	d := new(hz_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base, info.OtherInfo)

	getData, ok := d.OtherInfo["zhaMa"].(float64)
	if !ok {
		log.Error("hzmj at NewDataMgr [zhaMa] error")
		return nil
	}

	//TODO 客户端发的个数有误，暂时强制改掉
	getData = 2

	d.ZhuaHuaCnt = int(getData)

	return d
}

type hz_data struct {
	*mj_base.RoomData
	ZhuaHuaCnt   int //扎花个数
	ZhuaHuaScore int //扎花分数
}

func (room *hz_data) BeforeStartGame(UserCnt int) {
	room.InitRoom(UserCnt)
}

func (room *hz_data) InitRoom(UserCnt int) {
	//初始化
	log.Debug("初始化红中房间")
	room.RepertoryCard = make([]int, room.GetCfg().MaxRepertory)
	room.CardIndex = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.CardIndex[i] = make([]int, room.GetCfg().MaxIdx)
	}
	room.FlowerCnt = [4]int{}
	room.ChiHuKind = make([]int, UserCnt)
	room.ChiPengCount = make([]int, UserCnt)
	room.GangCard = make([]bool, UserCnt) //杠牌状态
	room.GangCount = make([]int, UserCnt)
	room.Ting = make([]bool, UserCnt)
	room.UserAction = make([]int, UserCnt)
	room.DiscardCard = make([][]int, UserCnt)
	room.UserGangScore = make([]int, UserCnt)
	room.WeaveItemArray = make([][]*msg.WeaveItem, UserCnt)
	room.ChiHuRight = make([]int, UserCnt)
	room.HeapCardInfo = make([][]int, UserCnt)
	room.IsResponse = make([]bool, UserCnt)
	room.PerformAction = make([]int, UserCnt)
	room.OperateCard = make([][]int, UserCnt)
	for i := 0; i < UserCnt; i++ {
		room.HeapCardInfo[i] = make([]int, 2)
	}
	room.UserActionDone = false
	room.SendStatus = Not_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideGangUser = INVALID_CHAIR
	room.MinusLastCount = 0
	room.MinusHeadCount = room.GetCfg().MaxRepertory
	room.OutCardCount = 0
	//扎码数
	room.EndLeftCount = room.ZhuaHuaCnt
}

//抓花
func (room *hz_data) OnZhuaHua(winUser []int) (CardData [][]int, BuZhong []int) {
	count := room.ZhuaHuaCnt
	if count == 0 {
		return
	}

	isWin := false
	for chairId, v := range room.UserAction {
		if v&WIK_CHI_HU != 0 && chairId == winUser[0] {
			isWin = true
		}
	}

	if !isWin {
		return
	}
	CardData = make([][]int, count)
	//抓花规则
	getInedx := [3]int{1, 5, 9}
	for i := 0; i < count; i++ {
		cardData := room.GetHeadCard()
		cardColor := room.MjBase.LogicMgr.GetCardColor(cardData)
		cardValue := room.MjBase.LogicMgr.GetCardValue(cardData)
		if cardColor == 3 {
			if cardValue >= 5 {
				//中发白
				temp := cardValue - 4
				if temp == getInedx[0] || temp == getInedx[1] || temp == getInedx[2] {
					CardData[0] = append(CardData[0], cardData)
					room.ZhuaHuaScore++
				}
			}
		} else if cardColor >= 0 && cardColor <= 2 {
			if cardValue == getInedx[0] || cardValue == getInedx[1] || cardValue == getInedx[2] {
				CardData[0] = append(CardData[0], cardData)
				room.ZhuaHuaScore++
			}
		}
		BuZhong = append(BuZhong, cardData)
	}
	return
}
