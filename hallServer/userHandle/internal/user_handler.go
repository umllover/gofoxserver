package internal

import (
	"errors"
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/hallServer/db/model"
	"mj/hallServer/db/model/base"
	"mj/hallServer/gameList"
	"mj/hallServer/user"
	"reflect"
	"time"

	"mj/hallServer/common"

	"mj/hallServer/conf"

	"mj/hallServer/idGenerate"

	"github.com/lovelly/leaf/gate"
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
	m.ChanRPC.Register("GetUser", m.GetUser)
	m.ChanRPC.Register("SrarchTableResult", m.SrarchTableResult)

	//c2s
	handlerC2S(m, &msg.C2L_Login{}, m.handleMBLogin)
	handlerC2S(m, &msg.C2L_Regist{}, m.handleMBRegist)
	handlerC2S(m, &msg.C2L_User_Individual{}, m.GetUserIndividual)

	handlerC2S(m, &msg.C2L_CreateTable{}, m.CreateRoom)
}

//连接进来的通知
func (m *UserModule) NewAgent(args []interface{}) error {
	log.Debug("at hall NewAgent")
	return nil
}

//连接关闭的同喜
func (m *UserModule) CloseAgent(args []interface{}) error {
	log.Debug("at hall CloseAgent")
	agent := args[0].(gate.Agent)
	u, ok := agent.UserData().(*user.User)
	if !ok {
		return nil
	}
	DelUser(u.Id)
	m.Close(common.UserOffline)
	return nil
}

func (m *UserModule) handleMBLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_Login)
	retMsg := &msg.L2C_LogonSuccess{}
	agent := m.a
	retcode := 0
	defer func() {
		if retcode != 0 {
			str := fmt.Sprintf("登录失败, 错误码: %d", retcode)
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
		retcode = ErrUserDoubleLogin
		return
	}

	//if accountData.PasswordID != recvMsg.Password {
	//	sendErrFunc("password is error")
	//	return
	//}

	user := user.NewUser(accountData.UserID)
	user.Id = accountData.UserID
	lok := loadUser(user)
	if !lok {
		retcode = LoadUserInfoError
		return
	}

	user.Agent = agent
	AddUser(user.Id, user)
	agent.SetUserData(user)
	BuildClientMsg(retMsg, user, accountData)
	gameList.ChanRPC.Go("sendGameList", agent)
}

func (m *UserModule) handleMBRegist(args []interface{}) {
	retcode := 0
	recvMsg := args[0].(*msg.C2L_Regist)
	agent := args[1].(gate.Agent)
	retMsg := &msg.L2C_LogonSuccess{}
	var accountData *model.Accountsinfo
	defer func() {
		if retcode != 0 {
			model.AccountsinfoOp.DeleteByMap(map[string]interface{}{
				"Accounts": recvMsg.Accounts,
			})
			if accountData != nil {
				model.AccountsmemberOp.Delete(accountData.UserID)
				model.GamescorelockerOp.Delete(accountData.UserID)
				model.GamescoreinfoOp.Delete(accountData.UserID)
				model.UserattrOp.Delete(accountData.UserID)
				model.UserextrainfoOp.Delete(accountData.UserID)
				model.UsertokenOp.Delete(accountData.UserID)
			}
			agent.WriteMsg(&msg.L2C_LogonFailure{ResultCode: retcode, DescribeString: "登录失败"})
		} else {
			agent.WriteMsg(retMsg)
		}
	}()

	var ok error
	accountData, ok = model.AccountsinfoOp.GetByMap(map[string]interface{}{
		"Accounts": recvMsg.Accounts,
	})
	if ok != nil || accountData != nil {
		log.Debug("errpr == %v", ok)
		retcode = AlreadyExistsAccount
		return
	}

	//todo 名字排重等等等 验证
	now := time.Now()
	accInfo := &model.Accountsinfo{
		Gender:           recvMsg.Gender,   //用户性别
		Accounts:         recvMsg.Accounts, //登录帐号
		LogonPass:        recvMsg.LogonPass,
		InsurePass:       recvMsg.InsurePass,
		NickName:         recvMsg.NickName, //用户昵称
		GameLogonTimes:   1,
		LastLogonIP:      agent.RemoteAddr().String(),
		LastLogonMobile:  recvMsg.MobilePhone,
		LastLogonMachine: recvMsg.MachineID,
		RegisterMobile:   recvMsg.MobilePhone,
		RegisterMachine:  recvMsg.MachineID,
		RegisterDate:     &now,
		RegisterIP:       agent.RemoteAddr().String(), //连接地址
	}

	lastid, err := model.AccountsinfoOp.Insert(accInfo)
	if err != nil {
		retcode = InsertAccountError
		return
	}
	accInfo.UserID = int(lastid)

	user, cok := createUser(accInfo.UserID, accInfo)
	if !cok {
		retcode = CreateUserError
		return
	}
	user.Agent = agent
	AddUser(user.Id, user)
	agent.SetUserData(user)
	BuildClientMsg(retMsg, user, accInfo)
}

