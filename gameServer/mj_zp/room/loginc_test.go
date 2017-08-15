package room

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/common/room_base"
	"mj/gameServer/conf"
	"mj/gameServer/db"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"mj/common/msg/mj_zp_msg"

	"fmt"

	"github.com/lovelly/leaf/chanrpc"
	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

var (
	room *ZP_base //Mj_base
	u1   *user.User
	u2   *user.User
	u3   *user.User
	u4   *user.User
)

var Wg sync.WaitGroup

func TestGameStart_1(t *testing.T) {
	room.UserReady([]interface{}{nil, u1})
}

func TestIndex(t *testing.T) {
	fmt.Println("11111111111111111111 ")
	var cards = []int{38, 37, 7, 25, 1, 18, 21, 20, 4, 20, 25, 70, 25, 5, 20, 5, 39, 22, 24, 5, 18, 69, 40, 2, 6, 9, 19, 36, 17, 2, 33, 39, 38, 3, 18, 19,
		67, 3, 36, 68, 6, 6, 2,
		22, 17, 35, 3, 36, 4, 35, 35, 17, 19, 66, 1, 4, 71, 39, 5, 6, 9, 20, 2, 41, 21, 23, 8, 33, 9, 1, 22, 33, 24, 41, 24, 1, 8, 3, 7, 35, 8, 23, 9, 37,
		4, 7, 21, 37, 7, 33, 21, 34, 40, 25, 40, 34, 40, 18, 65, 17, 23, 8, 24, 41, 41, 34, 38, 37, 38, 72, 19, 36, 23, 34, 22, 39}
	logic := room.LogicMgr.(*ZP_Logic)
	data := room.DataMgr.(*ZP_RoomData)
	for _, card := range cards {
		idx := logic.SwitchToIdx(card)
		fmt.Println(idx)
		data.CardIndex[0][idx]++
	}
}

func TestZP_RoomData_StartDispatchCard(t *testing.T) {

	//data := room.DataMgr.(*ZP_RoomData)
	//data.RepertoryCard = make([]int, 144)
	//data.StartDispatchCard()

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

func TestOutCard(t *testing.T) {
	Wg.Add(1)
	time.Sleep(3 * time.Second)
	//a := []int{}
	//room.DataMgr.CalHuPaiScore(a)
	data := &mj_zp_msg.C2G_ZPMJ_OperateCard{}
	data.OperateCard = append(data.OperateCard, 5)
	data.OperateCard = append(data.OperateCard, 0)
	data.OperateCard = append(data.OperateCard, 5)
	data.OperateCode = 64
	if room != nil {
		room.GetChanRPC().Go("OperateCard", u1, data.OperateCode, data.OperateCard)
	}

	Wg.Wait()
}

//func TestGameLogic_OutCard(t *testing.T) {
//	user := room.GetUserByChairId(0)
//	if user == nil {
//		t.Error("not foud t")
//	}
//
//	var cardidx int
//	var cnt int
//	for cardidx, cnt = range room.CardIndex[0] {
//		if cnt > 0 {
//			break
//		}
//	}
//
//	card := room.gameLogic.SwitchToCardData(int(cardidx))
//	dt := &msg.C2G_HZMJ_HZOutCard{CardData: card}
//	room.OutCard([]interface{}{dt, user})
//}
//
//func TestRoomUserOperateCard(t *testing.T) {
//	user := room.GetUserByChairId(0)
//	if user == nil {
//		t.Error("not foud t")
//	}
//
//	var cardidx int
//	var cnt int
//	for cardidx, cnt = range room.CardIndex[0] {
//		if cnt > 0 {
//			break
//		}
//	}
//
//	card := room.gameLogic.SwitchToCardData(int(cardidx))
//	dt := &msg.C2G_HZMJ_OperateCard{OperateCard: []int{card, card, card}, OperateCode: WIK_PENG}
//	room.UserOperateCard([]interface{}{dt, user})
//}

func TestGameConclude(t *testing.T) {

}

func TestDispatchCardData(t *testing.T) {

}

func TestAnalyseCard(t *testing.T) {

}

func init() {
	Wg.Add(1)
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
	InitLog()

	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()

	temp, ok := base.GameServiceOptionCache.Get(389, 1)
	if !ok {
		return
	}
	temp.OutCardTime = 2

	info := &msg.L2G_CreatorRoom{
		RoomID:       777777,
		MaxPlayerCnt: 4,
		KindId:       391,
		ServiceId:    1,
		PlayCnt:      8,
	}

	//游戏配置
	type gameCfg struct {
		ZhuaHua    int
		WithZiCard bool
		ScoreType  int
	}
	setCfg := map[string]interface{}{
		"zhuaHua": float64(16),
		"wanFa":   false,
		"suanFen": float64(1),
		"chaHua":  false,
	}

	info.OtherInfo = setCfg

	base := room_base.NewRoomBase()

	userg := room_base.NewRoomUserMgr(info, temp)

	u1 = newTestUser(1)
	u1.ChairId = 0
	userg.Users[0] = u1
	r := NewMJBase(info)
	datag := NewZPDataMgr(info, u1.Id, mj_base.IDX_ZPMJ, "", temp, r)
	if datag == nil {
		log.Error("测试错误，退出程序")
		os.Exit(0)
	}
	cfg := &mj_base.NewMjCtlConfig{
		BaseMgr:  base,
		DataMgr:  datag,
		UserMgr:  userg,
		LogicMgr: NewBaseLogic(mj_base.IDX_ZPMJ),
		TimerMgr: room_base.NewRoomTimerMgr(info.PlayCnt, temp, base.GetSkeleton()),
	}
	r.Init(cfg)
	RegisterHandler(r)
	RoomMgr.AddRoom(r)
	room = r

	var userCnt = 4

	for i := 1; i < userCnt; i++ {
		u := newTestUser(int64(i + 1))
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
	u.Agent = NewAgent()

	return u
}

type TestUser struct {
	*user.User
}

func (t *TestUser) WriteMsg(msg interface{}) {

}

func NewAgent() *TAgent {
	a := new(TAgent)
	a.Ch = chanrpc.NewServer(100000)
	go func() {
		for v := range a.Ch.ChanCall {
			fmt.Println(v)
		}
	}()
	return a
}

type TAgent struct {
	Ch *chanrpc.Server
}

func (t *TAgent) WriteMsg(msg interface{})     {}
func (t *TAgent) Destroy()                     {}
func (t *TAgent) LocalAddr() net.Addr          { return nil }
func (t *TAgent) Close()                       {}
func (t *TAgent) RemoteAddr() net.Addr         { return nil }
func (t *TAgent) UserData() interface{}        { return nil }
func (t *TAgent) SetUserData(data interface{}) {}
func (t *TAgent) Skeleton() *module.Skeleton   { return nil }
func (t *TAgent) ChanRPC() *chanrpc.Server     { return t.Ch }
func (t *TAgent) SetReason(int)                {}

func InitLog() {
	logger, err := log.New(conf.Server.LogLevel, "", conf.LogFlag)
	if err != nil {
		panic(err)
	}
	log.Export(logger)
}
