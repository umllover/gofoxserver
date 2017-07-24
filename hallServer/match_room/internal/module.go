package internal

import (
	"container/list"
	"mj/common/msg"
	"mj/hallServer/base"
	"mj/hallServer/game_list"
	"sort"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
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
}

func (m *MatchModule) OnDestroy() {

}

func (m *MatchModule) GetRooms() map[int][]*msg.RoomInfo {
	rooms, err := game_list.ChanRPC.TimeOutCall1("GetMatchRooms", 5)
	if err != nil {
		return make(map[int][]*msg.RoomInfo)
	}
	return rooms.(map[int][]*msg.RoomInfo)
}

func (m *MatchModule) GetRoomsByKind(kind int) []*msg.RoomInfo {
	log.Debug("beginc GetRoomsByKind %d", kind)
	rooms, err := game_list.ChanRPC.TimeOutCall1("GetMatchRoomsByKind", 5, kind)
	log.Debug("end GetRoomsByKind %d， rooms:%v", kind, rooms)
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
			for i := len(r.MachPlayer); i < r.MaxPlayerCnt; i++ {
				if li.Len() < 1 {
					bk = true
					break
				}
				v1 := li.Front()
				li.Remove(v1)
				player := v1.Value.(*MachPlayer)
				r.MachPlayer[player.Uid] = struct{}{}
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
