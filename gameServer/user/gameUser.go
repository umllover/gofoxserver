package user

import (
	"mj/common/msg"
	"sync"

	"mj/gameServer/RoomMgr"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

type User struct {
	gate.Agent
	robot      *Robot //设置玩家代理
	Id         int64  //唯一id
	NickName   string //名字
	RoomId     int    // roomId 就是tableid
	Status     int    //当前游戏状态
	offline    bool   //玩家是否在线github.com/labstack/gommon/color
	ChairId    int    //当前椅子
	UserLimit  int64  //限制行为
	ChatRoomId int    //聊天房间ID
	//Currency     int    //游戏豆
	//RoomCard     int    //房卡数
	FaceID      int8   // 头像标识
	KindID      int    // 房间索引
	ServerID    int    // 游戏标识
	CustomID    int    // 自定标识
	HeadImgUrl  string // 头像
	Experience  int    // 经验数值
	Gender      int8   // 性别
	WinCount    int    // 胜局数目
	LostCount   int    // 输局数目
	DrawCount   int    // 和局数目
	FleeCount   int    // 逃局数目
	UserRight   int    // 用户权限
	Score       int64  // 用户积分（货币）
	Revenue     int64  // 游戏税收
	InsureScore int64  // 银行金币
	MemberOrder int8   // 会员标识
	HallNodeId  int    //大厅服务器节点名字
	IconID      int    //头像id
	Sign        string //个性签名
	Star        int    //被点赞数
	mu          sync.RWMutex
}

func NewUser(UserId int64) *User {
	return &User{Id: UserId}
}

func (u *User) GetUid() int64 {
	return u.Id
}

func (u *User) SendSysMsg(ty int, context string) {
	u.WriteMsg(&msg.SysMsg{
		ClientID: u.Id,
		Type:     ty,
		Context:  context,
	})
}

func (u *User) WriteMsg(msg interface{}) {
	if u.robot != nil {
		u.robot.WriteMsg(msg)
	} else {
		u.Agent.WriteMsg(msg)
	}
}

/////////////////////////

func (u *User) IsOffline() bool {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.offline
}

func (u *User) SetRoomId(id int) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.RoomId = id
}

func (u *User) RunRobot() {
	room := RoomMgr.GetRoom(u.RoomId)
	if room == nil {
		log.Error("at RunRobot, cannt find room, roomid=%d", u.RoomId)
		return
	}
	u.robot = NewRobot(u.Id, u, room)
}

func (u *User) StopRobot() {
	if u.robot != nil {
		u.robot.Close()
		u.robot = nil
		log.Debug("at StopRobot stop success, uid=%d", u.Id)
	} else {
		log.Debug("at StopRobot u.robot is nil, uid=%d", u.Id)
	}
}
