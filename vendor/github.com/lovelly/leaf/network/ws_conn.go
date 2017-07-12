package network

import (
	"errors"
	"net"
	"sync"
	"time"

	"fmt"

	"github.com/gorilla/websocket"
	"github.com/lovelly/leaf/log"
	"github.com/lovelly/leaf/util"
)

type WebsocketConnSet map[*websocket.Conn]struct{}

type WSConn struct {
	sync.Mutex
	conn      *websocket.Conn
	writeChan chan []byte
	maxMsgLen uint32
	closeFlag bool
}

func newWSConn(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *WSConn {
	wsConn := new(WSConn)
	wsConn.conn = conn
	wsConn.writeChan = make(chan []byte, pendingWriteNum)
	wsConn.maxMsgLen = maxMsgLen

	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			conn.Close()
			wsConn.Lock()
			wsConn.closeFlag = true
			wsConn.Unlock()
		}()

		for {
			select {
			case b, ok := <-wsConn.writeChan:
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if b == nil || !ok {
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				err := conn.WriteMessage(websocket.BinaryMessage, b)
				if err != nil {
					log.Error("write msg error :%s", err.Error())
					conn.WriteMessage(websocket.CloseMessage, []byte{})
					break
				}
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
					return
				}
			}
		}
	}()

	return wsConn
}

func (wsConn *WSConn) doDestroy() {
	wsConn.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	wsConn.conn.Close()

	if !wsConn.closeFlag {
		close(wsConn.writeChan)
		wsConn.closeFlag = true
	}
}

func (wsConn *WSConn) Destroy() {
	wsConn.Lock()
	defer wsConn.Unlock()

	wsConn.doDestroy()
}

func (wsConn *WSConn) Close() {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return
	}

	wsConn.doWrite(nil)
	wsConn.closeFlag = true
}

func (wsConn *WSConn) doWrite(b []byte) {
	if len(wsConn.writeChan) == cap(wsConn.writeChan) {
		log.Debug("close conn: channel full")
		wsConn.doDestroy()
		return
	}

	wsConn.writeChan <- b
}

func (wsConn *WSConn) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

func (wsConn *WSConn) RemoteAddr() net.Addr {
	return wsConn.conn.RemoteAddr()
}

// goroutine not safe
func (wsConn *WSConn) ReadMsg() ([]byte, error) {
	_, b, err := wsConn.conn.ReadMessage()
	if err != nil {
		return b, err
	}
	str, err := util.DesDecrypt(b, []byte("mqjx@mqc"))
	if err != nil {
		return b, errors.New(fmt.Sprintf("at ws ReadMsg msg DesDecrypt error :%s", err.Error()))
	}
	return str, err
}

// args must not be modified by the others goroutines
func (wsConn *WSConn) WriteMsg(args ...[]byte) error {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return nil
	}

	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	// check len
	if msgLen > wsConn.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < 1 {
		return errors.New("message too short")
	}

	// don't copy
	if len(args) == 1 {
		str, err := util.DesEncrypt(args[0], []byte("mqjx@mqc"))
		if err != nil {
			return errors.New(fmt.Sprintf("at ws write msg DesEncrypt error :%s", err.Error()))
		}

		wsConn.doWrite(str)
		return nil
	}

	// merge the args
	msg := make([]byte, msgLen)
	l := 0
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}
	str, err := util.DesEncrypt(msg, []byte("mqjx@mqc"))
	if err != nil {
		return errors.New(fmt.Sprintf("at ws write msg DesEncrypt error :%s", err.Error()))
	}

	wsConn.doWrite(str)

	return nil
}
