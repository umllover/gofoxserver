package room

import (
	"mj/common/msg"
	. "mj/gameServer/common/mj"
	"mj/gameServer/common/mj/mj_base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

func GetCardWordArray(index []int) bool {
	CardWordArray := []string{
		"一筒", "二筒", "三筒", "四筒", "五筒", "六筒", "七筒", "八筒", "九筒",
		"一万", "二万", "三万", "四万", "五万", "六万", "七万", "八万", "九万",
		"一条", "二条", "三条", "四条", "五条", "六条", "七条", "八条", "九条",
		"东", "南", "西", "北", "中", "发", "白",
		"春", "夏", "秋", "冬", "梅", "兰", "竹", "菊",
	}

	var data string
	//var data2 []int
	for k, v := range index {
		if v > 0 {
			for i := 0; i < v; i++ {
				data = data + "," + CardWordArray[k]
				//data2 = append(data2, k)
			}
		}
	}
	log.Debug("手牌：%s", data)
	//log.Debug("手牌:%d", len(data2))
	return true
}

//func GetCardWord(index int) string {
//	CardWordArray := []string{
//		"一万", "二万", "三万", "四万", "五万", "六万", "七万", "八万", "九万",
//		"一筒", "二筒", "三筒", "四筒", "五筒", "六筒", "七筒", "八筒", "九筒",
//		"一条", "二条", "三条", "四条", "五条", "六条", "七条", "八条", "九条",
//		"东", "南", "西", "北", "中", "发", "白",
//		"春", "夏", "秋", "冬", "梅", "兰", "竹", "菊",
//	}
//	return CardWordArray[index]
//}

type ZP_Logic struct {
	*mj_base.BaseLogic
}

func NewBaseLogic(ConfIdx int) *ZP_Logic {
	bl := new(ZP_Logic)
	bl.BaseLogic = mj_base.NewBaseLogic(ConfIdx)
	bl.BaseLogic.CheckValid = IsValidCard
	bl.BaseLogic.SwitchToCard = SwitchToCardData
	bl.BaseLogic.SwitchToIdx = SwitchToCardIndex
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

//扑克转换
func SwitchToCardData(cbCardIndex int) int {
	if cbCardIndex >= 34 {
		return (4 << 4) | (cbCardIndex - 34 + 1)
	} else {
		return (cbCardIndex/9)<<4 | (cbCardIndex%9 + 1)
	}
}

//有效判断
func IsValidCard(cbCardData int) bool {
	var cbValue = int(cbCardData & MASK_VALUE)
	var cbColor = int((cbCardData & MASK_COLOR) >> 4)
	return ((cbValue >= 1) && (cbValue <= 9) && (cbColor <= 2)) || ((cbValue >= 1) && (cbValue <= 7) && (cbColor == 3) || ((cbValue >= 1) && (cbValue <= 8) && (cbColor == 4)))
}

func (lg *ZP_Logic) SwitchToCardIndex(cbCardData int) int {
	return lg.SwitchToIdx(cbCardData)
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
func (lg *ZP_Logic) AnalyseCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem) (bool, []*TagAnalyseItem) {

	TagAnalyseItemArray := make([]*TagAnalyseItem, 0)
	cbWeaveCount := len(WeaveItem)
	//计算数目
	cbCardCount := lg.GetCardCount(cbCardIndex)

	//效验数目
	if (cbCardCount < 2) || (cbCardCount > lg.GetCfg().MaxCount) || ((cbCardCount-2)%3 != 0) {
		//GetCardWordArray(cbCardIndex) //todo,测试代码
		log.Debug("at AnalyseCard (cbCardCount < 2) || (cbCardCount > room.GetCfg().MaxCount) || ((cbCardCount-2)mod3 != 0) %v,%d %v ", cbCardCount, lg.GetCfg().MaxCount, (cbCardCount-2)%3)
		return false, nil
	}

	//变量定义
	cbKindItemCount := 0
	KindItem := make([]*TagKindItem, 0)

	//需求判断
	cbLessKindItem := (cbCardCount - 2) / 3
	//单吊判断
	if cbLessKindItem == 0 {
		//牌眼判断
		for i := 0; i < lg.GetCfg().MaxIdx; i++ {
			if cbCardIndex[i] == 2 {
				//变量定义
				analyseItem := &TagAnalyseItem{WeaveKind: make([]int, lg.GetCfg().MaxWeave), CenterCard: make([]int, lg.GetCfg().MaxWeave), CardData: make([][]int, lg.GetCfg().MaxIdx),
					IsAnalyseGet: make([]bool, lg.GetCfg().MaxWeave), Param: make([]int, lg.GetCfg().MaxWeave)}
				for i := range analyseItem.CardData {
					analyseItem.CardData[i] = make([]int, 4)
				}

				//设置结果
				for j := 0; j < cbWeaveCount; j++ {
					log.Debug("############### len1:%d len2:%d j:%d", len(analyseItem.WeaveKind), len(WeaveItem), j)
					analyseItem.Param[j] = WeaveItem[j].Param
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
		for i := 0; i < lg.GetCfg().MaxIdx-lg.GetCfg().HuaIndex; i++ { //不计算花牌
			//同牌判断
			if cbCardIndex[i] >= 3 {
				tg := &TagKindItem{CardIndex: make([]int, 4)}
				tg.CenterCard = i
				tg.CardIndex[0] = i
				tg.CardIndex[1] = i
				tg.CardIndex[2] = i
				tg.IsAnalyseGet = true
				tg.WeaveKind = WIK_PENG
				KindItem = append(KindItem, tg)
				cbKindItemCount++
			}

			//连牌判断
			if (i < (lg.GetCfg().MaxIdx - 2 - 15)) && (cbCardIndex[i] > 0) && ((i % 9) < 7) {
				for j := 1; j <= cbCardIndex[i]; j++ {
					if (cbCardIndex[i+1] >= j) && (cbCardIndex[i+2] >= j) {
						tg := &TagKindItem{CardIndex: make([]int, 4)}
						tg.CenterCard = i
						tg.CardIndex[0] = i
						tg.CardIndex[1] = i + 1
						tg.CardIndex[2] = i + 2
						tg.WeaveKind = WIK_LEFT
						tg.IsAnalyseGet = true
						KindItem = append(KindItem, tg)
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
					analyseItem := &TagAnalyseItem{WeaveKind: make([]int, lg.GetCfg().MaxWeave), CenterCard: make([]int, lg.GetCfg().MaxWeave), CardData: make([][]int, lg.GetCfg().MaxIdx),
						IsAnalyseGet: make([]bool, lg.GetCfg().MaxWeave), Param: make([]int, lg.GetCfg().MaxWeave)}
					for i := 0; i < lg.GetCfg().MaxWeave; i++ {
						analyseItem.CardData[i] = make([]int, lg.GetCfg().MaxWeave)
					}
					//设置组合
					for i := 0; i < cbWeaveCount; i++ {
						analyseItem.Param[i] = WeaveItem[i].Param
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
						analyseItem.Param[i+SetWeaveCount] = 0
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

	return true, TagAnalyseItemArray
}

//吃胡分析
func (lg *ZP_Logic) AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbCurrentCard int) (bool, []*TagAnalyseItem) {

	//构造扑克
	cbCardIndexTemp := make([]int, lg.GetCfg().MaxIdx)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)

	//cbCurrentCard一定不为0			!!!!!!!!!
	if cbCurrentCard == 0 {
		return false, nil
	}

	//记录卡牌
	lg.HuOfCard = cbCurrentCard

	//插入扑克
	if cbCurrentCard != 0 {
		cbCardIndexTemp[lg.SwitchToIdx(cbCurrentCard)]++
	}

	//分析扑克
	_, TagAnalyseItemArray := lg.AnalyseCard(cbCardIndexTemp, WeaveItem)

	if len(TagAnalyseItemArray) > 0 {
		return true, TagAnalyseItemArray
	}

	return false, nil
}

//分析听牌
func (lg *ZP_Logic) AnalyseTingCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbOutCardData, cbHuCardCount []int, cbHuCardData [][]int, MaxCount int) int {

	cbOutCount := 0
	cbCardIndexTemp := make([]int, lg.GetCfg().MaxIdx)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)

	cbCardCount := lg.GetCardCount(cbCardIndexTemp)

	if cbOutCardData == nil || cbHuCardCount == nil || cbHuCardData == nil {
		if (cbCardCount+1)%3 == 0 {
			for i := 0; i < lg.GetCfg().MaxIdx-lg.GetCfg().HuaIndex; i++ {
				if cbCardIndexTemp[i] == 0 {
					continue
				}
				cbCardIndexTemp[i]--
				for j := 0; j < lg.GetCfg().MaxIdx-lg.GetCfg().HuaIndex; j++ {
					cbCurrentCard := lg.SwitchToCardData(j)
					hu, _ := lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, cbCurrentCard)
					if hu {
						return WIK_LISTEN
					}
				}
			}
		}
		return cbCardCount
	}

	if (cbCardCount-2)%3 == 0 {
		for i := 0; i < lg.GetCfg().MaxIdx-lg.GetCfg().HuaIndex; i++ {
			if cbCardIndexTemp[i] == 0 {
				continue
			}
			cbCardIndexTemp[i]--

			bAdd := false
			nCount := 0
			for j := 0; j < lg.GetCfg().MaxIdx-lg.GetCfg().HuaIndex; j++ {
				cbCurrentCard := lg.SwitchToCard(j)
				hu, _ := lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, cbCurrentCard)
				if hu {
					if bAdd == false {
						bAdd = true
						cbOutCardData[cbOutCount] = lg.SwitchToCard(i)
						cbOutCount++
					}
					if len(cbHuCardData[cbOutCount-1]) < 1 {
						cbHuCardData[cbOutCount-1] = make([]int, lg.GetCfg().MaxIdx-lg.GetCfg().HuaIndex)
					}
					cbHuCardData[cbOutCount-1][nCount] = lg.SwitchToCard(j)
					nCount++
				}
			}
			if bAdd {
				cbHuCardCount[cbOutCount-1] = nCount
			}

			cbCardIndexTemp[i]++
		}
	} else {
		cbCount := 0
		for j := 0; j < lg.GetCfg().MaxIdx; j++ {
			cbCurrentCard := lg.SwitchToCard(j)
			hu, _ := lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, cbCurrentCard)
			if hu {
				log.Debug("cbCount === %v", cbHuCardData)
				if len(cbHuCardData[0]) < 1 {
					cbHuCardData[0] = make([]int, lg.GetCfg().MaxIdx)
				}

				cbHuCardData[0][cbCount] = cbCurrentCard
				cbCount++
			}
		}
		cbHuCardCount[0] = cbCount
	}

	return cbOutCount
}

//扑克转换
func (lg *ZP_Logic) GetUserCards(cbCardIndex []int) (cbCardData []int) {
	//转换扑克
	for i := 0; i < lg.GetCfg().MaxIdx; i++ {
		if cbCardIndex[i] != 0 {
			for j := 0; j < cbCardIndex[i]; j++ { //牌展开
				cbCardData = append(cbCardData, lg.SwitchToCard(i))
			}
		}
	}
	return cbCardData
}

func (lg *ZP_Logic) GetHuOfCard() int {
	return lg.HuOfCard
}
