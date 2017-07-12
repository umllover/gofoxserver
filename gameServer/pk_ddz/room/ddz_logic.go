package room

import (
	"encoding/json"
	"mj/gameServer/common/pk/pk_base"

	"mj/gameServer/db/model"

	"mj/common/msg/pk_ddz_msg"

	"github.com/lovelly/leaf/log"
)

func NewDDZLogic(ConfigIdx int, info *model.CreateRoomInfo) *ddz_logic {
	l := new(ddz_logic)
	l.BaseLogic = pk_base.NewBaseLogic(ConfigIdx)

	var setInfo pk_ddz_msg.C2G_DDZ_CreateRoomInfo
	if err := json.Unmarshal([]byte(info.OtherInfo), &setInfo); err == nil {
		l.GameType = setInfo.GameType
	}

	return l
}

type ddz_logic struct {
	*pk_base.BaseLogic
	GameType int
	LizeCard int
}

const (
	// 牌类型
	CT_ERROR          = 0  // 错误类型
	CT_SINGLE         = 1  // 单张牌（散牌）
	CT_DOUBLE         = 2  // 对子牌
	CT_THREE          = 3  // 三张牌
	CT_SINGLE_LINE    = 4  // 单顺子
	CT_DOUBLE_LINE    = 5  // 双顺子
	CT_THREE_LINE     = 6  // 三顺子
	CT_THREE_TAKE_ONE = 7  // 三带一
	CT_THREE_TAKE_TWO = 8  // 三带二
	CT_FOUR_TAKE_ONE  = 9  // 四带两单
	CT_FOUR_TAKE_TWO  = 10 // 四带两对
	CT_BOMB_CARD      = 11 // 炸弹类型
	CT_King_2         = 12 // 双王类型
	CT_King_3         = 13 // 三王类型
	CT_King_4         = 14 // 四王类型
	CT_King_5         = 15 // 五王类型
	CT_King_6         = 16 // 六王类型
	CT_King_7         = 17 // 七王类型
	CT_King_8         = 18 // 八王类型

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
	cbBlockCount [4]int            //扑克数目
	cbCardData   [4][MAX_COUNT]int //扑克数据
}

//出牌结果
type tagOutCardResult struct {
	cbCardCount  int            //扑克数目
	cbResultCard [MAX_COUNT]int //结果扑克
}

//分布信息
type tagDistributing struct {
	cbCardCount    int        //扑克数目
	cbDistributing [15][6]int //分布信息
}

//搜索结果
type tagSearchCardResult struct {
	cbSearchCount int                       //结果数目
	cbCardCount   [MAX_COUNT]int            //扑克数目
	cbResultCard  [MAX_COUNT][MAX_COUNT]int //结果扑克
}

//获取类型
func (dg *ddz_logic) GetCardType(cbCardData []int) int {
	cbCardCount := len(cbCardData)
	dg.SortCardList(cbCardData, cbCardCount)
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
				return CT_King_2
			}
			if dg.GetCardLogicValue(cbCardData[0]) == dg.GetCardLogicValue(cbCardData[1]) {
				return CT_DOUBLE
			}

			return CT_ERROR
		}
	}

	kingCount := 0 // 判断王数量
	for _, v := range cbCardData {
		if v >= 0x4E && v <= 0x4F {
			kingCount++
		} else {
			break
		}
	}

	if kingCount == cbCardCount {
		return CT_King_2 + kingCount - 2
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
			for i := 1; int(i) < AnalyseResult.cbBlockCount[2]; i++ {
				cbCardData := AnalyseResult.cbCardData[2][i*3]
				if cbFirstLogicValue != (dg.GetCardLogicValue(cbCardData) + int(i)) {
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
			if cbFirstLogicValue != (dg.GetCardLogicValue(cbCardData) + int(i)) {
				return CT_ERROR
			}
		}

		// 二连判断
		if (AnalyseResult.cbBlockCount[1] * 2) == cbCardCount {
			return CT_DOUBLE_LINE
		}

		return CT_ERROR
	}

	// 全部无重复，都是单张
	if (AnalyseResult.cbBlockCount[0] >= 5) && (AnalyseResult.cbBlockCount[0] == cbCardCount) {
		// 变量定义
		cbCardData := AnalyseResult.cbCardData[0][0]
		cbFirstLogicValue := dg.GetCardLogicValue(cbCardData)

		//错误过虑
		if cbFirstLogicValue >= 15 {
			return CT_ERROR
		}

		//连牌判断
		for i := 1; i < AnalyseResult.cbBlockCount[0]; i++ {
			cbCardData := AnalyseResult.cbCardData[0][i]
			if cbFirstLogicValue != (dg.GetCardLogicValue(cbCardData) + i) {
				return CT_ERROR
			}
		}

		return CT_SINGLE_LINE
	}

	return CT_ERROR
}

