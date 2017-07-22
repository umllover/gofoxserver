package pk_sss_msg

import (
	"mj/common/msg"
)

var (
	Processor = msg.Processor
)

func init() {
	//sss msg
	Processor.Register(&G2C_SSS_StatusFree{})
	Processor.Register(&G2C_SSS_SendCard{})
	Processor.Register(&CMD_SSS_GameEnd{})
	Processor.Register(&CMD_S_StatusPlay{})
	Processor.Register(&CMD_C_ShowCard{})
	Processor.Register(&C2G_SSS_Open_Card{})
	Processor.Register(&G2C_SSS_Open_Card{})
}

type G2C_SSS_StatusFree struct {
	//历史积分
	lTurnScore    []int //积分信息
	lCollectScore []int //积分信息

	wUserToltalChip []int
}

//发送扑克
type G2C_SSS_SendCard struct {
	//wCurrentUser int   //当前玩家
	CardData  []int //手上扑克
	CellScore int   //游戏底分
}

//用户摊牌
type C2G_SSS_Open_Card struct {
	FrontCard   []int //前墩扑克
	MidCard     []int //中墩扑克
	BackCard    []int //后墩扑克
	SpecialType bool  //是否是特殊牌
	SpecialData []int //特殊扑克
	Dragon      bool  //是否乌龙
}

//用户摊牌
type G2C_SSS_Open_Card struct {
	CurrentUser    int   //当前玩家
	FrontCard      []int //前墩扑克
	MidCard        []int //中墩扑克
	BackCard       []int //后墩扑克
	CanSeeShowCard bool  //能否看牌
	SpecialType    bool  //是否是特殊牌
	SpecialData    []int //特殊扑克
	ShowUser       int   //摊牌的玩家
	Dragon         bool  //是否乌龙
}

//游戏结束
type CMD_SSS_GameEnd struct {
	LGameTax               int        //游戏税收
	LGameEveryTax          []int      //每个玩家的税收
	LGameScore             []int      //游戏积分
	BEndMode               int        //结束方式
	CbCompareResult        [][]int    //每一道比较结果
	CbSpecialCompareResult []int      //特殊牌型比较结果
	CbCompareDouble        []int      //翻倍的道数
	CbUserOverTime         []int      //玩家超时得到的道数
	CbCardData             [][]int    //扑克数据
	BUnderScoreDescribe    [][]int    //底分描述
	BCompCardDescribe      [][][]int  //牌比描述
	BToltalWinDaoShu       []int      //总共道数
	LUnderScore            int        //底注分数
	BAllDisperse           []bool     //所有散牌
	BOverTime              []bool     //超时状态
	BUserLeft              []bool     //玩家逃跑
	BLeft                  bool       //
	LeftszName             [][]string //
	LeftChairID            []int      //
	BAllLeft               bool       //
	LeftScore              []int      //
	BSpecialCard           []bool     //是否为特殊牌
	BAllSpecialCard        bool       //全是特殊牌
	NTimer                 int        //结束后比牌、打枪时间
	ShootState             [][]int    //赢的玩家,输的玩家 2为赢的玩家，1为全输的玩家，0为没输没赢的玩家
	M_nXShoot              int        //几家打枪
	CbThreeKillResult      []int      //全垒打加减分
	BEnterExit             bool       //是否一进入就离开
	WAllUser               int        //全垒打用户
}

//游戏状态
type CMD_S_StatusPlay struct {
	wCurrentUser       int             //当前玩家
	lCellScore         int             //单元底分
	nChip              []int           //下注大小
	bHandCardData      []int           //扑克数据
	bSegmentCard       [][][]int       //分段扑克
	bFinishSegment     []bool          //完成分段
	wUserToltalChip    int             //总共金币
	bOverTime          []bool          //超时状态
	bSpecialTypeTable1 []bool          //是否特殊牌型
	bDragon1           []bool          //是否倒水
	bAllHandCardData   [][]int         //所有玩家的扑克数据
	CMD_S_GameEnd      CMD_SSS_GameEnd //游戏结束数据
}

//分段信息
type CMD_C_ShowCard struct {
	bFrontCard    []int //前墩扑克
	bMidCard      []int //中墩扑克
	bBackCard     []int //后墩扑克
	bSpecialType  bool  //是否是特殊牌
	btSpecialData []int //特殊扑克
	bDragon       bool  //是否乌龙
}
