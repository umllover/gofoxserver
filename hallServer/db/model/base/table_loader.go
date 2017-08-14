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
	IncAgentNumCache.LoadAll()
	PersonalTableFeeCache.LoadAll()
	PersonalTableFeeBak728Cache.LoadAll()
	RefreshInTimeCache.LoadAll()
	UpgradeCache.LoadAll()
	UpgradeAdvisorCache.LoadAll()
	db.BaseDataCaches["Activity"] = ActivityCache
	db.BaseDataCaches["FreeLimit"] = FreeLimitCache
	db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
	db.BaseDataCaches["GameTestpai"] = GameTestpaiCache
	db.BaseDataCaches["GlobalVar"] = GlobalVarCache
	db.BaseDataCaches["Goods"] = GoodsCache
	db.BaseDataCaches["IncAgentNum"] = IncAgentNumCache
	db.BaseDataCaches["PersonalTableFee"] = PersonalTableFeeCache
	db.BaseDataCaches["PersonalTableFeeBak728"] = PersonalTableFeeBak728Cache
	db.BaseDataCaches["RefreshInTime"] = RefreshInTimeCache
	db.BaseDataCaches["Upgrade"] = UpgradeCache
	db.BaseDataCaches["UpgradeAdvisor"] = UpgradeAdvisorCache
	log.Debug("loadBaseData %v  %v %v", 12, time.Now().UnixNano()-start, "ns")
}
