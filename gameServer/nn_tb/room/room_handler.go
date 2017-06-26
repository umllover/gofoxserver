package room

import (
	"math"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/common/msg/nn_tb_msg"
	"mj/gameServer/Chat"
	"mj/gameServer/db/model/base"
	client "mj/gameServer/user"
	"strconv"
	"time"

	. "mj/gameServer/common/mj_logic_base"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
	"mj/common/msg/mj_hz_msg"
)

func RegisterHandler(r *Room) {
	r.ChanRPC.Register("Sitdown", r.Sitdown)
	r.ChanRPC.Register("SetGameOption", r.SetGameOption)
	r.ChanRPC.Register("UserStandup", r.UserStandup)
	r.ChanRPC.Register("UserReady", r.UserReady)
	r.ChanRPC.Register("userOffline", r.UserOffline)
	r.ChanRPC.Register("userRelogin", r.UserReLogin)
	r.ChanRPC.Register("GetUserChairInfo", r.GetUserChairInfo)
	r.ChanRPC.Register("DissumeRoom", r.DissumeRoom)

}

func (room  *Room) CallScore(args []interface{})  {
	recvMsg := args[0].(*nn_tb_msg.C2G_TBNN_CallScore)
	user := args[1].(*client.User)
	retcode := 0

	log.Debug("Enter Room TBNN id=%d", room.GetRoomId())
	defer func() {
		if retcode != 0 {
			user.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	if recvMsg.CallScore >= 0 {
		room.CellScore = recvMsg.CallScore
	}	else {
		retcode = -1
		return
	}

}

func (room *Room) SetGameOption(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_GameOption)
	user := args[1].(*client.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			user.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	if user.ChairId == INVALID_CHAIR {
		retcode = ErrNoSitdowm
		return
	}

	//template, ok := base.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	_, ok := base.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	if !ok {
		retcode = ConfigError
		return
	}

	user.WriteMsg(&msg.G2C_GameStatus{
		GameStatus:  room.Status,
		AllowLookon: room.AllowLookon[user.ChairId],
	})

	if room.CreateUser == user.Id { //房主设置
		room.AllowLookon[user.ChairId] = recvMsg.AllowLookon
	}

	user.WriteMsg(&msg.G2C_PersonalTableTip{
		TableOwnerUserID:  room.CreateUser,                          //桌主 I D
		DrawCountLimit:    room.CountLimit,                          //局数限制
		DrawTimeLimit:     room.TimeLimit,                           //时间限制
		PlayCount:         room.PlayCount,                           //已玩局数
		PlayTime:          int(room.CreateTime - time.Now().Unix()), //已玩时间
		CellScore:         room.CellScore,                              //游戏底分
		IniScore:          room.InitCellScore,                           //初始分数
		ServerID:          strconv.Itoa(room.GetRoomId()),           //房间编号
		IsJoinGame:        0,                                        //是否参与游戏 todo  tagPersonalTableParameter
		IsGoldOrGameScore: room.IsGoldOrGameScore,                   //金币场还是积分场 0 标识 金币场 1 标识 积分场
	})

	/*if (template.ServerType & GAME_GENRE_PERSONAL) != 0 { //约战类型。。。
		user.WriteMsg(room.Record)
	}*/

	if room.Status == RoomStatusReady { // 没开始
		StatusFree := &msg.G2C_StatusFree{}
		StatusFree.CellScore = room.CellScore                //基础积分
		StatusFree.TimeOutCard = room.TimeOutCard         //出牌时间
		StatusFree.TimeOperateCard = room.TimeOperateCard //操作时间
		StatusFree.TimeStartGame = room.TimeStartGame     //开始时间
		StatusFree.PlayerCount = room.PlayCount           //玩家人数
		StatusFree.CountLimit = room.CountLimit           //局数限制
		user.WriteMsg(StatusFree)
	} else { //开始了
		StatusPlay := &msg.G2C_StatusPlay{}
		//自定规则
		StatusPlay.TimeOutCard = room.TimeOutCard
		StatusPlay.TimeOperateCard = room.TimeOperateCard
		StatusPlay.TimeStartGame = room.TimeStartGame

		room.OnUserTrustee(user.ChairId, false) //重入取消托管

		//规则
		StatusPlay.PlayerCount = int(room.PlayerCount)
		//游戏变量
		StatusPlay.BankerUser = room.BankerUser
		StatusPlay.CellScore = room.CellScore
		//StatusPlay.Trustee = room.Trustee

		//历史积分
		for j := 0; j < room.UserCnt; j++ {
			//设置变量
			StatusPlay.TurnScore[j] = room.HistoryScores[j].TurnScore
			StatusPlay.CollectScore[j] = room.HistoryScores[j].CollectScore
		}

		user.WriteMsg(StatusPlay)
	}
}

//起立
func (room *Room) UserStandup(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserStandup{})
	user := args[1].(*client.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			user.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	if room.Status == RoomStatusStarting {
		retcode = ErrGameIsStart
		return
	}

	room.setUsetStatus(user, US_FREE)
	room.LeaveRoom(user)
}

//坐下
func (room *Room) Sitdown(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_UserSitdown)
	user := args[1].(*client.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			user.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	oldUser := room.GetUserByChairId(recvMsg.ChairID)
	if oldUser != nil {
		retcode = ChairHasUser
		return
	}

	template, ok := base.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	if !ok {
		retcode = ConfigError
		return
	}

	if room.Status == RoomStatusStarting && template.DynamicJoin == 1 {
		retcode = GameIsStart
		return
	}

	if room.ChatRoomId == 0 {
		id, err := Chat.ChanRPC.Call1("createRoom", user.Agent)
		if err != nil {
			log.Error("create Chat Room faild")
			retcode = ErrCreateRoomFaild
		}

		room.ChatRoomId = id.(int)
	}

	_, chairId := room.GetUserByUid(user.Id)
	if chairId > 0 {
		room.LeaveRoom(user)
	}

	room.EnterRoom(recvMsg.ChairID, user)
	//把自己的信息推送给所有玩家
	room.SendMsgAllNoSelf(user.Id, &msg.G2C_UserEnter{
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

	//把所有玩家信息推送给自己
	room.ForEachUser(func(u *client.User) {
		if u.Id == user.Id {
			return
		}
		user.WriteMsg(&msg.G2C_UserEnter{
			UserID:      u.Id,          //用户 I D
			FaceID:      u.FaceID,      //头像索引
			CustomID:    u.CustomID,    //自定标识
			Gender:      u.Gender,      //用户性别
			MemberOrder: u.MemberOrder, //会员等级
			TableID:     u.RoomId,      //桌子索引
			ChairID:     u.ChairId,     //椅子索引
			UserStatus:  u.Status,      //用户状态
			Score:       u.Score,       //用户分数
			WinCount:    u.WinCount,    //胜利盘数
			LostCount:   u.LostCount,   //失败盘数
			DrawCount:   u.DrawCount,   //和局盘数
			FleeCount:   u.FleeCount,   //逃跑盘数
			Experience:  u.Experience,  //用户经验
			NickName:    u.NickName,    //昵称
			HeaderUrl:   u.HeadImgUrl,  //头像
		})
	})

	Chat.ChanRPC.Go("addRoomMember", room.ChatRoomId, user.Agent)
	room.setUsetStatus(user, US_SIT)
}

func (room *Room) UserReady(args []interface{}) {
	//recvMsg := args[0].(*msg.C2G_UserReady)
	user := args[1].(*client.User)
	if user.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	room.setUsetStatus(user, US_READY)
	if room.isAllReady() {
		room.StartGame()
	}
}

func (room *Room) UserOffline(args []interface{}) {
	user := args[0].(*client.User)
	if user.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	room.setUsetStatus(user, US_OFFLINE)
	if room.Temp.TimeOffLineCount != 0 {
		room.KickOut[user.Id] = room.Skeleton.AfterFunc(time.Duration(room.Temp.TimeOffLineCount)*time.Second, func() {
			room.OfflineKickOut(user)
		})
	} else {
		room.OfflineKickOut(user)
	}
}

func (room *Room) UserReLogin(args []interface{}) {
	user := args[0].(*client.User)
	if user.Status == US_READY {
		log.Debug("user status is ready at UserReady")
		return
	}

	tm, ok := room.KickOut[user.Id]
	if ok {
		tm.Stop()
		delete(room.KickOut, user.Id)
	}

	room.setUsetStatus(user, US_PLAYING)
}

/////////////////// help
func (room *Room) setUsetStatus(user *client.User, stu int) {
	user.Status = stu
	room.SendMsgAll(&msg.G2C_UserStatus{
		UserID: user.Id,
		UserStatus: &msg.UserStu{
			TableID:    room.GetRoomId(),
			ChairID:    user.ChairId,
			UserStatus: user.Status,
		},
	})
}

func (room *Room) isAllReady() bool {
	for _, u := range room.Users {
		if u == nil || u.Status != US_READY {
			return false
		}
	}
	return true
}

func (room *Room) StartGame() {
	log.Debug("begin start game TBNN")

/*	//初始化
	room.RepertoryCard = make([]int, MAX_REPERTORY)
	room.HandCardCount = make([]int, room.UserCnt)
	for i := 0 i < room.UserCnt i++ {
		room.CardIndex[i] = make([]int, MAX_INDEX)
	}
	room.ChiHuKind = make([]int, room.UserCnt)
	room.ChiPengCount = make([]int, room.UserCnt)
	room.GangCard = make([]bool, room.UserCnt) //杠牌状态
	room.GangCount = make([]int, room.UserCnt)
	room.Ting = make([]bool, room.UserCnt)
	room.UserAction = make([]int, room.UserCnt)
	room.PerformAction = make([]int, room.UserCnt)
	room.DiscardCard = make([][]int, room.UserCnt)
	room.DiscardCount = make([]int, room.UserCnt)
	room.UserGangScore = make([]int, room.UserCnt)
	room.WeaveItemArray = make([][]*msg.WeaveItem, room.UserCnt)
	for i, _ := range room.WeaveItemArray {
		room.WeaveItemArray[i] = make([]*msg.WeaveItem, MAX_WEAVE)
	}
	room.WeaveItemCount = make([]int, room.UserCnt)
	room.ChiHuRight = make([]int, room.UserCnt)

	for i := 0 i < room.UserCnt i++ {
		room.HeapCardInfo[i] = make([]int, 2)
	}

	room.Status = RoomStatusStarting
	Sice1 := util.RandInterval(1, 7)
	Sice2 := util.RandInterval(1, 7)
	minSice := int(math.Min(float64(Sice1), float64(Sice2)))
	room.LeftCardCount = MAX_REPERTORY
	room.SiceCount = Sice2<<8 | Sice1
	room.SendCardCount = 0
	room.UserActionDone = false
	room.SendStatus = Not_Send
	room.GangStatus = WIK_GANERAL
	room.ProvideGangUser = INVALID_CHAIR
	room.gameLogic.RandCardList(room.RepertoryCard, CardDataArray)

	//红中可以当财神
	room.MagicIndex = room.gameLogic.SwitchToCardIndex(0x35)
	room.gameLogic.SetMagicIndex(room.MagicIndex)
	room.PlayerCount = room.GetCurlPlayerCount()
	//分发扑克
	room.ForEachUser(func(u *client.User) {
		room.LeftCardCount -= (MAX_COUNT - 1)
		room.MinusHeadCount += (MAX_COUNT - 1)
		room.gameLogic.SwitchToCardIndex3(room.RepertoryCard[room.LeftCardCount:], MAX_COUNT-1, room.CardIndex[u.ChairId])
	})

	template, ok := base.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	if !ok {
		log.Error("not foud game template at hzmj KindID :%d,  ServerId:%d", room.Kind, room.ServerId)
		return
	}

	OwnerUser, _ := room.GetUserByUid(room.Owner)
	if room.BankerUser == INVALID_CHAIR && (template.ServerType&GAME_GENRE_PERSONAL) != 0 { //房卡模式下先把庄家给房主
		if OwnerUser != nil {
			room.BankerUser = OwnerUser.ChairId
		} else {
			log.Error("get bamkerUser error at StartGame")
		}
	}

	if room.BankerUser == INVALID_CHAIR {
		room.BankerUser = util.RandInterval(0, room.UserCnt-1)
	}

	if room.BankerUser >= room.UserCnt {
		log.Error(" room.BankerUser >= room.UserCnt %d,  %d", room.BankerUser, room.UserCnt)
	}

	room.MinusHeadCount++
	room.SendCardData = room.RepertoryCard[room.LeftCardCount]
	room.LeftCardCount--
	room.CardIndex[room.BankerUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]++
	room.ProvideCard = room.SendCardData
	room.ProvideUser = room.BankerUser
	room.CurrentUser = room.BankerUser

	//堆立信息
	SiceCount := LOBYTE(room.SiceCount) + HIBYTE(room.SiceCount)
	TakeChairID := (room.BankerUser + SiceCount - 1) % room.UserCnt
	TakeCount := MAX_REPERTORY - room.LeftCardCount
	for i := 0 i < room.UserCnt i++ {
		//计算数目
		var ValidCount int
		if i == 0 {
			ValidCount = HEAP_FULL_COUNT - room.HeapCardInfo[TakeChairID][1] - (minSice)*2
		} else {
			ValidCount = HEAP_FULL_COUNT - room.HeapCardInfo[TakeChairID][1]
		}

		RemoveCount := int(math.Min(float64(ValidCount), float64(TakeCount)))

		//提取扑克
		TakeCount -= RemoveCount
		if i == 0 {
			room.HeapCardInfo[TakeChairID][1] += RemoveCount
		} else {
			room.HeapCardInfo[TakeChairID][0] += RemoveCount
		}

		//完成判断
		if TakeCount == 0 {
			room.HeapHead = TakeChairID
			room.HeapTail = (room.BankerUser + SiceCount - 1) % room.UserCnt
			break
		}
		//切换索引
		TakeChairID = (TakeChairID + room.UserCnt - 1) % room.UserCnt
	}

	room.UserAction = make([]int, room.UserCnt)

	gangCardResult := &TagGangCardResult{}
	room.UserAction[room.BankerUser] |= room.gameLogic.AnalyseGangCardEx(room.CardIndex[room.BankerUser], nil, 0, 0, gangCardResult)

	//胡牌判断
	chr := 0
	room.CardIndex[room.BankerUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]--
	room.UserAction[room.BankerUser] |= room.gameLogic.AnalyseChiHuCard(room.CardIndex[room.BankerUser], nil, 0, room.SendCardData, chr, true)
	room.CardIndex[room.BankerUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]++
	room.HandCardCount[room.BankerUser]++

	//听牌判断
	Count := 0
	HuData := &msg.G2C_Hu_Data{OutCardData: make([]int, MAX_COUNT), HuCardCount: make([]int, MAX_COUNT), HuCardData: make([][]int, MAX_COUNT), HuCardRemainingCount: make([][]int, MAX_COUNT)}
	if room.Ting[room.BankerUser] == false {
		Count = room.gameLogic.AnalyseTingCard(room.CardIndex[room.BankerUser], []*msg.WeaveItem{}, 0, HuData.OutCardData, HuData.HuCardCount, HuData.HuCardData)
		HuData.OutCardCount = Count
		if Count > 0 {
			room.UserAction[room.BankerUser] |= WIK_LISTEN
			for i := 0 i < MAX_COUNT i++ {
				if HuData.HuCardCount[i] > 0 {
					for j := 0 j < HuData.HuCardCount[i] j++ {
						HuData.HuCardRemainingCount[i] = append(HuData.HuCardRemainingCount[i], room.GetRemainingCount(room.BankerUser, HuData.HuCardData[i][j]))
					}
				} else {
					break
				}
			}
			OwnerUser.WriteMsg(HuData)
		}
	}

	//构造变量
	GameStart := &mj_hz_msg.G2C_HZMG_GameStart{}
	GameStart.BankerUser = room.BankerUser
	GameStart.SiceCount = room.SiceCount
	GameStart.HeapHead = room.HeapHead
	GameStart.HeapTail = room.HeapTail
	GameStart.MagicIndex = room.MagicIndex
	GameStart.HeapCardInfo = room.HeapCardInfo
	GameStart.CardData = make([]int, MAX_COUNT)
	//发送数据
	OutCard := make([][]int, MAX_COUNT)
	for i, u := range room.Users {
		if u == nil {
			continue
		}

		GameStart.UserAction = room.UserAction[i]
		room.gameLogic.SwitchToCardData2(room.CardIndex[i], GameStart.CardData)
		GameStart.OutCardCount = 0

		if i == room.BankerUser && Count > 0 {
			GameStart.OutCardCount = int(Count)
			GameStart.OutCardData = OutCard
		}

		u.WriteMsg(GameStart)
	}

	if (template.ServerType & GAME_GENRE_MATCH) != 0 {
		//room.SetGameTimer(IDI_CHECK_TABLE, 30000, -1, NULL)
		//room.WaitTime = 0
	}

	room.EndTime.Stop()
	room.EndTime = room.Skeleton.AfterFunc(time.Duration(room.TimeLimit)*time.Second, room.AfterGameTimeOut)
	log.Debug("end startgame ... ")*/

	// 初始化
	for i, _ := range room.IsOpenCard {
		room.IsOpenCard[i] = false
	}

	for i, _ := range room.CardData {
		for j, _ := range room.CardData[i] {
			room.CardData[i][j] = 0
		}
	}

	//游戏状态
	room.Status = RoomStatusStarting

	//用户状态
	room.ForEachUser(func(u *client.User) {
		room.setUsetStatus(u, US_PLAYING)
	})

	for i,_ := range room.BuckleServiceCharge {
		room.BuckleServiceCharge[i] = false
	}

	//设置变量
	CallBanker := nn_tb_msg.G2C_TBNN_CallBanker{}
		CallBanker.Qiang_Start = false//不用抢庄 直接开始

	CallBanker.FirstTimes = true
	CallBanker.CallBanker = room.CurrentUser

	for _, u := range room.Users {
		if u == nil {
			continue
		}
		u.WriteMsg(CallBanker)
	}


	//m_pITableFrame->SendLookonData(INVALID_CHAIR,SUB_S_CALL_BANKER,&CallBanker,sizeof(CallBanker))
	//删除时间
	//m_pITableFrame->KillGameTimer(IDI_SO_OPERATE) 代打定时器
	//m_pITableFrame->SetGameTimer(IDI_SO_OPERATE,TIME_SO_OPERATE,1,0)


	//随机扑克
	AllCard	:= nn_tb_msg.G2C_TBNN_AllCard{}
	var TempArray [GAME_PLAYER*MAX_COUNT]int
	m_GameLogic.RandCardList(bTempArray,sizeof(bTempArray))
	for (WORD i=0i<m_wPlayerCounti++)
	{

		//派发扑克
		CopyMemory(m_cbHandCardData[i],&bTempArray[i*MAX_COUNT],MAX_COUNT)
		//CopyMemory(AllCard.cbCardData[i],&bTempArray[i*MAX_COUNT],MAX_COUNT)

		IServerUserItem *pIServerUser=m_pITableFrame->GetTableUserItem(i)
		if(pIServerUser==NULL)continue
		if(pIServerUser->IsAndroidUser())AllCard.bAICount[i] =true

		m_bBuckleServiceCharge[i]=true
	}

	//m_cbHandCardData[0][0]=0x34
	//m_cbHandCardData[0][1]=0x37
	//m_cbHandCardData[0][2]=0x39
	//m_cbHandCardData[0][3]=0x4E
	//m_cbHandCardData[0][4]=0x4F

	//m_cbHandCardData[2][0]=0x3B
	//m_cbHandCardData[2][1]=0x3D
	//m_cbHandCardData[2][2]=0x0D
	//m_cbHandCardData[2][3]=0x1D
	//m_cbHandCardData[2][4]=0x2D//5hua

	//m_cbHandCardData[1][0]=0x05
	//m_cbHandCardData[1][1]=0x38
	//m_cbHandCardData[1][2]=0x12
	//m_cbHandCardData[1][3]=0x0B
	//m_cbHandCardData[1][4]=0x28

	//m_cbHandCardData[3][0]=0x31
	//m_cbHandCardData[3][1]=0x3C
	//m_cbHandCardData[3][2]=0x3D
	//m_cbHandCardData[3][3]=0x06
	//m_cbHandCardData[3][4]=0x18

	////牛牛数据
	//BOOL bUserOxData[GAME_PLAYER]={0}
	//ZeroMemory(bUserOxData, sizeof(bUserOxData))
	//for (WORD i = 0 i < GAME_PLAYER i++)
	//{
	//	IServerUserItem * pIServerUserItem=m_pITableFrame->GetTableUserItem(i)
	//	if(pIServerUserItem==NULL) continue

	//	if(m_GameLogic.GetOxCard(m_cbHandCardData[i], MAX_COUNT))
	//	  bUserOxData[i]=TRUE
	//}
	////预分析牌型
	//WORD wWinUser = INVALID_CHAIR
	//WORD wLostUser = INVALID_CHAIR
	////查找数据
	//for (WORD i = 0 i < GAME_PLAYER i++)
	//{
	//	//用户过滤
	//	//if (m_cbPlayStatus[i] == FALSE)
	//	//	continue
	//	IServerUserItem * pIServerUserItem=m_pITableFrame->GetTableUserItem(i)
	//	if(pIServerUserItem==NULL) continue

	//	//设置用户
	//	if (wWinUser == INVALID_CHAIR)
	//	{
	//		wWinUser = i
	//		continue
	//	}

	//	//对比扑克
	//	if (m_GameLogic.CompareCard(m_cbHandCardData[i], m_cbHandCardData[wWinUser], MAX_COUNT, bUserOxData[i], bUserOxData[wWinUser]))
	//	{
	//		wWinUser = i
	//	}
	//}
	//for (WORD i = 0 i < GAME_PLAYER i++)
	//{
	//	//用户过滤
	//	//if (m_cbPlayStatus[i] == FALSE)
	//	//	continue
	//	IServerUserItem * pIServerUserItem=m_pITableFrame->GetTableUserItem(i)
	//	if(pIServerUserItem==NULL) continue

	//	//设置用户
	//	if (wLostUser == INVALID_CHAIR)
	//	{
	//		wLostUser = i
	//		continue
	//	}

	//	//对比扑克
	//	if (!m_GameLogic.CompareCard(m_cbHandCardData[i], m_cbHandCardData[wLostUser], MAX_COUNT, bUserOxData[i], bUserOxData[wLostUser]))
	//	{
	//		wLostUser = i
	//	}
	//}
	////黑白名单控制
	//CArray<int> canUsePlayerLst
	//canUsePlayerLst.RemoveAll()

	//if(m_pITableFrame)
	//{
	//	//-----------------------------------//
	//	int iWhiteUser=m_pITableFrame->GetBlackWhiteControl(canUsePlayerLst)
	//	if(iWhiteUser==INVALID_CHAIR && canUsePlayerLst.GetCount()>0) iWhiteUser=canUsePlayerLst.GetAt(rand()%canUsePlayerLst.GetCount()) //剩下玩家处理
	//	if(iWhiteUser!=INVALID_CHAIR && m_lStockScore > 0)
	//	{
	//		if(iWhiteUser!=wWinUser)
	//		{
	//			BYTE bCardData[MAX_COUNT] = {0}
	//			memcpy(bCardData,m_cbHandCardData[wWinUser],sizeof(bCardData))
	//			memcpy(m_cbHandCardData[wWinUser],m_cbHandCardData[iWhiteUser],sizeof(bCardData))
	//			memcpy(m_cbHandCardData[iWhiteUser],bCardData,sizeof(bCardData))
	//		}
	//	}
	//	//判断是否有黑名单
	//	int iBlackUser=INVALID_CHAIR
	//	for(int i=0i<m_wPlayerCounti++)
	//	{
	//		//用户过滤
	//		if (m_cbPlayStatus[i] == FALSE)
	//			continue
	//	    IServerUserItem *pCurUser=m_pITableFrame->GetTableUserItem(i)
	//		if(pCurUser==NULL) continue
	//		tagUserListProp curUserInf=pCurUser->GetUserListProp()
	//		if(curUserInf.iListProp==2)
	//		{
	//			iBlackUser=i
	//			break
	//		}
	//	}
	//	if(wLostUser!=INVALID_CHAIR && iBlackUser!=INVALID_CHAIR)
	//	{
	//		if(iBlackUser!=wLostUser)
	//		{
	//			BYTE bCardData[MAX_COUNT] = {0}
	//			memcpy(bCardData,m_cbHandCardData[wLostUser],sizeof(bCardData))
	//			memcpy(m_cbHandCardData[wLostUser],m_cbHandCardData[iBlackUser],sizeof(bCardData))
	//			memcpy(m_cbHandCardData[iBlackUser],bCardData,sizeof(bCardData))
	//		}
	//	}
	//}

	//canUsePlayerLst.RemoveAll()
	//重拷数据
	for (WORD i=0i<m_wPlayerCounti++)
	{
		//派发扑克
		CopyMemory(AllCard.cbCardData[i],m_cbHandCardData[i],MAX_COUNT)
	}
	//发送数据
	for (WORD i=0i<m_wPlayerCounti++)
	{
		IServerUserItem *pIServerUser=m_pITableFrame->GetTableUserItem(i)
		if(pIServerUser==NULL)continue
#ifndef _DEBUG
		if(CUserRight::IsGameCheatUser(pIServerUser->GetUserRight())==false || m_bSpecialClient[i]==false)continue
#endif
		//fdl modify
		//m_pITableFrame->SendTableData(i,SUB_S_ALL_CARD,&AllCard,sizeof(AllCard))
		//只发给机器人。不能发到客户端
		if(pIServerUser->IsAndroidUser())
		{
			m_pITableFrame->SendTableData(i,SUB_S_ALL_CARD,&AllCard,sizeof(AllCard))
		}
	}

    //记录房间局数
	m_HistoryScore.OnRecordGameCount()


	return
}

//游戏结束
func (room *Room) OnEventGameConclude(ChairId int, user *client.User, cbReason int) bool {
	template, ok := base.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	if !ok {
		log.Error("at OnEventGameConclude not foud tempplate")
		return false
	}
	if (template.ServerType & GAME_GENRE_MATCH) != 0 {
		//room.KillGameTimer(IDI_CHECK_TABLE)
	}

	switch cbReason {
	case GER_NORMAL: //常规结束
		//变量定义
		GameConclude := &mj_hz_msg.G2C_HZMJ_GameConclude{}
		GameConclude.ChiHuKind = make([]int, room.UserCnt)
		GameConclude.CardCount = make([]int, room.UserCnt)
		GameConclude.HandCardData = make([][]int, room.UserCnt)
		GameConclude.GameScore = make([]int, room.UserCnt)
		GameConclude.GangScore = make([]int, room.UserCnt)
		GameConclude.Revenue = make([]int, room.UserCnt)
		GameConclude.ChiHuRight = make([]int, room.UserCnt)
		GameConclude.MaCount = make([]int, room.UserCnt)
		GameConclude.MaData = make([]int, room.UserCnt)

		for i, _ := range GameConclude.HandCardData {
			GameConclude.HandCardData[i] = make([]int, MAX_COUNT)
		}

		GameConclude.SendCardData = room.SendCardData
		GameConclude.LeftUser = INVALID_CHAIR
		room.ChiHuKind = make([]int, room.UserCnt)
		//结束信息
		for i := 0 i < room.UserCnt i++ {
			GameConclude.ChiHuKind[i] = room.ChiHuKind[i]
			//权位过滤
			if room.ChiHuKind[i] == WIK_CHI_HU {
				room.FiltrateRight(i, &room.ChiHuRight[i])
				GameConclude.ChiHuRight[i] = room.ChiHuRight[i]
			}
			GameConclude.CardCount[i] = room.gameLogic.SwitchToCardData2(room.CardIndex[i], GameConclude.HandCardData[i])
		}

		//计算胡牌输赢分
		UserGameScore := make([]int, room.UserCnt)
		room.CalHuPaiScore(UserGameScore)

		//拷贝码数据
		GameConclude.MaCount = room.UserMaCount
		nCount := room.MaCount
		if nCount > 1 {
			nCount++
		}

		for i := 0 i < nCount i++ {
			GameConclude.MaData[i] = room.RepertoryCard[room.MinusLastCount+i]
		}

		//积分变量
		ScoreInfoArray := make([]*msg.TagScoreInfo, room.UserCnt)

		GameConclude.ProvideUser = room.ProvideUser
		GameConclude.ProvideCard = room.ProvideCard

		//统计积分
		for i, u := range room.Users {
			if u.Status != US_PLAYING {
				continue
			}
			GameConclude.GameScore[i] = UserGameScore[i]
			//胡牌分算完后再加上杠的输赢分就是玩家本轮最终输赢分
			GameConclude.GameScore[i] += room.UserGangScore[i]
			GameConclude.GangScore[i] = room.UserGangScore[i]

			//收税
			if GameConclude.GameScore[i] > 0 && (template.ServerType&GAME_GENRE_GOLD) != 0 {
				GameConclude.Revenue[i] = room.CalculateRevenue(i, GameConclude.GameScore[i])
				GameConclude.GameScore[i] -= GameConclude.Revenue[i]
			}

			ScoreInfoArray[i] = &msg.TagScoreInfo{}
			ScoreInfoArray[i].Revenue = GameConclude.Revenue[i]
			ScoreInfoArray[i].Score = GameConclude.GameScore[i]
			if ScoreInfoArray[i].Score > 0 {
				ScoreInfoArray[i].Type = SCORE_TYPE_WIN
			} else {
				ScoreInfoArray[i].Type = SCORE_TYPE_LOSE
			}

			//历史积分
			if room.HistoryScores[i] == nil {
				room.HistoryScores[i] = &HistoryScore{}
			}
			room.HistoryScores[i].TurnScore = GameConclude.GameScore[i]
			room.HistoryScores[i].CollectScore += GameConclude.GameScore[i]

			if room.Record.Count < 32 {
				if len(room.Record.DetailScore[i]) < 1 {
					room.Record.DetailScore[i] = make([]int, 32)
				}
				room.Record.DetailScore[i][room.Record.Count] = GameConclude.GameScore[i]
				room.Record.AllScore[i] += GameConclude.GameScore[i]
			}
		}
		room.Record.Count++
		if (template.ServerType & GAME_GENRE_PERSONAL) != 0 { //房卡模式
			room.SendMsgAll(room.Record)
		}

		//发送数据
		room.SendMsgAll(GameConclude)
		//todo
		//room.pITableFrame->SendLookonData(INVALID_CHAIR, SUB_S_GAME_CONCLUDE, &GameConclude, sizeof(GameConclude))

		//写入积分 todo
		room.WriteTableScore(ScoreInfoArray, room.PlayCount, HZMJ_CHANGE_SOURCE)

		//结束游戏
		room.GameEnd(false)

		if (template.ServerType & GAME_GENRE_PERSONAL) != 0 { //房卡模式
			if room.IsDissumGame { //当前朋友局解散清理记录
				room.Record = &msg.G2C_Record{HuCount: make([]int, room.UserCnt), MaCount: make([]int, room.UserCnt), AnGang: make([]int, room.UserCnt), MingGang: make([]int, room.UserCnt), AllScore: make([]int, room.UserCnt), DetailScore: make([][]int, room.UserCnt)}
			}

		}
		return true
	case GER_USER_LEAVE: //用户强退
		if (template.ServerType & GAME_GENRE_PERSONAL) != 0 { //房卡模式
			return true
		}
		//自动托管
		room.OnUserTrustee(user.ChairId, true)
		return true
	case GER_DISMISS: //游戏解散
		//变量定义

		GameConclude := &mj_hz_msg.G2C_HZMJ_GameConclude{}
		GameConclude.ChiHuKind = make([]int, room.UserCnt)
		GameConclude.CardCount = make([]int, room.UserCnt)
		GameConclude.HandCardData = make([][]int, room.UserCnt)
		GameConclude.GameScore = make([]int, room.UserCnt)
		GameConclude.GangScore = make([]int, room.UserCnt)
		GameConclude.Revenue = make([]int, room.UserCnt)
		GameConclude.ChiHuRight = make([]int, room.UserCnt)
		GameConclude.MaCount = make([]int, room.UserCnt)
		GameConclude.MaData = make([]int, room.UserCnt)
		for i, _ := range GameConclude.HandCardData {
			GameConclude.HandCardData[i] = make([]int, MAX_COUNT)
		}

		room.BankerUser = INVALID_CHAIR

		GameConclude.SendCardData = room.SendCardData

		//用户扑克
		for i := 0 i < room.UserCnt i++ {
			GameConclude.CardCount[i] = room.gameLogic.SwitchToCardData2(room.CardIndex[i], GameConclude.HandCardData[i])
		}

		//发送信息
		room.SendMsgAll(GameConclude)
		//todo
		//room.pITableFrame->SendLookonData(INVALID_CHAIR, SUB_S_GAME_CONCLUDE, &GameConclude, sizeof(GameConclude))

		//结束游戏
		room.GameEnd(true)

		if (template.ServerType & GAME_GENRE_PERSONAL) != 0 { //房卡模式
			if room.IsDissumGame { //当前朋友局解散清理记录
				room.Record = &msg.G2C_Record{HuCount: make([]int, room.UserCnt), MaCount: make([]int, room.UserCnt), AnGang: make([]int, room.UserCnt), MingGang: make([]int, room.UserCnt), AllScore: make([]int, room.UserCnt), DetailScore: make([][]int, room.UserCnt)}
			}
		}

		return true
	}

	log.Error("at OnEventGameConclude error  ")
	return false
}

//todo
func (room *Room) GameEnd(Forced bool) {
	if Forced {
		room.Destroy()
		return
	}
	room.ForEachUser(func(u *client.User) {
		room.setUsetStatus(u, US_FREE)
	})

	room.PlayCount++
	if room.PlayCount >= room.Temp.PlayTurnCount {
		room.Destroy()
	}
}

func (room *Room) GetRemainingCount(ChairId int, cbCardData int) int {
	cbIndex := room.gameLogic.SwitchToCardIndex(cbCardData)
	Count := 0
	for i := room.MinusLastCount i < MAX_REPERTORY-room.MinusHeadCount i++ {
		if room.RepertoryCard[i] == cbCardData {
			Count++
		}
	}
	for i := 0 i < room.UserCnt i++ {
		if i == ChairId {
			continue
		}
		Count += room.CardIndex[i][cbIndex]
	}
	return Count
}

//权位过滤
func (room *Room) FiltrateRight(wWinner int, chr *int) {
	//自摸
	if wWinner == room.ProvideUser {
		*chr |= CHR_ZI_MO
	} else if room.GangStatus == WIK_MING_GANG {
		*chr |= CHR_QIANG_GANG_HU
	} else {
		log.Error("AT FiltrateRight")
	}
	return
}

//算分
func (room *Room) CalHuPaiScore(EndScore []int) {
	room.UserMaCount = make([]int, room.UserCnt)
	CellScore := room.Source
	UserScore := make([]int, room.UserCnt) //玩家手上分
	for i, u := range room.Users {
		if u == nil {
			continue
		}

		if u.Status != US_PLAYING {
			continue
		}
		UserScore[i] = int(u.Score)
	}

	WinUser := make([]int, room.UserCnt)
	WinCount := 0

	for i := 0 i < room.UserCnt i++ {
		if WIK_CHI_HU == room.ChiHuKind[(room.BankerUser+i)%room.UserCnt] {
			WinUser[WinCount] = (room.BankerUser + i) % room.UserCnt
			WinCount++

			//统计胡牌次数
			room.Record.HuCount[(room.BankerUser+i)%room.UserCnt]++
		}
	}

	if WinCount > 0 {
		//有人胡牌
		bZiMo := (room.ProvideUser == WinUser[0])
		if bZiMo {
			//自摸
			cbTimes := room.GetTimes(WinUser[0])
			for i := 0 i < room.UserCnt i++ {

				if i != WinUser[0] {
					EndScore[i] -= cbTimes * CellScore
					EndScore[WinUser[0]] += cbTimes * CellScore
				}
			}
		} else {
			//抢杠
			for i := 0 i < WinCount i++ {
				cbTimes := room.GetTimes(WinUser[i])
				for j := 0 j < room.UserCnt j++ {
					if j != WinUser[i] {
						EndScore[WinUser[i]] += cbTimes * CellScore
					}
				}
				EndScore[room.ProvideUser] -= EndScore[WinUser[i]]
			}
		}

		//谁胡谁当庄
		room.BankerUser = WinUser[0]
		if WinCount > 1 { //多个玩家胡牌，放炮者当庄
			room.BankerUser = room.ProvideUser
		}
	} else { //荒庄
		room.BankerUser = room.LastCatchCardUser //最后一个摸牌的人当庄
	}
}

func (room *Room) GetTimes(wChairId int) int {
	cbScore := 0
	room.UserMaCount[wChairId] = room.MaCount
	if room.MaCount == 1 { //一码全中
		carddata := room.RepertoryCard[room.MinusLastCount]
		if room.gameLogic.GetCardColor(carddata) < 0x30 {
			cbScore = int(room.gameLogic.GetCardValue(carddata))
		} else { //红中10分
			cbScore = 10
		}
		//统计中码个数
		room.Record.MaCount[wChairId]++
	} else { //2-6码
		if room.CardIndex[wChairId][room.MagicIndex] == 0 && room.gameLogic.SwitchToCardIndex(room.ProvideCard) != room.MagicIndex { //胡牌手中没红中，加一个码
			room.UserMaCount[wChairId]++
		}

		for i := 0 i < room.UserMaCount[wChairId] i++ {
			carddata := room.RepertoryCard[room.MinusLastCount+i]

			if room.gameLogic.GetCardValue(carddata)%4 == 1 { //1,5,9,红中 算中码
				cbScore += 2
				//统计中码个数
				room.Record.MaCount[wChairId]++
			}
		}
	}

	return cbScore + 2 //基础倍数+2
}

//计算税收 //可以移植到base
func (room *Room) CalculateRevenue(ChairId, lScore int) int {
	//效验参数

	if ChairId >= room.UserCnt {
		return 0
	}

	template, ok := base.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	if !ok {
		log.Error("at CalculateRevenue no foud template ")
		return 0
	}

	//计算税收
	if (template.RevenueRatio > 0 || template.PersonalRoomTax > 0) && (lScore >= REVENUE_BENCHMARK) {
		//获取用户
		user := room.GetUserByChairId(ChairId)
		if user == nil {
			log.Error("at CalculateRevenue not foud user user.ChairId:%d", user.ChairId)
			return 0
		}

		//计算税收
		lRevenue := lScore * template.RevenueRatio / REVENUE_DENOMINATOR

		if (template.ServerType & GAME_GENRE_PERSONAL) != 0 {
			lRevenue = lScore * (template.RevenueRatio + template.PersonalRoomTax) / REVENUE_DENOMINATOR
		}
		return lRevenue
	}
	return 0
}

//响应判断
func (room *Room) EstimateUserRespond(wCenterUser int, cbCenterCard int, EstimatKind int) bool {
	//变量定义
	bAroseAction := false

	//用户状态
	room.Response = make([]bool, room.UserCnt)
	room.UserAction = make([]int, room.UserCnt)
	room.PerformAction = make([]int, room.UserCnt)

	//动作判断
	for i, u := range room.Users {
		if u == nil {
			continue
		}
		//用户过滤
		if wCenterUser == i || room.Trustee[i] {
			continue
		}

		//出牌类型
		if EstimatKind == EstimatKind_OutCard {
			//吃碰判断
			if u.UserLimit&LimitPeng == 0 {
				//碰牌判断
				room.UserAction[i] |= room.gameLogic.EstimatePengCard(room.CardIndex[i], cbCenterCard)
			}

			//杠牌判断
			if room.LeftCardCount > room.EndLeftCount && u.UserLimit&LimitGang == 0 {
				room.UserAction[i] |= room.gameLogic.EstimateGangCard(room.CardIndex[i], cbCenterCard)
			}
		}

		//检查抢杠胡
		if EstimatKind == EstimatKind_GangCard {
			//只有庄家和闲家之间才能放炮
			if room.MagicIndex == MAX_INDEX || (room.MagicIndex != MAX_INDEX && cbCenterCard != room.gameLogic.SwitchToCardData(room.MagicIndex)) {
				if u.UserLimit|LimitChiHu == 0 {
					//吃胡判断
					chr := 0
					cbWeaveCount := room.WeaveItemCount[i]
					room.UserAction[i] |= room.gameLogic.AnalyseChiHuCard(room.CardIndex[i], room.WeaveItemArray[i], cbWeaveCount, cbCenterCard, chr, false)
				}
			}
		}

		//结果判断
		if room.UserAction[i] != WIK_NULL {
			bAroseAction = true
		}
	}

	//结果处理
	if bAroseAction {
		//设置变量
		room.ProvideUser = wCenterUser
		room.ProvideCard = cbCenterCard
		room.ResumeUser = room.CurrentUser
		room.CurrentUser = INVALID_CHAIR

		//发送提示
		room.ForEachUser(func(u *client.User) {
			if room.UserAction[u.ChairId] != WIK_NULL {
				u.WriteMsg(&mj_hz_msg.G2C_HZMJ_OperateNotify{
					ActionMask: room.UserAction[u.ChairId],
					ActionCard: room.ProvideCard,
				})
			}
		})
		return true
	}

	return false
}

//派发扑克
func (room *Room) DispatchCardData(wCurrentUser int, bTail bool) bool {
	//状态效验
	if wCurrentUser == INVALID_CHAIR {
		return false
	}

	if room.SendStatus == Not_Send {
		log.Error("at DispatchCardData f room.SendStatus == Not_Send")
		return false
	}

	//丢弃扑克
	if (room.OutCardUser != INVALID_CHAIR) && (room.OutCardData != 0) {
		room.OutCardCount++
		if len(room.DiscardCard[room.OutCardUser]) < 1 {
			room.DiscardCard[room.OutCardUser] = make([]int, 60)
		}

		room.DiscardCard[room.OutCardUser][room.DiscardCount[room.OutCardUser]] = room.OutCardData
		room.DiscardCount[room.OutCardUser]++
	}

	macnt := room.MaCount
	if room.MaCount > 1 {
		macnt = room.MaCount + 1
	}
	//荒庄结束
	if room.LeftCardCount <= macnt { //2-6码要多留一颗，一码全中不用
		room.ChiHuCard = 0
		room.ProvideUser = INVALID_CHAIR
		room.OnEventGameConclude(room.ProvideUser, nil, GER_NORMAL)
		return true
	}

	//发送扑克
	room.ProvideCard = room.GetSendCard(bTail)
	room.SendCardData = room.ProvideCard
	room.LastCatchCardUser = wCurrentUser
	//清除禁止胡牌的牌

	user := room.GetUserByChairId(wCurrentUser)
	if user == nil {
		log.Error("at DispatchCardData not foud user ")
	}

	//清除禁止胡牌的牌
	user.UserLimit |= ^LimitChiHu
	user.UserLimit |= ^LimitPeng
	user.UserLimit |= ^LimitGang

	//设置变量
	room.OutCardUser = INVALID_CHAIR
	room.OutCardData = 0
	room.CurrentUser = wCurrentUser
	room.ProvideUser = wCurrentUser
	room.GangOutCard = false

	if bTail { //从尾部取牌，说明玩家杠牌了,计算分数
		room.CallGangScore()
	}

	//加牌
	room.CardIndex[wCurrentUser][room.gameLogic.SwitchToCardIndex(room.ProvideCard)]++
	//room.UserCatchCardCount[wCurrentUser]++

	if !room.Trustee[wCurrentUser] {
		//胡牌判断
		chr := 0
		room.CardIndex[wCurrentUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]--
		log.Debug("befer %v ", room.UserAction[wCurrentUser])
		room.UserAction[wCurrentUser] |= room.gameLogic.AnalyseChiHuCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser],
			room.WeaveItemCount[wCurrentUser], room.SendCardData, chr, false)
		log.Debug("afert %v ", room.UserAction[wCurrentUser])
		room.CardIndex[wCurrentUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]++

		//杠牌判断
		if (room.LeftCardCount > room.EndLeftCount) && !room.Ting[wCurrentUser] {
			GangCardResult := &TagGangCardResult{}
			room.UserAction[wCurrentUser] |= room.gameLogic.AnalyseGangCardEx(room.CardIndex[wCurrentUser],
				room.WeaveItemArray[wCurrentUser], room.WeaveItemCount[wCurrentUser], room.ProvideCard, GangCardResult)
		}
	}

	log.Debug("aaaaaaaaa %v", room.WeaveItemArray[wCurrentUser])
	//听牌判断
	HuData := &msg.G2C_Hu_Data{OutCardData: make([]int, MAX_COUNT), HuCardCount: make([]int, MAX_COUNT), HuCardData: make([][]int, MAX_COUNT), HuCardRemainingCount: make([][]int, MAX_COUNT)}
	if room.Ting[wCurrentUser] == false {
		cbCount := room.gameLogic.AnalyseTingCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser], room.WeaveItemCount[wCurrentUser], HuData.OutCardData, HuData.HuCardCount, HuData.HuCardData)
		HuData.OutCardCount = int(cbCount)
		if cbCount > 0 {
			room.UserAction[wCurrentUser] |= WIK_LISTEN

			for i := 0 i < MAX_COUNT i++ {
				if HuData.HuCardCount[i] > 0 {
					for j := 0 j < HuData.HuCardCount[i] j++ {
						HuData.HuCardRemainingCount[i] = append(HuData.HuCardRemainingCount[i], room.GetRemainingCount(wCurrentUser, HuData.HuCardData[i][j]))
					}
				} else {
					break
				}
			}

			user.WriteMsg(HuData)
		}
	}

	log.Debug("User Action === %v , %d", room.UserAction, room.UserAction[wCurrentUser])
	//构造数据
	SendCard := &mj_hz_msg.G2C_HZMJ_SendCard{}
	SendCard.SendCardUser = wCurrentUser
	SendCard.CurrentUser = wCurrentUser
	SendCard.Tail = bTail
	SendCard.ActionMask = room.UserAction[wCurrentUser]
	SendCard.CardData = room.ProvideCard

	//发送数据
	user.WriteMsg(SendCard)
	SendCard.CardData = 0
	room.SendMsgAllNoSelf(user.Id, SendCard)

	//todo
	//room.pITableFrame->SendLookonData(INVALID_CHAIR,SUB_S_SEND_CARD, &SendCard, sizeof(SendCard))

	room.UserActionDone = false
	if room.Trustee[wCurrentUser] {
		room.UserActionDone = true
		//room.pITableFrame->SetGameTimer(IDI_OUT_CARD,1000,1,0) todo
	}
	return true
}

