package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//agent_limit
//

// +gen
type AgentLimit struct {
	Level  int    `db:"level" json:"level"`   //
	Min    int    `db:"min" json:"min"`       //
	Max    int    `db:"max" json:"max"`       //
	Remark string `db:"remark" json:"remark"` //
}

var DefaultAgentLimit = AgentLimit{}

type agentLimitCache struct {
	objMap  map[int]*AgentLimit
	objList []*AgentLimit
}

var AgentLimitCache = &agentLimitCache{}

func (c *agentLimitCache) LoadAll() {
	sql := "select * from agent_limit"
	c.objList = make([]*AgentLimit, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*AgentLimit)
	log.Debug("Load all agent_limit success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.Level] = v
	}
}

func (c *agentLimitCache) All() []*AgentLimit {
	return c.objList
}

func (c *agentLimitCache) Count() int {
	return len(c.objList)
}

func (c *agentLimitCache) Get(level int) (*AgentLimit, bool) {
	return c.GetKey1(level)
}

func (c *agentLimitCache) GetKey1(level int) (*AgentLimit, bool) {
	v, ok := c.objMap[level]
	return v, ok
}
