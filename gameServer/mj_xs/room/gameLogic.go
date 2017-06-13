package room

import (
	"math"
	"mj/common/msg/mj_xs_msg"
	. "mj/gameServer/common/mj_logic_base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

//////////////////////////////////////////////////////////////////////////

//静态变量

var RightMask = make([]int64, MAX_RIGHT_COUNT)

//构造函数
func init() {
	for i := 0; i < MAX_RIGHT_COUNT; i++ {
		if 0 == i {
			RightMask[i] = 0
		} else {
			RightMask[i] = (int64(math.Pow(2, float64(i-1)))) << 28
		}
	}
}

//////////////////////////////////////////////////////////////////////////////////
type GameLogic struct {
	*BaseLogic
	CardDataArray []int //扑克数据
	MagicIndex    int   //钻牌索引
}

func NewGameLogic() *GameLogic {
	g := new(GameLogic)
	g.MagicIndex = MAX_INDEX
	return g
}

//有效判断
func (lg *GameLogic) IsValidCard(cbCardData int) bool {
	var cbValue = int(cbCardData & MASK_VALUE)
	var cbColor = int((cbCardData & MASK_COLOR) >> 4)
	return (((cbValue > 0) && (cbValue <= 9) && (cbColor <= 2)) || ((cbValue >= 1) && (cbValue <= 15) && (cbColor == 3)))
}

//胡牌等级
func (lg *GameLogic) GetChiHuActionRank(ChiHuResult *TagChiHuResult) int {
	//变量定义
	cbChiHuOrder := 0
	wChiHuRight := ChiHuResult.ChiHuRight
	wChiHuKind := (ChiHuResult.ChiHuKind & 0xFF00) >> 4

	//大胡升级
	for i := 0; i < 8; i++ {
		wChiHuKind >>= 1
		if (wChiHuKind & 0x0001) != 0 {
			cbChiHuOrder++
		}
	}

	//权位升级
	for i := 0; i < 16; i++ {
		wChiHuRight >>= 1
		if (wChiHuRight & 0x0001) != 0 {
			cbChiHuOrder++
		}
	}

	return cbChiHuOrder
}

//吃牌判断
func (lg *GameLogic) EstimateEatCard(cbCardIndex []int, cbCurrentCard int) int {
	//参数效验
	if cbCurrentCard >= 0x31 {
		return WIK_NULL
	}

	//变量定义
	var cbExcursion = []int{0, 1, 2}
	var cbItemKind = []int{WIK_LEFT, WIK_CENTER, WIK_RIGHT}

	//吃牌判断
	cbEatKind := 0
	cbFirstIndex := 0
	cbCurrentIndex := lg.SwitchToCardIndex(cbCurrentCard)
	for i := 0; i < len(cbItemKind); i++ {
		cbValueIndex := cbCurrentIndex % 9
		if (cbValueIndex >= cbExcursion[i]) && ((cbValueIndex - cbExcursion[i]) <= 6) {
			//吃牌判断
			cbFirstIndex = cbCurrentIndex - cbExcursion[i]
			if (cbCurrentIndex != cbFirstIndex) && (cbCardIndex[cbFirstIndex] == 0) {
				continue
			}

			if (cbCurrentIndex != (cbFirstIndex + 1)) && (cbCardIndex[cbFirstIndex+1] == 0) {
				continue
			}

			if (cbCurrentIndex != (cbFirstIndex + 2)) && (cbCardIndex[cbFirstIndex+2] == 0) {
				continue
			}

			//设置类型
			cbEatKind |= cbItemKind[i]
		}
	}

	return cbEatKind
}

//杠牌分析
func (lg *GameLogic) AnalyseGangCard(cbCardIndex []int, WeaveItem []*mj_xs_msg.TagWeaveItem, cbWeaveCount int, gangCardResult *TagGangCardResult) int {
	//设置变量
	cbActionMask := WIK_NULL

	//手上杠牌
	for i := 0; i < MAX_INDEX; i++ {
		if cbCardIndex[i] == 4 {
			cbActionMask |= WIK_GANG
			gangCardResult.CardData[gangCardResult.CardCount] = lg.SwitchToCardData(i)
			gangCardResult.CardCount++
		}
	}

	//组合杠牌
	for i := 0; i < cbWeaveCount; i++ {
		if WeaveItem[i].WeaveKind == WIK_PENG {
			if cbCardIndex[lg.SwitchToCardIndex(WeaveItem[i].CenterCard)] == 1 {
				cbActionMask |= WIK_GANG
				gangCardResult.CardData[gangCardResult.CardCount] = WeaveItem[i].CenterCard
				gangCardResult.CardCount++
			}
		}
	}

	return cbActionMask
}

//吃胡分析
func (lg *GameLogic) AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*mj_xs_msg.TagWeaveItem, cbWeaveCount, cbCurrentCard, wChiHuRight int,
	ChiHuResult *TagChiHuResult, cbHandtai, cbHandFeng *int, wFengQuan, wBankerUser int, bZimo bool) int {

	//变量定义
	wChiHuKind := WIK_NULL
	AnalyseItemArray := make([]*TagAnalyseItem, 0)

	//构造扑克
	cbCardIndexTemp := make([]int, MAX_INDEX)

	//插入扑克
	if cbCurrentCard != 0 {
		cbCardIndexTemp[lg.SwitchToCardIndex(cbCurrentCard)]++
	}

	//权位处理
	if cbWeaveCount == 0 {
		wChiHuRight |= CHR_MEN_QI
	}

	if lg.IsQingYiSe(cbCardIndexTemp, WeaveItem, cbWeaveCount) == true {
		wChiHuRight |= CHR_QING_YI_SE
	}

	if lg.IsHunYiSe(cbCardIndexTemp, WeaveItem, cbWeaveCount) == true {
		wChiHuRight |= CHR_HUN_YI_SE
	}

	if (cbWeaveCount == 4) && (cbCurrentCard != 0) && (bZimo == false) {
		wChiHuRight |= CHR_DA_DIAO
	}

	if (cbWeaveCount == 4) && (cbCurrentCard != 0) && (bZimo == true) {
		wChiHuRight |= CHR_DA_DIAO
	}

	if true == bZimo {
		wChiHuRight |= CHR_ZI_MO
	}

	//权位调整 杠胡+自摸 = 杠上花
	if (wChiHuRight&CHR_QIANG_GANG != 0) && (true == bZimo) {
		wChiHuRight &= ^CHR_QIANG_GANG
		wChiHuRight |= CHR_GANG_FLOWER
	}
	//杠上花 +1个杠 多加门清
	if (wChiHuRight&CHR_GANG_FLOWER != 0) && (cbWeaveCount == 1) {
		if WeaveItem[0].WeaveKind == WIK_GANG {
			wChiHuRight |= CHR_MEN_QI
		}
	}
	//分析扑克
	_, AnalyseItemArray = lg.AnalyseCard(cbCardIndexTemp, WeaveItem, cbWeaveCount, AnalyseItemArray)

	//胡牌分析
	alen := len(AnalyseItemArray)
	if alen > 0 {
		//牌型分析
		for i := 0; i < alen; i++ {
			//变量定义
			bLianCard := false
			bPengCard := false
			pAnalyseItem := AnalyseItemArray[i]

			//牌型分析
			for j := 0; j < len(pAnalyseItem.WeaveKind); j++ {
				cbWeaveKind := pAnalyseItem.WeaveKind[j]
				if (cbWeaveKind & (WIK_GANG | WIK_PENG)) != 0 {
					bPengCard = true
				}

				if (cbWeaveKind & (WIK_LEFT | WIK_CENTER | WIK_RIGHT)) != 0 {
					bLianCard = true
				}
			}

			//牌型判断

			//碰碰牌型
			if (bLianCard == false) && (bPengCard == true) {
				wChiHuKind |= CHK_PENG_PENG
			}

			if (bLianCard == true) && (bPengCard == true) {
				wChiHuKind |= CHK_JI_HU
			}

			if (bLianCard == true) && (bPengCard == false) {
				wChiHuKind |= CHK_PING_HU
			}

			//大三元
			if lg.IsDaSanYuan(pAnalyseItem) == true {
				*cbHandtai += 1
			}

		}
	}
	if wChiHuKind != 0 {
		//判断边 嵌 对倒 单吊
		if (cbWeaveCount < 4) && (alen > 0) && (cbCurrentCard != 0) {
			for i := 0; i < alen; i++ {
				//变量定义
				bLianCard := false
				bPengCard := false
				wTempChihuKind := CHK_NULL
				pAnalyseItem := AnalyseItemArray[i]

				//牌型分析
				for j := 0; j < len(pAnalyseItem.WeaveKind); j++ {
					cbWeaveKind := pAnalyseItem.WeaveKind[j]
					if (cbWeaveKind & (WIK_GANG | WIK_PENG)) != 0 {
						bPengCard = true
					}

					if (cbWeaveKind & (WIK_LEFT | WIK_CENTER | WIK_RIGHT)) != 0 {
						bLianCard = true
					}
				}

				//碰碰牌型
				if (bLianCard == false) && (bPengCard == true) {
					wTempChihuKind |= CHK_PENG_PENG
				}

				if (bLianCard == true) && (bPengCard == true) {
					wTempChihuKind |= CHK_JI_HU
				}

				if (bLianCard == true) && (bPengCard == false) {
					wTempChihuKind |= CHK_PING_HU
				}

				//判断中发白
				if wTempChihuKind != CHK_NULL {
					//牌型分析
					for j := cbWeaveCount; j < len(pAnalyseItem.WeaveKind); j++ {
						cbWeaveKind := pAnalyseItem.WeaveKind[j]
						cbCenterCard := pAnalyseItem.CenterCard[j]
						if (cbWeaveKind & (WIK_GANG | WIK_PENG)) != 0 {
							bPengCard = true
						} else {
							bPengCard = bPengCard
						}

						if (cbCenterCard == 0x35) || (cbCenterCard == 0x36) || (cbCenterCard == 0x37) {
							*cbHandtai++
						}

						if (cbCenterCard == 0x31) || (cbCenterCard == 0x32) || (cbCenterCard == 0x33) || (cbCenterCard == 0x34) {
							//圈风
							if ((cbCenterCard & MASK_VALUE) - 1) == wFengQuan {
								*cbHandFeng++
							}

							//位风
							if ((cbCenterCard & MASK_VALUE) - 1) == wBankerUser {
								*cbHandFeng++
							}

						}
					}
				}

				//判断边嵌对倒
				if wTempChihuKind != CHK_NULL {

					if pAnalyseItem.CardEye == cbCurrentCard { // 单吊
						wChiHuRight |= CHR_DAN_DIAO
						wChiHuKind &= (^CHK_PING_HU)
						wChiHuKind |= CHK_JI_HU
						break
					} else {
						for j := cbWeaveCount; j < len(pAnalyseItem.WeaveKind); j++ {
							if (pAnalyseItem.CardData[j][0] == cbCurrentCard) && (pAnalyseItem.WeaveKind[j] == WIK_PENG) {
								wChiHuRight |= CHR_DUI_DAO
								wChiHuKind &= (^CHK_PING_HU)
								wChiHuKind |= CHK_JI_HU
								break
							}
							if (pAnalyseItem.CardData[j][1] == cbCurrentCard) && (pAnalyseItem.WeaveKind[j] == WIK_LEFT) {
								wChiHuRight |= CHR_QIAN
								wChiHuKind &= (^CHK_PING_HU)
								wChiHuKind |= CHK_JI_HU
								break
							}
							//判断边 嵌 由于在AnalyCard中全是WIK_LEFT 所以导致判断边 嵌麻烦
							if (((pAnalyseItem.CardData[j][2] == cbCurrentCard) && ((cbCurrentCard & MASK_VALUE) == 3)) || ((pAnalyseItem.CardData[j][0] == cbCurrentCard) && ((cbCurrentCard & MASK_VALUE) == 7))) && (pAnalyseItem.WeaveKind[j] == WIK_LEFT) {
								wChiHuRight |= CHR_BIAN
								wChiHuKind &= (^CHK_PING_HU)
								wChiHuKind |= CHK_JI_HU
								break
							}
						}
						if wChiHuRight != 0 {
							break
						}
					}
					if wChiHuRight != 0 {
						break
					}
				}
			}
		}
	}

	//结果判断
	if wChiHuKind != 0 {
		//设置结果
		ChiHuResult.ChiHuKind = wChiHuKind
		ChiHuResult.ChiHuRight = wChiHuRight
		return WIK_CHI_HU
	}

	return WIK_NULL

}

