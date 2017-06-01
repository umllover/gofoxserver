package internal

import (
	"mj/common/msg"
	"github.com/lovelly/leaf/cluster"
	"reflect"
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}, fType int) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.RegisterFromType(id, f, fType)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init(){
	handlerC2S(&msg.C2G_REQUserInfo{}, GerUserInfo)

}


func GerUserInfo(args []interface{}) {

}


