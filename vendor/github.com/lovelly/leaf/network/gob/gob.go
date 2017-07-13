package gob

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
)

type Processor struct {
	msgInfo map[string]*MsgInfo
}

type Buffer struct {
	*bytes.Buffer
}

type Encoder struct {
	buffer   *Buffer
	coder    *gob.Encoder
	encMutex sync.Mutex
}

func NewEncoder() *Encoder {
	buff := &Buffer{}
	coder := gob.NewEncoder(buff)
	return &Encoder{buffer: buff, coder: coder}
}

type Decoder struct {
	buffer   *Buffer
	coder    *gob.Decoder
	decMutex sync.Mutex
}

func NewDecoder() *Decoder {
	buff := &Buffer{}
	coder := gob.NewDecoder(buff)
	return &Decoder{buffer: buff, coder: coder}
}

type MsgInfo struct {
	msgType       reflect.Type
	msgRouter     *chanrpc.Server
	msgHandler    MsgHandler
	msgRawHandler MsgHandler
}

type MsgHandler func([]interface{})

type MsgRaw struct {
	msgID      string
	msgRawData []byte
}

func NewProcessor() *Processor {
	p := new(Processor)
	p.msgInfo = make(map[string]*MsgInfo)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) Register(msg interface{}) string {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("gob message pointer required")
	}
	msgID := msgType.Elem().Name()
	if msgID == "" {
		log.Fatal("unnamed gob message")
	}
	if _, ok := p.msgInfo[msgID]; ok {
		log.Fatal("message %v is already registered", msgID)
	}

	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo[msgID] = i
	return msgID
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRouter(msg interface{}, msgRouter *chanrpc.Server) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("gob message pointer required")
	}
	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRouter = msgRouter
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetHandler(msg interface{}, msgHandler MsgHandler) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("gob message pointer required")
	}
	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgHandler = msgHandler
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetRawHandler(msg interface{}, msgRawHandler MsgHandler) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		log.Fatal("message %v not registered", msgID)
	}

	i.msgRawHandler = msgRawHandler
}

// goroutine safe
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	// raw
	if msgRaw, ok := msg.(MsgRaw); ok {
		i, ok := p.msgInfo[msgRaw.msgID]
		if !ok {
			return fmt.Errorf("message %v not registered", msgRaw.msgID)
		}
		if i.msgRawHandler != nil {
			i.msgRawHandler([]interface{}{msgRaw.msgID, msgRaw.msgRawData, userData})
		}
		return nil
	}

	// json
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return errors.New("gob message pointer required")
	}
	msgID := msgType.Elem().Name()
	i, ok := p.msgInfo[msgID]
	if !ok {
		return fmt.Errorf("message %v not registered", msgID)
	}
	if i.msgHandler != nil {
		i.msgHandler([]interface{}{msg, userData})
	}
	if i.msgRouter != nil {
		i.msgRouter.Go(msgType, msg, userData)
	} else if i.msgHandler == nil {
		log.Error("%v msg without any handler", msgID)
	}
	return nil
}

// goroutine safe
func (p *Processor) Unmarshal(dec *Decoder, data []byte) (interface{}, error) {
	dec.decMutex.Lock()
	defer dec.decMutex.Unlock()
	var msgID string
	dec.buffer.Buffer = bytes.NewBuffer(data)
	err := dec.coder.Decode(&msgID)
	if err != nil {
		return nil, err
	}

	i, ok := p.msgInfo[msgID]
	if !ok {
		return nil, fmt.Errorf("message %v not registered", msgID)
	}

	// msg
	if i.msgRawHandler != nil {
		return MsgRaw{msgID, dec.buffer.Bytes()}, nil
	} else {
		msg := reflect.New(i.msgType.Elem()).Interface()
		return msg, dec.coder.Decode(msg)
	}

	panic("bug")
}

// goroutine safe
func (p *Processor) Marshal(enc *Encoder, msg interface{}) ([][]byte, error) {
	enc.encMutex.Lock()
	defer enc.encMutex.Unlock()
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		return nil, errors.New("json message pointer required")
	}
	msgID := msgType.Elem().Name()
	if _, ok := p.msgInfo[msgID]; !ok {
		return nil, fmt.Errorf("message %v not registered", msgID)
	}

	// data
	enc.buffer.Buffer = &bytes.Buffer{}
	err := enc.coder.Encode(&msgID)
	if err != nil {
		return [][]byte{enc.buffer.Bytes()}, err
	}

	err = enc.coder.Encode(msg)
	return [][]byte{enc.buffer.Bytes()}, err
}
