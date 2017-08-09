package base

import (
	"mj/gameServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//upgrade_advisor
//

// +gen
type UpgradeAdvisor struct {
	AdvisorId int    `db:"advisor_id" json:"advisor_id"` //
	Question  string `db:"question" json:"question"`     //
	Answer    string `db:"answer" json:"answer"`         //
	Order     int    `db:"order" json:"order"`           //
}

var DefaultUpgradeAdvisor = UpgradeAdvisor{}

type upgradeAdvisorCache struct {
	objMap  map[int]*UpgradeAdvisor
	objList []*UpgradeAdvisor
}

var UpgradeAdvisorCache = &upgradeAdvisorCache{}

func (c *upgradeAdvisorCache) LoadAll() {
	sql := "select * from upgrade_advisor"
	c.objList = make([]*UpgradeAdvisor, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*UpgradeAdvisor)
	log.Debug("Load all upgrade_advisor success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.AdvisorId] = v
	}
}

func (c *upgradeAdvisorCache) All() []*UpgradeAdvisor {
	return c.objList
}

func (c *upgradeAdvisorCache) Count() int {
	return len(c.objList)
}

func (c *upgradeAdvisorCache) Get(advisor_id int) (*UpgradeAdvisor, bool) {
	return c.GetKey1(advisor_id)
}

func (c *upgradeAdvisorCache) GetKey1(advisor_id int) (*UpgradeAdvisor, bool) {
	v, ok := c.objMap[advisor_id]
	return v, ok
}