func (m *UserModule) GetUserIndividual(args []interface{}) {
	agent := args[1].(gate.Agent)
	user, ok := agent.UserData().(*user.User)
	if !ok {
		log.Debug("not foud user data")
		return
	}
	retmsg := &msg.L2C_UserIndividual{
		UserID:      user.Id,        //用户 I D
		NickName:    user.NickName,  //昵称
		WinCount:    user.WinCount,  //赢数
		LostCount:   user.LostCount, //输数
		DrawCount:   user.DrawCount, //平数
		Medal:       user.UserMedal,
		RoomCard:    user.RoomCard,    //房卡
		MemberOrder: user.MemberOrder, //会员等级
		Score:       user.Score,
		HeadImgUrl:  user.HeadImgUrl,
	}

	user.WriteMsg(retmsg)
}

func (m *UserModule) UserOffline() {

}

func (m *UserModule) CreateRoom(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_CreateTable)
	retMsg := &msg.L2C_CreateTableSucess{}
	agent := args[1].(gate.Agent)
	retCode := 0
	defer func() {
		if retCode == 0 {
			agent.WriteMsg(retMsg)
		} else {
			agent.WriteMsg(&msg.L2C_CreateTableFailure{ErrorCode: retCode, DescribeString: "创建房间失败"})
		}
	}()
	template, ok := base.GameServiceOptionCache.Get(recvMsg.Kind, recvMsg.ServerId)
	if !ok {
		retCode = NoFoudTemplate
		return
	}

	feeTemp, ok1 := base.PersonalTableFeeCache.Get(recvMsg.ServerId, recvMsg.Kind, recvMsg.DrawCountLimit)
	if !ok1 {
		log.Error("not foud PersonalTableFeeCache")
		retCode = NoFoudTemplate
		return
	}

	u := agent.UserData().(*user.User)
	if u.GetRoomCnt() > common.GetGlobalVarInt(MAX_CREATOR_ROOM_CNT) {
		retCode = ErrMaxRoomCnt
		return
	}

	rid, iok := idGenerate.GetRoomId(u.Id)
	if !iok {
		retCode = RandRoomIdError
		return
	}

	if recvMsg.JoinGamePeopleCount != 0 {
		if recvMsg.JoinGamePeopleCount > template.MaxPlayer || recvMsg.JoinGamePeopleCount < template.MaxPlayer {
			retCode = ErrParamError
			return
		}
	}

	host := gameList.GetSvrByKind(recvMsg.Kind)
	if host == "" {
		retCode = ErrNotFoudServer
		return
	}

	if recvMsg.PayType == SELF_PAY_TYPE {
		if u.SubCurrency(feeTemp.TableFee) {
			retCode = NotEnoughFee
			return
		}

		record := &model.TokenRecord{}
		record.UserId = u.Id
		record.Amount = feeTemp.TableFee
		record.TokenType = SELF_PAY_TYPE
		record.KindID = template.KindID
		_, err := model.TokenRecordOp.Insert(record)
		if err != nil {
			retCode = ErrServerError
			u.AddCurrency(feeTemp.TableFee)
			return
		}
	}

	//记录创建房间信息
	info := &model.CreateRoomInfo{}
	info.UserId = u.Id
	info.PayType = recvMsg.PayType
	info.MaxPlayerCnt = recvMsg.JoinGamePeopleCount
	info.RoomId = rid
	info.Num = recvMsg.DrawCountLimit

	//回给客户端的消息
	retMsg.TableID = rid
	retMsg.DrawCountLimit = info.Num
	retMsg.DrawTimeLimit = 0
	retMsg.Beans = feeTemp.TableFee
	retMsg.RoomCard = u.Currency
	retMsg.ServerIP = host
}

func (m *UserModule) SrarchTableResult(args []interface{}) {
	roomInfo := args[0].(*msg.RoomInfo)
	u := m.a.UserData().(*user.User)
	retMsg := &msg.L2C_SearchResult{}
	retcode := 0
	defer func() {
		if retcode != 0 {
			u.WriteMsg(RenderErrorMessage(retcode))
		} else {
			u.WriteMsg(retMsg)
		}
	}()

	template, ok := base.GameServiceOptionCache.Get(roomInfo.KindID, roomInfo.ServerID)
	if !ok {
		retcode = ConfigError
		return
	}

	feeTemp, ok1 := base.PersonalTableFeeCache.Get(roomInfo.ServerID, roomInfo.KindID, roomInfo.PayCnt)
	if !ok1 {
		log.Error("not foud PersonalTableFeeCache")
		retcode = NoFoudTemplate
		return
	}

	monrey := feeTemp.TableFee
	if roomInfo.PayType == AA_PAY_TYPE {
		monrey = feeTemp.TableFee / roomInfo.MaxCnt

		if u.SubCurrency(monrey) {
			retcode = NotEnoughFee
			return
		}
		record := &model.TokenRecord{}
		record.UserId = u.Id
		record.Amount = monrey
		record.TokenType = AA_PAY_TYPE
		record.KindID = template.KindID
		_, err := model.TokenRecordOp.Insert(record)
		if err != nil {
			retcode = ErrServerError
			u.AddCurrency(monrey)
			return
		}
	}

	retMsg.TableID = roomInfo.RoomID
	u.WriteMsg(retMsg)
}

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
		log.Debug("at loadUser not foud GamescorelockerOp by user  %d", u.Id)
		glInfo = &model.Gamescorelocker{}
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

	rooms, err := model.CreateRoomInfoOp.QueryByMap(map[string]interface{}{
		"user_id": u.Id,
	})
	if err != nil {
		log.Error("at loadUser not foud CreateRoomInfoOp by user  %d", u.Id)
		return false
	}
	for _, v := range rooms {
		u.AddRooms(v.RoomId, v)
	}
	return true
}

