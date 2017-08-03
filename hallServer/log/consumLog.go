package log

import (
	"mj/hallServer/db/model/stats"
	"time"
	"github.com/lovelly/leaf/log"
)

type ConsumLog struct{}

func (consumLog *ConsumLog) AddConsumLogInfo(userId int64, subType int, money int) {
	now := time.Now()
	_,err:=stats.ConsumLogOp.Insert(&stats.ConsumLog{
		UserId:     userId,
		ConsumType: subType,
		ConsumNum:  money,
		ConsumTime: &now,
	})
	if err!=nil{
		log.Error("insert data into table Error:%v",err.Error())
	}
}
