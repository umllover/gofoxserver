
package main


import (
	"flag"
	//"fmt"
	//"mj/client/proto"
	//"github.com/lovelly/leaf/log"
	"mj/client/proto"
	"github.com/lovelly/leaf/log"
)
/*
	"WSAddr": "192.168.1.81:20001",
    "HttpAddr":"192.168.1.81:40001",
    "ListenAddr": "192.168.1.81:30001",
    */

var listenHost = flag.String("listenHost", "192.168.1.81", "listen host")
var listenPort = flag.Int("listenPort", 30001, "listen port")
var wsHost = flag.String("wsHost", "192.168.1.81", "ws host")
var wsPort = flag.Int("wsPort", 20001, "ws port")


func main() {

	/*flag.Parse()
	listenAddr := fmt.Sprintf("%s:%d", *listenHost, *listenPort)
	wsAddr := fmt.Sprintf("%s:%d", *wsHost, *wsPort)
	log.Debug("listenAddr : %s, wsAddr : %s")*/
	tcpAddr := "192.168.1.81:20001"
	//wsAddr := "ws://192.168.1.81:20001"

	mockClient := &proto.MockClient{}
	ret := mockClient.ConnectWS(tcpAddr)
	if ret == false {
		log.Debug("connect ws failed")
		return
	}

	mockClient.UserLogin()

	requests := make([]*proto.Request, 0)
	//requests = append(requests, &clientSimulator.Request{FuncName: "Login", Params: []interface{}{mockClient.UserId, mockClient.SessionKey, 32}})
	requests = append(requests, proto.MakeRequest()...)

	mockClient.Run(requests)
}
