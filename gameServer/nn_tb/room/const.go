package room

const (
	//用户状态
	US_NULL    = 0x00 //没有状态
	US_FREE    = 0x01 //站立状态
	US_SIT     = 0x02 //坐下状态
	US_READY   = 0x03 //同意状态
	US_LOOKON  = 0x04 //旁观状态
	US_PLAYING = 0x05 //游戏状态
	US_OFFLINE = 0x06 //断线状态
)

const (
	//房间状态
	RoomStatusReady    = 0
	RoomStatusStarting = 1
	RoomStatusEnd      = 2
)

const (
	//发牌状态

	Not_Send     = iota //无
	OutCard_Send        //出牌后发牌
	Gang_Send           //杠牌后发牌
	BuHua_Send          //补花后发牌

)



//常量定义
const (
	MAX_WEAVE     = 4   //最大组合
	MAX_COUNT     = 14  //最大数目
	MAX_REPERTORY = 112 //最大库存
	MAX_HUA_INDEX = 0   //花牌索引
	MAX_HUA_COUNT = 8   //花牌个数
)

const (
	//扑克定义
	HEAP_FULL_COUNT = 28 //堆立全牌
	MAX_RIGHT_COUNT = 1  //最大权位DWORD个数
)

const (
	//结束原因
	GER_NORMAL        = 0x00 //常规结束
	GER_DISMISS       = 0x01 //游戏解散
	GER_USER_LEAVE    = 0x02 //用户离开
	GER_NETWORK_ERROR = 0x03 //网络错误
)

const (
	//分数模式
	SCORE_GENRE_NORMAL   = 0x0100 //普通模式
	SCORE_GENRE_POSITIVE = 0x0200 //非负模式
)

const (
	//积分类型
	SCORE_TYPE_NULL    = 0x00 //无效积分
	SCORE_TYPE_WIN     = 0x01 //胜局积分
	SCORE_TYPE_LOSE    = 0x02 //输局积分
	SCORE_TYPE_DRAW    = 0x03 //和局积分
	SCORE_TYPE_FLEE    = 0x04 //逃局积分
	SCORE_TYPE_PRESENT = 0x10 //赠送积分
	SCORE_TYPE_SERVICE = 0x11 //服务积分
)

const (
	//税收定义
	REVENUE_BENCHMARK   = 0    //税收起点
	REVENUE_DENOMINATOR = 1000 //税收分母
	PERSONAL_ROOM_CHAIR = 8    //私人房间座子上椅子的最大数目
)

const (
	INDEX_REPLACE_CARD = 42

	//逻辑掩码

	//动作类型
	WIK_GANERAL   = 0x00 //普通操作
	WIK_MING_GANG = 0x01 //明杠（碰后再杠）
	WIK_FANG_GANG = 0x02 //放杠
	WIK_AN_GANG   = 0x03 //暗杠

	//胡牌定义
	CHR_PING_HU         = 0x00000001 //平胡
	CHR_PENG_PENG       = 0x00000002 //碰碰胡
	CHR_DAN_DIAN_QI_DUI = 0x00000004 //单点七对
	CHR_MA_QI_DUI       = 0x00000008 //麻七对
	CHR_MA_QI_WANG      = 0x00000010 //麻七王
	CHR_MA_QI_WZW       = 0x00000020 //麻七王中王
	CHR_SHI_SAN_LAN     = 0x00000040 //十三烂
	CHR_QX_SHI_SAN_LAN  = 0x00000080 //七星十三烂
	CHR_TIAN_HU         = 0x00000100 //天胡
	CHR_DI_HU           = 0x00000200 //地胡
	CHR_QI_SHOU_LISTEN  = 0x00000400 //起首听

	CHR_GANG_SHANG_HUA = 0x00800000 //杠上花
	CHR_GANG_SHANG_PAO = 0x01000000 //杠上炮
	CHR_QIANG_GANG_HU  = 0x02000000 //抢杠胡
	CHR_CHI_HU         = 0x04000000 //放炮
	CHR_ZI_MO          = 0x08000000 //自摸
)

//红中麻将限制行为
const (
	LimitChiHu = 1      //禁止吃胡
	LimitPeng  = 1 << 1 //禁止碰
	LimitGang  = 1 << 2 //禁止杠牌
)

//效验类型
const (
	EstimatKind_OutCard  = iota //出牌效验
	EstimatKind_GangCard        //杠牌效验
)

//麻将数据
var CardDataArray = []int{
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, //万子
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, //万子
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, //万子
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, //万子
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, //索子
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, //索子
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, //索子
	0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, //索子
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, //同子
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, //同子
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, //同子
	0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, //同子
	0x35, 0x35, 0x35, 0x35, //红中
}

type HistoryScore struct {
	TurnScore    int
	CollectScore int
}

//杠牌结果
type TagGangCardResult struct {
	CardCount int   //扑克数目
	CardData  []int //扑克数据
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
	ValidIndex []int //实际扑克索引 3
	MagicCount int   //财神牌数
}
