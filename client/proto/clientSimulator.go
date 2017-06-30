package proto


import (
	"encoding/json"
	//"errors"
	//"fmt"
	"net"
	//"net/url"
	"os"
	//"os/signal"
	//"runtime/debug"
	//"strconv"
	"sync"
	//"syscall"


	"time"
	"github.com/lovelly/leaf/log"
	"github.com/gorilla/websocket"

	/*"github.com/uber-go/zap" */

	//"mj/common/utils"
	"net/url"
	/*"io/ioutil"
	"mj/common/msg"
	"encoding/json"
	"runtime/debug"
	"os/signal"
	"syscall"
	"fmt"*/
	"mj/common/msg"
	"os/signal"
	"syscall"
)

var Wg sync.WaitGroup
var sigsub = make(chan os.Signal, 1)

var (
	readBufferSize  int = 1024
	writeBufferSize int = 1024
)

const (
	WRITE_TIMEOUT = 10 * time.Second
	PONG_TIMEOUT  = 30 * time.Second
	PING_PERIOD   = (PONG_TIMEOUT * 9) / 10
	MAX_MSG_SIZE  = 512 * 1000
)

var DATATYPE = websocket.BinaryMessage
//var DATATYPE = t

type Request struct {
	FuncName string        `json:"func_name"`
	Params   []interface{} `json:"params"`
}
type Data map[string]interface{}
type Response struct {
	FuncName string               `json:"func_name"`
	Data     Data                 `json:"data,omitempty"`
	//Error    game_error.GameError `json:"error,omitempty"`
}

type MockClient struct {
	conn         	*websocket.Conn
	IsEnterWorld 	bool
	UserName     	string
	Password     	string

	UserId       int
/*	SessionKey   string */

}



func (c *MockClient) SendReq(req *Request) {

	c.conn.SetReadLimit(MAX_MSG_SIZE)
	c.conn.SetReadDeadline(time.Now().Add(time.Second * 1))
	//c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(PONG_TIMEOUT)); return nil })
	dataBytes, err := json.Marshal(req)
	if err != nil {
		log.Debug("json marsha err %v", err)
		return
	}

	err = c.conn.WriteMessage(DATATYPE, dataBytes)
	if err != nil {
		log.Debug("write message err %v", err)
	}
}

func (c *MockClient) GetResp() []byte {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		log.Debug("ReadMessage error %v", err)
		//break
	}

	return message
}

func (mockClient *MockClient) UserLogin() {

	c2lLogin := &msg.C2L_Login{
		ModuleID:     1,
		PlazaVersion: 1,
		DeviceType:   1,
		LogonPass:    "e10adc3949ba59abbe56e057f20f883e",
		Accounts:     "jw1",
		MachineID:    "SDC641523AEE420A90A74DB40779D922",
		MobilePhone:  "1222222222",
	}

	msg := map[string]interface{}{
		"C2L_Login":c2lLogin,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		return
	}

	err = mockClient.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		log.Debug("user login write msg err : %s", err.Error())
		return
	}

	Wg.Add(1)
	go func(){
		for {
			_, b, err = mockClient.conn.ReadMessage()
			if err != nil {
				log.Debug(err.Error())
				log.Debug("user login get result failed : %s", err.Error())
				Wg.Done()
				return
			}

			if b != nil {
				log.Debug(string(b))
				Wg.Done()
				return
			}
		}
	}()
	Wg.Wait()
}

func (c *MockClient) ConnectWS(tcpAddr string) (ret bool)  {
	connectRet := false
	defer func() {
		if err := recover(); err != nil {
			log.Debug("游戏服务器报错 : %v", err)
		}
	}()

	l, err := net.Dial("tcp", tcpAddr)
	if err != nil {
		log.Debug("dial failed")
		log.Debug("err %v", err)
		return connectRet
	}
	wsUrl := "ws://" + tcpAddr
	u, err := url.Parse(wsUrl)
	if err != nil {
		log.Debug("url parse failed")
		log.Debug("err %v", err)
		return  connectRet
	}
	conn, response, err := websocket.NewClient(l, u, nil, readBufferSize, writeBufferSize)
	if err != nil {
		log.Debug("newclient failed")
		log.Debug("err %v,response %v", err, response)
		return connectRet
	}
	c.conn = conn
	connectRet = true
	log.Debug("connect ws sucess")
	return connectRet
}


func (c *MockClient) Run(reqs []*Request) {
	signal.Notify(sigsub, syscall.SIGINT)
	Wg.Add(1)

	go func() {

		ticker := time.NewTicker(time.Second * 2)

	DONE:
		for _, req := range reqs {
			select {
			case s, ok := <-sigsub:
				if ok {
					log.Debug("get signal %s from sigsub, will stop...\n", s)
				} else {
					log.Debug("get signal from sigsub error")
				}
				break DONE

			case <-ticker.C:
				c.SendReq(req)
				message := c.GetResp()
				//zaplogger.Info("recv:", zap.String("", string(message)))
				log.Debug("run recv %s", string(message))
			}
		}

		Wg.Done()

	}()

	Wg.Wait()
}

func MakeRequest() []*Request {
	requests := make([]*Request, 0)

	/*requests = append(requests, &Request{FuncName: "AddHeroTalent", Params: []interface{}{300003, 31}})
	requests = append(requests, &Request{FuncName: "AddHeroTalent", Params: []interface{}{300003, 31}})
	requests = append(requests, &Request{FuncName: "AddHeroTalent", Params: []interface{}{300003, 31}})
	requests = append(requests, &Request{FuncName: "AddHeroTalent", Params: []interface{}{300003, 31}})*/
	requests = append(requests, &Request{FuncName: "GetUserIndividual", Params: []interface{}{11}})

	return requests
}
