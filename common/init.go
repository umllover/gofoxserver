package common

import (
	"encoding/gob"
	"mj/common/msg"

	"gopkg.in/mgo.v2/bson"
)

func Init() {
	gob.Register(bson.NewObjectId())
	gob.Register([]bson.ObjectId{})
	gob.Register(map[string]string{})
	gob.Register([]*msg.TagGameServer{})
	gob.Register(&msg.RoomInfo{})
	gob.Register([]*msg.RoomInfo{})
}
