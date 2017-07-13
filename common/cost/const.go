package cost

import (
	"fmt"
	"mj/common/msg"
)

//login error code 0  ~ 100
const (
	NotFoudAccout        = 1 //没找到账号
	ParamError           = 2 //参数错误
	AlreadyExistsAccount = 3 //账号已经存在
	InsertAccountError   = 4 //服务器内部错误
	LoadUserInfoError    = 5 //玩家数据加载失败
	CreateUserError      = 6 // 创建玩家失败
	ErrUserDoubleLogin   = 7 //重复登录
	ErrPasswd            = 8 //密码错误
)

//房间错误码 100 ~ 200
const (
	RoomFull                = 101 //房间满了，不能再创建
	NotFoudGameType         = 102 //玩家不存在
	ErrParamError           = 103 //参数错误
	NoFoudTemplate          = 104 //配置没找到
	ConfigError             = 105 //配置错误
	NotEnoughFee            = 106 //代币不足
	RandRoomIdError         = 107 //生成房间id失败
	MaxSoucrce              = 108 // 低分太高
	ChairHasUser            = 109 //位置有玩家， 不能坐下
	GameIsStart             = 110 //游戏已经开始， 不能加入
	ErrNotOwner             = 111 // 不是房主 没权限操作
	ErrNoSitdowm            = 112 //请先坐下在操作
	ErrGameIsStart          = 113 //游戏已开始，不能离开房间
	ErrCreateRoomFaild      = 114 //创建聊天室失败
	NotOwner                = 115 //不是房主
	Errunlawful             = 116 //非法操作
	ErrMaxRoomCnt           = 117 //房间超限， 不能再创建了
	ErrServerError          = 118 //服务器内部错误
	ErrNotFoudServer        = 119 //没有找到可以服务的server
	ErrNoFoudRoom           = 120 //房间没有找到
	ErrNotFoundCreateRecord = 121 //没有找到房间记录
	ErrDoubleCreaterRoom    = 122 //重复创建房间
	ErrCreaterError         = 123 //创建错误
)

//红中麻将错误码
const (
	NotValidCard     = 201 //无效的牌
	ErrUserNotInRoom = 202 //玩家不在房间
	ErrNotFoudCard   = 203 //没找到牌
	ErrGameNotStart  = 204 //游戏没开始
	ErrNotSelfOut    = 205 //不是自己出牌
	ErrNoOperator    = 206 //没有操作
)

/// 活动领取错误码
const (
	ErrNotFoudTemplate = 301 //没有找到模板
	ErrMaxDrawTimes    = 302 //领取次数上线
)

//设置推举人的错误码
const (
	ErrNotFoudPlayer = 401 //没找到推举人
)

///////// 无效的数字
const (
	//无效数值
	INVALID_BYTE   = 0xFF       //无效数值
	INVALID_WORD   = 0xFFFF     //无效数值
	INVALID_DWORD  = 0xFFFFFFFF //无效数值
	INVALID_CHAIR  = 0xFFFF     //无效椅子
	INVALID_TABLE  = 0xFFFF     //无效桌子
	INVALID_SERVER = 0xFFFF     //无效房间
	INVALID_KIND   = 0xFFFF     //无效游戏
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

///////////////游戏模式.
const (
	GAME_GENRE_GOLD     = 0x0001 //金币类型
	GAME_GENRE_SCORE    = 0x0002 //点值类型
	GAME_GENRE_MATCH    = 0x0004 //比赛类型
	GAME_GENRE_EDUCATE  = 0x0008 //训练类型
	GAME_GENRE_PERSONAL = 0x0010 //约战类型
)

/// 通用状态
const (
	//用户状态
	US_NULL    = 0x00 //没有状态
	US_FREE    = 0x01 //站立状态
	US_SIT     = 0x02 //坐下状态
	US_READY   = 0x03 //同意状态
	US_LOOKON  = 0x04 //旁观状态
	US_PLAYING = 0x05 //游戏状态
	US_OFFLINE = 0x06 //断线状态
)

const (
	//房间状态
	RoomStatusReady    = 0
	RoomStatusStarting = 1
	RoomStatusEnd      = 2
)

const (
	//结束原因
	GER_NORMAL        = 0x00 //常规结束
	GER_DISMISS       = 0x01 //游戏解散
	GER_USER_LEAVE    = 0x02 //用户离开
	GER_NETWORK_ERROR = 0x03 //网络错误
)

//积分修改类型
const (
	HZMJ_CHANGE_SOURCE = 1
	ZPMJ_CHANGE_SOURCE = 2
)

//自己支付
const (
	SELF_PAY_TYPE = 1 //自己付钱
	AA_PAY_TYPE   = 2 //AA付钱
)

//////////////////////////////////////////////
//标识前缀
const (
	HallPrefix    = "HallSvr"    //房间服
	GamePrefix    = "GameSvr"    //游戏服
	HallPrefixFmt = "HallSvr_%d" //房间服
	GamePrefixFmt = "GameSvr_%d" //游戏服
)

func GetGameSvrTopc(id int) string {
	return fmt.Sprintf(GamePrefixFmt, id)
}

func GetHallSvrTopc(id int) string {
	return fmt.Sprintf(HallPrefixFmt, id)
}

func LOBYTE(w int) int {
	return w & 0xFF
}
func HIBYTE(w int) int {
	return w & 0xFF00
}

/////
func RenderErrorMessage(code int, Desc ...string) *msg.ShowErrCode {
	var des string
	if len(Desc) < 1 {
		des = fmt.Sprintf("请求错误, 错误码: %d", code)
	} else {
		des = fmt.Sprintf(Desc[0]+", 错误码: %d", code)
	}
	return &msg.ShowErrCode{
		ErrorCode:      code,
		DescribeString: des,
	}
}

func GetGameSvrName(sververId int) string {
	return fmt.Sprintf(GamePrefix+"_%d", sververId)
}
func GetHallSvrName(sververId int) string {
	return fmt.Sprintf(HallPrefix+"_%d", sververId)
}

///////////////// global 常量 ///////////////////////

const (
	MAX_CREATOR_ROOM_CNT = "MAX_CREATOR_ROOM_CNT"
	MAX_ELECT_AWARD      = "MAX_ELECT_AWARD"
)
