package game_list

import (
	"mj/hallServer/game_list/internal"

	"github.com/lovelly/leaf/chanrpc"
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

func SetMachRpc(rpc *chanrpc.Server) {
	internal.MatchRpc = rpc
}
