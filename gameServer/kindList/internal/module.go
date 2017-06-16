package internal

import (
	"mj/gameServer/base"
	"mj/gameServer/common"
	"mj/gameServer/conf"
	"mj/gameServer/db/model"
	"mj/gameServer/mj_hz"
	"mj/gameServer/mj_xs"

	"mj/gameServer/common/room_base"

	"github.com/lovelly/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	modules  = make(map[int]room_base.Module) //key kind
	KModule  = new(Module)
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton

}

func (m *Module) OnDestroy() {
	ClearLoocker()
}

func LoadAllModule() {
	if HasKind(common.KIND_TYPE_HZMJ) {
		modules[common.KIND_TYPE_HZMJ] = hzmj.Module
	}

	if HasKind(common.KIND_TYPE_XSMJ) {
		modules[common.KIND_TYPE_XSMJ] = mj_xs.Module
	}
}

func GetModules() []module.Module {
	ret := make([]module.Module, 0)
	for _, v := range modules {
		ret = append(ret, v)
	}
	ret = append(ret, KModule)
	return ret
}

func HasKind(kind int) bool {
	_, ok := conf.ValidKind[kind]
	return ok
}

func GetModByKind(kind int) (room_base.Module, bool) {
	mod, ok := modules[kind]
	return mod, ok
}

func ClearLoocker() {
	model.GamescorelockerOp.DeleteByMap(map[string]interface{}{
		"NodeID": conf.Server.NodeId,
	})
}
