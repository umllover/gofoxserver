package room

import (
	"mj/gameServer/common/pk/pk_base"

	"github.com/lovelly/leaf/log"
)

func NewDDZLogic(ConfigIdx int) *ddz_logic {
	l := new(ddz_logic)
	l.BaseLogic = pk_base.NewBaseLogic(ConfigIdx)
	return l
}

type ddz_logic struct {
	*pk_base.BaseLogic
}

const (
	// 牌类型
	CT_ERROR          = 0  // 错误类型
	CT_SINGLE         = 1  // 单牌类型
	CT_DOUBLE         = 2  // 对牌类型
	CT_THREE          = 3  // 三条类型
	CT_SINGLE_LINE    = 4  // 单连类型
	CT_DOUBLE_LINE    = 5  // 对连类型
	CT_THREE_LINE     = 6  // 三连类型
	CT_THREE_TAKE_ONE = 7  // 三带一单
	CT_THREE_TAKE_TWO = 8  // 三带一对
	CT_FOUR_TAKE_ONE  = 9  // 四带两单
	CT_FOUR_TAKE_TWO  = 10 // 四带两对
	CT_BOMB_CARD      = 11 // 炸弹类型
	CT_MISSILE_CARD   = 12 // 火箭类型

	// 数目定义
	MAX_COUNT  = 20 //最大数目
	FULL_COUNT = 54 //全牌数目

	// 逻辑数目
	NORMAL_COUNT    = 17 //常规数目
	DISPATCH_COUNT  = 51 //派发数目
	GOOD_CARD_COUTN = 38 //好牌数目

	// 排序类型
	ST_ORDER  = 1 //大小排序
	ST_COUNT  = 2 //数目排序
	ST_CUSTOM = 3 //自定排序

	// 索引变量
	cbIndexCount = 5

	// 游戏人数
	GAME_PLAYER = 3
)

type BaseLogic struct {
	CardDataArray []int //扑克数据
	MagicIndex    int   //钻牌索引
	ReplaceCard   int   //替换金牌的牌
	SwitchToIdx   func(int) int
	CheckValid    func(int) bool
	SwitchToCard  func(int) int
	ConfigIdx     int //配置文件索引
}

//分析结构
type tagAnalyseResult struct {
	cbBlockCount [4]uint8            //扑克数目
	cbCardData   [4][MAX_COUNT]uint8 //扑克数据
}

//出牌结果
type tagOutCardResult struct {
	cbCardCount  uint8            //扑克数目
	cbResultCard [MAX_COUNT]uint8 //结果扑克
}

//分布信息
type tagDistributing struct {
	cbCardCount    uint8        //扑克数目
	cbDistributing [15][6]uint8 //分布信息
}

//搜索结果
type tagSearchCardResult struct {
	cbSearchCount uint8                       //结果数目
	cbCardCount   [MAX_COUNT]uint8            //扑克数目
	cbResultCard  [MAX_COUNT][MAX_COUNT]uint8 //结果扑克
}

//获取类型
func (dg *ddz_logic) GetCardType(cbCardData []uint8, cbCardCount uint8) uint8 {
	//	简单牌型
	switch cbCardCount {
	case 0: //空牌
		return CT_ERROR
	case 1: //单牌
		return CT_SINGLE
	case 2: //对牌火箭
		{
			//牌型判断
			if (cbCardData[0] == 0x4F) && (cbCardData[1] == 0x4E) {
				return CT_MISSILE_CARD
			}
			if dg.GetCardLogicValue(cbCardData[0]) == dg.GetCardLogicValue(cbCardData[1]) {
				return CT_DOUBLE
			}

			return CT_ERROR
		}
	}

	//分析扑克
	var AnalyseResult tagAnalyseResult
	dg.AnalysebCardData(cbCardData, cbCardCount, &AnalyseResult)

	//四牌判断
	if AnalyseResult.cbBlockCount[3] > 0 {
		//牌型判断
		if (AnalyseResult.cbBlockCount[3] == 1) && (cbCardCount == 4) {
			return CT_BOMB_CARD
		}
		if (AnalyseResult.cbBlockCount[3] == 1) && (cbCardCount == 6) {
			return CT_FOUR_TAKE_ONE
		}
		if (AnalyseResult.cbBlockCount[3] == 1) &&
			(cbCardCount == 8) &&
			(AnalyseResult.cbBlockCount[1] == 2) {
			return CT_FOUR_TAKE_TWO
		}

		return CT_ERROR
	}

	// 三牌判断
	if AnalyseResult.cbBlockCount[2] > 0 {
		// 连牌判断
		if AnalyseResult.cbBlockCount[2] > 1 {
			//变量定义
			cbCardData := AnalyseResult.cbCardData[2][0]
			cbFirstLogicValue := dg.GetCardLogicValue(cbCardData)

			// 错误过虑
			if cbFirstLogicValue >= 15 {
				return CT_ERROR
			}

			//连牌判断
			for i := 1; uint8(i) < AnalyseResult.cbBlockCount[2]; i++ {
				cbCardData := AnalyseResult.cbCardData[2][i*3]
				if cbFirstLogicValue != (dg.GetCardLogicValue(cbCardData) + uint8(i)) {
					return CT_ERROR
				}
			}
		} else if cbCardCount == 3 {
			return CT_THREE
		}

		// 牌形判断
		if AnalyseResult.cbBlockCount[2]*3 == cbCardCount {
			return CT_THREE_LINE
		}
		if AnalyseResult.cbBlockCount[2]*4 == cbCardCount {
			return CT_THREE_TAKE_ONE
		}
		if (AnalyseResult.cbBlockCount[2]*5 == cbCardCount) &&
			(AnalyseResult.cbBlockCount[1] == AnalyseResult.cbBlockCount[2]) {
			return CT_THREE_TAKE_TWO
		}

		return CT_ERROR
	}

	//两张类型
	if AnalyseResult.cbBlockCount[1] >= 3 {
		// 变量定义
		cbCardData := AnalyseResult.cbCardData[1][0]
		cbFirstLogicValue := dg.GetCardLogicValue(cbCardData)

		// 错误过虑
		if cbFirstLogicValue >= 15 {
			return CT_ERROR
		}

		// 连牌判断
		for i := 1; i < int(AnalyseResult.cbBlockCount[1]); i++ {
			cbCardData := AnalyseResult.cbCardData[1][i*2]
			if cbFirstLogicValue != (dg.GetCardLogicValue(cbCardData) + uint8(i)) {
				return CT_ERROR
			}
		}

		// 二连判断
		if (AnalyseResult.cbBlockCount[1] * 2) == cbCardCount {
			return CT_DOUBLE_LINE
		}

		return CT_ERROR
	}

	//单张判断
	if (AnalyseResult.cbBlockCount[0] >= 5) && (AnalyseResult.cbBlockCount[0] == cbCardCount) {
		// 变量定义
		cbCardData := AnalyseResult.cbCardData[0][0]
		cbFirstLogicValue := dg.GetCardLogicValue(cbCardData)

		//错误过虑
		if cbFirstLogicValue >= 15 {
			return CT_ERROR
		}

		//连牌判断
		for i := 1; uint8(i) < AnalyseResult.cbBlockCount[0]; i++ {
			cbCardData := AnalyseResult.cbCardData[0][i]
			if cbFirstLogicValue != (dg.GetCardLogicValue(cbCardData) + uint8(i)) {
				return CT_ERROR
			}
		}

		return CT_SINGLE_LINE
	}

	return CT_ERROR
}

