package internal

import (
	"encoding/json"
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/register"
	"mj/hallServer/common"
	"mj/hallServer/conf"
	"mj/hallServer/db/model"
	"mj/hallServer/db/model/base"
	"mj/hallServer/game_list"
	"mj/hallServer/id_generate"
	"mj/hallServer/match_room"
	"mj/hallServer/user"
	"time"

	"mj/hallServer/center"

	"mj/common/utils"

	"github.com/lovelly/leaf/gate"
	"github.com/lovelly/leaf/log"
)

func RegisterHandler(m *UserModule) {
	reg := register.NewRegister(m.ChanRPC)
	//注册rpc 消息
	reg.RegisterRpc("handleMsgData", m.handleMsgData)
	reg.RegisterRpc("NewAgent", m.NewAgent)
	reg.RegisterRpc("CloseAgent", m.CloseAgent)
	reg.RegisterRpc("GetUser", m.GetUser)
	reg.RegisterRpc("SrarchTableResult", m.SrarchTableResult)
	reg.RegisterRpc("RoomCloseInfo", m.RoomCloseInfo)
	reg.RegisterRpc("restoreToken", m.restoreToken)
	reg.RegisterRpc("matchResult", m.matchResult)
	reg.RegisterRpc("LeaveRoom", m.leaveRoom)
	reg.RegisterRpc("JoinRoom", m.joinRoom)
	reg.RegisterRpc("Recharge", m.Recharge)
	reg.RegisterRpc("S2S_RenewalFeeFaild", m.RenewalFeeFaild)
	reg.RegisterRpc("S2S_OfflineHandler", m.HandlerOffilneEvent)
	reg.RegisterRpc("ForceClose", m.ForceClose)
	//c2s
	reg.RegisterC2S(&msg.C2L_Login{}, m.handleMBLogin)
	reg.RegisterC2S(&msg.C2L_Regist{}, m.handleMBRegist)
	reg.RegisterC2S(&msg.C2L_User_Individual{}, m.GetUserIndividual)
	reg.RegisterC2S(&msg.C2L_CreateTable{}, m.CreateRoom)
	reg.RegisterC2S(&msg.C2L_ReqCreatorRoomRecord{}, m.GetCreatorRecord)
	reg.RegisterC2S(&msg.C2L_ReqRoomPlayerBrief{}, m.GetRoomPlayerBreif)
	reg.RegisterC2S(&msg.C2L_DrawSahreAward{}, m.DrawSahreAward)
	reg.RegisterC2S(&msg.C2L_SetElect{}, m.SetElect)
	reg.RegisterC2S(&msg.C2L_DeleteRoom{}, m.DeleteRoom)
	reg.RegisterC2S(&msg.C2L_SetPhoneNumber{}, m.SetPhoneNumber)
	reg.RegisterC2S(&msg.C2L_DianZhan{}, m.DianZhan)
	reg.RegisterC2S(&msg.C2L_RenewalFees{}, m.RenewalFees)
	reg.RegisterC2S(&msg.C2L_ChangeUserName{}, m.ChangeUserName)
	reg.RegisterC2S(&msg.C2L_ChangeSign{}, m.ChangeSign)
	reg.RegisterC2S(&msg.C2L_ReqBindMaskCode{}, m.ReqBindMaskCode)
	reg.RegisterC2S(&msg.C2L_RechangerOk{}, m.RechangerOk)
	reg.RegisterRpc("RoomEndInfo", m.RoomEndInfo)
}

//连接进来的通知
func (m *UserModule) NewAgent(args []interface{}) error {
	log.Debug("at hall NewAgent")
	return nil
}

//连接关闭的通知
func (m *UserModule) CloseAgent(args []interface{}) error {
	log.Debug("at hall CloseAgent")
	agent := args[0].(gate.Agent)
	Reason := args[1].(int)
	player, ok := agent.UserData().(*user.User)
	if !ok || player == nil {
		log.Error("at CloseAgent not foud user")
		return nil
	}

	m.UserOffline()
	if Reason != KickOutMsg { //重登踢出会覆盖， 所以这里不用删除
		DelUser(player.Id)
	}

	if Reason == UserOffline {
		m.Close(UserOffline)
	}
	log.Debug("CloseAgent ok")
	return nil
}

