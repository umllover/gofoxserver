package cluster

import (
	"errors"
	"time"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"

	"fmt"
	"sync"
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
	}
}

func Broadcast(serverName string, id interface{}, args ...interface{}) {
	msg := &S2S_NsqMsg{MsgID: id, CallType: callBroadcast, Args: args}
	Publish(serverName, msg)
}

func Go(serverName string, id interface{}, args ...interface{}) {
	msg := &S2S_NsqMsg{MsgID: id, ReqType: NsqMsgTypeReq, CallType: callNotForResult, Args: args}
	Publish(serverName, msg)
}

//timeOutCall 会丢弃执行结果
func TimeOutCall1(serverName string, id interface{}, t time.Duration, args ...interface{}) (interface{}, error) {
	chanSyncRet := make(chan *chanrpc.RetInfo, 1)

	request := &RequestInfo{chanRet: chanSyncRet}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, MsgID: id, CallType: callForResult, Args: args}
	Publish(serverName, msg)
	select {
	case ri := <-chanSyncRet:
		return ri.Ret, ri.Err
	case <-time.After(time.Second * t):
		popRequest(requestID)
		return nil, errors.New(fmt.Sprintf("time out at TimeOutCall1 function: %v", id))
	}
}

func Call0(serverName string, id interface{}, args ...interface{}) error {
	chanSyncRet := make(chan *chanrpc.RetInfo, 1)

	request := &RequestInfo{chanRet: chanSyncRet}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, MsgID: id, CallType: callForResult, Args: args}
	Publish(serverName, msg)

	ri := <-chanSyncRet
	return ri.Err
}

func Call1(serverName string, id interface{}, args ...interface{}) (interface{}, error) {
	chanSyncRet := make(chan *chanrpc.RetInfo, 1)

	request := &RequestInfo{chanRet: chanSyncRet}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, MsgID: id, CallType: callForResult, Args: args}
	Publish(serverName, msg)

	ri := <-chanSyncRet
	return ri.Ret, ri.Err
}

func CallN(serverName string, id interface{}, args ...interface{}) ([]interface{}, error) {
	chanSyncRet := make(chan *chanrpc.RetInfo, 1)

	request := &RequestInfo{chanRet: chanSyncRet}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, MsgID: id, CallType: callForResult, Args: args}
	Publish(serverName, msg)

	ri := <-chanSyncRet
	return chanrpc.Assert(ri.Ret), ri.Err
}

func AsynCall(serverName string, chanAsynRet chan *chanrpc.RetInfo, id interface{}, args ...interface{}) {
	if len(args) < 1 {
		panic(fmt.Sprintf("%v asyn call of callback function not found", id))
	}

	lastIndex := len(args) - 1
	cb := args[lastIndex]
	args = args[:lastIndex]

	var callType uint8
	switch cb.(type) {
	case func(error):
		callType = callForResult
	case func(interface{}, error):
		callType = callForResult
	case func([]interface{}, error):
		callType = callForResult
	default:
		panic(fmt.Sprintf("%v asyn call definition of callback function is invalid", id))
	}

	request := &RequestInfo{cb: cb, chanRet: chanAsynRet}
	requestID := registerRequest(request)
	msg := &S2S_NsqMsg{RequestID: requestID, ReqType: NsqMsgTypeReq, MsgID: id, CallType: callType, Args: args}
	Publish(serverName, msg)
}
