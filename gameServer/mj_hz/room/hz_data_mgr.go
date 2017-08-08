package room

import (
	. "mj/common/cost"
	"mj/common/msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/db/model/base"

	"mj/common/utils"

	"github.com/lovelly/leaf/log"
)

type hz_data struct {
	*mj_base.RoomData
	ZhuaHuaCnt   int        //扎花个数
	ZhuaHuaScore int        //扎花分数
	ZhuaHuaMap   []*HuaUser //插花数据
}

//抓花结构体子项
type HuaUser struct {
	ChairID int  //椅子号
	Card    int  //牌值
	IsZhong bool //是否中花
}

func NewHZDataMgr(id int, uid int64, configIdx int, name string, temp *base.GameServiceOption, base *hz_entry, info *msg.L2G_CreatorRoom) *hz_data {
	d := new(hz_data)
	d.RoomData = mj_base.NewDataMgr(id, uid, configIdx, name, temp, base.Mj_base, info)

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

//计算抓花
func (room *hz_data) CalcZhuahua(winUser []int) {
	if room.ZhuaHuaCnt == 0 {
		return
	}
	ZhongCard, BuZhong := room.OnZhuaHua(winUser)
	for k, v := range winUser {
		for _, cardV := range ZhongCard[k] {
			for {
				randV, randOk := utils.RandInt(0, room.ZhuaHuaCnt-1)
				if randOk == nil && room.ZhuaHuaMap[randV] == nil {
					huaUser := &HuaUser{}
					huaUser.Card = cardV
					log.Debug("中花：%d", cardV)
					huaUser.ChairID = v
					huaUser.IsZhong = true
					room.ZhuaHuaMap[randV] = huaUser
					break
				}
			}
		}
	}
	for _, cardV2 := range BuZhong {
		for {
			randV, randOk := utils.RandInt(0, room.ZhuaHuaCnt-1)
			if randOk == nil && room.ZhuaHuaMap[randV] == nil {
				huaUser := &HuaUser{}
				huaUser.Card = cardV2
				//huaUser.ChairID = v
				log.Debug("不中花：%d", cardV2)
				huaUser.IsZhong = false
				room.ZhuaHuaMap[randV] = huaUser
				break
			}
		}
	}
}

//抓花
func (room *hz_data) OnZhuaHua(winUser []int) (CardData [][]int, BuZhong []int) {
	log.Debug("OnZhuaHua Start...")
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

	log.Debug("===========isWin=%v", isWin)
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
		} else {
			BuZhong = append(BuZhong, cardData)
		}
	}
	return
}
