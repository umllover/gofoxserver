package room

import (
	"mj/gameServer/common/pk/pk_base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

// 十三水逻辑
const (
	CT_INVALID                     = iota //错误类型
	CT_SINGLE                             //单牌类型
	CT_ONE_DOUBLE                         //只有一对
	CT_FIVE_TWO_DOUBLE                    //两对牌型
	CT_THREE                              //三张牌型
	CT_FIVE_MIXED_FLUSH_NO_A              //没A杂顺
	CT_FIVE_MIXED_FLUSH_FIRST_A           //A在前顺子
	CT_FIVE_MIXED_FLUSH_BACK_A            //A在后顺子
	CT_FIVE_FLUSH                         //同花五牌
	CT_FIVE_THREE_DEOUBLE                 //三条一对
	CT_FIVE_FOUR_ONE                      //四带一张
	CT_FIVE_STRAIGHT_FLUSH_NO_A           //没A同花顺
	CT_FIVE_STRAIGHT_FLUSH_FIRST_A        //A在前同花顺
	CT_FIVE_STRAIGHT_FLUSH_BACK_A         //A在后同花顺
)

//特殊牌型
const (
	CT_THIRTEEN_FLUSH      = 26 //同花十三水
	CT_THIRTEEN            = 25 //十三水
	CT_TWELVE_KING         = 24 //十二皇族
	CT_THREE_STRAIGHTFLUSH = 23 //三同花顺
	CT_THREE_BOMB          = 22 //三炸弹
	CT_ALL_BIG             = 21 //全大
	CT_ALL_SMALL           = 20 //全小
	CT_SAME_COLOR          = 19 //凑一色
	CT_FOUR_THREESAME      = 18 //四套冲三
	CT_FIVEPAIR_THREE      = 17 //五对冲三
	CT_SIXPAIR             = 16 //六对半
	CT_THREE_FLUSH         = 15 //三同花
	CT_THREE_STRAIGHT      = 14 //三顺子

	LX_ONEPARE       = 13 //一对
	LX_TWOPARE       = 14 //两对
	LX_THREESAME     = 15 //三条
	LX_STRAIGHT      = 16 //顺子
	LX_FLUSH         = 17 //同花
	LX_GOURD         = 18 //葫芦
	LX_FOURSAME      = 19 //铁支
	LX_STRAIGHTFLUSH = 20 //同花顺
)

//数值掩码
const (
	LOGIC_MASK_COLOR = 0xF0 //花色掩码
	LOGIC_MASK_VALUE = 0x0F //数值掩码
)

type TagAnalyseItem struct {
	bOneCount   int   //单张数目
	bTwoCount   int   //两张数目
	bThreeCount int   //三张数目
	bFourCount  int   //四张数目
	bFiveCount  int   //五张数目
	bOneFirst   []int //单牌位置
	bTwoFirst   []int //对牌位置
	bThreeFirst []int //三条位置
	bFourFirst  []int //四张位置
	bStraight   bool  //是否顺子
}

type tagAnalyseType struct {
}

//分析结构
type tagAnalyseData struct {
	bOneCount   int   //单张数目
	bTwoCount   int   //两张数目
	bThreeCount int   //三张数目
	bFourCount  int   //四张数目
	bFiveCount  int   //五张数目
	bOneFirst   []int //单牌位置
	bTwoFirst   []int //对牌位置
	bThreeFirst []int //三条位置
	bFourFirst  []int //四张位置
	bStraight   bool  //是否顺子
}

func NewSssZLogic(ConfigIdx int) *sss_logic {
	l := new(sss_logic)
	l.BtCardSpecialData = make([]int, 13)
	l.BaseLogic = pk_base.NewBaseLogic(ConfigIdx)
	return l
}

type sss_logic struct {
	*pk_base.BaseLogic
	BtCardSpecialData []int
}

func (lg *sss_logic) RandCardList(cbCardBuffer, OriDataArray []int) {

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
		cbPosition = int(util.RandInterval(0, int(cbBufferCount-cbRandCount-1)))
		cbCardBuffer[cbRandCount] = cbCardDataTemp[cbPosition]
		cbRandCount++
		cbCardDataTemp[cbPosition] = cbCardDataTemp[cbBufferCount-cbRandCount]
	}
	log.Debug("%d", cbCardBuffer)
	return
}

func (lg *sss_logic) RemoveCard(bRemoveCard []int, bRemoveCount int, bCardData []int, bCardCount int) bool {
	bDeleteCount := 0
	bTempCardData := make([]int, 13)
	if bCardCount > len(bTempCardData) {
		return false
	}
	//util.DeepCopy(&bTempCardData, &bCardData)
	copy(bTempCardData, bCardData)
	//置零扑克
	for i := 0; i < bRemoveCount; i++ {
		for j := 0; j < bCardCount; j++ {
			if bRemoveCard[i] == bTempCardData[j] {
				bDeleteCount++
				bTempCardData[j] = 0
				break
			}
		}
	}
	if bDeleteCount != bRemoveCount {
		return false
	}

	//清理扑克
	bCardPos := 0
	for i := 0; i < bCardCount; i++ {
		if bTempCardData[i] != 0 {
			bCardData[bCardPos] = bTempCardData[i]
			bCardPos++
		}
	}
	return true
}

