package consul

import (
	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
)

var (
	caches      map[string]*CacheInfo
	SelfId      string
	HookChanRpc *chanrpc.Server
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
	for id, svr := range caches {
		if _, ok := svrInfo[id]; !ok {
			SvrFaild(svr)
		}
	}
	for id, svr := range svrInfo {
		if _, ok := caches[id]; !ok && svr.Csid != SelfId {
			if HookChanRpc != nil {
				HookChanRpc.Go("ServerStart", svr)
			}
		}
	}

	caches = svrInfo
}

func NotifySvrFaild(args []interface{}) {
	log.Debug("at SvrFaild ==== %v", args)
	faildInfo := args[0].(map[string]string)
	for id, _ := range faildInfo {
		if svr, ok := caches[id]; ok {
			SvrFaild(svr)
		} else {
			log.Debug("no foud old svr %v", caches)
		}
	}
}

func SvrFaild(svr *CacheInfo) {
	log.Debug(" SvrFaild ==== :%s", svr.Csid)
	delete(caches, svr.Csid)
	if HookChanRpc != nil {
		HookChanRpc.Go("ServerFaild", svr)
	}
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