//取得扑克
func (room *Room) GetSendCard(bTail bool) int {
	//发送扑克
	room.SendCardCount++
	room.LeftCardCount--

	var cbSendCardData int
	var cbIndexCard int
	if bTail {
		cbSendCardData = room.RepertoryCard[room.MinusLastCount]
		room.MinusLastCount++
	} else {
		room.MinusHeadCount++
		log.Debug("aaaaaaaaaa ", MAX_REPERTORY-room.MinusHeadCount, room.LeftCardCount)
		cbIndexCard = MAX_REPERTORY - room.MinusHeadCount
		cbSendCardData = room.RepertoryCard[cbIndexCard]
	}

	//堆立信息

	if !bTail {
		//切换索引
		cbHeapCount := room.HeapCardInfo[room.HeapHead][0] + room.HeapCardInfo[room.HeapHead][1]
		if cbHeapCount == HEAP_FULL_COUNT {
			room.HeapHead = (room.HeapHead + room.UserCnt - 1) % len(room.HeapCardInfo)
		}
		room.HeapCardInfo[room.HeapHead][0]++
	} else {
		//切换索引
		cbHeapCount := room.HeapCardInfo[room.HeapTail][0] + room.HeapCardInfo[room.HeapTail][1]
		if cbHeapCount == HEAP_FULL_COUNT {
			room.HeapTail = (room.HeapTail + 1) % len(room.HeapCardInfo)
		}
		room.HeapCardInfo[room.HeapTail][1]++
	}

	return cbSendCardData
}

