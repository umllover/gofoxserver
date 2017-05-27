package internal

import (
	"github.com/lovelly/leaf/cluster"
)

var (

)

func handleRpc(id interface{}, f interface{}) {
	cluster.SetRoute(id, ChanRPC)
	skeleton.RegisterChanRPC(id, f)
}

func init() {

}
