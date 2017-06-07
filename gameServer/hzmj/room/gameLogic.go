package room

import (
	"math"
	"mj/common/msg"

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
	CardDataArray []int //扑克数据
	MagicIndex    int   //钻牌索引
}

func NewGameLogic() *GameLogic {
	g := new(GameLogic)
	g.MagicIndex = MAX_INDEX
	return g
}

func (lg *GameLogic) SetMagicIndex(cbMagicIndex int) {
	lg.MagicIndex = cbMagicIndex
}

//混乱扑克
func (lg *GameLogic) RandCardList(cbCardBuffer, OriDataArray []int) {
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
		cbPosition = int(util.RandInterval(0, int(cbBufferCount-cbRandCount)))
		cbCardBuffer[cbRandCount] = cbCardDataTemp[cbPosition]
		cbRandCount++
		cbCardDataTemp[cbPosition] = cbCardDataTemp[cbBufferCount-cbRandCount]
	}

	return
}

//删除扑克
func (lg *GameLogic) RemoveCardByArray(cbCardData []int, cbCardCount int, cbRemoveCard []int, cbRemoveCount int) bool {
	//效验扑克

	//定义变量
	var cbDeleteCount int
	var cbTempCardData []int
	if int(cbCardCount) > len(cbTempCardData) {
		return false
	}

	util.DeepCopy(&cbTempCardData, &cbCardData)

	//置零扑克
	for i := 0; i < cbRemoveCount; i++ {
		for j := 0; j < cbCardCount; j++ {
			if cbRemoveCard[i] == cbTempCardData[j] {
				cbDeleteCount++
				cbTempCardData[j] = 0
				break
			}
		}
	}

	//成功判断
	if cbDeleteCount != cbRemoveCount {
		return false
	}

	//清理扑克
	cbCardPos := 0
	for i := 0; i < cbCardCount; i++ {
		if cbTempCardData[i] != 0 {
			cbCardData[cbCardPos] = cbTempCardData[i]
			cbCardPos++
		}
	}

	return true
}

//删除扑克
func (lg *GameLogic) RemoveCardByCnt(cbCardIndex, cbRemoveCard []int, cbRemoveCount int) bool {
	//删除扑克
	for i := 0; i < cbRemoveCount; i++ {
		//效验扑克
		if !lg.IsValidCard(cbRemoveCard[i]) {
			log.Error("at RemoveCardByCnt card is Invalid %d", cbRemoveCard[i])
		}
		if cbCardIndex[lg.SwitchToCardIndex(cbRemoveCard[i])] < 0 {
			log.Error("at RemoveCardByCnt 11 card is Invalid %d", cbRemoveCard[i])
		}

		//删除扑克
		cbRemoveIndex := lg.SwitchToCardIndex(cbRemoveCard[i])
		if cbCardIndex[cbRemoveIndex] == 0 {

			//还原删除
			for j := 0; j < i; j++ {
				cbCardIndex[lg.SwitchToCardIndex(cbRemoveCard[j])]++
			}
			return false
		} else {
			//删除扑克
			cbCardIndex[cbRemoveIndex]--
		}
	}

	return true
}

//删除扑克
func (lg *GameLogic) RemoveCard(cbCardIndex []int, cbRemoveCard int) bool {
	//删除扑克
	cbRemoveIndex := lg.SwitchToCardIndex(cbRemoveCard)
	//效验扑克
	if !lg.IsValidCard(cbRemoveCard) {
		log.Error("at RemoveCard card is Invalid %d", cbRemoveCard)
	}
	if cbCardIndex[lg.SwitchToCardIndex(cbRemoveCard)] < 0 {
		log.Error("at RemoveCard 11 card is Invalid %d", cbRemoveCard)
	}
	if cbCardIndex[cbRemoveIndex] > 0 {
		cbCardIndex[cbRemoveIndex]--
		return true
	}

	return false
}

//财神判断
func (lg *GameLogic) IsMagicCard(cbCardData int) bool {
	if lg.MagicIndex != MAX_INDEX {
		return lg.SwitchToCardIndex(cbCardData) == lg.MagicIndex
	}

	return false
}

//花牌判断
func (lg *GameLogic) IsHuaCard(cbCardData int) bool {
	return cbCardData >= 0x38
}

//花牌判断
func (lg *GameLogic) HuaCardCnt(cbCardIndex []int) int {
	cbHuaCardCount := 0
	for i := MAX_INDEX - MAX_HUA_INDEX; i < MAX_INDEX; i++ {
		if cbCardIndex[i] > 0 {
			cbHuaCardCount += cbCardIndex[i]
		}
	}

	return cbHuaCardCount
}

//排序,根据牌值排序
func (lg *GameLogic) SortCardList(cbCardData []int, cbCardCount int) bool {
	//数目过虑
	if cbCardCount == 0 || cbCardCount > MAX_COUNT {
		return false
	}

	//排序操作
	bSorted := true
	var cbSwitchData int
	cbLast := cbCardCount - 1
	var cbCard1 = 0
	var cbCard2 = 0
	for {
		bSorted = true
		for i := 0; i < cbLast; i++ {
			cbCard1 = cbCardData[i]
			cbCard2 = cbCardData[i+1]
			//如果财神有代替牌，财神与代替牌转换
			if INDEX_REPLACE_CARD != MAX_INDEX && lg.MagicIndex != INDEX_REPLACE_CARD {
				if lg.SwitchToCardIndex(cbCard1) == INDEX_REPLACE_CARD {
					cbCard1 = lg.SwitchToCardData(lg.MagicIndex)
				}
				if lg.SwitchToCardIndex(cbCard2) == INDEX_REPLACE_CARD {
					cbCard2 = lg.SwitchToCardData(lg.MagicIndex)
				}
			}
			if cbCard1 > cbCard2 {
				//设置标志
				bSorted = false

				//扑克数据
				cbSwitchData = cbCardData[i]
				cbCardData[i] = cbCardData[i+1]
				cbCardData[i+1] = cbSwitchData
			}
		}
		cbLast--
		if bSorted == false {
			break
		}
	}

	return true
}

