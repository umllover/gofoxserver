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
)

//房间错误码 100 ~ 200
const (
	RoomFull = 101 //房间满了，不能再创建
	NotFoudGameType = 102 //玩家不存在
	CreateParamError = 103  //参数错误
	NoFoudTemplate = 104 //配置没找到
	ConfigError = 105 //配置错误
)




//////////////////////////////////////////////
//标识前缀
const (
	HallPrefix  = "HallSvr" //房间服
	GamePrefix = "GameSvr"
)



/////
func RenderErrorMessage(code int) *msg.ShowErrCode {
	return &msg.ShowErrCode{
		ErrCode:code,
	}
}

func GetGameSvrName(sververId int) string{
	return fmt.Sprintf(GamePrefix +"_%d", sververId)
}
func GetHallSvrName(sververId int) string{
	return fmt.Sprintf(HallPrefix +"_%d", sververId)
}

