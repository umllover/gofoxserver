package mj_base

import (
	"mj/common/msg"
	"mj/common/utils"
	. "mj/gameServer/common/mj"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

func IsValidCard(cbCardData int) bool {
	var cbValue = int(cbCardData & MASK_VALUE)
	var cbColor = int((cbCardData & MASK_COLOR) >> 4)
	return ((cbValue >= 1) && (cbValue <= 9) && (cbColor <= 2)) || ((cbValue >= 1) && (cbValue <= (7 + GetCfg(IDX_HZMJ).HuaIndex)) && (cbColor == 3))
}

//扑克转换
func SwitchToCardData(cbCardIndex int) int {
	if cbCardIndex < 34 { //花三种花色牌 3 * 9
		return ((cbCardIndex / 9) << 4) | (cbCardIndex%9 + 1)
	}
	return 48 | ((cbCardIndex-34)%8 + 8)
}

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

type mj_logic interface {
}

type BaseLogic struct {
	CardDataArray []int //扑克数据
	MagicIndex    int   //钻牌索引
	ReplaceCard   int   //替换金牌的牌
	SwitchToIdx   func(int) int
	CheckValid    func(int) bool
	SwitchToCard  func(int) int
	ConfigIdx     int //配置文件索引
}

func NewBaseLogic(ConfIdx int) *BaseLogic {
	bl := new(BaseLogic)
	bl.CheckValid = IsValidCard
	bl.SwitchToIdx = SwitchToCardIndex
	bl.SwitchToCard = SwitchToCardData
	bl.ConfigIdx = ConfIdx
	return bl
}

func (lg *BaseLogic) GetCfg() *MJ_CFG {
	return GetCfg(lg.ConfigIdx)
}

func (lg *BaseLogic) SwitchToCardData(cbCardIndex int) int {
	return lg.SwitchToCard(cbCardIndex)
}
func (lg *BaseLogic) SwitchToCardIndex(cbCardData int) int {
	return lg.SwitchToIdx(cbCardData)
}

func (lg *BaseLogic) GetMagicIndex() int {
	return lg.MagicIndex
}

func (lg *BaseLogic) SetMagicIndex(idx int) {
	lg.MagicIndex = idx
}

func (lg *BaseLogic) IsValidCard(card int) bool {
	return lg.CheckValid(card)
}

//混乱扑克
func (lg *BaseLogic) RandCardList(cbCardBuffer, OriDataArray []int) {
	//混乱准备
	cbBufferCount := int(len(cbCardBuffer))
	cbCardDataTemp := util.CopySlicInt(OriDataArray)
	//混乱扑克
	var cbRandCount int
	var cbPosition int
	for {
		if cbRandCount >= cbBufferCount {
			break
		}
		cbPosition, _ = utils.RandInt(0, cbBufferCount-cbRandCount)
		cbCardBuffer[cbRandCount] = cbCardDataTemp[cbPosition]
		cbRandCount++
		cbCardDataTemp[cbPosition] = cbCardDataTemp[cbBufferCount-cbRandCount]
	}

	log.Debug("rand card s ==: ", cbCardBuffer)
	return
}

//删除扑克
func (lg *BaseLogic) RemoveCardByArr(cbCardIndex, cbRemoveCard []int) bool {
	log.Debug("删除卡牌：%v", cbRemoveCard)
	//参数校验
	for _, card := range cbRemoveCard {
		//效验扑克
		if lg.CheckValid(card) == false {
			log.Debug("效验扑克CheckValid")
			return false
		}

		if cbCardIndex[lg.SwitchToIdx(card)] <= 0 {
			log.Debug("效验扑克cbCardIndex[lg.SwitchToIdx(card)] <= 0")
			return false
		}
	}
	//删除扑克
	for _, card := range cbRemoveCard {
		//删除扑克
		cbCardIndex[lg.SwitchToIdx(card)]--
	}

	return true
}

//删除扑克
func (lg *BaseLogic) RemoveCard(cbCardIndex []int, cbRemoveCard int) bool {
	log.Debug("用户卡牌数据：%v", cbCardIndex)
	//删除扑克
	cbRemoveIndex := lg.SwitchToIdx(cbRemoveCard)
	//效验扑克
	if !lg.CheckValid(cbRemoveCard) {
		log.Error("at RemoveCard card is Invalid %d", cbRemoveCard)
	}
	if cbCardIndex[lg.SwitchToIdx(cbRemoveCard)] < 0 {
		log.Error("at RemoveCard 11 card is Invalid %d", cbRemoveCard)
	}
	if cbCardIndex[cbRemoveIndex] > 0 {
		cbCardIndex[cbRemoveIndex]--
		return true
	}
	log.Debug("删除扑克用户卡牌数据：%v", cbCardIndex)
	return false
}

//扑克数目
func (lg *BaseLogic) GetCardCount(cbCardIndex []int) int {
	//数目统计
	cbCardCount := 0
	for i := 0; i < lg.GetCfg().MaxIdx; i++ {
		cbCardCount += cbCardIndex[i]
	}
	return cbCardCount
}

//获取组合
func (lg *BaseLogic) GetWeaveCard(cbWeaveKind, cbCenterCard int, cbCardBuffer []int) int {
	//组合扑克
	switch cbWeaveKind {
	case WIK_LEFT: //上牌操作
		//设置变量
		cbCardBuffer[0] = cbCenterCard
		cbCardBuffer[1] = cbCenterCard + 1
		cbCardBuffer[2] = cbCenterCard + 2
		return 3

	case WIK_RIGHT: //上牌操作
		//设置变量
		cbCardBuffer[0] = cbCenterCard - 2
		cbCardBuffer[1] = cbCenterCard - 1
		cbCardBuffer[2] = cbCenterCard
		return 3

	case WIK_CENTER: //上牌操作
		//设置变量
		cbCardBuffer[0] = cbCenterCard - 1
		cbCardBuffer[1] = cbCenterCard
		cbCardBuffer[2] = cbCenterCard + 1
		return 3

	case WIK_PENG: //碰牌操作
		//设置变量
		cbCardBuffer[0] = cbCenterCard
		cbCardBuffer[1] = cbCenterCard
		cbCardBuffer[2] = cbCenterCard
		return 3

	case WIK_GANG: //杠牌操作
		//设置变量
		cbCardBuffer[0] = cbCenterCard
		cbCardBuffer[1] = cbCenterCard
		cbCardBuffer[2] = cbCenterCard
		cbCardBuffer[3] = cbCenterCard
		return 4

	default:
	}

	return 0
}

//动作等级
func (lg *BaseLogic) GetUserActionRank(cbUserAction int) int {
	//胡牌等级
	if cbUserAction&WIK_CHI_HU != 0 {
		return 4
	}

	//杠牌等级
	if cbUserAction&WIK_GANG != 0 {
		return 3
	}

	//碰牌等级
	if cbUserAction&WIK_PENG != 0 {
		return 2
	}

	//上牌等级
	if cbUserAction&(WIK_RIGHT|WIK_CENTER|WIK_LEFT) != 0 {
		return 1
	}

	return 0
}

//碰牌判断
func (lg *BaseLogic) EstimatePengCard(cbCardIndex []int, cbCurrentCard int) int {
	log.Debug("at EstimatePengCard  cnt:%v, card:%v, allcard:%v", cbCardIndex[lg.SwitchToIdx(cbCurrentCard)], cbCurrentCard, cbCardIndex)
	if cbCardIndex[lg.SwitchToIdx(cbCurrentCard)] >= 2 {
		return WIK_PENG
	}

	return WIK_NULL
}

//杠牌判断
func (lg *BaseLogic) EstimateGangCard(cbCardIndex []int, cbCurrentCard int) int {
	log.Debug("at EstimateGangCard  cnt:%v, card:%v, allcard:%v", cbCardIndex[lg.SwitchToIdx(cbCurrentCard)], cbCurrentCard, cbCardIndex)
	if cbCardIndex[lg.SwitchToIdx(cbCurrentCard)] == 3 {
		return WIK_GANG
	}

	return WIK_NULL
}

func (lg *BaseLogic) GetCardColor(cbCardData int) int { return cbCardData & MASK_COLOR }
func (lg *BaseLogic) GetCardValue(cbCardData int) int { return cbCardData & MASK_VALUE }

//吃胡分析)
func (lg *BaseLogic) AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbCurrentCard int) (bool, []*TagAnalyseItem) {

	//构造扑克
	cbCardIndexTemp := make([]int, lg.GetCfg().MaxIdx)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)

	//cbCurrentCard一定不为0			!!!!!!!!!
	if cbCurrentCard == 0 {
		return false, nil
	}

	//插入扑克
	cbCardIndexTemp[lg.SwitchToIdx(cbCurrentCard)]++
	if lg.ConfigIdx == IDX_HZMJ && cbCardIndexTemp[31] == 4 { //四个红中直接胡牌
		return true, nil
	}
	//分析扑克
	_, TagAnalyseItemArray := lg.AnalyseCard(cbCardIndexTemp, WeaveItem)

	//胡牌分析
	if len(TagAnalyseItemArray) > 0 {
		log.Debug("hu hu hu hu hu le ")
		return true, TagAnalyseItemArray
	}

	return false, nil
}

