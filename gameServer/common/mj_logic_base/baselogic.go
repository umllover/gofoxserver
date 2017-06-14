package mj_logic_base

import (
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

type BaseLogic struct {
	SwitchInx  func(int) int
	CheckValid func(int) bool
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

//删除扑克
func (lg *BaseLogic) RemoveCardByArray(cbCardIndex []int, cbRemoveCard []int) bool {
	//参数校验
	for i := 0; i < len(cbRemoveCard); i++ {
		//效验扑克
		if lg.CheckValid(cbRemoveCard[i]) {
			return false
		}

		if cbCardIndex[lg.SwitchInx(cbRemoveCard[i])] <= 0 {
			return false
		}
	}
	//删除扑克
	for i := 0; i < len(cbRemoveCard); i++ {
		//删除扑克
		cbCardIndex[lg.SwitchInx(cbRemoveCard[i])]--
	}
	return true
}

//删除扑克
func (lg *BaseLogic) RemoveCardByCnt(cbCardIndex, cbRemoveCard []int, cbRemoveCount int) bool {
	//参数校验
	for i := 0; i < len(cbRemoveCard); i++ {
		//效验扑克
		if lg.CheckValid(cbRemoveCard[i]) {
			return false
		}

		if cbCardIndex[lg.SwitchInx(cbRemoveCard[i])] <= 0 {
			return false
		}
	}
	//删除扑克
	for i := 0; i < cbRemoveCount; i++ {
		//删除扑克
		cbRemoveIndex := lg.SwitchInx(cbRemoveCard[i])
		//删除扑克
		cbCardIndex[cbRemoveIndex]--
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
