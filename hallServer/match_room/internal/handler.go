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
	reg.RegisterC2S(&msg.C2L_SearchServerTable{}, SearchTable)

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
func SearchTable(args []interface{}) {
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
		log.Error("at SearchTable not foud room, %v", recvMsg)
		retcode = ErrNoFoudRoom
		return
	}
	_, has := roomInfo.MachPlayer[player.Id]
	if !has {
		if roomInfo.MachCnt >= roomInfo.MaxPlayerCnt {
			log.Debug("at SearchTable roomInfo.MachCnt >= roomInfo.MaxPlayerCnt, %v", recvMsg)
			retcode = ErrRoomFull
			return
		}

		if roomInfo.Status != RoomStatusReady {
			log.Debug("at SearchTable room is start , %v", recvMsg)
			retcode = ErrRoomIsPlaying
			return
		}

		cnt, err := IncRoomCnt(roomInfo.RoomID)
		if err != nil {
			log.Debug("Error === %s ", err.Error())
			retcode = ErrRoomFull
			return
		}

		if cnt > roomInfo.MaxPlayerCnt {
			log.Debug("at SearchTable roomInfo.MachCnt >= roomInfo.MaxPlayerCnt 222, %v", recvMsg)
			retcode = ErrRoomFull
			return
		}

		roomInfo.MachCnt = cnt
		roomInfo.MachPlayer[player.Id] = time.Now().Unix() + ResetMatchTime
	}

	agent.ChanRPC().Go("SearchTableResult", roomInfo)
	return
}

func delMatchPlayer(args []interface{}) {
	log.Debug("at del match player ")
	uid := args[0].(int64)
	roomInfo := args[1].(*msg.RoomInfo)
	if roomInfo.MachCnt > roomInfo.MaxPlayerCnt {
		roomInfo.MachCnt = roomInfo.MaxPlayerCnt
	}
	roomInfo.MachCnt -= 1
	delete(roomInfo.MachPlayer, uid)
	UpRoomCnt(roomInfo.RoomID, roomInfo.MachCnt)
}
