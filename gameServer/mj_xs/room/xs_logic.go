package room

import (
	"mj/common/msg"
	. "mj/gameServer/common/mj"

	"mj/gameServer/common/mj/mj_base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

//扑克转换
func SwitchToCardIndex(cbCardData int) int {
	//计算位置
	cbValue := cbCardData & MASK_VALUE
	cbColor := (cbCardData & MASK_COLOR) >> 4

	if cbColor >= 0x03 {
		return cbValue + 27 - 1
	}
	return cbColor*9 + cbValue - 1
}

//扑克转换
func SwitchToCardData(cbCardIndex int) int {
	if cbCardIndex >= 34 {
		return ((cbCardIndex / 9) << 4) | (cbCardIndex%9 + 1)
	} else {
		return 48 | ((cbCardIndex-34)%8 + 8)
	}
}

//有效判断
func IsValidCard(cbCardData int) bool {
	var cbValue = int(cbCardData & MASK_VALUE)
	var cbColor = int((cbCardData & MASK_COLOR) >> 4)
	return ((cbValue > 0) && (cbValue <= 9) && (cbColor <= 2)) || ((cbValue >= 1) && (cbValue <= 15) && (cbColor == 3))
}

func NewXSlogic(ConfIdx int) *xs_logic {
	l := new(xs_logic)
	l.BaseLogic = mj_base.NewBaseLogic(ConfIdx)
	l.SwitchToIdx = SwitchToCardIndex
	l.SwitchToCard = SwitchToCardData
	l.CheckValid = IsValidCard
	return l
}

type xs_logic struct {
	*mj_base.BaseLogic
}

func (lg *xs_logic) SwitchToCardIndex(cbCardData int) int {
	return lg.SwitchToIdx(cbCardData)
}

//吃牌判断
func (lg *xs_logic) EstimateEatCard(cbCardIndex []int, cbCurrentCard int) int {
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

//吃胡分析,cbCurrentCard 加入这个牌， 手牌数量必须是 3 3 3 2
func (lg *xs_logic) AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbCurrentCard int) (bool, []*TagAnalyseItem) {
	//构造扑克
	cbCardIndexTemp := make([]int, lg.GetCfg().MaxIdx)

	//插入扑克
	if cbCurrentCard != 0 {
		cbCardIndexTemp[lg.SwitchToCardIndex(cbCurrentCard)]++
	}

	//计算数目
	cbCardCount := lg.GetCardCount(cbCardIndex)
	cbWeaveCount := len(WeaveItem)
	//效验数目
	if (cbCardCount < 2) || (cbCardCount > lg.GetCfg().MaxCount) || ((cbCardCount-2)%3 != 0) {
		log.Debug("at AnalyseCard (cbCardCount < 2) || (cbCardCount > MAX_COUNT) || ((cbCardCount-2)mod3 != 0) %v, %v ", cbCardCount, (cbCardCount-2)%3)
		return false, nil
	}

	//需求判断
	TagAnalyseItemArray := make([]*TagAnalyseItem, 0)
	cbLessKindItem := (cbCardCount - 2) / 3
	log.Debug("cbLessKindItem ======= %v, %v ", cbCardCount, cbLessKindItem)
	//单吊判断
	if cbLessKindItem == 0 {
		//牌眼判断
		for i := 0; i < lg.GetCfg().MaxCount; i++ {
			if cbCardIndex[i] == 2 {
				//变量定义
				analyseItem := &TagAnalyseItem{WeaveKind: make([]int, lg.GetCfg().MaxWeave), CenterCard: make([]int, lg.GetCfg().MaxWeave),
					CardData: make([][]int, lg.GetCfg().MaxIdx), IsAnalyseGet: make([]bool, lg.GetCfg().MaxWeave)}
				for i, _ := range analyseItem.CardData {
					analyseItem.CardData[i] = make([]int, 4)
				}

				//设置结果
				for j := 0; j < cbWeaveCount; j++ {
					analyseItem.WeaveKind[j] = WeaveItem[j].WeaveKind
					analyseItem.CenterCard[j] = WeaveItem[j].CenterCard
				}
				analyseItem.CardEye = lg.SwitchToCardData(i)

				//插入结果
				TagAnalyseItemArray = append(TagAnalyseItemArray, analyseItem)
				return true, TagAnalyseItemArray
			}
		}
		return false, nil
	}

	//变量定义
	cbKindItemCount := 0
	KindItem := make([]*TagKindItem, 0)

	if cbCardCount >= 3 {
		for i := 0; i < lg.GetCfg().MaxIdx-lg.GetCfg().HuaIndex; i++ { //不计算花牌
			//同牌判断
			if cbCardIndex[i] >= 3 {
				KindItem = append(KindItem, &TagKindItem{CenterCard: i, CardIndex: []int{i, i, i}, WeaveKind: WIK_PENG})
				cbKindItemCount++
			}

			//连牌判断
			if (i < (lg.GetCfg().MaxIdx - 2 - 15)) && (cbCardIndex[i] > 0) && ((i % 9) < 7) {
				for j := 1; j <= cbCardIndex[i]; j++ {
					if (cbCardIndex[i+1] >= j) && (cbCardIndex[i+2] >= j) {
						KindItem = append(KindItem, &TagKindItem{CenterCard: i, CardIndex: []int{i, i + 1, i + 2}, WeaveKind: WIK_LEFT})
						cbKindItemCount++
					}
				}
			}
		}
	}

	//组合分析
	if cbKindItemCount >= cbLessKindItem {
		//变量定义
		cbCardIndexTemp := make([]int, lg.GetCfg().MaxIdx)
		var cbIndex []int
		Iterator := lg.GetIteratorFunc(cbLessKindItem, cbKindItemCount)
		pKindItem := make([]*TagKindItem, lg.GetCfg().MaxWeave)

		//开始组合
		for {
			cbIndex = Iterator()
			if cbIndex == nil {
				break
			}
			//设置变量
			cbCardIndexTemp = util.CopySlicInt(cbCardIndex)
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

				for i := 0; i < lg.GetCfg().MaxIdx; i++ {
					if cbCardIndexTemp[i] == 2 {
						cbCardEye = lg.SwitchToCard(i)
						break
					}
				}

				//组合类型
				if cbCardEye != 0 {
					//变量定义
					analyseItem := &TagAnalyseItem{WeaveKind: make([]int, lg.GetCfg().MaxWeave), CenterCard: make([]int, lg.GetCfg().MaxWeave), CardData: make([][]int, lg.GetCfg().MaxIdx), IsAnalyseGet: make([]bool, lg.GetCfg().MaxWeave)}
					for i := 0; i < lg.GetCfg().MaxWeave; i++ {
						analyseItem.CardData[i] = make([]int, lg.GetCfg().MaxWeave)
					}
					//设置组合
					for i := 0; i < cbWeaveCount; i++ {
						analyseItem.WeaveKind[i] = WeaveItem[i].WeaveKind
						analyseItem.CenterCard[i] = WeaveItem[i].CenterCard
						lg.GetWeaveCard(WeaveItem[i].WeaveKind, WeaveItem[i].CenterCard, analyseItem.CardData[i])
					}

					//设置牌型
					SetWeaveCount := 0
					if cbWeaveCount > 0 {
						SetWeaveCount = cbWeaveCount - 1
					}
					for i := 0; i < cbLessKindItem; i++ {
						analyseItem.IsAnalyseGet[i+SetWeaveCount] = pKindItem[i].IsAnalyseGet
						analyseItem.WeaveKind[i+SetWeaveCount] = pKindItem[i].WeaveKind
						cbCenterCard := lg.SwitchToCard(pKindItem[i].CenterCard)
						analyseItem.CenterCard[i+SetWeaveCount] = cbCenterCard
						lg.GetWeaveCard(pKindItem[i].WeaveKind, cbCenterCard, analyseItem.CardData[i+SetWeaveCount])
					}

					//设置牌眼
					analyseItem.CardEye = cbCardEye
					//插入结果
					TagAnalyseItemArray = append(TagAnalyseItemArray, analyseItem)
				}
			}
		}
	}

	if len(TagAnalyseItemArray) > 0 {
		log.Debug("hu hu hu hu hu le ")
		return true, TagAnalyseItemArray
	}

	return false, nil

}
