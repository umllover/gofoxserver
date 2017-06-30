package pk_base

import (
	/*"mj/common/msg"
	"mj/gameServer/common"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util" */
	//"fmt"
	"github.com/lovelly/leaf/util"
)

// 扑克通用逻辑
const (
	LOGIC_MASK_COLOR	=			0xF0								//花色掩码
	LOGIC_MASK_VALUE	=			0x0F								//数值掩码
)
//获取数值
func GetCardValue(CardData int) int {
	return CardData&LOGIC_MASK_VALUE
}
//获取花色
func GetCardColor(CardData int) int {
	return CardData&LOGIC_MASK_COLOR
}



//排列扑克
func  SortCardList(cardData []int, cardCount int)  {
	var logicValue []int
	for i:=0;i<cardCount;i++ {
		logicValue[i] = GetCardValue(cardData[i])
	}
	sorted := true
	last := cardCount -1
	for {
		sorted = true
		for i := 0; i < last; i++ {
			if (logicValue[i] < logicValue[i+1]) || (logicValue[i] == logicValue[i+1] && (cardData[i] < cardData[i+1])) {
				tempData := cardData[i]
				cardData[i] = cardData[i+1]
				cardData[i+1] = tempData
				tempData = logicValue[i]
				logicValue[i] = logicValue[i+1]
				logicValue[i+1] = tempData
				sorted = false
			}
		}
		last--
		if sorted==true {
			break
		}
	}
}

/*
//混乱扑克
func  RandCardList(cardBuffer []int, cardBufferCount int)  {
	cardData := GetNormalCards()
	randCount := 0
	position := 0

	for {

		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		r := random.Int()
		position = r % (len(cardData) - randCount)
		cardBuffer[randCount] = cardData[position]
		cardData[position] = cardData[len(cardData)-randCount]
		if (randCount >= cardBufferCount) {
			break
		}
	}

}
*/

//混乱扑克
func RandCardList(cbCardBuffer, OriDataArray []int) {
	//混乱准备
	cbBufferCount := int(len(cbCardBuffer))
	cbCardDataTemp := make([]int, cbBufferCount)
	util.DeepCopy(&cbCardDataTemp, &OriDataArray)

	//混乱扑克
	var cbRandCount int
	var cbPosition int
	for {
		if cbRandCount >= cbBufferCount {
			break
		}
		cbPosition = int(util.RandInterval(0, int(cbBufferCount-cbRandCount)))
		cbCardBuffer[cbRandCount] = cbCardDataTemp[cbPosition]
		cbRandCount++
		cbCardDataTemp[cbPosition] = cbCardDataTemp[cbBufferCount-cbRandCount]
	}

	return
}

