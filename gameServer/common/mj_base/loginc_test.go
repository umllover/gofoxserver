package mj_base

import (
	. "mj/common/cost"
	"mj/gameServer/common/room_base"
	"mj/gameServer/conf"
	"mj/gameServer/db"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"net"
	"testing"

	"fmt"

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
)

func TestGameStart_1(t *testing.T) {
	room.UserReady([]interface{}{nil, u1})
	room.DataMgr.SetUserCard(0, []int{
		0x1, 0x1, 0x1,
		0x2, 0x2, 0x2,
		0x3, 0x3, 0x3,
		0x4, 0x4,
		0x5, 0x5,
	})
}

func TestOutCard(t *testing.T) {
	ret := room.DataMgr.EstimateUserRespond(1, 0x4, EstimatKind_OutCard)
	fmt.Println("at EstimateUserRespond ret :", ret)
	room.OutCard([]interface{}{u1, 1})
}

func TestGameConclude(t *testing.T) {
	room.UserOperateCard([]interface{}{u1, 1, []int{1}})
}

func TestDispatchCardData(t *testing.T) {

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
	lconf.PendingWriteNum = conf.Server.PendingWriteNum
	lconf.HeartBeatInterval = conf.HeartBeatInterval
	leaf.InitLog()

	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()

	temp, ok := base.GameServiceOptionCache.Get(389, 1)
	if !ok {
		return
	}

	info := &model.CreateRoomInfo{
		RoomId:       777777,
		MaxPlayerCnt: 4,
		KindId:       389,
		ServiceId:    1,
	}

	base := room_base.NewRoomBase()

	userg := room_base.NewRoomUserMgr(info.RoomId, info.MaxPlayerCnt, temp)

	u1 = newTestUser(1)
	u1.ChairId = 0
	userg.Users[0] = u1
	r := NewMJBase(info)
	datag := NewDataMgr(info.RoomId, u1.Id, IDX_HZMJ, "", temp, r)
	cfg := &NewMjCtlConfig{
		BaseMgr:  base,
		DataMgr:  datag,
		UserMgr:  userg,
		LogicMgr: NewBaseLogic(IDX_HZMJ),
		TimerMgr: room_base.NewRoomTimerMgr(info.Num, temp),
	}
	r.Init(cfg)
	room = r
	var userCnt = 4

	for i := 1; i < userCnt; i++ {
		u := newTestUser(i + 1)
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

func newTestUser(uid int) *user.User {
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
