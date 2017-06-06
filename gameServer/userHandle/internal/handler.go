package internal

import (
	"mj/common/msg"
	"reflect"
	"github.com/lovelly/leaf/log"
	"errors"
	. "mj/common/cost"
	"mj/gameServer/db/model"
	"mj/gameServer/db/model/base"
	"mj/gameServer/user"
	"fmt"
	"mj/gameServer/center"
)


//注册 客户端消息调用
func handlerC2S(m *Module, msg interface{}, h interface{}) {
	m.ChanRPC.Register(reflect.TypeOf(msg), h)
}

func RegisterHandler(m *Module) {
	//注册rpc 消息
	m.ChanRPC.Register("handleMsgData", m.handleMsgData)
	m.ChanRPC.Register("NewAgent", m.NewAgent)
	m.ChanRPC.Register("CloseAgent", m.CloseAgent)
	m.ChanRPC.Register("WriteUserScore", m.WriteUserScore)

	//c2s
	handlerC2S(m, &msg.C2G_GR_LogonMobile{}, m.handleMBLogin)
	handlerC2S(m, &msg.C2G_REQUserInfo{}, m.GetUserInfo)
}

//连接进来的通知
func  (m *Module)NewAgent(args []interface{}) error {
	log.Debug("at game NewAgent")
	return nil
}

//连接关闭的通知
func  (m *Module)CloseAgent (args []interface{}) error {
	log.Debug("at game CloseAgent")
	agent :=  m.a
	user, ok := agent.UserData().(*user.User)
	if !ok {
		return nil
	}

	DelUser(user.Id)
	center.ChanRPC.Go("SelfNodeDelPlayer", user.Id)
	return nil
}


func (m *Module) GetUserInfo(args []interface{}) {
	log.Debug("at GetUserInfo ................ ")
}


func(m *Module) handleMBLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_GR_LogonMobile)
	retMsg := &msg.G2C_LogonFinish{}
	agent := m.a
	retcode := 0
	defer func() {
		if retcode != 0 {
			str := fmt.Sprintf("登录失败, 错误码: %d",retcode)
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

	if HasUser(accountData.UserID){
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
	user.Accountsinfo = accountData
	user.Id = accountData.UserID
	user.ChairId = INVALID_CHAIR
	user.RoomId = INVALID_CHAIR

	lok := loadUser(user)
	if !lok {
		retcode = LoadUserInfoError
		return
	}

	AddUser(user.Id)
	center.ChanRPC.Go("SelfNodeAddPlayer", user.Id, agent.ChanRPC())
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
		GameID : user.GameID,						//游戏 I D
		UserID : user.Id,							//用户 I D
		FaceID : user.FaceID,							//头像索引
		CustomID :user.CustomID,						//自定标识
		Gender :user.Gender,							//用户性别
		MemberOrder :user.Accountsinfo.MemberOrder,					//会员等级
		TableID : user.RoomId,							//桌子索引
		ChairID : user.ChairId,							//椅子索引
		UserStatus :user.Status,						//用户状态
		Score :user.Score,								//用户分数
		WinCount : user.WinCount,							//胜利盘数
		LostCount : user.LostCount,						//失败盘数
		DrawCount : user.DrawCount,						//和局盘数
		FleeCount : user.FleeCount,						//逃跑盘数
		Experience : user.Experience,						//用户经验
		NickName: user.NickName,				//昵称
		HeaderUrl :user.HeadImgUrl, 				//头像
	})
}

func (m *Module)WriteUserScore(args []interface{}){
	log.Debug("at WriteUserScore === %v", args)
}




















/////////////////////////////// help 函数
///////
func loadUser(u *user.User) ( bool){
	ainfo, aok := model.AccountsmemberOp.Get(u.Id, u.Accountsinfo.MemberOrder)
	if !aok {
		log.Error("at loadUser not foud AccountsmemberOp by user", u.Id)
		return false
	}

	log.Debug("load user : == %v", ainfo)
	u.Accountsmember = ainfo

	glInfo, glok := model.GamescorelockerOp.Get(u.Id)
	if !glok {
		log.Error("at loadUser not foud GamescorelockerOp by user %d", u.Id)
		return  false
	}
	u.Gamescorelocker = glInfo

	giInfom, giok := model.GamescoreinfoOp.Get(u.Id)
	if !giok {
		log.Error("at loadUser not foud GamescoreinfoOp by user  %d", u.Id)
		return  false
	}
	u.Gamescoreinfo = giInfom

	ucInfo, uok := model.UserattrOp.Get(u.Id)
	if !uok {
		log.Error("at loadUser not foud UserroomcardOp by user  %d", u.Id)
		return  false
	}
	u.Userattr = ucInfo

	uextInfo, ueok := model.UserextrainfoOp.Get(u.Id)
	if !ueok {
		log.Error("at loadUser not foud UserextrainfoOp by user  %d", u.Id)
		return  false
	}
	u.Userextrainfo = uextInfo
	return  true
}










/////主消息函数
func (m *Module) handleMsgData(args []interface{}) (error) {
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