func (lg *BaseLogic) AnalyseGangCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbProvideCard int, gangCardResult *TagGangCardResult) int {
	//设置变量
	cbActionMask := WIK_NULL
	cbWeaveCount := len(WeaveItem)
	gangCardResult.CardData = make([]int, lg.GetCfg().MaxWeave)
	//手上杠牌
	for i := 0; i < lg.GetCfg().MaxIdx; i++ {
		if i == lg.MagicIndex {
			continue
		}
		if cbCardIndex[i] == 4 {
			cbActionMask |= WIK_GANG
			gangCardResult.CardData[gangCardResult.CardCount] = lg.SwitchToCard(i)
			gangCardResult.CardCount++
		}
	}

	//组合杠牌
	for i := 0; i < cbWeaveCount; i++ {
		if WeaveItem[i].WeaveKind == WIK_PENG {
			if WeaveItem[i].CenterCard == cbProvideCard { //之后抓来的的牌才能和碰组成杠
				cbActionMask |= WIK_GANG
				gangCardResult.CardData[gangCardResult.CardCount] = WeaveItem[i].CenterCard
				gangCardResult.CardCount++
			}
			for j := 0; j < lg.GetCfg().MaxIdx; j++ { //碰来后的和手牌组成杠
				if cbCardIndex[j] > 0 {
					card := lg.SwitchToCard(j)
					if WeaveItem[i].CenterCard == card {
						cbActionMask |= WIK_GANG
						gangCardResult.CardData[gangCardResult.CardCount] = WeaveItem[i].CenterCard
						gangCardResult.CardCount++
						break
					}
				}
			}
		}
	}

	return cbActionMask
}

