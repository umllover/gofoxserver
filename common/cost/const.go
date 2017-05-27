package cost

import (
	"mj/common/msg"

	"fmt"
)

//login error code
const (
	NotFoudAccout = 1
	ParamError = 2
	AlreadyExistsAccount = 3
	InsertAccountError = 4
	LoadUserInfoError = 5
	CreateUserError = 6
)

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