func (room *Room) CallGangScore() {
	lcell := room.Source
	if room.GangStatus == WIK_FANG_GANG { //放杠一家扣分
		for i, u := range room.Users {
			if u.Status != US_PLAYING {
				continue
			}
			if i != room.CurrentUser {
				room.UserGangScore[room.ProvideGangUser] -= lcell
				room.UserGangScore[room.CurrentUser] += lcell
			}
		}
		//记录明杠次数
		room.Record.MingGang[room.CurrentUser]++
	} else if room.GangStatus == WIK_MING_GANG { //明杠每家出1倍
		for i, u := range room.Users {
			if u.Status != US_PLAYING {
				continue
			}
			if i != room.CurrentUser {
				room.UserGangScore[i] -= lcell
				room.UserGangScore[room.CurrentUser] += lcell
			}
		}
		//记录明杠次数
		room.Record.MingGang[room.CurrentUser]++
	} else if room.GangStatus == WIK_AN_GANG { //暗杠每家出2倍
		for i, u := range room.Users {
			if u.Status != US_PLAYING {
				continue
			}
			if i != room.CurrentUser {
				room.UserGangScore[i] -= 2 * lcell
				room.UserGangScore[room.CurrentUser] += 2 * lcell
			}
		}
		//记录暗杠次数
		room.Record.AnGang[room.CurrentUser]++
	}
}

