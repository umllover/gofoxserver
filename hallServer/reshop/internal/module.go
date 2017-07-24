package internal

import (
	"mj/hallServer/base"

	"mj/hallServer/db/model"

	"fmt"

	"github.com/lovelly/leaf/module"
)

var (
	skeleton   = base.NewSkeleton()
	ChanRPC    = skeleton.ChanRPCServer
	ShopsLives = make(map[int]*model.ShopLive)
)

type Module struct {
	*module.Skeleton
}

func (m *Module) OnInit() {
	m.Skeleton = skeleton
	LoadShopLive()
}

func (m *Module) OnDestroy() {

}

func LoadShopLive() {
	infos, err := model.ShopLiveOp.SelectAll()
	if err != nil {
		fmt.Println("Loading error :%s", err.Error())
	}
	for _, v := range infos {
		ShopsLives[v.Id] = v
	}
}

// 在handler中调用
func GetShopLive(id int) *model.ShopLive {
	return ShopsLives[id]
}