//排列扑克
func (lg *sss_logic) SSSSortCardList(cardData []int) {
	cardCount := len(cardData)
	logicValue := make([]int, cardCount)
	for i := 0; i < cardCount; i++ {
		logicValue[i] = lg.GetCardLogicValue(cardData[i])
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

//逻辑数值
func (lg *sss_logic) GetCardLogicValue(cardData int) int {
	//扑克属性
	cardValue := lg.GetCardValue(cardData)
	//转换数值
	if cardValue == 1 {
		cardValue += 13
	}
	return cardValue
}

//获取数值
func (lg *sss_logic) GetCardValue(bCardData int) int { return bCardData & LOGIC_MASK_VALUE } //十六进制前面四位表示牌的数值
//获取花色
func (lg *sss_logic) GetCardColor(bCardData int) int { return (bCardData & LOGIC_MASK_COLOR) >> 4 } //十六进制后面四位表示牌的花色

//分析牌
func (lg *sss_logic) AnalyseCard(metaCardData []int) *TagAnalyseItem {
	cardCount := len(metaCardData)
	cardData := make([]int, cardCount)
	copy(cardData, metaCardData)

	lg.SSSSortCardList(cardData)

	//变量定义
	bSameCount := 1
	bCardValueTemp := 0
	bSameColorCount := 1
	bFirstCardIndex := 0 //记录下标

	bLogicValue := lg.GetCardLogicValue(cardData[0])
	bCardColor := lg.GetCardColor(cardData[0])

	analyseItem := &TagAnalyseItem{bOneFirst: make([]int, 13), bTwoFirst: make([]int, 13), bThreeFirst: make([]int, 13), bFourFirst: make([]int, 13)}
	//扑克分析
	for i := 1; i < cardCount; i++ {
		//获取扑克
		bCardValueTemp = lg.GetCardLogicValue(cardData[i])

		if bCardValueTemp == bLogicValue {
			bSameCount++

		}
		if bCardValueTemp != bLogicValue || i == (cardCount-1) {
			switch bSameCount {
			case 2:
				analyseItem.bTwoFirst[analyseItem.bTwoCount] = bFirstCardIndex
				analyseItem.bTwoCount++
			case 3:
				analyseItem.bThreeFirst[analyseItem.bThreeCount] = bFirstCardIndex
				analyseItem.bThreeCount++
			case 4:
				analyseItem.bFourFirst[analyseItem.bFourCount] = bFirstCardIndex
				analyseItem.bFourCount++
			}
		}

		if bCardValueTemp != bLogicValue {
			if bSameCount == 1 {
				if i != cardCount-1 {
					analyseItem.bOneFirst[analyseItem.bOneCount] = bFirstCardIndex
					analyseItem.bOneCount++
				} else {
					analyseItem.bOneFirst[analyseItem.bOneCount] = bFirstCardIndex
					analyseItem.bOneCount++
					analyseItem.bOneFirst[analyseItem.bOneCount] = i
					analyseItem.bOneCount++
				}
			} else {
				if i == cardCount-1 {
					analyseItem.bOneFirst[analyseItem.bOneCount] = i
					analyseItem.bOneCount++
				}
			}
			bSameCount = 1
			bLogicValue = bCardValueTemp
			bFirstCardIndex = i
		}
		if lg.GetCardColor(cardData[i]) != bCardColor {
			bSameColorCount = 1
		} else {
			bSameColorCount++
		}
	}

	if cardCount == bSameColorCount {
		analyseItem.bStraight = true
	} else {
		analyseItem.bStraight = false
	}
	return analyseItem

}

func (lg *sss_logic) GetCardType(metaCardData []int) int {

	cardCount := len(metaCardData)

	if cardCount != 3 && cardCount != 5 && cardCount != 13 {
		return CT_INVALID
	}

	cardData := make([]int, cardCount)
	copy(cardData, metaCardData)
	lg.SSSSortCardList(cardData)

	TagAnalyseItemArray := new(TagAnalyseItem)
	TagAnalyseItemArray = lg.AnalyseCard(cardData)

	//开始分析
	switch cardCount {
	case 3: //三条类型
		//单牌类型
		if TagAnalyseItemArray.bOneCount == 3 {
			return CT_SINGLE
		}
		//对带一张
		if TagAnalyseItemArray.bTwoCount == 1 && TagAnalyseItemArray.bOneCount == 1 {
			return CT_ONE_DOUBLE
		}
		//三张牌型
		if TagAnalyseItemArray.bThreeCount == 1 {
			return CT_THREE
		}
		//错误类型
		return CT_INVALID
	case 5: //五张牌型
		bFlushNoA := false
		bFlushFirstA := false
		bFlushBackA := false
		//A连在后
		if lg.GetCardLogicValue(cardData[0]) == 14 && lg.GetCardLogicValue(cardData[4]) == 10 {
			bFlushBackA = true
		} else {
			bFlushNoA = true
		}
		for i := 0; i < 4; i++ {
			if lg.GetCardLogicValue(cardData[i])-lg.GetCardLogicValue(cardData[i+1]) != 1 {
				bFlushBackA = false
				bFlushNoA = false
			}
		}
		//A连在前
		if false == bFlushBackA && false == bFlushNoA && 14 == lg.GetCardLogicValue(cardData[0]) {
			bFlushFirstA = true
			for i := 1; i < 4; i++ {
				if 1 != lg.GetCardLogicValue(cardData[i])-lg.GetCardLogicValue(cardData[i+1]) {
					bFlushFirstA = false
				}
			}
			if lg.GetCardLogicValue(cardData[4]) != 2 {
				bFlushFirstA = false
			}
		}
		//同花五牌
		if false == bFlushBackA && false == bFlushNoA && false == bFlushFirstA {
			if true == TagAnalyseItemArray.bStraight {
				return CT_FIVE_FLUSH
			}
		} else if true == bFlushNoA {
			//杂顺类型
			if false == TagAnalyseItemArray.bStraight {
				return CT_FIVE_MIXED_FLUSH_NO_A
			} else { //同花顺牌
				return CT_FIVE_STRAIGHT_FLUSH_NO_A
			}
		} else if true == bFlushFirstA {
			//杂顺类型
			if false == TagAnalyseItemArray.bStraight {
				return CT_FIVE_MIXED_FLUSH_FIRST_A
			} else { //同花顺牌
				return CT_FIVE_STRAIGHT_FLUSH_FIRST_A
			}
		} else if true == bFlushBackA {
			//杂顺类型
			if false == TagAnalyseItemArray.bStraight {
				return CT_FIVE_MIXED_FLUSH_BACK_A
			} else { //同花顺牌
				return CT_FIVE_STRAIGHT_FLUSH_BACK_A
			}
		}
		//四带单张
		if 1 == TagAnalyseItemArray.bFourCount && 1 == TagAnalyseItemArray.bOneCount {
			return CT_FIVE_FOUR_ONE
		}
		//三条一对
		if 1 == TagAnalyseItemArray.bThreeCount && 1 == TagAnalyseItemArray.bTwoCount {
			return CT_FIVE_THREE_DEOUBLE
		}
		//三条带单
		if 1 == TagAnalyseItemArray.bThreeCount && 2 == TagAnalyseItemArray.bOneCount {
			return CT_THREE
		}
		//两对牌型
		if 2 == TagAnalyseItemArray.bTwoCount && 1 == TagAnalyseItemArray.bOneCount {
			return CT_FIVE_TWO_DOUBLE
		}
		//只有一对
		if 1 == TagAnalyseItemArray.bTwoCount && 3 == TagAnalyseItemArray.bOneCount {
			return CT_ONE_DOUBLE
		}
		//单牌类型
		if 5 == TagAnalyseItemArray.bOneCount && false == TagAnalyseItemArray.bStraight {
			return CT_SINGLE
		}
		//错误类型
		return CT_INVALID

	case 13: //13张特殊牌型
		//至尊清龙
		if 13 == TagAnalyseItemArray.bOneCount && true == TagAnalyseItemArray.bStraight {
			return CT_THIRTEEN_FLUSH
		}
		//一条龙
		if 13 == TagAnalyseItemArray.bOneCount {
			return CT_THIRTEEN
		}

		//三同花顺
		btCardData := make([]int, 13)
		copy(btCardData, cardData)
		lg.SSSSortCardList(btCardData)
		RbtCardData := make([]int, 13)
		StraightFlush1 := false
		StraightFlush2 := false
		StraightFlush3 := false
		StraightFlush := 1
		Number := 0
		Count1 := 0
		Count2 := 0
		Count3 := 0
		FCardData := lg.GetCardLogicValue(btCardData[0])
		SColor := lg.GetCardColor(btCardData[0])
		RbtCardData[Number] = btCardData[0]
		Number++

		for i := 1; i < 13; i++ {
			if FCardData == lg.GetCardLogicValue(btCardData[i])+1 && SColor == lg.GetCardColor(btCardData[i]) {
				StraightFlush++
				FCardData = lg.GetCardLogicValue(btCardData[i])
				RbtCardData[Number] = btCardData[i]
				Number++
			}

			if FCardData != lg.GetCardLogicValue(btCardData[i])+1 && FCardData != lg.GetCardLogicValue(btCardData[i]) {
				if 3 == StraightFlush {
					StraightFlush1 = true
					Count1 = 3
					lg.RemoveCard(RbtCardData, 3, btCardData, 13)
					//ZeroMemory(RbtCardData, sizeof(RbtCardData))
					RbtCardData = make([]int, 13)
					break
				}

			}

			if 5 == StraightFlush {
				StraightFlush1 = true
				Count1 = 5
				lg.RemoveCard(RbtCardData, 5, btCardData, 13)
				//ZeroMemory(RbtCardData, sizeof(RbtCardData))
				RbtCardData = make([]int, 13)
				break
			}

		}
		if StraightFlush1 {
			StraightFlush = 1
			Number = 0
			FCardData = lg.GetCardLogicValue(btCardData[0])
			SColor = lg.GetCardColor(btCardData[0])
			RbtCardData[Number] = btCardData[0]
			Number++
			for i := 1; i < 13-Count1; i++ {
				if FCardData == lg.GetCardLogicValue(btCardData[i])+1 && SColor == lg.GetCardColor(btCardData[i]) {
					StraightFlush++
					FCardData = lg.GetCardLogicValue(btCardData[i])
					RbtCardData[Number] = btCardData[i]
					Number++
				}
				if FCardData != lg.GetCardLogicValue(btCardData[i])+1 && FCardData != lg.GetCardLogicValue(btCardData[i]) {
					if 3 == StraightFlush {
						StraightFlush2 = true
						Count2 = 3
						lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
						//ZeroMemory(RbtCardData,sizeof(RbtCardData))
						RbtCardData = make([]int, 13)
					}
					break
				}
				if 5 == StraightFlush {
					StraightFlush2 = true
					Count2 = 5
					lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1)
					//ZeroMemory(RbtCardData,sizeof(RbtCardData));
					RbtCardData = make([]int, 13)
					break
				}
			}
		}
		if StraightFlush2 {
			StraightFlush = 1
			Number = 0
			//btCardData = btCardData[Count2:]
			FCardData = lg.GetCardLogicValue(btCardData[0])
			SColor = lg.GetCardColor(btCardData[0])
			RbtCardData[Number] = btCardData[0]
			Number++
			for i := 1; i < 13-Count1-Count2; i++ {
				if FCardData == lg.GetCardLogicValue(btCardData[i])+1 && SColor == lg.GetCardColor(btCardData[i]) {
					StraightFlush++
					FCardData = lg.GetCardLogicValue(btCardData[i])
					RbtCardData[Number] = btCardData[i]
					Number++
				}
				if FCardData != lg.GetCardLogicValue(btCardData[i])+1 && FCardData != lg.GetCardLogicValue(btCardData[i]) {
					if 3 == StraightFlush {
						StraightFlush3 = true
						Count3 = 3
						lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
						//ZeroMemory(RbtCardData,sizeof(RbtCardData))
						RbtCardData = make([]int, 13)
					}
					break
				}
				if 5 == StraightFlush {
					StraightFlush3 = true
					Count3 = 5
					lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1-Count2)
					//ZeroMemory(RbtCardData,sizeof(RbtCardData));
					RbtCardData = make([]int, 13)
					break
				}
			}
		}

		if StraightFlush1 && StraightFlush2 && StraightFlush3 && Count1+Count2+Count3 == 13 {
			return CT_THREE_STRAIGHTFLUSH
		}

		//三分天下
		if 3 == TagAnalyseItemArray.bFourCount {
			return CT_THREE_BOMB
		}

		//四套三条
		if 4 == TagAnalyseItemArray.bThreeCount {
			return CT_FOUR_THREESAME
		}

		//六对半
		if (6 == TagAnalyseItemArray.bTwoCount) || (4 == TagAnalyseItemArray.bTwoCount && 1 == TagAnalyseItemArray.bFourCount) ||
			(2 == TagAnalyseItemArray.bTwoCount && 2 == TagAnalyseItemArray.bFourCount) || (3 == TagAnalyseItemArray.bFourCount) {
			return CT_SIXPAIR
		}

		//三顺子
		nCount := 0
		for nCount < 4 {
			nCount++
			Straight1 := false
			Straight2 := false
			Straight3 := false
			Straight := 1
			Count1 = 0
			Count2 = 0
			Count3 = 0
			Number = 0
			//RbtCardData = RbtCardData[:0]
			RbtCardData = make([]int, 13)
			//util.DeepCopy(&btCardData, &cardData)
			copy(btCardData, cardData)
			lg.SSSSortCardList(btCardData)
			RbtCardData[Number] = btCardData[0]
			Number++
			FCardData = lg.GetCardLogicValue(btCardData[0])
			for i := 1; i < 13; i++ {
				if FCardData == lg.GetCardLogicValue(btCardData[i])+1 ||
					(FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 5) ||
					(FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 3) {
					Straight++
					RbtCardData[Number] = btCardData[i]
					Number++
					FCardData = lg.GetCardLogicValue(btCardData[i])

				} else if FCardData != lg.GetCardLogicValue(btCardData[i]) {
					if 3 == Straight {
						Straight1 = true
						Count1 = 3
						//util.DeepCopy(&btSpecialCard[10], RbtCardData)
						//copy(btSpecialCard[10:], RbtCardData[:Count1])
						lg.RemoveCard(RbtCardData, 3, btCardData, 13)
						RbtCardData[Number] = btCardData[0]
						break
					}
					Straight = 1
					Number = 0
					FCardData = lg.GetCardLogicValue(btCardData[i])
					RbtCardData[Number] = btCardData[i]
					Number++

				}
				if nCount == 0 || nCount == 1 {
					if i == 12 && 3 == Straight {
						Straight1 = true
						Count1 = 3
						//util.DeepCopy(&btSpecialCard[10], RbtCardData)
						//copy(btSpecialCard[10:], RbtCardData[:Count1])
						lg.RemoveCard(RbtCardData, 3, btCardData, 13)
						//RbtCardData = RbtCardData[:0]
						RbtCardData = make([]int, 13)
						break

					}
				} else if nCount == 2 || nCount == 3 {
					if 3 == Straight {

						Straight1 = true
						Count1 = 3
						//util.DeepCopy(&btSpecialCard[10], RbtCardData)
						//copy(btSpecialCard[10:], RbtCardData[:Count1])
						lg.RemoveCard(RbtCardData, 3, btCardData, 13)
						//RbtCardData = RbtCardData[:0]
						RbtCardData = make([]int, 13)
						break
					}
				}
				if 5 == Straight {
					Straight1 = true
					Count1 = 5
					//util.DeepCopy(&btSpecialCard[5], RbtCardData)
					//copy(btSpecialCard[5:], RbtCardData[:Count1])
					lg.RemoveCard(RbtCardData, 5, btCardData, 13)
					//RbtCardData = RbtCardData[:0]
					RbtCardData = make([]int, 13)
					break

				}
			}
			if Straight1 {
				Straight = 1
				Number = 0
				lg.SSSSortCardList(btCardData)
				RbtCardData[Number] = btCardData[0]
				Number++
				FCardData = lg.GetCardLogicValue(btCardData[0])
				for i := 1; i < 13-Count1; i++ {
					if FCardData == lg.GetCardLogicValue(btCardData[i])+1 || (FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 5) || (FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 3) {
						Straight++
						RbtCardData[Number] = btCardData[i]
						Number++
						FCardData = lg.GetCardLogicValue(btCardData[i])

					} else if FCardData != lg.GetCardLogicValue(btCardData[i]) {
						if 3 == Straight && Count1 != 3 {
							Straight2 = true
							Count2 = 3
							//util.DeepCopy(&btSpecialCard[10], RbtCardData)
							//copy(btSpecialCard[10:], RbtCardData[:Count2])
							lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
							//RbtCardData = RbtCardData[:0]
							RbtCardData = make([]int, 13)
							break
						}
						Straight = 1
						Number = 0
						FCardData = lg.GetCardLogicValue(btCardData[i])
						RbtCardData[Number] = btCardData[i]
						Number++
					}
					if nCount == 0 || nCount == 2 {
						if i == 13-Count1-1 && 3 == Straight && Count1 != 3 {
							Straight2 = true
							Count2 = 3
							//util.DeepCopy(&btSpecialCard[10], RbtCardData)
							//copy(btSpecialCard[10:], RbtCardData[:Count2])
							lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
							//RbtCardData = RbtCardData[:0]
							RbtCardData = make([]int, 13)
							break

						}
					} else if nCount == 1 || nCount == 3 {
						if 3 == Straight && Count1 != 3 {
							Straight2 = true
							Count2 = 3
							//util.DeepCopy(&btSpecialCard[10], RbtCardData)
							//copy(btSpecialCard[10:], RbtCardData[:Count2])
							lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
							//RbtCardData = RbtCardData[:0]
							RbtCardData = make([]int, 13)
							break

						}
					}
					if 5 == Straight {
						Straight2 = true
						Count2 = 5
						if Count1 == 5 {
							//util.DeepCopy(&btSpecialCard[0], RbtCardData)
							//copy(btSpecialCard, RbtCardData[:Count2])
						} else {
							//util.DeepCopy(&btSpecialCard[5], RbtCardData)
							//copy(btSpecialCard[5:], RbtCardData[:Count2])
						}

						lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1)
						//RbtCardData = RbtCardData[:0]
						RbtCardData = make([]int, 13)
						break
					}
				}
			}
			if Straight2 {
				Straight = 1
				Number = 0
				lg.SSSSortCardList(btCardData)
				RbtCardData[Number] = btCardData[0]
				Number++
				FCardData = lg.GetCardLogicValue(btCardData[0])
				for i := 1; i < 13-Count1-Count2; i++ {
					if FCardData == lg.GetCardLogicValue(btCardData[i])+1 || (FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 3) || (FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 5) {
						Straight++
						RbtCardData[Number] = btCardData[i]
						Number++
						FCardData = lg.GetCardLogicValue(btCardData[i])
					} else if FCardData != lg.GetCardLogicValue(btCardData[i]) {
						if 3 == Straight && Count1 != 3 && Count2 != 3 {
							Straight3 = true
							Count3 = 3
							//util.DeepCopy(&btSpecialCard[10], RbtCardData)
							//copy(btSpecialCard[10:], RbtCardData[:Count3])
							lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
							//RbtCardData = RbtCardData[:0]
							RbtCardData = make([]int, 13)
							break
						}
						Straight = 1
						Number = 0
						FCardData = lg.GetCardLogicValue(btCardData[i])
						RbtCardData[Number] = btCardData[i]
						Number++
					}
					if i == 13-Count1-Count2-1 && 3 == Straight && Count1 != 3 && Count2 != 3 {
						Straight3 = true
						Count3 = 3
						//util.DeepCopy(&btSpecialCard[10], RbtCardData)
						//copy(btSpecialCard[10:], RbtCardData[:Count3])
						lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
						//RbtCardData = RbtCardData[:0]
						RbtCardData = make([]int, 13)
						break
					}
					if 5 == Straight {
						Straight3 = true
						Count3 = 5
						//util.DeepCopy(&btSpecialCard[0], RbtCardData)
						//copy(btSpecialCard[:], RbtCardData[:Count3])
						lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1-Count2)
						//RbtCardData = RbtCardData[:0]
						RbtCardData = make([]int, 13)
						break
					}
				}
			}
			if Straight1 && Straight2 && Straight3 && Count1+Count2+Count3 == 13 {
				return CT_THREE_STRAIGHT
			}

		}

		//三同花
		Flush1 := false
		Flush2 := false
		Flush3 := false
		Flush := 1
		Count1 = 0
		Count2 = 0
		Count3 = 0
		Number = 0
		RbtCardData = make([]int, 13)
		//util.DeepCopy(&btCardData, &cardData)
		btCardData = make([]int, 13)
		copy(btCardData, cardData)
		RbtCardData[Number] = btCardData[0]
		Number++
		SColor = lg.GetCardColor(btCardData[0])
		for i := 1; i < 13; i++ {
			if SColor == lg.GetCardColor(btCardData[i]) {
				Flush++
				RbtCardData[Number] = btCardData[i]
				Number++
			}
			if 3 == Flush && i == 12 {
				Flush1 = true
				Count1 = 3
				//util.DeepCopy(&btSpecialCard[10], RbtCardData)
				//copy(btSpecialCard[10:], RbtCardData[:Count1])
				lg.RemoveCard(RbtCardData, 3, btCardData, 13)
				//RbtCardData = RbtCardData[:0]
				RbtCardData = make([]int, 13)
				break
			}
			if 5 == Flush {
				Flush1 = true
				Count1 = 5
				//util.DeepCopy(&btSpecialCard[5], RbtCardData)
				//copy(btSpecialCard[5:], RbtCardData[:Count1])
				lg.RemoveCard(RbtCardData, 5, btCardData, 13)
				//RbtCardData = RbtCardData[:0]
				RbtCardData = make([]int, 13)
				break
			}
		}
		if Flush1 {
			Flush = 1
			Number = 0

			RbtCardData[Number] = btCardData[0]
			Number++
			SColor = lg.GetCardColor(btCardData[0])
			for i := 1; i < 13-Count1; i++ {
				if SColor == lg.GetCardColor(btCardData[i]) {
					Flush++
					RbtCardData[Number] = btCardData[i]
					Number++
				}
				if 3 == Flush && i == 13-Count1-1 && Count1 != 3 {
					Flush2 = true
					Count2 = 3
					//util.DeepCopy(&btSpecialCard[10], RbtCardData)
					//copy(btSpecialCard[10:], RbtCardData[:Count2])
					lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
					//RbtCardData = RbtCardData[:0]
					RbtCardData = make([]int, 13)
					break
				}
				if 5 == Flush {
					Flush2 = true
					Count2 = 5
					if Count1 == 5 {
						//util.DeepCopy(&btSpecialCard[0], RbtCardData)
						//copy(btSpecialCard, RbtCardData[:Count1])
					} else if Count1 == 3 {
						//util.DeepCopy(&btSpecialCard[5], RbtCardData)
						//copy(btSpecialCard, RbtCardData[:Count1])
					}

					lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1)
					//RbtCardData = RbtCardData[:0]
					RbtCardData = make([]int, 13)
					break
				}
			}
		}
		if Flush2 {
			Flush = 1
			Number = 0
			RbtCardData[Number] = btCardData[0]
			Number++
			SColor = lg.GetCardColor(btCardData[0])
			for i := 1; i < 13-Count1-Count2; i++ {
				if SColor == lg.GetCardColor(btCardData[i]) {
					Flush++
					RbtCardData[Number] = btCardData[i]
					Number++
				}
				if 3 == Flush && i == 13-Count1-Count2-1 && Count1 != 3 && Count2 != 3 {
					Flush3 = true
					Count3 = 3
					//util.DeepCopy(&btSpecialCard[10], RbtCardData)
					//copy(btSpecialCard[10:], RbtCardData[:Count3])
					lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
					//RbtCardData = RbtCardData[:0]
					RbtCardData = make([]int, 13)
					break
				}
				if 5 == Flush {
					Flush3 = true
					Count3 = 5
					//util.DeepCopy(&btSpecialCard[0], RbtCardData)
					//copy(btSpecialCard, RbtCardData[:Count3])
					lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1-Count2)
					//RbtCardData = RbtCardData[:0]
					RbtCardData = make([]int, 13)
					break
				}
			}
		}
		if Flush1 && Flush2 && Flush3 && Count1+Count2+Count3 == 13 {
			nCount := 0
			for nCount < 4 {
				nCount++
				Straight1 := false
				Straight2 := false
				Straight3 := false
				Straight := 1
				Count1 = 0
				Count2 = 0
				Count3 = 0
				Number = 0
				//RbtCardData = RbtCardData[:0]
				RbtCardData = make([]int, 13)
				//util.DeepCopy(btCardData, cardData)
				copy(btCardData, cardData)
				lg.SSSSortCardList(btCardData)
				RbtCardData[Number] = btCardData[0]
				Number++
				FCardData = lg.GetCardLogicValue(btCardData[0])
				for i := 1; i < 13; i++ {
					if FCardData == lg.GetCardLogicValue(btCardData[i])+1 ||
						(FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 5) ||
						(FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 3) {
						Straight++
						RbtCardData[Number] = btCardData[i]
						Number++
						FCardData = lg.GetCardLogicValue(btCardData[i])

					} else if FCardData != lg.GetCardLogicValue(btCardData[i]) {
						if 3 == Straight {
							Straight1 = true
							Count1 = 3
							//util.DeepCopy(&btSpecialCard[10], RbtCardData)
							//copy(btSpecialCard[10:], RbtCardData[:Count1])
							lg.RemoveCard(RbtCardData, 3, btCardData, 13)
							RbtCardData[Number] = btCardData[0]
							break
						}
						Straight = 1
						Number = 0
						FCardData = lg.GetCardLogicValue(btCardData[i])
						RbtCardData[Number] = btCardData[i]
						Number++

					}
					if nCount == 0 || nCount == 1 {
						if i == 12 && 3 == Straight {
							Straight1 = true
							Count1 = 3
							//util.DeepCopy(&btSpecialCard[10], RbtCardData)
							//copy(btSpecialCard[10:], RbtCardData[:Count1])
							lg.RemoveCard(RbtCardData, 3, btCardData, 13)
							//RbtCardData = RbtCardData[:0]
							RbtCardData = make([]int, 13)
							break

						}
					} else if nCount == 2 || nCount == 3 {
						if 3 == Straight {

							Straight1 = true
							Count1 = 3
							//util.DeepCopy(&btSpecialCard[10], RbtCardData)
							//copy(btSpecialCard[10:], RbtCardData[:Count1])
							lg.RemoveCard(RbtCardData, 3, btCardData, 13)
							//RbtCardData = RbtCardData[:0]
							RbtCardData = make([]int, 13)
							break
						}
					}
					if 5 == Straight {
						Straight1 = true
						Count1 = 5
						//util.DeepCopy(&btSpecialCard[5], RbtCardData)
						//copy(btSpecialCard[5:], RbtCardData[:Count1])
						lg.RemoveCard(RbtCardData, 5, btCardData, 13)
						//RbtCardData = RbtCardData[:0]
						RbtCardData = make([]int, 13)
						break

					}
				}
				if Straight1 {
					Straight = 1
					Number = 0
					lg.SSSSortCardList(btCardData)
					RbtCardData[Number] = btCardData[0]
					Number++
					FCardData = lg.GetCardLogicValue(btCardData[0])
					for i := 1; i < 13-Count1; i++ {
						if FCardData == lg.GetCardLogicValue(btCardData[i])+1 ||
							(FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 5) ||
							(FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 3) {
							Straight++
							RbtCardData[Number] = btCardData[i]
							Number++
							FCardData = lg.GetCardLogicValue(btCardData[i])

						} else if FCardData != lg.GetCardLogicValue(btCardData[i]) {
							if 3 == Straight && Count1 != 3 {
								Straight2 = true
								Count2 = 3
								//util.DeepCopy(&btSpecialCard[10], RbtCardData)
								//copy(btSpecialCard[10:], RbtCardData[:Count2])
								lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
								//RbtCardData = RbtCardData[:0]
								RbtCardData = make([]int, 13)
								break
							}
							Straight = 1
							Number = 0
							FCardData = lg.GetCardLogicValue(btCardData[i])
							RbtCardData[Number] = btCardData[i]
							Number++
						}
						if nCount == 0 || nCount == 2 {
							if i == 13-Count1-1 && 3 == Straight && Count1 != 3 {
								Straight2 = true
								Count2 = 3
								//util.DeepCopy(&btSpecialCard[10], RbtCardData)
								//copy(btSpecialCard[10:], RbtCardData[:Count2])
								lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
								//RbtCardData = RbtCardData[:0]
								RbtCardData = make([]int, 13)
								break

							}
						} else if nCount == 1 || nCount == 3 {
							if 3 == Straight && Count1 != 3 {
								Straight2 = true
								Count2 = 3
								//util.DeepCopy(&btSpecialCard[10], RbtCardData)
								//copy(btSpecialCard[10:], RbtCardData[:Count2])
								lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
								//RbtCardData = RbtCardData[:0]
								RbtCardData = make([]int, 13)
								break

							}
						}
						if 5 == Straight {
							Straight2 = true
							Count2 = 5
							if Count1 == 5 {
								//util.DeepCopy(&btSpecialCard[0], RbtCardData)
								//copy(btSpecialCard, RbtCardData[:Count1])
							} else {
								//util.DeepCopy(&btSpecialCard[5], RbtCardData)
								//copy(btSpecialCard, RbtCardData[:Count1])
							}

							lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1)
							//RbtCardData = RbtCardData[:0]
							RbtCardData = make([]int, 13)
							break
						}
					}
				}
				if Straight2 {
					Straight = 1
					Number = 0
					lg.SSSSortCardList(btCardData)
					RbtCardData[Number] = btCardData[0]
					Number++
					FCardData = lg.GetCardLogicValue(btCardData[0])
					for i := 1; i < 13-Count1-Count2; i++ {
						if FCardData == lg.GetCardLogicValue(btCardData[i])+1 ||
							(FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 3) ||
							(FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 5) {
							Straight++
							RbtCardData[Number] = btCardData[i]
							Number++
							FCardData = lg.GetCardLogicValue(btCardData[i])
						} else if FCardData != lg.GetCardLogicValue(btCardData[i]) {
							if 3 == Straight && Count1 != 3 && Count2 != 3 {
								Straight3 = true
								Count3 = 3
								//util.DeepCopy(&btSpecialCard[10], RbtCardData)
								//copy(btSpecialCard[10:], RbtCardData[:Count3])
								lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
								//RbtCardData = RbtCardData[:0]
								RbtCardData = make([]int, 13)
								break
							}
							Straight = 1
							Number = 0
							FCardData = lg.GetCardLogicValue(btCardData[i])
							RbtCardData[Number] = btCardData[i]
							Number++
						}
						if i == 13-Count1-Count2-1 && 3 == Straight && Count1 != 3 && Count2 != 3 {
							Straight3 = true
							Count3 = 3
							//util.DeepCopy(&btSpecialCard[10], RbtCardData)
							//copy(btSpecialCard[10:], RbtCardData[:Count3])
							lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
							//RbtCardData = RbtCardData[:0]
							RbtCardData = make([]int, 13)
							break
						}
						if 5 == Straight {
							Straight3 = true
							Count3 = 5
							//util.DeepCopy(&btSpecialCard[0], RbtCardData)
							//copy(btSpecialCard, RbtCardData[:Count3])
							lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1-Count2)
							//RbtCardData = RbtCardData[:0]
							RbtCardData = make([]int, 13)
							break
						}
					}
				}
				if Straight1 && Straight2 && Straight3 && Count1+Count2+Count3 == 13 {
					return CT_THREE_STRAIGHTFLUSH
				}
			}
			return CT_THREE_FLUSH
		}

	}

	return CT_INVALID
}

