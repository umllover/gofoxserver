package base

//This file is generate by scripts,don't edit it
//

import (
	"mj/hallServer/db"
	"time"

	"github.com/lovelly/leaf/log"
)

func LoadBaseData() {
	var start = time.Now().UnixNano()
	GameServiceOptionCache.LoadAll()
	GameServiceSttribCache.LoadAll()
	db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
	db.BaseDataCaches["GameServiceSttrib"] = GameServiceSttribCache
	log.Debug("loadBaseData %v  %v %v", 2, time.Now().UnixNano()-start, "ns")
}
