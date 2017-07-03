package room

import (
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"
)

func NewHZDataMgr(id, uid, configIdx int, name string, temp *base.GameServiceOption, base *hz_entry) *hz_data {
	d := new(hz_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base)
	return d
}

type hz_data struct {
	*mj_base.RoomData
	ZhuaHuaCnt   int //扎花个数
	ZhuaHuaScore int //扎花分数
}

func (room *hz_data) OnZhuaHua(CenterUser int) (CardData []int, BuZhong []int) {
	count := room.ZhuaHuaCnt
	if count == 0 {
		return
	}

	isWin := false
	for chairId, v := range room.UserAction {
		if v&mj_base.WIK_CHI_HU != 0 && chairId == CenterUser {
			isWin = true
		}
	}

	if !isWin {
		return
	}

	//抓花规则
	getInedx := [3]int{1, 5, 9}
	for i := 0; i < count; i++ {
		room.LeftCardCount--
		cardData := room.RepertoryCard[room.LeftCardCount]
		cardColor := room.MjBase.LogicMgr.GetCardColor(cardData)
		cardValue := room.MjBase.LogicMgr.GetCardValue(cardData)
		if cardColor == 3 {
			if cardValue >= 5 {
				//中发白
				temp := cardValue - 4
				if temp == getInedx[0] || temp == getInedx[1] || temp == getInedx[2] {
					CardData = append(CardData, cardData)
					room.ZhuaHuaScore++
				}
			}
		} else if cardColor >= 0 && cardColor <= 2 {
			if cardValue == getInedx[0] || cardValue == getInedx[1] || cardValue == getInedx[2] {
				CardData = append(CardData, cardData)
				room.ZhuaHuaScore++
			}
		}
		BuZhong = append(BuZhong, cardData)
	}

	return
}
