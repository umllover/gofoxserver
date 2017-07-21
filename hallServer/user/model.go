package user

import (
	"mj/common/base"

	"github.com/lovelly/leaf/module"
)

var (
	UserMgr  = new(MgrModule)
	skeleton = base.NewSkeleton()
)

type MgrModule struct {
	*module.Skeleton
}

func (m *MgrModule) OnInit() {
	m.Skeleton = skeleton
}

func (m *MgrModule) OnDestroy() {

}
