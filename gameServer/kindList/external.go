package kindList

import (
	"mj/gameServer/kindList/internal"

	"github.com/lovelly/leaf/module"
)

var (
	ChanRPC = internal.ChanRPC
)

func GetModules() []module.Module {
	return internal.GetModules()
}

func Init() {
	internal.LoadAllModule()
}
