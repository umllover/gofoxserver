package internal

import (
	"regexp"

	"github.com/lovelly/leaf/cluster"
	"github.com/lovelly/leaf/log"
)

var (
	caches        map[string]*CacheInfo
	SelfId        string
	InitiativeSvr []string //需要注定去连接的类型
)

func init() {
	//init
	caches = make(map[string]*CacheInfo)
	//regist
	handleRpc("AddServerInfo", AddServerInfo)
	handleRpc("NotifySvrFaild", NotifySvrFaild)
	handleRpc("KvUpdate", KvUpdate)
	handleRpc("GetAllSvrInfo", GetAllSvrInfo)
}

func handleRpc(id interface{}, f interface{}) {
	ChanRPC.Register(id, f)
}

func AddServerInfo(args []interface{}) {
	svrInfo := args[0].(map[string]*CacheInfo)
	for id, _ := range caches {
		if _, ok := svrInfo[id]; !ok {
			SvrFaild(id)
		}
	}
	for id, svr := range svrInfo {
		if _, ok := caches[id]; !ok && svr.Csid != SelfId && len(InitiativeSvr) > 0 {
			for _, v := range InitiativeSvr {
				rok, err := regexp.MatchString(v, svr.Csid)
				if err != nil {
					log.Error("at AddServerInfo Error:%s", err.Error())
					continue
				}
				if rok {
					cluster.AddClient(svr.Csid, svr.Host)
				}
			}
		}
	}

	caches = svrInfo
}
func NotifySvrFaild(args []interface{}) {
	log.Debug("at SvrFaild ==== %v", args)
	faildInfo := args[0].(map[string]string)
	for id, _ := range faildInfo {
		if _, ok := caches[id]; ok {
			SvrFaild(id)
		} else {
			log.Debug("no foud old svr %v", caches)
		}
	}
}

func SvrFaild(id string) {
	log.Debug(" SvrFaild ==== :%s", id)
	delete(caches, id)
	cluster.RemoveClient(id)

}

func KvUpdate(args []interface{}) {
	KvInfo := args[0].(map[string]int)
	for k, v := range KvInfo {
		if svr, ok := caches[k]; ok {
			if svr.weight != v {
				svr.weight = v
			}
		} else {
			log.Error(" at KvUpdate no foud svr %s", k)
		}
	}
}

func GetAllSvrInfo(args []interface{}) (interface{}, error) {
	ret := make(map[string]*CacheInfo)
	for k, v := range caches {
		ret[k] = v
	}
	return ret, nil
}