func (m *UserModule) handleMBLogin(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_Login)
	retMsg := &msg.L2C_LogonSuccess{}
	agent := m.a
	retcode := 0

	log.Debug("enter mbLogin  user:%s", recvMsg.Accounts)

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
		if conf.Test {
			retcode, _, accountData = RegistUser(&msg.C2L_Regist{
				LogonPass:    recvMsg.LogonPass,
				Accounts:     recvMsg.Accounts,
				ModuleID:     recvMsg.ModuleID,
				PlazaVersion: recvMsg.PlazaVersion,
				MachineID:    recvMsg.MachineID,
				MobilePhone:  recvMsg.MobilePhone,
				NickName:     recvMsg.Accounts,
			}, agent)
			if retcode != 0 {
				return
			}
		} else {
			retcode = NotFoudAccout
			return
		}
	}

	if accountData.LogonPass != recvMsg.LogonPass {
		retcode = ErrPasswd
		return
	}

	player := user.NewUser(accountData.UserID)
	player.Id = accountData.UserID
	lok := loadUser(player)
	if !lok {
		retcode = LoadUserInfoError
		return
	}

	if player.Roomid != 0 {
		_, have := game_list.ChanRPC.Call1("HaseRoom", player.Roomid)
		if have != nil {
			log.Debug("user :%d room %d is close ", player.Id, player.Roomid)
			player.DelGameLockInfo()
		}
	}

	oldUser := getUser(accountData.UserID)
	if oldUser != nil {
		log.Debug("old user ====== %d  %d ", oldUser.KindID, oldUser.Roomid)
		m.KickOutUser(oldUser)
	}

	player.Agent = agent
	AddUser(player.Id, player)
	agent.SetUserData(player)
	player.LoadTimes()
	player.HallNodeID = conf.Server.NodeId
	model.GamescorelockerOp.UpdateWithMap(player.Id, map[string]interface{}{
		"HallNodeID": conf.Server.NodeId,
	})
	BuildClientMsg(retMsg, player, accountData)
	game_list.ChanRPC.Go("sendGameList", agent)

	m.Recharge(nil)
}

func (m *UserModule) handleMBRegist(args []interface{}) {
	retcode := 0
	recvMsg := args[0].(*msg.C2L_Regist)
	agent := args[1].(gate.Agent)
	retMsg := &msg.L2C_RegistResult{}
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
			agent.WriteMsg(RenderErrorMessage(retcode, "注册失败"))
		} else {
			agent.WriteMsg(retMsg)
		}
	}()

	retcode, _, _ = RegistUser(recvMsg, agent)
}