func (lg *BaseLogic) GetHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbHuCardData []int, MaxCount int) int {
	cbCardIndexTemp := make([]int, lg.GetCfg().MaxIdx)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)
	cbHuCardData = make([]int, lg.GetCfg().MaxIdx-lg.GetCfg().HuaIndex)

	count := 0
	cardCount := lg.GetCardCount(cbCardIndexTemp)
	if (cardCount-2)%3 != 0 {
		for i := 0; i < lg.GetCfg().MaxIdx-lg.GetCfg().HuaIndex; i++ {
			CurrentCard := lg.SwitchToCardData(i)
			hu, _ := lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, CurrentCard)
			if hu {
				cbHuCardData[count] = CurrentCard
				count++
			}
		}
	}
	if count > 0 {
		return count
	}

	return 0
}

func (lg *BaseLogic) AnalyseTingCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbOutCardData, cbHuCardCount []int, cbHuCardData [][]int) int {
	cbOutCount := 0
	cbCardIndexTemp := util.CopySlicInt(cbCardIndex)
	cbCardCount := lg.GetCardCount(cbCardIndexTemp)

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

//分析扑克
func (lg *BaseLogic) AnalyseCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem) (bool, []*TagAnalyseItem) {
	TagAnalyseItemArray := make([]*TagAnalyseItem, 0)
	cbWeaveCount := len(WeaveItem)
	//计算数目
	cbCardCount := lg.GetCardCount(cbCardIndex)

	//效验数目
	if (cbCardCount < 2) || (cbCardCount > lg.GetCfg().MaxCount) || ((cbCardCount-2)%3 != 0) {
		log.Debug("at AnalyseCard (cbCardCount < 2) || (cbCardCount > room.GetCfg().MaxCount) || ((cbCardCount-2)mod3 != 0) %v, %v ", cbCardCount, (cbCardCount-2)%3)
		return false, nil
	}

	//变量定义
	cbKindItemCount := 0
	KindItem := make([]*TagKindItem, 0)
	//万能牌
	cbMagicCount := 0
	cbMagicIndex := lg.GetMagicIndex()
	if cbMagicIndex != 0 {
		cbMagicCount = cbCardIndex[cbMagicIndex]
	}
	//需求判断
	cbLessKindItem := (cbCardCount - 2) / 3
	//单吊判断
	if cbLessKindItem == 0 {
		//牌眼判断
		for i := 0; i < lg.GetCfg().MaxIdx; i++ {
			if cbCardIndex[i] == 2 || (cbMagicIndex != 0 && i != cbMagicIndex && cbMagicCount+cbCardIndex[i] == 2) {
				//变量定义
				analyseItem := &TagAnalyseItem{WeaveKind: make([]int, lg.GetCfg().MaxWeave), CenterCard: make([]int, lg.GetCfg().MaxWeave), CardData: make([][]int, lg.GetCfg().MaxIdx), IsAnalyseGet: make([]bool, lg.GetCfg().MaxWeave)}
				for i := range analyseItem.CardData {
					analyseItem.CardData[i] = make([]int, lg.GetCfg().MaxWeave)
				}
				//设置结果
				for j := 0; j < cbWeaveCount; j++ {
					analyseItem.WeaveKind[j] = WeaveItem[j].WeaveKind
					analyseItem.CenterCard[j] = WeaveItem[j].CenterCard
				}
				analyseItem.MagicEye = cbCardIndex[i] < 2 || i == cbMagicIndex
				if cbCardIndex[i] == 0 {
					analyseItem.CardEye = lg.SwitchToCard(cbMagicIndex)
				} else {
					analyseItem.CardEye = lg.SwitchToCard(i)
				}
				//插入结果
				TagAnalyseItemArray = append(TagAnalyseItemArray, analyseItem)
				return true, TagAnalyseItemArray
			}
		}
		return false, nil
	}

	//多牌判断
	if cbCardCount >= 3 {
		for i := 0; i < lg.GetCfg().MaxIdx; i++ { //不计算花牌
			//同牌判断
			if cbCardIndex[i] >= 3 || (cbCardIndex[i]+cbMagicCount >= 3 && i != cbMagicIndex) {
				nTempCount := cbCardIndex[i]
				for {
					KindItem = lg.AddKindItem(KindItem, []int{i}, cbMagicIndex, nTempCount, WIK_PENG)
					cbKindItemCount++
					//当前索引牌未与万能牌组合且万能牌个数不为0
					if nTempCount >= 3 && cbMagicCount > 0 {
						nTempCount--
						//组合个数
						nRange := 1
						if cbMagicCount > 1 {
							nRange = 2
						}
						for t := 1; t <= nRange; t++ {
							KindItem = lg.AddKindItem(KindItem, []int{i}, cbMagicIndex, nTempCount, WIK_PENG)
							cbKindItemCount++
						}
						nTempCount++
					}
					if i == cbMagicIndex {
						break
					}
					nTempCount -= 3
					//如果刚好搭配全部，则退出，或者数量不足
					if nTempCount == 0 || nTempCount+cbMagicCount < 3 {
						break
					}
				}
			}
			//连牌判断
			if (i < (lg.GetCfg().MaxIdx - 2 - 15)) && ((i % 9) < 7) {
				if cbMagicCount+cbCardIndex[i]+cbCardIndex[i+1]+cbCardIndex[i+2] >= 3 {
					cbIndex := []int{cbCardIndex[i], cbCardIndex[i+1], cbCardIndex[i+2]}
					if cbIndex[0]+cbIndex[1]+cbIndex[2] == 0 {
						continue
					}
					nTempMagicCount := cbMagicCount
					for {
						if nTempMagicCount+cbIndex[0]+cbIndex[1]+cbIndex[2] < 3 {
							break
						}
						tempArray := [3]int{}
						tempCount := 0
						for j := 0; j <= 2; j++ {
							if cbIndex[j] > 0 {
								cbIndex[j]--
								tempArray[j] = i + j
							} else {
								nTempMagicCount--
								tempArray[j] = cbMagicIndex
								tempCount += 1
							}
						}
						if nTempMagicCount >= 0 {
							KindItem = lg.AddKindItem(KindItem, []int{i, tempArray[0], tempArray[1], tempArray[2]}, cbMagicIndex, tempCount, WIK_LEFT)
							cbKindItemCount++
						} else {
							break
						}
					}
				}
			}
		}
	}

	//log.Debug("-------cbKindItemCount=%d, cbLessKindItem=%d, KindItem=%d, cbMagicCount=%d", cbKindItemCount, cbLessKindItem, len(KindItem), cbMagicCount)
	//for _, tg := range KindItem {
	//	log.Debug("KindItem====%v %v %v %v", tg.CenterCard, tg.MagicCount, tg.WeaveKind, tg.CardIndex)
	//}
	//for _, wm := range WeaveItem {
	//	log.Debug("WeaveItem====%v %v %v %v %v %v", wm.PublicCard, wm.Param, wm.ActionMask, wm.CenterCard, wm.CardData, wm.WeaveKind)
	//}

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
			cbTempMagicCount := 0
			for i := 0; i < cbLessKindItem; i++ {
				cbTempMagicCount += KindItem[cbIndex[i]].MagicCount
			}
			if cbTempMagicCount <= cbMagicCount {
				cbCardIndexTemp = util.CopySlicInt(cbCardIndex)
				for i := 0; i < cbLessKindItem; i++ {
					pKindItem[i] = KindItem[cbIndex[i]]
				}
				bEnoughCard := true
				for i := 0; i < cbLessKindItem*3; i++ {
					cbCardIndex := pKindItem[i/3].CardIndex[i%3]
					if cbCardIndexTemp[cbCardIndex] == 0 {
						if cbMagicIndex != 0 && cbCardIndexTemp[cbMagicIndex] > 0 {
							pKindItem[i/3].CardIndex[i%3] = cbMagicIndex
							cbCardIndexTemp[cbMagicIndex]--
						} else {
							bEnoughCard = false
							break
						}
					} else {
						cbCardIndexTemp[cbCardIndex]--
					}
				}
				//胡牌判断
				if bEnoughCard == true {
					cbCardEye := 0
					bMagicEye := false
					if lg.GetCardCount(cbCardIndexTemp) == 2 {
						if cbMagicIndex != 0 && cbCardIndexTemp[cbMagicIndex] == 2 {
							cbCardEye = lg.SwitchToCard(cbMagicIndex)
							bMagicEye = true
						} else {
							for i := 0; i < lg.GetCfg().MaxIdx; i++ {
								if cbCardIndexTemp[i] == 2 {
									cbCardEye = lg.SwitchToCard(i)
									if cbMagicIndex != 0 && i == cbMagicIndex {
										bMagicEye = true
									}
									break
								} else if i != cbMagicIndex && cbMagicIndex != 0 && cbCardIndexTemp[i]+cbCardIndexTemp[cbMagicIndex] == 2 {
									cbCardEye = lg.SwitchToCard(i)
									bMagicEye = true
									break
								}
							}
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
						for i := 0; i < cbLessKindItem; i++ {
							analyseItem.IsAnalyseGet[i+cbWeaveCount] = KindItem[i].IsAnalyseGet
							analyseItem.WeaveKind[i+cbWeaveCount] = KindItem[i].WeaveKind
							cbCenterCard := lg.SwitchToCard(KindItem[i].CenterCard)
							analyseItem.CenterCard[i+cbWeaveCount] = cbCenterCard
							lg.GetWeaveCard(KindItem[i].WeaveKind, cbCenterCard, analyseItem.CardData[i+cbWeaveCount])
						}
						//设置牌眼
						analyseItem.CardEye = cbCardEye
						analyseItem.MagicEye = bMagicEye
						//插入结果
						TagAnalyseItemArray = append(TagAnalyseItemArray, analyseItem)
					}
				}
			}
		}
	}

	//log.Debug("--------TagAnalyseItemArray=%d", len(TagAnalyseItemArray))
	//for _, ana := range TagAnalyseItemArray {
	//	log.Debug("====%v %v %v %v %v", ana.CenterCard, ana.WeaveKind, ana.CardEye, ana.MagicEye, ana.CardData)
	//}

	return true, TagAnalyseItemArray
}

