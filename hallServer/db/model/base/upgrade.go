package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//upgrade
//

// +gen
type Upgrade struct {
	LevelId    int     `db:"level_id" json:"level_id"`     //
	Recharge   int     `db:"recharge" json:"recharge"`     //
	Commission int     `db:"commission" json:"commission"` //
	AgentNum   int     `db:"agent_num" json:"agent_num"`   //
	RateB      float64 `db:"rate_b" json:"rate_b"`         //
	RateC      float64 `db:"rate_c" json:"rate_c"`         //
	RateD      float64 `db:"rate_d" json:"rate_d"`         //
	IconId     int     `db:"icon_id" json:"icon_id"`       //
	Liveness   int     `db:"liveness" json:"liveness"`     //
}

var DefaultUpgrade = Upgrade{}

type upgradeCache struct {
	objMap  map[int]*Upgrade
	objList []*Upgrade
}

var UpgradeCache = &upgradeCache{}

func (c *upgradeCache) LoadAll() {
	sql := "select * from upgrade"
	c.objList = make([]*Upgrade, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*Upgrade)
	log.Debug("Load all upgrade success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.LevelId] = v
	}
}

func (c *upgradeCache) All() []*Upgrade {
	return c.objList
}

func (c *upgradeCache) Count() int {
	return len(c.objList)
}

func (c *upgradeCache) Get(level_id int) (*Upgrade, bool) {
	return c.GetKey1(level_id)
}

func (c *upgradeCache) GetKey1(level_id int) (*Upgrade, bool) {
	v, ok := c.objMap[level_id]
	return v, ok
}
