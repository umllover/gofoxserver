package http_svr

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/lovelly/leaf/log"
)

type HTTPServer struct {
	Addr        string
	MaxMsgLen   uint32
	HTTPTimeout time.Duration
	ln          net.Listener
	handler     http.Handler
	wg          sync.WaitGroup
	rpch        *RpcHelper
}

func NewHTTPServer(addr string, pa ...interface{}) *HTTPServer {
	if addr == "" {
		return nil
	}
	server := &HTTPServer{
		Addr:        "127.0.0.1:8080",
		MaxMsgLen:   65535,
		HTTPTimeout: 5,
	}

	server.Addr = addr
	if len(pa) > 0 {
		rpc, ok := pa[0].(*RpcHelper)
		if ok {
			server.rpch = rpc
		}
	}

	if len(pa) > 1 {
		t, ok := pa[1].(int)
		if ok && t != 0 {
			server.HTTPTimeout = time.Duration(t) * time.Second
		}
	}

	if len(pa) > 2 {
		hd, ok := pa[2].(http.Handler)
		if ok {
			server.handler = hd
		}
	}

	return server
}

func (server *HTTPServer) Serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("read data error:%v", err.Error())
		w.WriteHeader(500)
		w.Write([]byte("server error"))
		return
	}
	if len(requestData) < 1 {
		log.Error(" at HTTPServer not foud body")
		return
	}

	log.Debug("http in =:", string(requestData))
	method, params, game_err := server.rpch.Parse(&HandelData{MsgType: FROM_CLIENT_MSG, Data: requestData})
	if game_err != nil {
		log.Error("parse msg error arr = %v, error:%v", r.RemoteAddr, game_err.Error())
		return
	}

	server.wg.Add(1)
	defer server.wg.Done()

	result, game_err := server.rpch.Call(method, params)
	if game_err != nil {
		log.Error("call error addr:%v, error:%v", r.RemoteAddr, game_err.Error())
		result.Error = game_err.Error()
	}

	bytes, err := json.Marshal(result)
	log.Error("http  out =: ", string(bytes))
	if err != nil {
		log.Error("Marshal result error at httpsver addr:%v error:%v", r.RemoteAddr, "convert call result to json failed.")
	} else {
		w.Write(bytes)
	}
}

func (server *HTTPServer) Start() {
	log.Debug(" start HTTPServer at ", server.Addr)
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal("lister error %v", err)
	}

	server.ln = ln
	if server.handler == nil {
		mux := http.NewServeMux()
		mux.HandleFunc("/rpc", server.Serve)
		server.handler = mux
	}

	httpServer := &http.Server{
		Addr:           server.Addr,
		Handler:        server.handler,
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln)
}

func (server *HTTPServer) Close() {
	log.Debug("at HTTPServer close")
	server.ln.Close()
	server.wg.Wait()
}