//排列扑克
func (dg *ddz_logic) SortCardList([]int, int) {
	return
}

func (dg *ddz_logic) DDZSortCardList(cbCardData []uint8, cbCardCount uint8, cbSortType uint8) {
	// 数目过虑
	if cbCardCount == 0 {
		return
	}
	if cbSortType == ST_CUSTOM {
		return
	}
	// 转换数值
	var cbSortValue [MAX_COUNT]uint8
	for i := 0; uint8(i) < cbCardCount; i++ {
		cbSortValue[i] = dg.GetCardLogicValue(cbCardData[i])
	}

	// 排序操作
	bSorted := true
	var cbSwitchData uint8
	cbLast := cbCardCount - 1
	for bSorted {
		bSorted = true
		for i := 0; i < int(cbLast); i++ {
			if (cbSortValue[i] < cbSortValue[i+1]) ||
				((cbSortValue[i] == cbSortValue[i+1]) &&
					(cbCardData[i] < cbCardData[i+1])) {
				// 设置标志
				bSorted = false

				//扑克数据
				cbSwitchData = cbCardData[i]
				cbCardData[i] = cbCardData[i+1]
				cbCardData[i+1] = cbSwitchData

				//排序权位
				cbSwitchData = cbSortValue[i]
				cbSortValue[i] = cbSortValue[i+1]
				cbSortValue[i+1] = cbSwitchData
			}
		}
		cbLast--
	}
}

//删除扑克
func (dg *ddz_logic) RemoveCardList(cbRemoveCard []uint8, cbRemoveCount uint8, cbCardData []uint8) bool {
	// 检验数据
	if cbRemoveCount > uint8(len(cbCardData)) {
		log.Error("要删除的扑克数%i大于已有扑克数%i", cbRemoveCount, len(cbCardData))
		return false
	}

	// 备份
	var tmpCardData []uint8
	copy(tmpCardData, cbCardData)

	var u8DeleteCount uint8 // 记录删除记录

	for _, v1 := range cbRemoveCard {
		for j, v2 := range cbCardData {
			if v1 == v2 {
				copy(cbCardData[j:], cbCardData[j+1:])
				u8DeleteCount++
			}
		}
	}

	if u8DeleteCount != cbRemoveCount {
		// 删除数量不一，恢复数据
		log.Error("实际删除数量%与需要删除数量%i不一样", u8DeleteCount, cbRemoveCount)
		copy(cbCardData, tmpCardData)
		return false
	}

	return true
}

//删除扑克
func (dg *ddz_logic) RemoveCard(cbRemoveCard []uint8, cbRemoveCount uint8, cbCardData []uint8, cbCardCount uint8) bool {
	return dg.RemoveCardList(cbRemoveCard, cbRemoveCount, cbCardData)
}

// 排列出牌扑克
func (dg *ddz_logic) SortOutCardList(cbCardData []uint8, cbCardCount uint8) {

	// 获取牌型
	cbCardType := dg.GetCardType(cbCardData, cbCardCount)

	if cbCardType == CT_THREE_TAKE_ONE || cbCardType == CT_THREE_TAKE_TWO {
		//分析牌
		var AnalyseResult tagAnalyseResult
		dg.AnalysebCardData(cbCardData, cbCardCount, &AnalyseResult)

		cbCardCount = AnalyseResult.cbBlockCount[2] * 3
		copy(cbCardData, AnalyseResult.cbCardData[2][:cbCardCount])
		for i := 3; i >= 0; i-- {
			if i == 2 {
				continue
			}

			if AnalyseResult.cbBlockCount[i] > 0 {
				copy(cbCardData[cbCardCount:], AnalyseResult.cbCardData[i][:(i+1)*int(AnalyseResult.cbBlockCount[i])])
				cbCardCount += uint8(i+1) * AnalyseResult.cbBlockCount[i]
			}
		}
	} else if cbCardType == CT_FOUR_TAKE_ONE || cbCardType == CT_FOUR_TAKE_TWO {
		//分析牌
		var AnalyseResult tagAnalyseResult
		dg.AnalysebCardData(cbCardData, cbCardCount, &AnalyseResult)

		cbCardCount = AnalyseResult.cbBlockCount[3] * 4
		copy(cbCardData, AnalyseResult.cbCardData[3][:cbCardCount])
		for i := 3; i >= 0; i-- {
			if i == 3 {
				continue
			}

			if AnalyseResult.cbBlockCount[i] > 0 {
				copy(cbCardData[cbCardCount:], AnalyseResult.cbCardData[i][:uint8(i+1)*AnalyseResult.cbBlockCount[i]])
				cbCardCount += uint8(i+1) * AnalyseResult.cbBlockCount[i]
			}
		}
	}

	return
}

//逻辑数值
func (dg *ddz_logic) GetCardLogicValue(cbCardData uint8) uint8 {
	// 扑克属性
	cbCardColor := uint8(dg.GetCardColor(int(cbCardData)))
	cbCardValue := uint8(dg.GetCardValue(int(cbCardData)))

	if cbCardValue <= 0 || cbCardValue > (pk_base.LOGIC_MASK_VALUE&0x4f) {
		log.Error("求取逻辑数值的扑克数据有误%i", cbCardValue)
		return 0
	}

	// 转换数值
	if cbCardColor == 0x40 {
		return cbCardValue + 2
	}

	if cbCardValue <= 2 {
		return cbCardValue + 13
	} else {
		return cbCardValue
	}
}