//
func (lg *BaseLogic) AddKindItem(KindItem []*TagKindItem, Index []int, MagicIndex int, count int, opCode int) []*TagKindItem {
	tg := &TagKindItem{CardIndex: make([]int, 4)}
	switch opCode {
	case WIK_PENG:
		tg.CenterCard = Index[0]
		for t := 0; t <= 2; t++ {
			if count > t {
				tg.CardIndex[t] = Index[0]
			} else {
				tg.CardIndex[t] = MagicIndex
			}
		}
		tg.IsAnalyseGet = true
		tg.WeaveKind = WIK_PENG
	case WIK_LEFT:
		tg.CenterCard = Index[0]
		tg.CardIndex[0] = Index[1]
		tg.CardIndex[1] = Index[2]
		tg.CardIndex[2] = Index[3]
		tg.IsAnalyseGet = true
		tg.WeaveKind = WIK_LEFT
		tg.MagicCount = count
	}
	return append(KindItem, tg)
}

//扑克转换
func (lg *BaseLogic) GetUserCards(cbCardIndex []int) (cbCardData []int) {
	//转换扑克
	if lg.MagicIndex != lg.GetCfg().MaxIdx { //有财神， 把财神放进去
		for i := 0; i < cbCardIndex[lg.MagicIndex]; i++ {
			cbCardData = append(cbCardData, lg.SwitchToCard(lg.MagicIndex))
		}
	}
	for i := 0; i < lg.GetCfg().MaxIdx; i++ {
		if i == lg.MagicIndex && lg.MagicIndex != lg.ReplaceCard && lg.ReplaceCard != lg.GetCfg().MaxIdx {
			//如果财神有代替牌，则代替牌代替财神原来的位置
			for j := 0; j < cbCardIndex[lg.ReplaceCard]; j++ {
				cbCardData = append(cbCardData, lg.SwitchToCard(lg.ReplaceCard))
			}
			continue
		}

		if i == lg.ReplaceCard {
			continue
		}

		if cbCardIndex[i] != 0 {
			for j := 0; j < cbCardIndex[i]; j++ { //牌展开
				cbCardData = append(cbCardData, lg.SwitchToCard(i))
			}
		}
	}
	if len(cbCardData) < 14 {
		cbCardData = append(cbCardData, 0)
	}
	return cbCardData
}

