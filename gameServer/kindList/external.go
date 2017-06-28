package kindList

import (
	"mj/gameServer/kindList/internal"

	"mj/gameServer/common/room_base"

	"github.com/lovelly/leaf/module"
)

var (
	ChanRPC = internal.ChanRPC
)

func GetModules() []module.Module {
	return internal.GetModules()
}

func GetModByKind(kind int) (room_base.Module, bool) {
	return internal.GetModByKind(kind)
}

func Init() {
	internal.LoadAllModule()
	internal.Clears()
}