//清一色
func (lg *GameLogic) IsQingYiSe(cbCardIndex []int, WeaveItem []*mj_xs_msg.TagWeaveItem, cbItemCount int) bool {
	//变量定义
	cbCardColor := 0xFF
	for i := 0; i < MAX_WEAVE; i++ {
		if cbCardIndex[i] != 0 {
			if cbCardColor != 0xFF {
				return false
			}
			//设置花色
			cbCardColor = (lg.SwitchToCardData(i) & MASK_COLOR)
			//设置索引
			i = (i/9+1)*9 - 1
		}

	}

	//组合判断
	for i := 0; i < cbItemCount; i++ {
		cbCenterCard := WeaveItem[i].CenterCard
		if (cbCenterCard & MASK_COLOR) != cbCardColor {
			return false
		}
	}

	return true
}

//混一色
func (lg *GameLogic) IsHunYiSe(cbCardIndex []int, WeaveItem []*mj_xs_msg.TagWeaveItem, cbItemCount int) bool {
	//变量定义
	bColorFlags := make([]bool, 5)

	//扑克扑克
	for i := 0; i < MAX_INDEX; i++ {
		if cbCardIndex[i] != 0 {
			bColorFlags[i/9] = true
		}

	}

	//组合判断
	for i := 0; i < cbItemCount; i++ {
		cbCenterCard := WeaveItem[i].CenterCard
		bColorFlags[(cbCenterCard&MASK_COLOR)>>4] = true
	}

	//花色统计
	cbColorCount := 0
	for i := 0; i < 3; i++ {
		if bColorFlags[i] == true {
			cbColorCount++
		}
	}
	if (cbColorCount != 1) || (bColorFlags[3] == false) {
		return false
	}

	return true
}

