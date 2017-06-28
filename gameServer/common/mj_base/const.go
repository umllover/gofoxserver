package mj_base

const (

	//动作标志
	WIK_NULL     = 0x00 //没有类型
	WIK_LEFT     = 0x01 //左吃类型
	WIK_CENTER   = 0x02 //中吃类型
	WIK_RIGHT    = 0x04 //右吃类型
	WIK_PENG     = 0x08 //碰牌类型
	WIK_GANG     = 0x10 //杠牌类型
	WIK_LISTEN   = 0x20 //听牌类型
	WIK_CHI_HU   = 0x40 //吃胡类型
	WIK_FANG_PAO = 0x80 //放炮
)

//逻辑掩码
const (
	MASK_COLOR = 0xF0 //花色掩码
	MASK_VALUE = 0x0F //数值掩码
)

//麻将限制行为
const (
	LimitChiHu = 1      //禁止吃胡
	LimitPeng  = 1 << 1 //禁止碰
	LimitGang  = 1 << 2 //禁止杠牌
)

//发牌状态
const (
	Not_Send     = iota //无
	OutCard_Send        //出牌后发牌
	Gang_Send           //杠牌后发牌
	BuHua_Send          //补花后发牌
)

//效验类型
const (
	EstimatKind_OutCard  = iota //出牌效验
	EstimatKind_GangCard        //杠牌效验
)

const (
	//逻辑掩码
	INDEX_REPLACE_CARD = 42

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

const (
	//扑克定义
	HEAP_FULL_COUNT = 28 //堆立全牌
	MAX_RIGHT_COUNT = 1  //最大权位DWORD个数
)

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
