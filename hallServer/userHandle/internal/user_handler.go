package internal

import (
	"fmt"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/hallServer/common"
	"mj/hallServer/conf"
	"mj/hallServer/db/model"
	"mj/hallServer/db/model/base"
	"mj/hallServer/game_list"
	"mj/hallServer/id_generate"
	"mj/hallServer/user"
	"time"

	"encoding/json"

	"mj/common/register"

	"mj/hallServer/match_room"

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

	//c2s
	reg.RegisterC2S(&msg.C2L_Login{}, m.handleMBLogin)
	reg.RegisterC2S(&msg.C2L_Regist{}, m.handleMBRegist)
	reg.RegisterC2S(&msg.C2L_User_Individual{}, m.GetUserIndividual)
	reg.RegisterC2S(&msg.C2L_CreateTable{}, m.CreateRoom)
	reg.RegisterC2S(&msg.C2L_ReqCreatorRoomRecord{}, m.GetCreatorRecord)
	reg.RegisterC2S(&msg.C2L_ReqRoomPlayerBrief{}, m.GetRoomPlayerBreif)
	reg.RegisterC2S(&msg.C2L_DrawSahreAward{}, m.DrawSahreAward)
	reg.RegisterC2S(&msg.C2L_SetElect{}, m.SetElect)
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
	u, ok := agent.UserData().(*user.User)
	if !ok {
		log.Error("at CloseAgent not foud user")
		return nil
	}
	DelUser(u.Id)
	m.Close(common.UserOffline)
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
		retcode = NotFoudAccout
		return
	}

	if _, ok := Users[accountData.UserID]; ok {
		retcode = ErrUserDoubleLogin
		return
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
			player.KindID = 0
			player.ServerID = 0
			player.GameNodeID = 0
			player.EnterIP = ""
			player.Roomid = 0
			model.GamescorelockerOp.UpdateWithMap(player.Id, map[string]interface{}{
				"KindID":     0,
				"ServerID":   0,
				"GameNodeID": 0,
				"EnterIP":    "",
				"roomid":     0,
			})
		}
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
	if accountData != nil {
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
	accInfo.UserID = int64(lastid)

	player, cok := createUser(accInfo.UserID, accInfo)
	if !cok {
		retcode = CreateUserError
		return
	}

	player.HallNodeID = conf.Server.NodeId
	model.GamescorelockerOp.UpdateWithMap(player.Id, map[string]interface{}{
		"HallNodeID": conf.Server.NodeId,
	})
	player.Agent = agent
	//AddUser(user.Id, user)
	agent.SetUserData(player)
	BuildClientMsg(retMsg, player, accInfo)
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

	monrey := feeTemp.TableFee
	if recvMsg.PayType == AA_PAY_TYPE {
		monrey = feeTemp.TableFee / template.MaxPlayer
	}

	if !player.EnoughCurrency(monrey) {
		retCode = NotEnoughFee
		return
	}

	//记录创建房间信息
	info := &model.CreateRoomInfo{}
	info.UserId = player.Id
	info.PayType = recvMsg.PayType
	info.MaxPlayerCnt = template.MaxPlayer
	info.RoomId = rid
	info.NodeId = nodeId
	info.Num = template.MaxPlayer
	info.KindId = recvMsg.Kind
	info.ServiceId = recvMsg.ServerId
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

	monrey := feeTemp.TableFee
	if roomInfo.PayType == AA_PAY_TYPE {
		monrey = feeTemp.TableFee / roomInfo.MaxPlayerCnt
	}

	if !player.SubCurrency(feeTemp.TableFee) {
		retcode = NotEnoughFee
		return
	}

	record := &model.TokenRecord{}
	record.UserId = player.Id
	record.RoomId = roomInfo.RoomID
	record.Amount = monrey
	record.TokenType = AA_PAY_TYPE
	record.KindID = template.KindID
	if !player.AddRecord(record) {
		retcode = ErrServerError
		player.AddCurrency(monrey)
		return
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
	retMsg.RoomCard = user.Currency
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
