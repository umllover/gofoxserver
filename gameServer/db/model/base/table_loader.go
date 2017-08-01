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
	AchievementLimitCache.LoadAll()
	ActivityCache.LoadAll()
	AgentLimitCache.LoadAll()
	FreeLimitCache.LoadAll()
	GameServiceOptionCache.LoadAll()
	GameTestpaiCache.LoadAll()
	GlobalVarCache.LoadAll()
	GoodsCache.LoadAll()
	PersonalTableFeeCache.LoadAll()
	PersonalTableFeeBak728Cache.LoadAll()
	RechargeLimitCache.LoadAll()
	RefreshInTimeCache.LoadAll()
	ServerListCache.LoadAll()
	db.BaseDataCaches["AchievementLimit"] = AchievementLimitCache
	db.BaseDataCaches["Activity"] = ActivityCache
	db.BaseDataCaches["AgentLimit"] = AgentLimitCache
	db.BaseDataCaches["FreeLimit"] = FreeLimitCache
	db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
	db.BaseDataCaches["GameTestpai"] = GameTestpaiCache
	db.BaseDataCaches["GlobalVar"] = GlobalVarCache
	db.BaseDataCaches["Goods"] = GoodsCache
	db.BaseDataCaches["PersonalTableFee"] = PersonalTableFeeCache
	db.BaseDataCaches["PersonalTableFeeBak728"] = PersonalTableFeeBak728Cache
	db.BaseDataCaches["RechargeLimit"] = RechargeLimitCache
	db.BaseDataCaches["RefreshInTime"] = RefreshInTimeCache
	db.BaseDataCaches["ServerList"] = ServerListCache
	log.Debug("loadBaseData %v, %v %v", 13, time.Now().UnixNano()-start, "ns")
}
