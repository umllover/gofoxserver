package room

import (
	"mj/gameServer/common/pk/pk_base"

	"mj/gameServer/common/pk"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

// 十三水逻辑
const (
	CT_INVALID                     = 0  //错误类型
	CT_SINGLE                      = 1  //单牌类型
	CT_ONE_DOUBLE                  = 2  //只有一对
	CT_FIVE_TWO_DOUBLE             = 3  //两对牌型
	CT_THREE                       = 4  //三张牌型
	CT_FIVE_MIXED_FLUSH_NO_A       = 5  //没A杂顺
	CT_FIVE_MIXED_FLUSH_FIRST_A    = 6  //A在前顺子
	CT_FIVE_MIXED_FLUSH_BACK_A     = 7  //A在后顺子
	CT_FIVE_FLUSH                  = 8  //同花五牌
	CT_FIVE_THREE_DEOUBLE          = 9  //三条一对
	CT_FIVE_FOUR_ONE               = 10 //四带一张
	CT_FIVE_STRAIGHT_FLUSH_NO_A    = 11 //没A同花顺
	CT_FIVE_STRAIGHT_FLUSH_FIRST_A = 12 //A在前同花顺
	CT_FIVE_STRAIGHT_FLUSH_BACK_A  = 13 //A在后同花顺
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
	util.DeepCopy(&bTempCardData, &bCardData)
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

		return true
	}
	return true
}

//逻辑数值
func (lg *sss_logic) GetCardLogicValue(bCardData int) int {
	//扑克属性
	bCardValue := lg.GetCardValue(bCardData)

	//转换数值
	if bCardValue == 1 {
		bCardValue += 13
		return bCardValue
	} else {
		return bCardValue
	}
}

//获取数值
func (lg *sss_logic) GetCardValue(bCardData int) int { return bCardData & LOGIC_MASK_VALUE } //十六进制前面四位表示牌的数值
//获取花色
func (lg *sss_logic) GetCardColor(bCardData int) int { return (bCardData & LOGIC_MASK_COLOR) >> 4 } //十六进制后面四位表示牌的花色

