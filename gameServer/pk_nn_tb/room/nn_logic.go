package room

import (
	"mj/gameServer/common/pk/pk_base"
	"github.com/lovelly/leaf/log"
)


// 牛牛类逻辑
const (
	OX_VALUE0  =   0									//混合牌型
    OX_THREE_SAME  =   101                            //三条：有三张相同点数的牌；（3倍）
    OX_ORDER_NUMBER  =   102                           //顺子：五张牌是顺子，最小的顺子12345，最大的为910JQK；（3倍）
    OX_FIVE_SAME_FLOWER  =   103                       //同花：五张牌花色一样；（3倍）
    OX_THREE_SAME_TWAIN  =   104                       //葫芦：三张相同点数的牌+一对；（3倍）
    OX_FOUR_SAME  =   105								//炸弹：有4张相同点数的牌；（4倍）
    OX_STRAIGHT_FLUSH  =   106                          //同花顺：五张牌是顺子且是同一种花色；（4倍）
    OX_FIVE_KING  =   107								//五花：五张牌都是KQJ；（5倍）
    OX_FIVE_CALVES  =   108								//五小牛：5张牌都小于5点且加起来不超过10；（5倍）
	// 牛一到牛牛 ： 1 - 10
)



func NewNNTBZLogic(ConfigIdx int) *nntb_logic {
	l := new(nntb_logic)
	l.BaseLogic = pk_base.NewBaseLogic(ConfigIdx)
	return l
}

type nntb_logic struct {
	*pk_base.BaseLogic
}

func (lg *nntb_logic) CompareCard(firstCardData []int, lastCardData []int) bool  {

	return  false
}



//获取牛牛牌值
func (lg *nntb_logic) NNGetCardLogicValue(CardData int) int {
	//扑克属性
	//CardColor = GetCardColor(CardData)
	CardValue := lg.GetCardValue(CardData)

	//转换数值
	//return (CardValue>10)?(10):CardValue
	if CardValue > 10 {
		CardValue = 10
	}
	return CardValue
}


//获取牛牛牌型
func (lg *nntb_logic) NNGetCardType(CardData []int, CardCount int) int {

	if CardCount != lg.GetCfg().MaxCount {
		return 0
	}

	////炸弹牌型
	//SameCount := 0

	Temp := make([]int, lg.GetCfg().MaxCount)
	Sum := 0
	for i := 0; i < CardCount; i++ {
		Temp[i] = lg.NNGetCardLogicValue(CardData[i])
		log.Debug("%d", Temp[i])
		Sum += Temp[i]
	}
	log.Debug("%d", Sum)

	//王的数量
	KingCount := 0
	TenCount := 0

	for i := 0; i < CardCount; i++ {
		if lg.GetCardValue(CardData[i]) > 10 && CardData[i] != 0x4E && CardData[i] != 0x4F {
			KingCount++
		} else if lg.GetCardValue(CardData[i]) == 10 {
			TenCount++
		}
	}

	if KingCount == lg.GetCfg().MaxCount {
		return OX_FIVE_KING   //五花――5张牌都是10以上（不含10）的牌。。
	}

	Value := lg.NNGetCardLogicValue(CardData[3])
	Value += lg.NNGetCardLogicValue(CardData[4])

	if Value > 10 {
		if CardData[3] == 0x4E || CardData[4] == 0x4F || CardData[4] == 0x4E || CardData[3] == 0x4F {
			Value = 10
		} else {
			Value -= 10 //2.3
		}

	}

	return Value //OX_VALUE0
}

