package room

import (
	"math"
	. "mj/common/cost"
	"mj/common/msg"
	"mj/gameServer/Chat"
	. "mj/gameServer/common/mj_logic_base"
	"mj/gameServer/db/model/base"
	client "mj/gameServer/user"
	"strconv"
	"time"

	"mj/common/msg/mj_xs_msg"

	"fmt"

	"mj/common/msg/mj_hz_msg"

	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

func RegisterHandler(r *Room) {
	r.ChanRPC.Register("Sitdown", r.Sitdown)
	r.ChanRPC.Register("SetGameOption", r.SetGameOption)
	r.ChanRPC.Register("UserStandup", r.UserStandup)
	r.ChanRPC.Register("OutCard", r.OutCard)
	r.ChanRPC.Register("OperateCard", r.UserOperateCard)
	r.ChanRPC.Register("UserReady", r.UserReady)

}

func (room *Room) OutCard(args []interface{}) {
	recvMsg := args[0].(*mj_hz_msg.C2G_HZMJ_HZOutCard)
	user := args[1].(*client.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			user.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	retcode = room.OnUserOutCard(user.ChairId, recvMsg.CardData, false)
	return
}

func (room *Room) UserOperateCard(args []interface{}) {
	recvMsg := args[0].(*mj_xs_msg.C2G_MJXS_OperateCard)
	user := args[1].(*client.User)
	retcode := 0
	defer func() {
		if retcode != 0 {
			user.WriteMsg(RenderErrorMessage(retcode))
		}
	}()
	if user.ChairId >= room.UserCnt {
		log.Error("user not in room at OperateCard")
		retcode = ErrUserNotInRoom
		return
	}

	//if room.CurrentUser != user.ChairId && room.CurrentUser != INVALID_CHAIR {
	//	log.Error("CurrentUser != user.ChairId at OperateCard")
	//	retcode = ErrUserNotInRoom
	//	return
	//}

	if room.CurrentUser == INVALID_CHAIR {
		room.OnUserOperateCard(user, recvMsg.OperateCard, recvMsg.OperateCode)
	} else {
		room.OnUserOperateCard(user, recvMsg.OperateCard, recvMsg.OperateCode)
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
		CellScore:         room.Source,                              //游戏底分
		IniScore:          room.IniSource,                           //初始分数
		ServerID:          strconv.Itoa(room.GetRoomId()),           //房间编号
		IsJoinGame:        0,                                        //是否参与游戏 todo  tagPersonalTableParameter
		IsGoldOrGameScore: room.IsGoldOrGameScore,                   //金币场还是积分场 0 标识 金币场 1 标识 积分场
	})

	if room.Status == RoomStatusReady { // 没开始
		StatusFree := &mj_xs_msg.G2C_MJXS_StatusPlay{}
		StatusFree.CellScore = room.Source //基础积分
		StatusFree.BankerUser = room.BankerUser
		user.WriteMsg(StatusFree)
	} else { //开始了
		StatusPlay := &mj_xs_msg.G2C_MJXS_StatusPlay{}
		//自定规则
		StatusPlay.SiceCount = room.SiceCount
		StatusPlay.BankerUser = room.BankerUser
		StatusPlay.CurrentUser = room.CurrentUser
		StatusPlay.CellScore = room.Source
		StatusPlay.InitialBankerUser = room.InitialBankerUser
		StatusPlay.FengQuan = room.FengQuan

		//状态变量
		StatusPlay.ActionCard = room.ProvideCard
		StatusPlay.LeftCardCount = room.LeftCardCount
		StatusPlay.ActionMask = WIK_NULL
		if room.Response[user.ChairId] == false {
			StatusPlay.ActionMask = room.UserAction[user.ChairId]
		}

		//历史记录
		//出牌信息
		StatusPlay.EnjoinCardCount = room.EnjoinCardCount[room.CurrentUser]
		StatusPlay.EnjoinCardData = room.EnjoinCardData[room.CurrentUser]

		//历史记录
		StatusPlay.OutCardUser = room.OutCardUser
		StatusPlay.OutCardData = room.OutCardData
		StatusPlay.DiscardCard = room.DiscardCard
		StatusPlay.DiscardCount = room.DiscardCount
		StatusPlay.UserWindCount = room.UserWindCount
		StatusPlay.UserWindCardData = room.UserWindData

		//组合扑克
		StatusPlay.WeaveItemArray = room.WeaveItemArray
		StatusPlay.WeaveCount = room.WeaveItemCount

		//扑克数据
		StatusPlay.CardCount = room.gameLogic.SwitchToCardData2(room.CardIndex[user.ChairId], StatusPlay.CardData)

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

func (room *Room) StartGame() bool {
	log.Debug("begin start game hzmj")
	room.ForEachUser(func(u *client.User) {
		room.setUsetStatus(u, US_PLAYING)
	})

	//初始化
	room.RepertoryCard = make([]int, MAX_REPERTORY)
	for i := 0; i < room.UserCnt; i++ {
		room.CardIndex[i] = make([]int, MAX_INDEX)
	}
	room.ChiHuKind = make([]int, room.UserCnt)
	room.ChiPengCount = make([]int, room.UserCnt)
	room.Ting = make([]bool, room.UserCnt)
	room.UserAction = make([]int, room.UserCnt)
	room.PerformAction = make([]int, room.UserCnt)
	room.DiscardCard = make([][]int, room.UserCnt)
	room.DiscardCount = make([]int, room.UserCnt)
	room.WeaveItemArray = make([][]*mj_xs_msg.TagWeaveItem, room.UserCnt)
	for i, _ := range room.WeaveItemArray {
		room.WeaveItemArray[i] = make([]*mj_xs_msg.TagWeaveItem, MAX_WEAVE)
	}
	room.WeaveItemCount = make([]int, room.UserCnt)
	room.ChiHuRight = make([]int, room.UserCnt)

	room.Status = RoomStatusStarting
	Sice1 := util.RandInterval(1, 7)
	Sice2 := util.RandInterval(1, 7)
	room.LeftCardCount = MAX_REPERTORY
	room.SiceCount = Sice2<<8 | Sice1
	room.SendCardCount = 0
	room.gameLogic.RandCardList(room.RepertoryCard, CardDataArray)

	room.PlayerCount = room.GetCurlPlayerCount()
	//分发扑克
	room.ForEachUser(func(u *client.User) {
		room.LeftCardCount -= (MAX_COUNT - 1)
		room.gameLogic.SwitchToCardIndex3(room.RepertoryCard[room.LeftCardCount:], MAX_COUNT-1, room.CardIndex[u.ChairId])
	})

	room.SendCardCount++
	room.LeftCardCount--
	room.SendCardData = room.RepertoryCard[room.LeftCardCount]
	room.CardIndex[room.BankerUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]++

	//设置变量
	room.ProvideCard = room.SendCardData
	room.ProvideUser = INVALID_CHAIR
	room.CurrentUser = room.BankerUser

	template, ok := base.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	if !ok {
		log.Error("not foud game template at hzmj KindID :%d,  ServerId:%d", room.Kind, room.ServerId)
		return false
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
		room.InitialBankerUser = room.BankerUser
		room.FengQuan = 0 //东风圈
	}

	if room.BankerUser >= room.UserCnt {
		log.Error(" room.BankerUser >= room.UserCnt %d,  %d", room.BankerUser, room.UserCnt)
	}

	//动作分析
	bAroseAction := false
	for i := 0; i < room.PlayerCount; i++ {
		//庄家判断
		if i == room.BankerUser {
			//杠牌判断
			GangCardResult := &TagGangCardResult{}
			room.UserAction[i] |= room.gameLogic.AnalyseGangCard(room.CardIndex[i], nil, 0, GangCardResult)

			//胡牌判断
			ChiHuResult := &TagChiHuResult{}
			cbHandTai := 0
			cbHandFeng := 0
			room.UserAction[i] |= room.gameLogic.AnalyseChiHuCard(room.CardIndex[i], nil, 0, 0, 0, ChiHuResult, &cbHandTai, &cbHandFeng, room.FengQuan,
				0, true)

			cbHandTai += (room.Zfb[i] + room.UserWindCount[i] + room.Dnxb[i] + cbHandFeng)

			if (ChiHuResult.ChiHuKind&CHK_JI_HU != 0) && (ChiHuResult.ChiHuRight == 0 || ChiHuResult.ChiHuRight == CHR_QIANG_GANG) && (cbHandTai == 0) {
				room.UserAction[i] &= (^WIK_CHI_HU)
			}
		}

		//状态设置
		if (bAroseAction == false) && (i != room.BankerUser) && (room.UserAction[i] != WIK_NULL) {
			bAroseAction = true
			room.ResumeUser = room.CurrentUser
			room.CurrentUser = INVALID_CHAIR
		}
	}
	//构造变量
	GameStart := &mj_xs_msg.G2C_MJXS_GameStart{}
	GameStart.BankerUser = room.BankerUser
	GameStart.SiceCount = room.SiceCount
	GameStart.InitialBankerUser = room.InitialBankerUser
	GameStart.FengQuan = room.FengQuan
	GameStart.CurrentUser = INVALID_CHAIR
	GameStart.LeftCardCount = room.LeftCardCount
	GameStart.First = true
	GameStart.CardData = make([]int, MAX_COUNT)

	//花牌数据
	for i := 0; i < room.UserCnt; i++ {
		for j := MAX_INDEX - 8; j < MAX_INDEX; j++ {
			room.WindCount[i] += room.CardIndex[i][j]
		}
		room.SumWindCount += room.WindCount[i]
		room.AllWindCount += room.WindCount[i]
	}

	//预留扑克
	//状态判断
	if room.Status == RoomStatusStarting {
		//规则判断
		if room.AllWindCount < 8 {
			room.RemainCardCount = 14 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
		} else {
			room.RemainCardCount = 16 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
		}

	}
	//发送数据
	for i, u := range room.Users {
		if u == nil {
			continue
		}

		GameStart.UserAction = room.UserAction[i]
		room.gameLogic.SwitchToCardData2(room.CardIndex[i], GameStart.CardData)
		room.WindData[i] = GameStart.CardData[(MAX_COUNT - room.WindCount[i]):]

		GameStart.SunWindCount = room.SumWindCount
		u.WriteMsg(GameStart)
	}

	//花牌替换
	for {
		for i := 0; i < room.UserCnt; i++ {
			//用户索引
			wCurrentUser := (room.BankerUser - i + room.UserCnt) % room.UserCnt

			//花牌i出牌
			for j := 0; j < room.WindCount[wCurrentUser]; j++ {
				cbCardData := room.WindData[wCurrentUser][j]

				if room.gameLogic.IsValidCard(cbCardData) == false {
					return false
				}

				//删除扑克
				if room.gameLogic.RemoveCard(room.CardIndex[wCurrentUser], cbCardData) == false {
					return false
				}

				//出牌记录
				room.OutCardCount++
				room.OutCardUser = wCurrentUser
				room.OutCardData = cbCardData

				//构造数据
				OutCard := &mj_xs_msg.G2C_MJXS_OutCard{}
				OutCard.OutCardUser = wCurrentUser
				OutCard.OutCardData = cbCardData

				//发送消息
				room.SendMsgAll(OutCard)
				room.UserWindData[wCurrentUser][room.UserWindCount[wCurrentUser]] = cbCardData
				room.UserWindCount[wCurrentUser]++

				//替换花番
				if room.ForceDispatchCardData(wCurrentUser) == false {
					return false
				}
			}
			room.SumWindCount -= room.WindCount[wCurrentUser]
			room.WindCount[wCurrentUser] = room.TempWinCount[wCurrentUser]
			room.TempWinCount[wCurrentUser] = 0

		}
		if room.SumWindCount <= 0 {
			break
		}
	}

	//动作分析
	bAroseAction = false
	for i := 0; i < room.UserCnt; i++ {
		//庄家判断
		if i == room.BankerUser {
			//杠牌判断
			GangCardResult := &TagGangCardResult{}
			room.UserAction[i] |= room.gameLogic.AnalyseGangCard(room.CardIndex[i], nil, 0, GangCardResult)

			//胡牌判断
			ChiHuResult := &TagChiHuResult{}

			//八花胡牌类型
			if room.UserWindCount[room.BankerUser] == 8 {
				room.UserAction[i] |= WIK_CHI_HU

			}
			cbHandTai := 0
			cbHandFeng := 0
			room.UserAction[i] |= room.gameLogic.AnalyseChiHuCard(room.CardIndex[i], nil, 0, 0, 0,
				ChiHuResult, &cbHandTai, &cbHandFeng, room.FengQuan, 0, true)
			cbHandTai += (room.Zfb[i] + room.UserWindCount[i] + room.Dnxb[i] + cbHandFeng)

			if (ChiHuResult.ChiHuKind&CHK_JI_HU != 0) &&
				(ChiHuResult.ChiHuRight == 0 || ChiHuResult.ChiHuRight == CHR_QIANG_GANG) &&
				(cbHandTai == 0) {
				room.UserAction[i] &= (^WIK_CHI_HU)
			}

			if room.UserWindCount[room.BankerUser] == 8 {
				//八花胡牌权位
				if (ChiHuResult.ChiHuKind)&(^CHK_BA_HUA) != 0 {
					ChiHuResult.ChiHuRight |= CHR_BA_HUA
				} else {
					//八花胡牌类型不算自摸
					ChiHuResult.ChiHuRight &= ^CHR_ZI_MO

				}
			}
			//四花胡牌权位
			if room.UserWindCount[room.BankerUser] == 4 {
				bHaveSihuaCount := 0
				for bTemp := 0; bTemp < 4; bTemp++ {
					if (room.UserWindData[room.BankerUser][bTemp] - 0x38) < 4 {
						bHaveSihuaCount++
					}
				}
				if (bHaveSihuaCount == 4) || (bHaveSihuaCount == 0) {
					ChiHuResult.ChiHuRight |= CHR_SI_HUA
				}
			}
		}

		//状态设置
		if (bAroseAction == false) && (i != room.BankerUser) && (room.UserAction[i] != WIK_NULL) {
			bAroseAction = true
			room.ResumeUser = room.CurrentUser
			room.CurrentUser = INVALID_CHAIR
		}
	}

	//当前用户
	room.CurrentUser = room.BankerUser

	//构造数据
	GameStart.SiceCount = room.SiceCount
	GameStart.BankerUser = room.BankerUser
	GameStart.InitialBankerUser = room.InitialBankerUser
	GameStart.FengQuan = room.FengQuan
	GameStart.CurrentUser = room.CurrentUser
	GameStart.SunWindCount = room.SumWindCount
	GameStart.LeftCardCount = room.LeftCardCount
	GameStart.First = false

	//发送数据
	for i := 0; i < room.UserCnt; i++ {
		//设置变量
		GameStart.UserAction = room.UserAction[i]
		room.gameLogic.SwitchToCardData2(room.CardIndex[i], GameStart.CardData)
		//发送数据
		room.SendMsgAll(GameStart)
	}
	return true
}

//游戏结束
func (room *Room) OnEventGameConclude(wChairID int, user *client.User, cbReason int) bool {
	template, ok := base.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	if !ok {
		log.Error("at OnEventGameConclude not foud tempplate")
		return false
	}
	if (template.ServerType & GAME_GENRE_MATCH) != 0 {
		//room.KillGameTimer(IDI_CHECK_TABLE);
	}

	switch cbReason {
	case GER_NORMAL: //常规结束
		//变量定义
		GameEnd := &mj_xs_msg.G2C_MJXS_GameEnd{}

		//结束信息
		GameEnd.ProvideUser = wChairID
		GameEnd.ChiHuCard = room.ChiHuCard
		for i := 0; i < room.PlayerCount; i++ {
			//胡牌类型
			GameEnd.ChiHuKind[i] = room.ChiHuResult[i].ChiHuKind
			GameEnd.ChiHuRight[i] = room.ChiHuResult[i].ChiHuRight
		}

		//统计积分
		if room.ProvideUser != INVALID_CHAIR {

			//自摸类型
			if room.ChiHuResult[room.ProvideUser].ChiHuKind != CHK_NULL {
				//仅有8花
				if ((room.ChiHuResult[room.ProvideUser].ChiHuKind) & (^CHK_BA_HUA)) == 0 {
					GameEnd.ChiHuCard = 255
				}

				//风位计算
				cbFengcount := 0
				cbZi := 0
				cbHua := 0
				for nTemp := 0; nTemp < room.WeaveItemCount[room.ProvideUser]; nTemp++ {
					//中發白
					card := room.WeaveItemArray[room.ProvideUser][nTemp].CenterCard
					if room.gameLogic.isZFB(card) {
						cbZi++
					}

					//圈风
					qf := room.gameLogic.isFeng(card, room.FengQuan)
					if qf {
						cbFengcount += 1
					}

					if room.gameLogic.isWeiFeng(card, room.ProvideUser, room.PlayerCount, room.BankerUser) {
						//位风
						cbFengcount += 1
					}
				}
				//正花 野花
				for bTemp := 0; bTemp < room.UserWindCount[room.ProvideUser]; bTemp++ {
					if room.gameLogic.isZhengHua(room.UserWindData[room.ProvideUser][bTemp], room.ProvideUser, room.PlayerCount, room.BankerUser) {
						cbHua += 2
					} else { //其他花牌 野花
						cbHua++
					}
				}

				GameEnd.Feng[room.ProvideUser] = cbFengcount + room.RealDnxb[room.ProvideUser]
				GameEnd.Zi[room.ProvideUser] = cbZi + room.RealZfb[room.ProvideUser]
				GameEnd.Hua[room.ProvideUser] = cbHua

				bHasThree := false
				for j := 0; j < room.UserCnt; j++ {
					if j == room.ProvideUser {
						continue
					}
					if room.IsAfford(room.ProvideUser, j) > 1 {
						bHasThree = true
						break
					}
				}
				if !bHasThree {
					//循环累计
					for i := 0; i < room.PlayerCount; i++ {
						//赢家过滤
						if room.ChiHuResult[i].ChiHuKind != CHK_NULL {
							continue
						}
						cbAllTai := room.gameLogic.CalScore(room.ChiHuResult[room.ProvideUser]) + cbFengcount + cbHua + cbZi + room.RealZfb[room.ProvideUser] + room.RealDnxb[room.ProvideUser]
						cbAllTimes := 1

						//自摸算分
						GameEnd.GameScore[i] -= 1 * template.CellScore * cbAllTai * cbAllTimes * room.IsAfford(room.ProvideUser, i)
						GameEnd.GameScore[room.ProvideUser] += 1 * template.CellScore * cbAllTai * cbAllTimes * room.IsAfford(room.ProvideUser, i)

						//杠上开花
						if room.ChiHuResult[room.ProvideUser].ChiHuRight&CHR_GANG_FLOWER != 0 {
							//找放杠
							if (room.WeaveItemArray[room.ProvideUser][room.WeaveItemCount[room.ProvideUser]-1].ProvideUser == i) && (i != room.ProvideUser) {
								GameEnd.GameScore[i] -= 1 * template.CellScore * cbAllTai * cbAllTimes * 2 * room.IsAfford(room.ProvideUser, i)
								GameEnd.GameScore[room.ProvideUser] += 1 * template.CellScore * cbAllTai * 2 * cbAllTimes * room.IsAfford(room.ProvideUser, i)
							}
						}
					}
					GameEnd.All = GameEnd.GameScore[room.ProvideUser]
				} else {
					//循环累计
					for i := 0; i < room.PlayerCount; i++ {
						//赢家过滤
						if room.ChiHuResult[i].ChiHuKind != CHK_NULL {
							continue
						}
						cbAllTai := room.gameLogic.CalScore(room.ChiHuResult[room.ProvideUser]) + cbFengcount + cbHua + cbZi + room.RealZfb[room.ProvideUser] + room.RealDnxb[room.ProvideUser]
						cbAllTimes := 1
						cbAfford := room.IsAfford(room.ProvideUser, i)
						if cbAfford == 1 {
							cbAfford = 0
						}
						if cbAfford == 2 || cbAfford == 4 {
							cbAfford = 5
						}

						//自摸算分
						GameEnd.GameScore[i] -= 1 * template.CellScore * cbAllTai * cbAllTimes * cbAfford
						GameEnd.GameScore[room.ProvideUser] += 1 * template.CellScore * cbAllTai * cbAllTimes * cbAfford

						//杠上开花
						if room.ChiHuResult[room.ProvideUser].ChiHuRight&CHR_GANG_FLOWER != 0 {
							//找放杠
							if (room.WeaveItemArray[room.ProvideUser][room.WeaveItemCount[room.ProvideUser]-1].ProvideUser == i) && (i != room.ProvideUser) {
								GameEnd.GameScore[i] -= 1 * template.CellScore * cbAllTai * cbAllTimes * 2 * cbAfford
								GameEnd.GameScore[room.ProvideUser] += 1 * template.CellScore * cbAllTai * 2 * cbAllTimes * cbAfford
							}
						}
					}
					GameEnd.All = GameEnd.GameScore[room.ProvideUser]
				}

				//庄家设置
				if room.BankerUser == room.ProvideUser {
					room.BankerUser = room.ProvideUser
					room.Change = false
				} else {
					room.BankerUser = (room.BankerUser - 1 + room.PlayerCount) % room.PlayerCount
					room.Change = true
				}
				//风圈设置
				if (room.BankerUser == room.InitialBankerUser) && (room.Change == true) {
					room.FengQuan = (room.FengQuan + 1) % room.PlayerCount
				}
			}
			//捉炮类型
			if room.ChiHuResult[room.ProvideUser].ChiHuKind == CHK_NULL {
				//循环累计
				for i := 0; i < room.PlayerCount; i++ {
					//输家过滤
					if room.ChiHuResult[i].ChiHuKind == CHK_NULL {
						continue
					}
					//计算
					cbFengcount := 0
					cbZi := 0
					cbHua := 0
					for nTemp := 0; nTemp < room.WeaveItemCount[room.ProvideUser]; nTemp++ {
						//中發白
						card := room.WeaveItemArray[room.ProvideUser][nTemp].CenterCard
						if room.gameLogic.isZFB(card) {
							cbZi++
						}

						//圈风
						qf := room.gameLogic.isFeng(card, room.FengQuan)
						if qf {
							cbFengcount += 1
						}

						if room.gameLogic.isWeiFeng(card, room.ProvideUser, room.PlayerCount, room.BankerUser) {
							//位风
							cbFengcount += 1
						}
					}
					//正花 野花
					for bTemp := 0; bTemp < room.UserWindCount[room.ProvideUser]; bTemp++ {
						if room.gameLogic.isZhengHua(room.UserWindData[room.ProvideUser][bTemp], room.ProvideUser, room.PlayerCount, room.BankerUser) {
							cbHua += 2
						} else { //其他花牌 野花
							cbHua++
						}
					}

					cbAllTai := room.gameLogic.CalScore(room.ChiHuResult[i]) + cbFengcount + cbHua + cbZi + room.RealDnxb[i] + room.RealZfb[i]

					//数据校验
					cbAllTimes := 1
					if room.ChiHuResult[i].ChiHuRight&CHR_QIANG_GANG != 0 {
						cbAllTimes = 5
					}

					GameEnd.Feng[i] = cbFengcount + room.RealDnxb[i]
					GameEnd.Hua[i] = cbHua
					GameEnd.Zi[i] = cbZi + room.RealZfb[i]

					//放炮算分
					GameEnd.GameScore[room.ProvideUser] -= 2 * template.CellScore * cbAllTai * cbAllTimes * room.IsAfford(room.ProvideUser, i)
					GameEnd.GameScore[i] += 2 * template.CellScore * cbAllTai * cbAllTimes * room.IsAfford(room.ProvideUser, i)
					GameEnd.All += GameEnd.GameScore[i]

					//3包非放炮 同样输
					for j := 0; j < room.UserCnt; j++ {
						//过滤赢家 和放炮者
						if j == i || j == room.ProvideUser {
							continue
						}
						cbAfford := room.IsAfford(j, i)
						if cbAfford != 1 {
							//包输钱
							GameEnd.GameScore[j] -= 2 * template.CellScore * cbAllTai * cbAllTimes * room.IsAfford(room.ProvideUser, i)
							GameEnd.GameScore[i] += 2 * template.CellScore * cbAllTai * cbAllTimes * room.IsAfford(room.ProvideUser, i)
							GameEnd.All += GameEnd.GameScore[i]
						}
					}
				}

				//庄家设置
				if room.ChiHuResult[room.BankerUser].ChiHuKind == CHK_NULL {
					room.BankerUser = (room.BankerUser - 1 + room.PlayerCount) % room.PlayerCount
					room.Change = true
				} else {
					room.BankerUser = room.BankerUser
					room.Change = false
				}
				//风圈设置
				if (room.BankerUser == room.InitialBankerUser) && (room.Change == true) {
					room.FengQuan = (room.FengQuan + 1) % room.PlayerCount
				}
			}
		} else {
			//荒庄处理
			GameEnd.ChiHuCard = 0
			GameEnd.HaiDiCard = room.ProvideCard
			//设置庄家
			if room.GangCount > 0 {
				room.BankerUser = (room.BankerUser + room.PlayerCount - 1) % room.PlayerCount
				room.Change = true
			}
			//风圈设置
			if (room.BankerUser == room.InitialBankerUser) && (room.Change == true) {
				room.FengQuan = (room.FengQuan + 1) % room.PlayerCount
			}
		}
		//拷贝扑克
		for i := 0; i < room.PlayerCount; i++ {
			GameEnd.CardCount[i] = room.gameLogic.SwitchToCardData2(room.CardIndex[i], GameEnd.CardData[i])
		}

		//计算税收
		if template.ServerType&GAME_GENRE_GOLD != 0 {
			for i := 0; i < room.PlayerCount; i++ {
				if GameEnd.GameScore[i] >= 1000 {
					lTempTax := GameEnd.GameScore[i] * template.RevenueRatio / 1000
					GameEnd.GameTax += lTempTax
					GameEnd.GameScore[i] -= lTempTax
				}
			}
		}
		if template.ServerType&GAME_GENRE_GOLD != 0 { //积分过滤
			lDifference := float64(0)
			for i, u := range room.Users {
				if (u.Score <= int64(math.Abs(float64(GameEnd.GameScore[i])))) && (GameEnd.GameScore[i] < 0) && (u.Score >= 0) {
					lDifference += math.Abs(float64(GameEnd.GameScore[i])) - float64(u.Score)
					GameEnd.GameScore[i] = -int(u.Score)
				}
			}

			//不是荒庄
			if room.ProvideUser != INVALID_CHAIR {
				for i := 0; i < room.PlayerCount; i++ {
					if GameEnd.GameScore[i] > 0 {
						lTemp := math.Min(float64(GameEnd.GameScore[i]), lDifference)
						GameEnd.GameScore[i] -= int(lTemp)
						lDifference -= lTemp
					}
				}
			}
		}

		//发送信息
		room.SendMsgAll(GameEnd)

		//修改积分
		ScoreInfo := make([]*msg.TagScoreInfo, room.UserCnt)
		for i := 0; i < room.UserCnt; i++ {
			ScoreInfo[i].Score = GameEnd.GameScore[i]
		}
		//统计积分

		room.WriteTableScore(ScoreInfo, room.UserCnt, 1)

		//结束游戏
		room.GameEnd()
		return true
	case GER_USER_LEAVE: //用户强退
		//变量定义
		GameEnd := mj_xs_msg.G2C_MJXS_GameEnd{}

		//设置变量
		GameEnd.ChiHuCard = 255
		GameEnd.ProvideUser = INVALID_CHAIR
		GameEnd.GameScore[wChairID] = -12 * template.CellScore
		if template.ServerType&GAME_GENRE_GOLD != 0 {
			//积分过滤
			user := room.GetUserByChairId(wChairID)

			if user.Score < int64(math.Abs(float64(GameEnd.GameScore[wChairID]))) {
				GameEnd.GameScore[wChairID] = -int(user.Score)
			}
		}

		//通知消息
		szMessage := fmt.Sprintf("由于 [ %s ] 离开游戏，游戏结束", user.NickName)
		room.SendMsgAll(szMessage) //todo sysmsg

		//发送信息
		room.SendMsgAll(GameEnd)

		//修改积分
		ScoreInfo := &msg.TagScoreInfo{}
		ScoreInfo.Score = GameEnd.GameScore[wChairID]
		room.WriteTableScore([]*msg.TagScoreInfo{ScoreInfo}, room.UserCnt, 1)

		//结束游戏
		room.GameEnd()
		return true
	case GER_DISMISS: //用户强退
		//变量定义
		GameEnd := mj_xs_msg.G2C_MJXS_GameEnd{}

		//设置变量
		GameEnd.ChiHuCard = 255
		GameEnd.ProvideUser = INVALID_CHAIR
		GameEnd.GameScore[wChairID] = -12 * template.CellScore
		user := room.GetUserByChairId(wChairID)
		if template.ServerType&GAME_GENRE_GOLD != 0 {
			//积分过滤
			if user.Score < int64(math.Abs(float64(GameEnd.GameScore[wChairID]))) {
				GameEnd.GameScore[wChairID] = -int(user.Score)
			}
		}

		//通知消息

		szMessage := fmt.Sprintf("由于 [ %s ] 超时，游戏结束", user.NickName)

		room.SendMsgAll(szMessage) //todo

		//发送信息
		room.SendMsgAll(GameEnd)

		//修改积分
		ScoreInfo := &msg.TagScoreInfo{}
		ScoreInfo.Score = GameEnd.GameScore[wChairID]
		room.WriteTableScore([]*msg.TagScoreInfo{ScoreInfo}, room.UserCnt, 1)

		//结束游戏
		room.GameEnd()
		return true
	}

	return false
}

//todo
func (room *Room) GameEnd() {
	room.PlayCount++
}

func (room *Room) IsAfford(wUserProvider, wUserAccept int) int {
	cbAfford := 0 //1：不是包关系 2：包关系
	cbCount := 0
	//供三包
	for cbTemp := 0; cbTemp < room.WeaveItemCount[wUserAccept]; cbTemp++ {
		if room.WeaveItemArray[wUserAccept][cbTemp].ProvideUser == wUserProvider {
			cbCount++
		}

	}
	if cbCount >= 3 {
		cbAfford += 2
	}

	//求三包
	cbCount = 0
	for cbTemp := 0; cbTemp < room.WeaveItemCount[wUserProvider]; cbTemp++ {
		if room.WeaveItemArray[wUserProvider][cbTemp].ProvideUser == wUserAccept {
			cbCount++
		}
	}
	if cbCount >= 3 {
		cbAfford += 2
	}

	if cbAfford > 1 {
		return cbAfford
	}
	return 1
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
		if wCenterUser == i {
			continue
		}

		//出牌类型
		if EstimatKind == EstimatKind_OutCard {
			//吃碰判断
			if u.UserLimit&LimitPeng == 0 && room.LeftCardCount >= room.RemainCardCount+1 {
				//碰牌判断
				if cbCenterCard != room.EnjoinPengCard[i] {
					room.UserAction[i] |= room.gameLogic.EstimatePengCard(room.CardIndex[i], cbCenterCard)
				}

				//禁止碰拍
				if room.UserAction[i]&WIK_PENG != 0 {
					room.EnjoinPengCard[i] = cbCenterCard
				}

				//吃牌判断
				wEatUser := (wCenterUser + room.PlayerCount - 1) % room.PlayerCount
				if wEatUser == i {
					room.UserAction[i] |= room.gameLogic.EstimateEatCard(room.CardIndex[i], cbCenterCard)
				}
			}
			//杠牌判断
			room.UserAction[i] |= room.gameLogic.EstimateGangCard(room.CardIndex[i], cbCenterCard)
		}

		//胡牌判断
		if u.UserLimit&LimitChiHu == 0 {
			//牌型权位
			wChiHuRight := 0
			if room.GangStatus == true {
				wChiHuRight |= CHR_QIANG_GANG
			}

			if (room.SendCardCount == room.AllWindCount+1) && (room.OutCardCount == room.AllWindCount+1) {
				wChiHuRight |= CHR_DI
			}

			if (room.SendCardCount == room.AllWindCount+1) && (room.OutCardCount == room.AllWindCount) {
				wChiHuRight |= CHR_TIAN
			}

			//吃胡判断
			ChiHuResult := &TagChiHuResult{}
			cbWeaveCount := room.WeaveItemCount[i]
			if room.UserWindCount[i] == 8 {
				room.UserAction[i] |= WIK_CHI_HU
			}
			cbHandTai := 0
			cbHandFeng := 0
			if room.EnjoinHuCard[i] != cbCenterCard {
				room.UserAction[i] |= room.gameLogic.AnalyseChiHuCard(room.CardIndex[i], room.WeaveItemArray[i], cbWeaveCount, cbCenterCard, wChiHuRight,
					ChiHuResult, &cbHandTai, &cbHandFeng, room.FengQuan, (room.BankerUser-i+room.UserCnt)%room.UserCnt, false)
			}

			cbHandTai += (room.Zfb[i] + room.UserWindCount[i] + room.Dnxb[i] + cbHandFeng)

			if ChiHuResult.ChiHuKind == CHK_JI_HU && (ChiHuResult.ChiHuRight == 0 || ChiHuResult.ChiHuRight == CHR_QIANG_GANG) && (cbHandTai == 0) {
				room.UserAction[i] &= (^WIK_CHI_HU)
			}

			//禁止胡牌
			if room.UserAction[i]&WIK_CHI_HU != 0 {
				room.EnjoinHuCard[i] = cbCenterCard
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
		room.SendOperateNotify()
		return true
	}

	return false
}

//派发扑克
func (room *Room) DispatchCardData(wCurrentUser int, bGang bool) bool {
	//状态效验
	if wCurrentUser == INVALID_CHAIR {
		return false
	}

	if room.SendStatus == false {
		log.Error("at DispatchCardData f room.SendStatus == Not_Send")
		return false
	}

	//丢弃扑克
	if (room.OutCardUser != INVALID_CHAIR) && (room.OutCardData != 0) {
		if len(room.DiscardCard[room.OutCardUser]) < 1 {
			room.DiscardCard[room.OutCardUser] = make([]int, 60)
		}
		room.DiscardCard[room.OutCardUser][room.DiscardCount[room.OutCardUser]] = room.OutCardData
		room.DiscardCount[room.OutCardUser]++
	}

	//海底判断
	if (room.LeftCardCount <= room.RemainCardCount+1) && bGang {
		//发送扑克
		room.SendCardCount++
		room.LeftCardCount--
		room.SendCardData = room.RepertoryCard[room.LeftCardCount]

		//判断花牌
		if room.LeftCardCount > 0 {

			//判断花牌
			for (room.SendCardData >= 0x38) && (room.SendCardData <= 0x3F) {

				//强制出牌
				if room.gameLogic.IsValidCard(room.SendCardData) == false {
					return false
				}
				//插入数据
				room.CardIndex[wCurrentUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]++

				//构造数据
				SendCard := mj_xs_msg.G2C_MJXS_SendCard{}
				SendCard.CurrentUser = wCurrentUser
				SendCard.ActionMask = room.UserAction[wCurrentUser]
				SendCard.CardData = room.SendCardData
				SendCard.Gang = true

				//发送数据
				room.SendMsgAll(SendCard)
				//room.pITableFrame->SendLookonData(INVALID_CHAIR,SUB_S_FORCE_SEND_CARD,&SendCard,sizeof(SendCard));

				//删除扑克
				if room.gameLogic.RemoveCard(room.CardIndex[wCurrentUser], room.SendCardData) == false {
					return false
				}

				//出牌记录
				room.OutCardCount++
				room.OutCardUser = wCurrentUser
				room.OutCardData = room.SendCardData
				room.AllWindCount += 1
				room.UserWindData[wCurrentUser][room.UserWindCount[wCurrentUser]] = room.SendCardData
				room.UserWindCount[wCurrentUser]++

				//状态判断
				if room.Status == RoomStatusStarting {
					//规则判断
					if room.AllWindCount < 8 {
						room.RemainCardCount = 14 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
					} else {
						room.RemainCardCount = 16 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
					}
				}

				//构造数据
				OutCard := mj_xs_msg.G2C_MJXS_OutCard{}
				OutCard.OutCardUser = wCurrentUser
				OutCard.OutCardData = room.SendCardData

				//发送消息
				room.SendMsgAll(OutCard)
				//room.pITableFrame->SendLookonData(INVALID_CHAIR,SUB_S_FORCE_OUT_CARD,&OutCard,sizeof(OutCard));

				//重新发牌
				room.SendCardCount++
				room.LeftCardCount--
				room.SendCardData = room.RepertoryCard[room.LeftCardCount]
			}
		}
	}

	if (room.LeftCardCount <= room.RemainCardCount+1) && bGang == false {
		//设置变量
		room.ResumeUser = wCurrentUser
		room.CurrentUser = wCurrentUser
		room.ProvideUser = INVALID_CHAIR
		room.ProvideCard = room.SendCardData
		room.OnEventGameConclude(room.ProvideUser, nil, GER_NORMAL)
		return true
	}

	//荒庄结束
	if room.LeftCardCount <= 0 {
		room.ChiHuCard = 0
		room.ProvideUser = INVALID_CHAIR
		room.OnEventGameConclude(room.ProvideUser, nil, GER_NORMAL)
		return true
	}

	//发牌处理
	if room.SendStatus == true {
		//发送扑克
		room.SendCardCount++
		room.LeftCardCount--
		room.SendCardData = room.RepertoryCard[room.LeftCardCount]

		//判断花牌
		if room.LeftCardCount > 0 { ////此时可能还有杠

			//判断花牌
			for (room.SendCardData >= 0x38) && (room.SendCardData <= 0x3F) {

				//强制出牌
				if room.gameLogic.IsValidCard(room.SendCardData) == false {
					return false
				}

				//插入数据
				room.CardIndex[wCurrentUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]++

				//构造数据
				SendCard := mj_xs_msg.G2C_MJXS_SendCard{}
				SendCard.CurrentUser = wCurrentUser
				SendCard.ActionMask = room.UserAction[wCurrentUser]
				SendCard.CardData = room.SendCardData
				SendCard.Gang = true

				//发送数据
				room.SendMsgAll(SendCard)
				//room.pITableFrame->SendLookonData(INVALID_CHAIR,SUB_S_FORCE_SEND_CARD,&SendCard,sizeof(SendCard));

				//删除扑克
				if room.gameLogic.RemoveCard(room.CardIndex[wCurrentUser], room.SendCardData) == false {
					return false
				}

				//出牌记录
				room.OutCardCount++
				room.OutCardUser = wCurrentUser
				room.OutCardData = room.SendCardData
				room.AllWindCount += 1
				room.UserWindData[wCurrentUser][room.UserWindCount[wCurrentUser]] = room.SendCardData
				room.UserWindCount[wCurrentUser]++

				//状态判断
				if room.Status == RoomStatusStarting {
					//规则判断
					if room.AllWindCount < 8 {
						room.RemainCardCount = 14 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
					} else {
						room.RemainCardCount = 16 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
					}
				}

				//构造数据
				OutCard := &mj_xs_msg.G2C_MJXS_OutCard{}
				OutCard.OutCardUser = wCurrentUser
				OutCard.OutCardData = room.SendCardData

				//发送消息
				room.SendMsgAll(OutCard)
				//room.pITableFrame->SendLookonData(INVALID_CHAIR,SUB_S_FORCE_OUT_CARD,&OutCard,sizeof(OutCard));

				//重新发牌
				room.SendCardCount++
				room.LeftCardCount--
				room.SendCardData = room.RepertoryCard[room.LeftCardCount]
			}
		}

		room.CardIndex[wCurrentUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]++

		//设置变量
		room.ProvideUser = wCurrentUser
		room.ProvideCard = room.SendCardData

		GangCardResult := &TagGangCardResult{}
		room.UserAction[wCurrentUser] |= room.gameLogic.AnalyseGangCard(room.CardIndex[wCurrentUser],
			room.WeaveItemArray[wCurrentUser], room.WeaveItemCount[wCurrentUser], GangCardResult)

		//牌型权位
		wChiHuRight := 0
		if room.GangStatus == true {
			wChiHuRight |= CHR_QIANG_GANG
		}

		//胡牌判断
		ChiHuResult := &TagChiHuResult{}
		//八花胡牌类型
		if room.UserWindCount[wCurrentUser] == 8 {
			room.UserAction[wCurrentUser] |= WIK_CHI_HU
		}
		cbHandTai := 0
		cbHandFeng := 0
		room.UserAction[wCurrentUser] |= room.gameLogic.AnalyseChiHuCard(room.CardIndex[wCurrentUser], room.WeaveItemArray[wCurrentUser], room.WeaveItemCount[wCurrentUser],
			0, wChiHuRight, ChiHuResult, &cbHandTai, &cbHandFeng, room.FengQuan, (room.BankerUser-wCurrentUser+room.UserCnt)%room.UserCnt, true)
		cbHandTai += (room.Zfb[wCurrentUser] + room.UserWindCount[wCurrentUser] + room.Dnxb[wCurrentUser] + cbHandFeng)

		if (ChiHuResult.ChiHuKind&CHK_JI_HU != 0) && (ChiHuResult.ChiHuRight == 0 || ChiHuResult.ChiHuRight == CHR_QIANG_GANG) && (cbHandTai == 0) {
			room.UserAction[wCurrentUser] &= (^WIK_CHI_HU)
		}

		//八花胡牌权位
		if room.UserWindCount[wCurrentUser] == 8 {
			if (room.ChiHuResult[wCurrentUser].ChiHuKind)&(^CHK_BA_HUA) != 0 {
				room.ChiHuResult[wCurrentUser].ChiHuRight |= CHR_BA_HUA
			}
		}

		//四花胡牌权位
		if room.UserWindCount[wCurrentUser] == 4 {
			bHaveSihuaCount := 0
			for bTemp := 0; bTemp < 4; bTemp++ {
				if (room.UserWindData[wCurrentUser][bTemp] - 0x38) < 4 {
					bHaveSihuaCount++
				}
			}
			if (bHaveSihuaCount == 4) || (bHaveSihuaCount == 0) {
				room.ChiHuResult[wCurrentUser].ChiHuRight |= CHR_SI_HUA
			}
		}
	}

	room.CurrentUser = wCurrentUser

	SendCard := &mj_xs_msg.G2C_MJXS_SendCard{}
	SendCard.CurrentUser = wCurrentUser
	SendCard.ActionMask = room.UserAction[wCurrentUser]
	if room.SendStatus == true {
		SendCard.CardData = room.SendCardData
	} else {
		SendCard.CardData = 0
	}
	SendCard.Gang = bGang
	room.SendMsgAll(SendCard)

	//todo
	//room.pITableFrame->SendLookonData(INVALID_CHAIR,SUB_S_SEND_CARD, &SendCard, sizeof(SendCard));

	return true
}

func (room *Room) OnUserOperateCard(user *client.User, cbOperateCard int, cbOperateCode int) bool {

	//被动动作
	if room.CurrentUser == INVALID_CHAIR {
		//效验状态
		if room.Response[user.ChairId] {
			return true
		}
		if room.UserAction[user.ChairId] == WIK_NULL {
			return false
		}
		if (cbOperateCode != WIK_NULL) && ((room.UserAction[user.ChairId] & cbOperateCode) == 0) {
			return false
		}

		//变量定义
		wTargetUser := user.ChairId
		cbTargetAction := cbOperateCode

		//设置变量
		user.UserLimit |= ^LimitGang
		room.Response[wTargetUser] = true
		room.PerformAction[wTargetUser] = cbOperateCode
		if cbOperateCard == 0 {
			room.OperateCard[wTargetUser] = room.ProvideCard
		} else {
			room.OperateCard[wTargetUser] = cbOperateCard
		}

		//执行判断
		for i := 0; i < room.PlayerCount; i++ {
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

		//吃胡等待
		if cbTargetAction == WIK_CHI_HU {
			for i := 0; i < room.UserCnt; i++ {
				if (room.Response[i] == false) && (room.UserAction[i]&WIK_CHI_HU) != 0 {
					return true
				}
			}
		}

		//放弃操作
		if cbTargetAction == WIK_NULL {
			//用户状态
			room.Response = make([]bool, room.UserCnt)
			room.UserAction = make([]int, room.UserCnt)
			room.OperateCard = make([]int, room.UserCnt)
			room.PerformAction = make([]int, room.UserCnt)

			room.DispatchCardData(room.ResumeUser, false)
			return true
		}

		//变量定义
		cbTargetCard := room.OperateCard[wTargetUser]

		//出牌变量
		room.SendStatus = true
		//room.SendCardData = 0
		//room.OutCardUser = INVALID_CHAIR
		//room.OutCardData = 0

		//胡牌操作
		if cbTargetAction == WIK_CHI_HU {
			//结束信息
			wChiHuUser := room.BankerUser
			wChiHuRight := 0
			if room.GangStatus == true {
				wChiHuRight |= CHR_QIANG_GANG
			}

			if (room.SendCardCount == room.AllWindCount) && (room.OutCardCount == room.AllWindCount+1) {
				wChiHuRight |= CHR_DI
			}

			for i := 0; i < room.PlayerCount; i++ {
				wChiHuUser = (room.BankerUser + i) % room.PlayerCount
				//过虑判断
				if i == room.ProvideUser || ((room.PerformAction[wChiHuUser] & WIK_CHI_HU) == 0) {
					continue
				}

				//普通胡牌
				if room.ChiHuCard != 0 {
					//胡牌判断
					cbWeaveItemCount := room.WeaveItemCount[i]
					pWeaveItem := room.WeaveItemArray[i]
					zfb := room.RealZfb[i]
					dnxb := room.RealDnxb[i]
					room.gameLogic.AnalyseChiHuCard(room.CardIndex[i], pWeaveItem, cbWeaveItemCount, room.ChiHuCard, wChiHuRight, room.ChiHuResult[i],
						&zfb, &dnxb, room.FengQuan, (room.BankerUser-i+room.UserCnt)%room.UserCnt, false)
					room.RealZfb[i] = zfb
					room.RealDnxb[i] = dnxb
					//八花胡牌类型
					if room.UserWindCount[i] == 8 {
						room.ChiHuResult[i].ChiHuKind |= CHK_BA_HUA

						//八花胡牌权位
						if (room.ChiHuResult[i].ChiHuKind)&(^CHK_BA_HUA) != 0 {
							room.ChiHuResult[i].ChiHuRight |= CHR_BA_HUA
						}
					}
					//四花胡牌权位

					if room.UserWindCount[i] == 4 {
						bHaveSihuaCount := 0
						for bTemp := 0; bTemp < 4; bTemp++ {
							if (room.UserWindData[i][bTemp] - 0x38) < 4 {
								bHaveSihuaCount++
							}
						}
						if (bHaveSihuaCount == 4) || (bHaveSihuaCount == 0) {
							room.ChiHuResult[i].ChiHuRight |= CHR_SI_HUA
						}
					}

					//插入扑克
					if room.ChiHuResult[i].ChiHuKind != CHK_NULL {
						room.CardIndex[i][room.gameLogic.SwitchToCardIndex(room.ChiHuCard)]++
					}
				}
			}

			//构造结果
			OperateResult := mj_xs_msg.G2C_MJXS_OperateResult{}
			OperateResult.OperateUser = wTargetUser
			if room.ProvideUser == INVALID_CHAIR {
				OperateResult.ProvideUser = wTargetUser
			} else {
				OperateResult.ProvideUser = room.ProvideUser
			}

			OperateResult.OperateCode = WIK_CHI_HU
			OperateResult.OperateCard = cbOperateCard

			//结束游戏
			room.OnEventGameConclude(room.ProvideUser, nil, GER_NORMAL)
			return true
		}

		//组合扑克
		room.WeaveItemCount[wTargetUser]++
		wIndex := room.WeaveItemCount[wTargetUser]
		room.WeaveItemArray[wTargetUser][wIndex].PublicCard = true
		room.WeaveItemArray[wTargetUser][wIndex].CenterCard = cbTargetCard
		room.WeaveItemArray[wTargetUser][wIndex].WeaveKind = cbTargetAction
		if room.ProvideUser == INVALID_CHAIR {
			room.WeaveItemArray[wTargetUser][wIndex].ProvideUser = wTargetUser
		} else {
			room.WeaveItemArray[wTargetUser][wIndex].ProvideUser = room.ProvideUser
		}

		//删除扑克
		switch cbTargetAction {
		case WIK_LEFT: //上牌操作
			//删除扑克
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], []int{cbTargetCard + 1, cbTargetCard + 2}, 2) {
				log.Error("not foud card at Operater")
				return false
			}
			if room.WeaveItemCount[wTargetUser] < 3 {
				//禁止出牌数据和数目
				room.EnjoinCardData[wTargetUser][room.EnjoinCardCount[wTargetUser]] = cbTargetCard
				room.EnjoinCardCount[wTargetUser]++
				//过滤789
				if (cbTargetCard & MASK_VALUE) < 7 {
					room.EnjoinCardData[wTargetUser][room.EnjoinCardCount[wTargetUser]] = cbTargetCard + 3
					room.EnjoinCardCount[wTargetUser]++
				}
			}
			break
		case WIK_RIGHT: //上牌操作
			//删除扑克
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], []int{cbTargetCard - 2, cbTargetCard - 1}, 2) {
				log.Error("not foud card at Operater")
				return false
			}
			if room.WeaveItemCount[wTargetUser] < 3 {
				//禁止出牌数据和数目
				room.EnjoinCardData[wTargetUser][room.EnjoinCardCount[wTargetUser]] = cbTargetCard
				room.EnjoinCardCount[wTargetUser]++
				//过滤 1 2 3
				if (cbTargetCard & MASK_VALUE) < 3 {
					room.EnjoinCardData[wTargetUser][room.EnjoinCardCount[wTargetUser]] = cbTargetCard + 3
					room.EnjoinCardCount[wTargetUser]++
				}
			}

			break
		case WIK_CENTER: //上牌操作
			//删除扑克
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], []int{cbTargetCard - 1, cbTargetCard + 1}, 2) {
				log.Error("not foud card at Operater")
				return false
			}
			if room.WeaveItemCount[wTargetUser] < 3 {
				//禁止出牌数据和数目
				room.EnjoinCardData[wTargetUser][room.EnjoinCardCount[wTargetUser]] = cbTargetCard
				room.EnjoinCardCount[wTargetUser]++
			}

			break
		case WIK_PENG: //碰牌操作
			//删除扑克
			cbRemoveCard := []int{cbTargetCard, cbTargetCard}
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], cbRemoveCard, 2) {
				log.Error("not foud card at Operater")
				return false
			}
			//中发白
			if room.gameLogic.isZFB(cbTargetCard) {
				room.Zfb[wTargetUser]++
			}
			//东南西北
			if (cbTargetCard == 0x31) || (cbTargetCard == 0x32) || (cbTargetCard == 0x33) || (cbTargetCard == 0x34) {
				//圈风
				if ((cbTargetCard & MASK_VALUE) - 1) == room.FengQuan {
					room.Dnxb[wTargetUser]++
				}
				//位风
				if wTargetUser == ((cbTargetCard & MASK_VALUE) - 1) {
					room.Dnxb[wTargetUser]++
				}
			}
			break
		case WIK_GANG: //杠牌操作
			//删除扑克,被动动作只存在放杠
			cbRemoveCard := []int{cbTargetCard, cbTargetCard, cbTargetCard}
			if !room.gameLogic.RemoveCardByCnt(room.CardIndex[wTargetUser], cbRemoveCard, int(len(cbRemoveCard))) {
				log.Error("not foud card at Operater")
				return false
			}
			if room.gameLogic.isZFB(cbTargetCard) {
				room.Zfb[wTargetUser]++
			}
			//东南西北
			if (cbTargetCard == 0x31) || (cbTargetCard == 0x32) || (cbTargetCard == 0x33) || (cbTargetCard == 0x34) {
				//圈风
				if ((cbTargetCard & MASK_VALUE) - 1) == room.FengQuan {
					room.Dnxb[wTargetUser]++
				}
				//位风
				if wTargetUser == ((cbTargetCard & MASK_VALUE) - 1) {
					room.Dnxb[wTargetUser]++
				}
			}
			room.GangCount++
			//预留扑克
			//状态判断
			if room.Status == RoomStatusStarting {
				//规则判断
				if room.AllWindCount < 8 {
					room.RemainCardCount = 14 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
				} else {
					room.RemainCardCount = 16 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
				}
			}
			break
		}

		//脱牌处理 吃碰杠三次允许脱牌
		if room.WeaveItemCount[wTargetUser] > 2 {
			room.EnjoinCardCount[wTargetUser] = 0
			room.EnjoinCardData[wTargetUser] = make([]int, room.UserCnt)
		}

		//构造结果
		OperateResult := &mj_xs_msg.G2C_MJXS_OperateResult{}
		OperateResult.OperateUser = wTargetUser
		OperateResult.OperateCard = cbTargetCard
		OperateResult.OperateCode = cbTargetAction
		if room.ProvideUser == INVALID_CHAIR {
			OperateResult.ProvideUser = wTargetUser
		} else {
			OperateResult.ProvideUser = room.ProvideUser
		}
		room.SendMsgAll(OperateResult)

		//设置用户
		room.CurrentUser = wTargetUser

		//杠牌处理
		if cbTargetAction == WIK_GANG {
			//效验动作
			bAroseAction := room.EstimateUserRespond(wTargetUser, cbTargetCard, EstimatKind_GangCard)

			//发送扑克
			if bAroseAction == false {
				room.DispatchCardData(wTargetUser, true)
			}

			return true
		}

		//动作判断
		if room.LeftCardCount >= room.RemainCardCount+1 {
			//杠牌判断
			GangCardResult := &TagGangCardResult{}
			room.UserAction[room.CurrentUser] |= room.gameLogic.AnalyseGangCard(room.CardIndex[room.CurrentUser],
				room.WeaveItemArray[room.CurrentUser], room.WeaveItemCount[room.CurrentUser], GangCardResult)

			//结果处理
			if GangCardResult.CardCount > 0 {
				//设置变量
				room.UserAction[room.CurrentUser] |= WIK_GANG

				//发送动作
				room.SendOperateNotify()
			}
		}
		return true

	} else { //主动动作
		//扑克效验

		if (cbOperateCode != WIK_NULL) && (cbOperateCode != WIK_CHI_HU) && (!room.gameLogic.IsValidCard(cbOperateCard)) {
			return false
		}

		//设置变量
		room.UserAction[room.CurrentUser] = WIK_NULL
		room.PerformAction[room.CurrentUser] = WIK_NULL

		//执行动作
		switch cbOperateCode {
		case WIK_GANG: //杠牌操作
			//变量定义
			cbWeaveIndex := 0xFF
			cbCardIndex := room.gameLogic.SwitchToCardIndex(cbOperateCard)
			bpublic := true
			//杠牌处理
			if room.CardIndex[user.ChairId][cbCardIndex] == 1 {
				//寻找组合
				for i := 0; i < room.WeaveItemCount[user.ChairId]; i++ {
					cbWeaveKind := room.WeaveItemArray[user.ChairId][i].WeaveKind
					cbCenterCard := room.WeaveItemArray[user.ChairId][i].CenterCard
					if (cbCenterCard == cbOperateCard) && (cbWeaveKind == WIK_PENG) {
						cbWeaveIndex = i
						break
					}
				}

				//效验动作
				if cbWeaveIndex == 0xFF {
					return false
				}

				//组合扑克
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].PublicCard = true
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].WeaveKind = cbOperateCode
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].CenterCard = cbOperateCard
				bpublic = true
			} else {
				//扑克效验

				if room.CardIndex[user.ChairId][cbCardIndex] != 4 {
					return false
				}

				//设置变量
				room.WeaveItemCount[user.ChairId]++
				cbWeaveIndex := room.WeaveItemCount[user.ChairId]
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].PublicCard = false
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].ProvideUser = user.ChairId
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].WeaveKind = cbOperateCode
				room.WeaveItemArray[user.ChairId][cbWeaveIndex].CenterCard = cbOperateCard
				bpublic = false
			}

			//中发白
			if room.gameLogic.isZFB(cbOperateCard) {
				room.Zfb[user.ChairId]++
			}
			//东南西北
			if (cbOperateCard == 0x31) || (cbOperateCard == 0x32) || (cbOperateCard == 0x33) || (cbOperateCard == 0x34) {
				//圈风
				if ((cbOperateCard & MASK_VALUE) - 1) == room.FengQuan {
					room.Dnxb[user.ChairId]++
				}

				//位风
				if user.ChairId == ((cbOperateCard & MASK_VALUE) - 1) {
					room.Dnxb[user.ChairId]++
				}
			}

			//杠牌次数
			room.GangCount++
			//预留扑克
			//状态判断
			if room.Status == RoomStatusStarting {
				//规则判断
				if room.AllWindCount < 8 {
					room.RemainCardCount = 14 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
				} else {
					room.RemainCardCount = 16 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
				}
			}

			//脱牌处理 吃碰杠三次允许脱牌
			if room.WeaveItemCount[user.ChairId] > 2 {
				room.EnjoinCardCount[user.ChairId] = 0
				room.EnjoinCardData[user.ChairId] = make([]int, room.UserCnt)
			}

			//删除扑克
			room.CardIndex[user.ChairId][cbCardIndex] = 0
			//设置状态
			if cbOperateCode == WIK_GANG {
				room.GangStatus = true
				//room.bEnjoinChiPeng[wChairID]=true;
			}

			//构造结果
			OperateResult := &mj_hz_msg.G2C_HZMJ_OperateResult{}
			OperateResult.OperateUser = user.ChairId
			OperateResult.ProvideUser = user.ChairId
			OperateResult.OperateCode = cbOperateCode
			OperateResult.OperateCard[0] = cbOperateCard

			//发送消息
			room.SendMsgAll(OperateResult)
			//room.pITableFrame->SendLookonData(INVALID_CHAIR, SUB_S_OPERATE_RESULT, &OperateResult, sizeof(OperateResult));

			//效验动作
			bAroseAction := false
			if bpublic {
				bAroseAction = room.EstimateUserRespond(user.ChairId, cbOperateCard, EstimatKind_GangCard)
			}

			//发送扑克
			if !bAroseAction {
				room.DispatchCardData(user.ChairId, true)
			}
			return true
		case WIK_CHI_HU: //自摸
			//吃牌权位
			wChiHuRight := 0
			if room.GangStatus == true {
				wChiHuRight |= CHR_QIANG_GANG
			}

			if (room.SendCardCount == room.AllWindCount+1) && (room.OutCardCount == room.AllWindCount+1) {
				wChiHuRight |= CHR_DI
			}

			if (room.SendCardCount == room.AllWindCount+1) && (room.OutCardCount == room.AllWindCount) {
				room.ProvideUser = room.CurrentUser
				wChiHuRight |= CHR_TIAN
			}
			if room.LeftCardCount <= room.RemainCardCount {
				room.ProvideUser = room.CurrentUser
				wChiHuRight |= CHR_HAI_DI
			}
			//普通胡牌
			cbWeaveItemCount := room.WeaveItemCount[user.ChairId]
			pWeaveItem := room.WeaveItemArray[user.ChairId]
			if !room.gameLogic.RemoveCard(room.CardIndex[user.ChairId], room.SendCardData) {
				log.Error("not foud card at Operater")
				return false
			}
			zfb := room.RealZfb[user.ChairId]
			dnxb := room.RealZfb[user.ChairId]
			room.gameLogic.AnalyseChiHuCard(room.CardIndex[user.ChairId], pWeaveItem, cbWeaveItemCount, room.SendCardData, wChiHuRight, room.ChiHuResult[user.ChairId],
				&zfb, &dnxb, room.FengQuan, (room.BankerUser-user.ChairId+room.UserCnt)%room.UserCnt, true)
			room.RealZfb[user.ChairId] = zfb
			room.RealDnxb[user.ChairId] = dnxb

			room.CardIndex[user.ChairId][room.gameLogic.SwitchToCardIndex(room.SendCardData)]++

			//四花胡牌权位
			if room.UserWindCount[user.ChairId] == 4 {
				bHaveSihuaCount := 0
				for bTemp := 0; bTemp < 4; bTemp++ {
					if (room.UserWindData[user.ChairId][bTemp] - 0x38) < 4 {
						bHaveSihuaCount++
					}
				}
				if (bHaveSihuaCount == 4) || (bHaveSihuaCount == 0) {
					room.ChiHuResult[user.ChairId].ChiHuRight |= CHR_SI_HUA
				}
			}

			//八花胡牌类型
			if room.UserWindCount[user.ChairId] == 8 {
				room.ChiHuResult[user.ChairId].ChiHuKind |= CHK_BA_HUA

				//八花胡牌权位
				if room.ChiHuResult[user.ChairId].ChiHuKind&(^CHK_BA_HUA) != 0 {
					room.ChiHuResult[user.ChairId].ChiHuRight |= CHR_BA_HUA
				} else {
					//八花胡牌类型不算自摸
					room.ChiHuResult[user.ChairId].ChiHuRight &= ^CHR_ZI_MO
				}
			}

			//结束信息
			room.ChiHuCard = room.ProvideCard

			//构造结果
			OperateResult := mj_xs_msg.G2C_MJXS_OperateResult{}
			OperateResult.OperateUser = user.ChairId
			OperateResult.ProvideUser = user.ChairId
			OperateResult.OperateCode = WIK_CHI_HU
			OperateResult.OperateCard = cbOperateCard
			room.SendMsgAll(OperateResult)
			//结束游戏
			room.OnEventGameConclude(room.ProvideUser, nil, GER_NORMAL)
			return true
		}
		return true
	}
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

	//设置变量
	room.SendStatus = true
	user.UserLimit |= ^LimitChiHu
	room.UserAction[wChairID] = WIK_NULL
	room.PerformAction[wChairID] = WIK_NULL
	room.EnjoinCardCount[wChairID] = 0
	for i, _ := range room.EnjoinCardData[wChairID] {
		if room.EnjoinCardData[wChairID][i] != 0 {
			room.EnjoinCardData[wChairID][i] = 0
		}
	}

	//吃胡清空
	for i, u := range room.Users {
		if u == nil {
			continue
		}
		if i == wChairID {
			continue
		}
		user.UserLimit |= ^LimitChiHu
	}

	//出牌记录
	//出牌记录
	room.OutCardCount++
	room.OutCardUser = wChairID
	room.OutCardData = cbCardData

	//构造数据
	OutCard := &mj_xs_msg.G2C_MJXS_OutCard{}
	OutCard.OutCardUser = wChairID
	OutCard.OutCardData = cbCardData

	//发送消息
	room.SendMsgAll(OutCard)
	//room.SendLookonData(INVALID_CHAIR, SUB_S_OUT_CARD, &OutCard, sizeof(OutCard));

	room.ProvideUser = wChairID
	room.ProvideCard = cbCardData
	room.CurrentUser = (wChairID + 1) % room.PlayerCount
	//抢杆设置
	if room.GangStatus == true {
		for i := 0; i < room.UserCnt; i++ {
			if (room.UserAction[i] & WIK_CHI_HU) != 0 {
				break
			}
			if i == room.PlayerCount {
				room.GangStatus = false
			}
		}
	}
	//响应判断
	bAroseAction := room.EstimateUserRespond(wChairID, cbCardData, EstimatKind_OutCard)

	//派发扑克
	if !bAroseAction {
		room.DispatchCardData(room.CurrentUser, false)
	}

	return 0
}