func RegistUser(recvMsg *msg.C2L_Regist, agent gate.Agent) (int, *user.User, *model.Accountsinfo) {
	accountData, _ := model.AccountsinfoOp.GetByMap(map[string]interface{}{
		"Accounts": recvMsg.Accounts,
	})
	if accountData != nil {
		return AlreadyExistsAccount, nil, nil
	}

	//todo 名字排重等等等 验证
	now := time.Now()
	accInfo := &model.Accountsinfo{
		UserID:           user.GetUUID(),
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

	_, err := model.AccountsinfoOp.Insert(accInfo)
	if err != nil {
		log.Error("RegistUser err :%s", err.Error())
		return InsertAccountError, nil, nil
	}

	player, cok := createUser(accInfo.UserID, accInfo)
	if !cok {
		return CreateUserError, nil, nil
	}
	return 0, player, accInfo
}

//获取个人信息
func (m *UserModule) GetUserIndividual(args []interface{}) {
	agent := args[1].(gate.Agent)
	player, ok := agent.UserData().(*user.User)
	if !ok {
		log.Debug("not foud user data")
		return
	}
	retmsg := &msg.L2C_UserIndividual{
		UserID:      player.Id,        //用户 I D
		NickName:    player.NickName,  //昵称
		WinCount:    player.WinCount,  //赢数
		LostCount:   player.LostCount, //输数
		DrawCount:   player.DrawCount, //平数
		Medal:       player.UserMedal,
		RoomCard:    player.Currency,    //房卡
		MemberOrder: player.MemberOrder, //会员等级
		Score:       player.Score,
		HeadImgUrl:  player.HeadImgUrl,
	}

	player.WriteMsg(retmsg)
	player.SendActivityInfo()
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

	feeTemp, ok1 := base.PersonalTableFeeCache.Get(recvMsg.Kind, recvMsg.ServerId, recvMsg.DrawCountLimit)
	if !ok1 {
		log.Error("not foud PersonalTableFeeCache")
		retCode = NoFoudTemplate
		return
	}

	player := agent.UserData().(*user.User)
	if player.GetRoomCnt() >= common.GetGlobalVarInt(MAX_CREATOR_ROOM_CNT) {
		retCode = ErrMaxRoomCnt
		return
	}

	host, nodeId := game_list.GetSvrByKind(recvMsg.Kind)
	if host == "" {
		retCode = ErrNotFoudServer
		return
	}

	rid, iok := id_generate.GenerateRoomId(nodeId)
	if !iok {
		retCode = RandRoomIdError
		return
	}

	//检测是否有限时免费
	//if !player.CheckFree() {
	money := feeTemp.TableFee
	if recvMsg.PayType == AA_PAY_TYPE {
		money = feeTemp.TableFee / template.MaxPlayer
	}
	if !player.EnoughCurrency(money) {
		retCode = NotEnoughFee
		return
	}
	//}

	//记录创建房间信息
	info := &model.CreateRoomInfo{}
	info.UserId = player.Id
	info.PayType = recvMsg.PayType
	info.MaxPlayerCnt = template.MaxPlayer
	info.RoomId = rid
	info.NodeId = nodeId
	info.Num = recvMsg.DrawCountLimit
	info.KindId = recvMsg.Kind
	info.ServiceId = recvMsg.ServerId
	now := time.Now()
	info.CreateTime = &now
	if recvMsg.Public {
		info.Public = 1
	} else {
		info.Public = 0
	}

	by, err := json.Marshal(recvMsg.OtherInfo)
	if err != nil {
		log.Error("at CreateRoom json.Marshal(recvMsg.OtherInfo) error:%s", err.Error())
		retCode = ErrParamError
		return
	}
	info.OtherInfo = string(by)
	if recvMsg.RoomName != "" {
		info.RoomName = recvMsg.RoomName
	} else {
		info.RoomName = template.RoomName
	}

	player.AddRooms(info)

	roomInfo := &msg.RoomInfo{}
	roomInfo.KindID = info.KindId
	roomInfo.ServerID = info.ServiceId
	roomInfo.RoomID = info.RoomId
	roomInfo.NodeID = info.NodeId
	roomInfo.SvrHost = host
	roomInfo.PayType = info.PayType
	roomInfo.CreateTime = time.Now().Unix()
	roomInfo.CreateUserId = player.Id
	roomInfo.IsPublic = recvMsg.Public
	roomInfo.MachPlayer = make(map[int64]struct{})
	roomInfo.Players = make(map[int64]*msg.PlayerBrief)
	roomInfo.MaxPlayerCnt = info.MaxPlayerCnt
	roomInfo.PayCnt = info.Num
	roomInfo.RoomName = info.RoomName
	game_list.ChanRPC.Go("addyNewRoom", roomInfo)

	//回给客户端的消息
	retMsg.TableID = rid
	retMsg.DrawCountLimit = info.Num
	retMsg.DrawTimeLimit = 0
	retMsg.Beans = feeTemp.TableFee
	retMsg.RoomCard = player.Currency
	retMsg.ServerIP = host
}

func (m *UserModule) SrarchTableResult(args []interface{}) {
	roomInfo := args[0].(*msg.RoomInfo)
	player := m.a.UserData().(*user.User)
	retMsg := &msg.L2C_SearchResult{}
	retcode := 0
	defer func() {
		if retcode != 0 {
			if roomInfo.CreateUserId == player.Id {
				//todo  delte room ???
			}
			match_room.ChanRPC.Go("delMatchPlayer", player.Id, roomInfo)
			player.WriteMsg(RenderErrorMessage(retcode))
		} else {
			player.WriteMsg(retMsg)
		}
	}()

	template, ok := base.GameServiceOptionCache.Get(roomInfo.KindID, roomInfo.ServerID)
	if !ok {
		retcode = ConfigError
		return
	}

	feeTemp, ok1 := base.PersonalTableFeeCache.Get(roomInfo.KindID, roomInfo.ServerID, roomInfo.PayCnt)
	if !ok1 {
		log.Error("not foud PersonalTableFeeCache kindId:%d, serverID:%d, payCnt:%d", roomInfo.KindID, roomInfo.ServerID, roomInfo.PayCnt)
		retcode = NoFoudTemplate
		return
	}

	host := game_list.GetSvrByNodeID(roomInfo.NodeID)
	if host == "" {
		retcode = ErrNotFoudServer
		return
	}

	money := feeTemp.TableFee
	if roomInfo.PayType == AA_PAY_TYPE {
		money = feeTemp.AATableFee
	}

	//非限时免费 并且 不是全付方式 并且 钱大于零
	if !player.CheckFree() && roomInfo.PayType != SELF_PAY_TYPE && money > 0 {
		if !player.SubCurrency(money) {
			retcode = NotEnoughFee
			return
		}
	}

	if !player.HasRecord(roomInfo.RoomID) {
		record := &model.TokenRecord{}
		record.UserId = player.Id
		record.RoomId = roomInfo.RoomID
		record.Amount = money
		record.TokenType = AA_PAY_TYPE
		record.KindID = template.KindID
		if !player.AddRecord(record) {
			retcode = ErrServerError
			player.AddCurrency(money)
			return
		}
	} else { //已近口过钱了， 还来搜索房间
		log.Debug("player %d double srach room: %d", player.Id, roomInfo.RoomID)
	}

	player.KindID = roomInfo.KindID
	player.ServerID = roomInfo.ServerID
	player.Roomid = roomInfo.RoomID
	player.GameNodeID = roomInfo.NodeID
	player.EnterIP = host

	model.GamescorelockerOp.UpdateWithMap(player.Id, map[string]interface{}{
		"KindID":     player.KindID,
		"ServerID":   player.ServerID,
		"GameNodeID": roomInfo.NodeID,
		"EnterIP":    host,
		"roomid":     roomInfo.RoomID,
	})

	retMsg.TableID = roomInfo.RoomID
	retMsg.ServerIP = host
	return
}

//获取自己创建的房间
func (m *UserModule) GetCreatorRecord(args []interface{}) {
	//recvMsg := args[0].(*msg.C2L_ReqCreatorRoomRecord)
	retMsg := &msg.L2C_CreatorRoomRecord{}
	u := m.a.UserData().(*user.User)
	retMsg.Records = u.GetRoomInfo()
	u.WriteMsg(retMsg)
}

//获取某个房间内的玩家信息
func (m *UserModule) GetRoomPlayerBreif(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_ReqRoomPlayerBrief)
	u := m.a.UserData().(*user.User)
	r := u.GetRoom(recvMsg.RoomId)
	if r == nil {
		u.WriteMsg(&msg.L2C_RoomPlayerBrief{})
	} else {
		game_list.ChanRPC.Go("SendPlayerBrief", recvMsg.RoomId, u)
	}
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
		log.Error("at loadUser not foud GamescoreinfoOp by user  %d", u.Id)
		//return false
		glInfo = &model.Gamescorelocker{
			UserID: u.Id,
		}
		model.GamescorelockerOp.Insert(glInfo)
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
		log.Debug("load creator room info ok ... %v", v)
		u.Rooms[v.RoomId] = v
	}

	tokenRecords, terr := model.TokenRecordOp.QueryByMap(map[string]interface{}{
		"user_id": u.Id,
	})
	if terr != nil {
		log.Error("at loadUser not foud CreateRoomInfoOp by user  %d", u.Id)
		return false
	}

	//加载扣钱记录
	now := time.Now().Unix()
	for _, v := range tokenRecords {
		temp, ok := base.GameServiceOptionCache.Get(v.KindID, v.ServerId)
		if ok {
			if v.CreatorTime.Unix()+int64(temp.TimeNotBeginGame) < now && v.Status == 0 { //没开始返回钱
				if u.AddCurrency(v.Amount) {
					model.TokenRecordOp.Delete(v.RoomId, v.UserId)
				}
			}

			if v.CreatorTime.Unix()+86400 < now { // 一天了还没删除？？？ 在这里删除。安全处理
				model.TokenRecordOp.Delete(v.RoomId, v.UserId)
			}
		}
		u.Records[v.RoomId] = v
	}
	return true
}