//对比扑克
func (dg *ddz_logic) CompareCard(cbFirstCard []uint8, cbNextCard []uint8, cbFirstCount uint8, cbNextCount uint8) bool {
	// 获取类型
	cbNextType := dg.GetCardType(cbNextCard, cbNextCount)
	cbFirstType := dg.GetCardType(cbFirstCard, cbFirstCount)

	// 类型判断
	if cbNextType == CT_ERROR {
		return false
	}
	if cbNextType == CT_MISSILE_CARD {
		return true
	}

	// 炸弹判断
	if (cbFirstType != CT_BOMB_CARD) && (cbNextType == CT_BOMB_CARD) {
		return true
	}
	if (cbFirstType == CT_BOMB_CARD) && (cbNextType != CT_BOMB_CARD) {
		return false
	}

	// 规则判断
	if (cbFirstType != cbNextType) || (cbFirstCount != cbNextCount) {
		return false
	}

	// 开始对比
	switch cbNextType {
	case CT_SINGLE:
	case CT_DOUBLE:
	case CT_THREE:
	case CT_SINGLE_LINE:
	case CT_DOUBLE_LINE:
	case CT_THREE_LINE:
	case CT_BOMB_CARD:
		{
			// 获取数值
			cbNextLogicValue := dg.GetCardLogicValue(cbNextCard[0])
			cbFirstLogicValue := dg.GetCardLogicValue(cbFirstCard[0])

			// 对比扑克
			return cbNextLogicValue > cbFirstLogicValue
		}
	case CT_THREE_TAKE_ONE:
	case CT_THREE_TAKE_TWO:
		{
			// 分析扑克
			var NextResult tagAnalyseResult
			var FirstResult tagAnalyseResult
			dg.AnalysebCardData(cbNextCard, cbNextCount, &NextResult)
			dg.AnalysebCardData(cbFirstCard, cbFirstCount, &FirstResult)

			// 获取数值
			cbNextLogicValue := dg.GetCardLogicValue(NextResult.cbCardData[2][0])
			cbFirstLogicValue := dg.GetCardLogicValue(FirstResult.cbCardData[2][0])

			// 对比扑克
			return cbNextLogicValue > cbFirstLogicValue
		}
	case CT_FOUR_TAKE_ONE:
	case CT_FOUR_TAKE_TWO:
		{
			// 分析扑克
			var NextResult tagAnalyseResult
			var FirstResult tagAnalyseResult
			dg.AnalysebCardData(cbNextCard, cbNextCount, &NextResult)
			dg.AnalysebCardData(cbFirstCard, cbFirstCount, &FirstResult)

			//获取数值
			cbNextLogicValue := dg.GetCardLogicValue(NextResult.cbCardData[3][0])
			cbFirstLogicValue := dg.GetCardLogicValue(FirstResult.cbCardData[3][0])

			//对比扑克
			return cbNextLogicValue > cbFirstLogicValue
		}
	}
	return false
}

//构造扑克
func (dg *ddz_logic) MakeCardData(cbValueIndex uint8, cbColorIndex uint8) uint8 {
	return (cbColorIndex << 4) | (cbValueIndex + 1)
}

//分析扑克
func (dg *ddz_logic) AnalysebCardData(cbCardData []uint8, cbCardCount uint8, AnalyseResult *tagAnalyseResult) {
	// 设置结果
	//ZeroMemory(&AnalyseResult,sizeof(AnalyseResult));

	// 扑克分析
	for i := 0; uint8(i) < cbCardCount; i++ {
		// 变量定义
		cbSameCount := 1
		cbLogicValue := dg.GetCardLogicValue(cbCardData[i])

		// 搜索同牌
		for j := i + 1; uint8(j) < cbCardCount; j++ {
			// 获取扑克
			if dg.GetCardLogicValue(cbCardData[j]) != cbLogicValue {
				break
			}

			// 设置变量
			cbSameCount++
		}

		if cbSameCount > 4 {
			// 设置结果
			//ZeroMemory(&AnalyseResult, sizeof(AnalyseResult));
			log.Error("相同数量不可能大于4")
			return
		}

		// 设置结果
		cbIndex := AnalyseResult.cbBlockCount[cbSameCount-1]
		AnalyseResult.cbBlockCount[cbSameCount-1]++
		for j := 0; j < cbSameCount; j++ {
			AnalyseResult.cbCardData[cbSameCount-1][int(cbIndex)*cbSameCount+j] = cbCardData[i+j]
		}

		// 设置索引
		i += cbSameCount - 1
	}
}

// 分析分布
func (dg *ddz_logic) AnalysebDistributing(cbCardData []uint8, cbCardCount uint8, Distributing *tagDistributing) {
	// 设置变量
	//	ZeroMemory(&Distributing,sizeof(Distributing));

	//设置变量
	for i := 0; uint8(i) < cbCardCount; i++ {
		if cbCardData[i] == 0 {
			continue
		}

		//获取属性
		cbCardColor := dg.GetCardColor(int(cbCardData[i]))
		cbCardValue := dg.GetCardValue(int(cbCardData[i]))

		//分布信息
		Distributing.cbCardCount++
		Distributing.cbDistributing[cbCardValue-1][cbIndexCount]++
		Distributing.cbDistributing[cbCardValue-1][cbCardColor>>4]++
	}
	return
}