func (lg *sss_logic) SSSCompareCard(bInFirstList []int, bInNextList []int) int {
	bFirstCount := len(bInFirstList)
	bNextCount := len(bInNextList)

	if bFirstCount != bNextCount {
		//todo验证
		return 0
	}

	FirstAnalyseData := new(TagAnalyseItem)
	NextAnalyseData := new(TagAnalyseItem)

	bFirstList := make([]int, bFirstCount)
	bNextList := make([]int, bNextCount)

	copy(bFirstList, bInFirstList)
	copy(bNextList, bInNextList)

	// lg.SSSSortCardList(bFirstList)
	// lg.SSSSortCardList(bNextList)

	FirstAnalyseData = lg.AnalyseCard(bFirstList)
	NextAnalyseData = lg.AnalyseCard(bNextList)

	bNextType := lg.GetCardType(bNextList)
	bFirstType := lg.GetCardType(bFirstList)

	if CT_INVALID == bFirstType || CT_INVALID == bNextType {
		return -1
	}
	//三张牌型
	if 3 == bFirstCount {
		//开始对比
		if bNextType == bFirstType {
			switch bFirstType {
			case CT_SINGLE: //单牌类型
				for i := 0; i < 3; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				return 0
			case CT_ONE_DOUBLE: //对带一张
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return -1
				}
				return 0
			case CT_THREE: //三张牌型
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return -1
				}
				return 0
			}
		} else {
			if bNextType > bFirstType {
				return 1
			} else {
				return -1
			}
		}
	}
	//五张牌型
	if 5 == bFirstCount {
		if bNextType == bFirstType {
			switch bFirstType {
			case CT_SINGLE: //单牌类型
				for i := 0; i < 5; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				return 0
			case CT_ONE_DOUBLE: //对带三张
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return -1
				}
				for i := 0; i < 3; i++ {
					if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[i]]) != lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[i]]) {
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[i]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[i]]) {
							return 1
						} else {
							return -1
						}
					}
				}
				return 0
			case CT_FIVE_TWO_DOUBLE: //两对牌型
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return -1
				}
				return 0
			case CT_THREE: //三张牌型
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[1]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[1]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[1]]) {
					return -1
				}
				return 0
			case CT_FIVE_MIXED_FLUSH_NO_A, CT_FIVE_MIXED_FLUSH_FIRST_A, CT_FIVE_MIXED_FLUSH_BACK_A: //没A杂顺 A在前顺子 A在后顺子
				if lg.GetCardLogicValue(bNextList[0]) > lg.GetCardLogicValue(bFirstList[0]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[0]) < lg.GetCardLogicValue(bFirstList[0]) {
					return -1
				}
				return 0
			case CT_FIVE_FLUSH: //同花五牌
				for i := 0; i < 5; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				if lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) {
					return 1
				}
				if lg.GetCardColor(bNextList[0]) < lg.GetCardColor(bFirstList[0]) {
					return -1
				}
				return 0
			case CT_FIVE_THREE_DEOUBLE: //三条一对
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return -1
				}
				return 0
			case CT_FIVE_FOUR_ONE: //四带一张

				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return -1
				}
				return 0
			case CT_FIVE_STRAIGHT_FLUSH_NO_A, CT_FIVE_STRAIGHT_FLUSH_FIRST_A, CT_FIVE_STRAIGHT_FLUSH_BACK_A: //没A同花顺 A在前同花顺 A在后同花顺
				for i := 0; i < 5; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				if lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) {
					return 1
				}
				if lg.GetCardColor(bNextList[0]) < lg.GetCardColor(bFirstList[0]) {
					return -1
				}
				return 0

			default:
				return 0
			}
		} else {
			if bNextType > bFirstType {
				return 1
			} else {
				return -1
			}
		}
	}
	//13张牌型
	if 13 == bFirstCount {
		if bNextType == bFirstType {
			switch bFirstType {
			case CT_THIRTEEN_FLUSH:
				if lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) {
					return 1
				}
				if lg.GetCardColor(bNextList[0]) < lg.GetCardColor(bFirstList[0]) {
					return -1
				}
				return 0
			case CT_TWELVE_KING:
				for i := 0; i < 13; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				return 0
			case CT_THREE_STRAIGHTFLUSH:
				for i := 0; i < 13; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				return 0
				//todo 比花色？
			case CT_THREE_BOMB:
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[1]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[1]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[2]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[2]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return -1
				}
				return 0
			case CT_ALL_BIG, CT_ALL_SMALL, CT_SAME_COLOR:
				for i := 0; i < 13; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				return 0
			case CT_FOUR_THREESAME:
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[1]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[1]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[1]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[2]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[2]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[2]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[3]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[3]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[3]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[3]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return -1
				}
				return 0
			case CT_FIVEPAIR_THREE:
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return -1
				}
				return 0
			case CT_SIXPAIR:
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[5]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[5]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[5]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[5]]) {
					return -1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
					return -1
				}
				return 0
			case CT_THREE_FLUSH:
				for i := 0; i < 13; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				return 0
			case CT_THREE_STRAIGHT:
				for i := 0; i < 13; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				return 0
			}
		} else {
			if bNextType > bFirstType {
				return 1
			} else {
				return -1
			}
		}
	}

	return 0
}

func (*sss_logic) SortCardList(cardData []int, cardCount int) {}
