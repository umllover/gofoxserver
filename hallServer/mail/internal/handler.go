package internal

import (
	"mj/common/register"
)

const ()

func init() {
	reg := register.NewRegister(ChanRPC)
	reg.RegisterRpc("SendMail", SendMail)

}

//
func SendMail(args []interface{}) {
	//mail := args[0].(*model.Mail)
	//player := args[1].(*user.User)
	//mailID, err := model.MailOp.Insert(mail)
	//if err != nil {
	//
	//}
}
