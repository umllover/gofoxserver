package game_list

import (
	"mj/hallServer/game_list/internal"
)

var (
	Module  = new(internal.Module)
	ChanRPC = internal.ChanRPC
)

func init() {
	//cluster.AgentChanRPC = ChanRPC
}

func GetSvrByKind(kindId int) (string, int) {
	return internal.GetSvrByKind(kindId)
}

func GetSvrByNodeID(kindId int) string {
	return internal.GetSvrByNodeID(kindId)
}

func SetTest(v bool) {
	internal.Test = v
}
