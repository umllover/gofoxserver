package internal

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

func init() {
	skeleton.RegisterCommand("tickAccount", "Usage: tickAccount|accountId", tickAccount)
}

func tickAccount(args []interface{}) (ret interface{}, err error) {
	ret = ""
	if len(args) < 1 {
		err = errors.New("args len is less than 1")
		return
	}

	accountId := bson.ObjectIdHex(args[0].(string))
	ChanRPC.Go("KickAccount", accountId)
	return
}
