package internal

import (
	"mj/gameServer/base"
	"mj/gameServer/common"
	"mj/gameServer/common/room_base"
	"mj/gameServer/conf"
	"mj/gameServer/db/model"
	"mj/gameServer/mj_hz"
	"mj/gameServer/mj_zp"
	"mj/gameServer/pk_ddz"
	"mj/gameServer/pk_nn_tb"
	"mj/gameServer/pk_sss"

	"github.com/lovelly/leaf/module"
	//"mj/gameServer/pk_ddz_bak"
)

var (
	skeleton = base.NewSkeleton()
	ChanRPC  = skeleton.ChanRPCServer
	modules  = make(map[int]room_base.Module) //key kind
	KModule  = new(Module)
	Kinds    = map[int]room_base.Module{ // Register here
		common.KIND_TYPE_HZMJ: hzmj.Module,
		common.KIND_TYPE_ZPMJ: zpmj.Module,
		common.KIND_TYPE_TBNN: pk_nn_tb.Module,
		common.KIND_TYPE_DDZ:  pk_ddz.Module,
		common.KIND_TYPE_SSS:  pk_sss.Module,
	}
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton

}

func (m *Module) OnDestroy() {
	Clears()
}

func LoadAllModule() {
	for kind, m := range Kinds {
		if HasKind(kind) && m != nil {
			AddMoudle(kind, m)
		}
	}
}

func AddMoudle(kindID int, m room_base.Module) {
	modules[kindID] = m
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

func Clears() {
	ClearRoomId()
}

func ClearRoomId() {
	model.RoomIdOp.DeleteByMap(map[string]interface{}{
		"node_id": conf.Server.NodeId,
	})
}
