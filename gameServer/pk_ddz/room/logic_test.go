package room

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

	"sync"

	"os"

	"encoding/json"

	"mj/gameServer/common/pk/pk_base"

	"mj/common/msg/pk_ddz_msg"

	"github.com/lovelly/leaf/chanrpc"
	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
)

var (
	room   *DDZ_Entry //Pk_base
	ddzMrg *ddz_data_mgr
	u1     *user.User
	u2     *user.User
	u3     *user.User
)

var Wg sync.WaitGroup

func TestGameStart_1(t *testing.T) {
	room.UserReady([]interface{}{nil, u1})

}

func TestShowCard(t *testing.T) {
	args := []interface{}{nil, u2}
	room.ShowCard(args)
}

func TestCallScore(t *testing.T) {
	log.Debug("测试叫分")
	data := &pk_ddz_msg.C2G_DDZ_CallScore{
		CallScore: 1,
	}

	args := []interface{}{data, u1}
	room.CallScore(args)

	data.CallScore = 2
	args = []interface{}{data, u2}
	room.CallScore(args)

	data.CallScore = 3
	args = []interface{}{data, u3}
	room.CallScore(args)

	//Wg.Wait()
}

func TestTrustee(t *testing.T) {
	log.Debug("测试托管")
	data := &pk_ddz_msg.C2G_DDZ_TRUSTEE{
		Trustee: true,
	}

	args := []interface{}{data, u1}
	room.CTrustee(args)

	args = []interface{}{data, u2}
	room.CTrustee(args)

	args = []interface{}{data, u3}
	room.CTrustee(args)
}