//动作等级
func (lg *GameLogic) GetUserActionRank(cbUserAction int) int {
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

//胡牌等级
func (lg *GameLogic) GetChiHuActionRank(CChiHuRight int64) int {
	return 1
}

//胡牌倍数
func (lg *GameLogic) GetChiHuTime(ChiHuRight int64) int {
	wFanShu := 1 //平胡一倍
	if (ChiHuRight & CHR_TIAN_HU) != 0 {
		wFanShu = 16
	} else if (ChiHuRight & CHR_QI_SHOU_LISTEN) != 0 { //起首听
		wFanShu = 8
	} else if (ChiHuRight & CHR_PENG_PENG) != 0 { //碰碰胡
		wFanShu = 3
	} else if (ChiHuRight & CHR_DAN_DIAN_QI_DUI) != 0 { //单点大七对
		wFanShu = 6
	} else if (ChiHuRight & CHR_MA_QI_DUI) != 0 { //麻七对
		wFanShu = 5
	} else if (ChiHuRight & CHR_MA_QI_WANG) != 0 { //麻七王
		wFanShu = 10
	} else if (ChiHuRight & CHR_MA_QI_WZW) != 0 { //麻七王中王
		wFanShu = 20
	} else if (ChiHuRight & CHR_SHI_SAN_LAN) != 0 { //十三烂
		wFanShu = 2
	} else if (ChiHuRight & CHR_QX_SHI_SAN_LAN) != 0 { //七星十三烂
		wFanShu = 4
	}

	if (ChiHuRight & CHR_GANG_SHANG_HUA) != 0 { //杠上开花
		wFanShu *= 2
	} else if (ChiHuRight & CHR_QIANG_GANG_HU) != 0 { //抢杠
		wFanShu = 2*wFanShu + 1
	}

	if (ChiHuRight & CHR_ZI_MO) != 0 {
		wFanShu *= 2
	}

	return wFanShu
}

//自动出牌
func (lg *GameLogic) AutomatismOutCard(cbCardIndex, cbEnjoinOutCard []int, cbEnjoinOutCardCount int) int {
	// 先打财神
	if lg.MagicIndex != MAX_INDEX {
		if cbCardIndex[lg.MagicIndex] > 0 {
			return lg.SwitchToCardData(lg.MagicIndex)
		}
	}

	//而后打字牌，字牌打自己多的，数目一样就按东南西北中发白的顺序
	var cbCardData int
	var cbOutCardIndex = int(MAX_INDEX)
	var cbOutCardIndexCount int
	for i := int(MAX_INDEX - 7); i < MAX_INDEX-1; i++ {
		if cbCardIndex[i] > cbOutCardIndexCount {
			cbOutCardIndexCount = cbCardIndex[i]
			cbOutCardIndex = i
		}
	}

	if cbOutCardIndex != MAX_INDEX {
		cbCardData = lg.SwitchToCardData(cbOutCardIndex)
		bEnjoinCard := false
		for k := 0; k < cbEnjoinOutCardCount; k++ {
			if cbCardData == cbEnjoinOutCard[k] {
				bEnjoinCard = true
			}
		}
		if !bEnjoinCard {
			return cbCardData
		}
	}

	//没有字牌就打边张，1或9，顺序为万筒索，2，8，而后3，7，而后4，6，而后5
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			cbOutCardIndex = MAX_INDEX
			if cbCardIndex[j*9+i] > 0 {
				cbOutCardIndex = int(j*9 + i)
			} else if cbCardIndex[j*9+(9-i-1)] > 0 {
				cbOutCardIndex = int(j*9 + (9 - i - 1))
			}

			if cbOutCardIndex != MAX_INDEX {
				cbCardDataTemp := lg.SwitchToCardData(cbOutCardIndex)
				bEnjoinCard := false
				for k := 0; k < cbEnjoinOutCardCount; k++ {
					if cbCardDataTemp == cbEnjoinOutCard[k] {
						bEnjoinCard = true
					}
				}
				if !bEnjoinCard {
					return cbCardDataTemp
				} else {
					if cbCardData == 0 {
						cbCardData = cbCardDataTemp
					}
				}
			}
		}
	}
	return cbCardData
}

//吃牌判断
func (lg *GameLogic) EstimateEatCard(cbCardIndex []int, cbCurrentCard int) int {
	//参数效验
	cbCurrentIndex := lg.SwitchToCardIndex(cbCurrentCard)

	//过滤判断
	if cbCurrentIndex == lg.MagicIndex {
		return WIK_NULL
	}
	if cbCurrentIndex == INDEX_REPLACE_CARD && lg.MagicIndex >= 27 {
		return WIK_NULL
	}
	if cbCurrentCard >= 0x31 && cbCurrentIndex != INDEX_REPLACE_CARD {
		return WIK_NULL
	}

	//变量定义
	var cbExcursion = []int{0, 1, 2}
	var cbItemKind = []int{WIK_LEFT, WIK_CENTER, WIK_RIGHT}

	//拆分分析
	var cbMagicCardIndex []int
	util.DeepCopy(&cbMagicCardIndex, &cbCardIndex)

	//如果有财神
	cbMagicCardCount := 0
	if lg.MagicIndex != MAX_INDEX {
		cbMagicCardCount = cbCardIndex[lg.MagicIndex]
		//如果财神有代替牌，财神与代替牌转换
		if INDEX_REPLACE_CARD != MAX_INDEX {
			cbMagicCardIndex[lg.MagicIndex] = cbMagicCardIndex[INDEX_REPLACE_CARD]
			cbMagicCardIndex[INDEX_REPLACE_CARD] = cbMagicCardCount
		}
	}

	//吃牌判断
	cbEatKind := 0
	cbFirstIndex := 0
	if cbCurrentIndex == INDEX_REPLACE_CARD {
		cbCurrentIndex = lg.MagicIndex
	}
	for i := 0; i < len(cbItemKind); i++ {
		cbValueIndex := cbCurrentIndex % 9
		if (cbValueIndex >= cbExcursion[i]) && ((cbValueIndex - cbExcursion[i]) <= 6) {
			//吃牌判断
			cbFirstIndex = cbCurrentIndex - cbExcursion[i]

			//吃牌不能包含有财神
			if lg.MagicIndex != MAX_INDEX &&
				lg.MagicIndex >= cbFirstIndex && lg.MagicIndex <= cbFirstIndex+2 {
				continue
			}

			if (cbCurrentIndex != cbFirstIndex) && (cbMagicCardIndex[cbFirstIndex] == 0) {
				continue
			}

			if (cbCurrentIndex != (cbFirstIndex + 1)) && (cbMagicCardIndex[cbFirstIndex+1] == 0) {
				continue
			}

			if (cbCurrentIndex != (cbFirstIndex + 2)) && (cbMagicCardIndex[cbFirstIndex+2] == 0) {
				continue
			}

			//设置类型
			cbEatKind |= cbItemKind[i]
		}
	}

	return cbEatKind
}

