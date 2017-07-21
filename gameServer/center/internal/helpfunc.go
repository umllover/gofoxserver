package internal

import (
	"github.com/lovelly/leaf/log"
)

func SendMsgToSelfNotdeUser(args []interface{}) {
	uid := args[0].(int64)
	FuncName := args[1].(string)
	ch, ok := Users[uid]
	if ok {
		ch.Go(FuncName, args[2:]...)
		return
	} else {

	}
	log.Debug("at SendMsgToSelfNotdeUser player not in node")
	return
}
