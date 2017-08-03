package log

import (
	"mj/hallServer/db/model/stats"
	"time"
)

type DrawAwardLog struct{}

func (drarawardLog DrawAwardLog) AddDrawAdardLog(id, drawId int, description string, drawcount int64, drawType, amount, itemtype int) {
	now := time.Now()
	stats.DrawAwardLogOp.Insert(&stats.DrawAwardLog{
		Id:          id,
		DrawId:      drawId,
		Description: description,
		DrawCount:   drawcount,
		DrawType:    drawType,
		Amount:      amount,
		ItemType:    itemtype,
		DrawTime:    &now,
	})
}
