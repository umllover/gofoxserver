package room

import (
	"encoding/json"
	"mj/gameServer/db/model"

	"sync"

	"time"

	"github.com/lovelly/leaf/log"
)

var (
	HappyCards     HappySendCardList
	HappyCardsKing HappyKingSendCardList

	NoKingLock sync.RWMutex
	KingLock   sync.RWMutex
)

// 欢乐场发牌选择表
type HappySendCardList struct {
	index      int            // 当前更新到的索引
	count      int            // 当前表里有多少条记录
	StartTime  int            // 起始时间，小于该时间的不更新到数据库
	CreateTime [10000]int     // 创建时间
	Cards      [10000][54]int // 出牌记录列表
}

type HappyKingSendCardList struct {
	index      int            // 当前更新到的索引
	count      int            // 当前表里有多少条记录
	StartTime  int            // 起始时间，小于该时间的不更新到数据库
	CreateTime [10000]int     // 创建时间
	Cards      [10000][60]int // 出牌记录列表
}

// 初始化欢乐场发牌表
func InitHappyCardListData() {
	log.Debug("初始化欢乐场发牌表")
	// 取非八王表
	allData, err := model.RecordOutcardDdzOp.SelectAll()
	if err == nil {
		nCount := len(allData)
		if nCount > 0 {
			n := nCount - 10000
			if n < 0 {
				n = 0
				HappyCards.index = nCount
			}

			data := allData[n:]

			if nCount > 10000 {
				nCount = 10000
			}
			HappyCards.count = nCount
			for i := 0; i < nCount; i++ {
				HappyCards.CreateTime[i] = data[i].CreateTime
				json.Unmarshal([]byte(data[i].CardData), &HappyCards.Cards[i])
				log.Debug("非八王取到的牌%d,%v", HappyCards.CreateTime[i], HappyCards.Cards[i])
			}

			// 删除旧数据
			for i := 0; i < n; i++ {
				model.RecordOutcardDdzKingOp.Delete(allData[i].RecordID)
			}
		}
	}
	HappyCards.StartTime = int(time.Now().UnixNano() / 1000000000)
	// 取八王表
	allDataKing, err := model.RecordOutcardDdzKingOp.SelectAll()
	if err == nil {
		nCount := len(allDataKing)
		if nCount > 0 {
			n := nCount - 10000
			if n < 0 {
				n = 0
				HappyCardsKing.index = nCount
			}

			data := allDataKing[n:]

			if nCount > 10000 {
				nCount = 10000
			}
			HappyCardsKing.count = nCount
			for i := 0; i < nCount; i++ {
				HappyCardsKing.CreateTime[i] = data[i].CreateTime
				json.Unmarshal([]byte(data[i].CardData), &HappyCardsKing.Cards[i])
				log.Debug("八王取到的牌%d,%v", HappyCardsKing.CreateTime[i], HappyCardsKing.Cards[i])
			}

			// 删除旧数据
			for i := 0; i < n; i++ {
				model.RecordOutcardDdzKingOp.Delete(allDataKing[i].RecordID)
			}
		}
	}
	HappyCardsKing.StartTime = int(time.Now().UnixNano() / 1000000000)
}

// 更新非八王表
func UpdateHappyCardList(cards []int) {
	if len(cards) != 54 {
		return
	}
	NoKingLock.Lock()
	log.Debug("%d更新前的数据%d,%v", HappyCards.index, HappyCards.CreateTime[HappyCards.index], HappyCards.Cards[HappyCards.index])
	nIndex := HappyCards.index
	HappyCards.index = (HappyCards.index + 1) % 10000
	if HappyCards.count < 10000 {
		HappyCards.count++
	}
	for i := 0; i < 54; i++ {
		HappyCards.Cards[nIndex][i] = cards[i]
	}
	HappyCards.CreateTime[nIndex] = int(time.Now().UnixNano() / 1000000000)
	log.Debug("%d更新后的数据%d,%v", HappyCards.index, HappyCards.CreateTime[nIndex], HappyCards.Cards[nIndex])
	NoKingLock.Unlock()
}

// 更新八王表
func UpdatehappyKingCardList(cards []int) {
	if len(cards) != 60 {
		return
	}
	KingLock.Lock()
	log.Debug("%d更新前的数据%d,%v", HappyCardsKing.index, HappyCardsKing.CreateTime[HappyCardsKing.index], HappyCardsKing.Cards[HappyCardsKing.index])
	nIndex := HappyCardsKing.index
	HappyCardsKing.index = (HappyCardsKing.index + 1) % 10000
	if HappyCardsKing.count < 10000 {
		HappyCardsKing.count++
	}
	for i := 0; i < 60; i++ {
		HappyCardsKing.Cards[nIndex][i] = cards[i]
	}
	HappyCardsKing.CreateTime[nIndex] = int(time.Now().UnixNano() / 1000000000)
	log.Debug("%d更新后的数据%d,%v", HappyCardsKing.index, HappyCardsKing.CreateTime[nIndex], HappyCardsKing.Cards[nIndex])
	KingLock.Unlock()
}

// 服务器关闭的时候，选1000条存到数据库里
func SaveDataToDB() {
	log.Debug("非八王场个数%d,八王个数%d", HappyCards.count, HappyCardsKing.count)
	if HappyCards.count > 0 {
		nCount := HappyCards.index
		if nCount > 1000 {
			nCount = 1000
		}
		var nIndex int
		var dbCardData model.RecordOutcardDdz
		for i := 0; i < 1000; i++ {
			nIndex = (HappyCards.index + 10000 + i - nCount) % 10000
			//log.Debug("循环%d,%d,%d,%d", nIndex, HappyCards.CreateTime[nIndex], HappyCards.StartTime, HappyCards.index)
			if HappyCards.CreateTime[nIndex] < HappyCards.StartTime {
				// 旧数据不存表，避免重复
				continue
			}
			cardData, err := json.Marshal(HappyCards.Cards[nIndex])
			if err == nil {
				dbCardData.CardData = string(cardData)
				dbCardData.CreateTime = HappyCards.CreateTime[nIndex]
				model.RecordOutcardDdzOp.Insert(&dbCardData)
				log.Debug("更新非八王表%v", dbCardData)
			}
		}
	}
	if HappyCardsKing.count > 0 {
		nCount := HappyCards.index
		if nCount > 1000 {
			nCount = 1000
		}
		var nIndex int
		var dbCardData model.RecordOutcardDdzKing
		for i := 0; i < 1000; i++ {
			nIndex = (HappyCardsKing.index + 10000 + i - nCount) % 10000
			if HappyCardsKing.CreateTime[nIndex] < HappyCardsKing.StartTime {
				// 旧数据不存表，避免重复
				continue
			}
			cardData, err := json.Marshal(HappyCardsKing.Cards[nIndex])
			if err == nil {
				dbCardData.CardData = string(cardData)
				dbCardData.CreateTime = HappyCardsKing.CreateTime[nIndex]
				model.RecordOutcardDdzKingOp.Insert(&dbCardData)
				log.Debug("更新八王表%v", dbCardData)
			}
		}
	}
}
