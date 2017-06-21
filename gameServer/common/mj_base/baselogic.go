package mj_base

import (
	"mj/common/msg"
	"mj/gameServer/common"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

type mj_logic interface {
}

type BaseLogic struct {
	CardDataArray []int //扑克数据
	MagicIndex    int   //钻牌索引
	ReplaceCard   int   //替换金牌的牌
	SwitchInx     func(int) int
	CheckValid    func(int) bool
}

func NewBaseLogic(cb1 func(int) int, cb2 func(int) bool) *BaseLogic {
	bl := new(BaseLogic)
	bl.SwitchInx = cb1
	bl.CheckValid = cb2
	return bl
}

//////////// 上面函数必须重写  下面通用函数

//混乱扑克
func (lg *BaseLogic) RandCardList(cbCardBuffer, OriDataArray []int) {
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

//有效判断
func (lg *BaseLogic) IsValidCard(cbCardData int) bool {
	var cbValue = int(cbCardData & MASK_VALUE)
	var cbColor = int((cbCardData & MASK_COLOR) >> 4)
	return ((cbValue >= 1) && (cbValue <= 9) && (cbColor <= 2)) || ((cbValue >= 1) && (cbValue <= (7 + MAX_HUA_INDEX)) && (cbColor == 3))
}

//删除扑克
func (lg *BaseLogic) RemoveCardByArr(cbCardIndex, cbRemoveCard []int) bool {
	//参数校验
	for _, card := range cbRemoveCard {
		//效验扑克
		if lg.IsValidCard(card) {
			return false
		}

		if cbCardIndex[lg.SwitchToCardIndex(card)] <= 0 {
			return false
		}
	}
	//删除扑克
	for _, card := range cbRemoveCard {
		//删除扑克
		cbCardIndex[lg.SwitchToCardIndex(card)]--
	}

	return true
}

//删除扑克
func (lg *BaseLogic) RemoveCard(cbCardIndex []int, cbRemoveCard int) bool {
	//删除扑克
	cbRemoveIndex := lg.SwitchInx(cbRemoveCard)
	//效验扑克
	if !lg.CheckValid(cbRemoveCard) {
		log.Error("at RemoveCard card is Invalid %d", cbRemoveCard)
	}
	if cbCardIndex[lg.SwitchInx(cbRemoveCard)] < 0 {
		log.Error("at RemoveCard 11 card is Invalid %d", cbRemoveCard)
	}
	if cbCardIndex[cbRemoveIndex] > 0 {
		cbCardIndex[cbRemoveIndex]--
		return true
	}

	return false
}

//扑克数目
func (lg *BaseLogic) GetCardCount(cbCardIndex []int) int {
	//数目统计
	cbCardCount := 0
	for i := 0; i < MAX_INDEX; i++ {
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
	if cbCardIndex[lg.SwitchInx(cbCurrentCard)] >= 2 {
		return WIK_PENG
	}

	return WIK_NULL
}

//杠牌判断
func (lg *BaseLogic) EstimateGangCard(cbCardIndex []int, cbCurrentCard int) int {
	if cbCardIndex[lg.SwitchInx(cbCurrentCard)] == 3 {
		return WIK_GANG
	}

	return WIK_NULL
}

func (lg *BaseLogic) GetCardColor(cbCardData int) int { return cbCardData & MASK_COLOR }
func (lg *BaseLogic) GetCardValue(cbCardData int) int { return cbCardData & MASK_VALUE }

//扑克转换
func (lg *BaseLogic) SwitchToCardIndex(cbCardData int) int {
	//计算位置
	cbValue := cbCardData & MASK_VALUE
	cbColor := (cbCardData & MASK_COLOR) >> 4

	if cbColor >= 0x03 {
		return cbValue + 27 - 1
	}
	return cbColor*9 + cbValue - 1
}

//吃胡分析
func (lg *BaseLogic) AnalyseChiHuCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbCurrentCard int, ChiHuRight int, b4HZHu bool) int {

	cbWeaveCount := len(WeaveItem)
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
		cbCardIndexTemp[lg.SwitchToCardIndex(cbCurrentCard)]++
	}

	if b4HZHu && cbCardIndexTemp[31] == 4 { //四个红中直接胡牌
		return WIK_CHI_HU
	}
	//分析扑克
	_, TagAnalyseItemArray = lg.AnalyseCard(cbCardIndexTemp, WeaveItem, cbWeaveCount, TagAnalyseItemArray)

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

func (lg *BaseLogic) AnalyseGangCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbWeaveCount, cbProvideCard int, gangCardResult *common.TagGangCardResult) int {

	//设置变量
	cbActionMask := WIK_NULL
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

func (lg *BaseLogic) AnalyseTingCard(cbCardIndex []int, WeaveItem []*msg.WeaveItem, cbOutCardData, cbHuCardCount []int, cbHuCardData [][]int) int {

	cbOutCount := 0
	cbWeaveCount := len(WeaveItem)
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