//出牌搜索
func (dg *ddz_logic) SearchOutCard(cbHandCardData []uint8, cbHandCardCount uint8, cbTurnCardData []uint8, cbTurnCardCount uint8, pSearchCardResult *tagSearchCardResult) uint8 {
	// 设置结果
	if pSearchCardResult == nil {
		return 0
	}

	// 变量定义
	var cbResultCount uint8
	var tmpSearchCardResult tagSearchCardResult

	// 构造扑克
	var cbCardData [MAX_COUNT]uint8
	cbCardCount := cbHandCardCount
	copy(cbCardData[:], cbHandCardData[:cbHandCardCount])

	// 排列扑克
	dg.DDZSortCardList(cbCardData[:], cbCardCount, ST_ORDER)

	// 获取类型
	cbTurnOutType := dg.GetCardType(cbTurnCardData, cbTurnCardCount)

	// 出牌分析
	switch cbTurnOutType {
	case CT_ERROR:
		{ //错误类型
			// 提取各种牌型一组

			// 是否一手出完
			if dg.GetCardType(cbCardData[:], cbCardCount) != CT_ERROR {
				pSearchCardResult.cbCardCount[cbResultCount] = cbCardCount
				copy(pSearchCardResult.cbResultCard[cbResultCount][:], cbCardData[:cbCardCount])
				cbResultCount++
			}

			// 如果最小牌不是单牌，则提取
			var cbSameCount uint8
			if cbCardCount > 1 && dg.GetCardValue(int(cbCardData[cbCardCount-1])) == dg.GetCardValue(int(cbCardData[cbCardCount-2])) {
				cbSameCount = 1
				pSearchCardResult.cbResultCard[cbResultCount][0] = cbCardData[cbCardCount-1]
				cbCardValue := dg.GetCardValue(int(cbCardData[cbCardCount-1]))
				for i := cbCardCount - 2; i >= 0; i-- {
					if dg.GetCardValue(int(cbCardData[i])) == cbCardValue {
						pSearchCardResult.cbResultCard[cbResultCount][cbSameCount] = cbCardData[i]
						cbSameCount++
					} else {
						break
					}
				}

				pSearchCardResult.cbCardCount[cbResultCount] = cbSameCount
				cbResultCount++
			}

			// 单牌
			var cbTmpCount uint8
			if cbSameCount != 1 {
				cbTmpCount = dg.SearchSameCard(cbCardData[:], cbCardCount, 0, 1, &tmpSearchCardResult)
				if cbTmpCount > 0 {
					pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0]
					copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[0][:tmpSearchCardResult.cbCardCount[0]])
					cbResultCount++
				}
			}

			// 对牌
			if cbSameCount != 2 {
				cbTmpCount = dg.SearchSameCard(cbCardData[:], cbCardCount, 0, 2, &tmpSearchCardResult)
				if cbTmpCount > 0 {
					pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0]
					copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[0][:tmpSearchCardResult.cbCardCount[0]])
					cbResultCount++
				}
			}

			// 三条
			if cbSameCount != 3 {
				cbTmpCount = dg.SearchSameCard(cbCardData[:], cbCardCount, 0, 3, &tmpSearchCardResult)
				if cbTmpCount > 0 {
					pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0]
					copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[0][:tmpSearchCardResult.cbCardCount[0]])
					cbResultCount++
				}
			}

			// 三带一单
			cbTmpCount = dg.SearchTakeCardType(cbCardData[:], cbCardCount, 0, 3, 1, &tmpSearchCardResult)
			if cbTmpCount > 0 {
				pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0]
				copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[0][:tmpSearchCardResult.cbCardCount[0]])
				cbResultCount++
			}

			// 三带一对
			cbTmpCount = dg.SearchTakeCardType(cbCardData[:], cbCardCount, 0, 3, 2, &tmpSearchCardResult)
			if cbTmpCount > 0 {
				pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0]
				copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[0][:tmpSearchCardResult.cbCardCount[0]])
				cbResultCount++
			}

			//单连
			cbTmpCount = dg.SearchLineCardType(cbCardData[:], cbCardCount, 0, 1, 0, &tmpSearchCardResult)
			if cbTmpCount > 0 {
				pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0]
				copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[0][:tmpSearchCardResult.cbCardCount[0]])
				cbResultCount++
			}

			// 连对
			cbTmpCount = dg.SearchLineCardType(cbCardData[:], cbCardCount, 0, 2, 0, &tmpSearchCardResult)
			if cbTmpCount > 0 {
				pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0]
				copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[0][:tmpSearchCardResult.cbCardCount[0]])
				cbResultCount++
			}

			// 三连
			cbTmpCount = dg.SearchLineCardType(cbCardData[:], cbCardCount, 0, 3, 0, &tmpSearchCardResult)
			if cbTmpCount > 0 {
				pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0]
				copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[0][:tmpSearchCardResult.cbCardCount[0]])
				cbResultCount++
			}

			////飞机
			//cbTmpCount = SearchThreeTwoLine( cbCardData,cbCardCount,&tmpSearchCardResult );
			//if( cbTmpCount > 0 )
			//{
			//	pSearchCardResult->cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0];
			//	CopyMemory( pSearchCardResult->cbResultCard[cbResultCount],tmpSearchCardResult.cbResultCard[0],
			//		sizeof(BYTE)*tmpSearchCardResult.cbCardCount[0] );
			//	cbResultCount++;
			//}

			// 炸弹
			if cbSameCount != 4 {
				cbTmpCount = dg.SearchSameCard(cbCardData[:], cbCardCount, 0, 4, &tmpSearchCardResult)
				if cbTmpCount > 0 {
					pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[0]
					copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[0][:tmpSearchCardResult.cbCardCount[0]])
					cbResultCount++
				}
			}

			// 搜索火箭
			if (cbCardCount >= 2) && (cbCardData[0] == 0x4F) && (cbCardData[1] == 0x4E) {
				// 设置结果
				pSearchCardResult.cbCardCount[cbResultCount] = 2
				pSearchCardResult.cbResultCard[cbResultCount][0] = cbCardData[0]
				pSearchCardResult.cbResultCard[cbResultCount][1] = cbCardData[1]
				cbResultCount++
			}

			pSearchCardResult.cbSearchCount = cbResultCount
			return cbResultCount
		}
	case CT_SINGLE: //单牌类型
	case CT_DOUBLE: //对牌类型
	case CT_THREE:
		{ //三条类型
			// 变量定义
			cbReferCard := cbTurnCardData[0]
			var cbSameCount uint8
			cbSameCount = 1
			if cbTurnOutType == CT_DOUBLE {
				cbSameCount = 2
			} else if cbTurnOutType == CT_THREE {
				cbSameCount = 3
			}

			// 搜索相同牌
			cbResultCount = dg.SearchSameCard(cbCardData[:], cbCardCount, cbReferCard, cbSameCount, pSearchCardResult)

			break
		}

	case CT_SINGLE_LINE: //单连类型
	case CT_DOUBLE_LINE: //对连类型
	case CT_THREE_LINE:
		{ //三连类型
			// 变量定义
			var cbBlockCount uint8
			cbBlockCount = 1
			if cbTurnOutType == CT_DOUBLE_LINE {
				cbBlockCount = 2
			} else if cbTurnOutType == CT_THREE_LINE {
				cbBlockCount = 3
			}

			cbLineCount := cbTurnCardCount / cbBlockCount

			// 搜索边牌
			cbResultCount = dg.SearchLineCardType(cbCardData[:], cbCardCount, cbTurnCardData[0], cbBlockCount, cbLineCount, pSearchCardResult)

			break
		}
	case CT_THREE_TAKE_ONE: //三带一单
	case CT_THREE_TAKE_TWO:
		{ //三带一对
			// 效验牌数
			if cbCardCount < cbTurnCardCount {
				break
			}

			// 如果是三带一或三带二
			if cbTurnCardCount == 4 || cbTurnCardCount == 5 {
				var cbTakeCardCount uint8
				if cbTurnOutType == CT_THREE_TAKE_ONE {
					cbTakeCardCount = 1
				} else {
					cbTakeCardCount = 2
				}

				// 搜索三带牌型
				cbResultCount = dg.SearchTakeCardType(cbCardData[:], cbCardCount, cbTurnCardData[2], 3, cbTakeCardCount, pSearchCardResult)
			} else {
				// 变量定义
				var cbBlockCount uint8
				cbBlockCount = 3

				var cbLineCount uint8
				if cbTurnOutType == CT_THREE_TAKE_ONE {
					cbLineCount = 4
				} else {
					cbLineCount = 5
				}

				var cbTakeCardCount uint8
				if cbTurnOutType == CT_THREE_TAKE_ONE {
					cbTakeCardCount = 1
				} else {
					cbTakeCardCount = 2
				}

				// 搜索连牌
				var cbTmpTurnCard [MAX_COUNT]uint8
				copy(cbTmpTurnCard[:], cbTurnCardData[:cbTurnCardCount])
				dg.SortOutCardList(cbTmpTurnCard[:], cbTurnCardCount)
				cbResultCount = dg.SearchLineCardType(cbCardData[:], cbCardCount, cbTmpTurnCard[0], cbBlockCount, cbLineCount, pSearchCardResult)

				//提取带牌
				bAllDistill := true
				for i := 0; uint8(i) < cbResultCount; i++ {
					cbResultIndex := cbResultCount - uint8(i) - 1

					// 变量定义
					var cbTmpCardData [MAX_COUNT]uint8
					cbTmpCardCount := cbCardCount

					//删除连牌
					copy(cbTmpCardData[:], cbCardData[:cbCardCount])

					if dg.RemoveCard(pSearchCardResult.cbResultCard[cbResultIndex][:], pSearchCardResult.cbCardCount[cbResultIndex], cbTmpCardData[:], cbTmpCardCount) {

					}
					//VERIFY( );
					cbTmpCardCount -= pSearchCardResult.cbCardCount[cbResultIndex]

					// 分析牌
					var TmpResult tagAnalyseResult
					dg.AnalysebCardData(cbTmpCardData[:], cbTmpCardCount, &TmpResult)

					// 提取牌
					var cbDistillCard [MAX_COUNT]uint8
					var cbDistillCount uint8
					for j := cbTakeCardCount - 1; j < 4; j++ {
						if TmpResult.cbBlockCount[j] > 0 {
							if j+1 == cbTakeCardCount && TmpResult.cbBlockCount[j] >= cbLineCount {
								cbTmpBlockCount := TmpResult.cbBlockCount[j]
								copy(cbDistillCard[:], TmpResult.cbCardData[j][(cbTmpBlockCount-cbLineCount)*(j+1):(cbTmpBlockCount-cbLineCount)*(j+1)+(j+1)*cbLineCount])
								cbDistillCount = (j + 1) * cbLineCount
								break
							} else {
								for k := 0; uint8(k) < TmpResult.cbBlockCount[j]; k++ {
									cbTmpBlockCount := TmpResult.cbBlockCount[j]
									copy(cbDistillCard[cbDistillCount:], TmpResult.cbCardData[j][(cbTmpBlockCount-uint8(k)-1)*(j+1):(cbTmpBlockCount-uint8(k)-1)*(j+1)+cbTakeCardCount])
									cbDistillCount += cbTakeCardCount
									// 提取完成
									if cbDistillCount == cbTakeCardCount*cbLineCount {
										break
									}
								}
							}
						}

						// 提取完成
						if cbDistillCount == cbTakeCardCount*cbLineCount {
							break
						}
					}

					// 提取完成
					if cbDistillCount == cbTakeCardCount*cbLineCount {
						// 复制带牌
						cbCount := pSearchCardResult.cbCardCount[cbResultIndex]
						copy(pSearchCardResult.cbResultCard[cbResultIndex][cbCount:], cbDistillCard[:cbDistillCount])
						pSearchCardResult.cbCardCount[cbResultIndex] += cbDistillCount
					} else { // 否则删除连牌
						bAllDistill = false
						pSearchCardResult.cbCardCount[cbResultIndex] = 0
					}
				}

				// 整理组合
				if !bAllDistill {
					pSearchCardResult.cbSearchCount = cbResultCount
					cbResultCount = 0
					for i := 0; uint8(i) < pSearchCardResult.cbSearchCount; i++ {
						if pSearchCardResult.cbCardCount[i] != 0 {
							tmpSearchCardResult.cbCardCount[cbResultCount] = pSearchCardResult.cbCardCount[i]
							copy(tmpSearchCardResult.cbResultCard[cbResultCount][:], pSearchCardResult.cbResultCard[i][:pSearchCardResult.cbCardCount[i]])
							cbResultCount++
						}
					}
					tmpSearchCardResult.cbSearchCount = cbResultCount
					// 拷贝结构体
					//CopyMemory( pSearchCardResult,&tmpSearchCardResult,sizeof(tagSearchCardResult) );
				}
			}

			break
		}
	case CT_FOUR_TAKE_ONE: //四带两单
	case CT_FOUR_TAKE_TWO:
		{ //四带两双

			var cbTakeCount uint8
			if cbTurnOutType == CT_FOUR_TAKE_ONE {
				cbTakeCount = 1
			} else {
				cbTakeCount = 2
			}

			var cbTmpTurnCard [MAX_COUNT]uint8
			copy(cbTmpTurnCard[:], cbTurnCardData[:cbTurnCardCount])
			dg.SortOutCardList(cbTmpTurnCard[:], cbTurnCardCount)

			// 搜索带牌
			cbResultCount = dg.SearchTakeCardType(cbCardData[:], cbCardCount, cbTmpTurnCard[0], 4, cbTakeCount, pSearchCardResult)

			break
		}
	}

	// 搜索炸弹
	if (cbCardCount >= 4) && (cbTurnOutType != CT_MISSILE_CARD) {
		// 变量定义
		var cbReferCard uint8
		if cbTurnOutType == CT_BOMB_CARD {
			cbReferCard = cbTurnCardData[0]
		}

		// 搜索炸弹
		cbTmpResultCount := dg.SearchSameCard(cbCardData[:], cbCardCount, cbReferCard, 4, &tmpSearchCardResult)
		for i := 0; uint8(i) < cbTmpResultCount; i++ {
			pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[i]
			copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[i][:tmpSearchCardResult.cbCardCount[i]])
			cbResultCount++

		}
	}

	// 搜索火箭
	if cbTurnOutType != CT_MISSILE_CARD && (cbCardCount >= 2) && (cbCardData[0] == 0x4F) && (cbCardData[1] == 0x4E) {
		// 设置结果
		pSearchCardResult.cbCardCount[cbResultCount] = 2
		pSearchCardResult.cbResultCard[cbResultCount][0] = cbCardData[0]
		pSearchCardResult.cbResultCard[cbResultCount][1] = cbCardData[1]

		cbResultCount++
	}

	pSearchCardResult.cbSearchCount = cbResultCount
	return cbResultCount
}

