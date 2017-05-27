package user

import (
	"sync"
	"net"
	"github.com/lovelly/leaf/gate"
)

type WAgent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	UserData() interface{}
}

type User struct {
	WAgent
	Id int
	RoomId int  //当前在哪个房间
	sync.RWMutex
}

func NewUser(UserId int,a gate.Agent ) *User {
	return &User{Id : UserId, WAgent:a}
}
