package mj_xs_msg

import (
	"mj/common/msg"
)

func init() {
	msg.Processor.Register(&C2G_MJXS_OperateCard{})
	msg.Processor.Register(&C2G_MJXS_OutCard{})

	msg.Processor.Register(&G2C_GameStart{})
	msg.Processor.Register(&G2C_OutCard{})
	msg.Processor.Register(&G2C_GameConclude{})
	msg.Processor.Register(&G2C_SendCard{})
	msg.Processor.Register(&G2C_OperateResult{})
	msg.Processor.Register(&G2C_OperateNotify{})
	msg.Processor.Register(&G2C_StatusFree{})
	msg.Processor.Register(&G2C_StatusPlay{})
}

//出操作
type C2G_MJXS_OperateCard struct {
	OperateCode int //操作代码
	OperateCard int //操作扑克
}

//出操作
type C2G_MJXS_OutCard struct {
	CardData int //扑克数据
}

//发送扑克
type G2C_GameStart struct {
	SiceCount         int   //骰子点数
	BankerUser        int   //庄家用户
	CurrentUser       int   //当前用户
	UserAction        int   //用户动作
	SunWindCount      int   //总花番
	CardData          []int //扑克列表
	LeftCardCount     int   //剩余牌数
	First             bool  //是否首发
	FengQuan          int   //风圈
	InitialBankerUser int   //初始庄家
	RepertoryCard     []int //所有牌
}

type G2C_OutCard struct {
	OutCardUser int //出牌用户
	OutCardData int //出牌扑克
}

type G2C_GameConclude struct { //游戏结束
	GameTax     int     //游戏税收
	ChiHuCard   int     //吃胡扑克
	ProvideUser int     //点炮用户
	HaiDiCard   int     //海底扑克
	Feng        []int   //风番台数
	Hua         []int   //花台
	Zi          []int   //子台数
	All         int     //总台数
	GameScore   []int   //游戏积分
	ChiHuKind   []int   //胡牌类型
	CardCount   []int   //扑克数目
	CardData    [][]int //扑克数据
	ChiHuRight  []int   //胡牌权位
}

//发送扑克
type G2C_SendCard struct {
	CardData    int  //扑克数据
	ActionMask  int  //动作掩码
	CurrentUser int  //当前用户
	Gang        bool //是否杠牌
}

//操作命令
type G2C_OperateResult struct {
	OperateUser int //操作用户
	ProvideUser int //供应用户
	OperateCode int //操作代码
	OperateCard int //操作扑克
	ActionMask  int //操作码
}

//操作提示
type G2C_OperateNotify struct {
	ResumeUser int //还原用户
	ActionMask int //动作掩码
	ActionCard int //动作扑克
}

type G2C_StatusFree struct {
	CellScore  int //基础金币
	BankerUser int //庄家用户
}

//游戏状态
type G2C_StatusPlay struct {
	//游戏变量
	CellScore         int //单元积分
	SiceCount         int //骰子点数
	BankerUser        int //庄家用户
	CurrentUser       int //当前用户
	FengQuan          int //风圈
	InitialBankerUser int //初始庄家

	//状态变量
	ActionCard    int //动作扑克
	ActionMask    int //动作掩码
	LeftCardCount int //剩余数目

	//出牌信息
	EnjoinCardCount  int     //禁止吃牌
	EnjoinCardData   []int   //禁止吃牌
	OutCardUser      int     //出牌用户
	OutCardData      int     //出牌扑克
	DiscardCount     []int   //丢弃数目
	DiscardCard      [][]int //丢弃记录
	UserWindCount    []int   //风牌记录
	UserWindCardData [][]int //风牌记录

	//扑克数据
	CardCount int   //扑克数目
	CardData  []int //扑克列表

	//组合扑克
	WeaveCount     []int             //组合数目
	WeaveItemArray [][]*TagWeaveItem //组合扑克
}

//组合子项
type TagWeaveItem struct {
	WeaveKind   int  //组合类型
	CenterCard  int  //中心扑克
	PublicCard  bool //公开标志
	ProvideUser int  //供应用户
}
