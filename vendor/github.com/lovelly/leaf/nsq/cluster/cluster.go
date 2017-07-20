package cluster

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"runtime/debug"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
)

var (
	clientsMutex sync.Mutex
	clients      = make(map[string]*NsqClient)
)

type NsqClient struct {
	Addr       string
	ServerName string
}

func GetGameServerName(id int) string {
	return fmt.Sprintf("GameSvr_%d", id)
}

func GetHallServerName(id int) string {
	return fmt.Sprintf("HallSvr_%d", id)
}

func AddClient(c *NsqClient) {
	log.Debug("at cluster AddClient %s, %s", c.ServerName, c.Addr)
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	clients[c.ServerName] = c
}

func RemoveClient(serverName string) {
	_, ok := clients[serverName]
	if ok {
		log.Debug("at cluster _removeClient %s", serverName)
		delete(clients, serverName)
		ForEachRequest(func(id int64, request *RequestInfo) {
			if request.serverName != serverName {
				return
			}
			ret := &chanrpc.RetInfo{Ret: nil, Cb: request.cb}
			ret.Err = fmt.Errorf("at call %s server is close", serverName)
			request.chanRet <- ret
			delete(requestMap, id)
		})
		log.Debug("at cluster _removeClient ok %s", serverName)
	}
}

func Broadcast(serverName string, args interface{}) {
	bstr, err := Processor.Marshal(args)
	if err != nil {
		log.Error("CallN Marshal error:%s, stask:%S", err.Error(), string(debug.Stack()))
		return
	}
	msg := &S2S_NsqMsg{CallType: callBroadcast, SrcServerName: SelfName, DstServerName: serverName, Args: bstr[0]}
	Publish(msg)
}

func Go(serverName string, args interface{}) {
	bstr, err := Processor.Marshal(args)
	if err != nil {
		log.Error("CallN Marshal error:%s, stack:%s", err.Error(), string(debug.Stack()))
		return
	}
	msg := &S2S_NsqMsg{ReqType: NsqMsgTypeReq, CallType: callNotForResult, SrcServerName: SelfName, DstServerName: serverName, Args: bstr[0]}
	Publish(msg)
}

//timeOutCall 会丢弃执行结果
func TimeOutCall1(serverName string, t time.Duration, args interface{}) (interface{}, error) {
	bstr, err := Processor.Marshal(args)
	if err != nil {
		log.Error("CallN Marshal error:%s", err.Error())
		return nil, err
	}
	chanSyncRet := make(chan *chanrpc.RetInfo, 1)

	request := &RequestInfo{chanRet: chanSyncRet, serverName: serverName}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, CallType: callForResult, SrcServerName: SelfName, DstServerName: serverName, Args: bstr[0]}
	Publish(msg)
	select {
	case ri := <-chanSyncRet:
		log.Debug("222222222222222222222222")
		return ri.Ret, ri.Err
	case <-time.After(time.Second * t):
		log.Debug("3333333333333333333")
		popRequest(requestID)
		return nil, errors.New(fmt.Sprintf("time out at TimeOutCall1 msg: %v", args))
	}
}

func Call0(serverName string, args interface{}) error {
	bstr, err := Processor.Marshal(args)
	if err != nil {
		log.Error("CallN Marshal error:%s", err.Error())
		return err
	}
	chanSyncRet := make(chan *chanrpc.RetInfo, 1)

	request := &RequestInfo{chanRet: chanSyncRet, serverName: serverName}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, CallType: callForResult, SrcServerName: SelfName, DstServerName: serverName, Args: bstr[0]}
	Publish(msg)

	ri := <-chanSyncRet
	return ri.Err
}

func Call1(serverName string, args interface{}) (interface{}, error) {
	bstr, err := Processor.Marshal(args)
	if err != nil {
		log.Error("CallN Marshal error:%s, %s", err.Error(), string(debug.Stack()))
		return nil, err
	}
	chanSyncRet := make(chan *chanrpc.RetInfo, 1)

	request := &RequestInfo{chanRet: chanSyncRet, serverName: serverName}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, CallType: callForResult, SrcServerName: SelfName, DstServerName: serverName, Args: bstr[0]}
	Publish(msg)

	ri := <-chanSyncRet
	return ri.Ret, ri.Err
}

func CallN(serverName string, args interface{}) ([]interface{}, error) {
	bstr, err := Processor.Marshal(args)
	if err != nil {
		log.Error("CallN Marshal error:%s", err.Error())
		return nil, err
	}
	chanSyncRet := make(chan *chanrpc.RetInfo, 1)

	request := &RequestInfo{chanRet: chanSyncRet, serverName: serverName}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, CallType: callForResult, SrcServerName: SelfName, DstServerName: serverName, Args: bstr[0]}
	Publish(msg)

	ri := <-chanSyncRet
	return chanrpc.Assert(ri.Ret), ri.Err
}

func AsynCall(serverName string, chanAsynRet chan *chanrpc.RetInfo, args interface{}, cb interface{}) {
	bstr, err := Processor.Marshal(args)
	if err != nil {
		log.Error("AsynCall Marshal error:%s", err.Error())
		return
	}

	var callType uint8
	switch cb.(type) {
	case func(error):
		callType = callForResult
	case func(interface{}, error):
		callType = callForResult
	case func([]interface{}, error):
		callType = callForResult
	default:
		panic(fmt.Sprintf("%v asyn call definition of callback function is invalid", args))
	}

	request := &RequestInfo{cb: cb, chanRet: chanAsynRet, serverName: serverName}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, CallType: callType, SrcServerName: SelfName, DstServerName: serverName, Args: bstr[0]}
	Publish(msg)
}
