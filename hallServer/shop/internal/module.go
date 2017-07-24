package internal

import (
	"mj/hallServer/base"
	"mj/hallServer/db/model"

	base2 "mj/gameServer/db/model/base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

var (
	skeleton   = base.NewSkeleton()
	ChanRPC    = skeleton.ChanRPCServer
	GoodsLives = make(map[int]*model.GoodsLive)
	GoodsType  = make(map[int]*base2.Goods)
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
	infos, err := model.GoodsLiveOp.SelectAll()
	if err != nil {
		log.Fatal("LoadShopLive error:%s ", err.Error())
		return
	}

	for _, v := range infos {
		GoodsLives[v.Id] = v
	}
}

func GetGoodsLive(id int) *model.GoodsLive {
	return GoodsLives[id]
}

func GetGoodsType(id int) *base2.Goods {
	return GoodsType[id]
}
