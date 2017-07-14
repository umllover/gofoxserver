package internal

import (
	"mj/common/register"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	reg.RegisterRpc("offLineHanle", offLineHanle)
}

func offLineHanle(args []interface{}) {

}