//排列扑克
func (dg *ddz_logic) SortCardList(cbCardData []int, cbCount int) {
	dg.DDZSortCardList(cbCardData, len(cbCardData), 1)
}

func (dg *ddz_logic) DDZSortCardList(arry []int, cbCardCount int, cbSortType int) {
	// 数目过虑
	if cbCardCount == 0 {
		return
	}
	if cbSortType == ST_CUSTOM {
		return
	}

	startValue := [...]int{0, 11, 11, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 0, 0}

	for num := 0; num < len(arry); num++ {
		var a = arry[num] % 0x10
		arry[num] += startValue[a]
	}

	var arrLen = len(arry)
	Inum, Jnum := 0, 0
	for Inum = 0; Inum < arrLen; Inum++ {
		for Jnum = 0; Jnum < arrLen-1; Jnum++ {
			if (arry[Jnum] % 16) < (arry[Jnum+1] % 16) {
				arry[Jnum] = arry[Jnum] + arry[Jnum+1]
				arry[Jnum+1] = arry[Jnum] - arry[Jnum+1]
				arry[Jnum] = arry[Jnum] - arry[Jnum+1]
			}
			if (arry[Jnum] % 16) == (arry[Jnum+1] % 16) {
				if (arry[Jnum] / 16) < (arry[Jnum+1] / 16) {
					arry[Jnum] = arry[Jnum] + arry[Jnum+1]
					arry[Jnum+1] = arry[Jnum] - arry[Jnum+1]
					arry[Jnum] = arry[Jnum] - arry[Jnum+1]
				}
			}
		}
	}

	endValue := [...]int{0, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 11, 11, 0, 0}
	for num := 0; num < len(arry); num++ {
		var a = arry[num] % 0x10
		arry[num] -= endValue[a]
	}
}

//删除扑克
func (dg *ddz_logic) RemoveCardList(cbRemoveCard []int, cbCardData []int) bool {
	cbRemoveCount := len(cbRemoveCard)
	// 检验数据
	if cbRemoveCount > int(len(cbCardData)) {
		log.Error("要删除的扑克数%i大于已有扑克数%i", cbRemoveCount, len(cbCardData))
		return false
	}

	// 备份
	var tmpCardData []int
	copy(tmpCardData, cbCardData)

	var u8DeleteCount int // 记录删除记录

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
func (dg *ddz_logic) RemoveCard(cbRemoveCard []int, cbRemoveCount int, cbCardData []int, cbCardCount int) bool {
	return dg.RemoveCardList(cbRemoveCard, cbCardData)
}

// 排列出牌扑克
func (dg *ddz_logic) SortOutCardList(cbCardData []int, cbCardCount int) {

	// 获取牌型
	cbCardType := dg.GetCardType(cbCardData)

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
				cbCardCount += int(i+1) * AnalyseResult.cbBlockCount[i]
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
				copy(cbCardData[cbCardCount:], AnalyseResult.cbCardData[i][:int(i+1)*AnalyseResult.cbBlockCount[i]])
				cbCardCount += int(i+1) * AnalyseResult.cbBlockCount[i]
			}
		}
	}

	return
}

