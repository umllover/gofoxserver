package internal

import (
	"errors"
	"mj/common/msg"
	"reflect"

	"github.com/lovelly/leaf/chanrpc"
)

/////主消息函数
func (m *UserModule) handleMsgData(args []interface{}) error {
	if msg.Processor != nil {
		str := args[0].([]byte)
		data, err := msg.Processor.Unmarshal(str)
		if err != nil {
			return err
		}

		msgType := reflect.TypeOf(data)
		if msgType == nil || msgType.Kind() != reflect.Ptr {
			return errors.New("json message pointer required 11")
		}

		if m.a == nil && msgType.Elem().Name() != "C2G_GR_LogonMobile" {
			return errors.New("user not login")
		}

		f, ok := m.ChanRPC.HasFunc(msgType)
		if ok {
			m.ChanRPC.Exec(chanrpc.BuildGoCallInfo(f, data, m.a))
			return nil
		}

		err = msg.Processor.RouteByType(msgType, data, m.a)
		if err != nil {
			return err
		}
	}
	return nil
}
