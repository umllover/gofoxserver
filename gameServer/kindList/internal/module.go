package internal

import (
	"github.com/lovelly/leaf/module"
	"mj/gameServer/base"
	"mj/gameServer/common"
	"mj/gameServer/conf"
	"mj/gameServer/hzmj"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	modules = make(map[int]common.Module) //key kind
	KModule  = new(Module)
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	loadAllModule()
}

func (m *Module) OnDestroy() {

}

func loadAllModule(){
	if HasKind(common.KIND_TYPE_HZMJ) {
		modules[common.KIND_TYPE_HZMJ] = hzmj.Module
	}
}

func GetModules()[]module.Module {
	ret := make([]module.Module, 0)
	for _, v := range modules {
		ret = append(ret, v)
	}
	ret = append(ret, KModule)
	return ret
}

func HasKind(kind int) bool{
	_, ok := conf.ValidKind[kind]
	return ok
}