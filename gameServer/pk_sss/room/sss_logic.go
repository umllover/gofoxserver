package room

import (
	"mj/gameServer/common/pk/pk_base"

	//dbg "github.com/funny/debug"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

// 十三水逻辑
const (
	CT_INVALID        = iota //错误类型
	CT_SINGLE                //散牌
	CT_ONEPAIR               //对子
	CT_TWOPAIR               //两对
	CT_THREESAME             //三条
	CT_STRAIGHT              //顺子
	CT_FLUSH                 //同花
	CT_GOURD                 //葫芦
	CT_FOURSAME              //铁支
	CT_STRAIGHT_FLUSH        //同花顺
	CT_FIVE_SAME             //五同
)

//特殊牌型
const (
	CT_THIRTEEN_FLUSH       = 26 //至尊清龙
	CT_THIRTEEN             = 25 //一条龙
	CT_TWELVE_KING          = 24 //十二皇族
	CT_THREE_STRAIGHT_FLUSH = 23 //三同花顺
	CT_THREE_BOMB           = 22 //三分天下
	CT_ALL_BIG              = 21 //全大
	CT_ALL_SMALL            = 20 //全小
	CT_SAME_COLOR           = 19 //凑一色
	CT_FOUR_THREESAME       = 18 //四套三条
	CT_FIVEPAIR_THREE       = 17 //五对冲三
	CT_SIXPAIR              = 16 //六对半
	CT_THREE_FLUSH          = 15 //三同花
	CT_THREE_STRAIGHT       = 14 //三顺子
)

//数值掩码
const (
	LOGIC_MASK_COLOR = 0xF0 //花色掩码
	LOGIC_MASK_VALUE = 0x0F //数值掩码
)

type TagAnalyseItem struct {
	cardData    []int //排序后的牌数据
	laiZi       []int //癞子
	bOneCount   int   //单张数目
	bTwoCount   int   //两张数目
	bThreeCount int   //三张数目
	bFourCount  int   //四张数目
	bFiveCount  int   //五张数目
	bOneFirst   []int //单牌位置
	bTwoFirst   []int //对牌位置
	bThreeFirst []int //三条位置
	bFourFirst  []int //四张位置
	bFiveFirst  []int //五张位置
	bstraight   bool  //是否顺子
	bflush      bool  //是否同花
}

// type tagAnalyseType struct {
// }

// //分析结构
// type tagAnalyseData struct {
// 	bOneCount   int   //单张数目
// 	bTwoCount   int   //两张数目
// 	bThreeCount int   //三张数目
// 	bFourCount  int   //四张数目
// 	bFiveCount  int   //五张数目
// 	bOneFirst   []int //单牌位置
// 	bTwoFirst   []int //对牌位置
// 	bThreeFirst []int //三条位置
// 	bFourFirst  []int //四张位置
// 	bflush   bool  //是否顺子
// }

func NewSssZLogic(ConfigIdx int) *sss_logic {
	l := new(sss_logic)
	//l.BtCardSpecialData = make([]int, 13)
	l.BaseLogic = pk_base.NewBaseLogic(ConfigIdx)
	l.LaiZhiSubstitute = -1
	return l
}

type sss_logic struct {
	*pk_base.BaseLogic
	//BtCardSpecialData []int
	//UniversalCards []int //万能牌
	LaiZhiSubstitute int
}

func (lg *sss_logic) RandCardList(cbCardBuffer, OriDataArray []int) {

	//混乱准备
	cbBufferCount := len(cbCardBuffer)
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
	if cardData < 0 {
		return cardData
	}
	//扑克属性
	cardValue := lg.GetCardValue(cardData)
	if cardValue == 14 || cardValue == 15 {
		return 15
	}
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
	cardData := []int{}
	laiZi := []int{}

	for _, v := range metaCardData {
		if v != lg.LaiZhiSubstitute {
			cardData = append(cardData, v)
		} else {
			laiZi = append(laiZi, v)
		}
	}

	cardCount := len(cardData)

	lg.SSSSortCardList(cardData)

	//变量定义
	bSameCount := 1
	bCardValueTemp := 0
	bSameColorCount := 1
	bFirstCardIndex := 0 //记录下标

	bLogicValue := lg.GetCardLogicValue(cardData[0])
	bCardColor := lg.GetCardColor(cardData[0])

	analyseItem := &TagAnalyseItem{bOneFirst: make([]int, 13), bTwoFirst: make([]int, 13), bThreeFirst: make([]int, 13), bFourFirst: make([]int, 13), bFiveFirst: make([]int, 13)}
	analyseItem.cardData = cardData
	analyseItem.laiZi = laiZi
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
			case 5:
				analyseItem.bFiveFirst[analyseItem.bFiveCount] = bFirstCardIndex
				analyseItem.bFiveCount++
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
		analyseItem.bflush = true
	} else {
		analyseItem.bflush = false
	}
	return analyseItem

}