//扑克转换
func (lg *GameLogic) SwitchToCardData(cbCardIndex int) int {
	if cbCardIndex < 34 { //花三种花色牌 3 * 9
		return ((cbCardIndex / 9) << 4) | (cbCardIndex%9 + 1)
	}
	return 48 | ((cbCardIndex-34)%8 + 8)
}

//扑克转换
func (lg *GameLogic) SwitchToCardIndex(cbCardData int) int {
	//计算位置
	cbValue := cbCardData & MASK_VALUE
	cbColor := (cbCardData & MASK_COLOR) >> 4

	if cbColor >= 0x03 {
		return cbValue + 27 - 1
	}
	return cbColor*9 + cbValue - 1
}

//扑克转换
func (lg *GameLogic) SwitchToCardData2(cbCardIndex, cbCardData []int) int {
	//转换扑克
	cbPosition := 0
	for i := 0; i < MAX_INDEX; i++ {
		if cbCardIndex[i] != 0 {
			for j := 0; j < cbCardIndex[i]; j++ {
				cbCardData[cbPosition] = lg.SwitchToCardData(i)
				cbPosition++
			}
		}
	}

	return cbPosition
}

//扑克转换
func (lg *GameLogic) SwitchToCardIndex3(cbCardData []int, cbCardCount int, cbCardIndex []int) int {
	//转换扑克
	for i := 0; i < cbCardCount; i++ {
		cbCardIndex[lg.SwitchToCardIndex(cbCardData[i])]++
	}

	return cbCardCount
}

