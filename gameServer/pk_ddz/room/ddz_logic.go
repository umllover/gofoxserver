package room

import (
	"encoding/json"
	"mj/common/msg/pk_ddz_msg"
	"mj/gameServer/common/pk/pk_base"
	"mj/gameServer/db/model"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
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

	// 数目定义
	MAX_COUNT = 22 //最大数目
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

// 判断是否是火箭
func (dg *ddz_logic) isRocketType(cardArr []int) (int, bool) {
	if len(cardArr) <= 1 {
		log.Debug("火箭至少得两张牌以上")
		return CT_ERROR, false
	}
	var smallKing int // 小王数量
	var largeKing int // 大王数量
	for _, v := range cardArr {
		if v == 0x4e {
			smallKing++
		} else if v == 0x4f {
			largeKing++
		} else {
			break
		}
	}

	if (smallKing + largeKing) == len(cardArr) {
		var kingType int
		kingType |= smallKing
		kingType |= largeKing << 4
		return CT_KING | kingType, true
	}
	return CT_ERROR, false
}

// 判断是否是炸弹
func (dg *ddz_logic) isBombType(cardArr []int) (int, bool) {
	// 不是4张牌
	if len(cardArr) != 4 {
		log.Debug("炸弹得是4张牌")
		return CT_ERROR, false
	}

	// 有王肯定不是炸弹
	if dg.getKingCount(cardArr) > 0 {
		log.Debug("炸弹不能有王")
		return CT_ERROR, false
	}

	var AnalyseResult tagAnalyseResult
	dg.AnalysebCardData(cardArr, len(cardArr), &AnalyseResult)
	// 四张一样的，肯定是炸弹
	if AnalyseResult.cbBlockCount[3] > 0 {
		return CT_BOMB_CARD | (dg.GetCardLogicValue(AnalyseResult.cbCardData[3][0]) << 4), true
	}
	// 检查癞子是否能匹配
	nCount := dg.getLaiziCount(cardArr)
	if nCount > 0 {
		// 有三张，肯定满足
		if AnalyseResult.cbBlockCount[2] > 0 {
			return CT_BOMB_CARD | (dg.getMaxLogicCardValueWithoutLaizi(cardArr) << 4) | nCount, true
		}
		// 两个对子
		if AnalyseResult.cbBlockCount[1] == 2 {
			return CT_BOMB_CARD | (dg.getMaxLogicCardValueWithoutLaizi(cardArr) << 4) | nCount, true
		}
		// 其它貌似不满足了
	}
	return CT_ERROR, false
}

// 判断是否是四带二
func (dg *ddz_logic) isFourTakeTwo(cardArr []int) (int, bool) {
	nCount := len(cardArr)
	if nCount != 6 && nCount != 8 {
		log.Debug("四带二要么是6，要么是8张")
		return CT_ERROR, false
	}

	// 把王去掉
	tmpCard := util.CopySlicInt(cardArr)
	tmpCard, nkingCount := dg.removeKingFromCard(tmpCard)
	// 把癞子去掉
	tmpCard, nLaiziCount := dg.removeValuesFromCard(cardArr, dg.LizeCard)

	// 有王就不可能形成对子
	if nkingCount > 0 && nCount == 8 {
		log.Debug("如果是8张牌，只能是四带两对，那么有王就不可能成为对子")
		return CT_ERROR, false
	}

	// 分析扑克
	var AnalyseResult tagAnalyseResult
	dg.AnalysebCardData(tmpCard, len(tmpCard), &AnalyseResult)
	// 有4根癞子
	if nLaiziCount == 4 {
		// 六根就是4个癞子加两根其它
		if nCount == 6 {
			return CT_FOUR_TAKE_TWO | dg.GetCardLogicValue(dg.LizeCard), true
		}
		// 4个癞子加两个对子
		if AnalyseResult.cbBlockCount[1] == 2 {
			return CT_FOUR_TAKE_TWO | dg.GetCardLogicValue(dg.LizeCard) | 1, true
		}
		// 4个癞子加4个其它
		if AnalyseResult.cbBlockCount[3] > 0 {
			var maxValue int
			if dg.GetCardValue(AnalyseResult.cbCardData[3][0]) != dg.LizeCard {
				maxValue = dg.GetCardLogicValue(AnalyseResult.cbCardData[3][0])
			} else {
				maxValue = dg.GetCardLogicValue(AnalyseResult.cbCardData[3][4])
			}

			return CT_FOUR_TAKE_TWO | (maxValue << 4) | 1, true
		}

		return CT_ERROR, false
	}

	// 有一对非癞子炸
	if AnalyseResult.cbBlockCount[3] > 0 {
		// 6根为4带两根单
		if nCount == 6 {
			return CT_FOUR_TAKE_TWO | (dg.GetCardLogicValue(AnalyseResult.cbCardData[3][0]) << 4) | 0, true
		}
		// 8根牌
		if nCount == 8 {
			// 两个炸弹也满足
			if AnalyseResult.cbBlockCount[3] == 2 {
				maxValue := dg.maxValue(dg.GetCardLogicValue(AnalyseResult.cbCardData[3][0]), dg.GetCardLogicValue(AnalyseResult.cbCardData[3][4]))
				return CT_FOUR_TAKE_TWO | (maxValue << 4) | 1, true

			}
			// 有一个三张+一张癞子
			if AnalyseResult.cbBlockCount[2] == 1 && nLaiziCount == 1 {
				return CT_FOUR_TAKE_TWO | (dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0]) << 4) | 1, true
			}
			// 四带两对
			if AnalyseResult.cbBlockCount[1] == 2 {
				return CT_FOUR_TAKE_TWO | (dg.GetCardLogicValue(AnalyseResult.cbCardData[3][0]) << 4) | 1, true
			}
			// 四带一对+至少一根癞子
			if AnalyseResult.cbBlockCount[1] == 1 {
				if nLaiziCount >= 1 {
					return CT_FOUR_TAKE_TWO | (dg.GetCardLogicValue(AnalyseResult.cbCardData[3][0]) << 4) | 1, true
				}
				return CT_ERROR, false
			}
			// 一张散牌+三张癞子或者两张散牌+两张癞子
			if (AnalyseResult.cbBlockCount[0] == 1 && nLaiziCount == 3) || (AnalyseResult.cbBlockCount[0] == 2 && nLaiziCount == 2) {
				return CT_FOUR_TAKE_TWO | (dg.GetCardLogicValue(AnalyseResult.cbCardData[3][0]) << 4) | 1, true
			}
			// 其它情况不符合
		}
	} else if AnalyseResult.cbBlockCount[2] == 2 {
		// 两对三张+两根癞子
		if nLaiziCount == 2 {
			maxValue := dg.maxValue(dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0]), dg.GetCardLogicValue(AnalyseResult.cbCardData[2][3]))

			return CT_FOUR_TAKE_TWO | (maxValue << 4) | 1, true
		}
		// 其它不满足
	} else if AnalyseResult.cbBlockCount[2] == 1 {
		// 只有一对三张
		// 6张
		if nCount == 6 {
			// 只要有一张癞子就符合条件
			if nLaiziCount > 0 {
				return CT_FOUR_TAKE_TWO | (dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0])) | 0, true
			}
			return CT_ERROR, false
		}

		// 有三张癞子，随便组合都符合
		if nLaiziCount == 3 {
			return CT_FOUR_TAKE_TWO | (dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0])) | 1, true
		}
		// 三张+两对+一癞子
		if nCount == 8 && nLaiziCount > 0 && AnalyseResult.cbBlockCount[1] == 2 {
			return CT_FOUR_TAKE_TWO | (dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0])) | 1, true
		}
		// 三张+一对+三张癞子/三张+一对+两根癞子+一张非王
		if AnalyseResult.cbBlockCount[1] == 1 && nLaiziCount >= 2 {
			return CT_FOUR_TAKE_TWO | (dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0])) | 1, true
		}
		return CT_ERROR, false
	} else if nLaiziCount >= 2 && AnalyseResult.cbBlockCount[0] <= 1 {
		// 有两张以上的癞子，散牌数量小于等于1张
		maxValue := dg.getMaxLogicCardValueWithoutLaizi(tmpCard)
		return CT_FOUR_TAKE_TWO | (maxValue << 4) | 1, true
	}

	return CT_ERROR, false
}