//获取牛牛倍数
func (lg *nntb_logic) NNGetTimes(cardData []int, cardCount int, niu int) int {
	if niu != 1 {
		return 1
	}
	if cardCount != lg.GetCfg().MaxCount {
		return 0
	}
	times := lg.NNGetCardType(cardData, lg.GetCfg().MaxCount)
	log.Debug("times %d", times)

	/*if(bTimes<7)return 1;
	else if(bTimes==7)return 1;
	else if(bTimes==8)return 2;
	else if(bTimes==9)return 3;
	else if(bTimes==10)return 4;*/
	//else if(bTimes==OX_THREE_SAME)return 5;
	//else if(bTimes==OX_FOUR_SAME)return 5;
	//else if(bTimes==OX_FOURKING)return 5;
	//else if(bTimes==OX_FIVEKING)return 5;

	if times < 7 {
		return 1
	} else if times >= 7 && times <= 10 {
		return times - 6
	} else if times == OX_FIVE_KING {
		return 5
	}
	return 0
}

// 获取牛牛
func (lg *nntb_logic) NNGetOxCard(cardData []int, cardCount int) bool {
	if cardCount != lg.GetCfg().MaxCount {
		return false
	}

	temp := make([]int, lg.GetCfg().MaxCount)
	sum := 0
	for i := 0; i < lg.GetCfg().MaxCount; i++ {
		temp[i] = lg.NNGetCardLogicValue(cardData[i])
		sum += temp[i]
	}
	//王的数量
	kingCount := 0
	tenCount := 0

	for i := 0; i < lg.GetCfg().MaxCount; i++ {
		if cardData[i] == 0x4E || cardData[i] == 0x4F {
			kingCount++
		} else if lg.GetCardValue(cardData[i]) == 10 {
			tenCount++
		}
	}
	maxNiuZi := 0
	maxNiuPos := 0
	niuTemp := make([][]int, 30,lg.GetCfg().MaxCount)
	var isKingPai [30]bool

	niuCount := 0
	haveKing := false
	//查找牛牛
	for i := 0; i < cardCount-1; i++ {
		for j := 0; j < cardCount; j++ {
			haveKing = false
			left := (sum - temp[i] - temp[j]) % 10
			if left > 0 && kingCount > 0 {
				for k := 0; k < cardCount; k++ {
					if k != i && k != j {
						if cardData[k] == 0x4E || cardData[k] == 0x4F {
							haveKing = true
						}
					}
				}
			}
			if (sum-temp[i]-temp[j])%10 == 0 || haveKing { ////如果减去2个剩下3个是10的倍数
				count := 0
				for k := 0; k < cardCount; k++ {
					if k != i && k != j {
						niuTemp[niuCount][count] = cardData[k]
						count++
					}
				}
				if count != 3 {
					log.Debug("NNGetOxCard err not 3")
					return false
				}
				niuTemp[niuCount][count] = cardData[i]
				count++
				niuTemp[niuCount][count] = cardData[j]
				count++
				value := temp[i]
				value += temp[j]
				if value > 10 {
					if cardData[i] == 0x4E || cardData[j] == 0x4F || cardData[i] == 0x4F || cardData[j] == 0x4E {
						haveKing = true
						value = 10
					} else {
						value -= 10
					}
				}
				isKingPai[niuCount] = haveKing
				if value > maxNiuZi {
					maxNiuZi = value     //最大牛数量
					maxNiuPos = niuCount //记录最大牛牌的位置
				}
				niuCount++
				continue
			}
		}
	}
	if niuCount > 0 {
		for i := 0; i < cardCount; i++ {
			cardData[i] = niuTemp[maxNiuPos][i]
		}
		return true
	}
	return false
}

// 牛牛获取整数
func (lg *nntb_logic) NNIsIntValue(cardData []int, cardCount int) bool {
	sum := 0
	for i := 0; i < cardCount; i++ {
		sum += lg.NNGetCardLogicValue(cardData[i])
	}
	if !(sum > 0) {
		return false
	}
	return (sum%10 == 0)
}