func (room *Room) ForceDispatchCardData(wCurrentUser int) bool {
	//状态效验
	if wCurrentUser == INVALID_CHAIR {
		return false
	}

	//发送扑克
	room.SendCardCount++
	room.LeftCardCount--
	room.SendCardData = room.RepertoryCard[room.LeftCardCount]

	//判断花牌
	if room.SendCardData >= 0x38 && room.SendCardData <= 0x3F {
		room.SumWindCount += 1
		room.WindData[wCurrentUser][room.TempWinCount[wCurrentUser]] = room.SendCardData
		room.TempWinCount[wCurrentUser]++
		room.AllWindCount += 1

		//预留牌数
		//状态判断
		if room.Status == RoomStatusReady {
			//规则判断
			if room.AllWindCount < 8 {
				room.RemainCardCount = 14 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
			} else {
				room.RemainCardCount = 16 + (room.GangCount+room.AllWindCount)%2 + room.GangCount*2
			}
		}
	}
	room.CardIndex[wCurrentUser][room.gameLogic.SwitchToCardIndex(room.SendCardData)]++

	//设置变量
	room.ProvideUser = wCurrentUser
	room.ProvideCard = room.SendCardData

	//设置变量
	room.CurrentUser = wCurrentUser

	//构造数据
	SendCard := &mj_xs_msg.G2C_MJXS_SendCard{}
	SendCard.CurrentUser = wCurrentUser
	SendCard.ActionMask = room.UserAction[wCurrentUser]
	SendCard.CardData = room.SendCardData
	SendCard.Gang = true

	//发送数据
	room.SendMsgAll(SendCard)

	return true
}

func (room *Room) SendOperateNotify() {
	//发送提示
	OperateNotify := mj_xs_msg.G2C_MJXS_OperateNotify{}
	for i, u := range room.Users {
		if u == nil {
			continue
		}
		if room.UserAction[i] != WIK_NULL {
			//构造数据
			OperateNotify.ResumeUser = room.ResumeUser
			OperateNotify.ActionCard = room.ProvideCard
			OperateNotify.ActionMask = room.UserAction[i]
			u.WriteMsg(OperateNotify)
		}
	}

}