//同牌搜索
func (dg *ddz_logic) SearchSameCard(cbHandCardData []uint8, cbHandCardCount uint8, cbReferCard uint8, cbSameCardCount uint8, pSearchCardResult *tagSearchCardResult) uint8 {
	// 设置结果
	var cbResultCount uint8

	// 构造扑克
	var cbCardData [MAX_COUNT]uint8
	cbCardCount := cbHandCardCount
	copy(cbCardData[:], cbHandCardData[:cbHandCardCount])

	// 排列扑克
	dg.DDZSortCardList(cbCardData[:], cbCardCount, ST_ORDER)

	//分析扑克
	var AnalyseResult tagAnalyseResult
	dg.AnalysebCardData(cbCardData[:], cbCardCount, &AnalyseResult)

	var cbReferLogicValue uint8
	if cbReferCard == 0 {
		cbReferLogicValue = 0
	} else {
		cbReferLogicValue = dg.GetCardLogicValue(cbReferCard)
	}

	cbBlockIndex := cbSameCardCount - 1
	for cbBlockIndex < 4 {
		for i := 0; uint8(i) < AnalyseResult.cbBlockCount[cbBlockIndex]; i++ {
			cbIndex := (AnalyseResult.cbBlockCount[cbBlockIndex] - uint8(i) - 1) * (cbBlockIndex + 1)
			if dg.GetCardLogicValue(AnalyseResult.cbCardData[cbBlockIndex][cbIndex]) > cbReferLogicValue {
				if pSearchCardResult == nil {
					return 1
				}

				if cbResultCount >= 20 {

				}

				// 复制扑克
				copy(pSearchCardResult.cbResultCard[cbResultCount][:], AnalyseResult.cbCardData[cbBlockIndex][cbIndex:cbIndex+cbSameCardCount])

				pSearchCardResult.cbCardCount[cbResultCount] = cbSameCardCount

				cbResultCount++
			}
		}

		cbBlockIndex++
	}

	if pSearchCardResult != nil {
		pSearchCardResult.cbSearchCount = cbResultCount
	}
	return cbResultCount
}