//碰牌判断
func (lg *GameLogic) EstimatePengCard(cbCardIndex []int, cbCurrentCard int) int {
	//参数效验

	//过滤判断
	if lg.IsMagicCard(cbCurrentCard) || lg.IsHuaCard(cbCurrentCard) {
		return WIK_NULL
	}

	//碰牌判断
	if cbCardIndex[lg.SwitchToCardIndex(cbCurrentCard)] >= 2 {
		return WIK_PENG
	}

	return WIK_NULL
}

//杠牌判断
func (lg *GameLogic) EstimateGangCard(cbCardIndex []int, cbCurrentCard int) int {
	//参数效验

	//过滤判断
	if lg.IsMagicCard(cbCurrentCard) || lg.IsHuaCard(cbCurrentCard) {
		return WIK_NULL
	}

	//杠牌判断
	if cbCardIndex[lg.SwitchToCardIndex(cbCurrentCard)] == 3 {
		return WIK_GANG
	}

	return WIK_NULL

}

//杠牌分析
func (lg *GameLogic) AnalyseGangCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount int, gangCardResult *TagGangCardResult) int {
	//设置变量
	cbActionMask := int(WIK_NULL)

	//手上杠牌
	for i := 0; i < MAX_INDEX; i++ {
		if i == lg.MagicIndex {
			continue
		}
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

func (lg *GameLogic) AnalyseGangCardEx(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount, cbProvideCard int, gangCardResult *TagGangCardResult) int {
	//设置变量
	cbActionMask := int(WIK_NULL)
	gangCardResult.CardData = make([]int, MAX_WEAVE)
	//手上杠牌
	for i := 0; i < MAX_INDEX; i++ {
		if i == lg.MagicIndex {
			continue
		}
		if cbCardIndex[i] == 4 {
			cbActionMask |= WIK_GANG
			gangCardResult.CardData[gangCardResult.CardCount] = lg.SwitchToCardData(i)
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
		}
	}

	return cbActionMask
}

func (lg *GameLogic) GetCardColor(cbCardData int) int { return cbCardData & MASK_COLOR }
func (lg *GameLogic) GetCardValue(cbCardData int) int { return cbCardData & MASK_VALUE }

//吃胡分析
func (lg *GameLogic) AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount, cbCurrentCard int, ChiHuRight int, b4HZHu bool) int {
	//变量定义
	cbChiHuKind := int(WIK_NULL)
	TagAnalyseItemArray := make([]*TagAnalyseItem, MAX_INDEX) //todo MAX_INDEX CTagAnalyseItemArray

	//构造扑克
	cbCardIndexTemp := make([]int, MAX_INDEX)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)

	//cbCurrentCard一定不为0			!!!!!!!!!
	//ASSERT(cbCurrentCard != 0);
	if cbCurrentCard == 0 {
		return WIK_NULL
	}

	//插入扑克
	if cbCurrentCard != 0 {
		cbCardIndexTemp[lg.SwitchToCardIndex(cbCurrentCard)]++
	}

	if b4HZHu && cbCardIndexTemp[31] == 4 { //四个红中直接胡牌
		return WIK_CHI_HU
	}
	//分析扑克
	lg.AnalyseCard(cbCardIndexTemp, WeaveItem, cbWeaveCount, TagAnalyseItemArray)

	//胡牌分析
	if len(TagAnalyseItemArray) > 0 {
		//牌型分析
		// 		for  i := 0; i< len(TagAnalyseItemArray); i++ {
		// 			//变量定义
		// 			tagTagAnalyseItem * pTagAnalyseItem=&TagAnalyseItemArray[i];
		//
		//  		//碰碰胡
		//  		if(IsPengPeng(pTagAnalyseItem))
		// 				ChiHuRight |= CHR_PENG_PENG;
		// 			}
		// 		}

		ChiHuRight |= CHR_PING_HU
	}

	if ChiHuRight != 0 {
		cbChiHuKind = WIK_CHI_HU
	}

	return cbChiHuKind
}

//听牌分析
func (lg *GameLogic) AnalyseTingCard2(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount int) int {
	//复制数据
	cbCardIndexTemp := make([]int, MAX_INDEX)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)

	cbCardCount := lg.GetCardCount(cbCardIndexTemp)
	chr := 0

	if (cbCardCount+1)%3 == 0 {
		for i := 0; i < MAX_INDEX-MAX_HUA_INDEX; i++ {
			if cbCardIndexTemp[i] == 0 {
				continue
			}
			cbCardIndexTemp[i]--

			for j := 0; j < MAX_INDEX-MAX_HUA_INDEX; j++ {
				cbCurrentCard := lg.SwitchToCardData(j)
				if WIK_CHI_HU == lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, cbWeaveCount, cbCurrentCard, chr, false) {
					return WIK_LISTEN
				}

			}

			cbCardIndexTemp[i]++
		}
	} else {
		for j := 0; j < MAX_INDEX-MAX_HUA_INDEX; j++ {
			cbCurrentCard := lg.SwitchToCardData(j)
			if WIK_CHI_HU == lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, cbWeaveCount, cbCurrentCard, chr, false) {
				return WIK_LISTEN
			}
		}
	}

	return WIK_NULL
}

