package nn_tb_msg

import (
	"mj/common/msg"
)


func init() {
	//-----------s2c-----------------
	msg.Processor.Register(&G2C_TBNN_StatusFree{})
	msg.Processor.Register(&G2C_TBNN_StatusCall{})
	msg.Processor.Register(&G2C_TBNN_StatusScore{})
	msg.Processor.Register(&G2C_TBNN_StatusPlay{})
	msg.Processor.Register(&G2C_TBNN_CallScoreEnd{})
	msg.Processor.Register(&G2C_TBNN_GameStart{})
	msg.Processor.Register(&G2C_TBNN_GameEnd{})
	msg.Processor.Register(&G2C_TBNN_AddScore{})
	msg.Processor.Register(&G2C_TBNN_SendCard{})
	msg.Processor.Register(&G2C_TBNN_AllCard{})
	msg.Processor.Register(&G2C_TBNN_PublicCard{})
	msg.Processor.Register(&G2C_TBNN_PlayerExit{})
	msg.Processor.Register(&G2C_TBNN_Open_Card{})
	msg.Processor.Register(&G2C_TBNN_CalScore{})
	msg.Processor.Register(&G2C_TBNN_CallScore{})

	// ----------c2s------------
	msg.Processor.Register(&C2G_TBNN_CallScore{})
	msg.Processor.Register(&C2G_TBNN_AddScore{})
	msg.Processor.Register(&C2G_TBNN_OpenCard{})
}

// ------------ s2c ----------------
//---------- 游戏状态-------
type G2C_TBNN_StatusFree struct {
	CellScore	 			int						//基础积分

	//历史积分
	TurnScore 				[]int			//积分信息
	CollectScore 			[]int			//积分信息
	GameRoomName			string						//房间名称

	CtrFlag					int							//操作标志
	MaxScoreTimes			int		 					//最大倍数


	TimeOutCard				int							//出牌时间
	TimeOperateCard			int 						//操作时间
	TimeStartGame			int64 						//开始时间

	PlayerCount				int							//玩家人数
	TimesCount				int								//倍数
	PlayMode				int								//游戏模式
	CountLimit				int								//局数限制

	CurrentPlayCount		int							    //房间已玩局数
	EachRoundScore			[][]int			//房间每局游戏比分
}

type G2C_TBNN_StatusCall struct {
	CallBanker				int						//叫庄用户
	DynamicJoin				int                      //动态加入
	PlayStatus				[]int          //用户状态

	//历史积分
	TurnScore 				[]int64			//积分信息
	CollectScore 			[]int64			//积分信息
	GameRoomName			string						//房间名称
}

type G2C_TBNN_StatusScore struct {
	//下注信息
	PlayStatusi     		[]int          //用户状态
	DynamicJoin     		int                      //动态加入
	TurnMaxScore			int64					//最大下注
	TableScore      		[]int64			//下注数目
	BankerUser				int					//庄家用户
	TurnScore 				[]int64			//积分信息
	CollectScore 			[]int64			//积分信息
	GameRoomName			string						//房间名称
}

type G2C_TBNN_StatusPlay struct {
	CellScore       		int							//基础积分

	PlayStatus      		[]int          //用户状态
	DynamicJoin     		int                 //动态加入
	TurnMaxScore			int64					//最大下注
	TableScore      		[]int64			//下注数目
	BankerUser				int					//庄家用户

	HandCardData  			[][]int         //桌面扑克
	OxCard        			[]int				//通比牛牛数据

	TurnScore 				[]int			//积分信息
	CollectScore 			[]int			//积分信息
	GameRoomName			string						//房间名称

	CtrFlag					int							//操作标志
	MaxScoreTimes			int 					//最大倍数

	IsOpenCard				[]bool			//用户是否摊牌
	CurrentPlayCount		int							    //房间已玩局数
	EachRoundScore			[][]int			//房间每局游戏比分
}

//叫分结果
type G2C_TBNN_CallScoreEnd struct {
	Banker     int //庄家用户
	ScoreTimes int //倍数
}

//游戏开始
type G2C_TBNN_GameStart struct {
	CellScore				int64							//单元下注

	DrawMaxScore			int64							//最大下注
	TurnMaxScore			int64			//最大下注
	BankerUser				int				//庄家用户
}

//广播用户下注
type G2C_TBNN_AddScore struct {
	ChairID 				int
	AddScoreCount			int				//加注数目
}

//游戏结束
type G2C_TBNN_GameEnd	struct {

	GameTax					[]int				//游戏税收
	GameScore				[]int64			//游戏得分
	CardData				[]int			//用户扑克
	AllbCardValue			[]int
	MMcbCardData			[][]int     	//用户扑克
}

// 比牌结果
type G2C_TBNN_CalScore struct {
	GameScore 			int 	//得分
	CardData  			[]int 	//手牌
}

//发牌数据包
type G2C_TBNN_SendCard struct {
	CardData				[][]int     	//用户扑克
}

//发牌数据包
type G2C_TBNN_AllCard struct {
	CardData				[][]int			//用户扑克
}

// 公共牌数据
type G2C_TBNN_PublicCard struct {
	PublicCardData			[]int			//公共牌
}

// 最后一张牌
type G2C_TBNN_LastCard struct {
	LastCard 				int 		// 最后一张牌
}

//用户退出
type G2C_TBNN_PlayerExit struct {
	PlayerID				int			//退出用户
}

//用户摊牌
type G2C_TBNN_Open_Card struct {
	ChairID					int			//摊牌用户
	CardType				int 		//牌型
	CardData				[]int				//牌数据
}

//广播用户叫分
type G2C_TBNN_CallScore struct {
	ChairID					int 		//叫分用户
	CallScore				int 		//叫分数目
}



// ----------c2s----------------
//用户叫分
type C2G_TBNN_CallScore struct {
	CallScore				int 		//叫分数目
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
type C2G_TBNN_AddScore	struct {
	Score					int			//加注数目
}

//用户摊牌
type C2G_TBNN_OpenCard struct {
	CardType				int 					//牌型
	CardData				[]int					//用户扑克
}