// 牛牛比牌
func (lg *nntb_logic) NNCompareCard(firstData []int, nextData []int)  bool {
	/*
	if firstOX != nextOX {
		if firstOX {
			return true
		} else {
			return false
		}
	}
	if lg.NNGetCardType(firstData, cardCount) == OX_FIVE_KING && lg.NNGetCardType(nextData, cardCount) != OX_FIVE_KING {
		return true
	}
	if lg.NNGetCardType(firstData, cardCount) != OX_FIVE_KING && lg.NNGetCardType(nextData, cardCount) == OX_FIVE_KING {
		return false
	}
	//比较牛大小
	if firstOX == true {
		//获取点数
		firstType := 0
		nextType := 0

		value := lg.NNGetCardLogicValue(nextData[3])
		value += lg.NNGetCardLogicValue(nextData[4])

		firstKing := false
		nextKing := false

		firstDa := false
		nextDa := false //nextDa是判断4,5有没有利用大王的

		if value > 10 {
			if nextData[3] == 0x4E || nextData[4] == 0x4F || nextData[4] == 0x4E || nextData[3] == 0x4F {
				left := 0
				value = 0
				for i := 3; i < 5; i++ {
					value += lg.NNGetCardLogicValue(nextData[i])
				}
				left = value % 10
				if left > 0 {
					nextDa = true
				}
				value = 10
			} else {
				value -= 10
			}
		}
		nextType = value
		kingCount := 0
		for i := 0; i < 3; i++ {
			if nextData[i] == 0x4E || nextData[i] == 0x4F {
				kingCount++
			}
		}
		if kingCount > 0 {
			value = 0
			left := 0
			for i := 0; i < 3; i++ {
				value += lg.NNGetCardLogicValue(nextData[i])
			}
			left = value % 10
			if left > 10 {
				nextKing = true
			}
		}
		value = 0
		value = lg.NNGetCardLogicValue(firstData[3])
		value += lg.NNGetCardLogicValue(firstData[4])
		if value > 10 {
			if firstData[3] == 0x4E || firstData[4] == 0x4F || firstData[4] == 0x4E || firstData[3] == 0x4F {
				left := 0
				value = 0
				for i := 3; i < 5; i++ {
					value += lg.NNGetCardLogicValue(firstData[i])
				}
				left = value % 10
				if left > 0 {
					firstDa = true
				}
				value = 10
			} else {
				value -= 10
			}
		}
		firstType = value
		kingCount = 0
		for i := 0; i < 3; i++ {
			if firstData[i] == 0x4E || firstData[i] == 0x4F {
				kingCount++
			}
		}
		if kingCount > 0 {
			value = 0
			left := 0
			for i := 0; i < 3; i++ {
				value += lg.NNGetCardLogicValue(firstData[i])
			}
			left = value % 10
			if left > 0 {
				firstKing = true
			}
		}
		if firstType == nextType {
			//同点数大王>小王>...
			firstKingPoint := 10
			nextKingPoint := 10
			for i := 0; i < 5; i++ {
				if firstData[i] == 0x4E {
					firstKingPoint = 11
				} else if firstData[i] == 0x4F {
					firstKingPoint = 12
				}
				if nextData[i] == 0x4E {
					nextKingPoint = 11
				} else if nextData[i] == 0x4F {
					nextKingPoint = 12
				}
			}
			if firstKingPoint != nextKingPoint {
				return (firstKingPoint > nextKingPoint)
			}
			if firstKing || firstDa {
				return true
			} else if nextKing || nextDa {
				return false
			}
		}
		//点数判断
		if firstType != nextType {
			return (firstType > nextType)
		}
	}
	//排序大小
	var firstTemp []int
	var nextTemp []int
	util.DeepCopy(firstTemp, firstData)
	util.DeepCopy(nextTemp, nextData)
	lg.SortCardList(firstTemp, cardCount)
	lg.SortCardList(nextTemp, cardCount)
	//比较数值
	nextMaxValue := lg.GetCardValue(nextTemp[0])
	firstMaxValue := lg.GetCardValue(firstTemp[0])
	if nextMaxValue != firstMaxValue {
		return (firstMaxValue > nextMaxValue)
	}
	//比较颜色
	return (lg.GetCardColor(firstTemp[0]) > lg.GetCardColor(nextTemp[0]))
*/
	return false
}
