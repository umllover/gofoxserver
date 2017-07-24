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
	ActivityCache.LoadAll()
	FreeLimitCache.LoadAll()
	GameServiceOptionCache.LoadAll()
	GameTestpaiCache.LoadAll()
	GlobalVarCache.LoadAll()
	GoodsCache.LoadAll()
	PersonalTableFeeCache.LoadAll()
	RefreshInTimeCache.LoadAll()
	ServerListCache.LoadAll()
	db.BaseDataCaches["Activity"] = ActivityCache
	db.BaseDataCaches["FreeLimit"] = FreeLimitCache
	db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
	db.BaseDataCaches["GameTestpai"] = GameTestpaiCache
	db.BaseDataCaches["GlobalVar"] = GlobalVarCache
	db.BaseDataCaches["Goods"] = GoodsCache
	db.BaseDataCaches["PersonalTableFee"] = PersonalTableFeeCache
	db.BaseDataCaches["RefreshInTime"] = RefreshInTimeCache
	db.BaseDataCaches["ServerList"] = ServerListCache
	log.Debug("loadBaseData %v  %v %v", 9, time.Now().UnixNano()-start, "ns")
}
