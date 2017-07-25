package mj

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
	LimitChi   = 1 << 3 //禁止吃牌
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
	CHR_SAN_AN_KE       = 0x00000004 //三暗刻
	CHR_SI_AN_KE        = 0x00000008 //四暗刻
	CHR_WU_AN_KE        = 0x00000010 //五暗刻
	CHR_DAN_DIAO        = 0x00000020 //单吊
	CHR_ZI_PAI_GANG     = 0x00000040 //字牌杠
	CHR_ZI_YI_SE        = 0x00000080 //字一色
	CHR_TIAN_HU         = 0x00000100 //天胡
	CHR_DI_HU           = 0x00000200 //地胡
	CHR_QI_SHOU_LISTEN  = 0x00000400 //起首听
	CHR_HUA_SHANG_HUA   = 0x00000800 //花上开花
	CHR_HAI_DI_LAO_ZHEN = 0x00001000 //海底捞针
	CHR_ZI_KE_PAI       = 0x00002000 //字牌刻字
	CHR_HUA_GANG        = 0x00004000 //花杠
	CHR_WU_HUA_ZI       = 0x00008000 //无花字
	CHR_XIAO_SI_XI      = 0x00010000 //小四喜
	CHR_DA_SI_XI        = 0x00020000 //大四喜
	CHR_XIAO_SAN_YUAN   = 0x00040000 //小三元
	CHR_DA_SAN_YUAN     = 0x00080000 //大三元
	CHR_HUN_YI_SE       = 0x00100000 //混一色
	CHR_QING_YI_SE      = 0x00200000 //清一色
	CHR_HUA_YI_SE       = 0x00400000 //花一色
	CHR_BAI_LIU         = 0x00800000 //佰六
	CHR_GANG_SHANG_HUA  = 0x00800000 //杠上花
	CHR_GANG_SHANG_PAO  = 0x01000000 //杠上炮
	CHR_QIANG_GANG_HU   = 0x02000000 //抢杠胡
	CHR_CHI_HU          = 0x04000000 //放炮
	CHR_ZI_MO           = 0x08000000 //自摸
	CHR_QING_BAI_LIU    = 0x10000000 ////门清佰六
	CHR_WEI_ZHANG       = 0x20000000 //胡尾张
	CHR_JIE_TOU         = 0x40000000 //截头
	CHR_KONG_XIN        = 0x80000000 // 空心
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
	CardEye      int     //牌眼扑克
	MagicEye     bool    //牌眼是否是王霸
	WeaveKind    []int   //组合类型
	IsAnalyseGet []bool  //非打出组合
	CenterCard   []int   //中心扑克
	CardData     [][]int //实际扑克
	Param        []int   //类型标志
}

//类型子项
type TagKindItem struct {
	WeaveKind    int   //组合类型
	IsAnalyseGet bool  //非打出组合
	CenterCard   int   //中心扑克
	CardIndex    []int //扑克索引
	MagicCount   int   //
}
