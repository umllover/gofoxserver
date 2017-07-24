package internal

import (
	"github.com/lovelly/leaf/log"
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
	log.Debug("at server close offline user ")
}
