package internal

import (
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/kindList"
	client "mj/gameServer/user"

	"mj/common/register"

	"encoding/json"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/nsq/cluster"
)

func RegisterHandler(m *UserModule) {
	reg := register.NewRegister(m.ChanRPC)
	//注册rpc 消息
	reg.RegisterRpc("handleMsgData", m.handleMsgData)
	reg.RegisterRpc("NewAgent", m.NewAgent)
	reg.RegisterRpc("CloseAgent", m.CloseAgent)
	reg.RegisterRpc("WriteUserScore", m.WriteUserScore)
	reg.RegisterRpc("LeaveRoom", m.LeaveRoom)
	reg.RegisterRpc("ForceClose", m.ForceClose)

	//c2s
	reg.RegisterC2S(&msg.C2G_GR_LogonMobile{}, m.handleMBLogin)
	reg.RegisterC2S(&msg.C2G_REQUserInfo{}, m.GetUserInfo)
	reg.RegisterC2S(&msg.C2G_UserSitdown{}, m.UserSitdown)
	reg.RegisterC2S(&msg.C2G_GameOption{}, m.SetGameOption)
	reg.RegisterC2S(&msg.C2G_UserStandup{}, m.UserStandup)
	reg.RegisterC2S(&msg.C2G_REQUserChairInfo{}, m.GetUserChairInfo)
	reg.RegisterC2S(&msg.C2G_UserReady{}, m.UserReady)
	reg.RegisterC2S(&msg.C2G_GR_UserChairReq{}, m.UserChairReq)
	reg.RegisterC2S(&msg.C2G_HostlDissumeRoom{}, m.DissumeRoom)
	reg.RegisterC2S(&msg.C2G_LoadRoom{}, m.LoadRoom)

}

//连接进来的通知
func (m *UserModule) NewAgent(args []interface{}) error {
	log.Debug("at game NewAgent")
	return nil
}

//房间关闭的时候通知
func (m *UserModule) LeaveRoom(args []interface{}) error {
	log.Debug("at user LeaveRoom ...........")
	user := m.a.UserData().(*client.User)
	//if user.IsOffline() { //只有离线了， 才删除玩家 todo
	DelUser(user.Id)
	m.Close(common.UserOffline)
	//}
	return nil
}

//连接关闭的通知
func (m *UserModule) CloseAgent(args []interface{}) error {
	log.Debug("at game CloseAgent")
	agent := m.a
	user, ok := agent.UserData().(*client.User)
	if !ok {
		return nil
	}
	if user.SetOffline(true) {
		DelUser(user.Id)
		m.Close(common.UserOffline)
	} else {
		if user.RoomId != 0 {
			r := RoomMgr.GetRoom(user.RoomId)
			if r != nil {
				r.GetChanRPC().Go("userOffline", user)
			}
		}
	}

	return nil
}

func (m *UserModule) ForceClose(args []interface{}) {
	log.Debug("at ForceClose ..... ")
	m.Close(common.KickOutOffline)
}

func (m *UserModule) GetUserInfo(args []interface{}) {
	log.Debug("at GetUserInfo ................ ")
}

