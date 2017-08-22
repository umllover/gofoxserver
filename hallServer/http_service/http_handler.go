package http_service

import (
	"mj/hallServer/race_msg"

	"github.com/lovelly/leaf/log"
)

var DefaultHttpHandler = NewDefaultHttpHandler()

func NewDefaultHttpHandler() *HttpHandler {
	return new(HttpHandler)
}

type HttpHandler struct {
}

func (h *HttpHandler) GMNotice(sendTimes int, interval int, context string) {
	log.Debug("服务端接收到GM消息%d,%d,%s", sendTimes, interval, context)
	race_msg.GmNotify(sendTimes, interval, context)
}