// 是否是飞机带翅膀 9
func (dg *ddz_logic) isThreeLineTake(cardArr []int) (int, bool) {
	if len(cardArr) < 8 {
		log.Debug("飞机至少得有8张牌，当前只有%d", len(cardArr))
		return CT_ERROR, false
	}

	tmpArr, nLaiziCount := dg.removeValuesFromCard(cardArr, dg.LizeCard)

	nMax, b := dg.recursionIsPlane(tmpArr, nLaiziCount)
	if b {
		return CT_THREE_LINE_TAKE | nMax, true
	}

	return CT_ERROR, false
}

// 讲癞子牌递归插入到牌中并判断是否是飞机
func (dg *ddz_logic) recursionIsPlane(cardArr []int, nLaiziCount int) (int, bool) {
	if nLaiziCount == 0 {
		return dg.isPlane(cardArr)
	}

	nLaiziCount--
	for i := 14; i > 2; i-- {
		tmpArr := util.CopySlicInt(cardArr)
		tmpArr = append(tmpArr, i)
		nMax, b := dg.recursionIsPlane(tmpArr, nLaiziCount)
		if b {
			return nMax, true
		}
	}

	return 0, false
}

// 判断是否是飞机
func (dg *ddz_logic) isPlane(tmpArr []int) (int, bool) {

	cardArr := util.CopySlicInt(tmpArr)
	cardArr, nKingCount := dg.removeKingFromCard(cardArr)

	dg.SortCardList(cardArr, len(cardArr))

	var AnalyseResult tagAnalyseResult
	dg.AnalysebCardData(cardArr, len(cardArr), &AnalyseResult)
	if AnalyseResult.cbBlockCount[2] > 1 {

		var maxValue int
		for i := 0; i < AnalyseResult.cbBlockCount[2]-1; i++ {
			if dg.GetCardLogicValue(AnalyseResult.cbCardData[2][i*3])-dg.GetCardLogicValue(AnalyseResult.cbCardData[2][(i+1)*3]) != 1 {
				return CT_ERROR, false
			}
			maxValue = dg.maxValue(maxValue, dg.GetCardLogicValue(AnalyseResult.cbCardData[2][i*3]))
		}
		// 带单根的
		if len(tmpArr)-3*AnalyseResult.cbBlockCount[2] == AnalyseResult.cbBlockCount[2] {
			return CT_THREE_LINE_TAKE | (maxValue << 4) | (AnalyseResult.cbBlockCount[2] << 1), true
		}
		// 带对子的
		if len(cardArr)-3*AnalyseResult.cbBlockCount[2] == AnalyseResult.cbBlockCount[2]*2 && nKingCount == 0 {
			if AnalyseResult.cbBlockCount[1] == AnalyseResult.cbBlockCount[2] {
				return CT_THREE_LINE_TAKE | (maxValue << 4) | (AnalyseResult.cbBlockCount[2] << 1) | 1, true
			}
		}
	}
	return 0, false
}

