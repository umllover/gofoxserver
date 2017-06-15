package nn_tb_msg

import (
	"mj/common/msg"
)


func init() {

	msg.Processor.Register(&G2C_TBNN_StatusFree{})
	msg.Processor.Register(&G2C_TBNN_StatusCall{})
	msg.Processor.Register(&G2C_TBNN_StatusScore{})
	msg.Processor.Register(&G2C_TBNN_StatusPlay{})
	msg.Processor.Register(&G2C_TBNN_CallBanker{})
	msg.Processor.Register(&G2C_TBNN_GameStart{})
	msg.Processor.Register(&G2C_TBNN_GameEnd{})
	msg.Processor.Register(&G2C_TBNN_AddScore{})
	msg.Processor.Register(&G2C_TBNN_SendCard{})
	msg.Processor.Register(&G2C_TBNN_PlayerExit{})
	msg.Processor.Register(&G2C_TBNN_Open_Card{})


	msg.Processor.Register(&C2G_TBNN_CallScore{})
	msg.Processor.Register(&C2G_TBNN_AddScore{})
	msg.Processor.Register(&C2G_TBNN_CallBanker{})
	msg.Processor.Register(&C2G_TBNN_OxCard{})
	msg.Processor.Register(&C2G_TBNN_QIANG{})
	//msg.Processor.Register(&C2G_TBNN_SPECIAL_CLIENT_REPORT{})
}


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

/*
//发牌数据包
type G2C_TBNN_AllCard
{
bool								bAICount[]
BYTE								cbCardData[][]	//用户扑克
}*/

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

//用户叫分
type C2G_TBNN_CallScore struct {
	CallScore				int64 		//叫分数目
}

//#define DEFAULT_CELLSCORE_WAIT_TIME     20//20秒
/*
#define CELL_SET		1
#define CELL_SET_WAIT	2
#define CELL_MEET		3
#define CELL_NOT_MEET	4 */

type C2G_TBNN_QIANG	 struct {
	Qiang 					bool			//1 qiang    0 bu
}
//用户叫庄
type C2G_TBNN_CallBanker	struct {
	Banker					bool		//做庄标志
}

/*
//终端类型
type C2G_TBNN_SPECIAL_CLIENT_REPORT  struct {
WORD                                wUserChairID                       //用户方位
} */

//用户加注
type C2G_TBNN_AddScore	struct {
	Score					int64			//加注数目
}

//用户摊牌
type C2G_TBNN_OxCard struct {
	OX						bool		//通比牛牛标志
	CardData				[]int					//用户扑克
}


//////////////////////////////////////////////////////////////////////////
//type C2G_TBNN_AdminReq
//{
//	BYTE cbReqType
//#define RQ_SET_WIN_AREA	1
//#define RQ_RESET_CONTROL	2
//#define RQ_PRINT_SYN		3
//	BYTE cbExtendData[20]			//附加数据
//}

//请求回复
//type G2C_TBNN_CommandResult
//{
//	BYTE cbAckType					//回复类型
//#define ACK_SET_WIN_AREA  1
//#define ACK_PRINT_SYN     2
//#define ACK_RESET_CONTROL 3
//	BYTE cbResult
//#define CR_ACCEPT  2			//接受
//#define CR_REFUSAL 3			//拒绝
//	BYTE cbExtendData[20]			//附加数据
//}
/*
#define IDM_ADMIN_COMMDN WM_USER+1000

//控制区域信息
type tagControlInfo
{
INT  nAreaWin		//控制区域
}

//服务器控制返回
#define	 S_CR_FAILURE				0		//失败
#define  S_CR_UPDATE_SUCCES			1		//更新成功
#define	 S_CR_SET_SUCCESS			2		//设置成功
#define  S_CR_CANCEL_SUCCESS		3		//取消成功
type G2C_TBNN_ControlReturns
{
BYTE cbReturnsType				//回复类型
BYTE cbControlArea	//控制区域
BYTE cbControlTimes			//控制次数
}


//客户端控制申请
#define  C_CA_UPDATE				1		//更新
#define	 C_CA_SET					2		//设置
#define  C_CA_CANCELS				3		//取消
//type C2G_TBNN_ControlApplication
//{
//	BYTE cbControlAppType			//申请类型
//	BYTE cbControlArea	//控制区域
//	BYTE cbControlTim


/*
//出操作
type C2G_TBNN_MJXS_OperateCard type {
	OperateCode int //操作代码
	OperateCard int //操作扑克
}

//出操作
type C2G_TBNN_MJXS_OutCard type {
	CardData int //扑克数据
}

//发送扑克
type G2C_MJXS_GameStart type {
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

type G2C_MJXS_OutCard type {
	OutCardUser int //出牌用户
	OutCardData int //出牌扑克
}

type G2C_MJXS_GameEnd type { //游戏结束
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
type G2C_MJXS_SendCard type {
	CardData    int  //扑克数据
	ActionMask  int  //动作掩码
	CurrentUser int  //当前用户
	Gang        bool //是否杠牌
}

//操作命令
type G2C_MJXS_OperateResult type {
	OperateUser int //操作用户
	ProvideUser int //供应用户
	OperateCode int //操作代码
	OperateCard int //操作扑克
}

//操作提示
type G2C_MJXS_OperateNotify type {
	ResumeUser int //还原用户
	ActionMask int //动作掩码
	ActionCard int //动作扑克
}

type G2C_MJXS_StatusFree type {
	CellScore  int //基础金币
	BankerUser int //庄家用户
}

//游戏状态
type G2C_MJXS_StatusPlay type {
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
type TagWeaveItem type {
	WeaveKind   int  //组合类型
	CenterCard  int  //中心扑克
	PublicCard  bool //公开标志
	ProvideUser int  //供应用户
}

*/

