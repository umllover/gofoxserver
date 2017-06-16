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
	GlobalVarCache.LoadAll()
	PersonalTableFeeCache.LoadAll()
	ServerListCache.LoadAll()
	db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
	db.BaseDataCaches["GlobalVar"] = GlobalVarCache
	db.BaseDataCaches["PersonalTableFee"] = PersonalTableFeeCache
	db.BaseDataCaches["ServerList"] = ServerListCache
	log.Debug("loadBaseData %v  %v %v", 4, time.Now().UnixNano()-start, "ns")
}
