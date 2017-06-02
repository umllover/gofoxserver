package internal

import (
	"mj/common/msg"
	"github.com/lovelly/leaf/cluster"
	"reflect"
	"github.com/lovelly/leaf/gate"
	"github.com/name5566/leaf/log"
	"mj/gameServer/user"
)

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

func init(){
	handlerC2S(&msg.C2G_REQUserInfo{}, GetUserInfo)
	handlerC2S(&msg.C2G_REQUserChairInfo{}, GetUserChairInfo)
}


func GetUserInfo(args []interface{}) {


}

func GetUserChairInfo(args []interface{}) {
	recvMsg := args[0].(*msg.C2G_REQUserChairInfo)
	_ = recvMsg
	agent := args[1].(gate.Agent)
	user, ok := agent.UserData().(*user.User)
	if !ok {
		log.Error("at GerUserInfo user not logon")
		return
	}
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


