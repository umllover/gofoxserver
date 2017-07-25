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

	"mj/common/msg/pk_sss_msg"

	"github.com/lovelly/leaf/chanrpc"
	lconf "github.com/lovelly/leaf/conf"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"

	dbg "github.com/funny/debug"
)

var (
	room    *SSS_Entry //Pk_base
	dataMgr *sss_data_mgr
	u1      *user.User
	u2      *user.User
	u3      *user.User
)

var Wg sync.WaitGroup

func TestGameStart_1(t *testing.T) {
	room.UserReady([]interface{}{nil, u1})
}

//func TestShowCard(t *testing.T) {
//	log.Debug("测试摊牌")
//	data := &pk_sss_msg.C2G_SSS_Open_Card{
//		FrontCard:   []int{0x01, 0x12, 0x24},
//		MidCard:     []int{0x21, 0x05, 0x15, 0x13, 0x08},
//		BackCard:    []int{0x02, 0x32, 0x22, 0x19, 0x04},
//		SpecialType: false,
//		SpecialData: []int{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D},
//		Dragon:      false,
//	}
//	args := []interface{}{data, u2}
//	room.ShowSSsCard(args)
//	log.Debug("测试摊牌结束")
//}
//
//func TestShowCard_1(t *testing.T) {
//	log.Debug("测试摊牌")
//	data := &pk_sss_msg.C2G_SSS_Open_Card{
//		FrontCard:   []int{0x01, 0x02, 0x04},
//		MidCard:     []int{0x15, 0x25, 0x13, 0x14, 0x16},
//		BackCard:    []int{0x34, 0x14, 0x11, 0x21, 0x3D},
//		SpecialType: false,
//		SpecialData: []int{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D},
//		Dragon:      false,
//	}
//	args := []interface{}{data, u1}
//	room.ShowSSsCard(args)
//	log.Debug("测试摊牌结束")
//}

