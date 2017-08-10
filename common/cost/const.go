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
	ErrCreaterError         = 123 //创建房间失败
	ErrPlayerNotInRoom      = 124 //玩家不在房间内
	ErrLoveRoomFaild        = 125 //离开房间异常
	ErrPlayerIsReady        = 126 //玩家已经准备了
	ErrRenewalFee           = 127 //请先续费
	ErrRoomIsClose          = 128 //房间已经结束了
	ErrRoomFull             = 129 //房间已满
	ErrRenewalFeesFaild     = 130 //续费失败
	ErrRefuseLeave          = 131 //拒绝离开
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

// 斗地主错误代码
const (
	ErrDDZCSUser  = 501 // 叫分玩家错误
	ErrDDZCSValid = 502 // 叫分无效
)

//个人信息操作码
const (
	ErrNotFondCreatorRoom = 601 //没有找到要删除的房间
	ErrRoomIsStart        = 602 //房间已经开始了

	ErrFrequentAccess    = 603 //获取验证码太频繁了
	ErrRandMaskCodeError = 604 //获取验证码失败
	ErrMaskCodeNotFoud   = 605 //验证码没找到
	ErrMaskCodeError     = 606 //验证码失败
	ErrNotInRoom         = 607 //不在房间内
	ErrFindRoomError     = 608 //查找房间失败
	ErrConfigError       = 609 //配置错误
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
	GAME_GENRE_ZhuanShi = 1 // 比赛类型
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
	GER_NORMAL           = 0 //常规结束
	GER_DISMISS          = 1 //游戏解散
	USER_LEAVE           = 2 //玩家请求解散
	NO_START_GER_DISMISS = 3 //没开始就解散
)

const (
	//加入房间累型
	GIRPrivate = 0 //私房加入
	GIRPublic  = 1 //公房加入
)

const (
	//房间结束
	RoomErrorDismiss   = 1 //出错解散房间
	RoomNormalDistmiss = 2 //正常解散房间
)

const (
	//是否为他人开房
	CreateRoomForSelf   = 0 //为自己开房
	CreateRoomForOthers = 1 //为他人开房
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

//离线处理消息枚举
const (
	OfflineTypeDianZhan = "DianZhan"
	OfflineRoomEndInfo  = "RoomEndInfo"
	OfflineReturnMoney  = "RoomReturnMoney"
	OfflineAddElectId   = "OfflineAddElectId"
)

//踢出玩家原因
const (
	UserOffline        = 0 //socket 断开 主动断线
	ServerKick         = 1 //服务器主动踢出
	KickOutMsg         = 2 //踢号 重登
	KickOutGameEnd     = 3 //游戏结束，关闭房间踢出房间
	KickOutUnlawfulMsg = 4 //非法消息
)

//积分类型
const (
	IDX_SUB_SCORE_JC = 0 //基础分(底分1台)
	//桌面分
	IDX_SUB_SCORE_LZ  = 1 //连庄
	IDX_SUB_SCORE_HUA = 2 //花牌（补花）
	IDX_SUB_SCORE_AG  = 3 //暗杠
	IDX_SUB_SCORE_FB  = 4 //分饼
	IDX_SUB_SCORE_ZH  = 5 //抓花
	IDX_SUB_SCORE_ZF  = 6 //庄分
	//胡牌+分
	IDX_SUB_SCORE_ZPKZ = 7  //字牌刻字
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

//////////////////////////////////////////////
//标识前缀
const (
	HallPrefix     = "HallSvr"        //房间服
	GamePrefix     = "GameSvr"        //游戏服
	HallPrefixFmt  = "HallSvr_%d"     //房间服
	GamePrefixFmt  = "GameSvr_%d"     //游戏服
	GameChannelFmt = "GameChannel_%d" //nsq channel
	HallCahnnelFmt = "HallCahnnel_%d" //nsq channel
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
	return fmt.Sprintf(GamePrefixFmt, sververId)
}
func GetHallSvrName(sververId int) string {
	return fmt.Sprintf(HallPrefixFmt, sververId)
}

///////////////// global 常量 ///////////////////////

const (
	MAX_CREATOR_ROOM_CNT = "MAX_CREATOR_ROOM_CNT"
	MAX_ELECT_AWARD      = "MAX_ELECT_AWARD"
	MAX_SHOW_ENTRY       = "MAX_SHOW_ENTRY"
	MATCH_TIMEOUT        = "MATCH_TIMEOUT"
	MASK_CODE_TEXT       = "MASK_CODE_TEXT"
	DelayDestroyRoom     = "DelayDestroyRoom"
	LeaveRoomTimer       = "LeaveRoomTimer"
)

/////////////// Redis key //////////////////////
const (
	CreatorRoom = "CreatorRoom:%d"
)
