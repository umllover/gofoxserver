package cluster

import (
	"encoding/gob"
	"errors"
	"fmt"
	"mj/gameServer/conf"
	"sync"

	"reflect"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/network/json"
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
	Processor       = json.NewProcessor()
	RequestInfoLock sync.Mutex
	requestID       int64
	requestMap      = make(map[int64]*RequestInfo)
)

type RequestInfo struct {
	serverName string
	cb         interface{}
	chanRet    chan *chanrpc.RetInfo
}

func init() {
	gob.Register(map[string]string{})
	gob.Register(map[string]interface{}{})
	gob.Register(&S2S_NsqMsg{})

	Processor.Register(&S2S_NsqMsg{})
	Processor.Register(&chanrpc.RetInfo{})
}

func SetRouter(msgID interface{}, server *chanrpc.Server) {
	_, ok := routeMap[msgID]
	if ok {
		panic(fmt.Sprintf("function id %v: already set route", msgID))
	}

	routeMap[msgID] = server.Open(0)
}

type S2S_NsqMsg struct {
	RequestID     int64
	ReqType       int
	CallType      uint8
	SrcServerName string
	DstServerName string
	Args          []byte
	Err           string
}

func handleRequestMsg(recvMsg *S2S_NsqMsg) {
	sendMsg := &S2S_NsqMsg{ReqType: NsqMsgTypeRsp, DstServerName: recvMsg.SrcServerName, RequestID: recvMsg.RequestID}
	if isClose() && recvMsg.CallType == callForResult {
		sendMsg.Err = fmt.Sprintf("%v server is closing", conf.ServerName)
		Publish(sendMsg)
		return
	}

	msg, err := Processor.Unmarshal(recvMsg.Args)
	if err != nil && recvMsg.CallType == callForResult {
		sendMsg.Err = fmt.Sprintf("%v Unmarshal msg error:%s", conf.ServerName, err.Error())
		Publish(sendMsg)
		return
	}

	msgType := reflect.TypeOf(msg)
	if (msgType == nil || msgType.Kind() != reflect.Ptr) && recvMsg.CallType == callForResult {
		sendMsg.Err = fmt.Sprintf("json message pointer required")
		Publish(sendMsg)
		return
	}

	msgID := msgType.Elem().Name()
	client, ok := routeMap[msgID]
	if !ok {
		err := fmt.Sprintf("%v msg is not set route", msgID)
		log.Error(err)

		if recvMsg.CallType == callForResult {
			sendMsg.Err = err
			Publish(sendMsg)
		}
		return
	}

	args := []interface{}{msg}
	if recvMsg.CallType == callForResult {
		sendMsgFunc := func(ret *chanrpc.RetInfo) {
			data, err := Processor.Marshal(ret.Ret)
			if err == nil {
				sendMsg.Args = data[0]
			} else {
				log.Error("at handleRequestMsg  Processor.Marshal ret error:%s", err.Error())
				sendMsg.Err = err.Error()
			}

			if ret.Err != nil {
				sendMsg.Err = ret.Err.Error()
			}
			Publish(sendMsg)
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
		log.Error("%v: request id %v is not exist", msg.SrcServerName, msg.RequestID)
		return
	}

	ret := &chanrpc.RetInfo{Cb: request.cb}
	retMsg, err := Processor.Unmarshal(msg.Args)
	if err != nil {
		log.Error("handleResponseMsg Unmarshal msg error:%s", err.Error())
		ret.Err = fmt.Errorf("handleResponseMsg Unmarshal msg error:%s", err.Error())
		return
	}
	ret.Ret = retMsg
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

func ForEachRequest(f func(id int64, request *RequestInfo)) {
	RequestInfoLock.Lock()
	defer RequestInfoLock.Unlock()
	for id, v := range requestMap {
		f(id, v)
	}
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