func (lg *sss_logic) AnalyseCard(bCardDataList []int, bCardCount int, TagAnalyseItemArray *TagAnalyseItem) *TagAnalyseItem {

	cbBufferCount := int(len(bCardDataList))
	bCardData := make([]int, cbBufferCount)
	util.DeepCopy(&bCardData, &bCardDataList)

	//变量定义
	bSameCount := 1
	bCardValueTemp := 0
	bSameColorCount := 1
	bFirstCardIndex := 0 //记录下标

	bLogicValue := lg.GetCardLogicValue(bCardData[0])
	bCardColor := lg.GetCardColor(bCardData[0])

	analyseItem := &TagAnalyseItem{bOneFirst: make([]int, 13), bTwoFirst: make([]int, 13), bThreeFirst: make([]int, 13), bFourFirst: make([]int, 13)}
	//扑克分析
	for i := 1; i < bCardCount; i++ {
		//获取扑克
		bCardValueTemp = lg.GetCardLogicValue(bCardData[i])

		if bCardValueTemp == bLogicValue {
			bSameCount++

		}
		if bCardValueTemp != bLogicValue || i == (bCardCount-1) {
			switch bSameCount {
			case 1: //一张
			case 2: //两张
				{
					analyseItem.bTwoFirst[analyseItem.bTwoCount] = bFirstCardIndex
					analyseItem.bTwoCount++
				}
			case 3:
				{
					analyseItem.bThreeFirst[analyseItem.bThreeCount] = bFirstCardIndex
					analyseItem.bThreeCount++
				}
			case 4:
				{
					analyseItem.bFourFirst[analyseItem.bFourCount] = bFirstCardIndex
					analyseItem.bFourCount++
				}
			}
		}

		//设置变量
		if bCardValueTemp != bLogicValue {
			if bSameCount == 1 {
				if i != bCardCount-1 {
					analyseItem.bOneFirst[analyseItem.bOneCount] = bFirstCardIndex
					analyseItem.bOneCount++
				} else {
					analyseItem.bOneFirst[analyseItem.bOneCount] = bFirstCardIndex
					analyseItem.bOneCount++
					analyseItem.bOneFirst[analyseItem.bOneCount] = i
					analyseItem.bOneCount++
				}
			} else {
				if i == bCardCount-1 {
					analyseItem.bOneFirst[analyseItem.bOneCount] = i
					analyseItem.bOneCount++
				}
			}
			bSameCount = 1
			bLogicValue = bCardValueTemp
			bFirstCardIndex = i
		}
		if lg.GetCardColor(bCardData[i]) != bCardColor {
			bSameColorCount = 1
		} else {
			bSameColorCount++
		}
	}

	if bCardCount == bSameColorCount {
		analyseItem.bStraight = true
	} else {
		analyseItem.bStraight = false
	}
	return analyseItem
}
func (lg *sss_logic) GetSSSCardType(cardData []int, bCardCount int, btSpecialCard []int) int {
	CardCount := len(cardData)
	if CardCount != 3 && CardCount != 5 && CardCount != 13 {
		return CT_INVALID
	}

	TagAnalyseItemArray := new(TagAnalyseItem)
	TagAnalyseItemArray = lg.AnalyseCard(cardData, bCardCount, TagAnalyseItemArray)

	//开始分析
	switch bCardCount {
	case 3: //三条类型
		{
			//单牌类型
			if TagAnalyseItemArray.bOneCount == 3 {
				return CT_SINGLE
			}

			//对带一张
			if TagAnalyseItemArray.bTwoCount == 1 && 1 == TagAnalyseItemArray.bOneCount {
				return CT_ONE_DOUBLE
			}

			//三张牌型
			if TagAnalyseItemArray.bThreeCount == 1 {
				return CT_THREE
			}

			//错误类型
			return CT_INVALID
		}
	case 5: //五张牌型
		{

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
				if lg.GetCardLogicValue(cardData[4]) != 0 {
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
		}
	case 13: //13张特殊牌型
		{
			TwelveKing := false
			//同花十三水
			if 13 == TagAnalyseItemArray.bOneCount && true == TagAnalyseItemArray.bStraight {
				return CT_THIRTEEN_FLUSH
			}
			//十三水
			if 13 == TagAnalyseItemArray.bOneCount {
				return CT_THIRTEEN
			}

			TwelveKing = true
			for i := 0; i < 13; i++ {
				if lg.GetCardLogicValue(cardData[i]) < 11 {
					TwelveKing = false
					break
				}
			}
			if TwelveKing {
				return CT_TWELVE_KING
			}

			//三同花顺
			btCardData := make([]int, 13)
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
				if FCardData != lg.GetCardLogicValue(btCardData[i])+1 && SColor == lg.GetCardColor(btCardData[i]) {
					if 3 == StraightFlush {
						StraightFlush1 = true
						Count1 = 3
						lg.RemoveCard(RbtCardData, 3, btCardData, 13)
						//ZeroMemory(RbtCardData, sizeof(RbtCardData))
					}
					break
				}
				if 5 == StraightFlush {
					StraightFlush1 = true
					Count1 = 5
					lg.RemoveCard(RbtCardData, 5, btCardData, 13)
					//ZeroMemory(RbtCardData, sizeof(RbtCardData))
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
							StraightFlush1 = true
							Count2 = 3
							lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
							//ZeroMemory(RbtCardData,sizeof(RbtCardData))
						}
						break
					}
					if 5 == StraightFlush {
						StraightFlush2 = true
						Count2 = 5
						lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1)
						//ZeroMemory(RbtCardData,sizeof(RbtCardData));
						break
					}
				}
			}
			if StraightFlush2 {
				StraightFlush = 1
				Number = 0
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
							StraightFlush1 = true
							Count3 = 3
							lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
							//ZeroMemory(RbtCardData,sizeof(RbtCardData))
						}
						break
					}
					if 5 == StraightFlush {
						StraightFlush3 = true
						Count3 = 5
						lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1-Count2)
						//ZeroMemory(RbtCardData,sizeof(RbtCardData));
						break
					}
				}
			}
			if StraightFlush1 && StraightFlush2 && StraightFlush3 && Count1+Count2+Count3 == 13 {
				return CT_THREE_STRAIGHTFLUSH
			}
			//三炸弹
			if 3 == TagAnalyseItemArray.bFourCount {
				return CT_THREE_BOMB
			}
			//全大
			AllBig := true
			for i := 0; i < 13; i++ {
				if lg.GetCardLogicValue(cardData[i]) < 8 {
					AllBig = false
					break
				}
			}
			if AllBig {
				return CT_ALL_BIG
			}
			//全小
			AllSmall := true
			for i := 0; i < 13; i++ {
				if lg.GetCardLogicValue(cardData[i]) > 8 {
					AllSmall = false
					break
				}
			}
			if AllSmall {
				return CT_ALL_SMALL
			}
			//凑一色
			Flush := 1
			SColor = lg.GetCardColor(cardData[0]) & 0x01
			for i := 1; i < 13; i++ {
				if SColor == lg.GetCardColor(cardData[i])&0x01 {
					Flush++
				} else {
					break
				}
			}
			if 13 == Flush {
				return CT_SAME_COLOR
			}
			//四套冲三
			if 4 == TagAnalyseItemArray.bThreeCount {
				return CT_FOUR_THREESAME
			}
			//五对冲三
			if (5 == TagAnalyseItemArray.bTwoCount && 1 == TagAnalyseItemArray.bThreeCount) ||
				(3 == TagAnalyseItemArray.bTwoCount && 1 == TagAnalyseItemArray.bFourCount && 1 == TagAnalyseItemArray.bThreeCount) ||
				(1 == TagAnalyseItemArray.bTwoCount && 2 == TagAnalyseItemArray.bFourCount && 1 == TagAnalyseItemArray.bThreeCount) {
				return CT_FIVEPAIR_THREE
			}
			//六对半
			if (6 == TagAnalyseItemArray.bTwoCount) || (4 == TagAnalyseItemArray.bTwoCount && 1 == TagAnalyseItemArray.bFourCount) ||
				(2 == TagAnalyseItemArray.bTwoCount && 2 == TagAnalyseItemArray.bFourCount) || (3 == TagAnalyseItemArray.bFourCount) {
				return CT_SIXPAIR
			}
			//三同花
			Flush1 := false
			Flush2 := false
			Flush3 := false
			Flush = 1
			Count1 = 0
			Count2 = 0
			Count3 = 0
			Number = 0
			RbtCardData = make([]int, 13)
			util.DeepCopy(&btCardData, &cardData)

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
					copy(btSpecialCard[10:], RbtCardData[:Count1])
					lg.RemoveCard(RbtCardData, 3, btCardData, 13)
					//RbtCardData = RbtCardData[:0]
					RbtCardData = make([]int, 13)
					break
				}
				if 5 == Flush {
					Flush1 = true
					Count1 = 5
					//util.DeepCopy(&btSpecialCard[5], RbtCardData)
					copy(btSpecialCard[5:], RbtCardData[:Count1])
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
						util.DeepCopy(&btSpecialCard[10], RbtCardData)
						lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
						RbtCardData = RbtCardData[:0]
						break
					}
					if 5 == Flush {
						Flush2 = true
						Count2 = 5
						if Count1 == 5 {
							util.DeepCopy(&btSpecialCard[0], RbtCardData)
						} else if Count1 == 3 {
							util.DeepCopy(&btSpecialCard[5], RbtCardData)
						}

						lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1)
						RbtCardData = RbtCardData[:0]
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
						util.DeepCopy(&btSpecialCard[10], RbtCardData)
						lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
						RbtCardData = RbtCardData[:0]
						break
					}
					if 5 == Flush {
						Flush3 = true
						Count3 = 5
						util.DeepCopy(&btSpecialCard[0], RbtCardData)
						lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1-Count2)
						RbtCardData = RbtCardData[:0]
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
					RbtCardData = RbtCardData[:0]
					util.DeepCopy(btCardData, cardData)
					lg.SortCardList(btCardData, 13)
					RbtCardData[Number] = btCardData[0]
					Number++
					FCardData = lg.GetCardLogicValue(btCardData[0])
					for i := 1; i < 13; i++ {
						if FCardData == lg.GetCardLogicValue(btCardData[i])+1 || (FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 5) || (FCardData == 14 && lg.GetCardLogicValue(btCardData[i]) == 3) {
							Straight++
							RbtCardData[Number] = btCardData[i]
							Number++
							FCardData = lg.GetCardLogicValue(btCardData[i])

						} else if FCardData != lg.GetCardLogicValue(btCardData[i]) {
							if 3 == Straight {
								Straight1 = true
								Count1 = 3
								util.DeepCopy(&btSpecialCard[10], RbtCardData)
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
								util.DeepCopy(&btSpecialCard[10], RbtCardData)
								lg.RemoveCard(RbtCardData, 3, btCardData, 13)
								RbtCardData = RbtCardData[:0]
								break

							}
						} else if nCount == 2 || nCount == 3 {
							if 3 == Straight {

								Straight1 = true
								Count1 = 3
								util.DeepCopy(&btSpecialCard[10], RbtCardData)
								lg.RemoveCard(RbtCardData, 3, btCardData, 13)
								RbtCardData = RbtCardData[:0]
								break
							}
						}
						if 5 == Straight {
							Straight1 = true
							Count1 = 5
							util.DeepCopy(&btSpecialCard[5], RbtCardData)
							lg.RemoveCard(RbtCardData, 5, btCardData, 13)
							RbtCardData = RbtCardData[:0]
							break

						}
					}
					if Straight1 {
						Straight = 1
						Number = 0
						lg.SortCardList(btCardData, 13-Count1)
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
									util.DeepCopy(&btSpecialCard[10], RbtCardData)
									lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
									RbtCardData = RbtCardData[:0]
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
									util.DeepCopy(&btSpecialCard[10], RbtCardData)
									lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
									RbtCardData = RbtCardData[:0]
									break

								}
							} else if nCount == 1 || nCount == 3 {
								if 3 == Straight && Count1 != 3 {
									Straight2 = true
									Count2 = 3
									util.DeepCopy(&btSpecialCard[10], RbtCardData)
									lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1)
									RbtCardData = RbtCardData[:0]
									break

								}
							}
							if 5 == Straight {
								Straight2 = true
								Count2 = 5
								if Count1 == 5 {
									util.DeepCopy(&btSpecialCard[0], RbtCardData)
								} else {
									util.DeepCopy(&btSpecialCard[5], RbtCardData)
								}

								lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1)
								RbtCardData = RbtCardData[:0]
								break
							}
						}
					}
					if Straight2 {
						Straight = 1
						Number = 0
						lg.SortCardList(btCardData, 13-Count1-Count2)
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
									util.DeepCopy(&btSpecialCard[10], RbtCardData)
									lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
									RbtCardData = RbtCardData[:0]
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
								util.DeepCopy(&btSpecialCard[10], RbtCardData)
								lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
								RbtCardData = RbtCardData[:0]
								break
							}
							if 5 == Straight {
								Straight3 = true
								Count3 = 5
								util.DeepCopy(&btSpecialCard[0], RbtCardData)
								lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1-Count2)
								RbtCardData = RbtCardData[:0]
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
				util.DeepCopy(&btCardData, &cardData)
				lg.SortCardList(btCardData, 13)
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
							copy(btSpecialCard[10:], RbtCardData[:Count1])
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
							copy(btSpecialCard[10:], RbtCardData[:Count1])
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
							copy(btSpecialCard[10:], RbtCardData[:Count1])
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
						copy(btSpecialCard[5:], RbtCardData[:Count1])
						lg.RemoveCard(RbtCardData, 5, btCardData, 13)
						//RbtCardData = RbtCardData[:0]
						RbtCardData = make([]int, 13)
						break

					}
				}
				if Straight1 {
					Straight = 1
					Number = 0
					lg.SortCardList(btCardData, 13-Count1)
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
								util.DeepCopy(&btSpecialCard[10], RbtCardData)
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
								util.DeepCopy(&btSpecialCard[10], RbtCardData)
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
								copy(btSpecialCard[10:], RbtCardData[:Count2])
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
								copy(btSpecialCard, RbtCardData[:Count2])
							} else {
								//util.DeepCopy(&btSpecialCard[5], RbtCardData)
								copy(btSpecialCard[5:], RbtCardData[:Count2])
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
					lg.SortCardList(btCardData, 13-Count1-Count2)
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
								copy(btSpecialCard[10:], RbtCardData[:Count3])
								lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
								RbtCardData = RbtCardData[:0]
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
							copy(btSpecialCard[10:], RbtCardData[:Count3])
							lg.RemoveCard(RbtCardData, 3, btCardData, 13-Count1-Count2)
							RbtCardData = RbtCardData[:0]
							break
						}
						if 5 == Straight {
							Straight3 = true
							Count3 = 5
							//util.DeepCopy(&btSpecialCard[0], RbtCardData)
							copy(btSpecialCard[:], RbtCardData[:Count3])
							lg.RemoveCard(RbtCardData, 5, btCardData, 13-Count1-Count2)
							RbtCardData = RbtCardData[:0]
							break
						}
					}
				}
				if Straight1 && Straight2 && Straight3 && Count1+Count2+Count3 == 13 {
					return CT_THREE_STRAIGHT
				}
			}
		}
	}

	return CT_INVALID
}
func (lg *sss_logic) GetType(bCardData []int, bCardCount int) *pk.TagAnalyseType {
	CardData := make([]int, 5)
	Type := new(pk.TagAnalyseType)
	Type.CbOnePare = make([]int, 100)
	Type.CbTwoPare = make([]int, 100)
	Type.CbThreeSame = make([]int, 100)
	Type.CbStraight = make([]int, 100)
	Type.CbFlush = make([]int, 100)
	Type.CbStraightFlush = make([]int, 100)
	Type.CbFourSame = make([]int, 100)
	Type.CbGourd = make([]int, 100)

	Type.BbOnePare = make([]bool, 20)
	Type.BbTwoPare = make([]bool, 20)
	Type.BbThreeSame = make([]bool, 20)
	Type.BbStraight = make([]bool, 20)
	Type.BbFlush = make([]bool, 20)
	Type.BbGourd = make([]bool, 20)
	Type.BbFourSame = make([]bool, 20)
	Type.BbStraightFlush = make([]bool, 20)

	util.DeepCopy(&CardData, &bCardData)
	log.Debug("%d   xxxxxxxxxxxxx", CardData)
	lg.SortCardList(CardData, 5)
	Index := make([]int, 5)
	Number := 0
	SameValueCount := 1
	Num := make([]int, 8)
	bLogicValue := lg.GetCardLogicValue(CardData[0])
	Index[Number] = 0
	Number++
	for i := 1; i < bCardCount; i++ {
		if bLogicValue == lg.GetCardLogicValue(CardData[i]) {
			SameValueCount++
			Index[Number] = i
			Number++
		}
		if bLogicValue != lg.GetCardLogicValue(CardData[i]) || i == bCardCount-1 {
			if SameValueCount == 1 {

			} else if SameValueCount == 2 {
				Type.BOnePare = true
				Type.CbOnePare[Num[0]] = Index[SameValueCount-2]
				Num[0]++
				Type.CbOnePare[Num[0]] = Index[SameValueCount-1]
				Num[0]++
				Type.BtOnePare++
			} else if SameValueCount == 3 {
				Type.BOnePare = true
				Type.CbOnePare[Num[0]] = Index[SameValueCount-3]
				Num[0]++
				Type.CbOnePare[Num[0]] = Index[SameValueCount-2]
				Num[0]++
				Type.BThreeSame = true
				Type.CbThreeSame[Num[2]] = Index[SameValueCount-3]
				Num[2]++
				Type.CbThreeSame[Num[2]] = Index[SameValueCount-2]
				Num[2]++
				Type.CbThreeSame[Num[2]] = Index[SameValueCount-1]
				Num[2]++
				Type.BtThreeSame++
			} else if SameValueCount == 4 {
				Type.BOnePare = true
				Type.CbOnePare[Num[0]] = Index[SameValueCount-4]
				Num[0]++
				Type.CbOnePare[Num[0]] = Index[SameValueCount-3]
				Num[0]++
				Type.BThreeSame = true
				Type.CbThreeSame[Num[2]] = Index[SameValueCount-4]
				Num[2]++
				Type.CbThreeSame[Num[2]] = Index[SameValueCount-3]
				Num[2]++
				Type.CbThreeSame[Num[2]] = Index[SameValueCount-2]
				Num[2]++
				Type.BFourSame = true
				Type.CbFourSame[Num[6]] = Index[SameValueCount-4]
				Num[6]++
				Type.CbFourSame[Num[6]] = Index[SameValueCount-3]
				Num[6]++
				Type.CbFourSame[Num[6]] = Index[SameValueCount-2]
				Num[6]++
				Type.CbFourSame[Num[6]] = Index[SameValueCount-1]
				Num[6]++
				Type.BtFourSame++
			} else {

			}
			Number = 0
			//ZeroMemory(Index,sizeof(Index));
			Index[Number] = i
			Number++
			SameValueCount = 1
			bLogicValue = lg.GetCardLogicValue(CardData[i])
		}

	}
	//判断两对
	OnePareCount := Num[0] / 2
	ThreeSameCount := Num[2] / 3
	if OnePareCount >= 2 {
		Type.BTwoPare = true
		for i := 0; i < OnePareCount; i++ {
			for j := i + 1; j < OnePareCount; j++ {
				Type.CbTwoPare[Num[1]] = Type.CbOnePare[i*2]
				Num[1]++
				Type.CbTwoPare[Num[1]] = Type.CbOnePare[i*2+1]
				Num[1]++
				Type.CbTwoPare[Num[1]] = Type.CbOnePare[j*2]
				Num[1]++
				Type.CbTwoPare[Num[1]] = Type.CbOnePare[j*2+1]
				Num[1]++
				Type.BtTwoPare++
			}
		}
	}
	//判断葫芦
	if OnePareCount > 0 && ThreeSameCount > 0 {
		for i := 0; i < ThreeSameCount; i++ {
			for j := 0; j < OnePareCount; j++ {
				if lg.GetCardLogicValue(Type.CbThreeSame[i*3]) == lg.GetCardLogicValue(Type.CbOnePare[j*2]) {
					continue
				}
				Type.BGourd = true
				Type.CbGourd[Num[5]] = Type.CbThreeSame[i*3]
				Num[5]++
				Type.CbGourd[Num[5]] = Type.CbThreeSame[i*3+1]
				Num[5]++
				Type.CbGourd[Num[5]] = Type.CbThreeSame[i*3+2]
				Num[5]++
				Type.CbGourd[Num[5]] = Type.CbOnePare[j*2]
				Num[5]++
				Type.CbGourd[Num[5]] = Type.CbOnePare[j*2+1]
				Num[5]++
				Type.BtGourd++
			}
		}
	}
	//判断顺子及同花顺
	Number = 0
	//ZeroMemory(Index,sizeof(Index))
	Straight := 1
	bStraight := lg.GetCardLogicValue(CardData[0])
	Index[Number] = 0
	Number++
	if bStraight != 14 {
		for i := 1; i < bCardCount; i++ {
			if bStraight == lg.GetCardLogicValue(CardData[i])+1 {
				Straight++
				Index[Number] = i
				Number++
				bStraight = lg.GetCardLogicValue(CardData[i])
			}
			if bStraight > lg.GetCardLogicValue(CardData[i])+1 || i == bCardCount-1 {
				if Straight >= 5 {
					Type.BStraight = true
					for j := 0; j < Straight; j++ {
						if Straight-j >= 5 {
							Type.CbStraight[Num[3]] = Index[j]
							Num[3]++
							Type.CbStraight[Num[3]] = Index[j+1]
							Num[3]++
							Type.CbStraight[Num[3]] = Index[j+2]
							Num[3]++
							Type.CbStraight[Num[3]] = Index[j+3]
							Num[3]++
							Type.CbStraight[Num[3]] = Index[j+4]
							Num[3]++
							Type.BtStraight++
							//从手牌中找到和顺子5张中其中一张数值相同的牌，组成另一种顺子
							for k := j; k < j+5; k++ {
								for m := 0; m < bCardCount; m++ {
									if lg.GetCardLogicValue(CardData[Index[k]]) == lg.GetCardLogicValue(CardData[m]) && lg.GetCardColor(CardData[Index[k]]) != lg.GetCardColor(CardData[m]) {
										for n := j; n < j+5; n++ {
											if n == k {
												Type.CbStraight[Num[3]] = m
												Num[3]++
											} else {
												Type.CbStraight[Num[3]] = Index[n]
												Num[3]++
											}
										}
										Type.BtStraight++
									}
								}
							}
						} else {
							break
						}
					}

				}
				if bCardCount-i < 5 {
					break
				}
				bStraight = lg.GetCardLogicValue(CardData[i])
				Straight = 1
				Number = 0
				//ZeroMemory(Index,sizeof(Index));
				Index[Number] = i
				Number++
			}
		}

	}
	if bStraight == 14 {
		for i := 1; i < bCardCount; i++ {
			if bStraight == lg.GetCardLogicValue(CardData[i])+1 {
				Straight++
				Index[Number] = i
				Number++
				bStraight = lg.GetCardLogicValue(CardData[i])
			}
			if bStraight > lg.GetCardLogicValue(CardData[i])+1 || i == bCardCount-1 {
				if Straight >= 5 {
					Type.BStraight = true
					for j := 0; j < Straight; j++ {
						if Straight-j >= 5 {
							Type.CbStraight[Num[3]] = Index[j]
							Num[3]++
							Type.CbStraight[Num[3]] = Index[j+1]
							Num[3]++
							Type.CbStraight[Num[3]] = Index[j+2]
							Num[3]++
							Type.CbStraight[Num[3]] = Index[j+3]
							Num[3]++
							Type.CbStraight[Num[3]] = Index[j+4]
							Num[3]++
							Type.BtStraight++
							//从手牌中找到和顺子5张中其中一张数值相同的牌，组成另一种顺子
							for k := j; k < j+5; k++ {
								for m := 0; m < bCardCount; m++ {
									if lg.GetCardLogicValue(CardData[Index[k]]) == lg.GetCardLogicValue(CardData[m]) && lg.GetCardColor(CardData[Index[k]]) != lg.GetCardColor(CardData[m]) {
										for n := j; n < j+5; n++ {
											if n == k {
												Type.CbStraight[Num[3]] = m
												Num[3]++
											} else {
												Type.CbStraight[Num[3]] = Index[n]
												Num[3]++
											}
										}
										Type.BtStraight++
									}
								}
							}

						} else {
							break
						}
					}
				}
				if bCardCount-i < 5 {
					break
				}
				bStraight = lg.GetCardLogicValue(CardData[i])
				Straight = 1
				Number = 0
				//ZeroMemory(Index,sizeof(Index));
				Index[Number] = i
				Number++
			}
		}
		if lg.GetCardLogicValue(CardData[bCardCount-1]) == 2 {
			Number = 0
			BackA := 1
			FrontA := 1
			bStraight = lg.GetCardLogicValue(CardData[0])
			//ZeroMemory(Index,sizeof(Index));
			Index[Number] = 0
			Number++
			bStraight = lg.GetCardLogicValue(CardData[bCardCount-1])
			Index[Number] = bCardCount - 1
			Number++
			for i := bCardCount - 2; i >= 0; i-- {
				if bStraight == lg.GetCardLogicValue(CardData[i])-1 {
					FrontA++
					Index[Number] = i
					Number++
					bStraight = lg.GetCardLogicValue(CardData[i])
				}
			}
			if FrontA+BackA >= 5 {
				Type.BStraight = true
				for i := BackA; i > 0; i-- {
					for j := 1; j <= FrontA; j++ {
						if i+j == 5 {
							for k := 0; k < i; k++ {
								Type.CbStraight[Num[3]] = Index[k]
								Num[3]++
							}
							for k := 0; k < j; k++ {
								Type.CbStraight[Num[3]] = Index[k+BackA]
								Num[3]++
							}
							break
						}
					}
				}

			}
		}

	}
	//判断同花及同花顺
	Number = 0
	//ZeroMemory(Index,sizeof(Index));
	lg.SortCardList(CardData, bCardCount)
	cbCardData := make([]int, 13)
	util.DeepCopy(&cbCardData, &bCardData)
	lg.SortCardList(cbCardData, bCardCount)
	SameColorCount := 1
	bCardColor := lg.GetCardColor(CardData[0])
	Index[Number] = 0
	Number++
	for i := 1; i < bCardCount; i++ {
		if bCardColor == lg.GetCardColor(CardData[i]) {
			SameColorCount++
			Index[Number] = i
			Number++
		}
		if bCardColor != lg.GetCardColor(CardData[i]) || i == bCardCount-1 {
			if SameColorCount >= 5 {
				Type.BFlush = true

				for j := 0; j < SameColorCount; j++ {
					for k := 0; k < bCardCount; k++ {
						if lg.GetCardLogicValue(CardData[Index[j]]) == lg.GetCardLogicValue(cbCardData[k]) && lg.GetCardColor(CardData[Index[j]]) == lg.GetCardColor(cbCardData[k]) {
							Index[j] = k
							break
						}
					}
				}
				SaveIndex := 0
				for j := 0; j < SameColorCount; j++ {
					for k := j + 1; k < SameColorCount; k++ {
						if Index[j] > Index[k] {
							SaveIndex = Index[j]
							Index[j] = Index[k]
							Index[k] = SaveIndex
						}
					}
				}
				for j := 0; j < SameColorCount; j++ {
					if SameColorCount-j >= 5 {
						Type.CbFlush[Num[4]] = Index[j]
						Num[4]++
						Type.CbFlush[Num[4]] = Index[j+1]
						Num[4]++
						Type.CbFlush[Num[4]] = Index[j+2]
						Num[4]++
						Type.CbFlush[Num[4]] = Index[j+3]
						Num[4]++
						Type.CbFlush[Num[4]] = Index[j+4]
						Num[4]++
						Type.BtFlush++
						if lg.GetCardLogicValue(cbCardData[Index[j]]) == 14 {
							if lg.GetCardLogicValue(cbCardData[Index[j+1]]) == 5 && lg.GetCardLogicValue(cbCardData[Index[j+2]]) == 4 && lg.GetCardLogicValue(cbCardData[Index[j+3]]) == 3 && lg.GetCardLogicValue(cbCardData[Index[j+4]]) == 2 {
								Type.BStraightFlush = true
								Type.CbStraightFlush[Num[7]] = Index[j]
								Num[7]++
								Type.CbStraightFlush[Num[7]] = Index[j+1]
								Num[7]++
								Type.CbStraightFlush[Num[7]] = Index[j+2]
								Num[7]++
								Type.CbStraightFlush[Num[7]] = Index[j+3]
								Num[7]++
								Type.CbStraightFlush[Num[7]] = Index[j+4]
								Num[7]++
								Type.BtStraightFlush++
							}

						}
						if lg.GetCardLogicValue(cbCardData[Index[j]]) == lg.GetCardLogicValue(cbCardData[Index[j+1]])+1 &&
							lg.GetCardLogicValue(cbCardData[Index[j]]) == lg.GetCardLogicValue(cbCardData[Index[j+2]])+2 &&
							lg.GetCardLogicValue(cbCardData[Index[j]]) == lg.GetCardLogicValue(cbCardData[Index[j+3]])+3 &&
							lg.GetCardLogicValue(cbCardData[Index[j]]) == lg.GetCardLogicValue(cbCardData[Index[j+4]])+4 {
							Type.BStraightFlush = true
							Type.CbStraightFlush[Num[7]] = Index[j]
							Num[7]++
							Type.CbStraightFlush[Num[7]] = Index[j+1]
							Num[7]++
							Type.CbStraightFlush[Num[7]] = Index[j+2]
							Num[7]++
							Type.CbStraightFlush[Num[7]] = Index[j+3]
							Num[7]++
							Type.CbStraightFlush[Num[7]] = Index[j+4]
							Num[7]++
							Type.BtStraightFlush++
						}

					} else {
						break
					}
				}
			}
			if bCardCount-i < 5 {
				break
			}
			Number = 0
			SameColorCount = 1
			Index[Number] = i
			Number++
			bCardColor = lg.GetCardColor(CardData[i])
		}
	}
	return Type
}
func (lg *sss_logic) CompareSSSCard(bInFirstList []int, bInNextList []int, bFirstCount int, bNextCount int, bComPerWithOther bool) bool {

	FirstAnalyseData := new(TagAnalyseItem)
	NextAnalyseData := new(TagAnalyseItem)

	bFirstList := make([]int, 13)
	bNextList := make([]int, 13)

	util.DeepCopy(&bFirstList, &bInFirstList)
	util.DeepCopy(&bNextList, &bInNextList)

	lg.SortCardList(bFirstList, bFirstCount)
	lg.SortCardList(bNextList, bNextCount)

	FirstAnalyseData = lg.AnalyseCard(bFirstList, bFirstCount, FirstAnalyseData)
	NextAnalyseData = lg.AnalyseCard(bNextList, bNextCount, NextAnalyseData)


	if bFirstCount != (FirstAnalyseData.bOneCount + FirstAnalyseData.bTwoCount*2 + FirstAnalyseData.bThreeCount*3 + FirstAnalyseData.bFourCount*4 + FirstAnalyseData.bFiveCount*5) {
		return false
	}
	if bNextCount != (NextAnalyseData.bOneCount + NextAnalyseData.bTwoCount*2 + NextAnalyseData.bThreeCount*3 + NextAnalyseData.bFourCount*4 + NextAnalyseData.bFiveCount*5) {
		return false
	}
	if !((bFirstCount == bNextCount) || (bFirstCount != bNextCount && (3 == bFirstCount && 5 == bNextCount || 5 == bFirstCount && 3 == bNextCount))) {
		return false
	}
	bNextType := lg.GetSSSCardType(bNextList, bNextCount, lg.BtCardSpecialData)
	bFirstType := lg.GetSSSCardType(bFirstList, bFirstCount, lg.BtCardSpecialData)

	if CT_INVALID == bFirstType || CT_INVALID == bNextType {
		return false
	}
	//头段比较
	if true == bComPerWithOther {
		if 3 == bFirstCount {
			//开始对比
			if bNextType == bFirstType {
				switch bFirstType {
				case CT_SINGLE: //单牌类型
					{
						if bNextList[0] == bFirstList[0] {
							return false
						}
						bAllSame := true
						for i := 0; i < 3; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								bAllSame = false
								break
							}
						}
						if true == bAllSame {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) //比较花色
						} else {
							for i := 0; i < 3; i++ {
								if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
									return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
								}
							}
							return false
						}
						return false
					}
				case CT_ONE_DOUBLE: //对带一张
					{
						if bNextList[NextAnalyseData.bTwoFirst[0]] == bFirstList[FirstAnalyseData.bTwoFirst[0]] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
							if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) != lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
								return lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]])
							} else {
								return lg.GetCardColor(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较花色
							}
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]])
						}
					}
				case CT_THREE: //三张牌型
					{
						if bNextList[NextAnalyseData.bThreeFirst[0]] == bFirstList[FirstAnalyseData.bThreeFirst[0]] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
							return lg.GetCardColor(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较花色
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较数值
						}
					}
				}
			} else {
				return bNextType > bFirstType
			}
		} else if 5 == bFirstCount {
			//开始对比
			if bNextType == bFirstType {
				switch bFirstType {
				case CT_SINGLE: //单牌类型
					{
						if bNextList[0] == bFirstList[0] {
							return false
						}
						bAllSame := true
						for i := 0; i < 5; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								bAllSame = false
								break
							}
						}
						if true == bAllSame {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) //比较花色
						} else {
							for i := 0; i < 5; i++ {
								if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
									return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
								}
							}
							return false
						}
						return false
					}
				case CT_ONE_DOUBLE: //对带一张
					{
						if bNextList[NextAnalyseData.bTwoFirst[0]] == bFirstList[FirstAnalyseData.bTwoFirst[0]] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
							//对比单张
							for i := 0; i < 3; i++ {
								if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[i]]) != lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[i]]) {
									return lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[i]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[i]])
								}
							}
							return lg.GetCardColor(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较花色
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较数值
						}
					}
				case CT_FIVE_TWO_DOUBLE: //两对牌型
					{
						if bNextList[NextAnalyseData.bTwoFirst[0]] == bFirstList[FirstAnalyseData.bTwoFirst[0]] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
							if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
								if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) != lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
									return lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]])
								}
								return lg.GetCardColor(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较花色
							} else {
								return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) //比较数值
							}
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较数值
						}
					}
				case CT_THREE: //三张牌型
					{
						//数据验证
						if bNextList[NextAnalyseData.bThreeFirst[0]] == bFirstList[FirstAnalyseData.bThreeFirst[0]] {
						}
						if bNextList[NextAnalyseData.bThreeFirst[0]] == bFirstList[FirstAnalyseData.bThreeFirst[0]] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
							return lg.GetCardColor(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较花色
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较数值
						}
					}
				case CT_FIVE_MIXED_FLUSH_FIRST_A: //A在前顺子
					{
						if bNextList[0] == bFirstList[0] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[0]) == lg.GetCardLogicValue(bFirstList[0]) {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) //比较花色
						} else {
							return lg.GetCardLogicValue(bNextList[0]) > lg.GetCardLogicValue(bFirstList[0]) //比较数值
						}
					}
				case CT_FIVE_MIXED_FLUSH_NO_A: //没A杂顺
					{
						if bNextList[0] == bFirstList[0] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[0]) == lg.GetCardLogicValue(bFirstList[0]) {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) //比较花色
						} else {
							return lg.GetCardLogicValue(bNextList[0]) > lg.GetCardLogicValue(bFirstList[0]) //比较数值
						}
					}
				case CT_FIVE_MIXED_FLUSH_BACK_A: //A在后顺子
					{
						if bNextList[0] == bFirstList[0] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[0]) == lg.GetCardLogicValue(bFirstList[0]) {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) //比较花色
						} else {
							return lg.GetCardLogicValue(bNextList[0]) > lg.GetCardLogicValue(bFirstList[0]) //比较数值
						}
					}
				case CT_FIVE_FLUSH: //同花五牌
					{
						if bNextList[0] == bFirstList[0] {
							return false
						}
						//比较数值
						for i := 0; i < 5; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						//比较花色
						return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
					}
				case CT_FIVE_THREE_DEOUBLE: //三条一对
					{
						if bNextList[NextAnalyseData.bThreeFirst[0]] == bFirstList[FirstAnalyseData.bThreeFirst[0]] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
							return lg.GetCardColor(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较花色
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较数值
						}
					}
				case CT_FIVE_FOUR_ONE: //四带一张
					{
						if bNextList[NextAnalyseData.bFourFirst[0]] == bFirstList[FirstAnalyseData.bFourFirst[0]] {
							return false
						}
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) {
							return lg.GetCardColor(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bFourFirst[0]]) //比较花色
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) //比较数值
						}
					}
				case CT_FIVE_STRAIGHT_FLUSH_NO_A: //没A同花顺
				case CT_FIVE_STRAIGHT_FLUSH_FIRST_A: //A在前同花顺
				case CT_FIVE_STRAIGHT_FLUSH_BACK_A: //A在后同花顺
					{
						if bNextList[0] == bFirstList[0] {
							return false
						}
						//比较数值
						for i := 0; i < 5; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						//比较花色
						return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
					}
				default:
					return false
				}
			} else {
				return bNextType > bFirstType
			}
		} else {
			if bNextType == bFirstType {
				switch bFirstType {
				case CT_THIRTEEN_FLUSH:
					{
						return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
					}
				case CT_TWELVE_KING:
					{
						AllSame := true
						for i := 0; i < 13; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								AllSame = false
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						if AllSame {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
						}
						return false
					}
				case CT_THREE_STRAIGHTFLUSH:
					{
						AllSame := true
						for i := 0; i < 13; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								AllSame = false
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						if AllSame {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
						}
						return false
					}
				case CT_THREE_BOMB:
					{

						if bNextList[NextAnalyseData.bFourFirst[0]] == bFirstList[FirstAnalyseData.bFourFirst[0]] {
							return false
						}

						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) {
							if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[1]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[1]]) {
								if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[2]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[2]]) {
									if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
										return lg.GetCardColor(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bFourFirst[0]])
									} else {
										return lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]])
									}
								} else {
									return lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[2]])
								}
							} else {
								return lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[1]])
							}
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) //比较数值
						}

						return false
					}
				case CT_ALL_BIG:
					{
						AllSame := true
						for i := 0; i < 13; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								AllSame = false
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						if AllSame {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
						}
						return false
					}
				case CT_ALL_SMALL:
					{

						AllSame := true
						for i := 0; i < 13; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								AllSame = false
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						if AllSame {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
						}
						return false
					}
				case CT_SAME_COLOR:
					{
						AllSame := true
						for i := 0; i < 13; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								AllSame = false
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						if AllSame {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
						}
						return false
					}
				case CT_FOUR_THREESAME:
					{

						if bNextList[NextAnalyseData.bThreeFirst[0]] == bFirstList[FirstAnalyseData.bThreeFirst[0]] {
							return false
						}

						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
							if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[1]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[1]]) {
								if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[2]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[2]]) {
									if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[3]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[3]]) {
										if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
											return lg.GetCardColor(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bThreeFirst[0]])
										} else {
											return lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]])
										}

									} else {
										return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[3]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[3]])
									}
								} else {
									return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[2]])
								}
							} else {
								return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[1]])
							}
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较数值
						}

						return false
					}

				case CT_FIVEPAIR_THREE:
					{

						if bNextList[NextAnalyseData.bTwoFirst[0]] == bFirstList[FirstAnalyseData.bTwoFirst[0]] {
							return false
						}

						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
							if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
								if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]]) {
									if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]]) {
										if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]]) {
											if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
												return lg.GetCardColor(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bTwoFirst[0]])
											} else {
												return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]])
											}
										} else {
											return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]])
										}
									} else {
										return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]])
									}
								} else {
									return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]])
								}
							} else {
								return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]])
							}
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较数值
						}

						return false
					}

				case CT_SIXPAIR:
					{

						if bNextList[NextAnalyseData.bTwoFirst[0]] == bFirstList[FirstAnalyseData.bTwoFirst[0]] {
							return false
						}

						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
							if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
								if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]]) {
									if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]]) {
										if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]]) {
											if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[5]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[5]]) {
												if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
													return lg.GetCardColor(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bTwoFirst[0]])
												} else {
													return lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]])
												}
											} else {
												return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[5]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[5]])
											}
										} else {
											return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[4]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[4]])
										}
									} else {
										return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[3]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[3]])
									}
								} else {
									return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[2]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[2]])
								}
							} else {
								return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]])
							}
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较数值
						}

						return false
					}
				case CT_THREE_FLUSH:
					{
						AllSame := true
						for i := 0; i < 13; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								AllSame = false
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						if AllSame {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
						}
						return false
					}
				case CT_THREE_STRAIGHT:
					{
						AllSame := true
						for i := 0; i < 13; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								AllSame = false
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						if AllSame {
							return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
						}
						return false
					}
				}
			} else {
				return bNextType > bFirstType
			}
		}
	} else {
		//开始对比
		if bNextType == bFirstType {
			switch bFirstType {
			case CT_SINGLE: //单牌类型
				{
					if bNextList[0] == bFirstList[0] {
						return false
					}
					bAllSame := true
					for i := 0; i < 3; i++ {
						if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
							bAllSame = false
							break
						}
					}
					if true == bAllSame {
						return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) //比较花色
					} else {
						for i := 0; i < 3; i++ {
							if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
								return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
							}
						}
						return bNextCount < bFirstCount
					}
					return bNextCount < bFirstCount
				}
			case CT_ONE_DOUBLE: //对带一张
				{
					if bNextList[NextAnalyseData.bTwoFirst[0]] == bFirstList[FirstAnalyseData.bTwoFirst[0]] {
						return false
					}
					if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) != lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]])
						}

						return lg.GetCardColor(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较花色
					} else {
						return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较数值
					}
					return bNextCount < bFirstCount
				}
			case CT_FIVE_TWO_DOUBLE: //两对牌型
				{
					if bNextList[NextAnalyseData.bTwoFirst[0]] == bFirstList[FirstAnalyseData.bTwoFirst[0]] {
						return false
					}
					if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) {
						if lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) {
							//对比单牌
							if lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) != lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]]) {
								return lg.GetCardLogicValue(bNextList[NextAnalyseData.bOneFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bOneFirst[0]])
							}
							return lg.GetCardColor(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较花色
						} else {
							return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[1]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[1]]) //比较数值
						}
					} else {
						return lg.GetCardLogicValue(bNextList[NextAnalyseData.bTwoFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bTwoFirst[0]]) //比较数值
					}
				}
			case CT_THREE: //三张牌型
				{
					if bNextList[NextAnalyseData.bThreeFirst[0]] == bFirstList[FirstAnalyseData.bThreeFirst[0]] {
						return false
					}
					if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
						return lg.GetCardColor(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较花色
					} else {
						return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较数值
					}
					return bNextCount < bFirstCount
				}
			case CT_FIVE_MIXED_FLUSH_FIRST_A: //A在前顺子
				{
					if bNextList[0] == bFirstList[0] {
						return false
					}
					if lg.GetCardLogicValue(bNextList[0]) == lg.GetCardLogicValue(bFirstList[0]) {
						return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) //比较花色
					} else {
						return lg.GetCardLogicValue(bNextList[0]) > lg.GetCardLogicValue(bFirstList[0]) //比较数值
					}
				}
			case CT_FIVE_MIXED_FLUSH_NO_A: //没A杂顺
				{
					if bNextList[0] == bFirstList[0] {
						return false
					}
					if lg.GetCardLogicValue(bNextList[0]) == lg.GetCardLogicValue(bFirstList[0]) {
						return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) //比较花色
					} else {
						return lg.GetCardLogicValue(bNextList[0]) > lg.GetCardLogicValue(bFirstList[0]) //比较数值
					}
				}
			case CT_FIVE_MIXED_FLUSH_BACK_A: //A在后顺子
				{
					if bNextList[0] == bFirstList[0] {
						return false
					}
					if lg.GetCardLogicValue(bNextList[0]) == lg.GetCardLogicValue(bFirstList[0]) {
						return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0]) //比较花色
					} else {
						return lg.GetCardLogicValue(bNextList[0]) > lg.GetCardLogicValue(bFirstList[0]) //比较数值
					}
				}
			case CT_FIVE_FLUSH: //同花五牌
				{
					if bNextList[0] == bFirstList[0] {
						return false
					}
					//比较数值
					for i := 0; i < 5; i++ {
						if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
							return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
						}
					}
					//比较花色
					return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
				}
			case CT_FIVE_THREE_DEOUBLE: //三条一对
				{
					if bNextList[NextAnalyseData.bThreeFirst[0]] == bFirstList[FirstAnalyseData.bThreeFirst[0]] {
						return false
					}
					if lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) {
						return lg.GetCardColor(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较花色
					} else {
						return lg.GetCardLogicValue(bNextList[NextAnalyseData.bThreeFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bThreeFirst[0]]) //比较数值
					}
				}
			case CT_FIVE_FOUR_ONE: //四带一张
				{
					if bNextList[NextAnalyseData.bFourFirst[0]] == bFirstList[FirstAnalyseData.bFourFirst[0]] {
						return false
					}
					if lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) == lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) {
						return lg.GetCardColor(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardColor(bFirstList[FirstAnalyseData.bFourFirst[0]]) //比较花色
					} else {
						return lg.GetCardLogicValue(bNextList[NextAnalyseData.bFourFirst[0]]) > lg.GetCardLogicValue(bFirstList[FirstAnalyseData.bFourFirst[0]]) //比较数值
					}
				}
			case CT_FIVE_STRAIGHT_FLUSH_NO_A: //没A同花顺
			case CT_FIVE_STRAIGHT_FLUSH_FIRST_A: //A在前同花顺
			case CT_FIVE_STRAIGHT_FLUSH_BACK_A: //A在后同花顺
				{
					if bNextList[0] == bFirstList[0] {
						return false
					}
					//比较数值
					for i := 0; i < 5; i++ {
						if lg.GetCardLogicValue(bNextList[i]) != lg.GetCardLogicValue(bFirstList[i]) {
							return lg.GetCardLogicValue(bNextList[i]) > lg.GetCardLogicValue(bFirstList[i])
						}
					}
					//比较花色
					return lg.GetCardColor(bNextList[0]) > lg.GetCardColor(bFirstList[0])
				}
			default:
				return false
			}
		} else {
			return bNextType > bFirstType
		}
	}
	return false
}
