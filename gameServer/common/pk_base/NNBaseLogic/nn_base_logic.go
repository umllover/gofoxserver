package NNBaseLogic

//import "mj/gameServer/common/pk_base"
import (

	"mj/gameServer/common/pk_base"

	"github.com/lovelly/leaf/log"

	"math"
)


// 牛牛类通用逻辑

const  (
	OX_VALUE0     	= 				0				//混合牌型
	OX_THREE_SAME 	= 			105				//小牛牛――5张牌都小于5（含5），并且5张牌相加不大于10
	OX_FOUR_SAME  	= 			104				////炸弹――5张牌中有4张一样的牌。
	OX_FOURKING   	= 				102				//天王牌型四花
	OX_FIVEKING   	= 				103				//天王牌型五花
)

const (
	NN_GAME_PLAYER		 =				4					//游戏人数
	NN_MAX_COUNT		 =				5					//最大数目
)

//分析结构
type AnalyseResult struct {
	FourCount			 int			//四张数目
	ThreeCount			 int			//三张数目
	DoubleCount			 int			//两张数目
	SignedCount			 int			//单张数目
	FourLogicVolue        []	 int			//四张列表
	ThreeLogicVolue        []	 int			//三张列表
	DoubleLogicVolue        []	 int			//两张列表
	SignedLogicVolue        []	 int			//单张列表
	FourCardData        []		 int	//四张列表
	ThreeCardData        []		 int	//三张列表
	DoubleCardData        []	 int	//两张列表
	SignedCardData        []	 int	//单张数目

}

//获取牛牛牌值
func NNGetCardLogicValue( CardData int) int {
	//扑克属性
	//CardColor = GetCardColor(CardData)
	CardValue := pk_base.GetCardValue(CardData)

	//转换数值
	//return (CardValue>10)?(10):CardValue
	if CardValue > 10 {
		CardValue = 10
	}
	return CardValue
	math.Abs()
}




//获取牛牛牌型
func NNGetCardType(CardData []int, CardCount int) int {

	if CardCount != NN_MAX_COUNT {
		return 0
	}

	////炸弹牌型
	//SameCount := 0

	var Temp [NN_MAX_COUNT]int
	Sum :=0
	for i:=0; i<CardCount; i++ {
		Temp[i] = NNGetCardLogicValue(CardData[i])
		log.Debug("%d", Temp[i])
		Sum += Temp[i]
	}
	log.Debug("%d", Sum)

	//王的数量
	KingCount := 0
	TenCount := 0

	for i:=0; i<CardCount; i++ {
		if pk_base.GetCardValue(CardData[i])>10 && CardData[i]!=0x4E && CardData[i]!=0x4F {
			KingCount++
		}else if(pk_base.GetCardValue(CardData[i])==10) {
			TenCount++
		}
	}

	if KingCount == NN_MAX_COUNT {
		return OX_FIVEKING//五花――5张牌都是10以上（不含10）的牌。。
	}

	Value := NNGetCardLogicValue(CardData[3])
	Value += NNGetCardLogicValue(CardData[4])

	if Value>10 {
		if CardData[3]==0x4E||CardData[4]==0x4F||CardData[4]==0x4E||CardData[3]==0x4F {
			Value=10
		}else {
			Value-=10 //2.3
		}

	}

	return Value//OX_VALUE0
}



