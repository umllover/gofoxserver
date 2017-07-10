package pk_base


const (
	//税收定义
	REVENUE_BENCHMARK   = 0    //税收起点
	REVENUE_DENOMINATOR = 1000 //税收分母
	PERSONAL_ROOM_CHAIR = 8    //私人房间座子上椅子的最大数目
)

type HistoryScore struct {
	TurnScore    int
	CollectScore int
}

//分析子项
type TagAnalyseItem struct {
	CardEye    int     //牌眼扑克
	bMagicEye  bool    //牌眼是否是王霸
	WeaveKind  []int   //组合类型
	CenterCard []int   //中心扑克
	CardData   [][]int //实际扑克
}

//类型子项
type TagKindItem struct {
	WeaveKind  int   //组合类型
	CenterCard int   //中心扑克
	CardIndex  []int //扑克索引
}

const (
	NN_GAME_PLAYER = 4 //游戏人数
	NN_MAX_COUNT   = 5 //最大数目
)

//分析结构
type AnalyseResult struct {
	FourCount        int   //四张数目
	ThreeCount       int   //三张数目
	DoubleCount      int   //两张数目
	SignedCount      int   //单张数目
	FourLogicVolue   []int //四张列表
	ThreeLogicVolue  []int //三张列表
	DoubleLogicVolue []int //两张列表
	SignedLogicVolue []int //单张列表
	FourCardData     []int //四张列表
	ThreeCardData    []int //三张列表
	DoubleCardData   []int //两张列表
	SignedCardData   []int //单张数目

}


// 扑克通用逻辑
const (
	LOGIC_MASK_COLOR = 0xF0 //花色掩码
	LOGIC_MASK_VALUE = 0x0F //数值掩码
)

