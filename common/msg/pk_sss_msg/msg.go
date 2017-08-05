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
	Processor.Register(&G2C_SSS_COMPARE{})
	Processor.Register(&G2C_SSS_StatusPlay{})
	Processor.Register(&CMD_C_ShowCard{})
	Processor.Register(&C2G_SSS_Open_Card{})
	Processor.Register(&G2C_SSS_Open_Card{})
	Processor.Register(&G2C_SSS_Record{})

}

type G2C_SSS_StatusFree struct {
	SubCmd           int `json:"subCmd"`           //当前状态
	PlayerCount      int `json:"playerCount"`      //实际人数
	CurrentPlayCount int `json:"currentPlayCount"` //当前局数
	MaxPlayCount     int `json:"maxPlayCount"`     //总局数
}

//发送扑克
type G2C_SSS_SendCard struct {
	//wCurrentUser int   //当前玩家
	CardData    []int //手上扑克
	Laizi       []int //癞子牌
	PublicCards []int //公共牌
	CellScore   int   //游戏底分
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
	CurrentUser int //当前玩家
}

//游戏结束
type G2C_SSS_COMPARE struct {
	LGameTax               int        `json:"lGameTax"`               //游戏税收
	LGameEveryTax          []int      `json:"lGameEveryTax"`          //每个玩家的税收
	LGameScore             []int      `json:"lGameScore"`             //游戏积分
	BEndMode               int        `json:"bEndMode"`               //结束方式
	CbCompareResult        [][]int    `json:"cbCompareResult"`        //每一道比较结果
	CbSpecialCompareResult []int      `json:"cbSpecialCompareResult"` //特殊牌型比较结果
	CbCompareDouble        []int      `json:"cbCompareDouble"`        //翻倍的道数
	CbUserOverTime         []int      `json:"cbUserOverTime"`         //玩家超时得到的道数
	CbCardData             [][]int    `json:"cbCardData"`             //扑克数据
	BUnderScoreDescribe    []string   `json:"bUnderScoreDescribe"`    //底分描述
	BCompCardDescribe      [][]string `json:"bCompCardDescribe"`      //牌比描述
	BToltalWinDaoShu       []int      `json:"bToltalWinDaoShu"`       //总共道数
	LUnderScore            int        `json:"lUnderScore"`            //底注分数
	BAllDisperse           []bool     `json:"bAllDisperse"`           //所有散牌
	BOverTime              []bool     `json:"bOverTime"`              //超时状态
	BUserLeft              []bool     `json:"bUserLeft"`              //玩家逃跑
	BLeft                  bool       `json:"bLeft"`
	LeftszName             []string   `json:"leftszName"`
	LeftChairID            []int      `json:"leftChairID"`
	BAllLeft               bool       `json:"bAllLeft"`
	LeftScore              []int      `json:"leftScore"`
	BSpecialCard           []bool     `json:"bSpecialCard"`      //是否为特殊牌
	BAllSpecialCard        bool       `json:"bAllSpecialCard"`   //全是特殊牌
	NTimer                 int        `json:"nTimer"`            //结束后比牌、打枪时间
	ShootState             [][]int    `json:"shootState"`        //赢的玩家,输的玩家 2为赢的玩家，1为全输的玩家，0为没输没赢的玩家
	M_nXShoot              int        `json:"m_nXShoot"`         //几家打枪
	CbThreeKillResult      []int      `json:"cbThreeKillResult"` //全垒打加减分
	BEnterExit             bool       `json:"bEnterExit"`        //是否一进入就离开
	WAllUser               int        `json:"wAllUser"`          //全垒打用户
}

// 结算
type G2C_SSS_Record struct {
	AllResult [][]int `json:"allResult"` //每一局总分
	AllScore  []int   `json:"allScore"`  //总分
}

//游戏状态
type G2C_SSS_StatusPlay struct {
	WCurrentUser       int             `json:"wCurrentUser"`       //当前玩家
	LCellScore         int             `json:"lCellScore"`         //单元底分
	NChip              []int           `json:"nChip"`              //下注大小
	BHandCardData      []int           `json:"bHandCardData"`      //扑克数据
	BSegmentCard       [][]int         `json:"bSegmentCard"`       //分段扑克
	BFinishSegment     []bool          `json:"bFinishSegment"`     //完成分段
	WUserToltalChip    int             `json:"wUserToltalChip"`    //总共金币
	BOverTime          []bool          `json:"bOverTime"`          //超时状态
	BSpecialTypeTable1 []bool          `json:"bSpecialTypeTable1"` //是否特殊牌型
	BDragon1           []bool          `json:"bDragon1"`           //是否倒水
	BAllHandCardData   [][]int         `json:"bAllHandCardData"`   //所有玩家的扑克数据
	SGameEnd           G2C_SSS_COMPARE `json:"sGameEnd"`           //游戏结束数据
	PlayerCount        int             `json:"playerCount"`        //实际人数
	CurrentPlayCount   int             `json:"currentPlayCount"`   //当前局数
	MaxPlayCount       int             `json:"maxPlayCount"`       //总局数
	Laizi              []int           `json:"Laizi"`              //癞子牌
	PublicCards        []int           `json:"PublicCards"`        //公共牌
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
