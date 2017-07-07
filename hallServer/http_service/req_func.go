package http_service

import (
	"mj/common/msg"

	"github.com/lovelly/leaf/log"
)

func AuthUser(url string, data interface{}) (info *msg.AuthInfo) {
	info = &msg.AuthInfo{}
	Result := PostJSON(url, data)
	jsonMap, err := Result.Map()
	if err != nil {
		info.RetCode = 99
		return
	}

	log.Debug("get AuthUser info %v", jsonMap)
	return
}
