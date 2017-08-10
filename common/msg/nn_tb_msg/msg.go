package nn_tb_msg

import (
	"mj/common/msg"
)

func init() {
	//-----------s2c-----------------
	msg.Processor.Register(&G2C_TBNN_CallScoreEnd{})
	msg.Processor.Register(&G2C_TBNN_GameStart{})
	msg.Processor.Register(&G2C_TBNN_GameEnd{})
	msg.Processor.Register(&G2C_TBNN_AddScore{})
	msg.Processor.Register(&G2C_TBNN_SendCard{})
	msg.Processor.Register(&G2C_TBNN_AllCard{})
	msg.Processor.Register(&G2C_TBNN_PublicCard{})
	msg.Processor.Register(&G2C_TBNN_LastCard{})
	msg.Processor.Register(&G2C_TBNN_PlayerExit{})
	msg.Processor.Register(&G2C_TBNN_Open_Card{})
	msg.Processor.Register(&G2C_TBNN_CalScore{})
	msg.Processor.Register(&G2C_TBNN_CallScore{})

	// ----------c2s------------
	msg.Processor.Register(&C2G_TBNN_CallScore{})
	msg.Processor.Register(&C2G_TBNN_AddScore{})
	msg.Processor.Register(&C2G_TBNN_OpenCard{})
}

// ------------ g2c ----------------
//叫分结果
type G2C_TBNN_CallScoreEnd struct {
	Banker         int   //庄家用户
	ScoreTimes     int   //倍数
	ScoreTimesUser []int // 与专家叫一样分数的玩家
}

//游戏开始
type G2C_TBNN_GameStart struct {
	PlayerCount int //游戏人数
}

//广播用户下注
type G2C_TBNN_AddScore struct {
	ChairID       int
	AddScoreCount int //加注数目
}

//游戏结束
type G2C_TBNN_GameEnd struct {
	CurrentPlayCount int     //当前局数
	LimitPlayCount   int     //总局数
	InitScore        []int   //玩家积分
	EachRoundScore   [][]int //每局积分
	Reason           int     //结束原因
}

// 比牌结果
type G2C_TBNN_CalScore struct {
	GameTax   []int   //服务费
	GameScore []int   //得分
	CardType  []int   //牌型
	CardData  [][]int //手牌
	InitScore []int   //玩家积分信息
}

//发牌数据包
type G2C_TBNN_SendCard struct {
	CardData [][]int //用户扑克
}

//发牌数据包
type G2C_TBNN_AllCard struct {
	CardData [][]int //用户扑克
}

// 公共牌数据
type G2C_TBNN_PublicCard struct {
	PublicCardData []int //公共牌
}

// 最后一张牌
type G2C_TBNN_LastCard struct {
	LastCard [][]int // 最后一张牌
}

//用户退出
type G2C_TBNN_PlayerExit struct {
	PlayerID int //退出用户
}

//用户摊牌
type G2C_TBNN_Open_Card struct {
	ChairID  int   //摊牌用户
	CardType int   //牌型
	CardData []int //牌数据
}

//广播用户叫分
type G2C_TBNN_CallScore struct {
	ChairID   int //叫分用户
	CallScore int //叫分数目
}

// ----------c2s----------------
//用户叫分
type C2G_TBNN_CallScore struct {
	CallScore int //叫分数目
}

/*
type C2G_TBNN_QIANG	 struct {
	Qiang 					bool			//1 qiang    0 bu
}
//用户叫庄
type C2G_TBNN_CallBanker	struct {
	Banker					bool		//做庄标志
}
*/

//用户加注
type C2G_TBNN_AddScore struct {
	Score int //加注数目
}

//用户摊牌
type C2G_TBNN_OpenCard struct {
	CardType int   //牌型
	CardData []int //用户扑克
}
