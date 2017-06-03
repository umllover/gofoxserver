package utils

import (
	"time"
	"net/http"
	"net/http/pprof"
	"fmt"
	"github.com/lovelly/leaf/log"
)
func CreatePrivateServer(prot int){
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	addr :=  fmt.Sprintf("127.0.0.1:%d", prot)
	s := &http.Server{
		Addr:addr,
		Handler: mux,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Release("Private Server listen at %v", addr)
	go s.ListenAndServe()
}