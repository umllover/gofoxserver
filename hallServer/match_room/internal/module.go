package internal

import (
	"container/list"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/hallServer/base"
	"mj/hallServer/game_list"
	"sort"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

const (
	ResetMatchTime = 45
)

type MachPlayer struct {
	ch      *chanrpc.Server
	EndTime int64
	Uid     int64
}

var (
	skeleton  = base.NewSkeleton()
	ChanRPC   = skeleton.ChanRPCServer
	MatchList = make(map[int]*list.List)

	DefaultMachModule = new(MatchModule)
)

type MatchModule struct {
	*module.Skeleton
	rooms map[int]*msg.RoomInfo
}

func (m *MatchModule) OnInit() {
	m.Skeleton = skeleton
	m.Skeleton.AfterFunc(2*time.Second, m.Match)
	game_list.SetMachRpc(ChanRPC)
}

func (m *MatchModule) OnDestroy() {

}

func (m *MatchModule) GetRoomsByKind(kind int) []*msg.RoomInfo {
	log.Debug("beginc GetRoomsByKind %d", kind)
	rooms, err := game_list.ChanRPC.TimeOutCall1("GetMatchRoomsByKind", 5, kind)
	if err != nil {
		log.Debug("at GetRoomsByKind error:%s", err.Error())
		return []*msg.RoomInfo{}
	}
	return rooms.([]*msg.RoomInfo)
}

func (m *MatchModule) GetRoomByRoomId(RoomId int) *msg.RoomInfo {
	log.Debug("beginc GetRoomByRoomId %d", RoomId)
	rooms, err := game_list.ChanRPC.TimeOutCall1("GetRoomByRoomId", 5, RoomId)
	log.Debug("end GetRoomByRoomId %d", RoomId)
	if err != nil {
		return nil
	}
	return rooms.(*msg.RoomInfo)
}

func (m *MatchModule) Match() {
	now := time.Now().Unix()
	defer m.Skeleton.AfterFunc(2*time.Second, m.Match)
	if len(MatchList) < 1 {
		return
	}

	for kindid, li := range MatchList {
		if li.Len() < 1 {
			continue
		}

		rooms := m.GetRoomsByKind(kindid)
		if len(rooms) < 1 {
			continue
		}

		sort.Slice(rooms, func(i, j int) bool {
			if rooms[i].CurCnt > rooms[j].CurCnt {
				return true
			} else if rooms[i].CurCnt == rooms[j].CurCnt && rooms[i].CreateTime < rooms[j].CreateTime {
				return true
			}
			return false
		})

		for _, r := range rooms {
			bk := false
			if bk {
				break
			}

			if r.MachCnt >= r.MaxPlayerCnt {
				continue
			}

			if r.Status != RoomStatusReady {
				continue
			}

			for i := len(r.MachPlayer); i < r.MaxPlayerCnt; i++ {
				if li.Len() < 1 {
					bk = true
					break
				}

				v1 := li.Front()
				player := v1.Value.(*MachPlayer)
				_, has := r.MachPlayer[player.Uid]
				if !has {
					cnt, err := IncRoomCnt(r.RoomID)
					if err != nil {
						break
					}

					if cnt > r.MaxPlayerCnt {
						log.Debug("at MatchModule roomInfo.MachCnt >= roomInfo.MaxPlayerCnt 222, %v", cnt)
						break
					}

					r.MachPlayer[player.Uid] = time.Now().Unix() + ResetMatchTime
					r.MachCnt = cnt
				}

				li.Remove(v1)
				log.Debug("player %d match ok ", player.Uid)
				player.ch.Go("matchResult", true, r)
			}
		}
	}

	for _, li := range MatchList {
		//检测匹配超时
		for e := li.Front(); e != nil; e = e.Next() {
			player := e.Value.(*MachPlayer)
			if player.EndTime < now {
				log.Debug("player %d match tmie out ", player.Uid)
				li.Remove(e)
				player.ch.Go("matchResult", false, nil)
			}
		}
	}

}

func (m *MatchModule) AddMatchPlayer(kindID int, p *MachPlayer) {
	li := MatchList[kindID]
	if li == nil {
		li = list.New()
		MatchList[kindID] = li
	}

	li.PushBack(p)
}