// 是否是三顺子
func (dg *ddz_logic) isThreeLine(cardArr []int) (int, bool) {
	if len(cardArr) < 6 || dg.getKingCount(cardArr) > 0 {
		log.Debug("三顺子至少得是6张牌，并且不能有王")
		return CT_ERROR, false
	}

	// 有2也不满足
	nCount := dg.getCountWithCardValue(cardArr, 2)
	if nCount > 0 && dg.LizeCard != 2 {
		log.Debug("癞子不为2的情况下，三顺子不能有2")
		return CT_ERROR, false
	}

	tmpArr, nLaiziCount := dg.removeValuesFromCard(cardArr, dg.LizeCard)

	nMax, b := dg.recursionIsLine(tmpArr, nLaiziCount, 3)
	if b {
		return CT_THREE_LINE | nMax, b
	}

	return CT_ERROR, false
}

// 是否双顺子
func (dg *ddz_logic) isDoubleLine(cardArr []int) (int, bool) {
	if len(cardArr) < 6 || dg.getKingCount(cardArr) > 0 {
		log.Debug("双顺子至少得有六张，并且不能有王")
		return CT_ERROR, false
	}

	// 有2也不满足
	nCount := dg.getCountWithCardValue(cardArr, 2)
	if nCount > 0 && dg.LizeCard != 2 {
		log.Debug("癞子不为2的情况下，双顺子不能有2")
		return CT_ERROR, false
	}

	tmpArr, nLaiziCount := dg.removeValuesFromCard(cardArr, dg.LizeCard)

	nMax, b := dg.recursionIsLine(tmpArr, nLaiziCount, 2)
	if b {
		return CT_DOUBLE_LINE | nMax, b
	}

	return CT_ERROR, false
}

