package gate

import (
	"fmt"
	"io"
	"mj/common/msg"
	"net"
	"reflect"
	"time"

	"regexp"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/module"
	"github.com/lovelly/leaf/network"
)

var (
	filter = []string{
		"C2L_GetRoomList",
		"L2C_GetRoomList",
	}
)

func filterMsg(b []byte) bool {
	for _, msg := range filter {
		if ok, _ := regexp.Match(msg, b); ok {
			return true
		}
	}
	return false
}

type IdUser interface {
	GetUid() int64
}

type UserHandler interface {
	OnInit()
	OnDestroy()
	Run()
	GetChanRPC() *chanrpc.Server
}

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	Processor       network.Processor

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration
	CertFile    string
	KeyFile     string

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool

	// agent
	GoLen              int
	TimerDispatcherLen int
	AsynCallLen        int
	NewChanRPCFunc     func(Agent) UserHandler
	OnAgentInit        func(Agent)
	OnAgentDestroy     func(Agent)
}

func (gate *Gate) Run(closeSig chan bool) {
	time.Sleep(4 * time.Second)
	newAgent := func(conn network.Conn) network.Agent {
		a := &agent{conn: conn, gate: gate}
		if gate.NewChanRPCFunc != nil {
			a.userHandler = gate.NewChanRPCFunc(a)
			a.chanRPC = a.userHandler.GetChanRPC()
		}
		if a.chanRPC != nil {
			a.chanRPC.Go("NewAgent", a)
		}
		return a
	}

	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.CertFile = gate.CertFile
		wsServer.KeyFile = gate.KeyFile
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			return newAgent(conn)
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			return newAgent(conn)
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

type agent struct {
	conn        network.Conn
	userHandler UserHandler
	chanRPC     *chanrpc.Server
	gate        *Gate
	userData    interface{}
	Reason      int
}

func (a *agent) Run() {
	fmt.Println("at aget run .... ")
	defer func() {
		if r := recover(); r != nil {
			log.Recover(r)
		}

		if a.chanRPC != nil {
			err := a.chanRPC.Call0("CloseAgent", a, a.Reason)
			if err != nil {
				log.Error("chanrpc error: %v", err)
			}
		}
	}()

	handleMsgData := func(args []interface{}) error {
		if a.gate.Processor != nil {
			data := args[0].([]byte)
			msg, err := msg.Processor.Unmarshal(data)
			if err != nil {
				return err
			}

			err = a.gate.Processor.Route(msg, a)
			if err != nil {
				return err
			}
		}
		return nil
	}

	if a.chanRPC != nil {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Recover(r)
				}

				if a.gate.OnAgentDestroy != nil {
					a.gate.OnAgentDestroy(a)
				}
			}()

			if a.gate.OnAgentInit != nil {
				a.gate.OnAgentInit(a)
			}

			a.userHandler.Run()
		}()
	}

	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			if err != io.EOF {
				log.Debug("read message: %v", err)
			}
			break
		}
		var userId int64
		var ok bool
		if userId, ok = a.userData.(int64); ok {
		} else if user, ok1 := a.userData.(IdUser); ok1 {
			userId = user.GetUid()
		}

		if !filterMsg(data) {
			log.Debug("IN msg =: %s, userId:%v", string(data), userId)
		}

		if a.chanRPC == nil {
			err = handleMsgData([]interface{}{data})
		} else {
			err = a.chanRPC.Call0("handleMsgData", data)
		}
		if err != nil {
			log.Error("handle message: %v", err)
			break
		}
	}
}

func (a *agent) OnClose() {

}

func (a *agent) SetReason(r int) {
	a.Reason = r
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}

		var userId int64
		var ok bool
		if userId, ok = a.userData.(int64); ok {
		} else if user, ok1 := a.userData.(IdUser); ok1 {
			userId = user.GetUid()
		}
		if !filterMsg(data[0]) {
			log.Debug("OUT msg =: %s, userId:%v", string(data[0]), userId)
		}

		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
		}
	}
}

func (a *agent) LocalAddr() net.Addr {
	return a.conn.LocalAddr()
}

func (a *agent) RemoteAddr() net.Addr {
	return a.conn.RemoteAddr()
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}

func (a *agent) Skeleton() *module.Skeleton {
	return nil
}

func (a *agent) ChanRPC() *chanrpc.Server {
	return a.chanRPC
}
