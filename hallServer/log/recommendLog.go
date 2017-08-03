package log

import (
	"mj/hallServer/db/model/stats"
	"time"

	"github.com/lovelly/leaf/log"
)

type RecommendLog struct{}

func (recommendLog *RecommendLog) AddRecommendLog(userId, ElectUid int64) {
	now := time.Now()
	_, err := stats.RecommendLogOp.Insert(&stats.RecommendLog{
		SubElectUid: userId,
		ElectUid:    ElectUid,
		ElectTime:   &now,
	})
	if err != nil {
		log.Error("get data from RecommendLog Table!!!")
	}
}