func TestShowCard_Special(t *testing.T) {
	//log.Debug("测试摊牌_特殊牌")
	//data := &pk_sss_msg.C2G_SSS_Open_Card{}
	//args := []interface{}{}
	//至尊清龙
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x31, 0x32, 0x33},
	//	MidCard:     []int{0x34, 0x35, 0x36, 0x37, 0x38},
	//	BackCard:    []int{0x39, 0x3A, 0x3B, 0x3C, 0x3D},
	//	SpecialType: true,
	//	SpecialData: []int{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//一条龙
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x01, 0x02, 0x03},
	//	MidCard:     []int{0x34, 0x35, 0x36, 0x37, 0x38},
	//	BackCard:    []int{0x39, 0x3A, 0x3B, 0x3C, 0x3D},
	//	SpecialType: true,
	//	SpecialData: []int{0x01, 0x02, 0x03, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//十二皇族
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x0B, 0x0B, 0x0C},
	//	MidCard:     []int{0x0D, 0x1B, 0x1C, 0x1D, 0x2B},
	//	BackCard:    []int{0x2C, 0x2D, 0x3B, 0x3C, 0x3D},
	//	SpecialType: true,
	//	SpecialData: []int{0x0B, 0x0B, 0x0C, 0x0D, 0x1B, 0x1C, 0x1D, 0x2B, 0x2C, 0x2D, 0x3B, 0x3C, 0x3D},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//三同花顺
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x01, 0x02, 0x03},
	//	MidCard:     []int{0x13, 0x14, 0x15, 0x16, 0x17},
	//	BackCard:    []int{0x27, 0x28, 0x29, 0x2A, 0x2B},
	//	SpecialType: true,
	//	SpecialData: []int{0x01, 0x02, 0x03, 0x13, 0x14, 0x15, 0x16, 0x17, 0x27, 0x28, 0x29, 0x2A, 0x2B},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//三分天下
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x01, 0x11, 0x21},
	//	MidCard:     []int{0x31, 0x02, 0x12, 0x22, 0x32},
	//	BackCard:    []int{0x03, 0x13, 0x23, 0x33, 0x2B},
	//	SpecialType: true,
	//	SpecialData: []int{0x01, 0x11, 0x21, 0x31, 0x02, 0x12, 0x22, 0x32, 0x03, 0x13, 0x23, 0x33, 0x2B},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//全大
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x09, 0x19, 0x29},
	//	MidCard:     []int{0x0C, 0x1C, 0x2c, 0x3A, 0x1A},
	//	BackCard:    []int{0x2C, 0x1D, 0x39, 0x2A, 0x2B},
	//	SpecialType: true,
	//	SpecialData: []int{0x09, 0x19, 0x29, 0x0C, 0x1C, 0x2c, 0x3A, 0x1A, 0x2C, 0x1D, 0x39, 0x2A, 0x2B},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//全小
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x02, 0x12, 0x22},
	//	MidCard:     []int{0x07, 0x17, 0x27, 0x36, 0x16},
	//	BackCard:    []int{0x26, 0x16, 0x35, 0x25, 0x25},
	//	SpecialType: true,
	//	SpecialData: []int{0x02, 0x12, 0x22, 0x07, 0x17, 0x27, 0x36, 0x16, 0x26, 0x16, 0x35, 0x25, 0x25},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//凑一色
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x01, 0x22, 0x03},
	//	MidCard:     []int{0x05, 0x25, 0x07, 0x29, 0x0A},
	//	BackCard:    []int{0x2B, 0x0C, 0x2D, 0x04, 0x26},
	//	SpecialType: true,
	//	SpecialData: []int{0x01, 0x22, 0x03, 0x05, 0x25, 0x07, 0x29, 0x0A, 0x2B, 0x0C, 0x2D, 0x04, 0x26},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//四套冲三
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x01, 0x11, 0x21},
	//	MidCard:     []int{0x02, 0x12, 0x22, 0x03, 0x13},
	//	BackCard:    []int{0x23, 0x04, 0x14, 0x24, 0x2A},
	//	SpecialType: true,
	//	SpecialData: []int{0x01, 0x11, 0x21, 0x02, 0x12, 0x22, 0x03, 0x13, 0x23, 0x04, 0x14, 0x24, 0x2A},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//五对冲三
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x01, 0x11, 0x21},
	//	MidCard:     []int{0x02, 0x12, 0x03, 0x23, 0x04},
	//	BackCard:    []int{0x24, 0x05, 0x25, 0x06, 0x26},
	//	SpecialType: true,
	//	SpecialData: []int{0x01, 0x11, 0x21, 0x02, 0x12, 0x03, 0x23, 0x04, 0x24, 0x05, 0x25, 0x06, 0x26},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//六对半
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x01, 0x11, 0x02},
	//	MidCard:     []int{0x12, 0x03, 0x13, 0x04, 0x14},
	//	BackCard:    []int{0x05, 0x15, 0x06, 0x16, 0x17},
	//	SpecialType: true,
	//	SpecialData: []int{0x01, 0x11, 0x02, 0x12, 0x03, 0x13, 0x04, 0x14, 0x05, 0x15, 0x06, 0x16, 0x17},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//三顺子
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x02, 0x13, 0x24},
	//	MidCard:     []int{0x24, 0x15, 0x26, 0x37, 0x18},
	//	BackCard:    []int{0x07, 0x18, 0x19, 0x2A, 0x3B},
	//	SpecialType: true,
	//	SpecialData: []int{0x02, 0x13, 0x24, 0x24, 0x15, 0x26, 0x37, 0x18, 0x07, 0x18, 0x19, 0x2A, 0x3B},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)
	//三同花
	//data = &pk_sss_msg.C2G_SSS_Open_Card{
	//	FrontCard:   []int{0x01, 0x02, 0x03},
	//	MidCard:     []int{0x13, 0x15, 0x17, 0x18, 0x19},
	//	BackCard:    []int{0x29, 0x2A, 0x2B, 0x2C, 0x2D},
	//	SpecialType: true,
	//	SpecialData: []int{0x01, 0x11, 0x02, 0x12, 0x03, 0x13, 0x04, 0x14, 0x05, 0x15, 0x06, 0x16, 0x17},
	//	Dragon:      false,
	//}
	//args = []interface{}{data, u1}
	//room.ShowSSsCard(args)

	//log.Debug("测试摊牌_特殊牌结束")

}