func (m *UserModule) handleMBLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_GR_LogonMobile)
	retMsg := &msg.G2C_LogonFinish{}
	agent := m.a
	retcode := 0
	defer func() {
		if retcode != 0 {
			str := fmt.Sprintf("登录失败, 错误码: %d", retcode)
			agent.WriteMsg(&msg.G2C_LogonFailure{ResultCode: retcode, DescribeString: str})
		} else {

		}
	}()

	if recvMsg.UserID == 0 {
		retcode = ParamError
		return
	}

	accountData, ok := model.AccountsinfoOp.Get(recvMsg.UserID)
	if !ok || accountData == nil {
		retcode = NotFoudAccout
		return
	}

	template, ok := base.GameServiceOptionCache.Get(recvMsg.KindID, recvMsg.ServerID)
	if !ok {
		retcode = NoFoudTemplate
		return
	}

	user, ok := getUser(accountData.UserID)
	if ok && !user.IsOffline() {
		retcode = ErrUserDoubleLogin
		return
	}

	//if accountData.PasswordID != recvMsg.Password {
	// retcode = ErrPasswd
	//	return
	//}

	if user == nil {
		user = client.NewUser(accountData.UserID)
		user.KindID = recvMsg.KindID
		user.ServerID = recvMsg.ServerID
		user.Id = accountData.UserID
		user.HallNodeName = GetHallSvrName(recvMsg.HallNodeID)
		lok := loadUser(user)
		if !lok {
			retcode = LoadUserInfoError
			return
		}
		user.ChairId = INVALID_CHAIR
	} else {
		log.Debug("old user ====== %d  %d ", user.KindID, user.RoomId)
		if user.KindID != 0 && user.RoomId != 0 {
			r := RoomMgr.GetRoom(user.RoomId)
			if r != nil {
				r.GetChanRPC().Go("userRelogin", user)
			}
		}
		user.ChanRPC().Go("ForceClose")
		user.HallNodeName = GetHallSvrName(recvMsg.HallNodeID)
	}

	user.Agent = agent
	AddUser(user.Id, user)

	agent.SetUserData(user)
	agent.WriteMsg(&msg.G2C_ConfigServer{
		TableCount: common.TableFullCount,
		ChairCount: 4,
		ServerType: template.ServerType,
		ServerRule: 0, //废弃字段
	})

	agent.WriteMsg(&msg.G2C_ConfigFinish{})

	agent.WriteMsg(&msg.G2C_UserEnter{
		UserID:      user.Id,          //用户 I D
		FaceID:      user.FaceID,      //头像索引
		CustomID:    user.CustomID,    //自定标识
		Gender:      user.Gender,      //用户性别
		MemberOrder: user.MemberOrder, //会员等级
		TableID:     user.RoomId,      //桌子索引
		ChairID:     user.ChairId,     //椅子索引
		UserStatus:  user.Status,      //用户状态
		Score:       user.Score,       //用户分数
		WinCount:    user.WinCount,    //胜利盘数
		LostCount:   user.LostCount,   //失败盘数
		DrawCount:   user.DrawCount,   //和局盘数
		FleeCount:   user.FleeCount,   //逃跑盘数
		Experience:  user.Experience,  //用户经验
		NickName:    user.NickName,    //昵称
		HeaderUrl:   user.HeadImgUrl,  //头像
	})

	agent.WriteMsg(retMsg)
}

////////////////////// help
func (m *UserModule) UserOffline() {

}

func (m *UserModule) WriteUserScore(args []interface{}) {
	log.Debug("at WriteUserScore === %v", args)
	info := args[0].(*msg.TagScoreInfo)
	Type := args[1].(int)
	user := m.a.UserData().(*client.User)
	user.Score += int64(info.Score)
	user.Revenue += int64(info.Revenue)
	user.InsureScore += 0 //todo
	if info.IsWin == 1 {  //1 胜利 2失败 3逃跑
		user.WinCount += 1
	} else if info.IsWin == 2 {
		user.LostCount += 1
	} else if info.IsWin == 3 {
		user.FleeCount += 1
	} else {
		user.DrawCount += 1
	}

	model.GamescoreinfoOp.UpdateWithMap(user.Id, map[string]interface{}{
		"Score":       user.Score,
		"Revenue":     user.Revenue,
		"InsureScore": user.InsureScore,
		"WinCount":    user.WinCount,
		"LostCount":   user.LostCount,
		"FleeCount":   user.FleeCount,
		"DrawCount":   user.DrawCount,
	})

	//todo log
	_ = Type

}

func (m *UserModule) UserSitdown(args []interface{}) {
	player := m.a.UserData().(*client.User)
	recvMsg := args[0].(*msg.C2G_UserSitdown)
	if player.KindID == 0 {
		log.Error("at UserSitdown not foud module userid:%d", player.Id)
		return
	}

	roomid := recvMsg.TableID
	r := RoomMgr.GetRoom(recvMsg.TableID)
	if r == nil {
		if player.RoomId != 0 {
			roomid = player.RoomId
			m.LoadRoom([]interface{}{&msg.C2G_LoadRoom{RoomID: player.RoomId}})
			r = RoomMgr.GetRoom(player.RoomId)
		}
		if r == nil {
			log.Error("at UserSitdown not foud roomd userid:%d, roomId: %d", player.Id, roomid)
			return
		}
	}

	r.GetChanRPC().Go("Sitdown", recvMsg.ChairID, player)
}

func (m *UserModule) SetGameOption(args []interface{}) {
	user := m.a.UserData().(*client.User)
	if user.KindID == 0 {
		log.Error("at UserSitdown not foud module userid:%d", user.Id)
		return
	}

	if user.RoomId == 0 {
		log.Error("at UserSitdown not foud roomd id userid:%d", user.Id)
		return
	}

	r := RoomMgr.GetRoom(user.RoomId)
	if r == nil {
		log.Error("at UserSitdown not foud roomd:%v, userid:%d", user.RoomId, user.Id)
		return
	}

	r.GetChanRPC().Go("SetGameOption", args[0], user)
}

