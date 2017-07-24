package internal

import (
	"reflect"

	"github.com/lovelly/leaf/gate"

	"mj/common/msg"
	"mj/hallServer/db/model/base"
	"mj/hallServer/user"
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
	handlerC2S(&msg.C2L_TradeGoods{}, ExchangeGoods)
}

//获取商店信息
func GetShopInfo(args []interface{}) {

	//recvMsg := args[0].(*msg.CL2_ReqShopInfo)
	agent := args[1].(gate.Agent)
	retMsg := &msg.L2C_RspShopInfo{}
	for _, v := range base.GoodsCache.All() {
		item := &msg.ShopItem{}
		item.Id = v.GoodsId
		item.Name = v.Name
		shopLive := GetGoodsLive(v.GoodsId)
		item.LeftAmount = shopLive.LeftAmount
		retMsg.Items = append(retMsg.Items, item)
	}

	agent.WriteMsg(retMsg)
}

func ExchangeGoods(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_TradeGoods)
	agent := args[1].(gate.Agent)
	player := agent.UserData().(*user.User)
	retMsg := &msg.L2C_RspTradeShopInfo{}
	goods, _ := base.GoodsCache.Get(recvMsg.ShopID)
	if player.SubCurrency(goods.Rmb) {
		return
	}

	switch goods.GoodsType {
	case 1: //vip
		player.SetVip()
	}

	agent.WriteMsg(retMsg)
}