func createUser(UserID int64, accountData *model.Accountsinfo) (*user.User, bool) {
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

	if conf.Test {
		U.Usertoken.Currency = 9999999
		U.Usertoken.RoomCard = 1000000
	}

	U.Gamescorelocker = &model.Gamescorelocker{
		UserID: UserID,
	}
	_, err = model.GamescorelockerOp.Insert(U.Gamescorelocker)
	if err != nil {
		log.Error("at createUser insert Gamescorelocker error")
		return nil, false
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
	//retMsg.UserScore = user.Score
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
	//retMsg.RoomCard = user.Currency
	retMsg.UserScore = user.Currency
	retMsg.LockServerID = user.ServerID
	retMsg.KindID = user.KindID
	retMsg.LockServerID = user.ServerID
	retMsg.ServerIP = user.EnterIP
}

//房间结束了
func (m *UserModule) RoomCloseInfo(args []interface{}) {
	info := args[0].(*msg.RoomEndInfo)
	player := m.a.UserData().(*user.User)
	if info.Status == 0 { //没开始就结束
		record := player.GetRecord(info.RoomId)
		if record != nil { //还原扣的钱
			err := player.DelRecord(record.RoomId)
			if err == nil {
				player.AddCurrency(record.Amount)
			} else {
				log.Error("at restoreToken not DelRecord error uid:%d", player.Id)
			}

		} else {
			log.Error("at restoreToken not foud record uid:%d", player.Id)
		}
	}
	player.DelRooms(info.RoomId)
	player.DelGameLockInfo()
	return
}

//离开房间还原
func (m *UserModule) restoreToken(args []interface{}) {
	player := m.a.UserData().(*user.User)
	RoomId := args[0].(int)
	record := player.GetRecord(RoomId)
	if record != nil { //还原扣的钱
		err := player.DelRecord(record.RoomId)
		if err == nil {
			player.AddCurrency(record.Amount)
		} else {
			log.Error("at restoreToken not DelRecord error uid:%d", player.Id)
		}

	} else {
		log.Error("at restoreToken not foud record uid:%d", player.Id)
	}
}

func (m *UserModule) matchResult(args []interface{}) {
	ret := args[0].(bool)
	retMsg := &msg.L2C_SearchResult{}
	u := m.a.UserData().(*user.User)
	if ret {
		r := args[1].(*msg.RoomInfo)
		retMsg.TableID = r.RoomID
		retMsg.ServerIP = r.SvrHost
	} else {
		retMsg.TableID = INVALID_TABLE
	}
	u.WriteMsg(retMsg)
}

func (m *UserModule) leaveRoom(args []interface{}) {
	u := m.a.UserData().(*user.User)
	log.Debug("at hall server leaveRoom uid:%v", u.Id)
}

func (m *UserModule) joinRoom(args []interface{}) {
	room := args[0].(*msg.RoomInfo)
	u := m.a.UserData().(*user.User)
	log.Debug("at hall server joinRoom uid:%v", u.Id)
	u.KindID = room.KindID
	u.ServerID = room.ServerID
	u.GameNodeID = room.NodeID
	u.EnterIP = room.SvrHost
}

func (m *UserModule) Recharge(args []interface{}) {
	u := m.a.UserData().(*user.User)
	orders := GetOrders(u.Id)
	for _, v := range orders {
		goods, ok := base.GoodsCache.Get(v.GoodsID)
		if !ok {
			log.Error("at Recharge error")
			continue
		}

		if UpdateOrderStats(v.OnLineID) {
			u.AddCurrency(goods.Diamond)
		}
	}
}

//离线通知时间
func (m *UserModule) HandlerOffilneEvent(args []interface{}) {
	recvMsg := args[0].(*msg.S2S_OfflineHandler)
	player := m.a.UserData().(*user.User)
	h, ok := model.UserOfflineHandlerOp.Get(recvMsg.EventID)
	if ok {
		handlerEventFunc(player, h)
	}
}

func (m *UserModule) KickOutUser(player *user.User) {
	player.ChanRPC().Go("ForceClose")
}

func (m *UserModule) ForceClose(args []interface{}) {
	log.Debug("at ForceClose ..... ")
	m.Close(KickOutMsg)
}

//删除自己创建的房间
func (m *UserModule) DeleteRoom(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_DeleteRoom)
	player := m.a.UserData().(*user.User)

	info := player.GetRoom(recvMsg.RoomId)
	if info != nil {
		player.WriteMsg(&msg.L2C_DeleteRoomResult{Code: ErrNotFondCreatorRoom})
		return
	}

	center.AsynCallGame(info.NodeId, m.Skeleton.GetChanAsynRet(), &msg.S2S_CloseRoom{RoomID: recvMsg.RoomId}, func(data interface{}, err error) {
		if err != nil {
			player.WriteMsg(&msg.L2C_DeleteRoomResult{Code: ErrRoomIsStart})
		} else {
			player.DelRooms(recvMsg.RoomId)
			player.WriteMsg(&msg.L2C_DeleteRoomResult{})
		}
	})

}

