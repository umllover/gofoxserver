package common

import (
	"encoding/gob"
	"gopkg.in/mgo.v2/bson"
	"mj/common/msg"
)

func Init() {
	gob.Register(bson.NewObjectId())
	gob.Register([]bson.ObjectId{})
	gob.Register(map[string]string{})
	gob.Register([]*msg.TagGameServer{})
	gob.Register(&msg.RoomInfo{})
}