//带牌类型搜索(三带一，四带一等)
func (dg *ddz_logic) SearchTakeCardType(cbHandCardData []uint8, cbHandCardCount uint8, cbReferCard uint8, cbSameCount uint8, cbTakeCardCount uint8, pSearchCardResult *tagSearchCardResult) uint8 {

	// 设置结果
	var cbResultCount uint8

	// 效验
	if cbSameCount != 3 && cbSameCount != 4 {
		log.Error("cuowu")
		return cbResultCount
	}
	if cbTakeCardCount != 1 && cbTakeCardCount != 2 {
		log.Error("cuowu")
		return cbResultCount
	}

	// 长度判断
	if cbSameCount == 4 && cbHandCardCount < cbSameCount+cbTakeCardCount*2 || cbHandCardCount < cbSameCount+cbTakeCardCount {
		return cbResultCount
	}

	// 构造扑克
	var cbCardData [MAX_COUNT]uint8
	cbCardCount := cbHandCardCount
	copy(cbCardData[:], cbHandCardData[:cbHandCardCount])

	// 排列扑克
	dg.DDZSortCardList(cbCardData[:], cbCardCount, ST_ORDER)

	//搜索同张
	var SameCardResult tagSearchCardResult
	cbSameCardResultCount := dg.SearchSameCard(cbCardData[:], cbCardCount, cbReferCard, cbSameCount, &SameCardResult)

	if cbSameCardResultCount > 0 {
		// 分析扑克
		var AnalyseResult tagAnalyseResult
		dg.AnalysebCardData(cbCardData[:], cbCardCount, &AnalyseResult)

		// 需要牌数
		cbNeedCount := cbSameCount + cbTakeCardCount
		if cbSameCount == 4 {
			cbNeedCount += cbTakeCardCount
		}

		// 提取带牌
		for i := 0; uint8(i) < cbSameCardResultCount; i++ {
			bMerge := false

			for j := cbTakeCardCount - 1; j < 4; j++ {
				for k := 0; uint8(k) < AnalyseResult.cbBlockCount[j]; k++ {
					// 从小到大
					cbIndex := (AnalyseResult.cbBlockCount[j] - uint8(k) - 1) * (j + 1)

					// 过滤相同牌
					if dg.GetCardValue(int(SameCardResult.cbResultCard[i][0])) ==
						dg.GetCardValue(int(AnalyseResult.cbCardData[j][cbIndex])) {
						continue
					}

					// 复制带牌
					cbCount := SameCardResult.cbCardCount[i]
					copy(SameCardResult.cbResultCard[i][cbCount:], AnalyseResult.cbCardData[j][cbIndex:cbIndex+cbTakeCardCount])
					SameCardResult.cbCardCount[i] += cbTakeCardCount

					if SameCardResult.cbCardCount[i] < cbNeedCount {
						continue
					}

					if pSearchCardResult == nil {
						return 1
					}

					// 复制结果
					copy(pSearchCardResult.cbResultCard[cbResultCount][:], SameCardResult.cbResultCard[i][:SameCardResult.cbCardCount[i]])
					pSearchCardResult.cbCardCount[cbResultCount] = SameCardResult.cbCardCount[i]
					cbResultCount++

					bMerge = true

					// 下一组合
					break
				}

				if bMerge {
					break
				}
			}
		}
	}

	if pSearchCardResult != nil {
		pSearchCardResult.cbSearchCount = cbResultCount
	}
	return cbResultCount
}

