package base

//This file is generate by scripts,don't edit it
//

import (
	"mj/hallServer/db"
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
    PersonalTableFeeBak-7-28Cache.LoadAll()
    RefreshInTimeCache.LoadAll()
    ServerListCache.LoadAll()
    db.BaseDataCaches["Activity"] = ActivityCache
    db.BaseDataCaches["FreeLimit"] = FreeLimitCache
    db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
    db.BaseDataCaches["GameTestpai"] = GameTestpaiCache
    db.BaseDataCaches["GlobalVar"] = GlobalVarCache
    db.BaseDataCaches["Goods"] = GoodsCache
    db.BaseDataCaches["PersonalTableFee"] = PersonalTableFeeCache
    db.BaseDataCaches["PersonalTableFeeBak-7-28"] = PersonalTableFeeBak-7-28Cache
    db.BaseDataCaches["RefreshInTime"] = RefreshInTimeCache
    db.BaseDataCaches["ServerList"] = ServerListCache
    log.Debug("loadBaseData %v  %v %v", 10,  time.Now().UnixNano()-start, "ns")
}