func (lg *sss_logic) GetCardType(metaCardData []int) int {
	return 0
}

func (lg *sss_logic) SSSGetCardType(metaCardData []int) (int, *TagAnalyseItem) {
	metaCount := len(metaCardData)
	if metaCount != 3 && metaCount != 5 && metaCount != 13 {
		return CT_INVALID, new(TagAnalyseItem)
	}

	cardData := make([]int, metaCount)
	copy(cardData, metaCardData)

	lg.SSSSortCardList(cardData)

	TagAnalyseItemArray := new(TagAnalyseItem)
	TagAnalyseItemArray = lg.AnalyseCard(cardData)

	//开始分析
	switch metaCount {
	case 3: //三条类型
		switch len(TagAnalyseItemArray.laiZi) {
		// case 3:
		// 	cardData = []int{0x31, 0x31, 0x31}
		// 	return lg.SSSGetCardType(cardData)
		case 2:
			cardData = []int{cardData[0], cardData[0], cardData[0]}
			return lg.SSSGetCardType(cardData)
		case 1:
			if TagAnalyseItemArray.bTwoCount == 1 {
				cardData = []int{cardData[0], cardData[0], cardData[0]}
				return lg.SSSGetCardType(cardData)
			}
			if TagAnalyseItemArray.bOneCount == 2 {
				cardData = []int{cardData[0], cardData[0], cardData[1]}
				return lg.SSSGetCardType(cardData)
			}
		}

		//散牌
		if TagAnalyseItemArray.bOneCount == 3 {
			return CT_SINGLE, TagAnalyseItemArray
		}
		//对子
		if TagAnalyseItemArray.bTwoCount == 1 && TagAnalyseItemArray.bOneCount == 1 {
			return CT_ONEPAIR, TagAnalyseItemArray
		}
		//三条
		if TagAnalyseItemArray.bThreeCount == 1 {
			return CT_THREESAME, TagAnalyseItemArray
		}
		//错误类型
		return CT_INVALID, TagAnalyseItemArray
	case 5: //五张牌型
		switch len(TagAnalyseItemArray.laiZi) {
		// case 5: //最大五同
		// 	cardData = []int{0x31, 0x31, 0x31, 0x31, 0x31}
		// 	return lg.SSSGetCardType(cardData)
		case 4: //五同
			cardData = []int{cardData[0], cardData[0], cardData[0], cardData[0], cardData[0]}
			return lg.SSSGetCardType(cardData)
		case 3:
			//五同
			if TagAnalyseItemArray.bTwoCount == 1 {
				cardData = []int{cardData[0], cardData[0], cardData[0], cardData[0], cardData[0]}
				return lg.SSSGetCardType(cardData)
			}
			//同花顺
			if TagAnalyseItemArray.bflush {
				a := lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[0]]) - lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[1]])
				if a > 0 && a < 5 {
					if lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[1]]) > 9 {
						color := lg.GetCardColor(cardData[TagAnalyseItemArray.bOneFirst[1]])
						tempCardData := color<<4 + 10
						cardData = []int{tempCardData - 9, tempCardData, tempCardData + 1, tempCardData + 2, tempCardData + 3}
					} else {
						cardData = []int{cardData[1], cardData[1] + 1, cardData[1] + 2, cardData[1] + 3, cardData[1] + 4}
					}
					return lg.SSSGetCardType(cardData)
				}
				a = lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[1]]) - lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[0]])
				if a > 0 && a < 5 {
					cardData = []int{cardData[1], cardData[1] - 1, cardData[1] - 2, cardData[1] - 3, cardData[1] - 4}
					return lg.SSSGetCardType(cardData)
				}
			}
			//铁支
			cardData = []int{cardData[0], cardData[0], cardData[0], cardData[0], cardData[1]}
			return lg.SSSGetCardType(cardData)
		case 2:
			//五同
			if TagAnalyseItemArray.bThreeCount == 1 {
				cardData = []int{cardData[0], cardData[0], cardData[0], cardData[0], cardData[0]}
				return lg.SSSGetCardType(cardData)
			}
			if TagAnalyseItemArray.bOneCount == 3 {
				if TagAnalyseItemArray.bflush {
					//同花顺
					a := lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[0]]) - lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[2]])
					if a > 0 && a < 5 {
						if lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[2]]) > 9 {
							color := lg.GetCardColor(cardData[TagAnalyseItemArray.bOneFirst[2]])
							tempCardData := color<<4 + 10
							cardData = []int{tempCardData - 9, tempCardData, tempCardData + 1, tempCardData + 2, tempCardData + 3}
						} else {
							cardData = []int{cardData[2], cardData[2] + 1, cardData[2] + 2, cardData[2] + 3, cardData[2] + 4}
						}
						return lg.SSSGetCardType(cardData)
					}
					a = lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[2]]) - lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[0]])
					if a > 0 && a < 5 {
						cardData = []int{cardData[2], cardData[2] - 1, cardData[2] - 2, cardData[2] - 3, cardData[2] - 4}
						return lg.SSSGetCardType(cardData)
					}
					//同花
					cardData = []int{cardData[0], cardData[0], cardData[0], cardData[1], cardData[2]}
					return lg.SSSGetCardType(cardData)
				}

				//顺子
				a := lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[0]]) - lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[2]])
				if a > 0 && a < 5 {
					if lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[2]]) > 9 {
						cardData = []int{0<<4 + 1, 1<<4 + 10, 2<<4 + 11, 3<<4 + 12, 0<<4 + 13}
					} else {
						tv := lg.GetCardLogicValue(cardData[2])
						cardData = []int{cardData[2], 0<<4 + tv + 1, 1<<4 + tv + 2, 2<<4 + tv + 3, 3<<4 + tv + 4}
					}
					return lg.SSSGetCardType(cardData)
				}
				a = lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[2]]) - lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[0]])
				if a > 0 && a < 5 {
					tv := lg.GetCardLogicValue(cardData[2])
					cardData = []int{cardData[2], 0<<4 + tv - 1, 1<<4 + tv - 2, 2<<4 + tv - 3, 3<<4 + tv - 4}
					return lg.SSSGetCardType(cardData)
				}

				//三条
				cardData = []int{cardData[0], cardData[0], cardData[0], cardData[1], cardData[2]}
				return lg.SSSGetCardType(cardData)
			}
			//铁支
			if TagAnalyseItemArray.bTwoCount == 1 && TagAnalyseItemArray.bOneCount == 1 {
				tempCardData := cardData[TagAnalyseItemArray.bTwoFirst[0]]
				cardData = []int{tempCardData, tempCardData, tempCardData, tempCardData, cardData[TagAnalyseItemArray.bOneFirst[0]]}
				return lg.SSSGetCardType(cardData)
			}

		case 1:
			//五同
			if TagAnalyseItemArray.bFourCount == 1 {
				cardData = []int{cardData[0], cardData[0], cardData[0], cardData[0], cardData[0]}
				return lg.SSSGetCardType(cardData)
			}
			//铁支
			if TagAnalyseItemArray.bThreeCount == 1 {
				tempCardData := cardData[TagAnalyseItemArray.bThreeFirst[0]]
				cardData = []int{tempCardData, tempCardData, tempCardData, tempCardData, cardData[TagAnalyseItemArray.bOneFirst[0]]}
				return lg.SSSGetCardType(cardData)
			}
			//葫芦
			if TagAnalyseItemArray.bTwoCount == 2 {
				tempCardData := cardData[TagAnalyseItemArray.bTwoFirst[0]]
				cardData = []int{tempCardData, tempCardData, tempCardData, cardData[TagAnalyseItemArray.bTwoFirst[1]], cardData[TagAnalyseItemArray.bTwoFirst[1]]}
				return lg.SSSGetCardType(cardData)
			}

			if TagAnalyseItemArray.bOneCount == 4 {
				//同花顺
				if TagAnalyseItemArray.bflush {
					a := lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[0]]) - lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[3]])
					if a > 0 && a < 5 {
						if lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[3]]) > 9 {
							color := lg.GetCardColor(cardData[TagAnalyseItemArray.bOneFirst[3]])
							tempCardData := color<<4 + 10
							cardData = []int{tempCardData - 9, tempCardData, tempCardData + 1, tempCardData + 2, tempCardData + 3}
						} else {
							cardData = []int{cardData[3], cardData[3] + 1, cardData[3] + 2, cardData[3] + 3, cardData[3] + 4}
						}
						return lg.SSSGetCardType(cardData)
					}
					a = lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[3]]) - lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[0]])
					if a > 0 && a < 5 {
						cardData = []int{cardData[3], cardData[3] - 1, cardData[3] - 2, cardData[3] - 3, cardData[3] - 4}
						return lg.SSSGetCardType(cardData)
					}
					//同花
					cardData = []int{cardData[0], cardData[0], cardData[1], cardData[2], cardData[3]}
					return lg.SSSGetCardType(cardData)
				}

				//顺子
				a := lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[0]]) - lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[3]])
				if a > 0 && a < 5 {
					if lg.GetCardLogicValue(cardData[TagAnalyseItemArray.bOneFirst[3]]) > 9 {
						cardData = []int{0<<4 + 1, 1<<4 + 10, 2<<4 + 11, 3<<4 + 12, 0<<4 + 13}
					} else {
						tv := lg.GetCardLogicValue(cardData[3])
						cardData = []int{cardData[2], 0<<4 + tv + 1, 1<<4 + tv + 2, 2<<4 + tv + 3, 3<<4 + tv + 4}
					}
					return lg.SSSGetCardType(cardData)
				}
				a = lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[3]]) - lg.GetCardValue(cardData[TagAnalyseItemArray.bOneFirst[0]])
				if a > 0 && a < 5 {
					tv := lg.GetCardLogicValue(cardData[3])
					cardData = []int{cardData[2], 0<<4 + tv - 1, 1<<4 + tv - 2, 2<<4 + tv - 3, 3<<4 + tv - 4}
					return lg.SSSGetCardType(cardData)
				}

				//对子
				cardData = []int{cardData[TagAnalyseItemArray.bOneFirst[0]], cardData[TagAnalyseItemArray.bOneFirst[0]], cardData[TagAnalyseItemArray.bOneFirst[1]], cardData[TagAnalyseItemArray.bOneFirst[2]], cardData[TagAnalyseItemArray.bOneFirst[3]]}
				return lg.SSSGetCardType(cardData)
			}

			//三条
			if TagAnalyseItemArray.bTwoCount == 1 {
				tempCardData := cardData[TagAnalyseItemArray.bTwoFirst[0]]
				cardData = []int{tempCardData, tempCardData, tempCardData, cardData[TagAnalyseItemArray.bOneFirst[0]], cardData[TagAnalyseItemArray.bOneFirst[1]]}
				return lg.SSSGetCardType(cardData)
			}

		}

		//五同
		if TagAnalyseItemArray.bFiveCount == 1 {
			return CT_FIVE_SAME, TagAnalyseItemArray
		}

		//同花顺
		tempCardData := make([]int, metaCount)
		copy(tempCardData, cardData)
		if lg.IsLine(tempCardData, metaCount, true) {
			return CT_STRAIGHT_FLUSH, TagAnalyseItemArray
		}
		//铁支
		if 1 == TagAnalyseItemArray.bFourCount && 1 == TagAnalyseItemArray.bOneCount {
			return CT_FOURSAME, TagAnalyseItemArray
		}
		//葫芦
		if 1 == TagAnalyseItemArray.bThreeCount && 1 == TagAnalyseItemArray.bTwoCount {
			return CT_GOURD, TagAnalyseItemArray
		}
		//同花
		if TagAnalyseItemArray.bflush {
			return CT_FLUSH, TagAnalyseItemArray
		}
		//顺子
		tempCardData = make([]int, metaCount)
		copy(tempCardData, cardData)
		if lg.IsLine(tempCardData, metaCount, false) {
			return CT_STRAIGHT, TagAnalyseItemArray
		}
		//三条
		if 1 == TagAnalyseItemArray.bThreeCount && 2 == TagAnalyseItemArray.bOneCount {
			return CT_THREESAME, TagAnalyseItemArray
		}
		//两对
		if 2 == TagAnalyseItemArray.bTwoCount && 1 == TagAnalyseItemArray.bOneCount {
			return CT_TWOPAIR, TagAnalyseItemArray
		}
		//对子
		if 1 == TagAnalyseItemArray.bTwoCount && 3 == TagAnalyseItemArray.bOneCount {
			return CT_ONEPAIR, TagAnalyseItemArray
		}
		//散牌
		if 5 == TagAnalyseItemArray.bOneCount && false == TagAnalyseItemArray.bflush {
			return CT_SINGLE, TagAnalyseItemArray
		}
		//错误类型
		return CT_INVALID, TagAnalyseItemArray

	case 13: //13张特殊牌型
		//至尊清龙
		if 13 == TagAnalyseItemArray.bOneCount && true == TagAnalyseItemArray.bflush {
			return CT_THIRTEEN_FLUSH, TagAnalyseItemArray
		}
		//一条龙
		if 13 == TagAnalyseItemArray.bOneCount {
			return CT_THIRTEEN, TagAnalyseItemArray
		}

		//三同花顺
		btCardData := make([]int, metaCount)
		copy(btCardData, cardData)
		if lg.IsAllLine(btCardData, len(btCardData), true) {
			return CT_THREE_STRAIGHT_FLUSH, TagAnalyseItemArray
		}

		//三分天下
		if 3 == TagAnalyseItemArray.bFourCount {
			return CT_THREE_BOMB, TagAnalyseItemArray
		}

		//四套三条
		if 4 == TagAnalyseItemArray.bThreeCount {
			return CT_FOUR_THREESAME, TagAnalyseItemArray
		}

		//六对半
		if (6 == TagAnalyseItemArray.bTwoCount) ||
			(4 == TagAnalyseItemArray.bTwoCount && 1 == TagAnalyseItemArray.bFourCount) ||
			(2 == TagAnalyseItemArray.bTwoCount && 2 == TagAnalyseItemArray.bFourCount) ||
			(3 == TagAnalyseItemArray.bFourCount) {
			return CT_SIXPAIR, TagAnalyseItemArray
		}

		//三同花
		bThree_C := true
		for i := 0; i < 4; i++ {
			GetOutNum := lg.GetColorCardNum(cardData, i)
			if GetOutNum != 0 && GetOutNum != 3 && GetOutNum != 5 && GetOutNum != 8 && GetOutNum != 10 && GetOutNum != 13 {
				bThree_C = false
			}
		}
		for i := 0; i < len(cardData); i++ {
			if cardData[i] == 0x4E || cardData[i] == 0x4F {
				bThree_C = false
				break
			}
		}
		if bThree_C {
			return CT_THREE_FLUSH, TagAnalyseItemArray
		}

		//三顺子
		copy(btCardData, cardData)
		if lg.IsAllLine(btCardData, len(btCardData), false) {
			return CT_THREE_STRAIGHT, TagAnalyseItemArray
		}

	}

	return CT_INVALID, TagAnalyseItemArray
}