func (m *UserModule) UserReady(args []interface{}) {
	user := m.a.UserData().(*client.User)
	if user.KindID == 0 {
		log.Error("at UserSitdown not foud module userid:%d", user.Id)
		return
	}

	if user.RoomId == 0 {
		log.Error("at UserSitdown not foud roomd id userid:%d", user.Id)
		return
	}

	r := RoomMgr.GetRoom(user.RoomId)
	if r == nil {
		log.Error("at UserSitdown not foud roomd userid:%d", user.Id)
		return
	}

	r.GetChanRPC().Go("UserReady", args[0], user)

}
func (m *UserModule) GetUserChairInfo(args []interface{}) {
	user := m.a.UserData().(*client.User)
	if user.KindID == 0 {
		log.Error("at UserSitdown not foud module userid:%d", user.Id)
		return
	}

	if user.RoomId == 0 {
		log.Error("at UserSitdown not foud roomd id userid:%d", user.Id)
		return
	}

	r := RoomMgr.GetRoom(user.RoomId)
	if r == nil {
		log.Error("at UserSitdown not foud roomd userid:%d", user.Id)
		return
	}

	r.GetChanRPC().Go("GetUserChairInfo", args[0], user)
}

//起立
func (m *UserModule) UserStandup(args []interface{}) {
	user := m.a.UserData().(*client.User)
	if user.KindID == 0 {
		log.Error("at UserSitdown not foud module userid:%d", user.Id)
		return
	}

	if user.RoomId == 0 {
		log.Error("at UserSitdown not foud roomd id userid:%d", user.Id)
		return
	}

	r := RoomMgr.GetRoom(user.RoomId)
	if r == nil {
		log.Error("at UserSitdown not foud roomd userid:%d", user.Id)
		return
	}

	r.GetChanRPC().Go("UserStandup", args[0], user)

}

//客户端请求更换椅子
func (m *UserModule) UserChairReq(args []interface{}) {

}

//创建房间
func (m *UserModule) LoadRoom(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_LoadRoom)
	retMsg := &msg.G2C_LoadRoomOk{}
	player := m.a.UserData().(*client.User)
	retCode := -1
	defer func() {
		if retCode != 0 {
			player.WriteMsg(&msg.G2C_InitRoomFailure{ErrorCode: retCode, DescribeString: "创建房间失败"})
		} else {
			player.WriteMsg(retMsg)
		}
	}()
	info, err := model.CreateRoomInfoOp.GetByMap(map[string]interface{}{
		"room_id": recvMsg.RoomID,
	})
	if err != nil || info == nil {
		log.Error("at LoadRoom error :%v", err)
		retCode = ErrNotFoundCreateRecord
		return
	}

	b, _ := json.Marshal(info)
	log.Debug("at LoadRoom Info == %v", string(b))
	if info.Status != 0 {
		retCode = ErrDoubleCreaterRoom
		return
	}

	mod, ok := kindList.GetModByKind(info.KindId)
	if !ok {
		retCode = ErrNotFoundCreateRecord
		return
	}

	log.Debug("begin CreateRoom.....")
	ok1 := mod.CreateRoom(info, player)
	if !ok1 {
		retCode = ErrCreaterError
		return
	}

	retCode = 0
	return
}

//解散房间
func (m *UserModule) DissumeRoom(args []interface{}) {
	user := m.a.UserData().(*client.User)
	if user.KindID == 0 {
		log.Error("at DissumeRoom not foud module userid:%d", user.Id)
		return
	}

	if user.RoomId == 0 {
		log.Error("at DissumeRoom not foud roomd id userid:%d", user.Id)
		return
	}

	r := RoomMgr.GetRoom(user.RoomId)
	if r == nil {
		log.Error("at DissumeRoom not foud roomd userid:%d", user.Id)
		return
	}

	r.GetChanRPC().Go("DissumeRoom", user)
}

/////////////////////////////// help 函数
///////
func loadUser(u *client.User) bool {
	data, err := cluster.Call1(u.HallNodeName, &msg.S2S_GetPlayerInfo{Uid: u.Id})
	if err != nil {
		log.Error("get room data error :%v", err.Error())
		return false
	}

	info, ok := data.(*msg.S2S_GetPlayerInfoResult)
	if !ok {
		log.Error("loadUser data is error")
		return false
	}

	log.Debug("get user data == %v", info)
	u.Id = info.Id
	u.NickName = info.NickName
	u.Currency = info.Currency
	u.RoomCard = info.RoomCard
	u.FaceID = info.FaceID
	u.CustomID = info.CustomID
	u.HeadImgUrl = info.HeadImgUrl
	u.Experience = info.Experience
	u.Gender = info.Gender
	u.WinCount = info.WinCount
	u.LostCount = info.LostCount
	u.DrawCount = info.DrawCount
	u.FleeCount = info.FleeCount
	u.UserRight = info.UserRight
	u.Score = info.Score
	u.Revenue = info.Revenue
	u.InsureScore = info.InsureScore
	u.MemberOrder = info.MemberOrder
	u.RoomId = info.RoomId
	return true
}
