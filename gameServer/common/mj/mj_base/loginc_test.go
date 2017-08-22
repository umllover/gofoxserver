package mj_base

import (
	. "mj/common/cost"
	"mj/gameServer/common/room_base"
	"mj/gameServer/conf"
	"mj/gameServer/db"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"net"
	"testing"

	"fmt"

	"sync"

	"mj/common/msg"

	"encoding/json"

	"github.com/lovelly/leaf"
	"github.com/lovelly/leaf/chanrpc"
	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/module"
)

var (
	room *Mj_base
	u1   *user.User
	u2   *user.User
	u3   *user.User
	u4   *user.User
	wg   sync.WaitGroup
)

type RpcRequest struct {
	Method string        `json:"func_name"`
	Params []interface{} `json:"params"`
}

func TestGameStart_1(t *testing.T) {

	r := &RpcRequest{}
	r.Method = "GMNotice"
	r.Params = []interface{}{1, 2, "sssss"}
	b, err := json.Marshal(r)
	fmt.Println(string(b), err)
	room.UserReady([]interface{}{nil, u1})
}

func TestGameAddPlayCnt(t *testing.T) {
	room.GetChanRPC().Register("RenewalFeesSetInfo", room.RenewalFeesSetInfo)
	_, err := room.GetChanRPC().Call1("RenewalFeesSetInfo", 8)
	if err != nil {
		fmt.Println("AAAAAAAAAAAAA ", err.Error())
		return
	}
}

func TestOutCard(t *testing.T) {
	//room.GetChanRPC().Go("OutCard", user, 5)
	//ret := room.DataMgr.EstimateUserRespond(1, 0x4, EstimatKind_OutCard)
	//log.Debug("at EstimateUserRespond ret :%v", ret)
	//room.OutCard([]interface{}{u1, 1})
}

func TestMj_base_DissumeRoom(t *testing.T) {
	//base := Mj_base{}
	//args := *user.User{}
	//base.DissumeRoom()
}

func TestBaseLogic_ReplaceCard(t *testing.T) {
	//m := GetCardByIdx(0)
	//log.Debug("库存的牌%v", m)
	//TmpRepertoryCard := []int{1, 1, 3, 17, 25, 24}
	//log.Debug("TmpRepertoryCardAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	////tempCard := make([]int, len(m))
	//
	////room.LogicMgr.RandCardList(tempCard, m)
	//
	//log.Debug("删除前 %d, %v", len(m), m)
	//for _, v := range TmpRepertoryCard {
	//	for idx, v1 := range m {
	//		if v == v1 {
	//			m = utils.IntSliceDelete(m, idx)
	//			break
	//		}
	//	}
	//}
	//log.Debug("删除后%d  %v", len(m), m)
}

func TestBaseLogic_AnalyseCard(t *testing.T) {
	fmt.Println("===========================================")
	lg := room.LogicMgr.(*BaseLogic)
	hzIndex := lg.SwitchToCardIndex(0x35)
	cbCardIndexTemp := make([]int, lg.GetCfg().MaxIdx)
	WraveItem := []*msg.WeaveItem{}

	/*cbCardIndexTemp[0x3] = 1
	cbCardIndexTemp[0x4] = 1
	cbCardIndexTemp[0x5] = 1
	cbCardIndexTemp[0x6] = 1
	cbCardIndexTemp[0x7] = 1
	cbCardIndexTemp[0x8] = 1*/
	cbCardIndexTemp[0x3] = 3
	cbCardIndexTemp[0x6] = 3
	cbCardIndexTemp[0x18] = 1
	cbCardIndexTemp[hzIndex] = 1
	//lg.SetMagicIndex(hzIndex)

	wrave := &msg.WeaveItem{}
	wrave.PublicCard = false
	wrave.Param = 0
	wrave.ActionMask = 0
	wrave.CenterCard = 9
	wrave.CardData = []int{9, 9, 9, 0}
	wrave.WeaveKind = 8
	wrave.ProvideUser = 1
	WraveItem = append(WraveItem, wrave)

	wrave2 := &msg.WeaveItem{}
	wrave2.PublicCard = false
	wrave2.Param = 0
	wrave2.ActionMask = 0
	wrave2.CenterCard = 7
	wrave2.CardData = []int{7, 7, 7, 0}
	wrave2.WeaveKind = 8
	wrave2.ProvideUser = 1
	WraveItem = append(WraveItem, wrave2)

	hu, cards := lg.AnalyseCard(cbCardIndexTemp, WraveItem)
	//hu, cards := lg.AnalyseCard(cbCardIndexTemp, []*msg.WeaveItem{})
	fmt.Println(hu, cards)
	fmt.Println("===========================================")
}

