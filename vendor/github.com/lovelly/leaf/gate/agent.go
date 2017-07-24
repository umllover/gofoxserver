package gate

import (
	"net"

	"github.com/lovelly/leaf/chanrpc"
	"github.com/lovelly/leaf/module"
)

type Agent interface {
	WriteMsg(msg interface{})
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	SetReason( int)
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	Skeleton() *module.Skeleton
	ChanRPC() *chanrpc.Server
}
