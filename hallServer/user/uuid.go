package user

import (
	"mj/hallServer/conf"
	"time"

	"github.com/lovelly/leaf/log"
)

var key int64

func GetKey() int64 {
	key = key + 1
	if key > 0xFFF {
		key = 0
	}
	return key
}

type Uuid struct {
	uid int64
}

func (u *Uuid) GetUUid() int64 {
	return u.uid
}

func (u *Uuid) SetTimestamp(ti int64) {
	if ti < 0 {
		log.Error("SetTimestamp ti < 0 ")
	}
	ti = ti << 22
	u.uid |= ti
}

func (u *Uuid) SetMachineKey(key int64) {
	key = key << 12
	key = key & 0X3FF000
	u.uid |= key
}

func (u *Uuid) SetSerial(s int64) {
	s = s & 0xFFF
	u.uid |= s
}

func NewUUid() *Uuid {
	return new(Uuid)
}

func GetUUID() int64 {
	time.Sleep(1 * time.Millisecond)
	timeline := time.Now().UnixNano() / 1e6
	uuid := NewUUid()
	uuid.SetTimestamp(timeline)
	uuid.SetMachineKey(int64(conf.Server.NodeId))
	uuid.SetSerial(GetKey())
	return uuid.GetUUid()
}