// 是否单顺子
func (dg *ddz_logic) isSingleLine(cardArr []int) (int, bool) {
	log.Debug("进入判断是否是单顺子")
	if len(cardArr) < 5 || dg.getKingCount(cardArr) > 0 {
		log.Debug("单顺子至少得5张，并且不能有王")
		return CT_ERROR, false
	}

	// 有2也不满足
	nCount := dg.getCountWithCardValue(cardArr, 2)
	if nCount > 0 && dg.LizeCard != 2 {
		log.Debug("单顺子也不能有2")
		return CT_ERROR, false
	}

	tmpArr, nLaiziCount := dg.removeValuesFromCard(cardArr, dg.LizeCard)

	nMax, b := dg.recursionIsLine(tmpArr, nLaiziCount, 1)
	if b {
		return CT_SINGLE_LINE | nMax, b
	}

	return CT_ERROR, false
}

// 递归判断是否是顺子
func (dg *ddz_logic) recursionIsLine(cardArr []int, nLaiziCount int, nType int) (int, bool) {

	if nLaiziCount == 0 {
		return dg.isLine(cardArr, nType)
	}
	nLaiziCount--
	for i := 14; i > 2; i-- {
		tmpArr := util.CopySlicInt(cardArr)
		tmpArr = append(tmpArr, i)
		nMax, b := dg.recursionIsLine(tmpArr, nLaiziCount, nType)
		if b {
			return nMax, true
		}
	}
	return 0, false
}

// 是否是顺子，nType为单顺子、双顺子、三顺子
func (dg *ddz_logic) isLine(cardArr []int, nType int) (int, bool) {

	lineArr := make([]int, 15)

	var nMax int
	nMin := 14

	// 把每张牌值的数量加到一个表里
	for _, v := range cardArr {
		v1 := dg.GetCardLogicValue(v)
		if v1 > 14 {
			return 0, false
		}
		if v1 > nMax {
			nMax = v1
		}
		if v1 < nMin {
			nMin = v1
		}

		lineArr[v1]++
	}

	tmpArr := lineArr[nMin : nMax+1]
	if len(tmpArr) < 2 {
		return 0, false
	}

	for i := 0; i < len(tmpArr); i++ {
		if tmpArr[i] != nType {
			return 0, false
		}
	}

	nCount := len(tmpArr)
	if nType == 1 {
		return nMax, nCount >= 5
	} else if nType == 2 {
		return nMax, nCount >= 3
	} else if nCount == 3 {
		return nMax, nCount >= 2
	}
	return nMax, true
}

