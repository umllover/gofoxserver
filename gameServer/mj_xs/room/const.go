package room

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
	//非胡类型
	CHK_NULL = 0x0000 //非胡类型

	INDEX_REPLACE_CARD = 42

	//动作类型
	WIK_GANERAL   = 0x00 //普通操作
	WIK_MING_GANG = 0x01 //明杠（碰后再杠）
	WIK_FANG_GANG = 0x02 //放杠
	WIK_AN_GANG   = 0x03 //暗杠

	//胡牌牌型
	CHK_JI_HU     = 0x0001 //鸡胡类型
	CHK_PING_HU   = 0x0002 //平胡类型
	CHK_PENG_PENG = 0x0004 //碰碰胡牌
	CHK_BA_HUA    = 0x0008 //八花

	//胡牌权位
	CHR_NULL        = 0x0000 //无权位
	CHR_HUN_YI_SE   = 0x0001 //混一色
	CHR_GANG_FLOWER = 0x0002 //杠上开花
	CHR_HAI_DI      = 0x0004 //海底权位
	CHR_ZI_MO       = 0x0008 //自摸
	CHR_QING_YI_SE  = 0x0010 //清色权位
	CHR_MEN_QI      = 0x0020 //门清
	CHR_DI          = 0x0040 //地胡权位
	CHR_TIAN        = 0x0080 //天胡权位
	CHR_QIANG_GANG  = 0x0100 //抢杆权位
	CHR_DA_DIAO     = 0x0200 //大吊车
	CHR_BIAN        = 0x0400 //边
	CHR_QIAN        = 0x0800 //嵌
	CHR_DUI_DAO     = 0x1000 //对倒
	CHR_DAN_DIAO    = 0x2000 //单吊
	CHR_SI_HUA      = 0x4000 //四花
	CHR_BA_HUA      = 0x8000 //八花

	//胡牌类型  算台
	KIND_JI_HU     = 0  //鸡胡
	KIND_PING_HU   = 1  //平胡
	KIND_PENG_PEMG = 4  //碰碰胡
	KIND_BA_HUA    = 13 //八花胡

	//胡牌权位 算台
	RIGHT_HUN_YI_SE   = 3  //混一色
	RIGHT_GANG_FLOWER = 1  //杠上花
	RIGHT_HAI_DI      = 1  //海底捞月
	RIGHT_ZI_MO       = 1  //自摸
	RIGHT_QING_YI_SE  = 6  //清一色
	RIGHT_MEN_QI      = 1  //门清
	RIGHT_DI          = 7  //地胡
	RIGHT_TIAN        = 13 //天胡
	RIGHT_QIANG_GANG  = 0  //抢杠
	RIGHT_DA_DIAO     = 2  //大吊车
	RIGHT_BIAN        = 1  //边
	RIGHT_QIAN        = 1  //嵌
	RIGHT_DUI_DAO     = 1  //对倒
	RIGHT_DAN_DIAO    = 1  //单吊
	RIGHT_SI_HUA      = 7  //四花
	RIGHT_BA_HUA      = 16 //八花

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
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, //番子
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, //番子
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, //番子
	0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, //番子
	0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, //花子
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

//胡牌结果
type TagChiHuResult struct {
	ChiHuKind  int //吃胡类型
	ChiHuRight int //胡牌权位
}
