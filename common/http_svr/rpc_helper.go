package http_svr

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"runtime/debug"

	"unicode"
	"unicode/utf8"

	"errors"

	"github.com/lovelly/leaf/log"
)

const (
	FROM_SVR_MSG    = 1
	FROM_CLIENT_MSG = 2
	FROM_CLOSE_MSG  = 4
)

type RpcRequest struct {
	Method *string        `json:"func_name"`
	Params *[]interface{} `json:"params"`
}

type RpcResponse struct {
	Method string      `json:"func_name"`
	Result interface{} `json:"data"`
	Error  string      `json:"error"`
}

type RpcMethod struct {
	Method reflect.Method
	host   reflect.Value
	idx    int
}

type RpcHelper struct {
	Methods map[string]*RpcMethod
}

type HandelData struct {
	MsgType int
	Data    interface{}
}

// JSON standard : all number are Number type, that is float64 in golang.
func convertParam(v interface{}, target_type reflect.Type) (new_v reflect.Value, ok bool) {
	defer func() {
		if re := recover(); re != nil {
			ok = false
			log.Error("convertParam %v", re)
		}
	}()

	ok = true
	if target_type.Kind() == reflect.Interface {
		new_v = reflect.ValueOf(v)
	} else if reflect.TypeOf(v).Kind() == reflect.Float64 {
		f := v.(float64)
		switch target_type.Kind() {
		case reflect.Int:
			new_v = reflect.ValueOf(int(f))
		case reflect.Uint8:
			new_v = reflect.ValueOf(uint8(f))
		case reflect.Uint16:
			new_v = reflect.ValueOf(uint16(f))
		case reflect.Uint32:
			new_v = reflect.ValueOf(uint32(f))
		case reflect.Uint64:
			new_v = reflect.ValueOf(uint64(f))
		case reflect.Int8:
			new_v = reflect.ValueOf(int8(f))
		case reflect.Int16:
			new_v = reflect.ValueOf(int16(f))
		case reflect.Int32:
			new_v = reflect.ValueOf(int32(f))
		case reflect.Int64:
			new_v = reflect.ValueOf(int64(f))
		case reflect.Float32:
			new_v = reflect.ValueOf(float32(f))
		default:
			ok = false
		}
	} else if reflect.TypeOf(v).Kind() == target_type.Kind() {
		new_v = reflect.ValueOf(v)
	} else if target_type.Kind() == reflect.Ptr { //if it is pointer, get it element type
		new_v = reflect.ValueOf(&v) //target_type.Elem()
	} else {
		ok = false
	}

	return
}

var DefaultHttpRpcHelpr = NewRpcHelper()

func NewRpcHelper() *RpcHelper {
	return &RpcHelper{make(map[string]*RpcMethod)}
}

// 解析客户端请求，
// default_params : 服务器端调用时自动带入的参数, 和客户端请求的参数共同组成method的参数。
func (h *RpcHelper) Parse(msg *HandelData, default_params ...interface{}) (method *RpcMethod, params []reflect.Value, err error) {
	req := &RpcRequest{}

	json_data, ok := (msg.Data).([]byte)
	if !ok {
		err = errors.New("not foud data at parse")
		return
	}
	if json_err := json.Unmarshal(json_data, req); json_err != nil {
		err = errors.New("Unmarshal error at parse")
		return
	}

	if req.Method == nil || *req.Method == "" {
		err = errors.New(" not foud module ")
		return
	}

	method, ok = h.Methods[*req.Method]
	if !ok {
		err = errors.New(" not foud module ")
		return
	}

	default_params_len := len(default_params)
	//长度应减去method的receiver

	lens := len(*req.Params)
	if lens != (method.Method.Type.NumIn() - default_params_len - 1) {
		err = errors.New(fmt.Sprintf("params not matched. found %d, need %d.", lens, method.Method.Type.NumIn()-default_params_len-1))
		return
	}

	params = make([]reflect.Value, lens+default_params_len)
	//第一个参数是*net.Client
	for idx, hdn_param := range default_params {
		params[idx] = reflect.ValueOf(hdn_param)
	}

	for i := 0; i < lens; i++ {
		target_type := method.Method.Type.In(i + 1 + default_params_len) //跳过receiver和default_params
		new_param, ok := convertParam((*req.Params)[i], target_type)
		if !ok {
			err = errors.New(fmt.Sprintf("convert param faild. expect %s, found=%v value=%v.", target_type, reflect.TypeOf(((*req.Params)[i])), (*req.Params)[i]))
			return
		}
		params[i+default_params_len] = new_param
	}
	return
}

func (h *RpcHelper) Call(method *RpcMethod, params []reflect.Value) (resp *RpcResponse, game_err error) {

	resp = &RpcResponse{}

	resp.Method = method.Method.Name

	defer func() {
		if re := recover(); re != nil {
			switch re.(type) {
			case runtime.Error:
				log.Error("runtime error %v\n %v", re.(error).Error(), string(debug.Stack()))
				game_err = errors.New(re.(error).Error())
			case string:
				log.Error("runtime error %v\n %v", re.(string), string(debug.Stack()))
				game_err = errors.New(re.(string))
			case error:
				log.Error("runtime error %v\n %v", re.(error).Error(), string(debug.Stack()))
				game_err = re.(error)
			default:
				debug.PrintStack()
			}
			//r.Error = err
		}
	}()

	result := method.host.Method(method.idx).Call(params)
	if len(result) > 0 {
		resp.Result = result[0].Interface()
	}
	return
}

func (h *RpcHelper) RegisterMethod(v interface{}) {
	reflectType := reflect.TypeOf(v)
	host := reflect.ValueOf(v)
	for i := 0; i < reflectType.NumMethod(); i++ {
		m := reflectType.Method(i)
		char, _ := utf8.DecodeRuneInString(m.Name)
		if !unicode.IsUpper(char) {
			continue
		}
		h.Methods[m.Name] = &RpcMethod{m, host, m.Index}
	}
}