// 是否三带二
func (dg *ddz_logic) isThreeTakeTwo(cardArr []int) (int, bool) {
	// 不是5张牌或者有一张王
	if len(cardArr) != 5 || dg.getKingCount(cardArr) > 0 {
		log.Debug("三带二至少得有5张，且不能有王")
		return CT_ERROR, false
	}

	tmpArr, nLaiziCount := dg.removeValuesFromCard(cardArr, dg.LizeCard)

	// 四张癞子
	if nLaiziCount == 4 {
		maxValue := dg.GetCardLogicValue(tmpArr[0])
		maxValue = dg.maxValue(maxValue, dg.GetCardLogicValue(dg.LizeCard))
		return CT_THREE_TAKE_TWO | maxValue, true
	}

	var AnalyseResult tagAnalyseResult
	dg.AnalysebCardData(tmpArr, len(tmpArr), &AnalyseResult)
	// 三张癞子
	if nLaiziCount == 3 {
		var maxValue int

		if AnalyseResult.cbBlockCount[1] == 1 {
			// 三癞子+一个对子
			maxValue = dg.maxValue(dg.GetCardLogicValue(dg.LizeCard), dg.GetCardLogicValue(AnalyseResult.cbCardData[1][0]))
		} else {
			// 三癞子+两张散牌
			maxValue = dg.maxValue(dg.GetCardLogicValue(AnalyseResult.cbCardData[0][1]), dg.GetCardLogicValue(AnalyseResult.cbCardData[0][0]))
			maxValue = dg.maxValue(maxValue, dg.GetCardLogicValue(dg.LizeCard))
		}

		return CT_THREE_TAKE_TWO | maxValue, true
	}
	// 两张癞子+一个对子
	if nLaiziCount == 2 && AnalyseResult.cbBlockCount[1] == 1 {
		maxValue := dg.maxValue(dg.GetCardLogicValue(AnalyseResult.cbCardData[1][0]), dg.GetCardLogicValue(AnalyseResult.cbCardData[0][0]))
		return CT_THREE_TAKE_TWO | maxValue, true
	}
	// 一张癞子
	if nLaiziCount == 1 {
		// 一对三张
		if AnalyseResult.cbBlockCount[2] == 1 {
			return CT_THREE_TAKE_TWO | dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0]), true
		}
		// 两个对子
		if AnalyseResult.cbBlockCount[1] == 2 {
			maxValue := dg.maxValue(dg.GetCardLogicValue(AnalyseResult.cbCardData[1][0]), dg.GetCardLogicValue(AnalyseResult.cbCardData[1][2]))
			return CT_THREE_TAKE_TWO | maxValue, true
		}
	}

	log.Debug("分析完后%v", AnalyseResult)
	// 无癞子
	if AnalyseResult.cbBlockCount[2] == 1 && AnalyseResult.cbBlockCount[1] == 1 {
		return CT_THREE_TAKE_TWO | dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0]), true
	}

	return CT_ERROR, false
}

// 是否三带一
func (dg *ddz_logic) isThreeTakeOne(cardArr []int) (int, bool) {
	if len(cardArr) == 4 {
		tmpArr, nKingCount := dg.removeKingFromCard(cardArr)
		// 超过一张王就不可能符合
		if nKingCount > 1 {
			log.Debug("三带一不能超过一张王")
			return CT_ERROR, false
		}

		tmpArr, nLaiziCount := dg.removeValuesFromCard(tmpArr, dg.LizeCard)
		// 三张或四张癞子都符合
		if nLaiziCount >= 3 {
			if nLaiziCount == 3 && nKingCount == 0 {
				return CT_THREE_TAKE_ONE | dg.maxValue(dg.GetCardLogicValue(dg.LizeCard), dg.GetCardLogicValue(tmpArr[0])), true
			}
			return CT_THREE_TAKE_ONE | dg.GetCardLogicValue(dg.LizeCard), true
		}

		// 两张癞子
		if nLaiziCount == 2 {
			return CT_THREE_TAKE_ONE | dg.getMaxLogicCardValueWithoutLaizi(tmpArr), true
		}

		var AnalyseResult tagAnalyseResult
		dg.AnalysebCardData(tmpArr, len(tmpArr), &AnalyseResult)
		// 一张癞子+一个对子
		if nLaiziCount == 1 && AnalyseResult.cbBlockCount[1] == 1 {
			return CT_THREE_TAKE_ONE | dg.GetCardLogicValue(AnalyseResult.cbCardData[1][0]), true
		}
		// 三张非癞子
		if AnalyseResult.cbBlockCount[2] == 1 {
			return CT_THREE_TAKE_ONE | dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0]), true
		}
	}
	return CT_ERROR, false
}

// 是否三张牌
func (dg *ddz_logic) isThree(cardArr []int) (int, bool) {
	if len(cardArr) == 3 && dg.getKingCount(cardArr) == 0 {
		tmpArr, nCount := dg.removeValuesFromCard(cardArr, dg.LizeCard)
		var AnalyseResult tagAnalyseResult
		dg.AnalysebCardData(tmpArr, len(tmpArr), &AnalyseResult)
		// 三张一样的
		if AnalyseResult.cbBlockCount[2] > 0 {

			return CT_THREE | dg.GetCardLogicValue(AnalyseResult.cbCardData[2][0]), true
		}

		// 一张癞子+一个对子，对子不为王
		if nCount == 1 && AnalyseResult.cbBlockCount[1] > 0 {
			return CT_THREE | dg.GetCardLogicValue(AnalyseResult.cbCardData[1][0]), true
		}
		// 两张癞子+一张非癞子（不为王）
		if nCount == 2 {
			return CT_THREE | dg.GetCardLogicValue(AnalyseResult.cbCardData[0][0]), true
		}

		// 三张癞子
		if nCount == 3 {
			return CT_THREE | dg.GetCardLogicValue(dg.LizeCard), true
		}
	}
	return CT_ERROR, false
}