//逻辑数值
func (dg *ddz_logic) GetCardLogicValue(cbCardData int) int {
	// 扑克属性
	cbCardColor := int(dg.GetCardColor(int(cbCardData)))
	cbCardValue := int(dg.GetCardValue(int(cbCardData)))

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
func (dg *ddz_logic) CompareCard(cbFirstCard []int, cbNextCard []int) bool {
	cbFirstCount := len(cbFirstCard)
	cbNextCount := len(cbNextCard)
	// 获取类型
	cbNextType := dg.GetCardType(cbNextCard)
	cbFirstType := dg.GetCardType(cbFirstCard)

	// 类型判断
	if cbNextType == CT_ERROR {
		return false
	}
	if cbNextType >= CT_King_2 {
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
func (dg *ddz_logic) MakeCardData(cbValueIndex int, cbColorIndex int) int {
	return (cbColorIndex << 4) | (cbValueIndex + 1)
}

//分析扑克
func (dg *ddz_logic) AnalysebCardData(cbCardData []int, cbCardCount int, AnalyseResult *tagAnalyseResult) {
	// 设置结果
	//ZeroMemory(&AnalyseResult,sizeof(AnalyseResult));

	// 扑克分析
	for i := 0; int(i) < cbCardCount; i++ {
		// 变量定义
		cbSameCount := 1
		cbLogicValue := dg.GetCardLogicValue(cbCardData[i])

		// 搜索同牌
		for j := i + 1; int(j) < cbCardCount; j++ {
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
func (dg *ddz_logic) AnalysebDistributing(cbCardData []int, cbCardCount int, Distributing *tagDistributing) {
	// 设置变量
	//	ZeroMemory(&Distributing,sizeof(Distributing));

	//设置变量
	for i := 0; int(i) < cbCardCount; i++ {
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
func (dg *ddz_logic) SearchOutCard(cbHandCardData []int, cbHandCardCount int, cbTurnCardData []int, cbTurnCardCount int, pSearchCardResult *tagSearchCardResult) int {
	// 设置结果
	if pSearchCardResult == nil {
		return 0
	}

	// 变量定义
	var cbResultCount int
	var tmpSearchCardResult tagSearchCardResult

	// 构造扑克
	var cbCardData [MAX_COUNT]int
	cbCardCount := cbHandCardCount
	copy(cbCardData[:], cbHandCardData[:cbHandCardCount])

	// 排列扑克
	dg.DDZSortCardList(cbCardData[:], cbCardCount, ST_ORDER)

	// 获取类型
	cbTurnOutType := dg.GetCardType(cbTurnCardData)

	// 出牌分析
	switch cbTurnOutType {
	case CT_ERROR:
		{ //错误类型
			// 提取各种牌型一组

			// 是否一手出完
			if dg.GetCardType(cbCardData[:]) != CT_ERROR {
				pSearchCardResult.cbCardCount[cbResultCount] = cbCardCount
				copy(pSearchCardResult.cbResultCard[cbResultCount][:], cbCardData[:cbCardCount])
				cbResultCount++
			}

			// 如果最小牌不是单牌，则提取
			var cbSameCount int
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
			var cbTmpCount int
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
			var cbSameCount int
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
			var cbBlockCount int
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
				var cbTakeCardCount int
				if cbTurnOutType == CT_THREE_TAKE_ONE {
					cbTakeCardCount = 1
				} else {
					cbTakeCardCount = 2
				}

				// 搜索三带牌型
				cbResultCount = dg.SearchTakeCardType(cbCardData[:], cbCardCount, cbTurnCardData[2], 3, cbTakeCardCount, pSearchCardResult)
			} else {
				// 变量定义
				var cbBlockCount int
				cbBlockCount = 3

				var cbLineCount int
				if cbTurnOutType == CT_THREE_TAKE_ONE {
					cbLineCount = 4
				} else {
					cbLineCount = 5
				}

				var cbTakeCardCount int
				if cbTurnOutType == CT_THREE_TAKE_ONE {
					cbTakeCardCount = 1
				} else {
					cbTakeCardCount = 2
				}

				// 搜索连牌
				var cbTmpTurnCard [MAX_COUNT]int
				copy(cbTmpTurnCard[:], cbTurnCardData[:cbTurnCardCount])
				dg.SortOutCardList(cbTmpTurnCard[:], cbTurnCardCount)
				cbResultCount = dg.SearchLineCardType(cbCardData[:], cbCardCount, cbTmpTurnCard[0], cbBlockCount, cbLineCount, pSearchCardResult)

				//提取带牌
				bAllDistill := true
				for i := 0; int(i) < cbResultCount; i++ {
					cbResultIndex := cbResultCount - int(i) - 1

					// 变量定义
					var cbTmpCardData [MAX_COUNT]int
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
					var cbDistillCard [MAX_COUNT]int
					var cbDistillCount int
					for j := cbTakeCardCount - 1; j < 4; j++ {
						if TmpResult.cbBlockCount[j] > 0 {
							if j+1 == cbTakeCardCount && TmpResult.cbBlockCount[j] >= cbLineCount {
								cbTmpBlockCount := TmpResult.cbBlockCount[j]
								copy(cbDistillCard[:], TmpResult.cbCardData[j][(cbTmpBlockCount-cbLineCount)*(j+1):(cbTmpBlockCount-cbLineCount)*(j+1)+(j+1)*cbLineCount])
								cbDistillCount = (j + 1) * cbLineCount
								break
							} else {
								for k := 0; int(k) < TmpResult.cbBlockCount[j]; k++ {
									cbTmpBlockCount := TmpResult.cbBlockCount[j]
									copy(cbDistillCard[cbDistillCount:], TmpResult.cbCardData[j][(cbTmpBlockCount-int(k)-1)*(j+1):(cbTmpBlockCount-int(k)-1)*(j+1)+cbTakeCardCount])
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
					for i := 0; int(i) < pSearchCardResult.cbSearchCount; i++ {
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

			var cbTakeCount int
			if cbTurnOutType == CT_FOUR_TAKE_ONE {
				cbTakeCount = 1
			} else {
				cbTakeCount = 2
			}

			var cbTmpTurnCard [MAX_COUNT]int
			copy(cbTmpTurnCard[:], cbTurnCardData[:cbTurnCardCount])
			dg.SortOutCardList(cbTmpTurnCard[:], cbTurnCardCount)

			// 搜索带牌
			cbResultCount = dg.SearchTakeCardType(cbCardData[:], cbCardCount, cbTmpTurnCard[0], 4, cbTakeCount, pSearchCardResult)

			break
		}
	}

	// 搜索炸弹
	if (cbCardCount >= 4) && (cbTurnOutType < CT_King_2) {
		// 变量定义
		var cbReferCard int
		if cbTurnOutType == CT_BOMB_CARD {
			cbReferCard = cbTurnCardData[0]
		}

		// 搜索炸弹
		cbTmpResultCount := dg.SearchSameCard(cbCardData[:], cbCardCount, cbReferCard, 4, &tmpSearchCardResult)
		for i := 0; int(i) < cbTmpResultCount; i++ {
			pSearchCardResult.cbCardCount[cbResultCount] = tmpSearchCardResult.cbCardCount[i]
			copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpSearchCardResult.cbResultCard[i][:tmpSearchCardResult.cbCardCount[i]])
			cbResultCount++

		}
	}

	// 搜索火箭
	if cbTurnOutType < CT_King_2 && (cbCardCount >= 2) && (cbCardData[0] == 0x4F) && (cbCardData[1] == 0x4E) {
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
func (dg *ddz_logic) SearchSameCard(cbHandCardData []int, cbHandCardCount int, cbReferCard int, cbSameCardCount int, pSearchCardResult *tagSearchCardResult) int {
	// 设置结果
	var cbResultCount int

	// 构造扑克
	var cbCardData [MAX_COUNT]int
	cbCardCount := cbHandCardCount
	copy(cbCardData[:], cbHandCardData[:cbHandCardCount])

	// 排列扑克
	dg.DDZSortCardList(cbCardData[:], cbCardCount, ST_ORDER)

	//分析扑克
	var AnalyseResult tagAnalyseResult
	dg.AnalysebCardData(cbCardData[:], cbCardCount, &AnalyseResult)

	var cbReferLogicValue int
	if cbReferCard == 0 {
		cbReferLogicValue = 0
	} else {
		cbReferLogicValue = dg.GetCardLogicValue(cbReferCard)
	}

	cbBlockIndex := cbSameCardCount - 1
	for cbBlockIndex < 4 {
		for i := 0; int(i) < AnalyseResult.cbBlockCount[cbBlockIndex]; i++ {
			cbIndex := (AnalyseResult.cbBlockCount[cbBlockIndex] - int(i) - 1) * (cbBlockIndex + 1)
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
func (dg *ddz_logic) SearchTakeCardType(cbHandCardData []int, cbHandCardCount int, cbReferCard int, cbSameCount int, cbTakeCardCount int, pSearchCardResult *tagSearchCardResult) int {

	// 设置结果
	var cbResultCount int

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
	var cbCardData [MAX_COUNT]int
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
		for i := 0; int(i) < cbSameCardResultCount; i++ {
			bMerge := false

			for j := cbTakeCardCount - 1; j < 4; j++ {
				for k := 0; int(k) < AnalyseResult.cbBlockCount[j]; k++ {
					// 从小到大
					cbIndex := (AnalyseResult.cbBlockCount[j] - int(k) - 1) * (j + 1)

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
func (dg *ddz_logic) SearchLineCardType(cbHandCardData []int, cbHandCardCount int, cbReferCard int, cbBlockCount int, cbLineCount int, pSearchCardResult *tagSearchCardResult) int {
	// 设置结果
	var cbResultCount int
	var cbLessLineCount int

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

	var cbReferIndex int
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
	var cbCardData [MAX_COUNT]int
	cbCardCount := cbHandCardCount
	copy(cbCardData[:], cbHandCardData[:cbHandCardCount])

	// 排列扑克
	dg.DDZSortCardList(cbCardData[:], cbCardCount, ST_ORDER)

	// 分析扑克
	var Distributing tagDistributing
	dg.AnalysebDistributing(cbCardData[:], cbCardCount, &Distributing)

	// 搜索顺子
	var cbTmpLinkCount int
	var cbValueIndex int
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
			var cbCount int
			for cbIndex := cbValueIndex + 1 - cbTmpLinkCount; cbIndex <= cbValueIndex; cbIndex++ {
				var cbTmpCount int
				for cbColorIndex := 0; cbColorIndex < 4; cbColorIndex++ {
					for cbColorCount := 0; int(cbColorCount) < Distributing.cbDistributing[cbIndex][3-cbColorIndex]; cbColorCount++ {
						pSearchCardResult.cbResultCard[cbResultCount][cbCount] = dg.MakeCardData(cbIndex, 3-int(cbColorIndex))
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
			var cbCount int
			var cbTmpCount int
			for cbIndex := cbValueIndex - cbTmpLinkCount; cbIndex < 13; cbIndex++ {
				cbTmpCount = 0
				for cbColorIndex := 0; cbColorIndex < 4; cbColorIndex++ {
					for cbColorCount := 0; int(cbColorCount) < Distributing.cbDistributing[cbIndex][3-cbColorIndex]; cbColorCount++ {
						pSearchCardResult.cbResultCard[cbResultCount][cbCount] = dg.MakeCardData(cbIndex, 3-int(cbColorIndex))
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
					for cbColorCount := 0; int(cbColorCount) < Distributing.cbDistributing[0][3-cbColorIndex]; cbColorCount++ {
						pSearchCardResult.cbResultCard[cbResultCount][cbCount] = dg.MakeCardData(0, 3-int(cbColorIndex))
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
func (dg *ddz_logic) SearchThreeTwoLine(cbHandCardData []int, cbHandCardCount int, pSearchCardResult *tagSearchCardResult) int {

	// 变量定义
	var tmpSearchResult tagSearchCardResult
	var tmpSingleWing tagSearchCardResult
	var tmpDoubleWing tagSearchCardResult
	var cbTmpResultCount int

	// 搜索连牌
	cbTmpResultCount = dg.SearchLineCardType(cbHandCardData, cbHandCardCount, 0, 3, 0, &tmpSearchResult)

	if cbTmpResultCount > 0 {
		//提取带牌
		for i := 0; int(i) < cbTmpResultCount; i++ {
			// 变量定义
			var cbTmpCardData [MAX_COUNT]int
			var cbTmpCardCount = cbHandCardCount

			// 不够牌
			if cbHandCardCount-tmpSearchResult.cbCardCount[i] < tmpSearchResult.cbCardCount[i]/3 {
				var cbNeedDelCount int
				cbNeedDelCount = 3
				for cbHandCardCount+int(cbNeedDelCount)-tmpSearchResult.cbCardCount[i] < (tmpSearchResult.cbCardCount[i]-int(cbNeedDelCount))/3 {
					cbNeedDelCount += 3
				}

				// 不够连牌
				if (tmpSearchResult.cbCardCount[i]-int(cbNeedDelCount))/3 < 2 {
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
				var cbNeedDelCount int
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
			var cbDistillCard [MAX_COUNT]int
			var cbDistillCount int
			cbLineCount := tmpSearchResult.cbCardCount[i] / 3
			for j := 1; j < 4; j++ {
				if TmpResult.cbBlockCount[j] > 0 {
					if j+1 == 2 && TmpResult.cbBlockCount[j] >= cbLineCount {
						cbTmpBlockCount := TmpResult.cbBlockCount[j]
						copy(cbDistillCard[:], TmpResult.cbCardData[j][(cbTmpBlockCount-cbLineCount)*int(j+1):(cbTmpBlockCount-cbLineCount)*int(j+1)+int(j+1)*cbLineCount])
						cbDistillCount = int(j+1) * cbLineCount
						break
					} else {
						for k := 0; int(k) < TmpResult.cbBlockCount[j]; k++ {
							cbTmpBlockCount := TmpResult.cbBlockCount[j]
							copy(cbDistillCard[cbDistillCount:], TmpResult.cbCardData[j][(cbTmpBlockCount-int(k)-1)*int(j+1):(cbTmpBlockCount-int(k)-1)*int(j+1)+2])
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
		for i := 0; int(i) < tmpDoubleWing.cbSearchCount; i++ {
			cbResultCount := pSearchCardResult.cbSearchCount
			pSearchCardResult.cbSearchCount++
			copy(pSearchCardResult.cbResultCard[cbResultCount][:], tmpDoubleWing.cbResultCard[i][:tmpDoubleWing.cbCardCount[i]])
			pSearchCardResult.cbCardCount[cbResultCount] = tmpDoubleWing.cbCardCount[i]
		}
		for i := 0; int(i) < tmpSingleWing.cbSearchCount; i++ {
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

// 设置癞子牌
func (dg *ddz_logic) SetParamToLogic(args interface{}) {
	dg.LizeCard = args.(int)
}
