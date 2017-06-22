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
	GameTestpaiCache.LoadAll()
	GlobalVarCache.LoadAll()
	PersonalTableFeeCache.LoadAll()
	RefreshInTimeCache.LoadAll()
	ServerListCache.LoadAll()
	TokenRecordCache.LoadAll()
	db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
	db.BaseDataCaches["GameTestpai"] = GameTestpaiCache
	db.BaseDataCaches["GlobalVar"] = GlobalVarCache
	db.BaseDataCaches["PersonalTableFee"] = PersonalTableFeeCache
	db.BaseDataCaches["RefreshInTime"] = RefreshInTimeCache
	db.BaseDataCaches["ServerList"] = ServerListCache
	db.BaseDataCaches["TokenRecord"] = TokenRecordCache
	log.Debug("loadBaseData %v  %v %v", 7, time.Now().UnixNano()-start, "ns")
}
