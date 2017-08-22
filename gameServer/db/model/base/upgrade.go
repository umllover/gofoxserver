package base

import (
	"mj/gameServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//upgrade
//

// +gen
type Upgrade struct {
	UpId           int     `db:"up_id" json:"up_id"`                     //
	Level          int     `db:"level" json:"level"`                     //
	InstructorFees int     `db:"instructor_fees" json:"instructor_fees"` //
	InstructorRate int     `db:"instructor_rate" json:"instructor_rate"` //
	TotalIncome    string  `db:"total_income" json:"total_income"`       //
	AgentNum       int     `db:"agent_num" json:"agent_num"`             //
	RateB          float64 `db:"rate_b" json:"rate_b"`                   //
	RateC          float64 `db:"rate_c" json:"rate_c"`                   //
	RateD          float64 `db:"rate_d" json:"rate_d"`                   //
	IconId         int     `db:"icon_id" json:"icon_id"`                 //
	Liveness       int     `db:"liveness" json:"liveness"`               //
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
		c.objMap[v.UpId] = v
	}
}

func (c *upgradeCache) All() []*Upgrade {
	return c.objList
}

func (c *upgradeCache) Count() int {
	return len(c.objList)
}

func (c *upgradeCache) Get(up_id int) (*Upgrade, bool) {
	return c.GetKey1(up_id)
}

func (c *upgradeCache) GetKey1(up_id int) (*Upgrade, bool) {
	v, ok := c.objMap[up_id]
	return v, ok
}
