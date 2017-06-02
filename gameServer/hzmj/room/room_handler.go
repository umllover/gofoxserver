package room

import (
	"github.com/lovelly/leaf/chanrpc"
	"mj/common/msg"
	"mj/gameServer/user"
	. "mj/common/cost"
	"mj/gameServer/db/model/base"
	"time"
	"strconv"
	"fmt"
)

func RegisterHandler(r *Room) {
	r.ChanRPC.RegisterFromType("EnterRoom", OutCard, chanrpc.FuncThis, r)
	r.ChanRPC.RegisterFromType("Sitdown", Sitdown, chanrpc.FuncThis, r)
	r.ChanRPC.RegisterFromType("SetGameOption", SetGameOption, chanrpc.FuncThis, r)
}

func OutCard(args []interface{}) (error) {
	card := args[0].(int)
	room := args[len(args) - 1].(*Room)
	room.SendMsgAll(card )
	return nil
}

func SetGameOption(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_GameOption)
	user := args[1].(*user.User)
	room := args[len(args) - 1].(*Room)
	retcode := 0
	defer func(){
		if retcode != 0 {
			user.WriteMsg(RenderErrorMessage(retcode))
		}
	}()

	if room.CreateUser != user.Id {
		retcode = ErrNotOwner
		return
	}

	if user.ChairId == INVALID_CHAIR {
		retcode = ErrNoSitdowm
		return
	}

	template, ok := base.GameServiceOptionCache.Get(room.Kind, room.ServerId)
	if !ok {
		retcode = ConfigError
		return
	}

	user.WriteMsg(&msg.G2C_GameStatus{
		GameStatus: room.Status,
		AllowLookon:room.AllowLookon[user.ChairId],
	})

	room.AllowLookon[user.ChairId] = recvMsg.AllowLookon
	user.WriteMsg(&msg.G2C_PersonalTableTip{
		TableOwnerUserID: room.CreateUser,			//桌主 I D
		DrawCountLimit: room.CountLimit,				//局数限制
		DrawTimeLimit: room.TimeLimit,				//时间限制
		PlayCount: room.PlayCount,					//已玩局数
		PlayTime: int(room.CreateTime - time.Now().Unix()),		//已玩时间
		CellScore: room.Source,					//游戏底分
		IniScore : room.IniSource,					//初始分数
		ServerID : strconv.Itoa(room.GetRoomId()),					//房间编号
		IsJoinGame :0,					//是否参与游戏 todo  tagPersonalTableParameter
		IsGoldOrGameScore: room.IsGoldOrGameScore,			//金币场还是积分场 0 标识 金币场 1 标识 积分场
	})

	if (template.ServerType &GAME_GENRE_PERSONAL) != 0 { //约战类型。。。
		user.WriteMsg(room.Record)
	}

	if room.Status == RoomStatusReady { // 没开始
		StatusFree := &msg.G2C_StatusFree{}
		StatusFree.CellScore = room.Source				//基础积分
		StatusFree.TimeOutCard  = room.TimeOutCard			//出牌时间
		StatusFree.TimeOperateCard = room.TimeOutCard				//操作时间
		StatusFree.TimeStartGame = room.TimeStartGame				//开始时间
		StatusFree.TurnScore = room.TurnScore					//积分信息
		StatusFree.CollectScore = room.CollectScore				//积分信息
		StatusFree.PlayerCount = room.PlayCount					//玩家人数
		StatusFree.MaCount = room.MaCount						//码数
		StatusFree.CountLimit = room.CountLimit               	//局数限制
		user.WriteMsg(StatusFree)
	}else { //开始了
		StatusPlay := &msg.G2C_StatusPlay{}
		//自定规则
		StatusPlay.TimeOutCard = room.TimeOutCard
		StatusPlay.TimeOperateCard = room.TimeOperateCard
		StatusPlay.TimeStartGame = room.TimeStartGame

		OnUserTrustee(user.ChairId,false)//重入取消托管

		//规则
		StatusPlay.MaCount = room.MaCount
		StatusPlay.PlayerCount = room.PlayerCount
		//游戏变量
		StatusPlay.BankerUser = room.BankerUser
		StatusPlay.CurrentUser = room.OutCardUser
		StatusPlay.CellScore =  room.Source
		StatusPlay.MagicIndex = room.MagicIndex
		StatusPlay.Trustee = room.Trustee

		//状态变量
		StatusPlay.ActionCard = room.ProvideCard
		StatusPlay.LeftCardCount = room.LeftCardCount
		if !room.Response[user.ChairId] {
			StatusPlay.ActionMask = room.UserAction[user.ChairId]
		}else {
			StatusPlay.ActionMask = WIK_NULL
		}

		StatusPlay.Ting = room.Ting
		//当前能胡的牌
		StatusPlay.OutCardCount = room.gameLogic.AnalyseTingCard(room.CardIndex[user.ChairId], room.WeaveItemArray[user.ChairId ],
			room.WeaveItemCount[user.ChairId], StatusPlay.OutCardCount, StatusPlay.OutCardDataEx, StatusPlay.HuCardCount, StatusPlay.HuCardData)

		//历史记录
		StatusPlay.OutCardUser = room.OutCardUser
		StatusPlay.OutCardData = room.OutCardData
		StatusPlay.DiscardCard =  room.DiscardCard
		StatusPlay.DiscardCount = room.DiscardCount

		//组合扑克
		StatusPlay.WeaveItemArray = room.WeaveItemArray
		StatusPlay.WeaveItemCount = room.WeaveItemCount

		//堆立信息
		StatusPlay.HeapHead = room.HeapHead
		StatusPlay.HeapTail = room.HeapTail
		StatusPlay.HeapCardInfo = room.HeapCardInfo

		//扑克数据
		var j int8 = 0
		for ; j < room.UserCnt; j++{
			StatusPlay.CardCount[j] = room.gameLogic.GetCardCount(room.CardIndex[j])
		}
		room.gameLogic.SwitchToCardData(room.CardIndex[user.ChairId], StatusPlay.CardData)
	 	if room.CurrentUser == user.ChairId {
			 StatusPlay.SendCardData = room.SendCardData
		}else {
		 	StatusPlay.SendCardData =0x00
		 }

		//历史积分
		for j = 0; j < room.UserCnt; j++{
			//设置变量
			StatusPlay.TurnScore[j] = room.HistoryScores[j].TurnScore
			StatusPlay.CollectScore[j] =room.HistoryScores[j].CollectScore
		}

		user.WriteMsg(StatusPlay)
	}
}

func Sitdown(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_UserSitdown)
	user := args[1].(*user.User)
	room := args[len(args) - 1].(*Room)
	retcode := 0
	defer func(){
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

	if room.GetRoomStatus() == RoomStatusStarting && template.DynamicJoin == 1{
		retcode = GameIsStart
		return
	}

	fmt.Println(user.Id)
	fmt.Println(room.RoomInfo)
	_, chairId := room.GetUserByUid(user.Id)
	if chairId > 0 {
		room.LeaveRoom(chairId)
	}

	room.EnterRoom(recvMsg.ChairID, user)
	user.Status = US_SIT
	room.SendMsgAll(&msg.G2C_UserStatus{
		UserID:user.Id,
		UserStatus:&msg.UserStu{
			TableID: room.id,
			ChairID: user.ChairId,
			UserStatus:user.Status,
		},
	})
}



/////////////////// help
//托管
func OnUserTrustee(chairId int, trusteeship bool) {

}