//吃牌判断
func (lg *BaseLogic) EstimateEatCard(cbCardIndex []int, cbCurrentCard int) int {
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
	CurrentIndex := lg.SwitchToIdx(cbCurrentCard)
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

func (lg *BaseLogic) GetIteratorFunc(needCnt, allCnt int) func() []int {
	cbIndex := make([]int, 0)
	needCnt -= 1
	allCnt -= 1
	return func() []int {
		if len(cbIndex) < 1 {
			for i := 0; i <= needCnt; i++ {
				cbIndex = append(cbIndex, i)
			}
			return cbIndex
		}

		if cbIndex[needCnt] == allCnt {
			i := needCnt
			for ; i > 0; i-- {
				if (cbIndex[i-1] + 1) != cbIndex[i] {
					cbNewIndex := cbIndex[i-1]
					for j := (i - 1); j <= needCnt; j++ {
						cbIndex[j] = cbNewIndex + j - i + 2
					}
					break
				}
			}
			if i == 0 {
				return nil
			}
		} else {
			cbIndex[needCnt]++
		}
		return cbIndex
	}
}

func (lg *BaseLogic) IsZFB(card int) bool {
	return card == 0x35 || card == 0x36 || card == 0x37
}

func (lg *BaseLogic) IsFeng(card, quanfeng int) bool {
	return (card == 0x31 || card == 0x32 || card == 0x33 || card == 0x34) && ((card&MASK_VALUE)-1) == quanfeng
}

func (lg *BaseLogic) IsZhengHua(card, ProvideUser, playerCnt, BankerUser int) bool {
	zHua1, zHua2 := lg.GetZhengHuaCard(ProvideUser, playerCnt, BankerUser)
	return card == zHua1 || card == zHua2
}

func (lg *BaseLogic) GetZhengHuaCard(ProvideUser, PlayerCount, BankerUser int) (int, int) {
	if ProvideUser == BankerUser { //东风位
		return 0x38, 0x3C
	}

	if ProvideUser == (BankerUser+PlayerCount-1)%PlayerCount { //南风位
		return 0x39, 0x3D
	}

	if ProvideUser == (BankerUser+PlayerCount-2)%PlayerCount { //西风位
		return 0x3A, 0x3E
	}

	if ProvideUser == (BankerUser+PlayerCount-2)%PlayerCount { //北风位
		return 0x3B, 0x3F
	}

	log.Error("at GetZhengHuaCard error ..... ")
	return 0, 0
}

func (lg *BaseLogic) IsWeiFeng(card, ProvideUser, PlayerCount, BankerUser int) bool {
	if ProvideUser == BankerUser { //东风位
		return card == 0x31
	}

	if ProvideUser == (BankerUser+PlayerCount-1)%PlayerCount { //南风位
		return card == 0x32
	}

	if ProvideUser == (BankerUser+PlayerCount-2)%PlayerCount { //西风位
		return card == 0x33
	}

	if ProvideUser == (BankerUser+PlayerCount-2)%PlayerCount { //北风位
		return card == 0x34
	}

	return false
}
