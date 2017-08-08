package internal

import (
	"mj/common/base"
	"mj/hallServer/center"

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
	center.SetOfflineHandler(AddOfflineHandler)
}

func (m *MgrModule) OnDestroy() {

}