// 是否对子牌
func (dg *ddz_logic) isDouble(cardArr []int) (int, bool) {
	if len(cardArr) == 2 && dg.getKingCount(cardArr) == 0 {
		if dg.GetCardValue(cardArr[0]) == dg.GetCardValue(cardArr[1]) {
			return CT_DOUBLE | dg.GetCardLogicValue(cardArr[0]), true
		} else if dg.getLaiziCount(cardArr) > 0 {
			return CT_DOUBLE | dg.getMaxLogicCardValueWithoutLaizi(cardArr), true
		}
	}
	return CT_ERROR, false
}

// 癞子牌的数量
func (dg *ddz_logic) getLaiziCount(cardArr []int) int {
	var nCount int
	for _, v := range cardArr {
		if dg.GetCardValue(v) == dg.LizeCard {
			nCount++
		}
	}
	return nCount
}

// 王的数量
func (dg *ddz_logic) getKingCount(cardArr []int) int {
	var nCount int
	for _, v := range cardArr {
		if v >= 0x4e {
			nCount++
		}
	}
	return nCount
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
		return CT_SINGLE | dg.GetCardLogicValue(cbCardData[0])
	case 2: //对牌
		{

			if nType, isType := dg.isDouble(cbCardData); isType {
				return nType
			}
		}
	}

	log.Debug("是否火箭")
	// 判断是否是火箭12
	nType, isType := dg.isRocketType(cbCardData)
	if isType {
		return nType
	}
	log.Debug("是否炸弹")
	// 两张牌的已经判断完毕
	if cbCardCount == 2 {
		return CT_ERROR
	}

	// 判断是否是炸弹11
	if nType, isType = dg.isBombType(cbCardData); isType {
		return nType
	}
	log.Debug("是否4带2")
	// 判断是否是4带2 10
	if nType, isType = dg.isFourTakeTwo(cbCardData); isType {
		return nType
	}
	log.Debug("是否飞机")
	// 飞机 9
	if nType, isType = dg.isThreeLineTake(cbCardData); isType {
		return nType
	}
	log.Debug("是否三顺子")
	// 判断是否是三顺子 8
	if nType, isType = dg.isThreeLine(cbCardData); isType {
		return nType
	}
	log.Debug("是否双顺子")
	// 判断是否是双顺子 7
	if nType, isType = dg.isDoubleLine(cbCardData); isType {
		return nType
	}
	log.Debug("是否单顺子")
	// 判断是否是单顺子 6
	if nType, isType = dg.isSingleLine(cbCardData); isType {
		return nType
	}
	log.Debug("是否三带二")
	// 三带二 5
	if nType, isType = dg.isThreeTakeTwo(cbCardData); isType {
		return nType
	}
	log.Debug("是否三带一")
	// 判断是否是三带一 4
	if nType, isType = dg.isThreeTakeOne(cbCardData); isType {
		return nType
	}
	log.Debug("是否三牌")
	// 判断是否是三牌 3
	if nType, isType = dg.isThree(cbCardData); isType {
		return nType
	}
	log.Debug("无效")
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
func (dg *ddz_logic) RemoveCardList(cbRemoveCard []int, cbCardData []int) ([]int, bool) {
	cbRemoveCount := len(cbRemoveCard)
	// 检验数据
	if cbRemoveCount > int(len(cbCardData)) {
		log.Error("要删除的扑克数%d大于已有扑克数%d", cbRemoveCount, len(cbCardData))
		return cbCardData, false
	}

	// 备份
	var tmpCardData []int
	copy(tmpCardData, cbCardData)

	cardArr := util.CopySlicInt(cbCardData)

	var u8DeleteCount int // 记录删除记录

	for _, v1 := range cbRemoveCard {
		for j, v2 := range cardArr {
			if v1 == v2 {
				copy(cardArr[j:], cardArr[j+1:])
				cardArr = cardArr[:len(cardArr)-1]
				u8DeleteCount++
				break
			}
		}
	}

	if u8DeleteCount != cbRemoveCount {
		// 删除数量不一，恢复数据
		log.Error("实际删除数量%d与需要删除数量%d不一样", u8DeleteCount, cbRemoveCount)
		copy(cardArr, tmpCardData)
		return cardArr, false
	}

	return cardArr, true
}

