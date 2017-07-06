package internal

import (
	"mj/gameServer/base"

	"time"

	"container/list"

	"mj/common/msg"

	"mj/hallServer/game_list"

	"sort"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

type MachPlayer struct {
	ch      *chanrpc.Server
	EndTime int64
	Uid     int
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
		log.Error("at GetRooms error:%s", err.Error())
		return make(map[int][]*msg.RoomInfo)
	}
	return rooms.(map[int][]*msg.RoomInfo)
}

func (m *MatchModule) Match() {
	now := time.Now().Unix()
	defer m.Skeleton.AfterFunc(2*time.Second, m.Match)
	for kind, v := range m.GetRooms() {
		li := MatchList[kind]
		if li.Len() < 1 {
			continue
		}
		if len(v) > 1 {
			sort.Slice(v, func(i, j int) bool {
				if v[i].CurCnt > v[j].CurCnt {
					return true
				} else if v[i].CurCnt == v[j].CurCnt && v[i].CreateTime < v[j].CreateTime {
					return true
				}
				return false
			})

		}

		for _, r := range v {
			bk := false
			if bk {
				break
			}
			for i := len(r.MachPlayer); i < r.MaxCnt; i++ {
				if li.Len() < 1 {
					bk = true
					break
				}
				v1 := li.Front()
				li.Remove(v1)
				player := v1.Value.(*MachPlayer)
				r.MachPlayer[player.Uid] = struct{}{}
				player.ch.Go("matchResult", true, r)
			}
		}
	}

	for _, li := range MatchList {
		//检测匹配超时
		for e := li.Front(); e != nil; e = e.Next() {
			player := e.Value.(*MachPlayer)
			if player.EndTime < now {
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
