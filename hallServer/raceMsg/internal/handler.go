package internal

import (
	"math/rand"
	"reflect"
	"time"

	"mj/common/msg"
	"mj/hallServer/user"
	"mj/hallServer/userHandle"

	"mj/hallServer/db/model"

	"github.com/lovelly/leaf/log"
)

////注册rpc 消息
func handleRpc(id interface{}, f interface{}) {
	ChanRPC.Register(id, f)
}

//注册 客户端消息调用
func handlerC2S(m interface{}, h interface{}) {
	msg.Processor.SetRouter(m, ChanRPC)
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {

}

// 接收到消息，存表
func ReciveMsg(beginTime int, endTime int, interval int, context string) {

	var raceMsginfo model.RaceMsgInfo
	raceMsginfo.Context = context
	raceMsginfo.EndTime = endTime
	raceMsginfo.StartTime = beginTime
	raceMsginfo.IntervalTime = interval

	log.Debug("准备插入数据库")
	lastId, rerror := model.RaceMsgInfoOp.Insert(&raceMsginfo)
	if rerror != nil {
		log.Debug("插入失败")
	}
	log.Debug("插入数据库结果:", lastId)
}

// 启动跑马灯
func StartHorseRaceLamp() {
	go startTimer(5 * 60)
}

// 启动定时器
func startTimer(nTime int) {
	log.Debug("定时器启动了")
	timerMsg := time.NewTicker(time.Duration(nTime) * time.Second)
	for {
		select {
		case <-timerMsg.C:
			log.Debug("定时器时间到了，开始发消息")
			msgInfo, err := model.RaceMsgInfoOp.SelectAll()
			if err != nil {
				log.Error("race_msg_info查找所有数据失败,error=%i", err)
				break
			}

			msgCount := len(msgInfo)
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			index := r.Intn(msgCount)

			//timerMsg.Stop()
			SendMsgToAll(msgInfo[index].Context)
		default:
		}
	}
}

//发送消息给所有人
func SendMsgToAll(data interface{}) {

	log.Debug("即将发送消息给所有人：%v", data)
	userHandle.UserMgr.ForEachUser(func(u *user.User) {
		u.WriteMsg(data)
	})
}

//发送消息给某人
func sendMsgToUser(args []interface{}) {

}
