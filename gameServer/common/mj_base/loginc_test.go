package mj_base

import (
	"mj/common/msg"
	"mj/gameServer/common/mj_base"
	"mj/gameServer/common/room_base"
	"mj/gameServer/conf"
	"mj/gameServer/db"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"net"
	"testing"

	"mj/gameServer/RoomMgr"

	"github.com/lovelly/leaf/chanrpc"
	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/module"
)

var (
	room RoomMgr.IRoom
	u1   *user.User
	u2   *user.User
	u3   *user.User
	u4   *user.User
)

func TestGameStart_1(t *testing.T) {

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

	db.InitDB(&conf.DBConfig{})
	base.LoadBaseData()

	u1 := newTestUser(1)
	cfg := &NewMjCtlConfig{
		BaseMgr:  room_base.NewRoomBase(),
		DataMgr:  NewDataMgr(info.RoomId, u.Id, temp.GameName, temp),
		UserMgr:  room_base.NewRoomUserMgr(info.RoomId, info.MaxPlayerCnt, temp),
		LogicMgr: NewBaseLogic(),
		TimerMgr: room_base.NewRoomTimerMgr(),
	}
	r := mj_base.NewMJBase(info, cfg)
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
		data := &msg.C2G_UserSitdown{}
		data.TableID = room.GetRoomId()
		data.ChairID = i
		room.GetChanRPC().Go("Sitdown", data, u)
	}

}

func newTestUser(uid int) *user.User {
	u := new(user.User)
	u.Id = uid
	u.RoomId = 1
	u.Status = 0
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