//绑定电话号码
func (m *UserModule) SetPhoneNumber(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_SetPhoneNumber)
	player := m.a.UserData().(*user.User)
	retCode := 0
	defer func() {
		player.WriteMsg(&msg.L2C_SetPhoneNumberRsp{Code: retCode})
	}()

	info, ok := model.UserMaskCodeOp.Get(player.Id)
	if !ok {
		retCode = ErrMaskCodeNotFoud
		return
	}

	if info.MaskCode != recvMsg.MaskCode {
		retCode = ErrMaskCodeError
		return
	}

	model.UserMaskCodeOp.Delete(player.Id)
	player.PhomeNumber = info.PhomeNumber
	model.UserattrOp.UpdateWithMap(player.Id, map[string]interface{}{
		"phome_number": info.PhomeNumber,
	})
}

//点赞
func (m *UserModule) DianZhan(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_DianZhan)
	//player := m.a.UserData().(*user.User)
	AddOfflineHandler(MailTypeDianZhan, recvMsg.UserID, nil)
}

//续费
func (m *UserModule) RenewalFees(args []interface{}) {
	//recvMsg := args[0].(*msg.C2L_RenewalFees)
	player := m.a.UserData().(*user.User)
	retCode := 0
	defer func() {
		player.WriteMsg(&msg.L2C_RenewalFeesRsp{Code: retCode})
	}()
	if player.Roomid == 0 {
		retCode = ErrNotInRoom
		return
	}

	info, err := game_list.ChanRPC.TimeOutCall1("GetRoomByRoomId", 5*time.Second, player.Roomid)
	if err != nil {
		retCode = ErrFindRoomError
		return
	}

	room := info.(*msg.RoomInfo)
	feeTemp, ok := base.PersonalTableFeeCache.Get(room.KindID, room.ServerID, room.PayCnt/room.RenewalCnt)
	if !ok {
		retCode = ErrConfigError
		return
	}

	monrey := feeTemp.TableFee
	if room.PayType == AA_PAY_TYPE {
		monrey = feeTemp.AATableFee
	}

	if !player.SubCurrency(feeTemp.TableFee) {
		retCode = NotEnoughFee
		return
	}

	if !player.HasRecord(room.RoomID) {
		record := &model.TokenRecord{}
		record.UserId = player.Id
		record.RoomId = room.RoomID
		record.Amount = monrey
		record.TokenType = AA_PAY_TYPE
		record.KindID = room.KindID
		if !player.AddRecord(record) {
			retCode = ErrServerError
			player.AddCurrency(monrey)
			return
		}
	} else { //已近口过钱了， 还来搜索房间
		log.Debug("player %d double srach room: %d", player.Id, room.RoomID)
	}
	room.PayCnt *= 2
	room.RenewalCnt++

	center.SendMsgToGame(room.NodeID, &msg.S2S_RenewalFee{RoomID: room.RoomID})
}

