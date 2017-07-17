package room

import (
	"encoding/json"
	. "mj/common/cost"
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common/mj/mj_base"
	"mj/gameServer/common/room_base"
	"mj/gameServer/conf"
	"mj/gameServer/db"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"net"
	"os"
	"sync"
	"testing"

	"mj/common/msg/mj_zp_msg"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

var (
	room *xs_entry //Mj_base
	u1   *user.User
	u2   *user.User
	u3   *user.User
	u4   *user.User
)

var Wg sync.WaitGroup

func TestGameStart_1(t *testing.T) {
	room.UserReady([]interface{}{nil, u1})

}

func TestOutCard(t *testing.T) {
	time.Sleep(3 * time.Second)
	a := []int{}
	room.DataMgr.CalHuPaiScore(a)
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
	lconf.PendingWriteNum = conf.Server.PendingWriteNum
	lconf.HeartBeatInterval = conf.HeartBeatInterval
	InitLog()

	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()

	temp, ok := base.GameServiceOptionCache.Get(389, 1)
	if !ok {
		return
	}
	temp.OutCardTime = 2

	info := &model.CreateRoomInfo{
		RoomId:       777777,
		MaxPlayerCnt: 4,
		KindId:       389,
		ServiceId:    1,
	}

	//游戏配置
	type gameCfg struct {
		ZhuaHua    int
		WithZiCard bool
		ScoreType  int
	}
	setCfg := map[string]interface{}{
		"ZhuaHua":    0,
		"WithZiCard": false,
		"ScoreType":  33,
	}
	myCfg, cfgOk := json.Marshal(setCfg)
	if cfgOk != nil {
		log.Error("测试错误，退出程序")
		os.Exit(0)
	}
	info.OtherInfo = string(myCfg)

	base := room_base.NewRoomBase()

	userg := room_base.NewRoomUserMgr(info.RoomId, info.MaxPlayerCnt, temp)

	u1 = newTestUser(1)
	u1.ChairId = 0
	userg.Users[0] = u1
	r := NewXSEntry(info)
	datag := NewXSDataMgr(info.RoomId, u1.Id, mj_base.IDX_ZPMJ, "", temp, r)
	if datag == nil {
		log.Error("测试错误，退出程序")
		os.Exit(0)
	}
	cfg := &mj_base.NewMjCtlConfig{
		BaseMgr:  base,
		DataMgr:  datag,
		UserMgr:  userg,
		LogicMgr: NewXSlogic(mj_base.IDX_ZPMJ),
		TimerMgr: room_base.NewRoomTimerMgr(info.Num, temp),
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
func InitLog() {
	logger, err := log.New(conf.Server.LogLevel, "", conf.LogFlag)
	if err != nil {
		panic(err)
	}
	log.Export(logger)
}