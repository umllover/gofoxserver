package internal

import (
	"errors"
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"reflect"

	"mj/gameServer/common"

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
	m.ChanRPC.Register("RoomCloseUserOffline", m.RoomClose)
	//c2s
	handlerC2S(m, &msg.C2G_GR_LogonMobile{}, m.handleMBLogin)
	handlerC2S(m, &msg.C2G_REQUserInfo{}, m.GetUserInfo)
}

//连接进来的通知
func (m *UserModule) NewAgent(args []interface{}) error {
	log.Debug("at game NewAgent")
	return nil
}

//房间关闭的时候通知
func (m *UserModule) RoomClose(args []interface{}) error {
	user := m.a.UserData().(*user.User)
	if user.IsOffline() {
		DelUser(user.Id)
		m.Close(common.UserOffline)
	}
	return nil
}

//连接关闭的通知
func (m *UserModule) CloseAgent(args []interface{}) error {
	log.Debug("at game CloseAgent")
	agent := m.a
	user, ok := agent.UserData().(*user.User)
	if !ok {
		return nil
	}
	if user.SetOffline(true) {
		DelUser(user.Id)
		m.Close(common.UserOffline)
	} else {
		//等待房间结束
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

	if HasUser(accountData.UserID) {
		retcode = ErrUserReLogin
		return
	}

	//if accountData.PasswordID != recvMsg.Password {
	// retcode = ErrPasswd
	//	return
	//}

	template, ok := base.GameServiceOptionCache.Get(recvMsg.KindID, recvMsg.ServerID)
	if !ok {
		retcode = NoFoudTemplate
		return
	}

	user := user.NewUser(accountData.UserID)
	user.Agent = agent
	user.Id = accountData.UserID
	user.ChairId = INVALID_CHAIR
	user.RoomId = INVALID_CHAIR

	lok := loadUser(user)
	if !lok {
		retcode = LoadUserInfoError
		return
	}

	AddUser(user.Id, user)
	user.KindID = recvMsg.KindID
	user.ServerID = recvMsg.ServerID

	agent.WriteMsg(retMsg)
	agent.SetUserData(user)

	agent.WriteMsg(&msg.G2C_ConfigServer{
		TableCount: template.TableCount,
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
	user := m.a.UserData().(*user.User)
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

/////////////////////////////// help 函数
///////
func loadUser(u *user.User) bool {
	ainfo, aok := model.AccountsmemberOp.Get(u.Id)
	if !aok {
		log.Error("at loadUser not foud AccountsmemberOp by user", u.Id)
		return false
	}

	log.Debug("load user : == %v", ainfo)
	u.Accountsmember = ainfo

	glInfo, glok := model.GamescorelockerOp.Get(u.Id)
	if !glok {
		log.Error("at loadUser not foud GamescorelockerOp by user %d", u.Id)
		return false
	}
	u.Gamescorelocker = glInfo

	giInfom, giok := model.GamescoreinfoOp.Get(u.Id)
	if !giok {
		log.Error("at loadUser not foud GamescoreinfoOp by user  %d", u.Id)
		return false
	}
	u.Gamescoreinfo = giInfom

	ucInfo, uok := model.UserattrOp.Get(u.Id)
	if !uok {
		log.Error("at loadUser not foud UserroomcardOp by user  %d", u.Id)
		return false
	}
	u.Userattr = ucInfo

	uextInfo, ueok := model.UserextrainfoOp.Get(u.Id)
	if !ueok {
		log.Error("at loadUser not foud UserextrainfoOp by user  %d", u.Id)
		return false
	}
	u.Userextrainfo = uextInfo

	userToken, tok := model.UsertokenOp.Get(u.Id)
	if !tok {
		log.Error("at loadUser not foud UsertokenOp by user  %d", u.Id)
		return false
	}
	u.Usertoken = userToken
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

		if m.ChanRPC.HasFunc(msgType) {
			m.ChanRPC.Go(msgType, data, m.a)
			return nil
		}

		err = msg.Processor.RouteByType(msgType, data, m.a)
		if err != nil {
			return err
		}
	}
	return nil
}