//连牌搜索
func (dg *ddz_logic) SearchLineCardType(cbHandCardData []uint8, cbHandCardCount uint8, cbReferCard uint8, cbBlockCount uint8, cbLineCount uint8, pSearchCardResult *tagSearchCardResult) uint8 {
	// 设置结果
	var cbResultCount uint8
	var cbLessLineCount uint8

	if cbLineCount == 0 {
		if cbBlockCount == 1 {
			cbLessLineCount = 5
		} else if cbBlockCount == 2 {
			cbLessLineCount = 3
		} else {
			cbLessLineCount = 2
		}
	} else {
		cbLessLineCount = cbLineCount
	}

	var cbReferIndex uint8
	cbReferIndex = 2
	if cbReferCard != 0 {
		if dg.GetCardLogicValue(cbReferCard)-cbLessLineCount >= 2 {

		}
		//ASSERT( GetCardLogicValue(cbReferCard)-cbLessLineCount>=2 );
		cbReferIndex = dg.GetCardLogicValue(cbReferCard) - cbLessLineCount + 1
	}
	// 超过A
	if cbReferIndex+cbLessLineCount > 14 {
		return cbResultCount
	}

	// 长度判断
	if cbHandCardCount < cbLessLineCount*cbBlockCount {
		return cbResultCount
	}

	// 构造扑克
	var cbCardData [MAX_COUNT]uint8
	cbCardCount := cbHandCardCount
	copy(cbCardData[:], cbHandCardData[:cbHandCardCount])

	// 排列扑克
	dg.DDZSortCardList(cbCardData[:], cbCardCount, ST_ORDER)

	// 分析扑克
	var Distributing tagDistributing
	dg.AnalysebDistributing(cbCardData[:], cbCardCount, &Distributing)

	// 搜索顺子
	var cbTmpLinkCount uint8
	var cbValueIndex uint8
	for cbValueIndex = cbReferIndex; cbValueIndex < 13; cbValueIndex++ {
		// 继续判断
		if Distributing.cbDistributing[cbValueIndex][cbIndexCount] < cbBlockCount {
			if cbTmpLinkCount < cbLessLineCount {
				cbTmpLinkCount = 0
				continue
			} else {
				cbValueIndex--
			}
		} else {
			cbTmpLinkCount++
			// 寻找最长连
			if cbLineCount == 0 {
				continue
			}
		}

		if cbTmpLinkCount >= cbLessLineCount {
			if pSearchCardResult == nil {
				return 1
			}

			if cbResultCount > 20 {

			}
			//ASSERT( cbResultCount < CountArray(pSearchCardResult.cbCardCount) );

			// 复制扑克
			var cbCount uint8
			for cbIndex := cbValueIndex + 1 - cbTmpLinkCount; cbIndex <= cbValueIndex; cbIndex++ {
				var cbTmpCount uint8
				for cbColorIndex := 0; cbColorIndex < 4; cbColorIndex++ {
					for cbColorCount := 0; uint8(cbColorCount) < Distributing.cbDistributing[cbIndex][3-cbColorIndex]; cbColorCount++ {
						pSearchCardResult.cbResultCard[cbResultCount][cbCount] = dg.MakeCardData(cbIndex, 3-uint8(cbColorIndex))
						cbCount++
						cbTmpCount++
						if cbTmpCount == cbBlockCount {
							break
						}
					}
					if cbTmpCount == cbBlockCount {
						break
					}
				}
			}

			// 设置变量
			pSearchCardResult.cbCardCount[cbResultCount] = cbCount
			cbResultCount++

			if cbLineCount != 0 {
				cbTmpLinkCount--
			} else {
				cbTmpLinkCount = 0
			}
		}
	}

	// 特殊顺子
	if cbTmpLinkCount >= cbLessLineCount-1 && cbValueIndex == 13 {
		if Distributing.cbDistributing[0][cbIndexCount] >= cbBlockCount || cbTmpLinkCount >= cbLessLineCount {
			if pSearchCardResult == nil {
				return 1
			}
			if cbResultCount > 20 {
				//ASSERT( cbResultCount < CountArray(pSearchCardResult.cbCardCount) );
			}

			// 复制扑克
			var cbCount uint8
			var cbTmpCount uint8
			for cbIndex := cbValueIndex - cbTmpLinkCount; cbIndex < 13; cbIndex++ {
				cbTmpCount = 0
				for cbColorIndex := 0; cbColorIndex < 4; cbColorIndex++ {
					for cbColorCount := 0; uint8(cbColorCount) < Distributing.cbDistributing[cbIndex][3-cbColorIndex]; cbColorCount++ {
						pSearchCardResult.cbResultCard[cbResultCount][cbCount] = dg.MakeCardData(cbIndex, 3-uint8(cbColorIndex))
						cbCount++
						cbTmpCount++
						if cbTmpCount == cbBlockCount {
							break
						}
					}
					if cbTmpCount == cbBlockCount {
						break
					}
				}
			}
			// 复制A
			if Distributing.cbDistributing[0][cbIndexCount] >= cbBlockCount {
				cbTmpCount = 0
				for cbColorIndex := 0; cbColorIndex < 4; cbColorIndex++ {
					for cbColorCount := 0; uint8(cbColorCount) < Distributing.cbDistributing[0][3-cbColorIndex]; cbColorCount++ {
						pSearchCardResult.cbResultCard[cbResultCount][cbCount] = dg.MakeCardData(0, 3-uint8(cbColorIndex))
						cbCount++
						cbTmpCount++
						if cbTmpCount == cbBlockCount {
							break
						}
					}
					if cbTmpCount == cbBlockCount {
						break
					}
				}
			}

			// 设置变量
			pSearchCardResult.cbCardCount[cbResultCount] = cbCount
			cbResultCount++
		}
	}

	if pSearchCardResult != nil {
		pSearchCardResult.cbSearchCount = cbResultCount
	}
	return cbResultCount
}