//听牌分析
func (lg *GameLogic) AnalyseTingCard1(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount int, cbOutCard [][]int) int {
	//复制数据
	cbOutCount := 0
	cbCardIndexTemp := make([]int, MAX_INDEX)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)

	cbCardCount := lg.GetCardCount(cbCardIndexTemp)
	chr := 0

	if (cbCardCount-2)%3 == 0 {
		for i := 0; i < MAX_INDEX-MAX_HUA_INDEX; i++ {
			if cbCardIndexTemp[i] == 0 {
				continue
			}
			cbCardIndexTemp[i]--

			bAdd := false
			nCount := 0
			for j := 0; j < MAX_INDEX-MAX_HUA_INDEX; j++ {
				cbCurrentCard := lg.SwitchToCardData(j)
				if WIK_CHI_HU == lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, cbWeaveCount, cbCurrentCard, chr, false) {
					if bAdd == false {
						bAdd = true
						cbOutCard[0][cbOutCount] = lg.SwitchToCardData(i)
						cbOutCount++
					}
					cbOutCard[cbOutCount][nCount] = lg.SwitchToCardData(j)
					nCount++
				}
			}

			cbCardIndexTemp[i]++
		}
	}

	return cbOutCount
}

func (lg *GameLogic) AnalyseTingCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount int, cbOutCardData, cbHuCardCount []int, cbHuCardData [][]int) int {

	//复制数据
	cbOutCount := 0
	cbCardIndexTemp := make([]int, MAX_INDEX)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)

	cbCardCount := lg.GetCardCount(cbCardIndexTemp)
	chr := 0

	if (cbCardCount-2)%3 == 0 {
		for i := 0; i < MAX_INDEX-MAX_HUA_INDEX; i++ {
			if cbCardIndexTemp[i] == 0 {
				continue
			}
			cbCardIndexTemp[i]--

			bAdd := false
			nCount := 0
			for j := 0; j < MAX_INDEX-MAX_HUA_INDEX; j++ {
				cbCurrentCard := lg.SwitchToCardData(j)
				if WIK_CHI_HU == lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, cbWeaveCount, cbCurrentCard, chr, false) {
					if bAdd == false {
						bAdd = true
						cbOutCardData[cbOutCount] = lg.SwitchToCardData(i)
						cbOutCount++
					}
					if len(cbHuCardData[cbOutCount-1]) < 1 {
						cbHuCardData[cbOutCount-1] = make([]int, MAX_INDEX-MAX_HUA_INDEX)
					}
					cbHuCardData[cbOutCount-1][nCount] = lg.SwitchToCardData(j)
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
		for j := 0; j < MAX_INDEX; j++ {
			cbCurrentCard := lg.SwitchToCardData(j)
			if WIK_CHI_HU == lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, cbWeaveCount, cbCurrentCard, chr, false) {
				if len(cbHuCardData[0]) < 1 {
					cbHuCardData[cbOutCount-1] = make([]int, MAX_INDEX)
				}
				cbHuCardData[0][cbCount] = cbCurrentCard
				cbCount++
			}
		}
		cbHuCardCount[0] = cbCount
	}

	return cbOutCount
}

func (lg *GameLogic) GetHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount int, cbHuCardData []int) int {
	//复制数据
	cbCardIndexTemp := make([]int, MAX_INDEX)
	util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)

	cbCount := 0
	//ZeroMemory(cbHuCardData,sizeof(cbHuCardData));

	cbCardCount := lg.GetCardCount(cbCardIndexTemp)
	chr := 0

	if (cbCardCount-2)%3 != 0 {
		for j := 0; j < MAX_INDEX; j++ {
			cbCurrentCard := lg.SwitchToCardData(j)
			if WIK_CHI_HU == lg.AnalyseChiHuCard(cbCardIndexTemp, WeaveItem, cbWeaveCount, cbCurrentCard, chr, false) {
				cbHuCardData[cbCount] = cbCurrentCard
				cbCount++
			}
		}
		if cbCount > 0 {
			return cbCount
		}
	}
	return 0
}

//扑克转换
func (lg *GameLogic) SwitchToCardData(cbCardIndex int) int {
	if cbCardIndex < 27 { //花三种花色牌 3 * 9
		return ((cbCardIndex / 9) << 4) | (cbCardIndex%9 + 1)
	}
	return (0x30 | (cbCardIndex - 27 + 1))
}

//扑克转换
func (lg *GameLogic) SwitchToCardIndex(cbCardData int) int {
	return ((cbCardData&MASK_COLOR)>>4)*9 + (cbCardData & MASK_VALUE) - 1 //高位乘以9 + 低位 9 是系数 用来区分index的花色 index/ 9取整就是花色 所有的牌都是不会大于9 是余数
}

//扑克转换
func (lg *GameLogic) SwitchToCardData2(cbCardIndex, cbCardData []int) int {
	//转换扑克
	cbPosition := 0
	//财神
	if lg.MagicIndex != MAX_INDEX {
		for i := 0; i < cbCardIndex[lg.MagicIndex]; i++ {
			cbCardData[cbPosition] = lg.SwitchToCardData(lg.MagicIndex)
			cbPosition++
		}

	}
	for i := 0; i < MAX_INDEX; i++ {
		if i == lg.MagicIndex && lg.MagicIndex != INDEX_REPLACE_CARD {
			//如果财神有代替牌，则代替牌代替财神原来的位置
			if INDEX_REPLACE_CARD != MAX_INDEX {
				for j := 0; j < cbCardIndex[INDEX_REPLACE_CARD]; j++ {
					cbCardData[cbPosition] = lg.SwitchToCardData(INDEX_REPLACE_CARD)
					cbPosition++
				}

			}
			continue
		}
		if i == INDEX_REPLACE_CARD {
			continue
		}
		if cbCardIndex[i] != 0 {
			for j := 0; j < cbCardIndex[i]; j++ { //牌展开
				//ASSERT(cbPosition<MAX_COUNT);
				cbCardData[cbPosition] = lg.SwitchToCardData(i)
				cbPosition++
			}
		}
	}

	return cbPosition
}

