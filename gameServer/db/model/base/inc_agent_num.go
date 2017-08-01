package base

import (
	"mj/gameServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//inc_agent_num
//

// +gen
type IncAgentNum struct {
	NodeId int   `db:"node_id" json:"node_id"` //
	Total  int64 `db:"total" json:"total"`     //
}

var DefaultIncAgentNum = IncAgentNum{}

type incAgentNumCache struct {
	objMap  map[int]*IncAgentNum
	objList []*IncAgentNum
}

var IncAgentNumCache = &incAgentNumCache{}

func (c *incAgentNumCache) LoadAll() {
	sql := "select * from inc_agent_num"
	c.objList = make([]*IncAgentNum, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*IncAgentNum)
	log.Debug("Load all inc_agent_num success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.NodeId] = v
	}
}

func (c *incAgentNumCache) All() []*IncAgentNum {
	return c.objList
}

func (c *incAgentNumCache) Count() int {
	return len(c.objList)
}

func (c *incAgentNumCache) Get(node_id int) (*IncAgentNum, bool) {
	return c.GetKey1(node_id)
}

func (c *incAgentNumCache) GetKey1(node_id int) (*IncAgentNum, bool) {
	v, ok := c.objMap[node_id]
	return v, ok
}
