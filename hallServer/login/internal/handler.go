package internal

import (
	. "mj/common/cost"
	"mj/common/msg"
	"mj/hallServer/UserData"
	"mj/hallServer/db/model"
	"mj/hallServer/user"
	"reflect"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/gate"
	"mj/hallServer/gameList"
)

var userDatach = UserData.ChanRPC

////注册rpc 消息
func handleRpc(id interface{}, f interface{}, fType int) {
	cluster.SetRoute(id, ChanRPC)
	ChanRPC.RegisterFromType(id, f, fType)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {

	handlerC2S(&msg.C2L_Login{}, handleMBLogin)
	handlerC2S(&msg.C2L_Regist{}, handleMBRegist)
}

func handleMBLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_Login)
	retMsg := &msg.CMD_MB_LogonSuccess{}
	agent := args[1].(gate.Agent)
	retcode := 0
	defer func() {
		if retcode != 0 {
			agent.WriteMsg(&msg.CMD_GP_LogonFailure{ResultCode: retcode, DescribeString: "登录失败"})
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

	//if accountData.PasswordID != recvMsg.Password {
	//	sendErrFunc("password is error")
	//	return
	//}

	user, lok := loadUser(accountData.UserID)
	if !lok {
		retcode = LoadUserInfoError
		return
	}

	user.Accountsinfo = accountData
	agent.SetUserData(accountData.UserID)
	BuildClientMsg(retMsg, user)
	userDatach.Go("addUser", user)
	gameList.ChanRPC.Go("sendGameList", agent)
}

func handleMBRegist(args []interface{}) {
	retcode := 0
	recvMsg := args[0].(*msg.C2L_Regist)
	agent := args[1].(gate.Agent)
	retMsg := &msg.CMD_MB_LogonSuccess{}
	defer func() {
		if retcode != 0 {
			agent.WriteMsg(&msg.CMD_GP_LogonFailure{ResultCode: retcode, DescribeString: "登录失败"})
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
	accInfo := &model.Accountsinfo{
		FaceID:   recvMsg.FaceID,   //头像标识
		Gender:   recvMsg.Gender,   //用户性别
		Accounts: recvMsg.Accounts, //登录帐号
		NickName: recvMsg.NickName, //用户昵称

		//密码变量
		LogonPass:  recvMsg.LogonPass,  //登录密码
		InsurePass: recvMsg.InsurePass, //银行密码

		//附加信息
		RegisterIP:      agent.RemoteAddr().String(), //连接地址
		RegisterMachine: recvMsg.MachineID,           //机器序列
		RegisterMobile:  recvMsg.MobilePhone,         //电话号码
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
	user.Accountsinfo = accInfo
	agent.SetUserData(accInfo.UserID)
	BuildClientMsg(retMsg, user)
	userDatach.Go("addUser", user)
}


///////
func loadUser(UserID int) (*user.User, bool){
	U := user.NewUser(UserID)

	ainfo, aok := model.AccountsmemberOp.Get(UserID, 1) // todo ...
	if !aok {
		return nil, false
	}
	U.Accountsmember = ainfo

	glInfo, glok := model.GamescorelockerOp.Get(UserID)
	if !glok {
		return nil, false
	}
	U.Gamescorelocker = glInfo

	giInfom, giok := model.GamescoreinfoOp.Get(UserID)
	if !giok {
		return nil, false
	}
	U.Gamescoreinfo = giInfom

	ucInfo, uok := model.UserroomcardOp.Get(UserID)
	if !uok {
		return nil, false
	}
	U.Userroomcard = ucInfo

	uextInfo, ueok := model.UserextrainfoOp.Get(UserID)
	if !ueok {
		return nil, false
	}
	U.Userextrainfo = uextInfo
	return U, true
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

	U.Gamescorelocker = &model.Gamescorelocker{
		UserID:UserID,
	}
	_, err = model.GamescorelockerOp.Insert(U.Gamescorelocker)
	if err != nil {
		log.Error("at createUser insert Gamescorelocker error")
		return nil, false
	}

	U.Gamescoreinfo = &model.Gamescoreinfo{
		UserID:UserID,
	}
	_, err = model.GamescoreinfoOp.Insert(U.Gamescoreinfo)
	if err != nil {
		log.Error("at createUser insert Gamescoreinfo error")
		return nil, false
	}

	U.Userroomcard = &model.Userroomcard{
		UserID:UserID,
	}
	_, err = model.UserroomcardOp.Insert(U.Userroomcard)
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

func BuildClientMsg(retMsg *msg.CMD_MB_LogonSuccess, user *user.User){
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
