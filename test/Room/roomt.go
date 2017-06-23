package Room

import (
	"fmt"
	"sync"
)

var Wg sync.WaitGroup

func NewRoom() *Room {
	r := new(Room)
	r.ChanRpc = make(chan []interface{}, 1)
	r.Funcs = make(map[string]func(interface{}))
	return r
}

type Room struct {
	ChanRpc chan []interface{}
	Funcs   map[string]func(interface{})
}

func (r *Room) Run() {
	Wg.Add(1)
	go func() {
		for v := range r.ChanRpc {
			r.Handler(v)
		}
		Wg.Done()
	}()
}

func (r *Room) RegistFuncs() {
	r.Regist("Test", r.TestFunc)
	r.Regist("SitDown", r.SitDownFunc)
}

func (r *Room) Go(args ...interface{}) {
	r.ChanRpc <- args
}

func (r *Room) Regist(funcName string, f func(interface{})) {
	r.Funcs[funcName] = f
}

func (r *Room) Handler(args []interface{}) {
	funcName := args[0].(string)
	f, ok := r.Funcs[funcName]
	if ok {
		f(args[1])
	}
}

func (r *Room) TestFunc(param interface{}) {
	fmt.Println("at test func ", param)
}

func (r *Room) SitDownFunc(param interface{}) {
	fmt.Println("at test func ", param)
}