func TestRandRandCard(t *testing.T) {
	/*for i := 0; i < 100000; i++ {
		m := make(map[int]int)
		newCard := make([]int, len(cards[1]))
		room.LogicMgr.RandCardList(newCard, cards[1])
		for _, v := range newCard {
			m[v]++
			if v <= 0x37 {
				if m[v] > 4 {
					log.Debug("cards  ==== card :%d  ## :%v", v, newCard)
				}
			}
			if v > 0x37 {
				if m[v] > 1 {
					log.Debug("cards  ==== card :%d  ## :%v", v, newCard)
				}
			}
		}
	}*/
}

func TestGameConclude(t *testing.T) {
	//room.UserOperateCard([]interface{}{u1, 1, []int{1}})
}

func TestStartDispatchCard(t *testing.T) {
	fmt.Println("===========================================")
	fmt.Println("===========================================")
}

func TestAnalyseCard(t *testing.T) {

}

func init() {
	conf.Init("./gameServer/gameApp/gameServer.json")
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ServerName = conf.ServerName()
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath
	lconf.ListenAddr = conf.Server.ListenAddr
	lconf.ConnAddrs = conf.Server.ConnAddrs
	conf.Test = true
	lconf.PendingWriteNum = conf.Server.PendingWriteNum
	leaf.InitLog()

	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()

	temp, ok := base.GameServiceOptionCache.Get(389, 1)
	if !ok {
		return
	}

	info := &msg.L2G_CreatorRoom{
		RoomID:       777777,
		MaxPlayerCnt: 4,
		KindId:       389,
		ServiceId:    1,
	}

	base := room_base.NewRoomBase()

	userg := room_base.NewRoomUserMgr(info, temp)

	u1 = newTestUser(1)
	u1.ChairId = 0
	userg.Users[0] = u1
	r := NewMJBase(info.KindId, info.ServiceId)
	datag := NewDataMgr(info.RoomID, u1.Id, IDX_HZMJ, "", temp, r, info)
	cfg := &NewMjCtlConfig{
		BaseMgr:  base,
		DataMgr:  datag,
		UserMgr:  userg,
		LogicMgr: NewBaseLogic(IDX_HZMJ),
		TimerMgr: room_base.NewRoomTimerMgr(info.PlayCnt, temp, base.GetSkeleton()),
	}
	r.Init(cfg)
	room = r
	var userCnt = 4
	for i := 1; i < userCnt; i++ {
		u := newTestUser((int64)(i + 1))
		if i == 1 {
			u2 = u
		} else if 1 == 2 {
			u3 = u
		} else if i == 3 {
			u4 = u
		}
		userg.Users[i] = u
		u.ChairId = i
	}
}

func newTestUser(uid int64) *user.User {
	u := new(user.User)
	u.Id = uid
	u.RoomId = 1
	if uid != 1 {
		u.Status = US_READY
	}

	u.ChairId = 0
	u.Agent = new(TAgent)
	return u
}

type TestUser struct {
	*user.User
}

func (t *TestUser) WriteMsg(msg interface{}) {

}

type TAgent struct {
}

func (t *TAgent) WriteMsg(msg interface{})     {}
func (t *TAgent) Destroy()                     {}
func (t *TAgent) LocalAddr() net.Addr          { return nil }
func (t *TAgent) Close()                       {}
func (t *TAgent) RemoteAddr() net.Addr         { return nil }
func (t *TAgent) UserData() interface{}        { return nil }
func (t *TAgent) SetUserData(data interface{}) {}
func (t *TAgent) Skeleton() *module.Skeleton   { return nil }
func (t *TAgent) ChanRPC() *chanrpc.Server     { return nil }
func (t *TAgent) SetReason(int)                {}
