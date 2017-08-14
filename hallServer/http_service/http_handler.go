package http_service

import (
	"github.com/lovelly/leaf/log"
)

var DefaultHttpHandler = NewDefaultHttpHandler()

func NewDefaultHttpHandler() *HttpHandler {
	return new(HttpHandler)
}

type HttpHandler struct {
	GMNoticeCallBack func(args []interface{}) // 跑马灯GM公告回调
}

func (h *HttpHandler) GMNotice(sendTimes int, interval int, context string) {
	log.Debug("服务端接收到GM消息%d,%d,%s", sendTimes, interval, context)
	h.GMNoticeCallBack([]interface{}{sendTimes, interval, context})
}

func (h *HttpHandler) SetGMNoticeCallBack(callBack func(args []interface{})) {
	log.Debug("设置GM回调")
	h.GMNoticeCallBack = callBack
}