func createUser(UserID int, accountData *model.Accountsinfo) (*user.User, bool) {
	U := user.NewUser(UserID)
	U.Accountsmember = &model.Accountsmember{
		UserID: UserID,
	}
	_, err := model.AccountsmemberOp.Insert(U.Accountsmember)
	if err != nil {
		log.Error("at createUser insert Accountsmember error")
		return nil, false
	}

	now := time.Now()
	U.Gamescoreinfo = &model.Gamescoreinfo{
		UserID:        UserID,
		LastLogonDate: &now,
	}
	_, err = model.GamescoreinfoOp.Insert(U.Gamescoreinfo)
	if err != nil {
		log.Error("at createUser insert Gamescoreinfo error")
		return nil, false
	}

	U.Userattr = &model.Userattr{
		UserID:     UserID,
		NickName:   accountData.NickName,
		Gender:     accountData.Gender,
		HeadImgUrl: accountData.HeadImgUrl,
	}
	_, err = model.UserattrOp.Insert(U.Userattr)
	if err != nil {
		log.Error("at createUser insert Userroomcard error")
		return nil, false
	}

	U.Userextrainfo = &model.Userextrainfo{
		UserId: UserID,
	}
	_, err = model.UserextrainfoOp.Insert(U.Userextrainfo)
	if err != nil {
		log.Error("at createUser insert Userroomcard error")
		return nil, false
	}

	U.Usertoken = &model.Usertoken{
		UserID: UserID,
	}

	_, err = model.UsertokenOp.Insert(U.Usertoken)
	if err != nil {
		log.Error("at createUser insert Userroomcard error")
		return nil, false
	}

	return U, true
}

func BuildClientMsg(retMsg *msg.L2C_LogonSuccess, user *user.User, acinfo *model.Accountsinfo) {
	retMsg.FaceID = user.FaceID //头像标识
	retMsg.Gender = user.Gender
	retMsg.UserID = user.Id
	retMsg.Spreader = acinfo.SpreaderID
	retMsg.Experience = user.Experience
	retMsg.LoveLiness = user.LoveLiness
	retMsg.NickName = user.NickName

	//用户成绩
	retMsg.UserScore = user.Score
	retMsg.UserInsure = user.InsureScore
	retMsg.Medal = user.UserMedal
	retMsg.UnderWrite = user.UnderWrite
	retMsg.WinCount = user.WinCount
	retMsg.LostCount = user.LostCount
	retMsg.DrawCount = user.DrawCount
	retMsg.FleeCount = user.FleeCount
	log.Debug("node id === %v", conf.Server.NodeId)
	retMsg.HallNodeID = conf.Server.NodeId
	tm := &msg.DateTime{}
	tm.Year = acinfo.RegisterDate.Year()
	tm.DayOfWeek = int(acinfo.RegisterDate.Weekday())
	tm.Day = acinfo.RegisterDate.Day()
	tm.Hour = acinfo.RegisterDate.Hour()
	tm.Second = acinfo.RegisterDate.Second()
	tm.Minute = acinfo.RegisterDate.Minute()
	retMsg.RegisterDate = tm
	//额外信息
	retMsg.MbTicket = user.MbTicket
	retMsg.MbPayTotal = user.MbPayTotal
	retMsg.MbVipLevel = user.MbVipLevel
	retMsg.PayMbVipUpgrade = user.PayMbVipUpgrade

	//约战房相关
	retMsg.RoomCard = user.RoomCard
	retMsg.LockServerID = user.ServerID
	retMsg.KindID = user.KindID
	retMsg.LockServerID = user.ServerID
}

/////////////////////////////// help 函数

func (m *UserModule) GetUser(args []interface{}) (interface{}, error) {
	u, ok := m.a.UserData().(*user.User)
	if !ok {
		return nil, errors.New("not foud user Data at GetUser")
	}
	return u, nil
}