func (lg *GameLogic) CalScore(ChiHuResult *TagChiHuResult) int {
	lGain := 0
	wChihuKind := ChiHuResult.ChiHuKind
	if wChihuKind&CHK_JI_HU != 0 {
		lGain += KIND_JI_HU
	}

	if wChihuKind&CHK_PING_HU != 0 {
		lGain += KIND_PING_HU
	}

	if wChihuKind&CHK_PENG_PENG != 0 {
		lGain += KIND_PENG_PEMG
	}

	if wChihuKind&CHK_BA_HUA != 0 {
		lGain += KIND_BA_HUA
	}

	wChihuRight := ChiHuResult.ChiHuRight
	if wChihuRight&CHR_HUN_YI_SE != 0 {
		lGain += RIGHT_HUN_YI_SE
	}
	if wChihuRight&CHR_GANG_FLOWER != 0 {
		lGain += RIGHT_GANG_FLOWER
	}
	if wChihuRight&CHR_HAI_DI != 0 {
		lGain += RIGHT_HAI_DI
	}
	if wChihuRight&CHR_ZI_MO != 0 {
		lGain += RIGHT_ZI_MO
	}
	if wChihuRight&CHR_QING_YI_SE != 0 {
		lGain += RIGHT_QING_YI_SE
	}
	if wChihuRight&CHR_MEN_QI != 0 {
		lGain += RIGHT_MEN_QI
	}
	if wChihuRight&CHR_DI != 0 {
		lGain += RIGHT_DI
	}
	if wChihuRight&CHR_TIAN != 0 {
		lGain += RIGHT_TIAN
	}
	if wChihuRight&CHR_QIANG_GANG != 0 {
		lGain += RIGHT_QIANG_GANG
	}
	if wChihuRight&CHR_DA_DIAO != 0 {
		lGain += RIGHT_DA_DIAO
	}
	if wChihuRight&CHR_BIAN != 0 {
		lGain += RIGHT_BIAN
	}
	if wChihuRight&CHR_QIAN != 0 {
		lGain += RIGHT_QIAN
	}
	if wChihuRight&CHR_DUI_DAO != 0 {
		lGain += RIGHT_DUI_DAO
	}
	if wChihuRight&CHR_DAN_DIAO != 0 {
		lGain += RIGHT_DAN_DIAO
	}
	if wChihuRight&CHR_SI_HUA != 0 {
		lGain += RIGHT_SI_HUA
	}
	if wChihuRight&CHR_BA_HUA != 0 {
		lGain += RIGHT_BA_HUA
	}
	return lGain
}

