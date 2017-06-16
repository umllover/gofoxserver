package http_service

import (
	"mj/common/http_svr"
	"mj/hallServer/conf"
	"net/http"
	"net/http/pprof"

	"github.com/lovelly/leaf/log"
)

func StartHttpServer() {
	rpc := http_svr.NewRpcHelper()
	rpc.RegisterMethod(DefaultHttpHandler)
	hsvr := http_svr.NewHTTPServer(conf.Server.HttpAddr, rpc)
	if hsvr != nil {
		hsvr.Start()
	} else {
		log.Debug(" http svr not start addr is nil .....")
	}
}

func StartPrivateServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	svr := http_svr.NewHTTPServer(conf.Server.WatchAddr, nil, 0, mux)
	if svr != nil {
		svr.Start()
	}
}
