package http_service

var DefaultHttpHandler = NewDefaultHttpHandler()

func NewDefaultHttpHandler() *HttpHandler {
	return new(HttpHandler)
}

type HttpHandler struct {
}

func (h *HttpHandler) GMNotice(beginTime, endTime, interval, context string) {

}