//分析扑克
func (lg *GameLogic) AnalyseCard(cbCardIndex []int, WeaveItem []*mj_xs_msg.TagWeaveItem, cbWeaveCount int, TagAnalyseItemArray []*TagAnalyseItem) (bool, []*TagAnalyseItem) { //todo , CTagAnalyseItemArray & TagAnalyseItemArray
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
				analyseItem := &TagAnalyseItem{WeaveKind: make([]int, 4), CenterCard: make([]int, 4), CardData: make([][]int, 4)}
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
		cbIndex := []int{0, 1, 2, 3}

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
						cbCardEye = lg.SwitchToCardData(i)
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
						cbCenterCard := lg.SwitchToCardData(pKindItem[i].CenterCard)
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

//大三元牌
func (lg *GameLogic) IsDaSanYuan(AnalyseItem *TagAnalyseItem) bool {
	bExist := []bool{false, false, false}
	for i := 0; i < len(AnalyseItem.WeaveKind); i++ {
		if AnalyseItem.CenterCard[i] == 0x35 {
			bExist[0] = true
		}
		if AnalyseItem.CenterCard[i] == 0x36 {
			bExist[1] = true
		}
		if AnalyseItem.CenterCard[i] == 0x37 {
			bExist[2] = true
		}
	}

	if bExist[0] && bExist[1] && bExist[2] {
		return true
	}

	return false
}

func (lg *GameLogic) isZFB(card int) bool {
	return card == 0x35 || card == 0x36 || card == 0x37
}

func (lg *GameLogic) isFeng(card, quanfeng int) bool {
	return (card == 0x31 || card == 0x32 || card == 0x33 || card == 0x34) && ((card&MASK_VALUE)-1) == quanfeng
}

func (lg *GameLogic) isZhengHua(card, ProvideUser, playerCnt, BankerUser int) bool {
	zHua1, zHua2 := lg.GetZhengHuaCard(ProvideUser, playerCnt, BankerUser)
	return card == zHua1 || card == zHua2
}

func (lg *GameLogic) GetZhengHuaCard(ProvideUser, PlayerCount, BankerUser int) (int, int) {
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

func (lg *GameLogic) isWeiFeng(card, ProvideUser, PlayerCount, BankerUser int) bool {
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
