package internal

import (
	"encoding/json"
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/register"
	"mj/gameServer/RoomMgr"
	"mj/gameServer/common"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/kindList"
	"mj/gameServer/user"
	client "mj/gameServer/user"

	"mj/gameServer/db/model/stats"
	"time"

	"github.com/lovelly/leaf/log"
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
	reg.RegisterRpc("SvrShutdown", m.SvrShutdown)

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
	reg.RegisterC2S(&msg.C2G_LeaveRoom{}, m.ReqLeaveRoom)
	reg.RegisterC2S(&msg.C2G_ReplyLeaveRoom{}, m.ReplyLeaveRoom)

}

//连接进来的通知
func (m *UserModule) NewAgent(args []interface{}) error {
	log.Debug("at game NewAgent")
	return nil
}

//房间关闭的时候通知
func (m *UserModule) LeaveRoom(args []interface{}) error {
	log.Debug("at user LeaveRoom ...........")
	m.Close(KickOutGameEnd)
	return nil
}

//连接关闭的通知
func (m *UserModule) CloseAgent(args []interface{}) error {
	log.Debug("at game CloseAgent")
	Reason := args[1].(int)
	agent := m.a
	player, ok := agent.UserData().(*client.User)
	if !ok || player == nil {
		log.Error("at CloseAgent not foud user")
		return nil
	}

	if player.RoomId != 0 {
		r := RoomMgr.GetRoom(player.RoomId)
		if r != nil {
			r.GetChanRPC().Go("userOffline", player)
		}
	}

	m.UserOffline()
	if Reason != KickOutMsg {
		DelUser(player.Id)
	}

	if Reason == UserOffline {
		m.Close(UserOffline)
	}

	return nil
}

func (m *UserModule) ForceClose(args []interface{}) {
	log.Debug("at ForceClose ..... ")
	m.Close(KickOutMsg)
}

