package internal

import (
	"mj/hallServer/db/model"
	"mj/hallServer/user"
	"time"

	"github.com/lovelly/leaf/log"
)

func init() {

}

// 接收到消息，存表
func ReciveGMMsg(sendTimes int, interval int, context string) {

	var raceMsginfo model.RaceMsgInfo
	raceMsginfo.Context = context
	raceMsginfo.SendTimes = sendTimes
	raceMsginfo.IntervalTime = interval

	SendMsgToAll(context)
	raceMsginfo.SendTimes--
	if raceMsginfo.SendTimes > 0 {
		go SendMsgTimer(raceMsginfo)
		log.Debug("GM消息准备插入数据库")
		lastId, rerror := model.RaceMsgInfoOp.Insert(&raceMsginfo)
		if rerror != nil {
			log.Error("GM消息插入数据库失败")
			return
		}
		log.Debug("GM消息插入数据库结果:", lastId)
	}
}

// 服务端启动，从数据库读取GM未发送完成的消息数据
func GetGMMsgFromDB() {
	// 先从数据库里取所有数据
	allMsg, err := model.RaceMsgInfoOp.SelectAll()
	if err != nil {
		log.Error("从race_msg_info表读取所有数据失败,error=%i", err)
		return
	}

	for _, value := range allMsg {
		go SendMsgTimer(*value)
	}
}

// 发数据
func SendMsgTimer(raceMsginfo1 model.RaceMsgInfo) {
	f := func() {
		SendMsgToAll(raceMsginfo1.Context)
		raceMsginfo1.SendTimes--
		if raceMsginfo1.SendTimes > 0 {
			SendMsgTimer(raceMsginfo1)
			model.RaceMsgInfoOp.Update(&raceMsginfo1)
		} else {
			log.Debug("msg为%v的消息发完了", raceMsginfo1.MsgID)
			// 删除数据库
			model.RaceMsgInfoOp.Delete(raceMsginfo1.MsgID)
		}
	}

	time.AfterFunc(time.Duration(raceMsginfo1.IntervalTime)*time.Second, f)
}

//发送消息给所有人
func SendMsgToAll(data interface{}) {
	log.Debug("即将发送消息给所有人：%v", data)
	user.ForEachUser(func(u *user.User) {
		u.WriteMsg(data)
	})
}
