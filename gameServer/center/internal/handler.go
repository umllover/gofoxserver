package internal

import (

	"github.com/lovelly/leaf/cluster"
)


func handleRpc(id interface{}, f interface{}, fType int) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.RegisterFromType(id, f, fType)
}

func init() {


}