//扑克转换
func (lg *GameLogic) SwitchToCardIndex3(cbCardData []int, cbCardCount int, cbCardIndex []int) int {
	//设置变量
	//ZeroMemory(cbCardIndex,sizeof(BYTE)*MAX_INDEX);

	//转换扑克
	for i := 0; i < cbCardCount; i++ {
		//ASSERT(lg.IsValidCard(cbCardData[i]));
		cbCardIndex[lg.SwitchToCardIndex(cbCardData[i])]++
	}

	return cbCardCount
}

//有效判断
func (lg *GameLogic) IsValidCard(cbCardData int) bool {
	var cbValue = int(cbCardData & MASK_VALUE)
	var cbColor = int((cbCardData & MASK_COLOR) >> 4)
	return (((cbValue >= 1) && (cbValue <= 9) && (cbColor <= 2)) || ((cbValue >= 1) && (cbValue <= (7 + MAX_HUA_INDEX)) && (cbColor == 3)))
}

//扑克数目
func (lg *GameLogic) GetCardCount(cbCardIndex []int) int {
	//数目统计
	cbCardCount := 0
	for i := 0; i < MAX_INDEX; i++ {
		cbCardCount += cbCardIndex[i]
	}
	return cbCardCount
}

//获取组合
func (lg *GameLogic) GetWeaveCard(cbWeaveKind, cbCenterCard int, cbCardBuffer []int) int {
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

func (lg *GameLogic) AddKindItem(TempKindItem *TagKindItem, KindItem []*TagKindItem, cbKindItemCount int, bMagicThree bool) bool { // todo BYTE &cbKindItemCount, bool &bMagicThree
	TempKindItem.MagicCount = 0
	if lg.MagicIndex == TempKindItem.ValidIndex[0] {
		TempKindItem.MagicCount++
	}
	if lg.MagicIndex == TempKindItem.ValidIndex[1] {
		TempKindItem.MagicCount++
	}
	if lg.MagicIndex == TempKindItem.ValidIndex[2] {
		TempKindItem.MagicCount++
	}

	if TempKindItem.MagicCount >= 3 {
		if !bMagicThree {
			bMagicThree = true
			//CopyMemory(&KindItem[cbKindItemCount++],&TempKindItem,sizeof(TempKindItem));
			return true
		}
		return false
	} else if TempKindItem.MagicCount == 2 {
		cbNoMagicIndex := 0
		cbNoTempMagicIndex := 0
		for i := 0; i < 3; i++ {
			if TempKindItem.ValidIndex[i] != lg.MagicIndex {
				cbNoTempMagicIndex = TempKindItem.ValidIndex[i]
				break
			}
		}
		bFind := false
		for j := 0; j < cbKindItemCount; j++ {
			for i := 0; i < 3; i++ {
				if KindItem[j].ValidIndex[i] != lg.MagicIndex {
					cbNoMagicIndex = KindItem[j].ValidIndex[i]
					break
				}
			}
			if cbNoMagicIndex == cbNoTempMagicIndex && cbNoMagicIndex != 0 {
				bFind = true
			}
		}

		if !bFind {
			util.DeepCopy(&KindItem[cbKindItemCount], &TempKindItem)
			cbKindItemCount++
			return true
		}
		return false
	} else if TempKindItem.MagicCount == 1 {
		cbTempCardIndex := []int{0, 0}
		cbCardIndex := []int{0, 0}
		cbCardCount := 0
		for i := 0; i < 3; i++ {
			if TempKindItem.ValidIndex[i] != lg.MagicIndex {
				cbTempCardIndex[cbCardCount] = TempKindItem.ValidIndex[i]
				cbCardCount++
			}
		}
		//ASSERT(cbCardCount == 2);

		for j := 0; j < cbKindItemCount; j++ {
			if 1 == KindItem[j].MagicCount {
				cbCardCount = 0
				for i := 0; i < 3; i++ {
					if KindItem[j].ValidIndex[i] != lg.MagicIndex {
						cbCardIndex[cbCardCount] = KindItem[j].ValidIndex[i]
						cbCardCount++
					}
				}
				//ASSERT(cbCardCount == 2);

				if cbTempCardIndex[0] == cbCardIndex[0] && cbTempCardIndex[1] == cbCardIndex[1] {
					return false
				}
			}
		}

		util.DeepCopy(&KindItem[cbKindItemCount], &TempKindItem)
		cbKindItemCount++
		return true
	} else {
		util.DeepCopy(&KindItem[cbKindItemCount], &TempKindItem)
		cbKindItemCount++
		return true
	}
}

//分析扑克
func (lg *GameLogic) AnalyseCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount int, TagAnalyseItemArray []*TagAnalyseItem) bool { //todo , CTagAnalyseItemArray & TagAnalyseItemArray
	//计算数目
	cbCardCount := lg.GetCardCount(cbCardIndex)

	//效验数目
	//ASSERT((cbCardCount>=2)&&(cbCardCount<=MAX_COUNT)&&((cbCardCount-2)%3==0));
	if (cbCardCount < 2) || (cbCardCount > MAX_COUNT) || ((cbCardCount-2)%3 != 0) {
		return false
	}

	//变量定义
	cbKindItemCount := 0
	KindItem := make([]*TagKindItem, 27*2+28+16)
	TempKindItem := &TagKindItem{ValidIndex: make([]int, 3)}

	bMagicThree := false

	//需求判断
	cbLessKindItem := int((cbCardCount - 2) / 3)
	//ASSERT((cbLessKindItem+cbWeaveCount)==MAX_WEAVE);

	//单吊判断
	if cbLessKindItem == 0 {
		//效验参数
		//ASSERT((cbCardCount==2)&&(cbWeaveCount==MAX_WEAVE));

		//牌眼判断
		for i := 0; i < MAX_INDEX; i++ {
			if cbCardIndex[i] == 2 || (lg.MagicIndex != MAX_INDEX && i != lg.MagicIndex && cbCardIndex[lg.MagicIndex]+cbCardIndex[i] == 2) {
				//变量定义
				analyseItem := &TagAnalyseItem{}
				//ZeroMemory(&TagAnalyseItem,sizeof(TagAnalyseItem));

				//设置结果
				for j := 0; j < cbWeaveCount; j++ {
					analyseItem.WeaveKind[j] = WeaveItem[j].WeaveKind
					analyseItem.CenterCard[j] = WeaveItem[j].CenterCard
					util.DeepCopy(analyseItem.CardData[j], WeaveItem[j].CardData)
				}
				if cbCardIndex[i] < 2 || i == lg.MagicIndex {
					analyseItem.bMagicEye = true
				} else {
					analyseItem.bMagicEye = false
				}
				if cbCardIndex[i] == 0 {
					analyseItem.CardEye = lg.SwitchToCardData(cbCardIndex[lg.MagicIndex])
				} else {
					analyseItem.CardEye = lg.SwitchToCardData(i)
				}

				//插入结果
				TagAnalyseItemArray = append(TagAnalyseItemArray, analyseItem)
				return true
			}
		}
		return false
	}

	//拆分分析
	cbMagicCardIndex := make([]int, MAX_INDEX)
	util.DeepCopy(&cbMagicCardIndex, &cbCardIndex)

	//如果有财神
	cbMagicCardCount := 0
	cbTempMagicCount := 0

	if lg.MagicIndex != MAX_INDEX {
		cbMagicCardCount = cbCardIndex[lg.MagicIndex]
		//如果财神有代替牌，财神与代替牌转换
		if INDEX_REPLACE_CARD != MAX_INDEX {
			cbMagicCardIndex[lg.MagicIndex] = cbMagicCardIndex[INDEX_REPLACE_CARD]
			cbMagicCardIndex[INDEX_REPLACE_CARD] = cbMagicCardCount
		}
	}

	if cbCardCount >= 3 {
		for i := 0; i < MAX_INDEX-MAX_HUA_INDEX; i++ {
			//同牌判断
			//如果是财神,并且财神数小于3,则不进行组合
			if cbMagicCardIndex[i] >= 3 || (cbMagicCardIndex[i]+cbMagicCardCount >= 3 &&
				((INDEX_REPLACE_CARD != MAX_INDEX && i != INDEX_REPLACE_CARD) || (INDEX_REPLACE_CARD == MAX_INDEX && i != lg.MagicIndex))) {
				nTempIndex := cbMagicCardIndex[i]
				for {
					//ASSERT(cbKindItemCount < len(KindItem));
					cbIndex := int(i)
					cbCenterCard := lg.SwitchToCardData(i)
					//如果是财神且财神有代替牌,则换成代替牌
					if i == lg.MagicIndex && INDEX_REPLACE_CARD != MAX_INDEX {
						cbIndex = INDEX_REPLACE_CARD
						cbCenterCard = lg.SwitchToCardData(INDEX_REPLACE_CARD)
					}
					TempKindItem.WeaveKind = WIK_PENG
					TempKindItem.CenterCard = cbCenterCard
					if nTempIndex > 0 {
						TempKindItem.ValidIndex[0] = cbIndex
					} else {
						TempKindItem.ValidIndex[0] = lg.MagicIndex
					}
					if nTempIndex > 1 {
						TempKindItem.ValidIndex[1] = cbIndex
					} else {
						TempKindItem.ValidIndex[1] = lg.MagicIndex
					}
					if nTempIndex > 2 {
						TempKindItem.ValidIndex[2] = cbIndex
					} else {
						TempKindItem.ValidIndex[2] = lg.MagicIndex
					}

					lg.AddKindItem(TempKindItem, KindItem, cbKindItemCount, bMagicThree)

					//当前索引牌未与财神组合 且财神个数不为0
					if nTempIndex >= 3 && cbMagicCardCount > 0 {
						nTempIndex--
						//1个财神与之组合
						TempKindItem.WeaveKind = WIK_PENG
						TempKindItem.CenterCard = cbCenterCard
						if nTempIndex > 0 {
							TempKindItem.ValidIndex[0] = cbIndex
						} else {
							TempKindItem.ValidIndex[0] = lg.MagicIndex
						}
						if nTempIndex > 1 {
							TempKindItem.ValidIndex[1] = cbIndex
						} else {
							TempKindItem.ValidIndex[1] = lg.MagicIndex
						}
						if nTempIndex > 2 {
							TempKindItem.ValidIndex[2] = cbIndex
						} else {
							TempKindItem.ValidIndex[2] = lg.MagicIndex
						}
						lg.AddKindItem(TempKindItem, KindItem, cbKindItemCount, bMagicThree)

						//两个财神与之组合
						if cbMagicCardCount > 1 {
							TempKindItem.WeaveKind = WIK_PENG
							TempKindItem.CenterCard = cbCenterCard
							if nTempIndex > 0 {
								TempKindItem.ValidIndex[0] = cbIndex
							} else {
								TempKindItem.ValidIndex[0] = lg.MagicIndex
							}
							if nTempIndex > 1 {
								TempKindItem.ValidIndex[1] = cbIndex
							} else {
								TempKindItem.ValidIndex[1] = lg.MagicIndex
							}
							if nTempIndex > 2 {
								TempKindItem.ValidIndex[2] = cbIndex
							} else {
								TempKindItem.ValidIndex[2] = lg.MagicIndex
							}
							lg.AddKindItem(TempKindItem, KindItem, cbKindItemCount, bMagicThree)
						}

						nTempIndex++
					}

					//如果是财神,则退出
					if i == INDEX_REPLACE_CARD || ((i == lg.MagicIndex) && INDEX_REPLACE_CARD == MAX_INDEX) {
						break
					}

					nTempIndex -= 3
					//如果刚好搭配全部，则退出
					if nTempIndex == 0 {
						break
					}

					if nTempIndex+cbMagicCardCount < 3 {
						break
					}
				}
			}

			//连牌判断
			if (i < (MAX_INDEX - MAX_HUA_INDEX - 9)) && ((i % 9) < 7) {
				//只要财神牌数加上3个顺序索引的牌数大于等于3,则进行组合
				if cbMagicCardCount+cbMagicCardIndex[i]+cbMagicCardIndex[i+1]+cbMagicCardIndex[i+2] >= 3 {
					var cbIndex = []int{cbMagicCardIndex[i], cbMagicCardIndex[i+1], cbMagicCardIndex[i+2]}

					if cbIndex[0]+cbIndex[1]+cbIndex[2] == 0 {
						continue
					}

					nMagicCountTemp := cbMagicCardCount

					cbValidIndex := make([]int, 3)
					for {
						if nMagicCountTemp+cbIndex[0]+cbIndex[1]+cbIndex[2] < 3 {
							break
						}
						maxlen := int(len(cbIndex))
						for j := 0; j < maxlen; j++ {
							if cbIndex[j] > 0 {
								cbIndex[j]--
								if (i+j == lg.MagicIndex) && INDEX_REPLACE_CARD != MAX_INDEX {
									cbValidIndex[j] = INDEX_REPLACE_CARD
								} else {
									cbValidIndex[j] = i + j
								}
							} else {
								nMagicCountTemp--
								cbValidIndex[j] = lg.MagicIndex
							}
						}
						if nMagicCountTemp >= 0 {
							//ASSERT(cbKindItemCount < len(KindItem));
							TempKindItem.WeaveKind = WIK_LEFT
							TempKindItem.CenterCard = lg.SwitchToCardData(i)
							util.DeepCopy(&TempKindItem.ValidIndex, &cbValidIndex)
							lg.AddKindItem(TempKindItem, KindItem, cbKindItemCount, bMagicThree)
						} else {
							break
						}
					}
				}
			}
		}
	}

	//组合分析
	if cbKindItemCount >= cbLessKindItem {
		//ASSERT(27*2+28+16 >= cbKindItemCount);
		//变量定义
		cbCardIndexTemp := make([]int, MAX_INDEX)
		//ZeroMemory(cbCardIndexTemp,sizeof(cbCardIndexTemp));

		//变量定义
		cbIndex := make([]int, MAX_WEAVE)
		for i := 0; i < int(len(cbIndex)); i++ {
			cbIndex[i] = i
		}

		pKindItem := make([]*TagKindItem, MAX_WEAVE)
		//ZeroMemory(&pKindItem,sizeof(pKindItem));

		KindItemTemp := make([]*TagKindItem, len(KindItem))

		//开始组合
		for {
			//如果四个组合中的混牌大于手上的混牌个数则重置索引
			cbTempMagicCount = 0
			for i := 0; i < cbLessKindItem; i++ {
				cbTempMagicCount += KindItem[cbIndex[i]].MagicCount
			}
			if cbTempMagicCount <= cbMagicCardCount {

				//设置变量
				util.DeepCopy(&cbCardIndexTemp, &cbCardIndex)
				util.DeepCopy(&KindItemTemp, &KindItem)

				for i := 0; i < cbLessKindItem; i++ {
					pKindItem[i] = KindItemTemp[cbIndex[i]]
				}

				//数量判断
				bEnoughCard := true

				for i := 0; i < cbLessKindItem*3; i++ {
					//存在判断
					cbCardIndex := pKindItem[i/3].ValidIndex[i%3]
					if cbCardIndexTemp[cbCardIndex] == 0 {
						if lg.MagicIndex != MAX_INDEX && cbCardIndexTemp[lg.MagicIndex] > 0 {
							pKindItem[i/3].ValidIndex[i%3] = lg.MagicIndex
							cbCardIndexTemp[lg.MagicIndex]--
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
					//牌眼判断
					cbCardEye := 0
					bMagicEye := false
					if lg.GetCardCount(cbCardIndexTemp) == 2 {
						if lg.MagicIndex != MAX_INDEX && cbCardIndexTemp[lg.MagicIndex] == 2 {
							cbCardEye = lg.SwitchToCardData(lg.MagicIndex)
							bMagicEye = true
						} else {
							for i := 0; i < MAX_INDEX; i++ {
								if cbCardIndexTemp[i] == 2 {
									cbCardEye = lg.SwitchToCardData(i)
									if lg.MagicIndex != MAX_INDEX && i == lg.MagicIndex {
										bMagicEye = true
									}
									break
								} else if i != lg.MagicIndex && lg.MagicIndex != MAX_INDEX && cbCardIndexTemp[i]+cbCardIndexTemp[lg.MagicIndex] == 2 {
									cbCardEye = lg.SwitchToCardData(i)
									bMagicEye = true
									break
								}
							}
						}
					}

					//组合类型
					if cbCardEye != 0 {
						//变量定义
						analyseItem := &TagAnalyseItem{}
						//ZeroMemory(&TagAnalyseItem,sizeof(TagAnalyseItem));

						//设置组合
						for i := 0; i < cbWeaveCount; i++ {
							analyseItem.WeaveKind[i] = WeaveItem[i].WeaveKind
							analyseItem.CenterCard[i] = WeaveItem[i].CenterCard
							lg.GetWeaveCard(WeaveItem[i].WeaveKind, WeaveItem[i].CenterCard, analyseItem.CardData[i])
						}

						//设置牌型
						for i := 0; i < cbLessKindItem; i++ {
							analyseItem.WeaveKind[i+cbWeaveCount] = pKindItem[i].WeaveKind
							analyseItem.CenterCard[i+cbWeaveCount] = pKindItem[i].CenterCard
							analyseItem.CardData[cbWeaveCount+i][0] = lg.SwitchToCardData(pKindItem[i].ValidIndex[0])
							analyseItem.CardData[cbWeaveCount+i][1] = lg.SwitchToCardData(pKindItem[i].ValidIndex[1])
							analyseItem.CardData[cbWeaveCount+i][2] = lg.SwitchToCardData(pKindItem[i].ValidIndex[2])
						}

						//设置牌眼
						analyseItem.CardEye = cbCardEye
						analyseItem.bMagicEye = bMagicEye

						//插入结果
						TagAnalyseItemArray = append(TagAnalyseItemArray, analyseItem)
					}
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

	return (len(TagAnalyseItemArray) > 0)
}

/*
// 胡法分析函数
*/

//碰碰和
func (lg *GameLogic) IsPengPeng(pTagAnalyseItem *TagAnalyseItem) bool {
	for i := 0; i < int(len(pTagAnalyseItem.WeaveKind)); i++ {
		if (pTagAnalyseItem.WeaveKind[i] & (WIK_LEFT | WIK_CENTER | WIK_RIGHT)) != 0 {
			return false
		}
	}
	return true
}

//是否麻七系列
func (lg *GameLogic) IsMaQi(cbCardIndex []int, cbWeaveCount int, ChiHuRight int64) bool {
	if cbWeaveCount != 0 {
		return false
	}

	cbGang := 0
	var cbMagicCount int
	if lg.MagicIndex != MAX_INDEX {
		cbMagicCount = cbCardIndex[lg.MagicIndex]
	}

	//变量定义
	for i := 0; i < MAX_INDEX; i++ {
		if cbCardIndex[i] != 0 && i != lg.MagicIndex {
			if cbCardIndex[i]%2 == 1 {
				if cbMagicCount >= 1 {
					cbMagicCount--
				} else {
					return false //有非对子，跳出
				}
			}
			cbGang += cbCardIndex[i] / 4
		}
	}

	if cbGang >= 2 { //手上有两个4张，王中王
		ChiHuRight |= CHR_MA_QI_WZW
		return true
	} else if cbGang == 1 { //有一个4张，麻七王
		ChiHuRight |= CHR_MA_QI_WANG
		return true
	} else { //麻七对
		ChiHuRight |= CHR_MA_QI_DUI
		return true
	}

	return false
}

//十三烂系列
func (lg *GameLogic) IsShiSanLan(cbCardIndex []int, cbWeaveCount int, ChiHuRight int64) bool {
	//组合判断
	if cbWeaveCount != 0 {
		return false
	}

	for i := 0; i < MAX_INDEX; i++ {
		if cbCardIndex[i] >= 2 { //不能有重复牌
			return false
		}

	}

	for j := 0; j < 3; j++ {
		for i := 0; i < 9-2; i++ {
			index := j*9 + i
			if cbCardIndex[index]+cbCardIndex[index+1]+cbCardIndex[index+2] > 1 {
				//if(cbCardIndex[index+1]>0 || cbCardIndex[index+2]>0)//间隔必须>=3
				return false
			}
		}
	}

	for i := 27; i < MAX_INDEX; i++ { //检查风牌
		if cbCardIndex[i] == 0 { //没有包含所有风牌，十三烂
			ChiHuRight |= CHR_SHI_SAN_LAN
			return true
		}
	}

	//所有风牌都有，七星十三烂
	ChiHuRight |= CHR_QX_SHI_SAN_LAN
	return true
}

//鸡胡
func (lg *GameLogic) IsJiHu(pTagAnalyseItem *TagAnalyseItem) bool {
	bPeng := false
	bLian := false
	for i := 0; i < int(len(pTagAnalyseItem.WeaveKind)); i++ {
		if (pTagAnalyseItem.WeaveKind[i] & (WIK_PENG | WIK_GANG)) != 0 {
			bPeng = true
		} else {
			bLian = true
		}
	}

	return bPeng && bLian
}

//平胡
func (lg *GameLogic) IsPingHu(pTagAnalyseItem *TagAnalyseItem) bool {
	//检查组合
	for i := 0; i < int(len(pTagAnalyseItem.WeaveKind)); i++ {
		if (pTagAnalyseItem.WeaveKind[i] & (WIK_PENG | WIK_GANG)) != 0 {
			return false
		}
	}
	return true
}

//清一色
func (lg *GameLogic) IsQingYiSe(pTagAnalyseItem *TagAnalyseItem, bQuanFan bool) bool {
	//参数校验
	if pTagAnalyseItem == nil {
		return false
	}

	//变量定义
	cbCardColor := pTagAnalyseItem.CardEye & MASK_COLOR
	for i := 0; i < MAX_WEAVE; i++ {
		if (pTagAnalyseItem.CenterCard[i] & MASK_COLOR) != cbCardColor {
			return false
		}
	}

	if 0x30 == cbCardColor {
		bQuanFan = true
	} else {
		bQuanFan = false
	}

	return true
}

//混一色
func (lg *GameLogic) IsHunYiSe(pTagAnalyseItem *TagAnalyseItem) bool {
	//参数校验
	if pTagAnalyseItem == nil {
		return false
	}

	//变量定义
	cbCardColor := (pTagAnalyseItem.CardEye & MASK_COLOR) >> 4
	//ASSERT(cbCardColor >= 0 && cbCardColor <= 3);
	cbColorCount := make([]int, 4)
	cbColorCount[cbCardColor] = 1
	for i := 0; i < MAX_WEAVE; i++ {
		cbCardColor = ((pTagAnalyseItem.CenterCard[i]) & MASK_COLOR) >> 4
		//ASSERT(cbCardColor >= 0 && cbCardColor <= 3);
		if 0 == cbColorCount[cbCardColor] {
			cbColorCount[cbCardColor] = 1
		}
	}

	if cbColorCount[0]+cbColorCount[1]+cbColorCount[2] == 1 && cbColorCount[3] == 1 {
		return true
	}

	return false
}
