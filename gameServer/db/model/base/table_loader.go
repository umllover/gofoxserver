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
	GameTestpaiCache.LoadAll()
	GlobalVarCache.LoadAll()
	PersonalTableFeeCache.LoadAll()
	RefreshInTimeCache.LoadAll()
	ServerListCache.LoadAll()
	db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
	db.BaseDataCaches["GameTestpai"] = GameTestpaiCache
	db.BaseDataCaches["GlobalVar"] = GlobalVarCache
	db.BaseDataCaches["PersonalTableFee"] = PersonalTableFeeCache
	db.BaseDataCaches["RefreshInTime"] = RefreshInTimeCache
	db.BaseDataCaches["ServerList"] = ServerListCache
	log.Debug("loadBaseData %v, %v %v", 6, time.Now().UnixNano()-start, "ns")
}
