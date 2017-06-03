package internal

import (
	"mj/common/msg"
	"reflect"
	"github.com/name5566/leaf/log"
	"github.com/lovelly/leaf/gate"
	"mj/hallServer/center"
	"errors"
	. "mj/common/cost"
	"fmt"
	"time"
	"mj/hallServer/gameList"
	"mj/hallServer/db/model"
	"mj/hallServer/user"
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

	//c2s
	handlerC2S(m, &msg.C2L_Login{}, m.handleMBLogin)
	handlerC2S(m, &msg.C2L_Regist{}, m.handleMBRegist)

}

//连接进来的通知
func (m *Module)NewAgent(args []interface{}) error{
	log.Debug("at hall NewAgent")
	return nil
}

//连接关闭的同喜
func (m *Module)CloseAgent (args []interface{}) error {
	log.Debug("at hall CloseAgent")
	agent:= args[0].(gate.Agent)
	id, ok := agent.UserData().(int)
	if !ok {
		return nil
	}
	m.OnDestroy()
	DelUser(id)
	center.ChanRPC.Go("SelfNodeDelPlayer", id)
	return nil
}


func  (m *Module)handleMBLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_Login)
	retMsg := &msg.L2C_LogonSuccess{}
	agent := m.a
	retcode := 0
	defer func() {
		if retcode != 0 {
			str := fmt.Sprintf("登录失败, 错误码: %d",retcode)
			agent.WriteMsg(&msg.L2C_LogonFailure{ResultCode: retcode, DescribeString: str})
		} else {
			agent.WriteMsg(retMsg)
		}
	}()

	if recvMsg.Accounts == "" {
		retcode = ParamError
		return
	}

	accountData, ok := model.AccountsinfoOp.GetByMap(map[string]interface{}{
		"Accounts": recvMsg.Accounts,
	})

	if ok != nil || accountData == nil {
		retcode = NotFoudAccout
		return
	}

	if _, ok := Users[accountData.UserID]; ok {
		retcode = ErrUserReLogin
		return
	}

	//if accountData.PasswordID != recvMsg.Password {
	//	sendErrFunc("password is error")
	//	return
	//}

	user := user.NewUser(accountData.UserID)
	user.Accountsinfo = accountData
	user.Id = accountData.UserID
	lok := loadUser(user)
	if !lok {
		retcode = LoadUserInfoError
		return
	}

	Users[user.Id] = struct {}{}
	agent.SetUserData(accountData.UserID)
	BuildClientMsg(retMsg, user)
	center.ChanRPC.Go("SelfNodeAddPlayer", user.Id, agent.ChanRPC())
	gameList.ChanRPC.Go("sendGameList", agent)
}




func  (m *Module)handleMBRegist(args []interface{}) {
	retcode := 0
	recvMsg := args[0].(*msg.C2L_Regist)
	agent := args[1].(gate.Agent)
	retMsg := &msg.L2C_LogonSuccess{}
	defer func() {
		if retcode != 0 {
			agent.WriteMsg(&msg.L2C_LogonFailure{ResultCode: retcode, DescribeString: "登录失败"})
		} else {
			agent.WriteMsg(retMsg)
		}
	}()

	accountData, ok := model.AccountsinfoOp.GetByMap(map[string]interface{}{
		"Accounts": recvMsg.Accounts,
	})
	if ok == nil && accountData != nil {
		retcode = AlreadyExistsAccount
		return
	}

	//todo 名字排重等等等 验证
	now := time.Now()
	accInfo := &model.Accountsinfo{
		FaceID:   recvMsg.FaceID,   //头像标识
		Gender:   recvMsg.Gender,   //用户性别
		Accounts: recvMsg.Accounts, //登录帐号
		RegAccounts: recvMsg.Accounts,
		LogonPass: recvMsg.LogonPass,
		InsurePass: recvMsg.InsurePass,
		NickName: recvMsg.NickName, //用户昵称
		GameLogonTimes:1,
		LastLogonIP:agent.RemoteAddr().String(),
		LastLogonMobile:recvMsg.MobilePhone,
		LastLogonMachine:recvMsg.MachineID,
		RegisterMobile:recvMsg.MobilePhone,
		RegisterMachine: recvMsg.MachineID,
		RegisterDate : &now,
		RegisterIP:      agent.RemoteAddr().String(), //连接地址
	}

	lastid, err := model.AccountsinfoOp.Insert(accInfo)
	if err != nil {
		retcode = InsertAccountError
		return
	}
	accInfo.UserID = int(lastid)

	user, cok := createUser(accInfo.UserID)
	if !cok {
		retcode = CreateUserError
		return
	}
	Users[user.Id] = struct {}{}
	user.Accountsinfo = accInfo
	agent.SetUserData(accInfo.UserID)
	BuildClientMsg(retMsg, user)
	center.ChanRPC.Go("SelfNodeAddPlayer", user.Id, agent.ChanRPC())
}


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