func (m *UserModule) RenewalFeeFaild(args []interface{}) {
	recvMsg := args[0].(*msg.S2S_RenewalFeeFaild)
	player := m.a.UserData().(*user.User)
	record := player.GetRecord(recvMsg.RecodeID)
	if record != nil {
		player.AddCurrency(record.Amount)
		player.DelRecord(record.RoomId)
	}
}

//改名字
func (m *UserModule) ChangeUserName(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_ChangeUserName)
	player := m.a.UserData().(*user.User)
	player.NickName = recvMsg.NewName

	model.UserattrOp.UpdateWithMap(player.Id, map[string]interface{}{
		"NickName": player.NickName,
	})

	player.WriteMsg(&msg.L2C_ChangeUserNameRsp{Code: 0, NewName: player.NickName})
}

//改签名
func (m *UserModule) ChangeSign(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_ChangeSign)
	player := m.a.UserData().(*user.User)

	player.Sign = recvMsg.Sign
	model.UserattrOp.UpdateWithMap(player.Id, map[string]interface{}{
		"Sign": player.Sign,
	})

	player.WriteMsg(&msg.L2C_ChangeSignRsp{Code: 0, NewSign: player.Sign})
}

//获取验证码
func (m *UserModule) ReqBindMaskCode(args []interface{}) {
	recvMsg := args[0].(*msg.C2L_ReqBindMaskCode)
	player := m.a.UserData().(*user.User)
	retCode := 0
	defer func() {
		player.WriteMsg(&msg.L2C_ReqBindMaskCodeRsp{Code: retCode})
	}()

	if player.MacKCodeTime != nil {
		if time.Now().After(*player.MacKCodeTime) {
			retCode = ErrFrequentAccess
			return
		}
	}
	code, _ := utils.RandInt(100000, 1000000)
	now := time.Now()
	player.MacKCodeTime = &now

	err := model.UserMaskCodeOp.InsertUpdate(&model.UserMaskCode{UserId: player.Id, PhomeNumber: recvMsg.PhoneNumber, MaskCode: code}, map[string]interface{}{
		"mask_code":    code,
		"phome_number": recvMsg.PhoneNumber,
	})
	if err == nil {
		retCode = ErrRandMaskCodeError
		return
	}

	ReqGetMaskCode(recvMsg.PhoneNumber, code)
}

func (m *UserModule) RechangerOk(args []interface{}) {
	//recvMsg := args[0].(*msg.C2L_RechangerOk)
	m.Recharge(nil)
}

/// 游戏服发来的结束消息
func (m *UserModule) RoomEndInfo(args []interface{}) {
}