func TestAll(t *testing.T) {
	var data *pk_sss_msg.C2G_SSS_Open_Card
	data = &pk_sss_msg.C2G_SSS_Open_Card{
		FrontCard:   []int{0x02, 0x13, 0x24},
		MidCard:     []int{0x24, 0x15, 0x26, 0x37, 0x18},
		BackCard:    []int{0x07, 0x18, 0x19, 0x2A, 0x3B},
		SpecialType: true,
		SpecialData: []int{0x02, 0x13, 0x24, 0x24, 0x15, 0x26, 0x37, 0x18, 0x07, 0x18, 0x19, 0x2A, 0x3B},
		Dragon:      false,
	}
	setCard(dataMgr, u1, data, true, false)

	data = &pk_sss_msg.C2G_SSS_Open_Card{
		FrontCard:   []int{0x01, 0x11, 0x02},
		MidCard:     []int{0x12, 0x03, 0x13, 0x04, 0x14},
		BackCard:    []int{0x05, 0x15, 0x06, 0x16, 0x17},
		SpecialType: true,
		SpecialData: []int{0x01, 0x11, 0x02, 0x12, 0x03, 0x13, 0x04, 0x14, 0x05, 0x15, 0x06, 0x16, 0x17},
		Dragon:      false,
	}
	setCard(dataMgr, u2, data, true, false)

	dataMgr.ComputeChOut()
	dataMgr.ComputeResult()

	dbg.Print(dataMgr)

}

func setCard(dataMgr *sss_data_mgr, u *user.User, data *pk_sss_msg.C2G_SSS_Open_Card, isSpecial bool, isDragon bool) {
	dataMgr.SpecialTypeTable[u] = isSpecial
	dataMgr.Dragon[u] = isDragon
	dataMgr.m_bSegmentCard[u] = append(dataMgr.m_bSegmentCard[u], data.FrontCard, data.MidCard, data.BackCard)
	dataMgr.m_bUserCardData[u] = make([]int, 0, 13)
	dataMgr.m_bUserCardData[u] = append(dataMgr.m_bUserCardData[u], data.FrontCard...)
	dataMgr.m_bUserCardData[u] = append(dataMgr.m_bUserCardData[u], data.MidCard...)
	dataMgr.m_bUserCardData[u] = append(dataMgr.m_bUserCardData[u], data.BackCard...)
}

func TestGameConclude(t *testing.T) {

}

func TestDispatchCardData(t *testing.T) {

}

func TestAnalyseCard(t *testing.T) {

}

func init() {
	Wg.Add(1)
	conf.Init("D:/go/src/mj/gameServer/gameApp/gameServer.json")
	//conf.Init("/new/src/mj/gameServer/gameApp/gameServer.json")
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

	temp, ok := base.GameServiceOptionCache.Get(30, 1)
	if !ok {
		return
	}

	//log.Debug("tmp=%v", temp)

	info := &model.CreateRoomInfo{
		RoomId:       777777,
		MaxPlayerCnt: 2,
		KindId:       30,
		ServiceId:    1,
		Num:          1,
	}

	setCfg := map[string]interface{}{
		"YanSe": 0,
	}
	myCfg, cfgOk := json.Marshal(setCfg)
	if cfgOk != nil {
		log.Error("测试错误，退出程序")
		os.Exit(0)
	}
	info.OtherInfo = string(myCfg)

	_roombase := room_base.NewRoomBase()

	//	userg := room_base.NewRoomUserMgr(info.RoomId, info.MaxPlayerCnt, temp)
	userg := room_base.NewRoomUserMgr(info, temp)

	u1 = newTestUser(1)
	u1.ChairId = 0
	userg.Users[0] = u1
	r := NewSSSEntry(info)
	datag := NewDataMgr(info.RoomId, u1.Id, pk_base.IDX_SSS, "", temp, r)
	if datag == nil {
		log.Error("测试错误，退出程序")
		os.Exit(0)
	}
	dataMgr = datag
	cfg := &pk_base.NewPKCtlConfig{
		BaseMgr:  _roombase,
		DataMgr:  datag,
		UserMgr:  userg,
		LogicMgr: NewSssZLogic(pk_base.IDX_SSS),
		TimerMgr: room_base.NewRoomTimerMgr(info.Num, temp),
	}
	r.Init(cfg)
	room = r

	var userCnt = 2

	room.DataMgr.InitRoom(userCnt)

	for i := 1; i < userCnt; i++ {
		u := newTestUser(int64(i) + 1)
		if i == 1 {
			u2 = u
		} else if i == 2 {
			u3 = u
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
func InitLog() {
	logger, err := log.New(conf.Server.LogLevel, "", conf.LogFlag)
	if err != nil {
		panic(err)
	}
	log.Export(logger)
}
