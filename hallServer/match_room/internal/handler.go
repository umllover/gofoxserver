package internal

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/register"

	"github.com/lovelly/leaf/log"

	"mj/hallServer/user"

	"mj/hallServer/common"
	"time"

	"github.com/lovelly/leaf/gate"
)

func init() {
	reg := register.NewRegister(ChanRPC)
	reg.RegisterC2S(&msg.C2L_QuickMatch{}, QuickMatch)
	reg.RegisterC2S(&msg.C2L_SearchServerTable{}, SrarchTable)

	reg.RegisterRpc("delMatchPlayer", delMatchPlayer)
}

func QuickMatch(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_QuickMatch)
	agent := args[1].(gate.Agent)
	player := agent.UserData().(*user.User)

	limitTime := common.GetGlobalVarInt(MATCH_TIMEOUT)
	matchPlayer := &MachPlayer{Uid: player.Id, ch: agent.ChanRPC(), EndTime: time.Now().Unix() + int64(limitTime)}
	DefaultMachModule.AddMatchPlayer(recvMsg.KindID, matchPlayer)
	agent.WriteMsg(&msg.L2C_QuickMatchOk{MatchTime: limitTime})
}

//玩家请求查找房间
func SrarchTable(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_SearchServerTable)
	agent := args[1].(gate.Agent)
	player := agent.UserData().(*user.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			agent.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	roomInfo := DefaultMachModule.GetRoomByRoomId(recvMsg.TableID)
	if roomInfo == nil {
		log.Error("at SrarchTable not foud room, %v", recvMsg)
		retcode = ErrNoFoudRoom
		return
	}

	//_, has := roomInfo.MachPlayer[player.Id]
	//if len(roomInfo.MachPlayer) >= roomInfo.MaxPlayerCnt && !has {
	//	retcode = ErrRoomFull
	//	return
	//}

	roomInfo.MachPlayer[player.Id] = struct{}{}

	agent.ChanRPC().Go("SrarchTableResult", roomInfo)
	return
}

func delMatchPlayer(args []interface{}) {
	uid := args[0].(int64)
	roomInfo := args[1].(*msg.RoomInfo)
	delete(roomInfo.MachPlayer, uid)
}