func (lg *sss_logic) SSSCompareCard(bInFirstList sssCardType, bInNextList sssCardType) int {
	bFirstCount := len(bInFirstList.Item.cardData)
	bNextCount := len(bInNextList.Item.cardData)

	if bFirstCount != bNextCount {
		//todo验证
		return 0
	}

	FirstAnalyseData := new(TagAnalyseItem)
	NextAnalyseData := new(TagAnalyseItem)

	bFirstList := make([]int, bFirstCount)
	bNextList := make([]int, bNextCount)

	copy(bFirstList, bInFirstList.Item.cardData)
	copy(bNextList, bInNextList.Item.cardData)

	// lg.SSSSortCardList(bFirstList)
	// lg.SSSSortCardList(bNextList)

	// FirstAnalyseData = lg.AnalyseCard(bFirstList)
	// NextAnalyseData = lg.AnalyseCard(bNextList)

	FirstAnalyseData = bInFirstList.Item
	NextAnalyseData = bInNextList.Item

	// bNextType := lg.GetCardType(bNextList)
	// bFirstType := lg.GetCardType(bFirstList)

	bNextType := bInNextList.CT
	bFirstType := bInFirstList.CT
	nextIsLaiZi := bInNextList.isLaiZi
	firstIsLaiZi := bInFirstList.isLaiZi

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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_ONEPAIR: //对带一张
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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_THREESAME: //三张牌型
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
					return -1
				}
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_ONEPAIR: //对带三张
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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_TWOPAIR: //两对牌型
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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_THREESAME: //三张牌型
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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_STRAIGHT: //顺子
				if lg.GetCardLogicValue(bNextList[0]) > lg.GetCardLogicValue(bFirstList[0]) {
					return 1
				}
				if lg.GetCardLogicValue(bNextList[0]) < lg.GetCardLogicValue(bFirstList[0]) {
					return -1
				}
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_FLUSH: //同花五牌
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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_GOURD: //三条一对
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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_FOURSAME: //四带一张

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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_STRAIGHT_FLUSH: //同花顺
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
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
					return -1
				}
				return 0
			case CT_FIVE_SAME: // 五同
				if lg.GetCardLogicValue(bNextList[0]) > lg.GetCardLogicValue(bFirstList[0]) {
					return 1
				} else {
					return -1
				}
				if nextIsLaiZi == false && firstIsLaiZi == true {
					return 1
				}
				if nextIsLaiZi == true && firstIsLaiZi == false {
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
			// case CT_THIRTEEN_FLUSH:
			// 	if lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) {
			// 		return 1
			// 	}
			// 	if lg.GetCardColor(bNextList[0]) < lg.GetCardColor(bFirstList[0]) {
			// 		return -1
			// 	}
			// 	return 0
			// case CT_TWELVE_KING:
			// 	for i := 0; i < 13; i++ {
			// 		if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
			// 			if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
			// 				return 1
			// 			} else {
			// 				return -1
			// 			}
			// 		}
			// 	}
			// 	return 0
			case CT_THREE_STRAIGHT_FLUSH:
				for i := 0; i < 13; i++ {
					if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
						if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
					if lg.GetCardColor(bNextList[i]) != lg.GetCardColor(bFirstList[i]) {
						if lg.GetCardColor(bNextList[i]) > lg.GetCardColor(bFirstList[i]) {
							return 1
						} else {
							return -1
						}
					}
				}
				return 0
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
			// case CT_ALL_BIG, CT_ALL_SMALL, CT_SAME_COLOR:
			// 	for i := 0; i < 13; i++ {
			// 		if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
			// 			if lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i]) {
			// 				return 1
			// 			} else {
			// 				return -1
			// 			}
			// 		}
			// 	}
			// 	return 0
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
			// case CT_FIVEPAIR_THREE:
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
			// 		return 1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
			// 		return -1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
			// 		return 1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
			// 		return -1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]]) {
			// 		return 1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]]) {
			// 		return -1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]]) {
			// 		return 1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]]) {
			// 		return -1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]]) {
			// 		return 1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]]) {
			// 		return -1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
			// 		return 1
			// 	}
			// 	if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) < lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
			// 		return -1
			// 	}
			// 	return 0
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

func (lg *sss_logic) IsAllLine(cbCard []int, cbCount int, bSameColor bool) bool {
	if cbCount == 0 {
		return true
	}

	if cbCount != 3 && cbCount != 5 && cbCount != 8 && cbCount != 10 && cbCount != 13 {
		return false
	}

	if cbCount == 3 || cbCount == 5 {
		return lg.IsLine(cbCard, cbCount, bSameColor)
	}

	cbTemp := make([]int, 13)

	for i := 0; i < cbCount; i++ {
		cbTemp[i] = cbCard[i]
	}

	cbIndex := []int{0, 1, 2, 3, 4}

	for {

		cbCarResult := make([]int, 5)
		for i := 0; i < 5; i++ {
			cbCarResult[i] = cbTemp[cbIndex[i]]
			cbTemp[cbIndex[i]] = 0
		}
		lg.SortCardByValue(cbTemp, cbCount)

		if lg.IsLine(cbCarResult, 5, bSameColor) && lg.IsAllLine(cbTemp, cbCount-5, bSameColor) {
			cbSortCount := 0
			for i := 0; i < 5; i++ {
				cbCard[cbSortCount] = cbCarResult[i]
				cbSortCount++
			}
			for i := 0; i < cbCount-5; i++ {
				cbCard[cbSortCount] = cbTemp[i]
				cbSortCount++
			}
			return true
		} else {
			for i := 0; i < cbCount; i++ {
				cbTemp[i] = cbCard[i]
			}
		}

		if cbIndex[4] == (cbCount - 1) {
			i := 4
			for ; i > 0; i-- {
				if (cbIndex[i-1] + 1) != cbIndex[i] {
					cbNewIndex := cbIndex[i-1]
					for j := i - 1; j < 5; j++ {
						cbIndex[j] = cbNewIndex + j - (i - 1) + 1
					}
					break
				}
			}
			if i == 0 {
				break
			}
		} else {
			cbIndex[4]++
		}

	}

	return false
}

func (lg *sss_logic) IsLine(cbCard []int, cbCount int, bSameColor bool) bool {

	lg.SortCardByValue(cbCard, cbCount)

	cbBossCount := 0
	cbFirst := 0xFF
	cbSecond := 0xFF
	cbLast := lg.GetCardLogicValue(cbCard[cbCount-1])
	cbSortCard := make([]int, len(cbCard))

	for i := 0; i < cbCount-1; i++ {
		if cbCard[i] >= 0x4E {
			cbBossCount++
			continue
		}

		if cbFirst == 0xFF {
			cbFirst = lg.GetCardLogicValue(cbCard[i])
		} else if cbSecond == 0xFF {
			cbSecond = lg.GetCardLogicValue(cbCard[i])
		}

		if bSameColor && lg.GetCardColor(cbCard[i]) != lg.GetCardColor(cbCard[i+1]) {
			return false
		}

		if lg.GetCardLogicValue(cbCard[i]) == lg.GetCardLogicValue(cbCard[i+1]) {
			return false
		}

	}

	if cbFirst == 14 && cbSecond < (2+cbCount-1) {
		cbSortCard[0] = lg.GetCardWithValue(cbCard, cbCount, 14)
		for i := 2; i < (2 + cbCount - 1); i++ {
			cbSortCard[i-1] = lg.GetCardWithValue(cbCard, cbCount, i)
		}
		for i := 0; i < cbCount; i++ {
			cbCard[i] = cbSortCard[i]
		}

		return true
	}

	if cbFirst-cbLast+1 <= cbCount {
		if cbFirst < (2 + cbCount - 1) {
			cbSortCard[0] = lg.GetCardWithValue(cbCard, cbCount, 14)
			for i := 2; i < (2 + cbCount - 1); i++ {
				cbSortCard[i-1] = lg.GetCardWithValue(cbCard, cbCount, i)
			}

		} else if cbLast > 14-cbCount {
			for i := 0; i < cbCount; i++ {
				cbSortCard[i] = lg.GetCardWithValue(cbCard, cbCount, 14-i)
			}

		} else {
			for i := 0; i < cbCount; i++ {
				cbSortCard[i] = lg.GetCardWithValue(cbCard, cbCount, cbLast+cbCount-1-i)
			}

		}

		for i := 0; i < cbCount; i++ {
			cbCard[i] = cbSortCard[i]
		}

		return true
	}

	return false
}

func (lg *sss_logic) SortCardByValue(cbCard []int, cbCount int) {
	for i := 0; i < cbCount-1; i++ {
		for j := i; j < cbCount; j++ {
			if lg.GetSortValue(cbCard[i]) < lg.GetSortValue(cbCard[j]) {
				cbTemp := cbCard[i]
				cbCard[i] = cbCard[j]
				cbCard[j] = cbTemp
			}
		}
	}
}

func (lg *sss_logic) GetSortValue(cbCard int) int {
	if cbCard >= 0x4E {
		return cbCard
	}
	return lg.GetCardLogicValue(cbCard)*4 + lg.GetCardColor(cbCard)
}

func (lg *sss_logic) GetCardWithValue(cbCard []int, cbCount int, cbValue int) int {
	cbTemp := 0
	for i := 0; i < cbCount; i++ {
		if lg.GetCardLogicValue(cbCard[i]) == cbValue {
			cbTemp = cbCard[i]
			cbCard[i] = 0
			return cbTemp
		}
	}

	for i := 0; i < cbCount; i++ {
		if lg.GetCardLogicValue(cbCard[i]) >= 15 {
			cbTemp = cbCard[i]
			cbCard[i] = 0
			return cbTemp
		}
	}

	return cbTemp
}

func (lg *sss_logic) GetColorCardNum(cbCard []int, cbColor int) int {
	GetOutNum := 0
	for i := 0; i < len(cbCard); i++ {
		if lg.GetCardColor(cbCard[i]) == cbColor {
			GetOutNum++
		}
	}
	return GetOutNum
}

func (lg *sss_logic) GetColorCard(cbCard []int, cbColor int) []int {
	colorCard := []int{}
	for i := 0; i < len(cbCard); i++ {
		if lg.GetCardColor(cbCard[i]) == cbColor {
			colorCard = append(colorCard, cbCard[i])
		}
	}
	return colorCard
}

func (lg *sss_logic) GetUniqueColorCard(cbCard []int, cbColor int) []int {
	allCards := make([]int, 16)
	uniqueColorCard := []int{}
	for i := 0; i < len(cbCard); i++ {
		if lg.GetCardColor(cbCard[i]) == cbColor {
			if allCards[lg.GetCardLogicValue(cbCard[i])] == 0 {
				uniqueColorCard = append(uniqueColorCard, cbCard[i])
				allCards[lg.GetCardLogicValue(cbCard[i])]++
			}
		}
	}
	return uniqueColorCard
}

func (lg *sss_logic) getUnUsedCard(cardData []int, usedCard []int) []int {
	tempCardData := []int{}
	for _, v := range cardData {
		exist := false
		for _, v1 := range usedCard {
			if v == v1 {
				exist = true
			}
		}
		if !exist {
			tempCardData = append(tempCardData, v)
		}
	}

	return tempCardData
}
