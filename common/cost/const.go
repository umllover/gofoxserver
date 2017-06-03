package cost

import (
	"mj/common/msg"

	"fmt"
)

//login error code 0  ~ 100
const (
	NotFoudAccout = 1 //没找到账号
	ParamError = 2  //参数错误
	AlreadyExistsAccount = 3 //账号已经存在
	InsertAccountError = 4 //服务器内部错误
	LoadUserInfoError = 5 //玩家数据加载失败
	CreateUserError = 6 // 创建玩家失败
	ErrUserReLogin = 7 //重复登录
	ErrPasswd = 8 //密码错误
)

//房间错误码 100 ~ 200
const (
	RoomFull = 101 //房间满了，不能再创建
	NotFoudGameType = 102 //玩家不存在
	CreateParamError = 103  //参数错误
	NoFoudTemplate = 104 //配置没找到
	ConfigError = 105 //配置错误
	NotEnoughFee = 106 //代币不足
	RandRoomIdError = 107 //生成房间id失败
	MaxSoucrce = 108 // 低分太高
	ChairHasUser = 109  //位置有玩家， 不能坐下
	GameIsStart = 110 //游戏已经开始， 不能加入
	ErrNotOwner = 111 // 不是房主 没权限操作
	ErrNoSitdowm = 112 //请先坐下在操作
	ErrGameIsStart = 113 //游戏已开始，不能离开房间
)


///////// 无效的数字
const (
	//参数定义
	INVALID_CHAIR		=		0xFFFF								//无效椅子
	INVALID_TABLE		=		0xFFFF								//无效桌子
	INVALID_SERVER		=		0xFFFF								//无效房间
	INVALID_KIND		=		0xFFFF								//无效游戏
)

///////////////游戏模式.
const (
	GAME_GENRE_GOLD			=	0x0001								//金币类型
	GAME_GENRE_SCORE		=	0x0002								//点值类型
	GAME_GENRE_MATCH		=	0x0004								//比赛类型
	GAME_GENRE_EDUCATE		=	0x0008								//训练类型
	GAME_GENRE_PERSONAL		=	0x0010								//约战类型
)


//////////////////////////////////////////////
//标识前缀
const (
	HallPrefix  = "HallSvr" //房间服
	GamePrefix = "GameSvr"
)



/////
func RenderErrorMessage(code int, Desc... string) *msg.ShowErrCode {
	var des string
	if len(Desc) < 1 {
		des = fmt.Sprintf("请求错误, 错误码: %d", code)
	}else{
		des = Desc[0]
	}
	return &msg.ShowErrCode{
		ErrorCode:code,
		DescribeString: des,
	}
}

func GetGameSvrName(sververId int) string{
	return fmt.Sprintf(GamePrefix +"_%d", sververId)
}
func GetHallSvrName(sververId int) string{
	return fmt.Sprintf(HallPrefix +"_%d", sververId)
}