func (m *UserModule) SvrShutdown(args []interface{}) {
	log.Debug("at SvrShutdown ..... ")
	m.Close(ServerKick)
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

	if accountData.LogonPass != recvMsg.Password {
		retcode = ErrPasswd
		return
	}

	user := client.NewUser(accountData.UserID)
	user.KindID = recvMsg.KindID
	user.ServerID = recvMsg.ServerID
	user.Id = accountData.UserID
	user.Status = US_FREE
	user.HallNodeName = GetHallSvrName(recvMsg.HallNodeID)
	lok := loadUser(user)
	if !lok {
		retcode = LoadUserInfoError
		return
	}
	user.ChairId = INVALID_CHAIR

	oldUser := getUser(accountData.UserID)
	if oldUser != nil {
		log.Debug("old user ====== %d  %d ", oldUser.KindID, oldUser.RoomId)
		oldUser.RoomId = 0
		user.Status = oldUser.Status
		user.ChatRoomId = oldUser.ChatRoomId
		m.KickOutUser(oldUser)
	}

	if user.RoomId != 0 {
		r := RoomMgr.GetRoom(user.RoomId)
		if r != nil { //原来房间没关闭，投递个消息看下原来是否在房间内
			r.GetChanRPC().Go("userRelogin", user)
		}
	}

	user.Agent = agent
	AddUser(user.Id, user)

	agent.SetUserData(user)
	agent.WriteMsg(&msg.G2C_ConfigServer{
		TableCount: common.TableFullCount,
		ChairCount: 4,
		ServerType: template.GameType,
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

func (m *UserModule) KickOutUser(player *user.User) {
	player.ChanRPC().Go("ForceClose")
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
	player, _ := m.a.UserData().(*client.User)
	recvMsg := args[0].(*msg.C2G_UserSitdown)
	retcode := 0
	defer func() {
		if retcode != 0 {
			player.WriteMsg(&msg.G2C_UserSitDownRst{Code: retcode})
		}
	}()
	if player.KindID == 0 {
		log.Error("at UserSitdown not foud module userid:%d", player.Id)
		retcode = ErrFindRoomError
		return
	}

	r := RoomMgr.GetRoom(recvMsg.TableID)
	if r == nil {
		if player.RoomId != 0 {
			r = RoomMgr.GetRoom(player.RoomId)
			if r == nil {
				m.LoadRoom([]interface{}{&msg.C2G_LoadRoom{RoomID: player.RoomId}})
				r = RoomMgr.GetRoom(player.RoomId)
			}
		}
		if r == nil {
			retcode = ErrFindRoomError
			log.Error("at UserSitdown not foud roomd userid:%d, roomId: %d and %d ", player.Id, player.RoomId, recvMsg.TableID)
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
	log.Debug("UserReady KindID=%d, RoomId=%d, userId=%d, ChairId=%d", user.KindID, user.RoomId, user.Id, user.ChairId)
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
	myLogInfo := make(map[string]interface{})
	logInfo := make(map[string]interface{})
	AddLogDb := stats.RoomLogOp
	logInfo["room_id"] = user.RoomId
	logInfo["kind_id"] = user.KindID
	logInfo["service_id"] = user.ServerID
	logData, err1 := AddLogDb.GetByMap(logInfo)
	if err1 != nil {
		log.Error("Select Data from recode Error:%v", err1.Error())
	}
	now := time.Now()
	myLogInfo["end_time"] = &now
	log.Debug("解散房间ddebug======================================================%d", user.RoomId)

	if user.KindID == 0 {
		log.Error("at DissumeRoom not foud module userid:%d", user.Id)
		myLogInfo["start_endError"] = 1
		err := AddLogDb.UpdateWithMap(logData.RecodeId, myLogInfo)
		if err != nil {
			log.Error("结束时间和结束状态记录更新失败：%s", err.Error())
		}
		return
	}

	if user.RoomId == 0 {
		log.Error("at DissumeRoom not foud roomd id userid:%d", user.Id)
		myLogInfo["start_endError"] = 1
		err := AddLogDb.UpdateWithMap(logData.RecodeId, myLogInfo)
		if err != nil {
			log.Error("结束时间和结束状态记录更新失败：%s", err.Error())
		}
		return
	}
	r := RoomMgr.GetRoom(user.RoomId)
	if r == nil {
		log.Error("at DissumeRoom not foud roomd userid:%d", user.Id)
		myLogInfo["start_endError"] = 1
		err := AddLogDb.UpdateWithMap(logData.RecodeId, myLogInfo)
		if err != nil {
			log.Error("结束时间和结束状态记录更新失败：%s", err.Error())
		}
		return
	}

	r.GetChanRPC().Go("DissumeRoom", user)
}

func (m *UserModule) ReqLeaveRoom(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_LeaveRoom)
	player := m.a.UserData().(*user.User)
	r := RoomMgr.GetRoom(player.RoomId)
	if r != nil {
		r.GetChanRPC().Go("ReqLeaveRoom", player)
	} else {
		player.WriteMsg(&msg.G2C_LeaveRoomRsp{Code: ErrPlayerNotInRoom})
	}
}

func (m *UserModule) ReplyLeaveRoom(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_ReplyLeaveRoom)
	player := m.a.UserData().(*user.User)
	r := RoomMgr.GetRoom(player.RoomId)
	if r != nil {
		r.GetChanRPC().Go("ReplyLeaveRoom", player, recvMsg.Agree, recvMsg.UserID)
	} else {
		log.Error("at ReplyLeaveRoom user not in room ")
	}
}

/////////////////////////////// help 函数
///////
func loadUser(u *client.User) bool {
	//data, err := cluster.TimeOutCall1(u.HallNodeName, 8, &msg.S2S_GetPlayerInfo{Uid: u.Id})
	//if err != nil {
	//	log.Error("get room data error :%v", err.Error())
	//	return false
	//}

	//info, ok := data.(*msg.S2S_GetPlayerInfoResult)
	//if !ok {
	//	log.Error("loadUser data is error")
	//	return false
	//}
	//log.Debug("get user data == %v", info)

	attr, ok := model.UserattrOp.Get(u.Id)
	if !ok {
		log.Error("loadUser data is error 11")
		return false
	}

	source, sok := model.GamescoreinfoOp.Get(u.Id)
	if !sok {
		log.Error("loadUser data is error source")
		return false
	}

	locker, lok := model.GamescorelockerOp.Get(u.Id)
	if !lok || locker.Roomid == 0 {
		log.Error("loadUser data is error locker .roomID :%v", locker.Roomid)
		return false
	}

	u.NickName = attr.NickName
	u.FaceID = attr.FaceID
	u.CustomID = attr.CustomID
	u.HeadImgUrl = attr.HeadImgUrl
	u.Experience = attr.Experience
	u.Gender = attr.Gender
	u.WinCount = source.WinCount
	u.LostCount = source.LostCount
	u.DrawCount = source.DrawCount
	u.FleeCount = source.FleeCount
	u.UserRight = attr.UserRight
	u.Score = 0
	u.Revenue = source.Revenue
	u.InsureScore = source.InsureScore
	u.MemberOrder = 0
	u.RoomId = locker.Roomid
	u.KindID = locker.KindID
	u.ServerID = locker.ServerID
	return true
}