func TestOutCard(t *testing.T) {
	log.Debug("测试出牌")
	data := &pk_ddz_msg.C2G_DDZ_OutCard{
		CardData: []int{ddzMrg.HandCardData[2][len(ddzMrg.HandCardData[2])-1]},
	}

	args := []interface{}{data, u3}
	room.OutCard(args)

	//reader := bufio.NewReader(os.Stdin)
	//line, _ := reader.ReadString('a')
	////line, _ = reader.ReadString('\n')
	//log.Debug("1111%s", line)

	//cardData, _, _ := reader.ReadLine()
	//log.Debug("sfd%v", cardData)

	//fmt.Print("请输入Í∑Í要打的牌")
	//reader := bufio.NewReader(os.Stdin)
	//
	//cardData, _, _ := reader.ReadLine()
	//fmt.Printf("dfsdfsd%s", cardData)
	//Wg.Wait()
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

//func TestCardType(t *testing.T) {
//	lg := new(ddz_logic)
//
//	var c0 []int
//	log.Debug("空牌-%d", lg.GetCardType(c0))
//	c1 := [...]int{0x01}
//	log.Debug("单牌%d", lg.GetCardType(c1[:]))
//	c2 := [...]int{0x03, 0x33}
//	log.Debug("对子-%d", lg.GetCardType(c2[:]))
//	c21 := [...]int{0x03, 0x31}
//	log.Debug("无效两根-%d", lg.GetCardType(c21[:]))
//	c3 := [...]int{0x03, 0x23, 0x33}
//	log.Debug("三根-%d", lg.GetCardType(c3[:]))
//	c31 := [...]int{0x04, 0x34, 0x24, 0x08}
//	log.Debug("三代一-%d", lg.GetCardType(c31[:]))
//	c32 := [...]int{0x04, 0x34, 0x24, 0x08, 0x18}
//	log.Debug("三代二-%d", lg.GetCardType(c32[:]))
//	c5 := [...]int{0x03, 0x34, 0x25, 0x06, 0x17}
//	log.Debug("顺子%d", lg.GetCardType(c5[:]))
//	c51 := [...]int{0x03, 0x34, 0x25, 0x06, 0x17, 0x02}
//	log.Debug("带2的顺子-%d", lg.GetCardType(c51[:]))
//	c4 := [...]int{0x03, 0x33, 0x24, 0x04}
//	log.Debug("两个连续对子%d", lg.GetCardType(c4[:]))
//	c6 := [...]int{0x03, 0x33, 0x22, 0x02, 0x14, 0x04}
//	log.Debug("带2连对%d", lg.GetCardType(c6[:]))
//	c61 := [...]int{0x03, 0x33, 0x25, 0x05, 0x14, 0x04}
//	log.Debug("连对%d", lg.GetCardType(c61[:]))
//	c62 := [...]int{0x03, 0x33, 0x23, 0x02, 0x12, 0x32}
//	log.Debug("带2三顺子%d", lg.GetCardType(c62[:]))
//	c63 := [...]int{0x03, 0x33, 0x23, 0x04, 0x14, 0x24}
//	log.Debug("三顺子%d", lg.GetCardType(c63[:]))
//	c64 := [...]int{0x03, 0x33, 0x23, 0x04, 0x14, 0x24, 0x01, 0x02}
//	log.Debug("飞机带两单%d", lg.GetCardType(c64[:]))
//	c65 := [...]int{0x03, 0x33, 0x23, 0x04, 0x14, 0x24, 0x01, 0x11, 0x02, 0x12}
//	log.Debug("飞机带两对%d", lg.GetCardType(c65[:]))
//	c41 := [...]int{0x03, 0x33, 0x23, 0x13, 0x14, 0x25}
//	log.Debug("四带两单%d", lg.GetCardType(c41[:]))
//	c42 := [...]int{0x03, 0x33, 0x23, 0x13, 0x14, 0x24, 0x15, 0x25}
//	log.Debug("四带两对%d", lg.GetCardType(c42[:]))
//	c40 := [...]int{0x03, 0x33, 0x23, 0x13}
//	log.Debug("炸弹%d", lg.GetCardType(c40[:]))
//
//	var ck []int
//	for i := 0; i < 8; i++ {
//		ck = append(ck, 0x4E+i%2)
//		log.Debug("八王类型%d", lg.GetCardType(ck[:]))
//	}
//}

func init() {
	Wg.Add(1)
	conf.Init("/Users/zhangyudong/Documents/GIT/src/mj/gameServer/gameApp/gameServer.json")
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

	temp, ok := base.GameServiceOptionCache.Get(29, 1)
	if !ok {
		return
	}

	log.Debug("tmp=%v", temp)

	info := &model.CreateRoomInfo{
		RoomId:       777777,
		MaxPlayerCnt: 3,
		KindId:       29,
		ServiceId:    1,
		Num:          1,
	}

	setCfg := map[string]interface{}{
		"EightKing": 1,
		"GameType":  2,
	}
	myCfg, cfgOk := json.Marshal(setCfg)
	if cfgOk != nil {
		log.Error("测试错误，退出程序")
		os.Exit(0)
	}
	info.OtherInfo = string(myCfg)

	_roombase := room_base.NewRoomBase()

	userg := room_base.NewRoomUserMgr(info.RoomId, info.MaxPlayerCnt, temp)

	u1 = newTestUser(1)
	u1.ChairId = 0
	userg.Users[0] = u1
	r := NewDDZEntry(info)
	ddzMrg = NewDDZDataMgr(info, u1.Id, pk_base.IDX_DDZ, "", temp, r)
	if ddzMrg == nil {
		log.Error("测试错误，退出程序")
		os.Exit(0)
	}
	cfg := &pk_base.NewPKCtlConfig{
		BaseMgr:  _roombase,
		DataMgr:  ddzMrg,
		UserMgr:  userg,
		LogicMgr: NewDDZLogic(pk_base.IDX_DDZ, info),
		TimerMgr: room_base.NewRoomTimerMgr(info.Num, temp),
	}
	r.Init(cfg)
	room = r

	var userCnt = 3

	for i := 1; i < userCnt; i++ {
		u := newTestUser(i + 1)
		if i == 1 {
			u2 = u
		} else if i == 2 {
			u3 = u
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
func InitLog() {
	logger, err := log.New(conf.Server.LogLevel, "", conf.LogFlag)
	if err != nil {
		panic(err)
	}
	log.Export(logger)
}
