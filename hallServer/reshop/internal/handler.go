package internal

import (
	"reflect"

	"mj/common/msg"

	"mj/hallServer/db/model/base"

	"github.com/lovelly/leaf/gate"
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}) {
	ChanRPC.Register(id, f)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handlerC2S(&msg.C2L_ReqShopInfo{}, GetShopInfo)
}

func GetShopInfo(arg []interface{}) {
	agent := arg[1].(gate.Agent)
	retMsg := &msg.L2C_RspShopInfo{}

	for _, v := range base.ShopCache.All() {
		item := &msg.ShopItem{}
		item.Id = v.Id
		item.Name = v.Name
		item.Price = v.Price
		ShopsLives := GetShopLive(v.Id)
		item.LeftAmount = ShopsLives.LeftAmount
		retMsg.Items = append(retMsg.Items, item)
	}

	agent.WriteMsg(retMsg)
}
