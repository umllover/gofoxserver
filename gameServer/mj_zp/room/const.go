package room

const (
	//动作标志
	WIK_NULL     = 0x00  //没有类型
	WIK_LEFT     = 0x01  //左吃类型
	WIK_CENTER   = 0x02  //中吃类型
	WIK_RIGHT    = 0x04  //右吃类型
	WIK_PENG     = 0x08  //碰牌类型
	WIK_GANG     = 0x10  //杠牌类型
	WIK_LISTEN   = 0x20  //听牌类型
	WIK_CHI_HU   = 0x40  //吃胡类型
	WIK_FANG_PAO = 0x80  //放炮
	WIK_CHI      = 0X100 //吃牌类型
)

//游戏积分制
const (
	GAME_TYPE_33 = 33
	GAME_TYPE_48 = 48
	GAME_TYPE_88 = 88
)

//积分类型
const (
	IDX_SUB_SCORE_JC = 0 //基础分(底分1台)
	//桌面分
	IDX_SUB_SCORE_LZ   = 1 //连庄
	IDX_SUB_SCORE_HUA  = 2 //花牌
	IDX_SUB_SCORE_AG   = 3 //暗杠
	IDX_SUB_SCORE_AK   = 4 //***
	IDX_SUB_SCORE_ZG   = 5 //***
	IDX_SUB_SCORE_ZPKZ = 6 //字牌刻字
	//胡牌+分
	IDX_SUB_SCORE_HP   = 7  //平胡
	IDX_SUB_SCORE_ZM   = 8  //自摸(自摸+1台)
	IDX_SUB_SCORE_HDLZ = 9  //海底捞针(算自摸，不能额外加自摸分)
	IDX_SUB_SCORE_GSKH = 10 //杠上开花(算自摸，不能额外加自摸分)
	IDX_SUB_SCORE_HSKH = 11 //花上开花(算自摸，不能额外加自摸分)
	//额外+分
	IDX_SUB_SCORE_QYS = 12 //清一色
	IDX_SUB_SCORE_HYS = 13 //花一色
	IDX_SUB_SCORE_CYS = 14 //混一色
	IDX_SUB_SCORE_DSY = 15 //大三元
	IDX_SUB_SCORE_XSY = 16 //小三元
	IDX_SUB_SCORE_DSX = 17 //大四喜
	IDX_SUB_SCORE_XSX = 18 //小四喜
	IDX_SUB_SCORE_WHZ = 19 //无花字
	IDX_SUB_SCORE_DDH = 20 //对对胡
	IDX_SUB_SCORE_MQQ = 21 //门前清
	IDX_SUB_SCORE_BL  = 22 //佰六

	IDX_SUB_SCORE_QGH = 23 //抢杠胡
	IDX_SUB_SCORE_DH  = 24 //地胡
	IDX_SUB_SCORE_TH  = 25 //天胡

	IDX_SUB_SCORE_DDPH = 26 //单吊平胡
	IDX_SUB_SCORE_WDD  = 27 //尾单吊

	IDX_SUB_SCORE_KX    = 28 //空心
	IDX_SUB_SCORE_JT    = 29 //截头
	IDX_SUB_SCORE_DDZM  = 30 //单吊自摸
	IDX_SUB_SCORE_MQBL  = 31 //门清佰六
	IDX_SUB_SCORE_SANAK = 32 //三暗刻
	IDX_SUB_SCORE_SIAK  = 33 //四暗刻
	IDX_SUB_SCORE_WUAK  = 34 //五暗刻

	COUNT_KIND_SCORE = 35 //分数子项个数

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
	//动作类型
	WIK_GANERAL   = 0x00 //普通操作
	WIK_MING_GANG = 0x01 //明杠（碰后再杠）

	//胡牌定义
	CHR_PING_HU = 0x00000001 //平胡
)

const (
	//扑克定义
	HEAP_FULL_COUNT = 28 //堆立全牌
)

//类型子项
type TagKindItem struct {
	WeaveKind    int   //组合类型
	IsAnalyseGet bool  //非打出组合
	CenterCard   int   //中心扑克
	CardIndex    []int //扑克索引
}