func (room *Room) Operater(user *client.User, cbOperateCard []int, cbOperateCode int, IsZd bool) bool {
	if !IsZd {
		//效验状态
		if room.Response[user.ChairId] {
			return true
		}
		if room.UserAction[user.ChairId] == WIK_NULL {
			return true
		}
		if (cbOperateCode != WIK_NULL) && ((room.UserAction[user.ChairId] & cbOperateCode) == 0) {
			return true
		}

		//变量定义
		wTargetUser := user.ChairId
		cbTargetAction := cbOperateCode

		//设置变量
		user.UserLimit |= ^LimitGang
		room.Response[wTargetUser] = true
		room.PerformAction[wTargetUser] = cbOperateCode
		if cbOperateCard[0] == 0 {
			room.OperateCard[wTargetUser][0] = room.ProvideCard
		} else {
			room.OperateCard[wTargetUser] = cbOperateCard
		}

		//放弃操作
		if cbTargetAction == WIK_NULL {
			////禁止这轮吃胡
			if (room.UserAction[wTargetUser] & WIK_CHI_HU) != 0 {
				user.UserLimit |= LimitChiHu
			}
		}

		//执行判断
		for i := 0 i < room.PlayerCount i++ {
			//获取动作
			cbUserAction := room.UserAction[i]
			if room.Response[i] {
				cbUserAction = room.PerformAction[i]
			}

			//优先级别
			cbUserActionRank := room.gameLogic.GetUserActionRank(cbUserAction)
			cbTargetActionRank := room.gameLogic.GetUserActionRank(cbTargetAction)

			//动作判断
			if cbUserActionRank > cbTargetActionRank {
				wTargetUser = i
				cbTargetAction = cbUserAction
			}
		}
		if !room.Response[wTargetUser] {
			return true
		}

		//放弃操作
		if cbTargetAction == WIK_NULL {
			//用户状态
			room.Response = make([]bool, room.UserCnt)
			room.UserAction = make([]int, room.UserCnt)
			room.OperateCard = make([][]int, room.UserCnt)
			room.PerformAction = make([]int, room.UserCnt)

			room.DispatchCardData(room.ResumeUser, room.GangStatus != WIK_GANERAL)
			return true
		}

		//变量定义
		cbTargetCard := room.OperateCard[wTargetUser][0]

		//出牌变量
		room.SendStatus = Gang_Send
		room.SendCardData = 0
		room.OutCardUser = INVALID_CHAIR
		room.OutCardData = 0

		//胡牌操作
		if cbTargetAction == WIK_CHI_HU {
			//结束信息
			room.ChiHuCard = cbTargetCard

			wChiHuUser := room.BankerUser
			for i := 0 i < room.PlayerCount i++ {
				wChiHuUser = (room.BankerUser + i) % room.PlayerCount
				//过虑判断
				if (room.PerformAction[wChiHuUser] & WIK_CHI_HU) == 0 {
					continue
				}

				//胡牌判断
				cbWeaveItemCount := room.WeaveItemCount[wChiHuUser]
				pWeaveItem := room.WeaveItemArray[wChiHuUser]
				chihuKind := room.gameLogic.AnalyseChiHuCard(room.CardIndex[wChiHuUser], pWeaveItem, cbWeaveItemCount, room.ChiHuCard, room.ChiHuRight[wChiHuUser], false)
				room.ChiHuKind[wChiHuUser] = int(chihuKind)
				//插入扑克
				if room.ChiHuKind[wChiHuUser] != WIK_NULL {
					wTargetUser = wChiHuUser
					//break
				}
			}

			//结束游戏
			room.OnEventGameConclude(room.ProvideUser, nil, GER_NORMAL)

			return true
		}

		//组合扑克
		room.WeaveItemCount[wTargetUser]++
		wIndex := room.WeaveItemCount[wTargetUser]
		room.WeaveItemArray[wTargetUser][wIndex].Param = WIK_GANERAL
		room.WeaveItemArray[wTargetUser][wIndex].CenterCard = cbTargetCard
		room.WeaveItemArray[wTargetUser][wIndex].WeaveKind = cbTargetAction
		if room.ProvideUser == INVALID_CHAIR {
			room.WeaveItemArray[wTargetUser][wIndex].ProvideUser = wTargetUser
		} else {
			room.WeaveItemArray[wTargetUser][wIndex].ProvideUser = room.ProvideUser
		}

		room.WeaveItemArray[wTargetUser][wIndex].CardData[0] = cbTargetCard
		if cbTargetAction&(WIK_LEFT|WIK_CENTER|WIK_RIGHT) != 0 {
			room.WeaveItemArray[wTargetUser][wIndex].CardData[1] = room.OperateCard[wTargetUser][1]
			room.WeaveItemArray[wTargetUser][wIndex].CardData[2] = room.OperateCard[wTargetUser][2]
		} else {
			room.WeaveItemArray[wTargetUser][wIndex].CardData[1] = cbTargetCard
			room.WeaveItemArray[wTargetUser][wIndex].CardData[2] = cbTargetCard
			if cbTargetAction&WIK_GANG != 0 {
				room.WeaveItemArray[wTargetUser][wIndex].Param = WIK_FANG_GANG
				room.WeaveItemArray[wTargetUser][wIndex].CardData[3] = cbTargetCard
			}
		}

		//删除扑克
		switch cbTargetAction {
		case WIK_LEFT: //上牌操作
			//删除扑克
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], room.OperateCard[wTargetUser][1:], 2) {
				log.Error("not foud card at Operater")
				return false
			}
			room.ChiPengCount[wTargetUser]++
			break
		case WIK_RIGHT: //上牌操作
			//删除扑克
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], room.OperateCard[wTargetUser][1:], 2) {
				log.Error("not foud card at Operater")
				return false
			}
			room.ChiPengCount[wTargetUser]++

			break
		case WIK_CENTER: //上牌操作
			//删除扑克
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], room.OperateCard[wTargetUser][1:], 2) {
				log.Error("not foud card at Operater")
				return false
			}
			room.ChiPengCount[wTargetUser]++
			break
		case WIK_PENG: //碰牌操作
			//删除扑克
			cbRemoveCard := []int{cbTargetCard, cbTargetCard}
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], cbRemoveCard, 2) {
				log.Error("not foud card at Operater")
				return false
			}
			room.ChiPengCount[wTargetUser]++
			break
		case WIK_GANG: //杠牌操作
			//删除扑克,被动动作只存在放杠
			cbRemoveCard := []int{cbTargetCard, cbTargetCard, cbTargetCard}
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], cbRemoveCard, int(len(cbRemoveCard))) {
				log.Error("not foud card at Operater")
				return false
			}

			break
		default:
			log.Error("not foud Operater at Operater")
			return false
		}

		//构造结果
		OperateResult := &mj_hz_msg.G2C_HZMJ_OperateResult{}
		OperateResult.OperateUser = wTargetUser
		OperateResult.OperateCode = cbTargetAction
		if room.ProvideUser == INVALID_CHAIR {
			OperateResult.ProvideUser = wTargetUser
		} else {
			OperateResult.ProvideUser = room.ProvideUser
		}

		OperateResult.OperateCard[0] = cbTargetCard
		if cbTargetAction&(WIK_LEFT|WIK_CENTER|WIK_RIGHT) != 0 {
			OperateResult.OperateCard[1] = room.OperateCard[wTargetUser][1]
		} else if cbTargetAction&WIK_PENG != 0 {
			OperateResult.OperateCard[1] = cbTargetCard
			OperateResult.OperateCard[2] = cbTargetCard
		}

		//用户状态
		//用户状态
		room.Response = make([]bool, room.UserCnt)
		room.UserAction = make([]int, room.UserCnt)
		room.PerformAction = make([]int, room.UserCnt)
		room.OperateCard = make([][]int, room.UserCnt)

		//如果非杠牌
		if cbTargetAction != WIK_GANG {
			room.ProvideUser = INVALID_CHAIR
			room.ProvideCard = 0

			gcr := &TagGangCardResult{}
			room.UserAction[wTargetUser] |= room.gameLogic.AnalyseGangCardEx(room.CardIndex[wTargetUser], room.WeaveItemArray[wTargetUser], room.WeaveItemCount[wTargetUser], 0, gcr)

			if room.Ting[wTargetUser] == false {

				HuData := &msg.G2C_Hu_Data{OutCardData: make([]int, MAX_COUNT), HuCardCount: make([]int, MAX_COUNT), HuCardData: make([][]int, MAX_COUNT), HuCardRemainingCount: make([][]int, MAX_COUNT)}
				cbCount := room.gameLogic.AnalyseTingCard(room.CardIndex[wTargetUser], room.WeaveItemArray[wTargetUser], room.WeaveItemCount[wTargetUser], HuData.OutCardData, HuData.HuCardCount, HuData.HuCardData)
				HuData.OutCardCount = cbCount
				if cbCount > 0 {
					room.UserAction[wTargetUser] |= WIK_LISTEN
					for i := 0 i < MAX_COUNT i++ {
						if HuData.HuCardCount[i] > 0 {
							for j := 0 j < HuData.HuCardCount[i] j++ {
								HuData.HuCardRemainingCount[i][j] = room.GetRemainingCount(wTargetUser, HuData.HuCardData[i][j])
							}
						} else {
							break
						}

					}
					user.WriteMsg(HuData)
				}
			}
			OperateResult.ActionMask |= room.UserAction[wTargetUser]
		}

		//发送消息
		room.SendMsgAll(OperateResult)
		//room.pITableFrame->SendLookonData(INVALID_CHAIR,SUB_S_OPERATE_RESULT, &OperateResult, sizeof(OperateResult))

		//设置用户
		room.CurrentUser = wTargetUser

		//杠牌处理
		if cbTargetAction == WIK_GANG {
			room.GangStatus = WIK_FANG_GANG
			if room.ProvideUser == INVALID_CHAIR {
				room.ProvideGangUser = wTargetUser
			} else {
				room.ProvideGangUser = room.ProvideUser
			}
			room.GangCard[wTargetUser] = true
			room.GangCount[wTargetUser]++
			room.DispatchCardData(wTargetUser, true)
		}
		return true
	} else { //主动动作
		//扑克效验

		if (cbOperateCode != WIK_NULL) && (cbOperateCode != WIK_CHI_HU) && (!room.gameLogic.IsValidCard(cbOperateCard[0])) {
			return false
		}

		//设置变量
		room.UserAction[room.CurrentUser] = WIK_NULL
		room.PerformAction[room.CurrentUser] = WIK_NULL

		//执行动作
		switch cbOperateCode {
		case WIK_GANG: //杠牌操作
			room.SendStatus = Gang_Send
			//变量定义
			cbWeaveIndex := int(0xFF)
			cbCardIndex := room.gameLogic.SwitchToCardIndex(cbOperateCard[0])
			wProvideUser := user.ChairId
			cbGangKind := int(WIK_MING_GANG)
			//杠牌处理
			if room.CardIndex[user.ChairId][cbCardIndex] == 1 {
				//寻找组合
				for i := 0 i < room.WeaveItemCount[user.ChairId] i++ {
					cbWeaveKind := room.WeaveItemArray[user.ChairId][i].WeaveKind
					cbCenterCard := room.WeaveItemArray[user.ChairId][i].CenterCard
					if (cbCenterCard == cbOperateCard[0]) && (cbWeaveKind == WIK_PENG) {
						cbWeaveIndex = i
						break
					}
				}

				//效验动作
				if cbWeaveIndex == 0xFF {
					return false
				}
				cbGangKind = WIK_MING_GANG

				//组合扑克
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].Param = WIK_MING_GANG
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].WeaveKind = cbOperateCode
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].CenterCard = cbOperateCard[0]
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].CardData[3] = cbOperateCard[0]

				//杠牌得分
				wProvideUser = room.WeaveItemArray[user.ChairId][cbWeaveIndex].ProvideUser
			} else {
				//扑克效验

				if room.CardIndex[user.ChairId][cbCardIndex] != 4 {
					return false
				}

				//设置变量
				room.WeaveItemCount[user.ChairId]++
				cbWeaveIndex := room.WeaveItemCount[user.ChairId]
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].Param = WIK_AN_GANG
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].ProvideUser = user.ChairId
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].WeaveKind = cbOperateCode
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].CenterCard = cbOperateCard[0]
				for j := 0 j < 4 j++ {
					room.WeaveItemArray[user.ChairId][cbWeaveIndex].CardData[j] = cbOperateCard[0]
				}
			}

			//删除扑克
			room.CardIndex[user.ChairId][cbCardIndex] = 0
			room.GangStatus = cbGangKind
			room.ProvideGangUser = wProvideUser
			room.GangCard[user.ChairId] = true
			room.GangCount[user.ChairId]++

			//构造结果
			OperateResult := &mj_hz_msg.G2C_HZMJ_OperateResult{}
			OperateResult.OperateUser = user.ChairId
			OperateResult.ProvideUser = wProvideUser
			OperateResult.OperateCode = cbOperateCode
			OperateResult.OperateCard[0] = cbOperateCard[0]

			//发送消息
			room.SendMsgAll(OperateResult)
			//room.pITableFrame->SendLookonData(INVALID_CHAIR, SUB_S_OPERATE_RESULT, &OperateResult, sizeof(OperateResult))

			//效验动作
			bAroseAction := false
			if cbGangKind == WIK_MING_GANG {
				bAroseAction = room.EstimateUserRespond(user.ChairId, cbOperateCard[0], EstimatKind_GangCard)
			}

			//发送扑克
			if !bAroseAction {
				room.DispatchCardData(user.ChairId, true)
			}
			return true
		case WIK_CHI_HU: //自摸
			//普通胡牌
			cbWeaveItemCount := room.WeaveItemCount[user.ChairId]
			pWeaveItem := room.WeaveItemArray[user.ChairId]
			if !room.gameLogic.RemoveCard(room.CardIndex[user.ChairId], room.SendCardData) {
				log.Error("not foud card at Operater")
				return false
			}
			kind := room.gameLogic.AnalyseChiHuCard(room.CardIndex[user.ChairId], pWeaveItem, cbWeaveItemCount, room.SendCardData, room.ChiHuRight[user.ChairId], false)
			room.ChiHuKind[user.ChairId] = int(kind)
			//结束信息
			room.ChiHuCard = room.SendCardData
			room.ProvideCard = room.SendCardData

			//结束游戏
			room.OnEventGameConclude(room.ProvideUser, nil, GER_NORMAL)
			return true
		}
		return true
	}
}

