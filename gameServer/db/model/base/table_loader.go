package base

//This file is generate by scripts,don't edit it
//

import (
	"mj/gameServer/db"
	"time"

	"github.com/lovelly/leaf/log"
)

func LoadBaseData() {
	var start = time.Now().UnixNano()
	GameServiceOptionCache.LoadAll()
	PersonalTableFeeCache.LoadAll()
	db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
	db.BaseDataCaches["PersonalTableFee"] = PersonalTableFeeCache
	log.Debug("loadBaseData %v, %v %v", 2, time.Now().UnixNano()-start, "ns")
}
