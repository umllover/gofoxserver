package base

//This file is generate by scripts,don't edit it
//

import (
	"mj/gameServer/db"
    "github.com/lovelly/leaf/log"
)

func LoadBaseData() {
	var start = time.Now().UnixNano()
    ActivityCache.LoadAll()
    GameServiceOptionCache.LoadAll()
    GameTestpaiCache.LoadAll()
    GlobalVarCache.LoadAll()
    db.BaseDataCaches["Activity"] = ActivityCache
    db.BaseDataCaches["GameServiceOption"] = GameServiceOptionCache
    db.BaseDataCaches["GameTestpai"] = GameTestpaiCache
    db.BaseDataCaches["GlobalVar"] = GlobalVarCache
    log.Debug("loadBaseData %v, %v %v", 4,  time.Now().UnixNano()-start, "ns")
}