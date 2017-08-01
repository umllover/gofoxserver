package base

import (
	"mj/gameServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//recharge_limit
//

// +gen
type RechargeLimit struct {
	Level  int    `db:"level" json:"level"`   //
	Min    int    `db:"min" json:"min"`       //
	Max    int    `db:"max" json:"max"`       //
	Remark string `db:"remark" json:"remark"` //
}

var DefaultRechargeLimit = RechargeLimit{}

type rechargeLimitCache struct {
	objMap  map[int]*RechargeLimit
	objList []*RechargeLimit
}

var RechargeLimitCache = &rechargeLimitCache{}

func (c *rechargeLimitCache) LoadAll() {
	sql := "select * from recharge_limit"
	c.objList = make([]*RechargeLimit, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*RechargeLimit)
	log.Debug("Load all recharge_limit success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.Level] = v
	}
}

func (c *rechargeLimitCache) All() []*RechargeLimit {
	return c.objList
}

func (c *rechargeLimitCache) Count() int {
	return len(c.objList)
}

func (c *rechargeLimitCache) Get(level int) (*RechargeLimit, bool) {
	return c.GetKey1(level)
}

func (c *rechargeLimitCache) GetKey1(level int) (*RechargeLimit, bool) {
	v, ok := c.objMap[level]
	return v, ok
}
