package gameList

import (
	"mj/hallServer/gameList/internal"
	"github.com/lovelly/leaf/cluster"
)

var (
	Module  = new(internal.Module)
	ChanRPC = internal.ChanRPC
)

func init(){
	cluster.AgentChanRPC = ChanRPC
}