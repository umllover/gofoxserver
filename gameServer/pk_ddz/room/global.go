package room

import (
	"encoding/json"
	"mj/gameServer/db/model"

	"github.com/lovelly/leaf/log"
)

var (
	HappyCards     HappySendCardList
	HappyCardsKing HappyKingSendCardList
)

// 欢乐场发牌选择表
type HappySendCardList struct {
	index int            // 当前更新到的索引
	count int            // 当前表里有多少条记录
	Cards [10000][54]int // 出牌记录列表
}

type HappyKingSendCardList struct {
	index int            // 当前更新到的索引
	count int            // 当前表里有多少条记录
	Cards [10000][60]int // 出牌记录列表
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
				json.Unmarshal([]byte(data[i].CardData), &HappyCards.Cards[i])
				log.Debug("非八王取到的牌%v", HappyCards.Cards[i])
			}

			// 删除旧数据
			for i := 0; i < n; i++ {
				model.RecordOutcardDdzKingOp.Delete(allData[i].RecordID)
			}
		}
	}
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
			log.Debug("data.len=%d,count=%d,n=%d", len(data), nCount, n)
			HappyCardsKing.count = nCount
			for i := 0; i < nCount; i++ {
				json.Unmarshal([]byte(data[i].CardData), &HappyCardsKing.Cards[i])
				log.Debug("%d八王取到的牌%v", i, HappyCardsKing.Cards[i])
			}

			// 删除旧数据
			for i := 0; i < n; i++ {
				model.RecordOutcardDdzKingOp.Delete(allDataKing[i].RecordID)
			}
		}
	}
}