//托管
func (room *Room) OnUserTrustee(wChairID int, bTrustee bool) bool {
	//效验状态
	if wChairID >= room.UserCnt {
		return false
	}

	room.Trustee[wChairID] = bTrustee

	room.SendMsgAll(&mj_hz_msg.G2C_HZMJ_Trustee{
		Trustee: bTrustee,
		ChairID: wChairID,
	})

	//m_pITableFrame->SendLookonData(INVALID_CHAIR,SUB_S_TRUSTEE,&Trustee,sizeof(Trustee))

	if bTrustee {
		if wChairID == room.CurrentUser && room.UserActionDone == false {
			cardindex := INVALID_BYTE
			if room.SendCardData != 0 {
				cardindex = room.gameLogic.SwitchToCardIndex(room.SendCardData)
			} else {
				for i := 0 i < MAX_INDEX i++ {
					if room.CardIndex[wChairID][i] > 0 {
						cardindex = i
						break
					}
				}
			}
			room.OnUserOutCard(wChairID, room.gameLogic.SwitchToCardData(cardindex), false)
		} else if room.CurrentUser == INVALID_CHAIR && room.UserActionDone == false {
			operatecard := make([]int, 3)
			user := room.GetUserByChairId(wChairID)
			if user == nil {
				return false
			}
			room.Operater(user, operatecard, WIK_NULL, false)
		}
	}
	return true
}

