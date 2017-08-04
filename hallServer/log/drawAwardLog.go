package log

import (
	"mj/hallServer/db/model/stats"
	"time"

	"github.com/lovelly/leaf/log"
)

type DrawAwardLog struct{}

func (drarawardLog DrawAwardLog) AddDrawAdardLog(id, drawId int, description string, drawcount int64, drawType, amount, itemtype int) {
	now := time.Now()
	_, err := stats.DrawAwardLogOp.Insert(&stats.DrawAwardLog{
		Id:          id,
		DrawId:      drawId,
		Description: description,
		DrawCount:   drawcount,
		DrawType:    drawType,
		Amount:      amount,
		ItemType:    itemtype,
		DrawTime:    &now,
	})
	if err != nil {
		log.Error("insert data into table draw_award_log Error:%v", err.Error())
	}
}
