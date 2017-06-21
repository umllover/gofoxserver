package base

import (
	"mj/hallServer/db"

	"github.com/lovelly/leaf/log"
)

//This file is generate by scripts,don't edit it

//refresh_in_time
//

// +gen
type RefreshInTime struct {
	Id               int    `db:"id" json:"id"`                                 //
	NodeId           int    `db:"node_id" json:"node_id"`                       //
	RefreshTableList string `db:"refresh_table_list" json:"refresh_table_list"` //
	Cnt              int    `db:"cnt" json:"cnt"`                               //
}

var DefaultRefreshInTime = RefreshInTime{}

type refreshInTimeCache struct {
	objMap  map[int]*RefreshInTime
	objList []*RefreshInTime
}

var RefreshInTimeCache = &refreshInTimeCache{}

func (c *refreshInTimeCache) LoadAll() {
	sql := "select * from refresh_in_time"
	c.objList = make([]*RefreshInTime, 0)
	err := db.BaseDB.Select(&c.objList, sql)
	if err != nil {
		log.Fatal(err.Error())
	}
	c.objMap = make(map[int]*RefreshInTime)
	log.Debug("Load all refresh_in_time success %v", len(c.objList))
	for _, v := range c.objList {
		c.objMap[v.Id] = v
	}
}

func (c *refreshInTimeCache) All() []*RefreshInTime {
	return c.objList
}

func (c *refreshInTimeCache) Count() int {
	return len(c.objList)
}

func (c *refreshInTimeCache) Get(id int) (*RefreshInTime, bool) {
	return c.GetKey1(id)
}

func (c *refreshInTimeCache) GetKey1(id int) (*RefreshInTime, bool) {
	v, ok := c.objMap[id]
	return v, ok
}
