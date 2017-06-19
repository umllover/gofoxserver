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
	msg.Processor.Register(&G2C_TBNN_CallBanker{})
	msg.Processor.Register(&G2C_TBNN_GameStart{})
	msg.Processor.Register(&G2C_TBNN_GameEnd{})
	msg.Processor.Register(&G2C_TBNN_AddScore{})
	msg.Processor.Register(&G2C_TBNN_SendCard{})
	msg.Processor.Register(&G2C_TBNN_AllCard{})
	msg.Processor.Register(&G2C_TBNN_PlayerExit{})
	msg.Processor.Register(&G2C_TBNN_Open_Card{})

	// ----------c2s------------
	msg.Processor.Register(&C2G_TBNN_CallScore{})
	msg.Processor.Register(&C2G_TBNN_AddScore{})
	msg.Processor.Register(&C2G_TBNN_CallBanker{})
	msg.Processor.Register(&C2G_TBNN_OxCard{})
	msg.Processor.Register(&C2G_TBNN_QIANG{})
}

// ------------ s2c ----------------
//---------- 游戏状态-------
type G2C_TBNN_StatusFree struct {
	CellScore	 			int64						//基础积分

	//历史积分
	TurnScore 				[]int64			//积分信息
	CollectScore 			[]int64			//积分信息
	GameRoomName			string						//房间名称

	CtrFlag					int							//操作标志
	MaxScoreTimes			int		 					//最大倍数

	//LONG								lAndroidMaxCellScore				//机器人可设置的最大底注
	//LONG								lAndroidMinCellScore				//机器人可设置的最小底注

	TimeOutCard				int							//出牌时间
	TimeOperateCard			int 						//操作时间
	TimeStartGame			int 						//开始时间

	PlayerCount				int							//玩家人数
	TimesCount				int								//倍数
	PlayMode				int								//游戏模式
	PlayCount				int								//游戏局数

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
	CellScore       		int64							//基础积分

	PlayStatus      		[]int          //用户状态
	DynamicJoin     		int                 //动态加入
	TurnMaxScore			int64					//最大下注
	TableScore      		[]int64			//下注数目
	BankerUser				int					//庄家用户

	HandCardData  			[][]int         //桌面扑克
	OxCard        			[]int				//通比牛牛数据

	TurnScore 				[]int64			//积分信息
	CollectScore 			[]int64			//积分信息
	GameRoomName			string						//房间名称

	CtrFlag					int							//操作标志
	MaxScoreTimes			int 					//最大倍数
	//lAndroidMaxCellScore				//机器人可设置的最大底注
	//lAndroidMinCellScore				//机器人可设置的最小底注

	IsOpenCard				[]bool			//用户是否摊牌
	CurrentPlayCount		int							    //房间已玩局数
	EachRoundScore			[][]int			//房间每局游戏比分
}

type G2C_TBNN_CallBanker struct {
	Qiang_Start				bool		//开始抢
	CallBanker				int			//叫庄用户
	FirstTimes				bool		//首次叫庄
}

//游戏开始
type G2C_TBNN_GameStart struct {
	CellScore				int64							//单元下注

	DrawMaxScore			int64							//最大下注
	TurnMaxScore			int64			//最大下注
	BankerUser				int				//庄家用户
}

//用户下注
type G2C_TBNN_AddScore struct {
	AddScoreUser			int				//加注用户
	AddScoreCount			int64			//加注数目
}

//游戏结束
type G2C_TBNN_GameEnd	struct {

	GameTax					[]int				//游戏税收
	GameScore				[]int64			//游戏得分
	CardData				[]int			//用户扑克
	AllbCardValue			[]int
	MMcbCardData			[][]int     	//用户扑克
}

//发牌数据包
type G2C_TBNN_SendCard struct {
	CardData				[][]int     	//用户扑克
}

//发牌数据包
type G2C_TBNN_AllCard struct {
	CardData				[][]int			//用户扑克
}

//用户退出
type G2C_TBNN_PlayerExit struct {
	PlayerID				int			//退出用户
}

//用户摊牌
type G2C_TBNN_Open_Card struct {
	PlayerID				int			//摊牌用户
	Open					int			//摊牌标志
	CardData				[][]int				//牌数据
}



// ----------c2s----------------
//用户叫分
type C2G_TBNN_CallScore struct {
	CallScore				int 		//叫分数目
}

type C2G_TBNN_QIANG	 struct {
	Qiang 					bool			//1 qiang    0 bu
}
//用户叫庄
type C2G_TBNN_CallBanker	struct {
	Banker					bool		//做庄标志
}

//用户加注
type C2G_TBNN_AddScore	struct {
	Score					int64			//加注数目
}

//用户摊牌
type C2G_TBNN_OxCard struct {
	OX						bool		//通比牛牛标志
	CardData				[]int					//用户扑克
}