//用户出牌
func (room *Room) OnUserOutCard(wChairID int, cbCardData int, bSysOut bool) int {
	//效验状态
	if room.Status != RoomStatusStarting {
		log.Error("at OnUserOutCard game status != RoomStatusStarting ")
		return ErrGameNotStart
	}

	//效验参数
	if wChairID != room.CurrentUser {
		log.Error("at OnUserOutCard not self out ")
		return ErrNotSelfOut
	}

	if !room.gameLogic.IsValidCard(cbCardData) {
		log.Error("at OnUserOutCard IsValidCard card ")
		return NotValidCard
	}

	//删除扑克
	if !room.gameLogic.RemoveCard(room.CardIndex[wChairID], cbCardData) {
		log.Error("at OnUserOutCard not have card ")
		return ErrNotFoudCard
	}

	//清除禁止
	user := room.GetUserByChairId(wChairID)
	if user == nil {
		log.Error("at OnUserOutCard not foud user ")
		return ErrUserNotInRoom
	}

	user.UserLimit |= ^LimitChiHu
	user.UserLimit |= ^LimitPeng
	user.UserLimit |= ^LimitGang

	//设置变量
	room.SendStatus = OutCard_Send
	room.SendCardData = 0
	room.UserAction[wChairID] = WIK_NULL
	room.PerformAction[wChairID] = WIK_NULL

	//出牌记录
	room.OutCardUser = wChairID
	room.OutCardData = cbCardData

	//构造数据
	OutCard := &mj_hz_msg.G2C_HZMJ_OutCard{}
	OutCard.OutCardUser = wChairID
	OutCard.OutCardData = cbCardData
	OutCard.SysOut = bSysOut
	//发送消息
	room.SendMsgAll(OutCard)
	//m_pITableFrame->SendLookonData(INVALID_CHAIR, SUB_S_OUT_CARD, &OutCard, sizeof(OutCard))

	room.ProvideUser = wChairID
	room.ProvideCard = cbCardData

	//用户切换
	room.CurrentUser = (wChairID + 1) % room.PlayerCount

	//响应判断
	bAroseAction := room.EstimateUserRespond(wChairID, cbCardData, EstimatKind_OutCard)

	if room.GangStatus != WIK_GANERAL {
		room.GangOutCard = true
		room.GangStatus = WIK_GANERAL
		room.ProvideGangUser = INVALID_CHAIR
	} else {
		room.GangOutCard = false
	}

	//派发扑克
	if !bAroseAction {
		room.DispatchCardData(room.CurrentUser, false)
	}

	return 0
}