//搜索飞机
func (dg *ddz_logic) SearchThreeTwoLine(cbHandCardData []uint8, cbHandCardCount uint8, pSearchCardResult *tagSearchCardResult) uint8 {

	// 变量定义
	var tmpSearchResult tagSearchCardResult
	var tmpSingleWing tagSearchCardResult
	var tmpDoubleWing tagSearchCardResult
	var cbTmpResultCount uint8

	// 搜索连牌
	cbTmpResultCount = dg.SearchLineCardType(cbHandCardData, cbHandCardCount, 0, 3, 0, &tmpSearchResult)

	if cbTmpResultCount > 0 {
		//提取带牌
		for i := 0; uint8(i) < cbTmpResultCount; i++ {
			// 变量定义
			var cbTmpCardData [MAX_COUNT]uint8
			var cbTmpCardCount = cbHandCardCount

			// 不够牌
			if cbHandCardCount-tmpSearchResult.cbCardCount[i] < tmpSearchResult.cbCardCount[i]/3 {
				var cbNeedDelCount uint8
				cbNeedDelCount = 3
				for cbHandCardCount+uint8(cbNeedDelCount)-tmpSearchResult.cbCardCount[i] < (tmpSearchResult.cbCardCount[i]-uint8(cbNeedDelCount))/3 {
					cbNeedDelCount += 3
				}

				// 不够连牌
				if (tmpSearchResult.cbCardCount[i]-uint8(cbNeedDelCount))/3 < 2 {
					// 废除连牌
					continue
				}

				//拆分连牌
				dg.RemoveCard(tmpSearchResult.cbResultCard[i][:], cbNeedDelCount, tmpSearchResult.cbResultCard[i][:], tmpSearchResult.cbCardCount[i])
				tmpSearchResult.cbCardCount[i] -= cbNeedDelCount
			}

			if pSearchCardResult == nil {
				return 1
			}

			// 删除连牌
			copy(cbTmpCardData[:], cbHandCardData[:cbHandCardCount])
			dg.RemoveCard(tmpSearchResult.cbResultCard[i][:], tmpSearchResult.cbCardCount[i], cbTmpCardData[:], cbTmpCardCount)
			//VERIFY( RemoveCard( tmpSearchResult.cbResultCard[i],tmpSearchResult.cbCardCount[i], cbTmpCardData,cbTmpCardCount ) );
			cbTmpCardCount -= tmpSearchResult.cbCardCount[i]

			// 组合飞机
			cbNeedCount := tmpSearchResult.cbCardCount[i] / 3
			if cbNeedCount <= cbTmpCardCount {
				return 0
			}
			//ASSERT( cbNeedCount <= cbTmpCardCount );

			cbResultCount := tmpSingleWing.cbSearchCount
			tmpSingleWing.cbSearchCount++
			copy(tmpSingleWing.cbResultCard[cbResultCount][:], tmpSearchResult.cbResultCard[i][:tmpSearchResult.cbCardCount[i]])
			copy(tmpSingleWing.cbResultCard[cbResultCount][tmpSearchResult.cbCardCount[i]:], cbTmpCardData[cbTmpCardCount-cbNeedCount:cbTmpCardCount-cbNeedCount+cbNeedCount])
			tmpSingleWing.cbCardCount[i] = tmpSearchResult.cbCardCount[i] + cbNeedCount

			// 不够带翅膀
			if cbTmpCardCount < tmpSearchResult.cbCardCount[i]/3*2 {
				var cbNeedDelCount uint8
				cbNeedDelCount = 3
				for cbTmpCardCount+cbNeedDelCount-tmpSearchResult.cbCardCount[i] < (tmpSearchResult.cbCardCount[i]-cbNeedDelCount)/3*2 {
					cbNeedDelCount += 3
				}
				// 不够连牌
				if (tmpSearchResult.cbCardCount[i]-cbNeedDelCount)/3 < 2 {
					//废除连牌
					continue
				}
				// 拆分连牌
				dg.RemoveCard(tmpSearchResult.cbResultCard[i][:], cbNeedDelCount, tmpSearchResult.cbResultCard[i][:], tmpSearchResult.cbCardCount[i])
				tmpSearchResult.cbCardCount[i] -= cbNeedDelCount

				// 重新删除连牌
				copy(cbTmpCardData[:], cbHandCardData[:cbHandCardCount])
				if dg.RemoveCard(tmpSearchResult.cbResultCard[i][:], tmpSearchResult.cbCardCount[i], cbTmpCardData[:], cbTmpCardCount) {

				}
				//VERIFY( RemoveCard( tmpSearchResult.cbResultCard[i],tmpSearchResult.cbCardCount[i],
				//cbTmpCardData,cbTmpCardCount ) );
				cbTmpCardCount = cbHandCardCount - tmpSearchResult.cbCardCount[i]
			}

			// 分析牌
			var TmpResult tagAnalyseResult
			dg.AnalysebCardData(cbTmpCardData[:], cbTmpCardCount, &TmpResult)

			//提取翅膀
			var cbDistillCard [MAX_COUNT]uint8
			var cbDistillCount uint8
			cbLineCount := tmpSearchResult.cbCardCount[i] / 3
			for j := 1; j < 4; j++ {
				if TmpResult.cbBlockCount[j] > 0 {
					if j+1 == 2 && TmpResult.cbBlockCount[j] >= cbLineCount {
						cbTmpBlockCount := TmpResult.cbBlockCount[j]
						copy(cbDistillCard[:], TmpResult.cbCardData[j][(cbTmpBlockCount-cbLineCount)*uint8(j+1):(cbTmpBlockCount-cbLineCount)*uint8(j+1)+uint8(j+1)*cbLineCount])
						cbDistillCount = uint8(j+1) * cbLineCount
						break
					} else {
						for k := 0; uint8(k) < TmpResult.cbBlockCount[j]; k++ {
							cbTmpBlockCount := TmpResult.cbBlockCount[j]
							copy(cbDistillCard[cbDistillCount:], TmpResult.cbCardData[j][(cbTmpBlockCount-uint8(k)-1)*uint8(j+1):(cbTmpBlockCount-uint8(k)-1)*uint8(j+1)+2])
							cbDistillCount += 2

							//提取完成
							if cbDistillCount == 2*cbLineCount {
								break
							}
						}
					}
				}
				// 提取完成
				if cbDistillCount == 2*cbLineCount {
					break
				}
			}

			//提取完成
			if cbDistillCount == 2*cbLineCount {
				// 复制翅膀
				cbResultCount = tmpDoubleWing.cbSearchCount
				tmpDoubleWing.cbSearchCount++
				copy(tmpDoubleWing.cbResultCard[cbResultCount][:], tmpSearchResult.cbResultCard[i][:tmpSearchResult.cbCardCount[i]])
				copy(tmpDoubleWing.cbResultCard[cbResultCount][tmpSearchResult.cbCardCount[i]:], cbDistillCard[:cbDistillCount])
				tmpDoubleWing.cbCardCount[i] = tmpSearchResult.cbCardCount[i] + cbDistillCount
			}
		}

		// 复制结果
		for i := 0; uint8(i) < tmpDoubleWing.cbSearchCount; i++ {
			cbResultCount := pSearchCardResult.cbSearchCount
			pSearchCardResult.cbSearchCount++
			copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpDoubleWing.cbResultCard[i][:tmpDoubleWing.cbCardCount[i]])
			pSearchCardResult.cbCardCount[cbResultCount] = tmpDoubleWing.cbCardCount[i]
		}
		for i := 0; uint8(i) < tmpSingleWing.cbSearchCount; i++ {
			cbResultCount := pSearchCardResult.cbSearchCount
			pSearchCardResult.cbSearchCount++
			copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSingleWing.cbResultCard[i][:tmpSingleWing.cbCardCount[i]])
			pSearchCardResult.cbCardCount[cbResultCount] = tmpSingleWing.cbCardCount[i]
		}
	}

	if pSearchCardResult == nil {
		return 0
	} else {
		return pSearchCardResult.cbSearchCount
	}
}

//扑克转换
func (dg *ddz_logic) GetUserCards(cbCardIndex []int) (cbCardData []int) {
	//转换扑克

	return cbCardData
}
