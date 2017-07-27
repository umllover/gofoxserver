package room

//游戏积分制
const (
	GAME_TYPE_33 = 0
	GAME_TYPE_48 = 1
	GAME_TYPE_88 = 2
)

//积分类型
const (
	IDX_SUB_SCORE_JC = 0 //基础分(底分1台)
	//桌面分
	IDX_SUB_SCORE_LZ  = 1 //连庄
	IDX_SUB_SCORE_HUA = 2 //花牌
	IDX_SUB_SCORE_AG  = 3 //暗杠
	IDX_SUB_SCORE_FB  = 4 //分饼
	IDX_SUB_SCORE_ZH  = 5 //抓花
	//胡牌+分
	IDX_SUB_SCORE_ZPKZ = 6  //字牌刻字
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
	IDX_SUB_SCORE_ZYS   = 35 //字一色
	IDX_SUB_SCORE_ZPG   = 36 //字牌杠

	COUNT_KIND_SCORE = 37 //分数子项个数

)
