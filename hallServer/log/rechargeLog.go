package log

import (
	"mj/hallServer/db/model/stats"
	"time"

	"github.com/lovelly/leaf/log"
)

type RechargeLog struct{}

func (rechargelog *RechargeLog) AddRechargeLogInfo(OnlineID, PayAmount int, userID int64, PayType string, GoodsID int) {
	now := time.Now()
	_, err := stats.RechargeLogOp.Insert(&stats.RechargeLog{
		OnLineID:     OnlineID,
		PayAmount:    PayAmount,
		UserID:       userID,
		PayType:      PayType,
		GoodsID:      GoodsID,
		RechangeTime: &now,
	})
	if err != nil {
		log.Error("insert rechargeInfo data into table recharge_log")
	}
}
