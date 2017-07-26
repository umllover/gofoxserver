package user

import (
	"mj/hallServer/conf"

	"mj/hallServer/db/model"

	"github.com/lovelly/leaf/log"
)

var key int64
var MaxInc = int64(1<<42 - 1)

func LoadIncId() {
	inc, err := model.IncUseridOp.Get(conf.Server.NodeId)
	if !err {
		model.IncUseridOp.Insert(&model.IncUserid{NodeId: conf.Server.NodeId, IncId: 0})
	}

	inc, err = model.IncUseridOp.Get(conf.Server.NodeId)
	if !err {
		log.Fatal("LoadIncId faild ")
	}

	key = inc.IncId
}

func GetKey() int64 {
	key = key + 1
	if key > MaxInc {
		key = 0
	}
	model.IncUseridOp.UpdateWithMap(conf.Server.NodeId, map[string]interface{}{
		"inc_id": key,
	})
	return key
}

type Uuid struct {
	uid int64
}

func (u *Uuid) GetUUid() int64 {
	return u.uid
}

func (u *Uuid) SetNodeId(ti int64) {
	if ti < 0 {
		log.Error("SetTimestamp ti < 0 ")
	}
	ti = ti << 43
	u.uid |= ti
}

func (u *Uuid) SetSerial(s int64) {
	u.uid |= s
}

func NewUUid() *Uuid {
	return new(Uuid)
}

func GetUUID() int64 {
	uuid := NewUUid()
	uuid.SetNodeId(int64(conf.Server.NodeId))
	uuid.SetSerial(GetKey())
	return uuid.GetUUid()
}
