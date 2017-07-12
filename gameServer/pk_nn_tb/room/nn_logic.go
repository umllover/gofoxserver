package room

import (
	"mj/gameServer/common/pk/pk_base"
	"github.com/lovelly/leaf/log"
)


// 牛牛类逻辑
const (
	OX_VALUE0  =   0									//混合牌型
    OX_THREE_SAME  =   11                            //三条：有三张相同点数的牌；（3倍）
    OX_ORDER_NUMBER  =   12                           //顺子：五张牌是顺子，最小的顺子12345，最大的为91JQK；（3倍）
    OX_FIVE_SAME_FLOWER  =   13                       //同花：五张牌花色一样；（3倍）
    OX_THREE_SAME_TWAIN  =   14                       //葫芦：三张相同点数的牌+一对；（3倍）
    OX_FOUR_SAME  =   15								//炸弹：有4张相同点数的牌；（4倍）
    OX_STRAIGHT_FLUSH  =   16                          //同花顺：五张牌是顺子且是同一种花色；（4倍）
    OX_FIVE_KING  =   17								//五花：五张牌都是KQJ；（5倍）
    OX_FIVE_CALVES  =   18								//五小牛：5张牌都小于5点且加起来不超过1；（5倍）
	// 牛一到牛牛 ： 1 - 10
	OX_NiuNiu  		= 10
)



func NewNNTBZLogic(ConfigIdx int) *nntb_logic {
	l := new(nntb_logic)
	l.BaseLogic = pk_base.NewBaseLogic(ConfigIdx)
	return l
}

type nntb_logic struct {
	*pk_base.BaseLogic
}



//获取牛牛牌值
func (lg *nntb_logic) GetCardLogicValue(CardData int) int {
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


//牛牛牌型 -------------
// 特殊牌型判断
// 五小 // 五张都小于5 加起来不超过10
func (lg *nntb_logic) IsWuXiao (cardData []int) bool   {
	if len(cardData) != 5 {
		return  false
	}
	sum := 0
	for i:=0; i<5; i++ {
		if lg.GetCardValue(cardData[i])<5 {
			sum += lg.GetCardValue(cardData[i])
		} else {
			return false
		}
	}
	if sum <= 10 {
		return true
	}
	return false
}

// 五花 // 全部是jqk
func (lg *nntb_logic) IsWuHua(cardData []int) bool {
	if len(cardData) != 5 {
		return  false
	}
	for i:=0; i<5; i++ {
		if !(lg.GetCardValue(cardData[i])>10 && lg.GetCardValue(cardData[i])<14) {
			return false
		}
	}
	return  true
}

//  同花顺
func (lg *nntb_logic) IsTongHuaShun(cardData []int) bool {
	if len(cardData) != 5{
		return false
	}
	//排序
	lg.SortCardList(cardData,len(cardData))
	for i:=0; i<4; i++ {
		if lg.GetCardValue(cardData[i])+1 != lg.GetCardValue(cardData[i+1]) ||
			lg.GetCardColor(cardData[i]) != lg.GetCardColor(cardData[i+1]) {
			return false
		}
	}
	return true
}

// 炸弹
func (lg *nntb_logic) IsAllCardValueSame(cardData []int) bool  {
	size := len(cardData)
	for i:=0; i<size-1; i++ {
		if cardData[i] != cardData[i+1] {
			return false
		}
	}
	return  true
}
func (lg *nntb_logic) IsBomb(cardData []int) bool {
	if len(cardData)!= 5 {
		return false
	}
	// 5选4
	for i:=0; i<5; i++ {
		cardDataTemp := make([]int, 4)
		iTemp := 0
		for j:=0; j<5; j++ {
			if j==i {
				continue
			}
			cardDataTemp[iTemp] = cardData[j]
			iTemp++
		}
		// 其中4张点数相同
		if lg.IsAllCardValueSame(cardDataTemp) {
			return true
		}
	}
	return false
}

// 葫芦
func (lg *nntb_logic) IsHuLu(cardData []int) bool {
	if len(cardData)!=5 {
		return false
	}
	//先选两张对子
	for i:=0;i<5;i++ {
		for j:=0;j<5;j++ {
			if j==i {
				continue
			}
			if lg.GetCardValue(cardData[i]) == lg.GetCardValue(cardData[j]) { // 对子
				// 再选三张
				tempCardData := make([]int, 3)
				indexTemp := 0
				for k:=0;k<5;k++ {
					if k==i || k==j {
						continue
					}
					tempCardData[indexTemp] = cardData[k]
					indexTemp++
				}
				if lg.IsAllCardValueSame(tempCardData) {
					return true
				}
			}
		}
	}
	return false
}

// 同花
func (lg *nntb_logic)  IsTongHua(cardData []int) bool {
	if len(cardData) !=5 {
		return false
	}
	for i:=0; i<4; i++ {
		if lg.GetCardColor(cardData[i]) != lg.GetCardColor(cardData[i+1]) {
			return false
		}
	}
	return true
}

// 顺子
func (lg *nntb_logic) IsShunZi(cardData []int) bool {
	if len(cardData)!=5 {
		return false
	}
	lg.SortCardList(cardData, 5)
	for i:=0;i<4;i++ {
		if lg.GetCardValue(cardData[i])+1 != lg.GetCardValue(cardData[i+1]) {
			return false
		}
	}
	return  true
}

// 三条
func (lg *nntb_logic) IsSanTiao(cardData []int) bool {
	for i:=0; i<5; i++ {
		for j:=0; j<5; j++ {
			if j==i {
				continue
			}
			for k:=0; k<5; k++ {
				if k==i || k==j {
					continue
				}
				if lg.GetCardValue(cardData[i]) == lg.GetCardValue(cardData[j]) &&
					lg.GetCardValue(cardData[j]) == lg.GetCardValue(cardData[k]) {
					return true
				}
			}
		}
	}
	return false
}



// 牛牛
func (lg *nntb_logic) IsNiuNiu(cardData []int) bool {
	if len(cardData) != 5{
		return false
	}
	sum := 0
	for i:=0;i<5;i++ {
		sum += lg.GetCardLogicValue(cardData[i])
	}
	if sum%10 == 0 {
		return true
	}
	return false
}


func (lg *nntb_logic) GetCardType(CardData []int) int {

	CardCount := len(CardData)
	if CardCount != lg.GetCfg().MaxCount {
		return  OX_VALUE0
	}
	//特殊牌型判断
	if lg.IsWuXiao(CardData) {
		return OX_FIVE_CALVES
	}

	if lg.IsWuHua(CardData) {
		return OX_FIVE_KING
	}

	if lg.IsTongHuaShun(CardData) {
		return OX_STRAIGHT_FLUSH
	}

	if lg.IsBomb(CardData) {
		return OX_FOUR_SAME
	}

	if lg.IsHuLu(CardData) {
		return OX_THREE_SAME_TWAIN
	}

	if lg.IsTongHua(CardData) {
		return OX_FIVE_SAME_FLOWER
	}

	if lg.IsShunZi(CardData) {
		return OX_ORDER_NUMBER
	}

	if lg.IsSanTiao(CardData) {
		return OX_THREE_SAME
	}

	if lg.IsNiuNiu(CardData) {
		return OX_NiuNiu
	}
	//普通牌型 选3张 有牛
	for i:=0;i<5;i++ {
		for j:=0;j<5;j++ {
			if j==i {
				continue
			}
			for k:=0;k<5;k++ {
				if k==i || k==j {
					continue
				}
				if (lg.GetCardLogicValue(CardData[i]) +
					lg.GetCardLogicValue(CardData[j]) +
					lg.GetCardLogicValue(CardData[k])) % 10 ==0 {
				// 有牛 再选两张
					sum := 0
					for n:=0;n<5;n++ {
						if n==i || n==j || n==k {
							continue
						}
						sum += lg.GetCardLogicValue(CardData[n])
					}
					return sum%10
				}
			}
		}
	}
	return OX_VALUE0
}




//获取牛牛倍数
func (lg *nntb_logic) NNGetTimes(cardData []int, cardCount int, niu int) int {
	if niu != 1 {
		return 1
	}
	if cardCount != lg.GetCfg().MaxCount {
		return 0
	}
	times := lg.GetCardType(cardData)
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
		temp[i] = lg.GetCardLogicValue(cardData[i])
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
		sum += lg.GetCardLogicValue(cardData[i])
	}
	if !(sum > 0) {
		return false
	}
	return (sum%10 == 0)
}