//删除扑克
func (dg *ddz_logic) RemoveCard(cbRemoveCard []int, cbRemoveCount int, cbCardData []int, cbCardCount int) bool {
	_, err := dg.RemoveCardList(cbRemoveCard, cbCardData)
	return err
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
	} else if cbCardType == CT_FOUR_TAKE_TWO {
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
		log.Error("求取逻辑数值的扑克数据有误%d", cbCardValue)
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

//构造扑克
func (dg *ddz_logic) MakeCardData(cbValueIndex int, cbColorIndex int) int {
	return (cbColorIndex << 4) | (cbValueIndex + 1)
}

//分析扑克
func (dg *ddz_logic) AnalysebCardData(cbCardData []int, cbCardCount int, AnalyseResult *tagAnalyseResult) {

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

// 设置癞子牌
func (dg *ddz_logic) SetParamToLogic(args interface{}) {
	dg.LizeCard = args.(int)
}

// 比牌
func (dg *ddz_logic) CompareCardWithParam(firstCardData []int, lastCardData []int, args []interface{}) (int, bool) {
	firstCount := len(firstCardData)
	nextCount := len(lastCardData)

	firstType := args[0].(int)

	nextType := dg.GetCardType(lastCardData)

	if firstType == CT_ERROR && nextType != CT_ERROR {
		return nextType, true
	}

	var nType int
	var isType bool

	log.Debug("比火箭")
	nType, isType = dg.isRocketType(lastCardData)
	// 前牌是火箭
	if firstType > CT_KING {
		if isType {
			if nextCount > firstCount {
				return nType, true
			}
			if nextCount == firstCount {
				return nType, nType > firstType
			}
		}
		return CT_ERROR, false
	}

	// 前牌非火箭，后牌是火箭火箭
	if isType {
		return nType, true
	}
	log.Debug("比炸弹")
	nType, isType = dg.isBombType(lastCardData)
	// 前牌是炸弹
	if firstType >= CT_BOMB_CARD && firstType < CT_KING {
		if isType {
			// 两者都有癞子，比逻辑牌
			if (firstType&0xF) > 0 && (nType&0xF) > 0 || (firstType&0xF) == 0 && (nType&0xF) == 0 {
				return nType, (nType & 0xF0) > (firstType & 0xF0)
			}
			// 有一个无癞子，则无癞子的大
			return nType, (firstType & 0xF) > 0
		}
		return CT_ERROR, false
	}

	// 前牌非炸弹，后牌是炸弹
	if isType {
		return nType, true
	}

	// 都非炸弹，进行牌型比较
	// 张数不同
	if firstCount != nextCount {
		return CT_ERROR, false
	}
	log.Debug("比4带2")
	// 四带二
	if firstType >= CT_FOUR_TAKE_TWO {
		if nType, isType = dg.isFourTakeTwo(lastCardData); isType {
			// 带单还是对子要一样
			if (firstType & 0xF) == (nType & 0xF) {
				return nType, (nType & 0xF0) > (firstType & 0xF0)
			}
		}

		return CT_ERROR, false
	}
	log.Debug("比飞机")
	// 飞机带翅膀
	if firstType >= CT_THREE_LINE_TAKE {
		if nType, isType = dg.isThreeLineTake(lastCardData); isType {
			if (firstType & 0xF) == (nType & 0xF) {
				return nType, (nType & 0xF0) > (firstType & 0xF0)
			}
		}

		return CT_ERROR, false
	}
	log.Debug("比三顺子")
	// 三顺子
	if firstType >= CT_THREE_LINE {
		if nType, isType = dg.isThreeLine(lastCardData); isType && ((firstType & 0xF0) == (nType & 0xF0)) {
			return nType, (nType & 0xf) > (firstType & 0xf)
		}

		return CT_ERROR, false
	}
	log.Debug("比双顺子")
	// 双顺子
	if firstType >= CT_DOUBLE_LINE {
		if nType, isType = dg.isDoubleLine(lastCardData); isType && ((firstType & 0xF0) == (nType & 0xF0)) {
			return nType, (nType & 0xf) > (firstType & 0xf)
		}
		return CT_ERROR, false
	}
	log.Debug("比单顺子")
	// 单顺子
	if firstType >= CT_SINGLE_LINE {
		if nType, isType = dg.isSingleLine(lastCardData); isType && ((firstType & 0xF0) == (nType & 0xF0)) {
			return nType, (nType & 0xf) > (firstType & 0xf)
		}
		return CT_ERROR, false
	}
	log.Debug("比三带二")
	// 三带二
	if firstType >= CT_THREE_TAKE_TWO {
		if nType, isType = dg.isThreeTakeTwo(lastCardData); isType {
			return nType, (nType & 0xFF) > (firstType & 0xFF)
		}
		return CT_ERROR, false
	}
	log.Debug("比3带1")
	// 三带一
	if firstType >= CT_THREE_TAKE_ONE {
		if nType, isType = dg.isThreeTakeOne(lastCardData); isType {
			return nType, (nType & 0xFF) > (firstType & 0xFF)
		}
		return CT_ERROR, false
	}
	log.Debug("比3张牌")
	// 三张牌
	if firstType >= CT_THREE {
		if nType, isType = dg.isThree(lastCardData); isType {
			return nType, (nType & 0xFF) > (firstType & 0xFF)
		}
		return CT_ERROR, false
	}
	log.Debug("比对子")
	// 对子
	if firstType >= CT_DOUBLE {
		if nType, isType = dg.isDouble(lastCardData); isType {
			return nType, (nType & 0xFF) > (firstType & 0xFF)
		}
	}
	log.Debug("比单牌")
	// 单
	if firstType >= CT_SINGLE {
		if dg.GetCardLogicValue(lastCardData[0]) > dg.GetCardLogicValue(firstCardData[0]) {
			return CT_SINGLE | dg.GetCardLogicValue(lastCardData[0]), true
		}
		return CT_ERROR, false
	}

	return CT_ERROR, false
}

// max函数
func (dg *ddz_logic) maxValue(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 获取某一牌值的个数
func (dg *ddz_logic) getCountWithCardValue(cardArr []int, v int) int {
	var nCount int
	for _, v1 := range cardArr {
		if dg.GetCardValue(v1) == v {
			nCount++
		}
	}
	return nCount
}

// 去掉牌中的大小王
func (dg *ddz_logic) removeKingFromCard(cardsArr []int) ([]int, int) {
	cardArr := util.CopySlicInt(cardsArr)
	var nCount int

	for i, v := range cardArr {
		if v >= 0x4e {
			nCount++
			copy(cardArr[i:], cardArr[i+1:])
			cardArr = cardArr[:len(cardArr)-1]
		}
	}

	return cardArr, nCount
}

// 去掉某个值的牌
func (dg *ddz_logic) removeValuesFromCard(cardArr []int, cardValue int) ([]int, int) {
	tmpArr := util.CopySlicInt(cardArr)
	var nCount int

	for i, v := range tmpArr {
		if dg.GetCardValue(v) == cardValue {
			nCount++
			copy(tmpArr[i:], tmpArr[i+1:])
			tmpArr = tmpArr[:len(tmpArr)-1]
		}
	}
	return tmpArr, nCount
}

// 获取除癞子牌外的最大逻辑值
func (dg *ddz_logic) getMaxLogicCardValueWithoutLaizi(cardArr []int) int {
	var cardValue int

	for _, v := range cardArr {
		v1 := dg.GetCardValue(v)
		if v1 != dg.LizeCard {
			cardValue = dg.maxValue(cardValue, dg.GetCardLogicValue(v))
		}
	}

	return cardValue
}
