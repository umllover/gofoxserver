package room

import (
	"mj/common/msg"
	"mj/gameServer/common/mj_base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

type ZP_Logic struct {
	*mj_base.BaseLogic
}

func NewBaseLogic() *ZP_Logic {
	bl := new(ZP_Logic)
	bl.BaseLogic = new(mj_base.BaseLogic)
	return bl
}

//扑克转换
func SwitchToCardIndex(cbCardData int) int {
	//计算位置
	cbValue := cbCardData & MASK_VALUE
	cbColor := (cbCardData & MASK_COLOR) >> 4

	if cbColor < 3 {
		return cbColor*9 + cbValue - 1
	} else if cbColor == 3 {
		return 3*9 + cbValue - 1
	} else if cbColor == 4 {
		return 3*9 + 1*7 + cbValue - 1
	} else {
		return 0
	}
}

//有效判断
func (lg *ZP_Logic) IsValidCard(cbCardData int) bool {
	var cbValue = int(cbCardData & MASK_VALUE)
	var cbColor = int((cbCardData & MASK_COLOR) >> 4)
	return ((cbValue >= 1) && (cbValue <= 9) && (cbColor <= 2)) || ((cbValue >= 1) && (cbValue <= 7) && (cbColor == 3) || ((cbValue >= 1) && (cbValue <= 8) && (cbColor == 4)))
}

//吃牌判断
func (lg *ZP_Logic) EstimateEatCard(cbCardIndex []int, cbCurrentCard int) int {
	//番子无连
	if cbCurrentCard >= 0x31 {
		return WIK_NULL
	}

	//变量定义
	cbExcursion := [3]int{0, 1, 2}
	cbItemKind := [3]int{WIK_LEFT, WIK_CENTER, WIK_RIGHT}

	//吃牌判断
	var i int
	var cbEatKind int
	CurrentIndex := SwitchToCardIndex(cbCurrentCard)
	for i = 0; i < len(cbItemKind); i++ {
		cbValueIndex := CurrentIndex % 9
		if cbValueIndex >= cbExcursion[i] && (cbValueIndex-cbExcursion[i]) <= 6 {
			//吃牌判断
			cbFirstIndex := CurrentIndex - cbExcursion[i]
			if CurrentIndex != cbFirstIndex && cbCardIndex[cbFirstIndex] == 0 {
				continue
			}
			if CurrentIndex != (cbFirstIndex+1) && cbCardIndex[cbFirstIndex+1] == 0 {
				continue
			}
			if CurrentIndex != (cbFirstIndex+2) && cbCardIndex[cbFirstIndex+2] == 0 {
				continue
			}
			//设置类型
			cbEatKind |= cbItemKind[i]
		}
	}

	return cbEatKind
}

//分析扑克
func (lg *ZP_Logic) AnalyseCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, TagAnalyseItemArray []*TagAnalyseItem) (bool, []*TagAnalyseItem) {
	cbWeaveCount := len(WeaveItem)
	log.Debug("at AnalyseChiHuCard %v, %v , %v ,%v ", cbCardIndex, WeaveItem, cbWeaveCount, TagAnalyseItemArray)
	//计算数目
	cbCardCount := lg.GetCardCount(cbCardIndex)

	//效验数目
	if (cbCardCount < 2) || (cbCardCount > MAX_COUNT) || ((cbCardCount-2)%3 != 0) {
		log.Debug("at AnalyseCard (cbCardCount < 2) || (cbCardCount > MAX_COUNT) || ((cbCardCount-2)mod3 != 0) %v, %v ", cbCardCount, (cbCardCount-2)%3)
		return false, nil
	}

	//变量定义
	cbKindItemCount := 0
	KindItem := make([]*TagKindItem, MAX_COUNT-2)

	//需求判断
	cbLessKindItem := (cbCardCount - 2) / 3
	log.Debug("cbLessKindItem ======= %v, %v ", cbCardCount, cbLessKindItem)
	//单吊判断
	if cbLessKindItem == 0 {
		//牌眼判断
		for i := 0; i < MAX_INDEX; i++ {
			if cbCardIndex[i] == 2 {
				//变量定义
				analyseItem := &TagAnalyseItem{WeaveKind: make([]int, MAX_WEAVE), CenterCard: make([]int, MAX_WEAVE), CardData: make([][]int, MAX_WEAVE)}
				for i, _ := range analyseItem.CardData {
					analyseItem.CardData[i] = make([]int, 4)
				}

				//设置结果
				for j := 0; j < cbWeaveCount; j++ {
					analyseItem.WeaveKind[j] = WeaveItem[j].WeaveKind
					analyseItem.CenterCard[j] = WeaveItem[j].CenterCard
				}
				analyseItem.CardEye = lg.SwitchToCard(i)

				//插入结果
				TagAnalyseItemArray = append(TagAnalyseItemArray, analyseItem)
				return true, TagAnalyseItemArray
			}
		}
		return false, nil
	}

	if cbCardCount >= 3 {
		for i := 0; i < MAX_INDEX; i++ { //不计算花牌
			//同牌判断
			if cbCardIndex[i] >= 3 {
				KindItem[cbKindItemCount].CenterCard = i
				KindItem[cbKindItemCount].CardIndex[0] = i
				KindItem[cbKindItemCount].CardIndex[1] = i
				KindItem[cbKindItemCount].CardIndex[2] = i
				KindItem[cbKindItemCount].WeaveKind = WIK_PENG
				cbKindItemCount++
			}

			//连牌判断
			if (i < (MAX_INDEX - 2 - 15)) && (cbCardIndex[i] > 0) && ((i % 9) < 7) {
				for j := 1; j <= cbCardIndex[i]; j++ {
					if (cbCardIndex[i+1] >= j) && (cbCardIndex[i+2] >= j) {
						KindItem[cbKindItemCount].CenterCard = i
						KindItem[cbKindItemCount].CardIndex[0] = i
						KindItem[cbKindItemCount].CardIndex[1] = i + 1
						KindItem[cbKindItemCount].CardIndex[2] = i + 2
						KindItem[cbKindItemCount].WeaveKind = WIK_LEFT
						cbKindItemCount++
					}
				}
			}
		}
	}

	//组合分析
	if cbKindItemCount >= cbLessKindItem {
		//变量定义
		cbCardIndexTemp := make([]int, MAX_INDEX)
		cbIndex := []int{0, 1, 2, 3, 4}

		pKindItem := make([]*TagKindItem, MAX_WEAVE)

		//开始组合
		for {
			//设置变量
			util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)
			for i := 0; i < cbLessKindItem; i++ {
				pKindItem[i] = KindItem[cbIndex[i]]
			}

			//数量判断
			bEnoughCard := true

			for i := 0; i < cbLessKindItem*3; i++ {
				//存在判断
				cbCardIndex := pKindItem[i/3].CardIndex[i%3]
				if cbCardIndexTemp[cbCardIndex] == 0 {
					bEnoughCard = false
					break
				} else {
					cbCardIndexTemp[cbCardIndex]--
				}
			}

			//胡牌判断
			if bEnoughCard == true {
				//牌眼判断
				cbCardEye := 0

				for i := 0; i < MAX_INDEX; i++ {
					if cbCardIndexTemp[i] == 2 {
						cbCardEye = lg.SwitchToCard(i)
						break
					}
				}

				//组合类型
				if cbCardEye != 0 {
					//变量定义
					analyseItem := &TagAnalyseItem{WeaveKind: make([]int, MAX_WEAVE), CenterCard: make([]int, MAX_WEAVE), CardData: make([][]int, MAX_WEAVE)}
					for i := 0; i < MAX_WEAVE; i++ {
						analyseItem.CardData[i] = make([]int, MAX_WEAVE)
					}
					//设置组合
					for i := 0; i < cbWeaveCount; i++ {
						analyseItem.WeaveKind[i] = WeaveItem[i].WeaveKind
						analyseItem.CenterCard[i] = WeaveItem[i].CenterCard
						lg.GetWeaveCard(WeaveItem[i].WeaveKind, WeaveItem[i].CenterCard, analyseItem.CardData[i])
					}

					//设置牌型
					for i := 0; i < cbLessKindItem; i++ {
						analyseItem.WeaveKind[i+cbWeaveCount] = pKindItem[i].WeaveKind
						cbCenterCard := lg.SwitchToCard(pKindItem[i].CenterCard)
						analyseItem.CenterCard[i+cbWeaveCount] = cbCenterCard
						lg.GetWeaveCard(pKindItem[i].WeaveKind, cbCenterCard, analyseItem.CardData[i+cbWeaveCount])
					}

					//设置牌眼
					analyseItem.CardEye = cbCardEye
					//插入结果
					TagAnalyseItemArray = append(TagAnalyseItemArray, analyseItem)
				}
			}

			//设置索引
			if cbIndex[cbLessKindItem-1] == (cbKindItemCount - 1) {
				i := cbLessKindItem - 1
				for ; i > 0; i-- {
					if (cbIndex[i-1] + 1) != cbIndex[i] {
						cbNewIndex := cbIndex[i-1]
						for j := (i - 1); j < cbLessKindItem; j++ {
							cbIndex[j] = cbNewIndex + j - i + 2
						}
						break
					}
				}
				if i == 0 {
					break
				}

			} else {
				cbIndex[cbLessKindItem-1]++
			}
		}
	}

	return true, TagAnalyseItemArray
}