// 牛牛比牌
func (lg *nntb_logic) CompareCard(firstData []int, nextData []int)  bool {

	firstType := lg.GetCardType(firstData)
	nextType := lg.GetCardType(nextData)

	// 先比牌型
	if firstType!= nextType {
		return firstType>nextType
	} else {
		// 牌型一样比点数跟花色 最多只需比到第三张的花色（共同用到两张公共牌）；
		lg.SortCardList(firstData, len(firstData))
		lg.SortCardList(nextData, len(nextData))
		// 两种组合
		if firstType == OX_FOUR_SAME || firstType == OX_THREE_SAME_TWAIN ||
			firstType == OX_THREE_SAME {
			return lg.CompareCardTwoType(firstData, nextData)
		}
		return lg.CompareCardOneType(firstData, nextData)
	}

	return false
}


// 同种牌型比牌
// 一种组合
func (lg *nntb_logic) CompareCardOneType(firstData []int, nextData []int) bool  {
	for i:=0;i<3;i++ {
		if lg.GetCardValue(firstData[i]) != lg.GetCardValue(nextData[i]) {
			return lg.GetCardValue(firstData[i])>lg.GetCardValue(nextData[i])
		} else {
			return lg.GetCardColor(firstData[i])>lg.GetCardValue(nextData[i])
		}
	}
	return false
}
// 两种组合
// 找出大牌点数
func (lg *nntb_logic)FindCardValue(cardData []int) int {
	for i:=0; i<5; i++ {
		for j:=0; j<5; j++ {
			if j==i {
				continue
			}
			for k:=0; k<5; k++ {
				if k==i || k==j {
					continue
				}
				if lg.GetCardValue(cardData[i]) == lg.GetCardValue(cardData[j]) &&
					lg.GetCardValue(cardData[j]) == lg.GetCardValue(cardData[k]) {
					return lg.GetCardValue(cardData[i])
				}
			}
		}
	}
	return 0
}
// 找出一张匹配点数的牌
func (lg *nntb_logic) FindCardWithValue(cardData []int, cardValue int) int {
	for i:=0;i<len(cardData);i++ {
		if lg.GetCardValue(cardData[i]) == cardValue {
			return cardData[i]
		}
	}
	return 0
}

func (lg *nntb_logic) CompareCardTwoType(firstData []int, nextData []int) bool  {
	// 大牌就两种可能 3张或4张 找出点数
	firstValue := lg.FindCardValue(firstData)
	nextValue := lg.FindCardValue(nextData)
	// 比点数
	if firstValue != nextValue {
		return firstValue>nextValue
	}
	// 点数一样比花色
	firstCard := lg.FindCardWithValue(firstData, firstValue)
	nextCard := lg.FindCardWithValue(nextData, nextValue)

	return lg.GetCardColor(firstCard) > lg.GetCardColor(nextCard)
}
