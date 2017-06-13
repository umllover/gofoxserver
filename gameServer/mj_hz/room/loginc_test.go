package room

import (
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/common"
	"mj/gameServer/conf"
	"mj/gameServer/db"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"net"
	"strconv"
	"testing"

	"mj/gameServer/common/room_base"

	"github.com/lovelly/leaf/chanrpc"
	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

var room = new(Room)

func TestGameStart_1(t *testing.T) {
	room.StartGame()
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
	user := room.GetUserByChairId(0)
	room.OnEventGameConclude(room.ProvideUser, user, GER_USER_LEAVE)
}

func TestDispatchCardData(t *testing.T) {
	room.SendStatus = OutCard_Send
	room.DispatchCardData(0, false)
}

func TestAnalyseCard(t *testing.T) {
	CardInx := make([]int, MAX_INDEX)
	CardInx[room.gameLogic.SwitchToCardData(0x01)] = 3
	CardInx[room.gameLogic.SwitchToCardData(0x02)] = 3
	CardInx[room.gameLogic.SwitchToCardData(0x03)] = 3
	CardInx[room.gameLogic.SwitchToCardData(0x04)] = 3
	CardInx[room.gameLogic.SwitchToCardData(0x05)] = 2
	WeaveItem := make([]*msg.WeaveItem, 0)
	tg := make([]*TagAnalyseItem, 0)
	bret, arr := room.gameLogic.AnalyseCard(CardInx, WeaveItem, 0, tg)
	log.Debug("TestAnalyseCard ret == %v ", bret)
	for _, v := range arr {
		log.Debug("aaa %v", v)
	}
}

func init() {
	conf.Init("E:/gowork/src/mj/gameServer/gameApp/gameServer.json")
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
	var userCnt = 4
	room.RoomBase = room_base.NewRoomBase(userCnt, 1, nil, "aa")
	room.Kind = 389
	room.ServerId = 1
	room.Name = fmt.Sprintf(strconv.Itoa(common.KIND_TYPE_HZMJ)+"_%v", room.GetRoomId())
	room.CloseSig = make(chan bool, 1)
	room.TimeLimit = 90
	room.CountLimit = 90
	room.Source = 1
	room.Password = ""
	room.JoinGamePeopleCount = 2
	room.CreateUser = 3
	room.Response = make([]bool, userCnt)
	room.gameLogic = NewGameLogic()
	room.Owner = 3
	room.BankerUser = INVALID_CHAIR
	room.Record = &msg.G2C_Record{HuCount: make([]int, room.UserCnt), MaCount: make([]int, room.UserCnt), AnGang: make([]int, room.UserCnt), MingGang: make([]int, room.UserCnt), AllScore: make([]int, room.UserCnt), DetailScore: make([][]int, room.UserCnt)}
	room.CardIndex = make([][]int, room.UserCnt)
	room.HeapCardInfo = make([][]int, room.UserCnt) //堆牌信息
	room.HistoryScores = make([]*HistoryScore, room.UserCnt)

	for i, _ := range room.Users {
		room.Users[i] = newTestUser(i)
		room.Users[i].ChairId = i
	}

}

func newTestUser(uid int) *user.User {
	u := new(user.User)
	u.Accountsinfo = &model.Accountsinfo{}
	u.Accountsmember = &model.Accountsmember{}
	u.Gamescorelocker = &model.Gamescorelocker{}
	u.Gamescoreinfo = &model.Gamescoreinfo{}
	u.Userextrainfo = &model.Userextrainfo{}
	u.Userattr = &model.Userattr{}
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
