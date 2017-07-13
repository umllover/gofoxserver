package cluster

import (
	"auth_server/module/log"
	"errors"
	"fmt"
	"mj/gameServer/conf"

	"sync"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/network/gob"
)

const (
	callBroadcast = iota
	callNotForResult
	callForResult
)

const (
	NsqMsgTypeReq = iota
	NsqMsgTypeRsp
)

var (
	routeMap        = map[interface{}]*chanrpc.Client{}
	Processor       = gob.NewProcessor()
	RequestInfoLock sync.Mutex
	requestID       int64
	requestMap      = make(map[int64]*RequestInfo)
)

type RequestInfo struct {
	cb      interface{}
	chanRet chan *chanrpc.RetInfo
}

func SetRoute(id interface{}, server *chanrpc.Server) {
	_, ok := routeMap[id]
	if ok {
		panic(fmt.Sprintf("function id %v: already set route", id))
	}

	routeMap[id] = server.Open(0)
}

type S2S_NsqMsg struct {
	RequestID  int64
	ReqType    int
	MsgID      interface{}
	CallType   uint8
	ServerName string
	Args       []interface{}
	Err        string
}

func handleRequestMsg(recvMsg *S2S_NsqMsg) {
	sendMsg := &S2S_NsqMsg{ReqType: NsqMsgTypeRsp, RequestID: recvMsg.RequestID}
	if isClose() && recvMsg.CallType == callForResult {
		sendMsg.Err = fmt.Sprintf("%v server is closing", conf.ServerName)
		Publish(recvMsg.ServerName, sendMsg)
		return
	}

	msgID := recvMsg.MsgID
	client, ok := routeMap[msgID]
	if !ok {
		err := fmt.Sprintf("%v msg is not set route", msgID)
		log.Error(err)

		if recvMsg.CallType == callForResult {
			sendMsg.Err = err
			Publish(recvMsg.ServerName, sendMsg)
		}
		return
	}

	args := recvMsg.Args
	if recvMsg.CallType == callForResult {
		sendMsgFunc := func(ret *chanrpc.RetInfo) {
			sendMsg.Args = []interface{}{ret.Ret}
			if ret.Err != nil {
				sendMsg.Err = ret.Err.Error()
			}
			Publish(recvMsg.ServerName, sendMsg)
		}

		args = append(args, sendMsgFunc)
		client.RpcCall(msgID, args...)
	} else {
		args = append(args, nil)
		client.RpcCall(msgID, args...)
	}
}

func handleResponseMsg(msg *S2S_NsqMsg) {
	request := popRequest(msg.RequestID)
	if request == nil {
		log.Error("%v: request id %v is not exist", msg.ServerName, msg.RequestID)
		return
	}

	ret := &chanrpc.RetInfo{Ret: msg.Args[0], Cb: request.cb}
	if msg.Err != "" {
		ret.Err = errors.New(msg.Err)
	}
	request.chanRet <- ret
}

func registerRequest(request *RequestInfo) int64 {
	RequestInfoLock.Lock()
	defer RequestInfoLock.Unlock()
	reqID := requestID
	requestMap[reqID] = request
	requestID += 1
	return reqID
}

func popRequest(requestID int64) *RequestInfo {
	RequestInfoLock.Lock()
	defer RequestInfoLock.Unlock()

	request, ok := requestMap[requestID]
	if ok {
		delete(requestMap, requestID)
		return request
	} else {
		return nil
	}
}
