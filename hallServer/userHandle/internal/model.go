package internal

import (
	"mj/common/base"

	"github.com/lovelly/leaf/module"
)

var (
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
