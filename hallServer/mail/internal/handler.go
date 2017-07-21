package internal

import (
	. "mj/common/cost"
	"mj/common/register"
	"mj/hallServer/db/model"
	"os/user"
)

const ()

func init() {
	reg := register.NewRegister(ChanRPC)
	reg.RegisterRpc("SendMail", SendMail)

}

//
func SendMail(args []interface{}) {
	mail := args[0].(*model.Mail)
	player := args[1].(*user.User)
	model.
}
