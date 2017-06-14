package internal

import (
	"errors"
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	client "mj/gameServer/user"
	"reflect"

	"mj/gameServer/common"

	"mj/gameServer/RoomMgr"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/log"
)

//注册 客户端消息调用
func handlerC2S(m *UserModule, msg interface{}, h interface{}) {
	m.ChanRPC.Register(reflect.TypeOf(msg), h)
}

func RegisterHandler(m *UserModule) {
	//注册rpc 消息
	m.ChanRPC.Register("handleMsgData", m.handleMsgData)
	m.ChanRPC.Register("NewAgent", m.NewAgent)
	m.ChanRPC.Register("CloseAgent", m.CloseAgent)
	m.ChanRPC.Register("WriteUserScore", m.WriteUserScore)
	m.ChanRPC.Register("LeaveRoom", m.LeaveRoom)
	//c2s
	handlerC2S(m, &msg.C2G_GR_LogonMobile{}, m.handleMBLogin)
	handlerC2S(m, &msg.C2G_REQUserInfo{}, m.GetUserInfo)
	handlerC2S(m, &msg.C2G_UserSitdown{}, m.UserSitdown)
	handlerC2S(m, &msg.C2G_GameOption{}, m.SetGameOption)
	handlerC2S(m, &msg.C2G_UserStandup{}, m.UserStandup)
	handlerC2S(m, &msg.C2G_REQUserChairInfo{}, m.GetUserChairInfo)
	handlerC2S(m, &msg.C2G_UserReady{}, m.UserReady)
	handlerC2S(m, &msg.C2G_GR_UserChairReq{}, m.UserChairReq)
	handlerC2S(m, &msg.C2G_HostlDissumeRoom{}, m.DissumeRoom)
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
			r.GetChanRPC().Go("userOffline", user)
		}
	}

	return nil
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
			agent.WriteMsg(&msg.G2C_LogonFailur{ResultCode: retcode, DescribeString: str})
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
		user.RoomId = INVALID_CHAIR
	} else {
		log.Debug("old user ====== %d  %d ", user.KindID, user.RoomId)
		if user.KindID != 0 && user.RoomId != 0 {
			r := RoomMgr.GetRoom(user.RoomId)
			if r != nil {
				r.GetChanRPC().Go("userRelogin", user)
			}
		}
		user.HallNodeName = GetHallSvrName(recvMsg.HallNodeID)
	}

	user.Agent = agent
	AddUser(user.Id, user)
	agent.WriteMsg(retMsg)
	agent.SetUserData(user)
	agent.WriteMsg(&msg.G2C_ConfigServer{
		TableCount: common.TableFullCount,
		ChairCount: 4,
		ServerType: template.ServerType,
		ServerRule: template.ServerRule,
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
}

////////////////////// help
func (m *UserModule) UserOffline() {

}

func (m *UserModule) WriteUserScore(args []interface{}) {
	log.Debug("at WriteUserScore === %v", args)
	info := args[0].(*msg.TagScoreInfo)
	Type := args[0].(int)
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
	user := m.a.UserData().(*client.User)
	recvMsg := args[0].(*msg.C2G_UserSitdown)
	if user.KindID == 0 {
		log.Error("at UserSitdown not foud module userid:%d", user.Id)
		return
	}

	if user.RoomId == 0 {
		log.Error("at UserSitdown not foud roomd id userid:%d", user.Id)
		return
	}

	r := RoomMgr.GetRoom(recvMsg.TableID)
	if r == nil {
		log.Error("at UserSitdown not foud roomd userid:%d", user.Id)
		return
	}

	r.GetChanRPC().Go("Sitdown", args[0], user)
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
		log.Error("at UserSitdown not foud roomd userid:%d", user.Id)
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

//解散房间
func (m *UserModule) DissumeRoom(args []interface{}) {
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

	r.GetChanRPC().Go("DissumeRoom", user)
}

/////////////////////////////// help 函数
///////
func loadUser(u *client.User) bool {
	data, err := cluster.Call1(u.HallNodeName, "GetPlayerInfo", u.Id)
	if err != nil {
		log.Error("get room data error :%v", err.Error())
		return false
	}

	info, ok := data.(map[string]interface{})
	if !ok {
		log.Error("loadUser data is error")
		return false
	}

	log.Debug("get user data == %v", info)

	u.Id = info["Id"].(int)
	u.NickName = info["NickName"].(string)
	u.Currency = info["Currency"].(int)
	u.RoomCard = info["RoomCard"].(int)
	u.FaceID = info["FaceID"].(int8)
	u.CustomID = info["CustomID"].(int)
	u.HeadImgUrl = info["HeadImgUrl"].(string)
	u.Experience = info["Experience"].(int)
	u.Gender = info["Gender"].(int8)
	u.WinCount = info["WinCount"].(int)
	u.LostCount = info["LostCount"].(int)
	u.DrawCount = info["DrawCount"].(int)
	u.FleeCount = info["FleeCount"].(int)
	u.UserRight = info["UserRight"].(int)
	u.Score = info["Score"].(int64)
	u.Revenue = info["Revenue"].(int64)
	u.InsureScore = info["InsureScore"].(int64)
	u.MemberOrder = info["MemberOrder"].(int8)
	return true
}

/////主消息函数
func (m *UserModule) handleMsgData(args []interface{}) error {
	if msg.Processor != nil {
		str := args[0].([]byte)
		data, err := msg.Processor.Unmarshal(str)
		if err != nil {
			return err
		}

		msgType := reflect.TypeOf(data)
		if msgType == nil || msgType.Kind() != reflect.Ptr {
			return errors.New("json message pointer required 11")
		}

		f, ok := m.ChanRPC.HasFunc(msgType)
		if ok {
			m.ChanRPC.Exec(chanrpc.BuildGoCallInfo(f, data, m.a))
			return nil
		}

		err = msg.Processor.RouteByType(msgType, data, m.a)
		if err != nil {
			return err
		}
	}
	return nil
}