//吃牌分析
func (lg *ZP_Logic) AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbCurrentCard int, ChiHuRight int, b4HZHu bool) int {
	//todo ,特殊牌组

	//变量定义
	cbChiHuKind := int(WIK_NULL)
	TagAnalyseItemArray := make([]*TagAnalyseItem, 0) //

	//构造扑克
	cbCardIndexTemp := make([]int, MAX_INDEX)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)

	//cbCurrentCard一定不为0			!!!!!!!!!
	if cbCurrentCard == 0 {
		return WIK_NULL
	}

	//插入扑克
	if cbCurrentCard != 0 {
		cbCardIndexTemp[lg.SwitchToIdx(cbCurrentCard)]++
	}

	if b4HZHu && cbCardIndexTemp[31] == 4 { //四个红中直接胡牌
		return WIK_CHI_HU
	}
	//分析扑克
	_, TagAnalyseItemArray = lg.AnalyseCard(cbCardIndexTemp, WeaveItem, TagAnalyseItemArray)

	//胡牌分析
	if len(TagAnalyseItemArray) > 0 {
		log.Debug("len(TagAnalyseItemArray) > 0 ")
		ChiHuRight |= CHR_PING_HU
	}

	if ChiHuRight != 0 {
		log.Debug("ChiHuRight != 0 ")
		cbChiHuKind = WIK_CHI_HU
	}

	return cbChiHuKind
}