//获取对方信息
func (room *Room) GetUserChairInfo(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_REQUserChairInfo)
	user := args[1].(*client.User)
	tagUser := room.GetUserByChairId(recvMsg.ChairID)
	if tagUser == nil {
		log.Error("at GetUserChairInfo no foud tagUser %v, userId:%d", args[0], user.Id)
		return
	}

	user.WriteMsg(&msg.G2C_UserEnter{
		UserID:      tagUser.Id,          //用户 I D
		FaceID:      tagUser.FaceID,      //头像索引
		CustomID:    tagUser.CustomID,    //自定标识
		Gender:      tagUser.Gender,      //用户性别
		MemberOrder: tagUser.MemberOrder, //会员等级
		TableID:     tagUser.RoomId,      //桌子索引
		ChairID:     tagUser.ChairId,     //椅子索引
		UserStatus:  tagUser.Status,      //用户状态
		Score:       tagUser.Score,       //用户分数
		WinCount:    tagUser.WinCount,    //胜利盘数
		LostCount:   tagUser.LostCount,   //失败盘数
		DrawCount:   tagUser.DrawCount,   //和局盘数
		FleeCount:   tagUser.FleeCount,   //逃跑盘数
		Experience:  tagUser.Experience,  //用户经验
		NickName:    tagUser.NickName,    //昵称
		HeaderUrl:   tagUser.HeadImgUrl,  //头像
	})
}

func (room *Room) DissumeRoom(args []interface{}) {
	user := args[0].(*client.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			user.WriteMsg(RenderErrorMessage(retcode, "解散房间失败."))
		}
	}()
	if user.Id != room.Owner {
		retcode = NotOwner
		return
	}

	Cance := &msg.G2C_CancelTable{}
	room.ForEachUser(func(u *client.User) {
		u.WriteMsg(Cance)
	})

	Diis := &msg.G2C_PersonalTableEnd{}
	room.ForEachUser(func(u *client.User) {
		u.WriteMsg(Diis)
	})

	room.Destroy()
}
*/