func createUser(UserID int)  (*user.User, bool) {
	U := user.NewUser(UserID)
	U.Accountsmember = &model.Accountsmember{
		UserID:UserID,
	}
	_, err := model.AccountsmemberOp.Insert(U.Accountsmember)
	if err != nil {
		log.Error("at createUser insert Accountsmember error")
		return nil, false
	}

	now := time.Now()
	U.Gamescorelocker = &model.Gamescorelocker{
		UserID:UserID,
		CollectDate : &now,
	}
	_, err = model.GamescorelockerOp.Insert(U.Gamescorelocker)
	if err != nil {
		log.Error("at createUser insert Gamescorelocker error")
		return nil, false
	}

	U.Gamescoreinfo = &model.Gamescoreinfo{
		UserID:UserID,
		LastLogonDate: &now,
	}
	_, err = model.GamescoreinfoOp.Insert(U.Gamescoreinfo)
	if err != nil {
		log.Error("at createUser insert Gamescoreinfo error")
		return nil, false
	}

	U.Userattr = &model.Userattr{
		UserID:UserID,
	}
	_, err = model.UserattrOp.Insert(U.Userattr)
	if err != nil {
		log.Error("at createUser insert Userroomcard error")
		return nil, false
	}

	U.Userextrainfo = &model.Userextrainfo{
		UserId:UserID,
	}
	_, err = model.UserextrainfoOp.Insert(U.Userextrainfo)
	if err != nil {
		log.Error("at createUser insert Userroomcard error")
		return nil, false
	}

	return U, true
}

func BuildClientMsg(retMsg *msg.L2C_LogonSuccess, user *user.User){
	retMsg.FaceID = user.FaceID	//头像标识
	retMsg.Gender  = user.Gender
	retMsg.UserID  = user.Id
	retMsg.Spreader = user.SpreaderID
	retMsg.GameID  = user.GameID
	retMsg.Experience  = user.Experience
	retMsg.LoveLiness  = user.LoveLiness
	retMsg.NickName  = user.NickName

	//用户成绩
	retMsg.UserScore  = user.Score
	retMsg.UserInsure  = user.InsureScore
	retMsg.Medal  = user.UserMedal
	retMsg.UnderWrite = user.UnderWrite
	retMsg.WinCount   = user.WinCount
	retMsg.LostCount  = user.LostCount
	retMsg.DrawCount  = user.DrawCount
	retMsg.FleeCount = user.FleeCount
	tm := &msg.DateTime{}
	tm.Year = user.RegisterDate.Year()
	tm.DayOfWeek = int(user.RegisterDate.Weekday())
	tm.Day = user.RegisterDate.Day()
	tm.Hour = user.RegisterDate.Hour()
	tm.Second = user.RegisterDate.Second()
	tm.Minute = user.RegisterDate.Minute()
	retMsg.RegisterDate =tm

	//额外信息
	retMsg.MbTicket  = user.MbTicket
	retMsg.MbPayTotal = user.MbPayTotal
	retMsg.MbVipLevel  = user.MbVipLevel
	retMsg.PayMbVipUpgrade = user.PayMbVipUpgrade

	//约战房相关
	retMsg.RoomCard  = user.RoomCard
	retMsg.LockServerID  = user.ServerID
	retMsg.KindID  = user.KindID
}
























/////////////////////////////// help 函数

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

