package pk_base

import (
	"github.com/lovelly/leaf/util"
	"time"
	"math/rand"
)

type BaseLogic struct {
	ConfigIdx int //配置索引
}

func NewBaseLogic(ConfigIdx int) *BaseLogic {
	bl := new(BaseLogic)
	bl.ConfigIdx = ConfigIdx
	return bl
}

func (lg *BaseLogic) GetCfg() *PK_CFG {
	return GetCfg(lg.ConfigIdx)
}

func (lg *BaseLogic) RandCardList(cbCardBuffer, OriDataArray []int) {

	if !(len(OriDataArray) >= len(cbCardBuffer)) {
		return
	}
	//混乱准备
	cbBufferCount := int(len(cbCardBuffer))
	cbCardDataTemp := make([]int, cbBufferCount)
	util.DeepCopy(&cbCardDataTemp, &OriDataArray)

	//随机交换两张牌
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i:=0; i<len(cbCardDataTemp); i++ {
		t1 := r.Int() % len(cbCardDataTemp)
		t2 := r.Int() % len(cbCardDataTemp)
		temp := cbCardDataTemp[t1]
		cbCardDataTemp[t1] = cbCardDataTemp[t2]
		cbCardDataTemp[t2] = temp
	}

	// 赋给输出缓冲
	for i:=0; i< len(cbCardBuffer); i++ {
		cbCardBuffer[i] = cbCardDataTemp[i]
	}
	return
}

//排列扑克
func (lg *BaseLogic) SortCardList(cardData []int, cardCount int) {
	logicValue := make([]int, cardCount)
	for i := 0; i < cardCount; i++ {
		logicValue[i] = lg.GetCardValue(cardData[i])
	}
	sorted := true
	last := cardCount - 1
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
		if sorted == true {
			break
		}
	}
}

//获取数值
func (lg *BaseLogic) GetCardValue(CardData int) int {
	return CardData & LOGIC_MASK_VALUE
}

//获取花色
func (lg *BaseLogic) GetCardColor(CardData int) int {
	return CardData & LOGIC_MASK_COLOR
}




func (lg *BaseLogic) CompareCard(firstCardData []int, lastCardData []int) bool {
	return false
}
func (lg *BaseLogic) GetCardType(cardData []int) int {
	return 0
}
