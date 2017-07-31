package internal

import (
	"errors"
	"mj/common/msg"
	"reflect"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/log"
)

/////主消息函数
func (m *UserModule) handleMsgData(args []interface{}) error {
	if msg.Processor != nil {
		str := args[0].([]byte)
		data, err := msg.Processor.Unmarshal(str)
		if err != nil {
			log.Debug("111111111111111111")
			return err
		}

		msgType := reflect.TypeOf(data)
		if msgType == nil || msgType.Kind() != reflect.Ptr {
			return errors.New("json message pointer required 11")
		}

		f, ok := m.ChanRPC.HasFunc(msgType)
		if ok {
			m.ChanRPC.Exec(chanrpc.BuildGoCallInfo(f, data, m.a))
			return nil
		} else {
			log.Debug("2222222222222")
		}

		err = msg.Processor.RouteByType(msgType, data, m.a)
		if err != nil {
			log.Debug("33333")
			return err
		}
	}
	return nil
}
